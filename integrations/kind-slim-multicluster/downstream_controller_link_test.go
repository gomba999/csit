// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package kindslimmulticluster

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/agntcy/csit/integrations/testutils/k8shelper"
)

const (
	labelKindSlimMulticluster = "kind-slim-multicluster"

	// After kubectl port-forward starts (or restarts), allow time for the local listen.
	northboundTCPWaitAfterPFStart = 90 * time.Second
	// At the beginning of each spec that uses slimctl via the forward, re-check the tunnel.
	northboundTCPWaitEachSpec = 30 * time.Second
)

func envOr(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}

func envEnabled(key string) bool {
	v := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	return v == "1" || v == "true" || v == "yes"
}

func slimctlPath() string {
	if p := strings.TrimSpace(os.Getenv("SLIMCTL_PATH")); p != "" {
		return p
	}
	return "slimctl"
}

func runSlimctl(ctx context.Context, controllerURL string, args ...string) ([]byte, error) {
	bin := slimctlPath()
	cmd := exec.CommandContext(ctx, bin, append([]string{"-s", controllerURL}, args...)...)
	cmd.Env = slimctlCleanEnv()
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return buf.Bytes(), fmt.Errorf("%w: %s", err, buf.String())
	}
	return buf.Bytes(), nil
}

func slimctlCleanEnv() []string {
	var out []string
	for _, e := range os.Environ() {
		k, _, ok := strings.Cut(e, "=")
		if !ok {
			out = append(out, e)
			continue
		}
		// clap global opts (agntcy/slim slimctl GlobalOpts); TLS flags break plain http:// northbound.
		if strings.HasPrefix(k, "SLIMCTL_COMMON_OPTS_") {
			continue
		}
		if k == "SLIMCTL_CONFIG" {
			continue
		}
		out = append(out, e)
	}
	out = append(out, "SLIMCTL_CONFIG=/dev/null")
	return out
}

func kubectl(ctx context.Context, kubeCtx string, args ...string) *exec.Cmd {
	full := append([]string{"--context", kubeCtx}, args...)
	return exec.CommandContext(ctx, "kubectl", full...)
}

// localTCPPortFromHTTPControllerURL returns the TCP port when the URL host is loopback
// so slimctl and kubectl port-forward can use the same local port.
func localTCPPortFromHTTPControllerURL(httpURL string) (int, bool) {
	u, err := url.Parse(httpURL)
	if err != nil || u.Host == "" {
		return 0, false
	}
	host := strings.ToLower(u.Hostname())
	if host != "localhost" && host != "127.0.0.1" && host != "::1" {
		return 0, false
	}
	raw := u.Port()
	if raw == "" {
		switch u.Scheme {
		case "https":
			return 443, true
		case "http":
			return 80, true
		default:
			return 0, false
		}
	}
	port, err := strconv.Atoi(raw)
	if err != nil {
		return 0, false
	}
	return port, true
}

