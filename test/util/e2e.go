/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"context"
	"fmt"
	"os"
	"time"

	cmv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	"github.com/google/go-cmp/cmp/cmpopts"
	kfmpi "github.com/kubeflow/mpi-operator/pkg/apis/kubeflow/v2beta1"
	kftraining "github.com/kubeflow/training-operator/pkg/apis/kubeflow.org/v1"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	awv1beta2 "github.com/project-codeflare/appwrapper/api/v1beta2"
	rayv1 "github.com/ray-project/kuberay/ray-operator/apis/ray/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	jobset "sigs.k8s.io/jobset/api/jobset/v1alpha2"
	leaderworkersetv1 "sigs.k8s.io/lws/api/leaderworkerset/v1"
	"sigs.k8s.io/yaml"

	configapi "sigs.k8s.io/kueue/apis/config/v1beta1"
	kueuealpha "sigs.k8s.io/kueue/apis/kueue/v1alpha1"
	kueue "sigs.k8s.io/kueue/apis/kueue/v1beta1"
	visibility "sigs.k8s.io/kueue/apis/visibility/v1beta1"
	kueueclientset "sigs.k8s.io/kueue/client-go/clientset/versioned"
	visibilityv1beta1 "sigs.k8s.io/kueue/client-go/clientset/versioned/typed/visibility/v1beta1"
	utiltesting "sigs.k8s.io/kueue/pkg/util/testing"
)

const (
	// E2eTestAgnHostImageOld is the image used for testing rolling update.
	E2eTestAgnHostImageOld = "registry.k8s.io/e2e-test-images/agnhost:2.52@sha256:b173c7d0ffe3d805d49f4dfe48375169b7b8d2e1feb81783efd61eb9d08042e6"
	// E2eTestAgnHostImage is the image used for testing.
	E2eTestAgnHostImage = "registry.k8s.io/e2e-test-images/agnhost:2.53@sha256:99c6b4bb4a1e1df3f0b3752168c89358794d02258ebebc26bf21c29399011a85"
	// E2eTestSleepImageOld is the image used for testing rolling update.
	E2eTestSleepImageOld = "gcr.io/k8s-staging-perf-tests/sleep:v0.0.3@sha256:00ae8e01dd4439edfb7eb9f1960ac28eba16e952956320cce7f2ac08e3446e6b"
	// E2eTestSleepImage is the image used for testing.
	E2eTestSleepImage = "gcr.io/k8s-staging-perf-tests/sleep:v0.1.0@sha256:8d91ddf9f145b66475efda1a1b52269be542292891b5de2a7fad944052bab6ea"
	// E2eTTestCurlImage is the image used for testing with curl execution.
	E2eTTestCurlImage = "curlimages/curl:8.11.0@sha256:6324a8b41a7f9d80db93c7cf65f025411f55956c6b248037738df3bfca32410c"
)

const (
	// The environment variable for namespace where Kueue is installed
	namespaceEnvVar = "NAMESPACE"

	// The namespace where kueue is installed in opendatahub
	odhNamespace = "opendatahub"

	// The namespace where kueue is installed in rhoai
	rhoaiNamespace = "redhat-ods-applications"

	// The default namespace where kueue is installed
	kueueNamespace = "kueue-system"

	undefinedNamespace = "undefined"
)

func CreateClientUsingCluster(kContext string) (client.WithWatch, *rest.Config) {
	cfg, err := config.GetConfigWithContext(kContext)
	if err != nil {
		fmt.Printf("unable to get kubeconfig for context %q: %s", kContext, err)
		os.Exit(1)
	}
	gomega.ExpectWithOffset(1, cfg).NotTo(gomega.BeNil())

	err = kueue.AddToScheme(scheme.Scheme)
	gomega.ExpectWithOffset(1, err).NotTo(gomega.HaveOccurred())

	err = cmv1.AddToScheme(scheme.Scheme)
	gomega.ExpectWithOffset(1, err).NotTo(gomega.HaveOccurred())

	err = kueuealpha.AddToScheme(scheme.Scheme)
	gomega.ExpectWithOffset(1, err).NotTo(gomega.HaveOccurred())

	err = visibility.AddToScheme(scheme.Scheme)
	gomega.ExpectWithOffset(1, err).NotTo(gomega.HaveOccurred())

	err = jobset.AddToScheme(scheme.Scheme)
	gomega.ExpectWithOffset(1, err).NotTo(gomega.HaveOccurred())

	err = kftraining.AddToScheme(scheme.Scheme)
	gomega.ExpectWithOffset(1, err).NotTo(gomega.HaveOccurred())

	err = kfmpi.AddToScheme(scheme.Scheme)
	gomega.ExpectWithOffset(1, err).NotTo(gomega.HaveOccurred())

	err = leaderworkersetv1.AddToScheme(scheme.Scheme)
	gomega.ExpectWithOffset(1, err).NotTo(gomega.HaveOccurred())

	cfg.APIPath = "/api"
	cfg.ContentConfig.GroupVersion = &schema.GroupVersion{Group: "", Version: "v1"}
	cfg.ContentConfig.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	err = awv1beta2.AddToScheme(scheme.Scheme)
	gomega.ExpectWithOffset(1, err).NotTo(gomega.HaveOccurred())

	err = rayv1.AddToScheme(scheme.Scheme)
	gomega.ExpectWithOffset(1, err).NotTo(gomega.HaveOccurred())

	client, err := client.NewWithWatch(cfg, client.Options{Scheme: scheme.Scheme})
	gomega.ExpectWithOffset(1, err).NotTo(gomega.HaveOccurred())
	return client, cfg
}

