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

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	devv1alpha1 "github.com/machinology/mach-k8s-devenv-controller/api/v1alpha1"
)

const devEnvFinalizer = "dev.machinology.dev/finalizer"

// DevEnvironmentReconciler reconciles a DevEnvironment object
type DevEnvironmentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=dev.machinology.dev,resources=devenvironments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dev.machinology.dev,resources=devenvironments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dev.machinology.dev,resources=devenvironments/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DevEnvironment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.22.4/pkg/reconcile
func (r *DevEnvironmentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = logf.FromContext(ctx)

	devEnv := &devv1alpha1.DevEnvironment{}
	if err := r.Get(ctx, req.NamespacedName, devEnv); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !devEnv.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, r.handleDeletion(ctx, devEnv)
	}

	if !controllerutil.ContainsFinalizer(devEnv, devEnvFinalizer) {
		controllerutil.AddFinalizer(devEnv, devEnvFinalizer)
		if err := r.Update(ctx, devEnv); err != nil {
			return ctrl.Result{}, err
		}
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      devEnv.Name + "-pod",
			Namespace: devEnv.Namespace,
		},
	}

	_, err := ctrl.CreateOrUpdate(ctx, r.Client, pod, func() error {
		pod.Spec = corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "devenv",
					Image:   "nixos/nix",
					Command: []string{"nix", "develop", devEnv.Spec.Repository},
				},
			},
		}
		return nil
	})

	return ctrl.Result{}, err
}

func (r *DevEnvironmentReconciler) handleDeletion(ctx context.Context, devEnv *devv1alpha1.DevEnvironment) error {
	pod := &corev1.Pod{}
	err := r.Get(ctx, client.ObjectKey{
		Name:      devEnv.Name + "-pod",
		Namespace: devEnv.Namespace,
	}, pod)

	if err == nil {
		if err := r.Delete(ctx, pod); err != nil {
			return err
		}
	}

	controllerutil.RemoveFinalizer(devEnv, devEnvFinalizer)
	return r.Update(ctx, devEnv)
}

// SetupWithManager sets up the controller with the Manager.
func (r *DevEnvironmentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&devv1alpha1.DevEnvironment{}).
		Owns(&corev1.Pod{}).
		Complete(r)
}