// waitLoopbackNorthboundTCP waits until something accepts TCP on the controller URL host:port
// when the URL targets loopback (typical kubectl port-forward). Skips when the URL is non-loopback
// (caller manages reachability).
func waitLoopbackNorthboundTCP(ctx context.Context, httpURL string, timeout time.Duration) error {
	u, err := url.Parse(httpURL)
	if err != nil {
		return fmt.Errorf("parse controller URL %q: %w", httpURL, err)
	}
	if u.Host == "" {
		return fmt.Errorf("parse controller URL %q: empty host", httpURL)
	}
	host := strings.ToLower(u.Hostname())
	if host != "localhost" && host != "127.0.0.1" && host != "::1" {
		return nil
	}
	portStr := u.Port()
	if portStr == "" {
		switch u.Scheme {
		case "https":
			portStr = "443"
		case "http":
			portStr = "80"
		default:
			return fmt.Errorf("controller URL needs explicit port for scheme %q", u.Scheme)
		}
	}
	addr := net.JoinHostPort(host, portStr)
	deadline := time.Now().Add(timeout)
	var lastErr error
	dialer := net.Dialer{Timeout: 1 * time.Second}
	for time.Now().Before(deadline) {
		if err := ctx.Err(); err != nil {
			return err
		}
		c, err := dialer.DialContext(ctx, "tcp", addr)
		if err == nil {
			_ = c.Close()
			return nil
		}
		lastErr = err
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("timeout waiting for northbound TCP %s (port-forward listening?): last error: %w", addr, lastErr)
}

func expectNorthboundPortOpen(ctx context.Context, controllerURL string, timeout time.Duration, svc, ns string) {
	gomega.Expect(waitLoopbackNorthboundTCP(ctx, controllerURL, timeout)).To(gomega.Succeed(),
		"kubectl port-forward should accept TCP on the controller URL %q (check svc/%q in namespace %q targets northbound :50051)",
		controllerURL, svc, ns)
}

func waitPodReady(ctx context.Context, clientset kubernetes.Interface, namespace, selector string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		list, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{LabelSelector: selector})
		if err != nil {
			return err
		}
		if len(list.Items) == 0 {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		ready := 0
		for _, p := range list.Items {
			for _, c := range p.Status.Conditions {
				if c.Type == "Ready" && c.Status == "True" {
					ready++
					break
				}
			}
		}
		if ready == len(list.Items) && ready > 0 {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("timeout waiting for pods with selector %q in %s", selector, namespace)
}

func slimSVCID(ctx context.Context, kubeCtx, namespace, helmInstance string) (string, error) {
	clientset, err := k8shelper.CreateK8sClientSetForContext(kubeCtx)
	if err != nil {
		return "", err
	}
	sel := fmt.Sprintf("app.kubernetes.io/instance=%s,app.kubernetes.io/name=slim", helmInstance)
	if err := waitPodReady(ctx, clientset, namespace, sel, 3*time.Minute); err != nil {
		return "", err
	}
	list, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{LabelSelector: sel})
	if err != nil {
		return "", err
	}
	if len(list.Items) == 0 {
		return "", fmt.Errorf("no slim pods for selector %q", sel)
	}
	podObj := &list.Items[0]
	// agntcy/slim chart sets SLIM_SVC_ID from fieldRef metadata.name (pod name).
	return podObj.Name, nil
}

type portForward struct {
	cmd    *exec.Cmd
	cancel context.CancelFunc
}

func startPortForward(parent context.Context, kubeCtx, ns, svc string, localPort int) (*portForward, error) {
	ctx, cancel := context.WithCancel(parent)
	cmd := exec.CommandContext(ctx, "kubectl", "--context", kubeCtx, "-n", ns,
		"port-forward", fmt.Sprintf("svc/%s", svc), fmt.Sprintf("%d:50051", localPort))
	if err := cmd.Start(); err != nil {
		cancel()
		return nil, err
	}
	return &portForward{cmd: cmd, cancel: cancel}, nil
}

func (p *portForward) Stop() {
	if p == nil || p.cmd == nil {
		return
	}
	if p.cmd.Process != nil {
		_ = p.cmd.Process.Kill()
	}
	p.cancel()
	_ = p.cmd.Wait()
}

var _ = ginkgo.Describe("Downstream SLIM controller link", ginkgo.Ordered, ginkgo.Label(labelKindSlimMulticluster), func() {
	var (
		ctxA          = envOr("KIND_CTX_A", "kind-csit-a")
		ctxB          = envOr("KIND_CTX_B", "kind-csit-b")
		nsAdmin       = envOr("NS_ADMIN", "admin")
		ctrlRelease   = envOr("CTRL_RELEASE", "slim-control")
		nsB           = envOr("NS_B", "cluster-b")
		slimReleaseB  = envOr("SLIM_RELEASE_B", "agntcy-cluster-b")
		nsA           = envOr("NS_A", "cluster-a")
		slimReleaseA  = envOr("SLIM_RELEASE_A", "agntcy-cluster-a")
		controllerURL string
		pf            *portForward
		pfPort        int
		nodeA         string
		nodeB         string
		linkID        string
	)

	ginkgo.BeforeAll(func() {
		if !envEnabled("CSIT_KIND_SLIM_MULTICLUSTER") {
			ginkgo.Skip("set CSIT_KIND_SLIM_MULTICLUSTER=1 to run this suite")
		}
		bin := slimctlPath()
		if filepath.IsAbs(bin) {
			if _, err := os.Stat(bin); err != nil {
				ginkgo.Skip("slimctl not found at SLIMCTL_PATH: " + bin)
			}
			return
		}
		if _, err := exec.LookPath(bin); err != nil {
			ginkgo.Skip("slimctl not on PATH; set SLIMCTL_PATH to the binary")
		}
	})

	ginkgo.It("waits for downstream slim pod to be ready", func(ctx ginkgo.SpecContext) {
		clientset, err := k8shelper.CreateK8sClientSetForContext(ctxB)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		sel := fmt.Sprintf("app.kubernetes.io/instance=%s,app.kubernetes.io/name=slim", slimReleaseB)
		gomega.Expect(waitPodReady(ctx, clientset, nsB, sel, 5*time.Minute)).To(gomega.Succeed())
	})

	ginkgo.It("resolves SLIM_SVC_ID for root and downstream nodes", func(ctx ginkgo.SpecContext) {
		var err error
		nodeB, err = slimSVCID(ctx, ctxB, nsB, slimReleaseB)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		nodeA, err = slimSVCID(ctx, ctxA, nsA, slimReleaseA)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		gomega.Expect(nodeB).NotTo(gomega.BeEmpty())
		gomega.Expect(nodeA).NotTo(gomega.BeEmpty())
		gomega.Expect(nodeB).NotTo(gomega.Equal(nodeA))
	})

	ginkgo.It("starts port-forward to controller northbound", func(ctx ginkgo.SpecContext) {
		var err error
		pfPort, err = strconv.Atoi(envOr("SLIMCTL_PF_LOCAL_PORT", "51551"))
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		controllerURL = envOr("SLIM_CONTROLLER_HTTP_URL", fmt.Sprintf("http://127.0.0.1:%d", pfPort))
		pf, err = startPortForward(ctx, ctxA, nsAdmin, ctrlRelease, pfPort)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		ginkgo.DeferCleanup(func() { pf.Stop() })
		expectNorthboundPortOpen(ctx, controllerURL, northboundTCPWaitAfterPFStart, ctrlRelease, nsAdmin)
	})	

	ginkgo.Context("slimctl via northbound port-forward", func() {
		ginkgo.BeforeEach(func(ctx ginkgo.SpecContext) {
			var err error
			pf, err = startPortForward(ctx, ctxA, nsAdmin, ctrlRelease, pfPort)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			ginkgo.DeferCleanup(func() { pf.Stop() })

			expectNorthboundPortOpen(ctx, controllerURL, northboundTCPWaitEachSpec, ctrlRelease, nsAdmin)
		})

		ginkgo.It("eventually shows downstream node as Connected on the controller", func(ctx ginkgo.SpecContext) {
			gomega.Eventually(func(g gomega.Gomega) {
				out, err := runSlimctl(ctx, controllerURL, "controller", "node", "list")
				g.Expect(err).NotTo(gomega.HaveOccurred(), string(out))
				ids := ParseConnectedNodeIDs(string(out))
				g.Expect(ids).To(gomega.ContainElement(nodeB))
			}).WithTimeout(5 * time.Minute).WithPolling(3 * time.Second).Should(gomega.Succeed())
		})

		ginkgo.It("shows an APPLIED link from downstream to root dataplane", func(ctx ginkgo.SpecContext) {
			var out []byte
			var err error
			gomega.Eventually(func(g gomega.Gomega) {
				out, err = runSlimctl(ctx, controllerURL, "controller", "link", "outline",
					"-o", nodeB, "-t", nodeA)
				g.Expect(err).NotTo(gomega.HaveOccurred(), string(out))
				id, ok := ParseAppliedLinkID(string(out), nodeB, nodeA)
				g.Expect(ok).To(gomega.BeTrue(), "link outline:\n%s", string(out))
				linkID = id
			}).WithTimeout(5 * time.Minute).WithPolling(3 * time.Second).Should(gomega.Succeed())
			gomega.Expect(linkID).NotTo(gomega.BeEmpty())
		})

		ginkgo.It("keeps the same link ID after controller rollout restart", func(ctx ginkgo.SpecContext) {
			before := linkID
			restartOut, err := kubectl(ctx, ctxA, "rollout", "restart", "deployment/"+ctrlRelease, "-n", nsAdmin).CombinedOutput()
			gomega.Expect(err).NotTo(gomega.HaveOccurred(), string(restartOut))

			waitOut, err := kubectl(ctx, ctxA, "rollout", "status", "deployment/"+ctrlRelease, "-n", nsAdmin, "--timeout=5m").CombinedOutput()
			gomega.Expect(err).NotTo(gomega.HaveOccurred(), string(waitOut))

			pf.Stop()
			time.Sleep(2 * time.Second)
			var errPF error
			pf, errPF = startPortForward(ctx, ctxA, nsAdmin, ctrlRelease, pfPort)
			gomega.Expect(errPF).NotTo(gomega.HaveOccurred())
			expectNorthboundPortOpen(ctx, controllerURL, northboundTCPWaitAfterPFStart, ctrlRelease, nsAdmin)

			var after string
			gomega.Eventually(func(g gomega.Gomega) {
				out, err := runSlimctl(ctx, controllerURL, "controller", "link", "outline",
					"-o", nodeB, "-t", nodeA)
				g.Expect(err).NotTo(gomega.HaveOccurred(), string(out))
				id, ok := ParseAppliedLinkID(string(out), nodeB, nodeA)
				g.Expect(ok).To(gomega.BeTrue(), "link outline after restart:\n%s", string(out))
				after = id
			}).WithTimeout(5 * time.Minute).WithPolling(3 * time.Second).Should(gomega.Succeed())
			gomega.Expect(after).To(gomega.Equal(before))
		})
	})
})
