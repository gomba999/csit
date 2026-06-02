// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package k8shelper

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type SecretVolume struct {
	SecretName string
	VolumeName string
	Path       string
}

type k8sHelper struct {
	clientset          kubernetes.Interface
	dynamicClient      dynamic.Interface
	name               string
	namespace          string
	imageName          string
	envVars            map[string]string
	command            []string
	args               []string
	containerPorts     []int32
	secretVolumes      []SecretVolume
	serviceAccountName string

	volumes       []corev1.Volume
	volumeMounts  []corev1.VolumeMount
	initContainer corev1.Container
}

func NewK8sHelper(name, namespace, imageName string, c kubernetes.Interface, dynamicClient dynamic.Interface) *k8sHelper {
	return &k8sHelper{
		clientset:     c,
		dynamicClient: dynamicClient,
		name:          name,
		namespace:     namespace,
		imageName:     imageName,
	}
}

func (k *k8sHelper) WithEnvVars(envVars map[string]string) *k8sHelper {
	k.envVars = envVars

	return k
}

func (k *k8sHelper) WithCommand(command []string) *k8sHelper {
	k.command = command

	return k
}

func (k *k8sHelper) WithServiceAccountName(name string) *k8sHelper {
	k.serviceAccountName = name

	return k
}

func (k *k8sHelper) WithArgs(args []string) *k8sHelper {
	k.args = args

	return k
}

func (k *k8sHelper) WithContainerPorts(ports []int32) *k8sHelper {
	k.containerPorts = ports

	return k
}

func (k *k8sHelper) WithSecretVolumes(secretVolumes []SecretVolume) *k8sHelper {
	k.secretVolumes = secretVolumes

	return k
}

func (k *k8sHelper) WithInitContainer(container corev1.Container) *k8sHelper {
	k.initContainer = container

	return k
}

func (k *k8sHelper) WithWithVolumeMounts(volumeMounts []corev1.VolumeMount) *k8sHelper {
	k.volumeMounts = volumeMounts

	return k
}

func (k *k8sHelper) WithVolumes(volumes []corev1.Volume) *k8sHelper {
	k.volumes = volumes

	return k
}

func kubeConfigPath() string {
	return filepath.Join(os.Getenv("HOME"), ".kube", "config")
}

// BuildRESTConfig returns a client config for the kubeconfig at ~/.kube/config.
// If kubeContext is non-empty it overrides the current context.
func BuildRESTConfig(kubeContext string) (*rest.Config, error) {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules.ExplicitPath = kubeConfigPath()
	var overrides *clientcmd.ConfigOverrides
	if kubeContext != "" {
		overrides = &clientcmd.ConfigOverrides{CurrentContext: kubeContext}
	}
	cc := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides)
	return cc.ClientConfig()
}

// CreateK8sClientSetForContext returns a Kubernetes clientset bound to the named kubeconfig context.
func CreateK8sClientSetForContext(kubeContext string) (*kubernetes.Clientset, error) {
	config, err := BuildRESTConfig(kubeContext)
	if err != nil {
		return nil, fmt.Errorf("unable to load kubeconfig for context %q: %w", kubeContext, err)
	}
	return kubernetes.NewForConfig(config)
}

func CreateK8sClientSet() (*kubernetes.Clientset, error) {
	config, err := BuildRESTConfig("")
	if err != nil {
		return nil, fmt.Errorf("unable to load kubeconfig %w", err)
	}
	gomega.Expect(err).NotTo(gomega.HaveOccurred(), "unable to load kubeconfig")

	return kubernetes.NewForConfig(config)
}

func CreateDynamicK8sClient() (*dynamic.DynamicClient, error) {
	config, err := BuildRESTConfig("")
	if err != nil {
		return nil, fmt.Errorf("unable to load kubeconfig %w", err)
	}
	gomega.Expect(err).NotTo(gomega.HaveOccurred(), "unable to load kubeconfig")

	// Create dynamic client for custom resources
	return dynamic.NewForConfig(config)
}
