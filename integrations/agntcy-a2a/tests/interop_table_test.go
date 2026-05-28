// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package tests

// Three nested DescribeTableSubtrees produce the full (client, server, transport)
// interop matrix with minimal repetition:
//
//   A2A interoperability
//     go (client)                              [CLIENT:go]
//       server go                              [SERVER:go]
//         transport JSON-RPC                   [TRANSPORT:jsonrpc]
//         transport REST                       [TRANSPORT:rest]
//         transport gRPC                       [TRANSPORT:grpc]
//       server rust                            [SERVER:rust]
//         ...
//       server dotnet                          [SERVER:dotnet]
//         transport JSON-RPC
//         transport REST
//         transport gRPC                       ← skipped: dotnet does not support gRPC
//       ...
//     rust (client)                            [CLIENT:rust]
//       ...
//     dotnet (client)                          [CLIENT:dotnet]
//       ...
//         transport gRPC                       ← skipped: dotnet does not support gRPC
//
// Each spec inherits all three category labels from its ancestor entries, enabling
// set-based label filtering:
//
//   --ginkgo.label-filter='CLIENT: consistsOf {go}'
//   --ginkgo.label-filter='CLIENT: containsAny {rust, go} && SERVER: consistsOf {dotnet}'
//   --ginkgo.label-filter='TRANSPORT: containsAny {grpc}'
//
// When either the client or the server SDK does not support gRPC, the gRPC
// BeforeAll calls Skip() so the specs appear as skipped rather than absent.
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

// ── server factory type and factories ─────────────────────────────────────────

// serverMaker creates a fresh fixtureServer for a single transport protocol.
// Called in BeforeAll once the Skip check has passed.
type serverMaker func(protocol transportProtocol) *fixtureServer

var (
	goServerMaker serverMaker = func(p transportProtocol) *fixtureServer {
		return newGoServer(getGoBinaries, true, p)
	}
	rustServerMaker serverMaker = func(p transportProtocol) *fixtureServer {
		return newRustServer(getRustBinaries, true, p)
	}
	pythonServerMaker serverMaker = func(p transportProtocol) *fixtureServer {
		return newPythonServer(getPythonAssets, true, p)
	}
	dotNetServerMaker serverMaker = func(p transportProtocol) *fixtureServer {
		return newDotNetServer(getDotNetBinaries, false, p)
	}
)

// ── interoperability matrix ───────────────────────────────────────────────────

var _ = DescribeTableSubtree(
	"A2A interoperability",
	func(newClient newClientFn, clientGRPC bool) {
		DescribeTableSubtree("server",
			func(makeServer serverMaker, serverGRPC bool, expectPushSupported bool) {
				// ContinueOnFailure is not set here: transport entries are nested
				// inside client entries (which are the outermost Ordered containers)
				// and inherit ContinueOnFailure behaviour from them.
				DescribeTableSubtree("transport",
					func(protocol transportProtocol) {
						var srv *fixtureServer
						BeforeAll(func() {
							if protocol == transportGRPC && (!clientGRPC || !serverGRPC) {
								Skip("gRPC transport is not supported by this client/server combination")
							}
							srv = makeServer(protocol)
							gomega.Expect(srv.start()).NotTo(gomega.HaveOccurred())
						})
						AfterAll(func() {
							if srv != nil {
								srv.stop()
							}
						})
						registerBehaviors(newClient, func() interopTarget {
							return srv.targetFor(protocol)()
						}, expectPushSupported)
					},
					Entry("JSON-RPC", Ordered, Label("TRANSPORT:jsonrpc"), transportJSONRPC),
					Entry("REST",     Ordered, Label("TRANSPORT:rest"),    transportREST),
					Entry("gRPC",     Ordered, Label("TRANSPORT:grpc"),    transportGRPC),
				)
			},
			// ── server entries ────────────────────────────────────────────────
			// ContinueOnFailure is inherited from the client entries (outermost Ordered).
			Entry("go",     Ordered, Label("SERVER:go"),     goServerMaker,     true,  true),
			Entry("rust",   Ordered, Label("SERVER:rust"),   rustServerMaker,   true,  true),
			Entry("python", Ordered, Label("SERVER:python"), pythonServerMaker, true,  true),
			Entry("dotnet", Ordered, Label("SERVER:dotnet"), dotNetServerMaker, false, false),
		)
	},
	// ── client entries ────────────────────────────────────────────────────────
	Entry("go",     Ordered, ContinueOnFailure, Label("CLIENT:go"),     goClientFn,     true),
	Entry("rust",   Ordered, ContinueOnFailure, Label("CLIENT:rust"),   rustClientFn,   true),
	Entry("python", Ordered, ContinueOnFailure, Label("CLIENT:python"), pythonClientFn, true),
	Entry("dotnet", Ordered, ContinueOnFailure, Label("CLIENT:dotnet"), dotNetClientFn, false),
)
