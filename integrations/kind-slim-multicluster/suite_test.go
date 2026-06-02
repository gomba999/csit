// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

// Integration Ginkgo specs are skipped unless CSIT_KIND_SLIM_MULTICLUSTER=1 and slimctl is
// resolvable (see downstream_controller_link_test.go BeforeAll). From the repo use:
//
//	task test:integrations:kind-slim-multicluster
//
// in kind-slim-multi-host/. The IDE debug console shows "SKIP kind-slim-multicluster: ..." on stderr when gated.
package kindslimmulticluster

import (
	"testing"

	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestKindSlimMulticluster(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "KinD multicluster SLIM")
}
