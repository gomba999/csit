// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package protocol

import (
	"encoding/json"
)

const (
	OpFinding     = "finding"
	OpFanout      = "fanout"
	OpStart       = "start"
	OpSnapshot    = "snapshot"
	OpReady       = "ready"
)

type Request struct {
	Op          string `json:"op"`
	FindingJSON string `json:"findingJson,omitempty"`
	Body        string `json:"body,omitempty"`
}

type Response struct {
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
	Body  string `json:"body,omitempty"`
}

func EncodeRequest(req Request) ([]byte, error) {
	return json.Marshal(req)
}

func DecodeRequest(text string) (Request, error) {
	var req Request
	err := json.Unmarshal([]byte(text), &req)
	return req, err
}

func EncodeResponse(resp Response) ([]byte, error) {
	return json.Marshal(resp)
}

func DecodeResponse(text string) (Response, error) {
	var resp Response
	err := json.Unmarshal([]byte(text), &resp)
	return resp, err
}

func EncodeResponseBytes(resp Response) (string, error) {
	data, err := json.Marshal(resp)
	return string(data), err
}

func SnapshotBody(snapshot any) (string, error) {
	data, err := json.Marshal(snapshot)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
