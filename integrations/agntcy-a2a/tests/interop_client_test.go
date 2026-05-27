// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package tests

// This file defines the probeClient interface and the native Go SDK implementation.
// All interop behavior specs are written against probeClient so the same spec tree
// exercises Go, Rust, Python, and .NET clients without duplication.

import (
	"context"

	"github.com/a2aproject/a2a-go/v2/a2a"
	"github.com/a2aproject/a2a-go/v2/a2aclient"
	"github.com/a2aproject/a2a-go/v2/a2aclient/agentcard"
	a2agrpc "github.com/a2aproject/a2a-go/v2/a2agrpc/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

// goProbeClient wraps the native Go A2A SDK client.
type goProbeClient struct {
	client *a2aclient.Client
}

func newGoProbeClient(ctx context.Context, baseURL string) (probeClient, error) {
	card, err := agentcard.DefaultResolver.Resolve(ctx, baseURL)
	if err != nil {
		return nil, err
	}

	client, err := a2aclient.NewFromCard(
		ctx,
		card,
		a2agrpc.WithGRPCTransport(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	if err != nil {
		return nil, err
	}

	return &goProbeClient{client: client}, nil
}

func (c *goProbeClient) SendMessage(ctx context.Context, req *a2a.SendMessageRequest) (a2a.SendMessageResult, error) {
	return c.client.SendMessage(ctx, req)
}

func (c *goProbeClient) SendStreamingMessage(ctx context.Context, req *a2a.SendMessageRequest) ([]a2a.Event, error) {
	var events []a2a.Event
	for event, err := range c.client.SendStreamingMessage(ctx, req) {
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (c *goProbeClient) GetTask(ctx context.Context, req *a2a.GetTaskRequest) (*a2a.Task, error) {
	return c.client.GetTask(ctx, req)
}

func (c *goProbeClient) CancelTask(ctx context.Context, req *a2a.CancelTaskRequest) (*a2a.Task, error) {
	return c.client.CancelTask(ctx, req)
}

func (c *goProbeClient) ListTasks(ctx context.Context, req *a2a.ListTasksRequest) (*a2a.ListTasksResponse, error) {
	return c.client.ListTasks(ctx, req)
}

func (c *goProbeClient) CreatePushConfig(ctx context.Context, cfg *a2a.PushConfig) (*a2a.PushConfig, error) {
	return c.client.CreateTaskPushConfig(ctx, cfg)
}

func (c *goProbeClient) GetPushConfig(ctx context.Context, req *a2a.GetTaskPushConfigRequest) (*a2a.PushConfig, error) {
	return c.client.GetTaskPushConfig(ctx, req)
}

func (c *goProbeClient) ListPushConfigs(ctx context.Context, req *a2a.ListTaskPushConfigRequest) ([]*a2a.PushConfig, error) {
	return c.client.ListTaskPushConfigs(ctx, req)
}

func (c *goProbeClient) DeletePushConfig(ctx context.Context, req *a2a.DeleteTaskPushConfigRequest) error {
	return c.client.DeleteTaskPushConfig(ctx, req)
}

func (c *goProbeClient) GetExtendedCard(ctx context.Context, req *a2a.GetExtendedAgentCardRequest) (*a2a.AgentCard, error) {
	return c.client.GetExtendedAgentCard(ctx, req)
}
