// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package tests

// This file is the Rust/.NET suite wrapper. Tests cover all four client/server
// pairings of Rust and .NET over the transports both SDKs support (JSON-RPC and REST).

import (
	"context"

	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Rust+.NET", ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("suite-rust-dotnet"), func() {
	var binaries rustDotNetFixtureBinaries

	ginkgo.BeforeAll(func() {
		var err error
		binaries, err = buildRustDotNetFixtureBinaries()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

	ginkgo.AfterAll(func() { removeTempDir(binaries.tempDir) })

	ginkgo.When("a Rust client calls a Rust server", ginkgo.Ordered, func() {
		server := newRustServer(func() fixtureBinaries { return binaries.rustAssets() }, true,
			transportJSONRPC, transportREST)

		ginkgo.BeforeAll(func() { gomega.Expect(server.start()).NotTo(gomega.HaveOccurred()) })
		ginkgo.AfterAll(func() { server.stop() })

		rustClient := func(ctx context.Context, url string) (probeClient, error) {
			return newRustProbeClient(func() fixtureBinaries { return binaries.rustAssets() }, url), nil
		}

		ginkgo.Context("over JSON-RPC", ginkgo.Ordered, ginkgo.Label("jsonrpc", "rust-rust"), func() {
			registerBehaviors(rustClient, server.targetFor(transportJSONRPC))
		})
		ginkgo.Context("over REST", ginkgo.Ordered, ginkgo.Label("rest", "rust-rust"), func() {
			registerBehaviors(rustClient, server.targetFor(transportREST))
		})
	})

	ginkgo.When("a Rust client calls a .NET server", ginkgo.Ordered, func() {
		server := newDotNetServer(func() dotNetFixtureBinaries { return binaries.dotNetAssets() }, false,
			transportJSONRPC, transportREST)

		ginkgo.BeforeAll(func() { gomega.Expect(server.start()).NotTo(gomega.HaveOccurred()) })
		ginkgo.AfterAll(func() { server.stop() })

		rustClient := func(ctx context.Context, url string) (probeClient, error) {
			return newRustProbeClient(func() fixtureBinaries { return binaries.rustAssets() }, url), nil
		}

		ginkgo.Context("over JSON-RPC", ginkgo.Ordered, ginkgo.Label("jsonrpc", "rust-dotnet"), func() {
			registerBehaviors(rustClient, server.targetFor(transportJSONRPC))
		})
		ginkgo.Context("over REST", ginkgo.Ordered, ginkgo.Label("rest", "rust-dotnet"), func() {
			registerBehaviors(rustClient, server.targetFor(transportREST))
		})
	})

	ginkgo.When("a .NET client calls a Rust server", ginkgo.Ordered, func() {
		server := newRustServer(func() fixtureBinaries { return binaries.rustAssets() }, true,
			transportJSONRPC, transportREST)

		ginkgo.BeforeAll(func() { gomega.Expect(server.start()).NotTo(gomega.HaveOccurred()) })
		ginkgo.AfterAll(func() { server.stop() })

		dotNetClient := func(ctx context.Context, url string) (probeClient, error) {
			return newDotNetProbeClient(func() dotNetFixtureBinaries { return binaries.dotNetAssets() }, url), nil
		}

		ginkgo.Context("over JSON-RPC", ginkgo.Ordered, ginkgo.Label("jsonrpc", "dotnet-rust"), func() {
			registerBehaviors(dotNetClient, server.targetFor(transportJSONRPC))
		})
		ginkgo.Context("over REST", ginkgo.Ordered, ginkgo.Label("rest", "dotnet-rust"), func() {
			registerBehaviors(dotNetClient, server.targetFor(transportREST))
		})
	})

	ginkgo.When("a .NET client calls a .NET server", ginkgo.Ordered, func() {
		server := newDotNetServer(func() dotNetFixtureBinaries { return binaries.dotNetAssets() }, false,
			transportJSONRPC, transportREST)

		ginkgo.BeforeAll(func() { gomega.Expect(server.start()).NotTo(gomega.HaveOccurred()) })
		ginkgo.AfterAll(func() { server.stop() })

		dotNetClient := func(ctx context.Context, url string) (probeClient, error) {
			return newDotNetProbeClient(func() dotNetFixtureBinaries { return binaries.dotNetAssets() }, url), nil
		}

		ginkgo.Context("over JSON-RPC", ginkgo.Ordered, ginkgo.Label("jsonrpc", "dotnet-dotnet"), func() {
			registerBehaviors(dotNetClient, server.targetFor(transportJSONRPC))
		})
		ginkgo.Context("over REST", ginkgo.Ordered, ginkgo.Label("rest", "dotnet-dotnet"), func() {
			registerBehaviors(dotNetClient, server.targetFor(transportREST))
		})
	})
})
