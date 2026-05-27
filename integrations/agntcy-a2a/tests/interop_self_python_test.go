// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"context"

	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Python", ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("suite-self-python"), func() {
	var assets pythonFixtureAssets

	ginkgo.BeforeAll(func() {
		var err error
		assets, err = buildPythonFixtureAssets()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

	ginkgo.When("a Python client calls a Python server", ginkgo.Ordered, func() {
		server := newPythonServer(func() pythonFixtureAssets { return assets }, true,
			transportJSONRPC, transportREST, transportGRPC)

		ginkgo.BeforeAll(func() { gomega.Expect(server.start()).NotTo(gomega.HaveOccurred()) })
		ginkgo.AfterAll(func() { server.stop() })

		pythonClient := func(ctx context.Context, url string) (probeClient, error) {
			return newPythonProbeClient(func() pythonFixtureAssets { return assets }, url), nil
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
