// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/a2aproject/a2a-go/v2/a2a"
	"github.com/a2aproject/a2a-go/v2/a2aclient"
	"github.com/a2aproject/a2a-go/v2/a2aclient/agentcard"
	a2agrpc "github.com/a2aproject/a2a-go/v2/a2agrpc/v1"
	"github.com/agntcy/csit/benchmarks/slim-vs-a2a-v2/a2a/internal/protocol"
	"github.com/agntcy/csit/benchmarks/slim-vs-a2a-v2/internal/benchlog"
	"github.com/agntcy/csit/benchmarks/slim-vs-a2a-v2/internal/consensus"
	"github.com/agntcy/csit/benchmarks/slim-vs-a2a-v2/internal/scenario"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Runtime struct {
	scenario    *scenario.ConsensusScenario
	agentIndex  int
	agentID     string
	isCoordinator bool

	engine *consensus.Engine

	mu            sync.Mutex
	clients       map[string]*a2aclient.Client
	coordinatorID string
	running       atomic.Bool
	done          chan struct{}
	startedAt     time.Time

	unicastRPCCount int
	coordFanoutMS   int64
	clientsOnce sync.Once
	clientsErr   error
}

func NewRuntime(s *scenario.ConsensusScenario, agentIndex int, agentID string, isCoordinator bool) *Runtime {
	coord := s.Coordinator().ID
	return &Runtime{
		scenario:      s,
		agentIndex:    agentIndex,
		agentID:       agentID,
		isCoordinator: isCoordinator,
		engine:        consensus.NewEngine(s.Spec, agentIndex),
		clients:       map[string]*a2aclient.Client{},
		coordinatorID: coord,
		done:          make(chan struct{}),
	}
}

func (r *Runtime) ensureClients(ctx context.Context) error {
	r.clientsOnce.Do(func() {
		r.clientsErr = r.setupClients(ctx)
	})
	return r.clientsErr
}

func (r *Runtime) setupClients(ctx context.Context) error {
	var agents []scenario.Agent
	if r.isCoordinator {
		for _, a := range r.scenario.Agents {
			if a.ID != r.agentID {
				agents = append(agents, a)
			}
		}
	} else {
		coord, ok := r.scenario.AgentByID(r.coordinatorID)
		if !ok {
			return fmt.Errorf("coordinator not found")
		}
		agents = []scenario.Agent{coord}
	}

	for _, a := range agents {
		var cli *a2aclient.Client
		var err error
		for attempt := 0; attempt < 50; attempt++ {
			cli, err = dialAgent(ctx, r.scenario, a)
			if err == nil {
				break
			}
			time.Sleep(200 * time.Millisecond)
		}
		if err != nil {
			return err
		}
		r.clients[a.ID] = cli
	}
	return nil
}

func dialAgent(ctx context.Context, s *scenario.ConsensusScenario, agent scenario.Agent) (*a2aclient.Client, error) {
	baseURL := s.CardBaseURL(agent)
	card, err := agentcard.DefaultResolver.Resolve(ctx, baseURL)
	if err != nil {
		return nil, fmt.Errorf("resolve card for %s: %w", agent.ID, err)
	}
	return a2aclient.NewFromCard(
		ctx,
		card,
		a2agrpc.WithGRPCTransport(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
}

func (r *Runtime) Handle(ctx context.Context, req protocol.Request) protocol.Response {
	switch req.Op {
	case protocol.OpStart:
		if err := r.ensureClients(ctx); err != nil {
			return protocol.Response{OK: false, Error: err.Error()}
		}
		r.StartRun()
		return protocol.Response{OK: true}
	case protocol.OpSnapshot:
		body, err := json.Marshal(r.engine.Snapshot())
		if err != nil {
			return protocol.Response{OK: false, Error: err.Error()}
		}
		return protocol.Response{OK: true, Body: string(body)}
	case protocol.OpFinding:
		if !r.isCoordinator {
			return protocol.Response{OK: false, Error: "not coordinator"}
		}
		f, err := consensus.DecodeFinding([]byte(req.FindingJSON))
		if err != nil {
			return protocol.Response{OK: false, Error: err.Error()}
		}
		if err := r.handleCoordinatorFinding(ctx, f); err != nil {
			return protocol.Response{OK: false, Error: err.Error()}
		}
		return protocol.Response{OK: true}
	case protocol.OpFanout:
		f, err := consensus.DecodeFinding([]byte(req.FindingJSON))
		if err != nil {
			return protocol.Response{OK: false, Error: err.Error()}
		}
		if f.AgentIndex == r.agentIndex {
			return protocol.Response{OK: true}
		}
		r.applyFinding(f)
		return protocol.Response{OK: true}
	default:
		return protocol.Response{OK: false, Error: "unknown op"}
	}
}

func (r *Runtime) StartRun() {
	if r.running.Swap(true) {
		return
	}
	r.startedAt = time.Now()
	benchlog.SetRunStart(r.startedAt)
	go r.runLoop(context.Background())
}

func (r *Runtime) runLoop(ctx context.Context) {
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
			if err := r.publishFinding(ctx, *finding); err != nil {
				benchlog.Finding(benchlog.ImplA2A, "publish_error", r.agentIndex, finding.FindingID,
					fmt.Sprintf("err=%v", err))
			}
			time.Sleep(emitGap)
		}

		if r.engine.HasLocalConsensus() {
			return
		}
	}
}

