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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	spectrocloudlabsgithubcomv1 "github.com/spectrocloud-labs/imageshift/api/v1"
)

var _ = Describe("Imageshift Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-resource"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "default",
		}
		imageshift := &spectrocloudlabsgithubcomv1.Imageshift{}

		BeforeEach(func() {
			By("creating the custom resource for the Kind Imageshift")
			err := k8sClient.Get(ctx, typeNamespacedName, imageshift)
			if err != nil && errors.IsNotFound(err) {
				resource := &spectrocloudlabsgithubcomv1.Imageshift{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: "default",
					},
					Spec: spectrocloudlabsgithubcomv1.ImageshiftSpec{
						Default:           "docker.io",
						NamespaceSelector: "imageshift.dev",
						Mappings: spectrocloudlabsgithubcomv1.ImageshiftMapping{
							Swap: []spectrocloudlabsgithubcomv1.ImageshiftSwap{
								{Registry: "docker.io", Target: "my-registry.example.com"},
							},
						},
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			resource := &spectrocloudlabsgithubcomv1.Imageshift{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			if err == nil {
				By("Cleanup the specific resource instance Imageshift")
				Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
			}
		})

		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &ImageshiftReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("When no Imageshift resources exist", func() {
		ctx := context.Background()

		It("should handle empty resource list gracefully", func() {
			By("Ensuring no Imageshift resources exist")
			// First, clean up any existing resources
			resources := &spectrocloudlabsgithubcomv1.ImageshiftList{}
			Expect(k8sClient.List(ctx, resources)).To(Succeed())
			for _, r := range resources.Items {
				Expect(k8sClient.Delete(ctx, &r)).To(Succeed())
			}

			By("Reconciling with no resources")
			controllerReconciler := &ImageshiftReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			result, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      "non-existent",
					Namespace: "default",
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Requeue).To(BeFalse())
		})
	})

	Context("When multiple Imageshift resources exist", func() {
		ctx := context.Background()

		It("should return error for multiple configs", func() {
			By("Creating first Imageshift resource")
			resource1 := &spectrocloudlabsgithubcomv1.Imageshift{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "config-1",
					Namespace: "default",
				},
				Spec: spectrocloudlabsgithubcomv1.ImageshiftSpec{
					Default:           "docker.io",
					NamespaceSelector: "imageshift.dev",
				},
			}
			Expect(k8sClient.Create(ctx, resource1)).To(Succeed())

			By("Creating second Imageshift resource")
			resource2 := &spectrocloudlabsgithubcomv1.Imageshift{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "config-2",
					Namespace: "default",
				},
				Spec: spectrocloudlabsgithubcomv1.ImageshiftSpec{
					Default:           "docker.io",
					NamespaceSelector: "imageshift.dev",
				},
			}
			Expect(k8sClient.Create(ctx, resource2)).To(Succeed())

			By("Reconciling with multiple resources")
			controllerReconciler := &ImageshiftReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      "config-1",
					Namespace: "default",
				},
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("could not determine which Imageshift Config to use"))

			By("Cleaning up resources")
			Expect(k8sClient.Delete(ctx, resource1)).To(Succeed())
			Expect(k8sClient.Delete(ctx, resource2)).To(Succeed())
		})
	})

	Context("When reconciling with labeled namespace", func() {
		const resourceName = "test-imageshift"
		const namespaceName = "test-imageshift-ns"

		ctx := context.Background()

		BeforeEach(func() {
			By("Creating labeled namespace")
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespaceName,
					Annotations: map[string]string{
						"imageshift.dev": "enabled",
					},
				},
			}
			err := k8sClient.Create(ctx, ns)
			if err != nil && !errors.IsAlreadyExists(err) {
				Expect(err).NotTo(HaveOccurred())
			}

			By("Creating Imageshift resource")
			resource := &spectrocloudlabsgithubcomv1.Imageshift{
				ObjectMeta: metav1.ObjectMeta{
					Name:      resourceName,
					Namespace: "default",
				},
				Spec: spectrocloudlabsgithubcomv1.ImageshiftSpec{
					Default:           "docker.io",
					NamespaceSelector: "imageshift.dev",
					Mappings: spectrocloudlabsgithubcomv1.ImageshiftMapping{
						Swap: []spectrocloudlabsgithubcomv1.ImageshiftSwap{
							{Registry: "docker.io", Target: "my-registry.example.com"},
						},
					},
				},
			}
			err = k8sClient.Create(ctx, resource)
			if err != nil && !errors.IsAlreadyExists(err) {
				Expect(err).NotTo(HaveOccurred())
			}
		})

		AfterEach(func() {
			By("Cleaning up Imageshift resource")
			resource := &spectrocloudlabsgithubcomv1.Imageshift{}
			err := k8sClient.Get(ctx, types.NamespacedName{Name: resourceName, Namespace: "default"}, resource)
			if err == nil {
				Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
			}

			By("Cleaning up namespace")
			ns := &corev1.Namespace{}
			err = k8sClient.Get(ctx, types.NamespacedName{Name: namespaceName}, ns)
			if err == nil {
				Expect(k8sClient.Delete(ctx, ns)).To(Succeed())
			}
		})

		It("should process namespaces with annotation", func() {
			By("Reconciling the resource")
			controllerReconciler := &ImageshiftReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      resourceName,
					Namespace: "default",
				},
			})
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
