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

package v1

import (
	"context"
	"fmt"

	imageshiftv1 "github.com/spectrocloud-labs/imageshift/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	swap "github.com/spectrocloud-labs/imageshift/pkg/swap"
)

// nolint:unused
// log is for logging in this package.
var podlog = logf.Log.WithName("pod-resource")

// SetupPodWebhookWithManager registers the webhook for Pod in the manager.
func SetupPodWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&corev1.Pod{}).
		WithDefaulter(&PodCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate--v1-pod,mutating=true,failurePolicy=ignore,sideEffects=None,groups="",resources=pods,verbs=create;update,versions=v1,name=mpod-v1.kb.io,admissionReviewVersions=v1

// PodCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind Pod when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type PodCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &PodCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind Pod.
func (d *PodCustomDefaulter) Default(ctx context.Context, obj runtime.Object) error {
	pod, ok := obj.(*corev1.Pod)

	if !ok {
		return fmt.Errorf("expected an Pod object but got %T", obj)
	}
	podlog.Info("Defaulting for Pod", "name", pod.GetName())

	// read only one imageshift config

	var config *rest.Config

	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	scheme := runtime.NewScheme()
	_ = imageshiftv1.AddToScheme(scheme) // Register the scheme

	controllerClient, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		return err
	}

	// TODO: better way to do this?
	// should the imageshift config be forced named?

	resources := &imageshiftv1.ImageshiftList{}
	if err := controllerClient.List(ctx, resources, &client.ListOptions{}); err != nil {
		return fmt.Errorf("failed to list Imageshift resources: %v", err)
	}

	mapping := resources.Items[0]

	hasChanged := false

	for i, container := range pod.Spec.Containers {
		img := swap.SwapImage(mapping, container.Image)

		if img != "" {
			hasChanged = true

			annotation := fmt.Sprintf("%s.container.imageshift.dev/original", pod.Spec.Containers[i].Name)
			pod.Annotations[annotation] = pod.Spec.Containers[i].Image

			pod.Spec.Containers[i].Image = img

			podlog.Info("Patched Container", "pod", container.Name, "reference", img)
		}
	}

	for i, container := range pod.Spec.InitContainers {
		img := swap.SwapImage(mapping, container.Image)

		if img != "" {
			hasChanged = true
			annotation := fmt.Sprintf("%s.initContainer.imageshift.dev/original", pod.Spec.Containers[i].Name)
			pod.Annotations[annotation] = pod.Spec.Containers[i].Image
			pod.Spec.InitContainers[i].Image = img
			podlog.Info("Patched initContainer", "pod", container.Name, "reference", img)
		}
	}

	if hasChanged {
		pod.Labels["imageshift.dev/mutated"] = "true"
	}

	return nil
}
