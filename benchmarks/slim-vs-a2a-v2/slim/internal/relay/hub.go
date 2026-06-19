// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package relay

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	slim "github.com/agntcy/slim-bindings-go"
	"github.com/agntcy/csit/benchmarks/slim-vs-a2a-v2/internal/consensus"
	"github.com/agntcy/csit/benchmarks/slim-vs-a2a-v2/internal/scenario"
	"github.com/agntcy/csit/benchmarks/slim-vs-a2a-v2/slim/internal/protocol"
	"github.com/agntcy/csit/benchmarks/slim-vs-a2a-v2/slim/internal/slimrpc"
)

const defaultSharedSecret = "demo-shared-secret-min-32-chars!!"

// Hub owns the shared SLIM group channel and multicasts findings to all agents.
type Hub struct {
	mu             sync.Mutex
	group          *slim.Channel
	app            *slim.App
	connID         uint64
	channels       map[string]*slim.Channel
	streamRPCCount int
	coordFanoutMS  int64
}

func NewHub(endpoint string, s *scenario.ConsensusScenario) (*Hub, *slim.App, error) {
	slim.InitializeWithDefaults()
	service := slim.GetGlobalService()

	runnerName := slim.NewName("agntcy", "bench-v2", "runner")
	app, err := service.CreateAppWithSecret(runnerName, defaultSharedSecret)
	if err != nil {
		return nil, nil, err
	}

	connID, err := service.Connect(slim.NewInsecureClientConfig(endpoint))
	if err != nil {
		app.Destroy()
		return nil, nil, err
	}
	if err := app.Subscribe(runnerName, &connID); err != nil {
		app.Destroy()
		return nil, nil, err
	}

	members := make([]*slim.Name, 0, len(s.Agents))
	for _, a := range s.Agents {
		n, err := nameFromString(a.SlimName)
		if err != nil {
			app.Destroy()
			return nil, nil, err
		}
		members = append(members, n)
		if err := app.SetRoute(n, connID); err != nil {
			app.Destroy()
			return nil, nil, fmt.Errorf("set route %s: %w", a.SlimName, err)
		}
	}

	group, err := slim.ChannelNewGroupWithConnection(app, members, &connID)
	if err != nil {
		app.Destroy()
		return nil, nil, err
	}

	hub := &Hub{
		group:    group,
		app:      app,
		connID:   connID,
		channels: map[string]*slim.Channel{},
	}
	for _, agent := range s.Agents {
		n, _ := nameFromString(agent.SlimName)
		hub.channels[agent.ID] = slim.ChannelNewWithConnection(app, n, &connID)
	}
	server := slim.ServerNewWithConnection(app, runnerName, &connID)
	server.RegisterUnaryUnary(slimrpc.ServiceName, slimrpc.MethodPublish, &publishHandler{hub: hub})

	go func() {
		if err := server.Serve(); err != nil {
			fmt.Printf("runner relay serve: %v\n", err)
		}
	}()

	return hub, app, nil
}

func (h *Hub) Close() {
	if h.group != nil {
		_ = h.group.Close(nil)
	}
	for _, ch := range h.channels {
		_ = ch.Close(nil)
	}
	if h.app != nil {
		h.app.Destroy()
	}
}

func (h *Hub) StartAll(agentIDs []string) error {
	for _, id := range agentIDs {
		if err := h.callControl(id, protocol.Request{Op: protocol.OpStart}); err != nil {
			return err
		}
	}
	return nil
}

func (h *Hub) Snapshot(agentID string) (consensus.AgentSnapshot, error) {
	resp, err := h.callControlWithResponse(agentID, protocol.Request{Op: protocol.OpSnapshot})
	if err != nil {
		return consensus.AgentSnapshot{}, err
	}
	var snap consensus.AgentSnapshot
	if err := json.Unmarshal([]byte(resp.Body), &snap); err != nil {
		return consensus.AgentSnapshot{}, err
	}
	return snap, nil
}

func (h *Hub) callControl(agentID string, req protocol.Request) error {
	_, err := h.callControlWithResponse(agentID, req)
	return err
}

func (h *Hub) callControlWithResponse(agentID string, req protocol.Request) (protocol.Response, error) {
	ch, ok := h.channels[agentID]
	if !ok {
		return protocol.Response{}, fmt.Errorf("unknown agent %q", agentID)
	}
	payload, err := protocol.EncodeRequest(req)
	if err != nil {
		return protocol.Response{}, err
	}
	timeout := 30 * time.Second
	reply, err := ch.CallUnary(slimrpc.ServiceName, slimrpc.MethodHandle, payload, &timeout, nil)
	if err != nil {
		return protocol.Response{}, err
	}
	return protocol.DecodeResponse(reply)
}

func (h *Hub) Stats() (streamRPCCount int, coordFanoutMS int64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.streamRPCCount, h.coordFanoutMS
}

func (h *Hub) MulticastFinding(payload []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	start := time.Now()
	timeout := 30 * time.Second
	reader, err := h.group.CallMulticastUnary(slimrpc.ServiceName, slimrpc.MethodFindingStream, payload, &timeout, nil)
	if err != nil {
		return err
	}
	defer reader.Destroy()

	for {
		switch msg := reader.Next().(type) {
		case slim.MulticastStreamMessageEnd:
			h.streamRPCCount++
			h.coordFanoutMS += time.Since(start).Milliseconds()
			return nil
		case slim.MulticastStreamMessageError:
			if msg.Error != nil {
				return msg.Error.AsError()
			}
			return fmt.Errorf("multicast error")
		case slim.MulticastStreamMessageData:
			// drain peer acks
		}
	}
}

type publishHandler struct {
	hub *Hub
}

func (h *publishHandler) Handle(request []byte, _ *slim.Context) ([]byte, error) {
	if _, err := consensus.DecodeFinding(request); err != nil {
		return nil, err
	}
	if err := h.hub.MulticastFinding(request); err != nil {
		return []byte(`{"ok":false}`), err
	}
	return []byte(`{"ok":true}`), nil
}

func nameFromString(value string) (*slim.Name, error) {
	parts := strings.Split(value, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid name format: %s", value)
	}
	return slim.NewName(parts[0], parts[1], parts[2]), nil
}
