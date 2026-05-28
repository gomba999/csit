// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package tests

// This file replaces the 10 per-SDK-pair suite wrapper files with a single
// DescribeTable that has one Entry per (client SDK, server SDK, transport) combination.
//
// Each SDK's fixture assets are built at most once per test run via package-level
// sync.Once caches. Each Entry is an independent ordered Context: its BeforeAll
// starts a dedicated server instance and its AfterAll stops it. The shared
// registerBehaviors call registers the full When/Context/It spec tree.

import (
	"context"
	"sync"

	ginkgo "github.com/onsi/ginkgo/v2"
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

var _ = ginkgo.AfterSuite(func() {
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

// ── interoperability table ────────────────────────────────────────────────────

var _ = ginkgo.DescribeTableSubtree(
	"A2A interoperability",
	func(newClient newClientFn, server *fixtureServer, protocol transportProtocol) {
		ginkgo.BeforeAll(func() {
			gomega.Expect(server.start()).NotTo(gomega.HaveOccurred())
		})
		ginkgo.AfterAll(func() { server.stop() })
		registerBehaviors(newClient, server.targetFor(protocol))
	},

	// ── Go client ─────────────────────────────────────────────────────────────
	ginkgo.Entry("go → go [jsonrpc]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("go-go", "jsonrpc"),
		goClientFn, newGoServer(getGoBinaries, true, transportJSONRPC), transportJSONRPC),
	ginkgo.Entry("go → go [rest]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("go-go", "rest"),
		goClientFn, newGoServer(getGoBinaries, true, transportREST), transportREST),
	ginkgo.Entry("go → go [grpc]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("go-go", "grpc"),
		goClientFn, newGoServer(getGoBinaries, true, transportGRPC), transportGRPC),

	ginkgo.Entry("go → rust [jsonrpc]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("go-rust", "jsonrpc"),
		goClientFn, newRustServer(getRustBinaries, true, transportJSONRPC), transportJSONRPC),
	ginkgo.Entry("go → rust [rest]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("go-rust", "rest"),
		goClientFn, newRustServer(getRustBinaries, true, transportREST), transportREST),
	ginkgo.Entry("go → rust [grpc]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("go-rust", "grpc"),
		goClientFn, newRustServer(getRustBinaries, true, transportGRPC), transportGRPC),

	ginkgo.Entry("go → python [jsonrpc]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("go-python", "jsonrpc"),
		goClientFn, newPythonServer(getPythonAssets, true, transportJSONRPC), transportJSONRPC),
	ginkgo.Entry("go → python [rest]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("go-python", "rest"),
		goClientFn, newPythonServer(getPythonAssets, true, transportREST), transportREST),
	ginkgo.Entry("go → python [grpc]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("go-python", "grpc"),
		goClientFn, newPythonServer(getPythonAssets, true, transportGRPC), transportGRPC),

	ginkgo.Entry("go → dotnet [jsonrpc]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("go-dotnet", "jsonrpc"),
		goClientFn, newDotNetServer(getDotNetBinaries, false, transportJSONRPC), transportJSONRPC),
	ginkgo.Entry("go → dotnet [rest]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("go-dotnet", "rest"),
		goClientFn, newDotNetServer(getDotNetBinaries, false, transportREST), transportREST),

	// ── Rust client ───────────────────────────────────────────────────────────
	ginkgo.Entry("rust → go [jsonrpc]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("rust-go", "jsonrpc"),
		rustClientFn, newGoServer(getGoBinaries, true, transportJSONRPC), transportJSONRPC),
	ginkgo.Entry("rust → go [rest]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("rust-go", "rest"),
		rustClientFn, newGoServer(getGoBinaries, true, transportREST), transportREST),
	ginkgo.Entry("rust → go [grpc]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("rust-go", "grpc"),
		rustClientFn, newGoServer(getGoBinaries, true, transportGRPC), transportGRPC),

	ginkgo.Entry("rust → rust [jsonrpc]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("rust-rust", "jsonrpc"),
		rustClientFn, newRustServer(getRustBinaries, true, transportJSONRPC), transportJSONRPC),
	ginkgo.Entry("rust → rust [rest]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("rust-rust", "rest"),
		rustClientFn, newRustServer(getRustBinaries, true, transportREST), transportREST),
	ginkgo.Entry("rust → rust [grpc]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("rust-rust", "grpc"),
		rustClientFn, newRustServer(getRustBinaries, true, transportGRPC), transportGRPC),

	ginkgo.Entry("rust → python [jsonrpc]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("rust-python", "jsonrpc"),
		rustClientFn, newPythonServer(getPythonAssets, true, transportJSONRPC), transportJSONRPC),
	ginkgo.Entry("rust → python [rest]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("rust-python", "rest"),
		rustClientFn, newPythonServer(getPythonAssets, true, transportREST), transportREST),
	ginkgo.Entry("rust → python [grpc]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("rust-python", "grpc"),
		rustClientFn, newPythonServer(getPythonAssets, true, transportGRPC), transportGRPC),

	ginkgo.Entry("rust → dotnet [jsonrpc]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("rust-dotnet", "jsonrpc"),
		rustClientFn, newDotNetServer(getDotNetBinaries, false, transportJSONRPC), transportJSONRPC),
	ginkgo.Entry("rust → dotnet [rest]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("rust-dotnet", "rest"),
		rustClientFn, newDotNetServer(getDotNetBinaries, false, transportREST), transportREST),

	// ── Python client ─────────────────────────────────────────────────────────
	ginkgo.Entry("python → go [jsonrpc]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("python-go", "jsonrpc"),
		pythonClientFn, newGoServer(getGoBinaries, true, transportJSONRPC), transportJSONRPC),
	ginkgo.Entry("python → go [rest]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("python-go", "rest"),
		pythonClientFn, newGoServer(getGoBinaries, true, transportREST), transportREST),
	ginkgo.Entry("python → go [grpc]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("python-go", "grpc"),
		pythonClientFn, newGoServer(getGoBinaries, true, transportGRPC), transportGRPC),

	ginkgo.Entry("python → rust [jsonrpc]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("python-rust", "jsonrpc"),
		pythonClientFn, newRustServer(getRustBinaries, true, transportJSONRPC), transportJSONRPC),
	ginkgo.Entry("python → rust [rest]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("python-rust", "rest"),
		pythonClientFn, newRustServer(getRustBinaries, true, transportREST), transportREST),
	ginkgo.Entry("python → rust [grpc]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("python-rust", "grpc"),
		pythonClientFn, newRustServer(getRustBinaries, true, transportGRPC), transportGRPC),

	ginkgo.Entry("python → python [jsonrpc]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("python-python", "jsonrpc"),
		pythonClientFn, newPythonServer(getPythonAssets, true, transportJSONRPC), transportJSONRPC),
	ginkgo.Entry("python → python [rest]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("python-python", "rest"),
		pythonClientFn, newPythonServer(getPythonAssets, true, transportREST), transportREST),
	ginkgo.Entry("python → python [grpc]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("python-python", "grpc"),
		pythonClientFn, newPythonServer(getPythonAssets, true, transportGRPC), transportGRPC),

	ginkgo.Entry("python → dotnet [jsonrpc]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("python-dotnet", "jsonrpc"),
		pythonClientFn, newDotNetServer(getDotNetBinaries, false, transportJSONRPC), transportJSONRPC),
	ginkgo.Entry("python → dotnet [rest]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("python-dotnet", "rest"),
		pythonClientFn, newDotNetServer(getDotNetBinaries, false, transportREST), transportREST),

	// ── .NET client ───────────────────────────────────────────────────────────
	ginkgo.Entry("dotnet → go [jsonrpc]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("dotnet-go", "jsonrpc"),
		dotNetClientFn, newGoServer(getGoBinaries, true, transportJSONRPC), transportJSONRPC),
	ginkgo.Entry("dotnet → go [rest]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("dotnet-go", "rest"),
		dotNetClientFn, newGoServer(getGoBinaries, true, transportREST), transportREST),

	ginkgo.Entry("dotnet → rust [jsonrpc]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("dotnet-rust", "jsonrpc"),
		dotNetClientFn, newRustServer(getRustBinaries, true, transportJSONRPC), transportJSONRPC),
	ginkgo.Entry("dotnet → rust [rest]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("dotnet-rust", "rest"),
		dotNetClientFn, newRustServer(getRustBinaries, true, transportREST), transportREST),

	ginkgo.Entry("dotnet → python [jsonrpc]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("dotnet-python", "jsonrpc"),
		dotNetClientFn, newPythonServer(getPythonAssets, true, transportJSONRPC), transportJSONRPC),
	ginkgo.Entry("dotnet → python [rest]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("dotnet-python", "rest"),
		dotNetClientFn, newPythonServer(getPythonAssets, true, transportREST), transportREST),

	ginkgo.Entry("dotnet → dotnet [jsonrpc]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("dotnet-dotnet", "jsonrpc"),
		dotNetClientFn, newDotNetServer(getDotNetBinaries, false, transportJSONRPC), transportJSONRPC),
	ginkgo.Entry("dotnet → dotnet [rest]",
		ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("dotnet-dotnet", "rest"),
		dotNetClientFn, newDotNetServer(getDotNetBinaries, false, transportREST), transportREST),
)
