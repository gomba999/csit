// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package tests

// This file isolates the native Go fixture startup path from the shared launcher
// helpers so the cross-SDK launcher layer stays focused on generic orchestration.

import "fmt"

func startGoFixture(protocol transportProtocol) (*fixtureProcess, string, error) {
	return startNativeFixture(
		fmt.Sprintf("go-%s-server", protocol),
		findFreePort(),
		protocol,
		"go", "run", "./fixtures/go-jsonrpc-server",
	)
}
