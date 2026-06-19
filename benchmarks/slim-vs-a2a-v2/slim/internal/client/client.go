// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	slim "github.com/agntcy/slim-bindings-go"
	"github.com/agntcy/csit/benchmarks/slim-vs-a2a-v2/internal/consensus"
	"github.com/agntcy/csit/benchmarks/slim-vs-a2a-v2/internal/scenario"
	"github.com/agntcy/csit/benchmarks/slim-vs-a2a-v2/slim/internal/protocol"
	"github.com/agntcy/csit/benchmarks/slim-vs-a2a-v2/slim/internal/slimrpc"
)

const defaultSharedSecret = "demo-shared-secret-min-32-chars!!"

type Client struct {
	app        *slim.App
	connID     uint64
	channels   map[string]*slim.Channel
	agentNames map[string]*slim.Name
}

func New(endpoint string, s *scenario.ConsensusScenario) (*Client, error) {
	slim.InitializeWithDefaults()
	service := slim.GetGlobalService()

	runnerName := slim.NewName("agntcy", "bench-v2", "runner")
	app, err := service.CreateAppWithSecret(runnerName, defaultSharedSecret)
	if err != nil {
		return nil, err
	}

	connID, err := service.Connect(slim.NewInsecureClientConfig(endpoint))
	if err != nil {
		app.Destroy()
		return nil, err
	}
	if err := app.Subscribe(runnerName, &connID); err != nil {
		app.Destroy()
		return nil, err
	}

	agentNames := map[string]*slim.Name{}
	channels := map[string]*slim.Channel{}
	for _, agent := range s.Agents {
		name, err := nameFromString(agent.SlimName)
		if err != nil {
			app.Destroy()
			return nil, err
		}
		agentNames[agent.ID] = name
		if err := app.SetRoute(name, connID); err != nil {
			app.Destroy()
			return nil, fmt.Errorf("set route %s: %w", agent.SlimName, err)
		}
		channels[agent.ID] = slim.ChannelNewWithConnection(app, name, &connID)
	}

	return &Client{app: app, connID: connID, channels: channels, agentNames: agentNames}, nil
}

func (c *Client) Close() {
	for _, ch := range c.channels {
		_ = ch.Close(nil)
	}
	if c.app != nil {
		c.app.Destroy()
	}
}

func (c *Client) StartAll(ctx context.Context, agentIDs []string) error {
	for _, id := range agentIDs {
		if err := c.call(ctx, id, protocol.Request{Op: protocol.OpStart}); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) Snapshot(ctx context.Context, agentID string) (consensus.AgentSnapshot, error) {
	resp, err := c.callWithResponse(ctx, agentID, protocol.Request{Op: protocol.OpSnapshot})
	if err != nil {
		return consensus.AgentSnapshot{}, err
	}
	var snap consensus.AgentSnapshot
	if err := json.Unmarshal([]byte(resp.Body), &snap); err != nil {
		return consensus.AgentSnapshot{}, err
	}
	return snap, nil
}

func (c *Client) call(ctx context.Context, agentID string, req protocol.Request) error {
	_, err := c.callWithResponse(ctx, agentID, req)
	return err
}

func (c *Client) callWithResponse(ctx context.Context, agentID string, req protocol.Request) (protocol.Response, error) {
	ch, ok := c.channels[agentID]
	if !ok {
		return protocol.Response{}, fmt.Errorf("unknown agent %q", agentID)
	}
	payload, err := protocol.EncodeRequest(req)
	if err != nil {
		return protocol.Response{}, err
	}
	timeout := 30 * time.Second
	if deadline, ok := ctx.Deadline(); ok {
		timeout = time.Until(deadline)
	}
	reply, err := ch.CallUnary(slimrpc.ServiceName, slimrpc.MethodHandle, payload, &timeout, nil)
	if err != nil {
		return protocol.Response{}, err
	}
	return protocol.DecodeResponse(reply)
}

func nameFromString(value string) (*slim.Name, error) {
	parts := strings.Split(value, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid name format: %s", value)
	}
	return slim.NewName(parts[0], parts[1], parts[2]), nil
}
