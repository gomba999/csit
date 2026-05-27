// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"context"

	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe(".NET", ginkgo.Ordered, ginkgo.ContinueOnFailure, ginkgo.Label("suite-self-dotnet"), func() {
	var assets dotNetFixtureBinaries

	ginkgo.BeforeAll(func() {
		var err error
		assets, err = buildDotNetFixtureBinaryOnly()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

	ginkgo.AfterAll(func() { removeTempDir(assets.tempDir) })

	ginkgo.When("a .NET client calls a .NET server", ginkgo.Ordered, func() {
		server := newDotNetServer(func() dotNetFixtureBinaries { return assets }, false,
			transportJSONRPC, transportREST)

		ginkgo.BeforeAll(func() { gomega.Expect(server.start()).NotTo(gomega.HaveOccurred()) })
		ginkgo.AfterAll(func() { server.stop() })

		dotNetClient := func(ctx context.Context, url string) (probeClient, error) {
			return newDotNetProbeClient(func() dotNetFixtureBinaries { return assets }, url), nil
		}

		ginkgo.Context("over JSON-RPC", ginkgo.Ordered, ginkgo.Label("jsonrpc", "dotnet-dotnet"), func() {
			registerBehaviors(dotNetClient, server.targetFor(transportJSONRPC))
		})
		ginkgo.Context("over REST", ginkgo.Ordered, ginkgo.Label("rest", "dotnet-dotnet"), func() {
			registerBehaviors(dotNetClient, server.targetFor(transportREST))
		})
	})
})
