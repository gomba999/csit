// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"context"

	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Rust+Go", ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("suite-rust-go"), func() {
	var binaries fixtureBinaries

	ginkgo.BeforeAll(func() {
		var err error
		binaries, err = buildFixtureBinaries()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

	ginkgo.AfterAll(func() { removeTempDir(binaries.tempDir) })

	ginkgo.When("a Go client calls a Rust server", ginkgo.Ordered, func() {
		server := newRustServer(func() fixtureBinaries { return binaries }, true,
			transportJSONRPC, transportREST, transportGRPC)

		ginkgo.BeforeAll(func() { gomega.Expect(server.start()).NotTo(gomega.HaveOccurred()) })
		ginkgo.AfterAll(func() { server.stop() })

		goClient := func(ctx context.Context, url string) (probeClient, error) {
			return newGoProbeClient(ctx, url)
		}

		ginkgo.Context("over JSON-RPC", ginkgo.Ordered, ginkgo.Label("jsonrpc", "go-rust"), func() {
			registerBehaviors(goClient, server.targetFor(transportJSONRPC))
		})
		ginkgo.Context("over REST", ginkgo.Ordered, ginkgo.Label("rest", "go-rust"), func() {
			registerBehaviors(goClient, server.targetFor(transportREST))
		})
		ginkgo.Context("over gRPC", ginkgo.Ordered, ginkgo.Label("grpc", "go-rust"), func() {
			registerBehaviors(goClient, server.targetFor(transportGRPC))
		})
	})

	ginkgo.When("a Rust client calls a Go server", ginkgo.Ordered, func() {
		server := newGoServer(func() fixtureBinaries { return binaries }, true,
			transportJSONRPC, transportREST, transportGRPC)

		ginkgo.BeforeAll(func() { gomega.Expect(server.start()).NotTo(gomega.HaveOccurred()) })
		ginkgo.AfterAll(func() { server.stop() })

		rustClient := func(ctx context.Context, url string) (probeClient, error) {
			return newRustProbeClient(func() fixtureBinaries { return binaries }, url), nil
		}

		ginkgo.Context("over JSON-RPC", ginkgo.Ordered, ginkgo.Label("jsonrpc", "rust-go"), func() {
			registerBehaviors(rustClient, server.targetFor(transportJSONRPC))
		})
		ginkgo.Context("over REST", ginkgo.Ordered, ginkgo.Label("rest", "rust-go"), func() {
			registerBehaviors(rustClient, server.targetFor(transportREST))
		})
		ginkgo.Context("over gRPC", ginkgo.Ordered, ginkgo.Label("grpc", "rust-go"), func() {
			registerBehaviors(rustClient, server.targetFor(transportGRPC))
		})
	})
})
