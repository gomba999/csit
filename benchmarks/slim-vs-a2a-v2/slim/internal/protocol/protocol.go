// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package protocol

import (
	"encoding/json"
)

const (
	OpStart    = "start"
	OpSnapshot = "snapshot"
)

type Request struct {
	Op string `json:"op"`
}

type Response struct {
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
	Body  string `json:"body,omitempty"`
}

func EncodeRequest(req Request) ([]byte, error) {
	return json.Marshal(req)
}

func DecodeRequest(data []byte) (Request, error) {
	var req Request
	err := json.Unmarshal(data, &req)
	return req, err
}

func EncodeResponse(resp Response) ([]byte, error) {
	return json.Marshal(resp)
}

func DecodeResponse(data []byte) (Response, error) {
	var resp Response
	err := json.Unmarshal(data, &resp)
	return resp, err
}
