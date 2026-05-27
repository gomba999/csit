// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package tests

// This file implements the external language probe clients: Rust, Python, and .NET.
// Each client invokes the corresponding probe binary as a subprocess with the new
// subcommand CLI, deserializes the JSON response from stdout, and converts stderr
// error JSON to a Go error on non-zero exit.

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/a2aproject/a2a-go/v2/a2a"
)

// probeError is returned when an external probe exits with code 1.
type probeError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *probeError) Error() string {
	return e.Message
}

func parseProbeError(stderr []byte) error {
	var pe probeError
	if err := json.Unmarshal(bytes.TrimSpace(stderr), &pe); err == nil && pe.Message != "" {
		return &pe
	}
	if len(stderr) > 0 {
		return fmt.Errorf("%s", strings.TrimSpace(string(stderr)))
	}
	return fmt.Errorf("probe exited with non-zero status")
}

// probeInvoker is a function that runs the probe binary and returns its stdout/stderr.
type probeInvoker func(ctx context.Context, args ...string) (stdout []byte, stderr []byte, err error)

// externalProbeClient dispatches A2A operations to an external language probe binary.
type externalProbeClient struct {
	baseURL string
	invoke  probeInvoker
}

// probeMsgEnvelope wraps a send-message result from the probe.
type probeMsgEnvelope struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// probeEventEnvelope wraps a streaming event from the probe.
type probeEventEnvelope = probeMsgEnvelope

func (c *externalProbeClient) SendMessage(ctx context.Context, req *a2a.SendMessageRequest) (a2a.SendMessageResult, error) {
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal send-message request: %w", err)
	}

	stdout, stderr, err := c.invoke(ctx, "send-message", "--message-json", string(reqJSON))
	if err != nil {
		return nil, parseProbeError(stderr)
	}

	var env probeMsgEnvelope
	if err := json.Unmarshal(stdout, &env); err != nil {
		return nil, fmt.Errorf("parse send-message response: %w\n%s", err, stdout)
	}

	switch env.Type {
	case "task":
		var task a2a.Task
		if err := json.Unmarshal(env.Data, &task); err != nil {
			return nil, fmt.Errorf("parse task result: %w", err)
		}
		return &task, nil
	case "message":
		var msg a2a.Message
		if err := json.Unmarshal(env.Data, &msg); err != nil {
			return nil, fmt.Errorf("parse message result: %w", err)
		}
		return &msg, nil
	default:
		return nil, fmt.Errorf("unknown send-message result type %q", env.Type)
	}
}

func (c *externalProbeClient) SendStreamingMessage(ctx context.Context, req *a2a.SendMessageRequest) ([]a2a.Event, error) {
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal send-streaming-message request: %w", err)
	}

	stdout, stderr, err := c.invoke(ctx, "send-streaming-message", "--message-json", string(reqJSON))
	if err != nil {
		return nil, parseProbeError(stderr)
	}

	var events []a2a.Event
	scanner := bufio.NewScanner(bytes.NewReader(stdout))
	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			continue
		}

		var env probeEventEnvelope
		if err := json.Unmarshal(line, &env); err != nil {
			return nil, fmt.Errorf("parse streaming event line: %w\n%s", err, line)
		}

		event, err := parseStreamEvent(env)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read streaming events: %w", err)
	}

	return events, nil
}

func parseStreamEvent(env probeEventEnvelope) (a2a.Event, error) {
	switch env.Type {
	case "task":
		var task a2a.Task
		if err := json.Unmarshal(env.Data, &task); err != nil {
			return nil, fmt.Errorf("parse task event: %w", err)
		}
		return &task, nil
	case "task-status-update":
		var event a2a.TaskStatusUpdateEvent
		if err := json.Unmarshal(env.Data, &event); err != nil {
			return nil, fmt.Errorf("parse task-status-update event: %w", err)
		}
		return &event, nil
	case "task-artifact-update":
		var event a2a.TaskArtifactUpdateEvent
		if err := json.Unmarshal(env.Data, &event); err != nil {
			return nil, fmt.Errorf("parse task-artifact-update event: %w", err)
		}
		return &event, nil
	case "message":
		var msg a2a.Message
		if err := json.Unmarshal(env.Data, &msg); err != nil {
			return nil, fmt.Errorf("parse message event: %w", err)
		}
		return &msg, nil
	default:
		return nil, fmt.Errorf("unknown stream event type %q", env.Type)
	}
}

func (c *externalProbeClient) GetTask(ctx context.Context, req *a2a.GetTaskRequest) (*a2a.Task, error) {
	stdout, stderr, err := c.invoke(ctx, "get-task", "--task-id", string(req.ID))
	if err != nil {
		return nil, parseProbeError(stderr)
	}

	var task a2a.Task
	if err := json.Unmarshal(stdout, &task); err != nil {
		return nil, fmt.Errorf("parse get-task response: %w\n%s", err, stdout)
	}
	return &task, nil
}

func (c *externalProbeClient) CancelTask(ctx context.Context, req *a2a.CancelTaskRequest) (*a2a.Task, error) {
	stdout, stderr, err := c.invoke(ctx, "cancel-task", "--task-id", string(req.ID))
	if err != nil {
		return nil, parseProbeError(stderr)
	}

	var task a2a.Task
	if err := json.Unmarshal(stdout, &task); err != nil {
		return nil, fmt.Errorf("parse cancel-task response: %w\n%s", err, stdout)
	}
	return &task, nil
}

