// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package tests

// This file is the Rust/Python suite wrapper. Tests cover all four client/server
// pairings of Rust and Python over JSON-RPC, REST, and gRPC.

import (
	"context"

	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Rust+Python", ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("suite-rust-python"), func() {
	var (
		rustAssets   fixtureBinaries
		pythonAssets pythonFixtureAssets
	)

	ginkgo.BeforeAll(func() {
		var err error
		rustAssets, err = buildRustFixtureBinaryOnly()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		pythonAssets, err = buildPythonFixtureAssets()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

	ginkgo.AfterAll(func() { removeTempDir(rustAssets.tempDir) })

	ginkgo.When("a Rust client calls a Rust server", ginkgo.Ordered, func() {
		server := newRustServer(func() fixtureBinaries { return rustAssets }, true,
			transportJSONRPC, transportREST, transportGRPC)

		ginkgo.BeforeAll(func() { gomega.Expect(server.start()).NotTo(gomega.HaveOccurred()) })
		ginkgo.AfterAll(func() { server.stop() })

		rustClient := func(ctx context.Context, url string) (probeClient, error) {
			return newRustProbeClient(func() fixtureBinaries { return rustAssets }, url), nil
		}

		ginkgo.Context("over JSON-RPC", ginkgo.Ordered, ginkgo.Label("jsonrpc", "rust-rust"), func() {
			registerBehaviors(rustClient, server.targetFor(transportJSONRPC))
		})
		ginkgo.Context("over REST", ginkgo.Ordered, ginkgo.Label("rest", "rust-rust"), func() {
			registerBehaviors(rustClient, server.targetFor(transportREST))
		})
		ginkgo.Context("over gRPC", ginkgo.Ordered, ginkgo.Label("grpc", "rust-rust"), func() {
			registerBehaviors(rustClient, server.targetFor(transportGRPC))
		})
	})

	ginkgo.When("a Rust client calls a Python server", ginkgo.Ordered, func() {
		server := newPythonServer(func() pythonFixtureAssets { return pythonAssets }, true,
			transportJSONRPC, transportREST, transportGRPC)

		ginkgo.BeforeAll(func() { gomega.Expect(server.start()).NotTo(gomega.HaveOccurred()) })
		ginkgo.AfterAll(func() { server.stop() })

		rustClient := func(ctx context.Context, url string) (probeClient, error) {
			return newRustProbeClient(func() fixtureBinaries { return rustAssets }, url), nil
		}

		ginkgo.Context("over JSON-RPC", ginkgo.Ordered, ginkgo.Label("jsonrpc", "rust-python"), func() {
			registerBehaviors(rustClient, server.targetFor(transportJSONRPC))
		})
		ginkgo.Context("over REST", ginkgo.Ordered, ginkgo.Label("rest", "rust-python"), func() {
			registerBehaviors(rustClient, server.targetFor(transportREST))
		})
		ginkgo.Context("over gRPC", ginkgo.Ordered, ginkgo.Label("grpc", "rust-python"), func() {
			registerBehaviors(rustClient, server.targetFor(transportGRPC))
		})
	})

	ginkgo.When("a Python client calls a Rust server", ginkgo.Ordered, func() {
		server := newRustServer(func() fixtureBinaries { return rustAssets }, true,
			transportJSONRPC, transportREST, transportGRPC)

		ginkgo.BeforeAll(func() { gomega.Expect(server.start()).NotTo(gomega.HaveOccurred()) })
		ginkgo.AfterAll(func() { server.stop() })

		pythonClient := func(ctx context.Context, url string) (probeClient, error) {
			return newPythonProbeClient(func() pythonFixtureAssets { return pythonAssets }, url), nil
		}

		ginkgo.Context("over JSON-RPC", ginkgo.Ordered, ginkgo.Label("jsonrpc", "python-rust"), func() {
			registerBehaviors(pythonClient, server.targetFor(transportJSONRPC))
		})
		ginkgo.Context("over REST", ginkgo.Ordered, ginkgo.Label("rest", "python-rust"), func() {
			registerBehaviors(pythonClient, server.targetFor(transportREST))
		})
		ginkgo.Context("over gRPC", ginkgo.Ordered, ginkgo.Label("grpc", "python-rust"), func() {
			registerBehaviors(pythonClient, server.targetFor(transportGRPC))
		})
	})

	ginkgo.When("a Python client calls a Python server", ginkgo.Ordered, func() {
		server := newPythonServer(func() pythonFixtureAssets { return pythonAssets }, true,
			transportJSONRPC, transportREST, transportGRPC)

		ginkgo.BeforeAll(func() { gomega.Expect(server.start()).NotTo(gomega.HaveOccurred()) })
		ginkgo.AfterAll(func() { server.stop() })

		pythonClient := func(ctx context.Context, url string) (probeClient, error) {
			return newPythonProbeClient(func() pythonFixtureAssets { return pythonAssets }, url), nil
		}

		ginkgo.Context("over JSON-RPC", ginkgo.Ordered, ginkgo.Label("jsonrpc", "python-python"), func() {
			registerBehaviors(pythonClient, server.targetFor(transportJSONRPC))
		})
		ginkgo.Context("over REST", ginkgo.Ordered, ginkgo.Label("rest", "python-python"), func() {
			registerBehaviors(pythonClient, server.targetFor(transportREST))
		})
		ginkgo.Context("over gRPC", ginkgo.Ordered, ginkgo.Label("grpc", "python-python"), func() {
			registerBehaviors(pythonClient, server.targetFor(transportGRPC))
		})
	})
})
