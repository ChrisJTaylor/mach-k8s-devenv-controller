/*
Copyright 2026.

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

package controller

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	devv1alpha1 "github.com/machinology/mach-k8s-devenv-controller/api/v1alpha1"
)

var _ = Describe("DevEnvironmentConfig Controller", func() {
	const timeout = time.Second * 10
	const interval = time.Millisecond * 250

	var configName string

	BeforeEach(func() {
		configName = "test-config" + randSuffix()
	})

	Context("When creating a DevEnvironmentConfig", func() {
		It("Should exist in the cluster", func() {
			ctx := context.Background()

			config := &devv1alpha1.DevEnvironmentConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configName,
					Namespace: "default",
				},
				Spec: devv1alpha1.DevEnvironmentConfigSpec{
					UserEnvironment: "git+https://github.com/machinology/nixvim-config",
					Tools:           []string{"git", "cocogitto"},
				},
			}

			Expect(k8sClient.Create(ctx, config)).To(Succeed())

			retrieved := &devv1alpha1.DevEnvironmentConfig{}

			Eventually(func() error {
				return k8sClient.Get(ctx, types.NamespacedName{
					Name:      configName,
					Namespace: "default",
				}, retrieved)
			}, timeout, interval).Should(Succeed())
		})
	})

})
