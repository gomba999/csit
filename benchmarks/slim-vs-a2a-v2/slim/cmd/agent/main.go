// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	slim "github.com/agntcy/slim-bindings-go"
	"github.com/agntcy/csit/benchmarks/slim-vs-a2a-v2/internal/scenario"
	"github.com/agntcy/csit/benchmarks/slim-vs-a2a-v2/slim/internal/agent"
)

func main() {
	slimName := flag.String("slim-name", "", "SLIM identity org/group/app")
	endpoint := flag.String("endpoint", "http://127.0.0.1:46357", "SLIM dataplane endpoint")
	scenarioFile := flag.String("scenario-file", "", "path to consensus scenario yaml")
	agentIndex := flag.Int("agent-index", 0, "agent index in scenario")
	flag.Parse()

	if *slimName == "" || *scenarioFile == "" {
		log.Fatal("--slim-name and --scenario-file are required")
	}

	s, err := scenario.LoadFile(*scenarioFile)
	if err != nil {
		log.Fatalf("load scenario: %v", err)
	}
	if *agentIndex < 0 || *agentIndex >= len(s.Agents) {
		log.Fatalf("agent-index out of range")
	}

	rt := agent.NewRuntime(s, *agentIndex, *slimName, *endpoint)
	if err := rt.Setup(); err != nil {
		log.Fatalf("setup: %v", err)
	}
	defer rt.Close()

	name, err := nameFromString(*slimName)
	if err != nil {
		log.Fatalf("parse name: %v", err)
	}

	server := slim.ServerNewWithConnection(rt.App(), name, rt.ConnID())
	rt.Register(server)

	fmt.Printf("SLIM_AGENT_READY name=%s index=%d scenario=%s\n", *slimName, *agentIndex, s.Metadata.Name)
	if err := server.Serve(); err != nil {
		log.Printf("serve: %v", err)
	}
}

func nameFromString(value string) (*slim.Name, error) {
	parts := strings.Split(value, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid name format: %s", value)
	}
	return slim.NewName(parts[0], parts[1], parts[2]), nil
}
