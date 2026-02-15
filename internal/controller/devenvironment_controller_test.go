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
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	devv1alpha1 "github.com/machinology/mach-k8s-devenv-controller/api/v1alpha1"
)

var _ = Describe("DevEnvironment Controller", func() {
	const timeout = time.Second * 10
	const interval = time.Millisecond * 250

	var devEnvName string

	BeforeEach(func() {
		devEnvName = "test-env-" + randSuffix()
	})

	Context("When creating a DevEnvironment", func() {
		BeforeEach(func() {
			ctx := context.Background()
			devEnv := &devv1alpha1.DevEnvironment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      devEnvName,
					Namespace: "default",
				},
				Spec: devv1alpha1.DevEnvironmentSpec{
					Repository: "git+https://git.machinology.local/myproject",
				},
			}
			Expect(k8sClient.Create(ctx, devEnv)).To(Succeed())
		})

		It("should create a Pod", func() {
			ctx := context.Background()
			pod := &corev1.Pod{}

			Eventually(func() error {
				return k8sClient.Get(ctx, types.NamespacedName{
					Name:      devEnvName + "-pod",
					Namespace: "default",
				}, pod)
			}, timeout, interval).Should(Succeed())
		})
	})

	Context("When deleting a DevEnvironment", func() {
		BeforeEach(func() {
			ctx := context.Background()
			devEnv := &devv1alpha1.DevEnvironment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      devEnvName,
					Namespace: "default",
				},
				Spec: devv1alpha1.DevEnvironmentSpec{
					Repository: "git+https://git.machinology.local/myproject",
				},
			}
			Expect(k8sClient.Create(ctx, devEnv)).To(Succeed())

			Eventually(func() bool {
				if err := k8sClient.Get(ctx, types.NamespacedName{
					Name:      devEnvName,
					Namespace: "default",
				}, devEnv); err != nil {
					return false
				}
				return controllerutil.ContainsFinalizer(devEnv, devEnvFinalizer)
			}, timeout, interval).Should(BeTrue())

			Expect(k8sClient.Delete(ctx, devEnv)).To(Succeed())
		})

		It("Should delete the Pod", func() {
			ctx := context.Background()

			Eventually(func() error {
				return k8sClient.Get(ctx, types.NamespacedName{
					Name:      devEnvName + "-pod",
					Namespace: "default",
				}, &corev1.Pod{})
			}, timeout, interval).Should(MatchError(ContainSubstring("not found")))
		})
	})
})

func randSuffix() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
