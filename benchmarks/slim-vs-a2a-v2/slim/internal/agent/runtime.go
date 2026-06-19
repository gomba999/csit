// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package agent

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	slim "github.com/agntcy/slim-bindings-go"
	"github.com/agntcy/csit/benchmarks/slim-vs-a2a-v2/internal/benchlog"
	"github.com/agntcy/csit/benchmarks/slim-vs-a2a-v2/internal/consensus"
	"github.com/agntcy/csit/benchmarks/slim-vs-a2a-v2/internal/scenario"
	"github.com/agntcy/csit/benchmarks/slim-vs-a2a-v2/slim/internal/protocol"
	"github.com/agntcy/csit/benchmarks/slim-vs-a2a-v2/slim/internal/slimrpc"
)

const (
	defaultSharedSecret = "demo-shared-secret-min-32-chars!!"
	runnerSlimName      = "agntcy/bench-v2/runner"
)

type Runtime struct {
	scenario   *scenario.ConsensusScenario
	agentIndex int
	slimName   string
	endpoint   string

	mu        sync.Mutex
	engine    *consensus.Engine
	connID    uint64
	app       *slim.App
	relayCh   *slim.Channel
	running   atomic.Bool
	done      chan struct{}
	startedAt time.Time
}

func NewRuntime(s *scenario.ConsensusScenario, agentIndex int, slimName, endpoint string) *Runtime {
	return &Runtime{
		scenario:   s,
		agentIndex: agentIndex,
		slimName:   slimName,
		endpoint:   endpoint,
		engine:     consensus.NewEngine(s.Spec, agentIndex),
		done:       make(chan struct{}),
	}
}

func (r *Runtime) Setup() error {
	slim.InitializeWithDefaults()
	service := slim.GetGlobalService()

	name, err := nameFromString(r.slimName)
	if err != nil {
		return err
	}

	app, err := service.CreateAppWithSecret(name, defaultSharedSecret)
	if err != nil {
		return err
	}
	r.app = app

	connID, err := service.Connect(slim.NewInsecureClientConfig(r.endpoint))
	if err != nil {
		app.Destroy()
		return err
	}
	r.connID = connID
	if err := app.Subscribe(name, &connID); err != nil {
		app.Destroy()
		return err
	}

	runnerName, err := nameFromString(runnerSlimName)
	if err != nil {
		app.Destroy()
		return err
	}
	if err := app.SetRoute(runnerName, connID); err != nil {
		app.Destroy()
		return fmt.Errorf("set route runner: %w", err)
	}
	r.relayCh = slim.ChannelNewWithConnection(app, runnerName, &connID)
	return nil
}

func (r *Runtime) Register(server *slim.Server) {
	server.RegisterUnaryUnary(slimrpc.ServiceName, slimrpc.MethodHandle, &controlHandler{runtime: r})
	server.RegisterStreamUnary(slimrpc.ServiceName, slimrpc.MethodFindingStream, &findingStreamHandler{runtime: r})
}

func (r *Runtime) StartRun() {
	if r.running.Swap(true) {
		return
	}
	r.startedAt = time.Now()
	benchlog.SetRunStart(r.startedAt)
	go r.runLoop()
}

func (r *Runtime) runLoop() {
	defer close(r.done)
	spec := r.scenario.Spec
	think := time.Duration(spec.ThinkTimeMs) * time.Millisecond
	emitGap := time.Duration(spec.FindingEmitDelayMs) * time.Millisecond

	for round := 0; round < spec.MaxRounds; round++ {
		select {
		case <-r.done:
			return
		default:
		}

		time.Sleep(think)

		finding, emit := r.engine.Think()
		if emit && finding != nil {
			if err := r.publishFinding(*finding); err != nil {
				benchlog.Finding(benchlog.ImplSLIMStream, "publish_error", r.agentIndex, finding.FindingID,
					fmt.Sprintf("err=%v", err))
			} else {
				benchlog.Finding(benchlog.ImplSLIMStream, "published", r.agentIndex, finding.FindingID)
			}
			time.Sleep(emitGap)
		}

		if r.engine.HasLocalConsensus() {
			return
		}
	}
}

func (r *Runtime) publishFinding(f consensus.Finding) error {
	payload, err := consensus.EncodeFinding(f)
	if err != nil {
		return err
	}
	timeout := 30 * time.Second
	_, err = r.relayCh.CallUnary(slimrpc.ServiceName, slimrpc.MethodPublish, payload, &timeout, nil)
	return err
}

func (r *Runtime) applyFinding(f consensus.Finding) {
	recvNs := time.Now().UnixNano()
	r.engine.ApplyFinding(f)
	if f.EmittedAt > 0 {
		r.engine.RecordPropagation(f.EmittedAt, recvNs)
	}
}

func (r *Runtime) Snapshot() consensus.AgentSnapshot {
	return r.engine.Snapshot()
}

func (r *Runtime) App() *slim.App { return r.app }
func (r *Runtime) ConnID() *uint64 {
	if r.connID == 0 {
		return nil
	}
	id := r.connID
	return &id
}

func (r *Runtime) Close() {
	if r.relayCh != nil {
		_ = r.relayCh.Close(nil)
	}
	if r.app != nil {
		r.app.Destroy()
	}
}

type controlHandler struct {
	runtime *Runtime
}

func (h *controlHandler) Handle(request []byte, _ *slim.Context) ([]byte, error) {
	req, err := protocol.DecodeRequest(request)
	if err != nil {
		return nil, err
	}
	switch req.Op {
	case protocol.OpStart:
		h.runtime.StartRun()
		return protocol.EncodeResponse(protocol.Response{OK: true})
	case protocol.OpSnapshot:
		body, err := json.Marshal(h.runtime.Snapshot())
		if err != nil {
			return nil, err
		}
		return protocol.EncodeResponse(protocol.Response{OK: true, Body: string(body)})
	default:
		return protocol.EncodeResponse(protocol.Response{OK: false, Error: "unknown op"})
	}
}

type findingStreamHandler struct {
	runtime *Runtime
}

func (h *findingStreamHandler) Handle(stream *slim.RequestStream, _ *slim.Context) ([]byte, error) {
	for {
		switch item := stream.Next().(type) {
		case slim.StreamMessageEnd:
			return []byte(`{"ok":true}`), nil
		case slim.StreamMessageError:
			if item.Field0 != nil {
				return nil, item.Field0.AsError()
			}
			return nil, fmt.Errorf("stream error")
		case slim.StreamMessageData:
			f, err := consensus.DecodeFinding(item.Field0)
			if err != nil {
				continue
			}
			if f.AgentIndex == h.runtime.agentIndex {
				continue
			}
			h.runtime.applyFinding(f)
			benchlog.Finding(benchlog.ImplSLIMStream, "received", h.runtime.agentIndex, f.FindingID,
				fmt.Sprintf("from=%d", f.AgentIndex))
		}
	}
}

func nameFromString(value string) (*slim.Name, error) {
	parts := strings.Split(value, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid name format: %s", value)
	}
	return slim.NewName(parts[0], parts[1], parts[2]), nil
}
