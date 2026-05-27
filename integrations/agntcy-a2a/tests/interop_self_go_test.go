// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"context"

	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Go", ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("suite-self-go"), func() {
	var binaries fixtureBinaries

	ginkgo.BeforeAll(func() {
		var err error
		binaries, err = buildGoFixtureBinaryOnly()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

	ginkgo.AfterAll(func() { removeTempDir(binaries.tempDir) })

	ginkgo.When("a Go client calls a Go server", ginkgo.Ordered, func() {
		server := newGoServer(func() fixtureBinaries { return binaries }, true,
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
})
