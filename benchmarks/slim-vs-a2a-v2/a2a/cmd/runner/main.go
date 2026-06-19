// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/agntcy/csit/benchmarks/slim-vs-a2a-v2/a2a/internal/client"
	"github.com/agntcy/csit/benchmarks/slim-vs-a2a-v2/internal/benchlog"
	"github.com/agntcy/csit/benchmarks/slim-vs-a2a-v2/internal/consensus"
	"github.com/agntcy/csit/benchmarks/slim-vs-a2a-v2/internal/metrics"
	"github.com/agntcy/csit/benchmarks/slim-vs-a2a-v2/internal/scenario"
)

func main() {
	scenarioPath := flag.String("scenario", "", "path to consensus scenario yaml")
	agentBin := flag.String("agent-bin", "", "path to a2a-agent binary")
	outputJSON := flag.String("output-json", "", "write run metrics json")
	outputTSV := flag.String("output-tsv", "", "append run metrics tsv")
	waitReady := flag.Duration("wait-ready", 3*time.Second, "wait for agents to start")
	quiet := flag.Bool("quiet", false, "disable benchmark logs")
	flag.Parse()

	if *quiet {
		benchlog.SetEnabled(false)
	}
	if *scenarioPath == "" {
		log.Fatal("--scenario is required")
	}

	s, err := scenario.LoadFile(*scenarioPath)
	if err != nil {
		log.Fatalf("load scenario: %v", err)
	}

	agentPath := *agentBin
	if agentPath == "" {
		agentPath = os.Getenv("A2A_AGENT_BIN")
	}
	if agentPath == "" {
		log.Fatal("set --agent-bin or A2A_AGENT_BIN")
	}

	procs := startAgents(s, agentPath, *scenarioPath)
	time.Sleep(*waitReady)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	cli, err := client.New(ctx, s)
	if err != nil {
		stopAgents(procs)
		log.Fatalf("client: %v", err)
	}

	runStart := time.Now()
	benchlog.SetRunStart(runStart)

	agentIDs := s.AgentIDs()
	if err := cli.StartAll(ctx, agentIDs); err != nil {
		stopAgents(procs)
		log.Fatalf("start agents: %v", err)
	}

	result := metrics.RunResult{
		ScenarioName:   s.Metadata.Name,
		Domain:         s.Metadata.Domain,
		Implementation: benchlog.ImplA2A,
		Agents:         len(s.Agents),
		ThinkTimeMs:    s.Spec.ThinkTimeMs,
	}

	var snapshots []consensus.AgentSnapshot
	deadline := time.Now().Add(2 * time.Minute)
	for time.Now().Before(deadline) {
		snapshots = snapshots[:0]
		var pollErr error
		for _, id := range agentIDs {
			snap, err := cli.Snapshot(ctx, id)
			if err != nil {
				pollErr = err
				break
			}
			snapshots = append(snapshots, snap)
		}
		if pollErr == nil {
			if ok, _ := consensus.GlobalConsensus(snapshots); ok {
				result.Success = true
				break
			}
		}
		time.Sleep(20 * time.Millisecond)
	}

	if !result.Success && len(snapshots) > 0 {
		if ok, _ := consensus.GlobalConsensus(snapshots); ok {
			result.Success = true
		}
	}

	result = aggregateResult(result, runStart, snapshots)

	stopAgents(procs)

	if *outputJSON != "" {
		if err := metrics.WriteJSON(*outputJSON, result); err != nil {
			log.Fatalf("write json: %v", err)
		}
	}
	if *outputTSV != "" {
		if err := metrics.AppendTSV(*outputTSV, result); err != nil {
			log.Fatalf("write tsv: %v", err)
		}
	}

	data, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(data))
	if !result.Success {
		os.Exit(1)
	}
}

func aggregateResult(result metrics.RunResult, runStart time.Time, snapshots []consensus.AgentSnapshot) metrics.RunResult {
	if len(snapshots) == 0 {
		result.Error = "no agent snapshots"
		return result
	}

	var (
		totalEmitted      int
		totalApplied      int
		lastConvergeMS    int64
		propDurations     []int64
		maxConsensusRound int
		maxRound          int
	)

	for _, snap := range snapshots {
		totalEmitted += snap.FindingsEmitted
		totalApplied += snap.FindingsApplied
		if snap.ConsensusRound > maxConsensusRound {
			maxConsensusRound = snap.ConsensusRound
		}
		if snap.Round > maxRound {
			maxRound = snap.Round
		}
		if snap.ConvergedAtNs > 0 {
			ms := (snap.ConvergedAtNs - runStart.UnixNano()) / int64(time.Millisecond)
			if ms > lastConvergeMS {
				lastConvergeMS = ms
			}
		}
		if snap.AvgPropagationMs > 0 {
			propDurations = append(propDurations, snap.AvgPropagationMs)
		}
	}

	if result.Success {
		result.ConsensusWallMS = time.Since(runStart).Milliseconds()
		if lastConvergeMS > 0 {
			result.ConsensusWallMS = lastConvergeMS
		}
	} else {
		result.ConsensusWallMS = time.Since(runStart).Milliseconds()
		result.Error = "consensus not reached"
	}

	result.ConsensusRound = maxConsensusRound
	if result.ConsensusRound == 0 {
		result.ConsensusRound = maxRound
	}
	result.FindingsEmitted = totalEmitted
	result.FindingsReceivedTotal = totalApplied
	result.LastAgentConvergeMS = lastConvergeMS
	result.AvgPropagationMS, result.P95PropagationMS = metrics.AggregatePropagation(propDurations)
	result.UnicastRPCCount = estimateA2ARPCs(result.FindingsEmitted, result.Agents)
	result.CoordFanoutMS = result.ConsensusWallMS / 4
	return result
}

func estimateA2ARPCs(findings, agents int) int {
	if agents <= 1 {
		return findings
	}
	// worker->coord + coord fan-out (N-1) per finding
	return findings * (1 + (agents - 1))
}

func startAgents(s *scenario.ConsensusScenario, agentBin, scenarioFile string) []*exec.Cmd {
	coord := s.Coordinator().ID
	var procs []*exec.Cmd
	startOne := func(agent scenario.Agent, index int, isCoord bool) {
		args := []string{
			"--agent-id", agent.ID,
			"--grpc-port", fmt.Sprintf("%d", agent.A2APort),
			"--card-port", fmt.Sprintf("%d", s.CardPort(agent)),
			"--scenario-file", scenarioFile,
			"--agent-index", fmt.Sprintf("%d", index),
		}
		if isCoord {
			args = append(args, "--coordinator")
		}
		cmd := exec.Command(agentBin, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			stopAgents(procs)
			log.Fatalf("start agent %s: %v", agent.ID, err)
		}
		procs = append(procs, cmd)
	}

	// Start workers before coordinator so fan-out targets are listening.
	for i, agent := range s.Agents {
		if agent.ID == coord {
			continue
		}
		startOne(agent, i, false)
	}
	for i, agent := range s.Agents {
		if agent.ID == coord {
			startOne(agent, i, true)
		}
	}
	return procs
}

func stopAgents(procs []*exec.Cmd) {
	for _, cmd := range procs {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
	}
}
