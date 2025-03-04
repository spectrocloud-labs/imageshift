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

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	imageshiftv1 "github.com/spectrocloud-labs/imageshift/api/v1"
)

// nolint:unused
// log is for logging in this package.
var imageshiftlog = logf.Log.WithName("imageshift-resource")

// SetupImageshiftWebhookWithManager registers the webhook for Imageshift in the manager.
func SetupImageshiftWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&imageshiftv1.Imageshift{}).
		WithValidator(&ImageshiftCustomValidator{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-spectrocloud-labs-github-com-v1-imageshift,mutating=false,failurePolicy=ignore,sideEffects=None,groups=imageshift.dev,resources=imageshifts,verbs=create,versions=v1,name=vimageshift-v1.kb.io,admissionReviewVersions=v1
//
// ImageshiftCustomValidator struct is responsible for validating the Imageshift resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type ImageshiftCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &ImageshiftCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type Imageshift.
func (v *ImageshiftCustomValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	imageshift, ok := obj.(*imageshiftv1.Imageshift)
	if !ok {
		return nil, fmt.Errorf("expected a Imageshift object but got %T", obj)
	}
	imageshiftlog.Info("Validation for Imageshift upon creation", "name", imageshift.GetName())

	var config *rest.Config

	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	scheme := runtime.NewScheme()
	_ = imageshiftv1.AddToScheme(scheme) // Register the scheme

	controllerClient, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		return nil, err
	}

	resources := &imageshiftv1.ImageshiftList{}
	if err := controllerClient.List(ctx, resources, &client.ListOptions{}); err != nil {
		return nil, fmt.Errorf("failed to list Imageshift resources: %v", err)
	}

	if len(resources.Items) > 0 {
		return nil, apierrors.NewAlreadyExists(schema.GroupResource{Group: "imageshift.dev", Resource: "imageshifts"}, "imageshift")
	}

	imageshiftlog.Info("Validation for Imageshift upon creation", "name", len(resources.Items))

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type Imageshift.
func (v *ImageshiftCustomValidator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type Imageshift.
func (v *ImageshiftCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	return nil, nil
}
