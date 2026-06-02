// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package kindslimmulticluster

import (
	"testing"

	"github.com/onsi/gomega"
)

func TestParseConnectedNodeIDs(t *testing.T) {
	g := gomega.NewWithT(t)
	out := `2 node(s) found
Node ID: aaa status: Connected
  Connection details:
  - Endpoint: https://example:50052
Node ID: bbb status: NotConnected
Node ID: ccc status: Connected
`
	ids := ParseConnectedNodeIDs(out)
	g.Expect(ids).To(gomega.Equal([]string{"aaa", "ccc"}))
}

func TestParseAppliedLinkID(t *testing.T) {
	g := gomega.NewWithT(t)
	out := `Outline links (origin:[x] target:[y])
Number of links: 1

  LINK_ID    SOURCE    DEST_NODE    DEST_ENDPOINT    STATUS    STATUS_MSG    DELETED    LAST_UPDATED
  ---------- --------- ------------ ---------------- --------- ------------- ---------- --------------
  link-uuid  src-node  dst-node     0.0.0.0:46357    APPLIED   -             No         2025-05-01T12:00:00Z
`
	id, ok := ParseAppliedLinkID(out, "src-node", "dst-node")
	g.Expect(ok).To(gomega.BeTrue())
	g.Expect(id).To(gomega.Equal("link-uuid"))
}

func TestLocalTCPPortFromHTTPControllerURL(t *testing.T) {
	g := gomega.NewWithT(t)
	p, ok := localTCPPortFromHTTPControllerURL("http://localhost:50051")
	g.Expect(ok).To(gomega.BeTrue())
	g.Expect(p).To(gomega.Equal(50051))
	_, ok = localTCPPortFromHTTPControllerURL("http://control.nb.cluster-a.csit.test:80")
	g.Expect(ok).To(gomega.BeFalse())
}
