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

package certmanager

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/kueue/apis/config/v1beta1"
	"sigs.k8s.io/kueue/test/util"
)

var (
	k8sClient       client.WithWatch
	ctx             context.Context
	defaultKueueCfg *v1beta1.Configuration
)

func TestAPIs(t *testing.T) {
	suiteName := "End To End Cert Manager Integration Suite"
	if ver, found := os.LookupEnv("E2E_KIND_VERSION"); found {
		suiteName = fmt.Sprintf("%s: %s", suiteName, ver)
	}
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t,
		suiteName,
	)
}

var _ = ginkgo.BeforeSuite(func() {
	util.SetupLogger()

	k8sClient, _ = util.CreateClientUsingCluster("")
	ctx = context.Background()

	waitForAvailableStart := time.Now()
	util.WaitForKueueAvailability(ctx, k8sClient)
	ginkgo.GinkgoLogr.Info("Kueue operator is available in the cluster", "waitingTime", time.Since(waitForAvailableStart))
	defaultKueueCfg = util.GetKueueConfiguration(ctx, k8sClient)
})

var _ = ginkgo.AfterSuite(func() {
	util.ApplyKueueConfiguration(ctx, k8sClient, defaultKueueCfg)
	util.RestartKueueController(ctx, k8sClient)
	ginkgo.GinkgoLogr.Info("Default Kueue configuration restored")
})