// CreateRestClient creates a *rest.RESTClient using the provided config.
func CreateRestClient(cfg *rest.Config) *rest.RESTClient {
	restClient, err := rest.RESTClientFor(cfg)
	gomega.ExpectWithOffset(1, err).Should(gomega.Succeed())
	gomega.ExpectWithOffset(1, restClient).NotTo(gomega.BeNil())

	return restClient
}

func CreateVisibilityClient(user string) visibilityv1beta1.VisibilityV1beta1Interface {
	cfg, err := config.GetConfigWithContext("")
	if err != nil {
		fmt.Printf("unable to get kubeconfig: %s", err)
		os.Exit(1)
	}
	gomega.ExpectWithOffset(1, cfg).NotTo(gomega.BeNil())

	if user != "" {
		cfg.Impersonate.UserName = user
	}

	kueueClient, err := kueueclientset.NewForConfig(cfg)
	if err != nil {
		fmt.Printf("unable to create kueue clientset: %s", err)
		os.Exit(1)
	}
	visibilityClient := kueueClient.VisibilityV1beta1()
	return visibilityClient
}

func rolloutOperatorDeployment(ctx context.Context, k8sClient client.Client, key types.NamespacedName) {
	deployment := &appsv1.Deployment{}
	var deploymentCondition *appsv1.DeploymentCondition
	expectedDeploymentCondition := &appsv1.DeploymentCondition{
		Type:   appsv1.DeploymentProgressing,
		Status: corev1.ConditionTrue,
		Reason: "NewReplicaSetAvailable",
	}

	gomega.EventuallyWithOffset(2, func(g gomega.Gomega) {
		g.Expect(k8sClient.Get(ctx, key, deployment)).To(gomega.Succeed())
		deploymentCondition = FindDeploymentCondition(deployment, appsv1.DeploymentProgressing)
		g.Expect(deploymentCondition).To(gomega.BeComparableTo(expectedDeploymentCondition, IgnoreDeploymentConditionTimestampsAndMessage))
	}, Timeout, Interval).Should(gomega.Succeed())
	beforeUpdateTime := deploymentCondition.LastUpdateTime

	gomega.EventuallyWithOffset(2, func(g gomega.Gomega) {
		deployment.Spec.Template.ObjectMeta.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)
		g.Expect(k8sClient.Update(ctx, deployment)).To(gomega.Succeed())
	}, Timeout, Interval).Should(gomega.Succeed())

	gomega.EventuallyWithOffset(2, func(g gomega.Gomega) {
		g.Expect(k8sClient.Get(ctx, key, deployment)).To(gomega.Succeed())
		deploymentCondition := FindDeploymentCondition(deployment, appsv1.DeploymentProgressing)
		g.Expect(deploymentCondition).To(gomega.BeComparableTo(expectedDeploymentCondition, IgnoreDeploymentConditionTimestampsAndMessage))
		afterUpdateTime := deploymentCondition.LastUpdateTime
		g.Expect(afterUpdateTime).NotTo(gomega.Equal(beforeUpdateTime))
	}, StartUpTimeout, Interval).Should(gomega.Succeed())
}

