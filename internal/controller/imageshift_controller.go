/*
Copyright 2025.

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

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	imageshiftv1 "github.com/spectrocloud-labs/imageshift/api/v1"
	"github.com/spectrocloud-labs/imageshift/pkg/swap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// ImageshiftReconciler reconciles a Imageshift object
type ImageshiftReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=imageshift.dev,resources=imageshifts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=imageshift.dev,resources=imageshifts/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=imageshift.dev,resources=imageshifts/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Imageshift object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.0/pkg/reconcile
func (r *ImageshiftReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	var namespaceList corev1.NamespaceList

	_ = imageshiftv1.AddToScheme(r.Scheme)

	resources := &imageshiftv1.ImageshiftList{}
	if err := r.Client.List(ctx, resources, &client.ListOptions{}); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to list Imageshift resources: %v", err)
	}

	if len(resources.Items) > 1 {
		return ctrl.Result{}, fmt.Errorf("could not determine which Imageshift Config to use")
	}

	config := resources.Items[0]

	if err := r.Client.List(ctx, &namespaceList); err != nil {
		return ctrl.Result{}, err
	}

	// Define the annotation key you want to check
	annotationKey := "imageshift.dev"

	for _, ns := range namespaceList.Items {
		if _, exists := ns.Annotations[annotationKey]; exists {

			var podList corev1.PodList
			if err := r.Client.List(ctx, &podList, client.InNamespace(ns.Namespace)); err != nil {
				continue
			}

			for _, pod := range podList.Items {
				shouldDelete := false

				for _, container := range pod.Spec.Containers {
					if swapContainer := swap.SwapImage(config, container.Image); container.Image != swapContainer {
						shouldDelete = true
					}
				}
				for _, container := range pod.Spec.InitContainers {
					if swapContainer := swap.SwapImage(config, container.Image); container.Image != swapContainer {
						shouldDelete = true
					}
				}

				if shouldDelete {
					if err := r.Client.Delete(ctx, &pod, &client.DeleteOptions{}); err != nil {
						logger.Error(err, "Failed to delete pod", "pod", pod.Name)
					}
				}
			}
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ImageshiftReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&imageshiftv1.Imageshift{}).
		Named("imageshift").
		Complete(r)
}
