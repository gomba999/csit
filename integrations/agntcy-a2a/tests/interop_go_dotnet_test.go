// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package tests

// This file is the Go/.NET suite wrapper. Tests cover all four client/server
// pairings of Go and .NET over the transports both SDKs support (JSON-RPC and REST).

import (
	"context"

	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Go+.NET", ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("suite-go-dotnet"), func() {
	var (
		goAssets     fixtureBinaries
		dotNetAssets dotNetFixtureBinaries
	)

	ginkgo.BeforeAll(func() {
		var err error
		goAssets, err = buildGoFixtureBinaryOnly()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		dotNetAssets, err = buildDotNetFixtureBinaryOnly()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

	ginkgo.AfterAll(func() {
		removeTempDir(goAssets.tempDir)
		removeTempDir(dotNetAssets.tempDir)
	})

	ginkgo.When("a Go client calls a Go server", ginkgo.Ordered, func() {
		server := newGoServer(func() fixtureBinaries { return goAssets }, true,
			transportJSONRPC, transportREST)

		ginkgo.BeforeAll(func() { gomega.Expect(server.start()).NotTo(gomega.HaveOccurred()) })
		ginkgo.AfterAll(func() { server.stop() })

		goClient := func(ctx context.Context, url string) (probeClient, error) {
			return newGoProbeClient(ctx, url)
		}

		ginkgo.Context("over JSON-RPC", ginkgo.Ordered, ginkgo.Label("jsonrpc", "go-go"), func() {
			registerBehaviors(goClient, server.targetFor(transportJSONRPC))
		})
		ginkgo.Context("over REST", ginkgo.Ordered, ginkgo.Label("rest", "go-go"), func() {
			registerBehaviors(goClient, server.targetFor(transportREST))
		})
	})

	ginkgo.When("a Go client calls a .NET server", ginkgo.Ordered, func() {
		server := newDotNetServer(func() dotNetFixtureBinaries { return dotNetAssets }, false,
			transportJSONRPC, transportREST)

		ginkgo.BeforeAll(func() { gomega.Expect(server.start()).NotTo(gomega.HaveOccurred()) })
		ginkgo.AfterAll(func() { server.stop() })

		goClient := func(ctx context.Context, url string) (probeClient, error) {
			return newGoProbeClient(ctx, url)
		}

		ginkgo.Context("over JSON-RPC", ginkgo.Ordered, ginkgo.Label("jsonrpc", "go-dotnet"), func() {
			registerBehaviors(goClient, server.targetFor(transportJSONRPC))
		})
		ginkgo.Context("over REST", ginkgo.Ordered, ginkgo.Label("rest", "go-dotnet"), func() {
			registerBehaviors(goClient, server.targetFor(transportREST))
		})
	})

	ginkgo.When("a .NET client calls a Go server", ginkgo.Ordered, func() {
		server := newGoServer(func() fixtureBinaries { return goAssets }, true,
			transportJSONRPC, transportREST)

		ginkgo.BeforeAll(func() { gomega.Expect(server.start()).NotTo(gomega.HaveOccurred()) })
		ginkgo.AfterAll(func() { server.stop() })

		dotNetClient := func(ctx context.Context, url string) (probeClient, error) {
			return newDotNetProbeClient(func() dotNetFixtureBinaries { return dotNetAssets }, url), nil
		}

		ginkgo.Context("over JSON-RPC", ginkgo.Ordered, ginkgo.Label("jsonrpc", "dotnet-go"), func() {
			registerBehaviors(dotNetClient, server.targetFor(transportJSONRPC))
		})
		ginkgo.Context("over REST", ginkgo.Ordered, ginkgo.Label("rest", "dotnet-go"), func() {
			registerBehaviors(dotNetClient, server.targetFor(transportREST))
		})
	})

	ginkgo.When("a .NET client calls a .NET server", ginkgo.Ordered, func() {
		server := newDotNetServer(func() dotNetFixtureBinaries { return dotNetAssets }, false,
			transportJSONRPC, transportREST)

		ginkgo.BeforeAll(func() { gomega.Expect(server.start()).NotTo(gomega.HaveOccurred()) })
		ginkgo.AfterAll(func() { server.stop() })

		dotNetClient := func(ctx context.Context, url string) (probeClient, error) {
			return newDotNetProbeClient(func() dotNetFixtureBinaries { return dotNetAssets }, url), nil
		}

		ginkgo.Context("over JSON-RPC", ginkgo.Ordered, ginkgo.Label("jsonrpc", "dotnet-dotnet"), func() {
			registerBehaviors(dotNetClient, server.targetFor(transportJSONRPC))
		})
		ginkgo.Context("over REST", ginkgo.Ordered, ginkgo.Label("rest", "dotnet-dotnet"), func() {
			registerBehaviors(dotNetClient, server.targetFor(transportREST))
		})
	})
})