func waitForOperatorAvailability(ctx context.Context, k8sClient client.Client, key types.NamespacedName) {
	deployment := &appsv1.Deployment{}
	pods := &corev1.PodList{}
	waitForAvailableStart := time.Now()
	ginkgo.By(fmt.Sprintf("Waiting for availability of deployment: %q", key))
	gomega.EventuallyWithOffset(2, func(g gomega.Gomega) error {
		g.Expect(k8sClient.Get(ctx, key, deployment)).To(gomega.Succeed())
		g.Expect(k8sClient.List(ctx, pods, client.InNamespace(GetNamespace()), client.MatchingLabels(deployment.Spec.Selector.MatchLabels))).To(gomega.Succeed())
		g.Expect(deployment.Status.Conditions).To(gomega.ContainElement(gomega.BeComparableTo(
			appsv1.DeploymentCondition{Type: appsv1.DeploymentAvailable, Status: corev1.ConditionTrue},
			cmpopts.IgnoreFields(appsv1.DeploymentCondition{}, "Reason", "Message", "LastUpdateTime", "LastTransitionTime")),
		))
		return nil
	}, StartUpTimeout, Interval).Should(gomega.Succeed())
	ginkgo.GinkgoLogr.Info("Deployment is available in the cluster", "deployment", key, "waitingTime", time.Since(waitForAvailableStart))
}

func WaitForKueueAvailability(ctx context.Context, k8sClient client.Client) {
	kcmKey := types.NamespacedName{Namespace: GetNamespace(), Name: "kueue-controller-manager"}
	waitForOperatorAvailability(ctx, k8sClient, kcmKey)
}

func WaitForAppWrapperAvailability(ctx context.Context, k8sClient client.Client) {
	awmKey := types.NamespacedName{Namespace: "appwrapper-system", Name: "appwrapper-controller-manager"}
	waitForOperatorAvailability(ctx, k8sClient, awmKey)
}

func WaitForJobSetAvailability(ctx context.Context, k8sClient client.Client) {
	_, skipJobsetAvailabilityCheck := os.LookupEnv("SKIP_JOB_SET_AVAILABILITY_CHECK")
	if !skipJobsetAvailabilityCheck {
		jcmKey := types.NamespacedName{Namespace: "jobset-system", Name: "jobset-controller-manager"}
		waitForOperatorAvailability(ctx, k8sClient, jcmKey)
	}
}

func WaitForLeaderWorkerSetAvailability(ctx context.Context, k8sClient client.Client) {
	jcmKey := types.NamespacedName{Namespace: "lws-system", Name: "lws-controller-manager"}
	waitForOperatorAvailability(ctx, k8sClient, jcmKey)
}

func WaitForKubeFlowTrainingOperatorAvailability(ctx context.Context, k8sClient client.Client) {
	_, skipTrainingOperatorAvailabilityCheck := os.LookupEnv("SKIP_TRAINING_OPERATOR_AVAILABILITY_CHECK")
	if !skipTrainingOperatorAvailabilityCheck {
		kftoKey := types.NamespacedName{Namespace: "kubeflow", Name: "training-operator"}
		waitForOperatorAvailability(ctx, k8sClient, kftoKey)
	}
}

func WaitForKubeFlowMPIOperatorAvailability(ctx context.Context, k8sClient client.Client) {
	_, skipMPIOperatorAvailabilityCheck := os.LookupEnv("SKIP_MPI_OPERATOR_AVAILABILITY_CHECK")
	if !skipMPIOperatorAvailabilityCheck {
		kftoKey := types.NamespacedName{Namespace: "mpi-operator", Name: "mpi-operator"}
		waitForOperatorAvailability(ctx, k8sClient, kftoKey)
	}
}

func WaitForKubeRayOperatorAvailability(ctx context.Context, k8sClient client.Client) {
	// TODO: use ray-system namespace instead.
	// See discussions https://github.com/kubernetes-sigs/kueue/pull/4568#discussion_r2001045775 and
	// https://github.com/ray-project/kuberay/pull/2624/files#r2001143254 for context.
	kroKey := types.NamespacedName{Namespace: "default", Name: "kuberay-operator"}
	waitForOperatorAvailability(ctx, k8sClient, kroKey)
}

func GetKueueConfiguration(ctx context.Context, k8sClient client.Client) *configapi.Configuration {
	var kueueCfg configapi.Configuration
	kcmKey := types.NamespacedName{Namespace: GetNamespace(), Name: "kueue-manager-config"}
	configMap := &corev1.ConfigMap{}

	gomega.Expect(k8sClient.Get(ctx, kcmKey, configMap)).To(gomega.Succeed())
	gomega.Expect(yaml.Unmarshal([]byte(configMap.Data["controller_manager_config.yaml"]), &kueueCfg)).To(gomega.Succeed())
	return &kueueCfg
}

func ApplyKueueConfiguration(ctx context.Context, k8sClient client.Client, kueueCfg *configapi.Configuration) {
	configMap := &corev1.ConfigMap{}
	kcmKey := types.NamespacedName{Namespace: GetNamespace(), Name: "kueue-manager-config"}
	config, err := yaml.Marshal(kueueCfg)

	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Eventually(func(g gomega.Gomega) {
		g.Expect(k8sClient.Get(ctx, kcmKey, configMap)).To(gomega.Succeed())
		configMap.Data["controller_manager_config.yaml"] = string(config)
		g.Expect(k8sClient.Update(ctx, configMap)).To(gomega.Succeed())
	}, Timeout, Interval).Should(gomega.Succeed())
}

