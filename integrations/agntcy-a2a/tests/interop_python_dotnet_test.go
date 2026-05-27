// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package tests

// This file is the Python/.NET suite wrapper. Tests cover all four client/server
// pairings of Python and .NET over the transports both SDKs support (JSON-RPC and REST).

import (
	"context"

	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Python+.NET", ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("suite-python-dotnet"), func() {
	var (
		pythonAssets pythonFixtureAssets
		dotNetAssets dotNetFixtureBinaries
	)

	ginkgo.BeforeAll(func() {
		var err error
		pythonAssets, err = buildPythonFixtureAssets()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		dotNetAssets, err = buildDotNetFixtureBinaryOnly()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

	ginkgo.AfterAll(func() { removeTempDir(dotNetAssets.tempDir) })

	ginkgo.When("a Python client calls a Python server", ginkgo.Ordered, func() {
		server := newPythonServer(func() pythonFixtureAssets { return pythonAssets }, true,
			transportJSONRPC, transportREST)

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
	})

	ginkgo.When("a Python client calls a .NET server", ginkgo.Ordered, func() {
		server := newDotNetServer(func() dotNetFixtureBinaries { return dotNetAssets }, false,
			transportJSONRPC, transportREST)

		ginkgo.BeforeAll(func() { gomega.Expect(server.start()).NotTo(gomega.HaveOccurred()) })
		ginkgo.AfterAll(func() { server.stop() })

		pythonClient := func(ctx context.Context, url string) (probeClient, error) {
			return newPythonProbeClient(func() pythonFixtureAssets { return pythonAssets }, url), nil
		}

		ginkgo.Context("over JSON-RPC", ginkgo.Ordered, ginkgo.Label("jsonrpc", "python-dotnet"), func() {
			registerBehaviors(pythonClient, server.targetFor(transportJSONRPC))
		})
		ginkgo.Context("over REST", ginkgo.Ordered, ginkgo.Label("rest", "python-dotnet"), func() {
			registerBehaviors(pythonClient, server.targetFor(transportREST))
		})
	})

	ginkgo.When("a .NET client calls a Python server", ginkgo.Ordered, func() {
		server := newPythonServer(func() pythonFixtureAssets { return pythonAssets }, true,
			transportJSONRPC, transportREST)

		ginkgo.BeforeAll(func() { gomega.Expect(server.start()).NotTo(gomega.HaveOccurred()) })
		ginkgo.AfterAll(func() { server.stop() })

		dotNetClient := func(ctx context.Context, url string) (probeClient, error) {
			return newDotNetProbeClient(func() dotNetFixtureBinaries { return dotNetAssets }, url), nil
		}

		ginkgo.Context("over JSON-RPC", ginkgo.Ordered, ginkgo.Label("jsonrpc", "dotnet-python"), func() {
			registerBehaviors(dotNetClient, server.targetFor(transportJSONRPC))
		})
		ginkgo.Context("over REST", ginkgo.Ordered, ginkgo.Label("rest", "dotnet-python"), func() {
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
