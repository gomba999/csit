// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package tests

// This file isolates the native Go fixture build and startup path from the shared launcher
// helpers so the cross-SDK launcher layer stays focused on generic orchestration and the
// test-only native Go startup path is easy to find.

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func buildGoFixture(root string, serverPath string, probePath string) error {
	buildCtx, cancel := context.WithTimeout(context.Background(), buildTimeout)
	defer cancel()

	goBuildServer := exec.CommandContext(buildCtx, "go", "build", "-o", serverPath, "./fixtures/go-jsonrpc-server")
	goBuildServer.Dir = root
	if output, err := goBuildServer.CombinedOutput(); err != nil {
		return fmt.Errorf("build go server fixture: %w\n%s", err, string(output))
	}

	goBuildProbe := exec.CommandContext(buildCtx, "go", "build", "-o", probePath, "./fixtures/go-probe")
	goBuildProbe.Dir = root
	if output, err := goBuildProbe.CombinedOutput(); err != nil {
		return fmt.Errorf("build go probe fixture: %w\n%s", err, string(output))
	}

	return nil
}

func buildGoFixtureBinaryOnly() (fixtureBinaries, error) {
	root := componentRoot()
	tempDir, err := os.MkdirTemp("", "agntcy-a2a-go-")
	if err != nil {
		return fixtureBinaries{}, fmt.Errorf("create temp dir: %w", err)
	}

	binaries := fixtureBinaries{
		tempDir:  tempDir,
		goServer: filepath.Join(tempDir, executableName("go-jsonrpc-server")),
		goProbe:  filepath.Join(tempDir, executableName("go-probe")),
	}

	if err := buildGoFixture(root, binaries.goServer, binaries.goProbe); err != nil {
		_ = os.RemoveAll(tempDir)
		return fixtureBinaries{}, err
	}

	return binaries, nil
}

func startGoFixture(binaries fixtureBinaries, port int, protocol transportProtocol) (*fixtureProcess, string, error) {
	return startNativeFixture(fmt.Sprintf("go-%s-server", protocol), binaries.goServer, port, protocol)
}

// newGoServer returns a fixtureServer that starts a Go server for each of the given protocols.
func newGoServer(getBinaries func() fixtureBinaries, expectPushSupported bool, protocols ...transportProtocol) *fixtureServer {
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
				return startGoFixture(getBinaries(), findFreePort(), p)
			},
		})
	}
	return s
}
