// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package tests

// This file defines the probeClient interface and the newClientFn type.
// All interop behavior specs are written against probeClient so the same spec tree
// exercises Go, Rust, Python, and .NET clients without duplication.
// Client implementations live in interop_probe_clients_test.go.

import (
	"context"

	"github.com/a2aproject/a2a-go/v2/a2a"
)

// probeClient is the unified A2A client interface used by all behavior specs.
// Both the native Go SDK and external language probes implement this.
type probeClient interface {
	SendMessage(ctx context.Context, req *a2a.SendMessageRequest) (a2a.SendMessageResult, error)
	SendStreamingMessage(ctx context.Context, req *a2a.SendMessageRequest) ([]a2a.Event, error)
	GetTask(ctx context.Context, req *a2a.GetTaskRequest) (*a2a.Task, error)
	CancelTask(ctx context.Context, req *a2a.CancelTaskRequest) (*a2a.Task, error)
	ListTasks(ctx context.Context, req *a2a.ListTasksRequest) (*a2a.ListTasksResponse, error)
	CreatePushConfig(ctx context.Context, cfg *a2a.PushConfig) (*a2a.PushConfig, error)
	GetPushConfig(ctx context.Context, req *a2a.GetTaskPushConfigRequest) (*a2a.PushConfig, error)
	ListPushConfigs(ctx context.Context, req *a2a.ListTaskPushConfigRequest) ([]*a2a.PushConfig, error)
	DeletePushConfig(ctx context.Context, req *a2a.DeleteTaskPushConfigRequest) error
	GetExtendedCard(ctx context.Context, req *a2a.GetExtendedAgentCardRequest) (*a2a.AgentCard, error)
}

// newClientFn creates a probeClient bound to a specific baseURL.
type newClientFn func(ctx context.Context, baseURL string) (probeClient, error)