func (r *Runtime) publishFinding(ctx context.Context, f consensus.Finding) error {
	if err := r.ensureClients(ctx); err != nil {
		return err
	}
	body, err := json.Marshal(f)
	if err != nil {
		return err
	}
	if r.isCoordinator {
		return r.handleCoordinatorFinding(ctx, f)
	}
	req := protocol.Request{Op: protocol.OpFinding, FindingJSON: string(body)}
	return r.send(ctx, r.coordinatorID, req)
}

func (r *Runtime) handleCoordinatorFinding(ctx context.Context, f consensus.Finding) error {
	if f.AgentIndex != r.agentIndex {
		r.applyFinding(f)
	}
	return r.fanout(ctx, f)
}

func (r *Runtime) fanout(ctx context.Context, f consensus.Finding) error {
	body, err := json.Marshal(f)
	if err != nil {
		return err
	}
	req := protocol.Request{Op: protocol.OpFanout, FindingJSON: string(body)}

	start := time.Now()
	var wg sync.WaitGroup
	errCh := make(chan error, len(r.clients))
	for agentID := range r.clients {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			if err := r.send(ctx, id, req); err != nil {
				errCh <- err
			}
			r.mu.Lock()
			r.unicastRPCCount++
			r.mu.Unlock()
		}(agentID)
	}
	wg.Wait()
	close(errCh)

	r.mu.Lock()
	r.coordFanoutMS += time.Since(start).Milliseconds()
	r.mu.Unlock()

	for err := range errCh {
		if err != nil {
			return err
		}
	}
	benchlog.Finding(benchlog.ImplA2A, "fanout", r.agentIndex, f.FindingID,
		fmt.Sprintf("targets=%d", len(r.clients)))
	return nil
}

func (r *Runtime) send(ctx context.Context, agentID string, req protocol.Request) error {
	cli, ok := r.clients[agentID]
	if !ok {
		return fmt.Errorf("unknown client %q", agentID)
	}
	text, err := json.Marshal(req)
	if err != nil {
		return err
	}
	msg := a2a.NewMessage(a2a.MessageRoleUser, a2a.NewTextPart(string(text)))
	_, err = cli.SendMessage(ctx, &a2a.SendMessageRequest{Message: msg})
	if !r.isCoordinator {
		r.mu.Lock()
		r.unicastRPCCount++
		r.mu.Unlock()
	}
	return err
}

func (r *Runtime) applyFinding(f consensus.Finding) {
	recvNs := time.Now().UnixNano()
	r.engine.ApplyFinding(f)
	if f.EmittedAt > 0 {
		r.engine.RecordPropagation(f.EmittedAt, recvNs)
	}
	benchlog.Finding(benchlog.ImplA2A, "received", r.agentIndex, f.FindingID,
		fmt.Sprintf("from=%d", f.AgentIndex))
}

func (r *Runtime) Snapshot() consensus.AgentSnapshot {
	return r.engine.Snapshot()
}

func (r *Runtime) Close() {
	for _, cli := range r.clients {
		_ = cli
	}
}

func DecodeRequestText(text string) (protocol.Request, error) {
	return protocol.DecodeRequest(text)
}

func EncodeResponse(resp protocol.Response) (string, error) {
	data, err := json.Marshal(resp)
	return string(data), err
}
