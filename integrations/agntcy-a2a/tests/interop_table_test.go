// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package tests

// Three nested DescribeTableSubtrees produce the full (client, server, transport)
// interop matrix with minimal repetition:
//
//   A2A interoperability
//     go (client)
//       server go
//         transport JSON-RPC  →  BeforeAll/registerBehaviors/AfterAll
//         transport REST
//         transport gRPC
//       server rust
//         ...
//       server dotnet
//         transport JSON-RPC
//         transport REST
//         transport gRPC       ← skipped: dotnet server does not support gRPC
//       ...
//     rust (client)
//       ...
//     dotnet (client)
//       ...
//         transport gRPC       ← skipped: dotnet client does not support gRPC
//
// All three transports are always present in the tree.  When either the client
// or the server SDK does not support gRPC, the gRPC BeforeAll calls Skip() so
// that the specs appear as skipped in the report rather than being absent.
//
// Each SDK's fixture assets are built at most once per test run via package-level
// sync.Once caches. An AfterSuite handler removes the temp directories.

import (
	"context"
	"sync"

	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

// ── per-SDK asset caches ──────────────────────────────────────────────────────

var (
	onceGo   sync.Once
	cachedGo fixtureBinaries
	errGo    error

	onceRust   sync.Once
	cachedRust fixtureBinaries
	errRust    error

	onceDotNet   sync.Once
	cachedDotNet dotNetFixtureBinaries
	errDotNet    error

	oncePython   sync.Once
	cachedPython pythonFixtureAssets
	errPython    error
)

func getGoBinaries() fixtureBinaries {
	onceGo.Do(func() { cachedGo, errGo = buildGoFixtureBinaryOnly() })
	gomega.Expect(errGo).NotTo(gomega.HaveOccurred(), "build Go fixtures")
	return cachedGo
}

func getRustBinaries() fixtureBinaries {
	onceRust.Do(func() { cachedRust, errRust = buildRustFixtureBinaryOnly() })
	gomega.Expect(errRust).NotTo(gomega.HaveOccurred(), "build Rust fixtures")
	return cachedRust
}

func getDotNetBinaries() dotNetFixtureBinaries {
	onceDotNet.Do(func() { cachedDotNet, errDotNet = buildDotNetFixtureBinaryOnly() })
	gomega.Expect(errDotNet).NotTo(gomega.HaveOccurred(), "build .NET fixtures")
	return cachedDotNet
}

func getPythonAssets() pythonFixtureAssets {
	oncePython.Do(func() { cachedPython, errPython = buildPythonFixtureAssets() })
	gomega.Expect(errPython).NotTo(gomega.HaveOccurred(), "build Python fixtures")
	return cachedPython
}

var _ = AfterSuite(func() {
	removeTempDir(cachedGo.tempDir)
	removeTempDir(cachedRust.tempDir)
	removeTempDir(cachedDotNet.tempDir)
})

// ── client factories ──────────────────────────────────────────────────────────

var (
	goClientFn newClientFn = func(ctx context.Context, url string) (probeClient, error) {
		return newGoProbeClient(ctx, url)
	}
	rustClientFn newClientFn = func(_ context.Context, url string) (probeClient, error) {
		return newRustProbeClient(getRustBinaries, url), nil
	}
	dotNetClientFn newClientFn = func(_ context.Context, url string) (probeClient, error) {
		return newDotNetProbeClient(getDotNetBinaries, url), nil
	}
	pythonClientFn newClientFn = func(_ context.Context, url string) (probeClient, error) {
		return newPythonProbeClient(getPythonAssets, url), nil
	}
)

// ── server factory type ───────────────────────────────────────────────────────

// serverMaker creates a fresh fixtureServer configured for a single transport
// protocol.  Called at tree-building time; the actual process is started in BeforeAll.
type serverMaker func(protocol transportProtocol) *fixtureServer

// ── interoperability matrix ───────────────────────────────────────────────────

var _ = DescribeTableSubtree(
	"A2A interoperability",
	func(clientSDK string, newClient newClientFn, clientGRPC bool) {
		DescribeTableSubtree("server",
			func(serverSDK string, makeServer serverMaker, serverGRPC bool) {
				label := clientSDK + "-" + serverSDK

				// ContinueOnFailure is not set here: transport entries are nested
				// inside client entries (which are the outermost Ordered containers)
				// and inherit ContinueOnFailure behaviour from them.
				DescribeTableSubtree("transport",
					func(protocol transportProtocol) {
						srv := makeServer(protocol)
						BeforeAll(func() {
							if protocol == transportGRPC && (!clientGRPC || !serverGRPC) {
								Skip("gRPC transport is not supported by this client/server combination")
							}
							gomega.Expect(srv.start()).NotTo(gomega.HaveOccurred())
						})
						AfterAll(func() { srv.stop() })
						registerBehaviors(newClient, srv.targetFor(protocol))
					},
					Entry("JSON-RPC", Ordered, Label(label, "jsonrpc"), transportJSONRPC),
					Entry("REST",     Ordered, Label(label, "rest"),    transportREST),
					Entry("gRPC",     Ordered, Label(label, "grpc"),    transportGRPC),
				)
			},
			// ── server entries ────────────────────────────────────────────────
			// ContinueOnFailure is inherited from the client entries (outermost Ordered).
			Entry("go", Ordered,
				"go", serverMaker(func(p transportProtocol) *fixtureServer {
					return newGoServer(getGoBinaries, true, p)
				}), true),
			Entry("rust", Ordered,
				"rust", serverMaker(func(p transportProtocol) *fixtureServer {
					return newRustServer(getRustBinaries, true, p)
				}), true),
			Entry("python", Ordered,
				"python", serverMaker(func(p transportProtocol) *fixtureServer {
					return newPythonServer(getPythonAssets, true, p)
				}), true),
			Entry("dotnet", Ordered,
				"dotnet", serverMaker(func(p transportProtocol) *fixtureServer {
					return newDotNetServer(getDotNetBinaries, false, p)
				}), false),
		)
	},
	// ── client entries ────────────────────────────────────────────────────────
	Entry("go", Ordered, ContinueOnFailure, "go", goClientFn, true),
	Entry("rust", Ordered, ContinueOnFailure, "rust", rustClientFn, true),
	Entry("python", Ordered, ContinueOnFailure, "python", pythonClientFn, true),
	Entry("dotnet", Ordered, ContinueOnFailure, "dotnet", dotNetClientFn, false),
)
