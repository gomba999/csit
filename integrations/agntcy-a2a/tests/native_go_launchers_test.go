// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package tests

// This file isolates the native Go fixture startup path from the shared launcher
// helpers so the cross-SDK launcher layer stays focused on generic orchestration.

import "fmt"

func startGoFixture(port int, protocol transportProtocol) (*fixtureProcess, string, error) {
	return startNativeFixture(
		fmt.Sprintf("go-%s-server", protocol),
		port,
		protocol,
		"go", "run", "./fixtures/go-jsonrpc-server",
	)
}

// newGoServer returns a fixtureServer that starts a Go server for each of the given protocols.
func newGoServer(expectPushSupported bool, protocols ...transportProtocol) *fixtureServer {
	s := &fixtureServer{
		serverPrefix:        "go",
		expectPushSupported: expectPushSupported,
		runtime:             newInteropSuiteRuntime(),
	}
	for _, p := range protocols {
		p := p
		s.fixtures = append(s.fixtures, interopSuiteFixtureSpec{
			label:    "go",
			protocol: p,
			start: func() (*fixtureProcess, string, error) {
				return startGoFixture(findFreePort(), p)
			},
		})
	}
	return s
}