func (c *externalProbeClient) ListTasks(ctx context.Context, req *a2a.ListTasksRequest) (*a2a.ListTasksResponse, error) {
	stdout, stderr, err := c.invoke(ctx, "list-tasks", "--context-id", req.ContextID)
	if err != nil {
		return nil, parseProbeError(stderr)
	}

	var resp a2a.ListTasksResponse
	if err := json.Unmarshal(stdout, &resp); err != nil {
		return nil, fmt.Errorf("parse list-tasks response: %w\n%s", err, stdout)
	}
	return &resp, nil
}

func (c *externalProbeClient) CreatePushConfig(ctx context.Context, cfg *a2a.PushConfig) (*a2a.PushConfig, error) {
	cfgJSON, err := json.Marshal(cfg)
	if err != nil {
		return nil, fmt.Errorf("marshal push config: %w", err)
	}

	stdout, stderr, err := c.invoke(ctx, "create-push-config", "--config-json", string(cfgJSON))
	if err != nil {
		return nil, parseProbeError(stderr)
	}

	var result a2a.PushConfig
	if err := json.Unmarshal(stdout, &result); err != nil {
		return nil, fmt.Errorf("parse create-push-config response: %w\n%s", err, stdout)
	}
	return &result, nil
}

func (c *externalProbeClient) GetPushConfig(ctx context.Context, req *a2a.GetTaskPushConfigRequest) (*a2a.PushConfig, error) {
	stdout, stderr, err := c.invoke(ctx, "get-push-config",
		"--task-id", string(req.TaskID),
		"--config-id", req.ID,
	)
	if err != nil {
		return nil, parseProbeError(stderr)
	}

	var result a2a.PushConfig
	if err := json.Unmarshal(stdout, &result); err != nil {
		return nil, fmt.Errorf("parse get-push-config response: %w\n%s", err, stdout)
	}
	return &result, nil
}

func (c *externalProbeClient) ListPushConfigs(ctx context.Context, req *a2a.ListTaskPushConfigRequest) ([]*a2a.PushConfig, error) {
	stdout, stderr, err := c.invoke(ctx, "list-push-configs", "--task-id", string(req.TaskID))
	if err != nil {
		return nil, parseProbeError(stderr)
	}

	var resp a2a.ListTaskPushConfigResponse
	if err := json.Unmarshal(stdout, &resp); err != nil {
		return nil, fmt.Errorf("parse list-push-configs response: %w\n%s", err, stdout)
	}
	return resp.Configs, nil
}

func (c *externalProbeClient) DeletePushConfig(ctx context.Context, req *a2a.DeleteTaskPushConfigRequest) error {
	_, stderr, err := c.invoke(ctx, "delete-push-config",
		"--task-id", string(req.TaskID),
		"--config-id", req.ID,
	)
	if err != nil {
		return parseProbeError(stderr)
	}
	return nil
}

func (c *externalProbeClient) GetExtendedCard(ctx context.Context, req *a2a.GetExtendedAgentCardRequest) (*a2a.AgentCard, error) {
	stdout, stderr, err := c.invoke(ctx, "get-extended-card")
	if err != nil {
		return nil, parseProbeError(stderr)
	}

	var card a2a.AgentCard
	if err := json.Unmarshal(stdout, &card); err != nil {
		return nil, fmt.Errorf("parse get-extended-card response: %w\n%s", err, stdout)
	}
	return &card, nil
}

func makeProbeInvoker(binary string, globalArgs ...string) probeInvoker {
	return func(ctx context.Context, args ...string) ([]byte, []byte, error) {
		allArgs := append(globalArgs, args...)
		cmd := exec.CommandContext(ctx, binary, allArgs...)
		cmd.Dir = componentRoot()

		var stdoutBuf, stderrBuf bytes.Buffer
		cmd.Stdout = &stdoutBuf
		cmd.Stderr = &stderrBuf

		err := cmd.Run()
		return stdoutBuf.Bytes(), stderrBuf.Bytes(), err
	}
}

func newRustProbeClient(getBinaries func() fixtureBinaries, baseURL string) probeClient {
	return &externalProbeClient{
		baseURL: baseURL,
		invoke: func(ctx context.Context, args ...string) ([]byte, []byte, error) {
			invoke := makeProbeInvoker(getBinaries().rustProbe, "--card-url", baseURL)
			return invoke(ctx, args...)
		},
	}
}

func newPythonProbeClient(getAssets func() pythonFixtureAssets, baseURL string) probeClient {
	return &externalProbeClient{
		baseURL: baseURL,
		invoke: func(ctx context.Context, args ...string) ([]byte, []byte, error) {
			assets := getAssets()
			allArgs := append([]string{"run", assets.probeScript, "--card-url", baseURL}, args...)
			cmd := exec.CommandContext(ctx, assets.uvCommand, allArgs...)
			cmd.Dir = assets.fixtureDir

			var stdoutBuf, stderrBuf bytes.Buffer
			cmd.Stdout = &stdoutBuf
			cmd.Stderr = &stderrBuf

			err := cmd.Run()
			return stdoutBuf.Bytes(), stderrBuf.Bytes(), err
		},
	}
}

func newDotNetProbeClient(getBinaries func() dotNetFixtureBinaries, baseURL string) probeClient {
	return &externalProbeClient{
		baseURL: baseURL,
		invoke: func(ctx context.Context, args ...string) ([]byte, []byte, error) {
			binaries := getBinaries()
			allArgs := append([]string{binaries.dotnetProbeDL, "--card-url", baseURL}, args...)
			cmd := exec.CommandContext(ctx, binaries.dotnetCommand, allArgs...)
			cmd.Dir = componentRoot()

			var stdoutBuf, stderrBuf bytes.Buffer
			cmd.Stdout = &stdoutBuf
			cmd.Stderr = &stderrBuf

			err := cmd.Run()
			return stdoutBuf.Bytes(), stderrBuf.Bytes(), err
		},
	}
}
