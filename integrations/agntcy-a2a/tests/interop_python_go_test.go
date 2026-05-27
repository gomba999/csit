// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package tests

// This file is the Python/Go suite wrapper. Tests cover all four client/server
// pairings of Python and Go over JSON-RPC, REST, and gRPC.

import (
	"context"

	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Python+Go", ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("suite-python-go"), func() {
	var (
		goAssets     fixtureBinaries
		pythonAssets pythonFixtureAssets
	)

	ginkgo.BeforeAll(func() {
		var err error
		goAssets, err = buildGoFixtureBinaryOnly()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		pythonAssets, err = buildPythonFixtureAssets()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

	ginkgo.AfterAll(func() { removeTempDir(goAssets.tempDir) })

	ginkgo.When("a Go client calls a Go server", ginkgo.Ordered, func() {
		server := newGoServer(func() fixtureBinaries { return goAssets }, true,
			transportJSONRPC, transportREST, transportGRPC)

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
		ginkgo.Context("over gRPC", ginkgo.Ordered, ginkgo.Label("grpc", "go-go"), func() {
			registerBehaviors(goClient, server.targetFor(transportGRPC))
		})
	})

	ginkgo.When("a Go client calls a Python server", ginkgo.Ordered, func() {
		server := newPythonServer(func() pythonFixtureAssets { return pythonAssets }, true,
			transportJSONRPC, transportREST, transportGRPC)

		ginkgo.BeforeAll(func() { gomega.Expect(server.start()).NotTo(gomega.HaveOccurred()) })
		ginkgo.AfterAll(func() { server.stop() })

		goClient := func(ctx context.Context, url string) (probeClient, error) {
			return newGoProbeClient(ctx, url)
		}

		ginkgo.Context("over JSON-RPC", ginkgo.Ordered, ginkgo.Label("jsonrpc", "go-python"), func() {
			registerBehaviors(goClient, server.targetFor(transportJSONRPC))
		})
		ginkgo.Context("over REST", ginkgo.Ordered, ginkgo.Label("rest", "go-python"), func() {
			registerBehaviors(goClient, server.targetFor(transportREST))
		})
		ginkgo.Context("over gRPC", ginkgo.Ordered, ginkgo.Label("grpc", "go-python"), func() {
			registerBehaviors(goClient, server.targetFor(transportGRPC))
		})
	})

	ginkgo.When("a Python client calls a Go server", ginkgo.Ordered, func() {
		server := newGoServer(func() fixtureBinaries { return goAssets }, true,
			transportJSONRPC, transportREST, transportGRPC)

		ginkgo.BeforeAll(func() { gomega.Expect(server.start()).NotTo(gomega.HaveOccurred()) })
		ginkgo.AfterAll(func() { server.stop() })

		pythonClient := func(ctx context.Context, url string) (probeClient, error) {
			return newPythonProbeClient(func() pythonFixtureAssets { return pythonAssets }, url), nil
		}

		ginkgo.Context("over JSON-RPC", ginkgo.Ordered, ginkgo.Label("jsonrpc", "python-go"), func() {
			registerBehaviors(pythonClient, server.targetFor(transportJSONRPC))
		})
		ginkgo.Context("over REST", ginkgo.Ordered, ginkgo.Label("rest", "python-go"), func() {
			registerBehaviors(pythonClient, server.targetFor(transportREST))
		})
		ginkgo.Context("over gRPC", ginkgo.Ordered, ginkgo.Label("grpc", "python-go"), func() {
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