func RestartKueueController(ctx context.Context, k8sClient client.Client) {
	kcmKey := types.NamespacedName{Namespace: GetNamespace(), Name: "kueue-controller-manager"}
	rolloutOperatorDeployment(ctx, k8sClient, kcmKey)
}

func WaitForActivePodsAndTerminate(ctx context.Context, k8sClient client.Client, restClient *rest.RESTClient, cfg *rest.Config, namespace string, activePodsCount, exitCode int, opts ...client.ListOption) {
	var activePods []corev1.Pod
	pods := corev1.PodList{}
	podListOpts := &client.ListOptions{}
	podListOpts.Namespace = namespace
	podListOpts.ApplyOptions(opts)
	gomega.EventuallyWithOffset(1, func(g gomega.Gomega) {
		g.Expect(k8sClient.List(ctx, &pods, podListOpts)).To(gomega.Succeed())
		activePods = make([]corev1.Pod, 0)
		for _, p := range pods.Items {
			if len(p.Status.PodIP) != 0 && p.Status.Phase == corev1.PodRunning {
				cmd := []string{"/bin/sh", "-c", fmt.Sprintf("curl \"http://%s:8080/readyz\"", p.Status.PodIP)}
				_, _, err := KExecute(ctx, cfg, restClient, namespace, p.Name, p.Spec.Containers[0].Name, cmd)
				g.Expect(err).ToNot(gomega.HaveOccurred())
				activePods = append(activePods, p)
			}
		}
		g.Expect(activePods).To(gomega.HaveLen(activePodsCount))
	}, LongTimeout, Interval).Should(gomega.Succeed())

	for _, p := range activePods {
		klog.V(3).Info("Terminating pod", "pod", klog.KObj(&p))
		cmd := []string{"/bin/sh", "-c", fmt.Sprintf("curl \"http://%s:8080/exit?code=%v&timeout=2s&wait=2s\"", p.Status.PodIP, exitCode)}
		_, _, err := KExecute(ctx, cfg, restClient, namespace, p.Name, p.Spec.Containers[0].Name, cmd)
		// TODO: remove the custom handling of 137 response once this is fixed in the agnhost image
		// We add the custom handling to protect in situation when the target pods completes with the expected
		// exit code but it terminates before it completes sending the response.
		if err != nil {
			gomega.ExpectWithOffset(1, err.Error()).To(gomega.ContainSubstring("137"))
		} else {
			gomega.ExpectWithOffset(1, err).ToNot(gomega.HaveOccurred())
		}
	}
}

func GetKuberayTestImage() string {
	kuberayTestImage, found := os.LookupEnv("KUBERAY_RAY_IMAGE")
	gomega.Expect(found).To(gomega.BeTrue())
	return kuberayTestImage
}

func CreateNamespaceWithLog(ctx context.Context, k8sClient client.Client, nsName string) *corev1.Namespace {
	ginkgo.GinkgoHelper()
	ns := utiltesting.MakeNamespace(nsName)
	gomega.Expect(k8sClient.Create(ctx, ns)).To(gomega.Succeed())
	ginkgo.GinkgoLogr.Info(fmt.Sprintf("Created namespace: %s", ns.Name))
	return ns
}

func CreateNamespaceFromPrefixWithLog(ctx context.Context, k8sClient client.Client, nsPrefix string) *corev1.Namespace {
	ginkgo.GinkgoHelper()
	ns := utiltesting.MakeNamespaceWithGenerateName(nsPrefix)
	gomega.Expect(k8sClient.Create(ctx, ns)).To(gomega.Succeed())
	ginkgo.GinkgoLogr.Info(fmt.Sprintf("Created namespace: %s", ns.Name))
	return ns
}
func GetNamespace() string {
	namespace, ok := os.LookupEnv(namespaceEnvVar)
	if !ok {
		fmt.Printf("Expected environment variable %s is unset, please use this environment variable to specify in which namespace Kueue is installed", namespaceEnvVar)
		os.Exit(1)
	}
	switch namespace {
	case "opendatahub":
		return odhNamespace
	case "redhat-ods-applications":
		return rhoaiNamespace
	case "kueue-system":
		return kueueNamespace
	default:
		fmt.Printf("Expected environment variable %s contains an incorrect value", namespaceEnvVar)
		os.Exit(1)
		return undefinedNamespace
	}
}
