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
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	testValueStr     = "test-value"
	newValueStr      = "new-value"
	existingValueStr = "existing-value"
)

func TestPodNilAnnotationsAndLabels(t *testing.T) {
	// Test that nil annotations and labels don't cause panic
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
			// Annotations and Labels are nil
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "test-container",
					Image: "nginx:latest",
				},
			},
		},
	}

	// Verify annotations and labels are nil initially
	if pod.Annotations != nil {
		t.Error("Expected Annotations to be nil initially")
	}
	if pod.Labels != nil {
		t.Error("Expected Labels to be nil initially")
	}

	// Initialize maps like the webhook does
	if pod.Annotations == nil {
		pod.Annotations = make(map[string]string)
	}
	if pod.Labels == nil {
		pod.Labels = make(map[string]string)
	}

	// Now we should be able to write to them without panic
	pod.Annotations["test-annotation"] = testValueStr
	pod.Labels["test-label"] = testValueStr

	if pod.Annotations["test-annotation"] != testValueStr {
		t.Error("Failed to set annotation after initialization")
	}
	if pod.Labels["test-label"] != testValueStr {
		t.Error("Failed to set label after initialization")
	}
}

func TestPodWithExistingAnnotationsAndLabels(t *testing.T) {
	// Test that existing annotations and labels are preserved
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
			Annotations: map[string]string{
				"existing-annotation": existingValueStr,
			},
			Labels: map[string]string{
				"existing-label": existingValueStr,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "test-container",
					Image: "nginx:latest",
				},
			},
		},
	}

	// The nil check should not overwrite existing maps
	if pod.Annotations == nil {
		pod.Annotations = make(map[string]string)
	}
	if pod.Labels == nil {
		pod.Labels = make(map[string]string)
	}

	// Add new entries
	pod.Annotations["new-annotation"] = newValueStr
	pod.Labels["new-label"] = newValueStr

	// Verify existing values are preserved
	if pod.Annotations["existing-annotation"] != existingValueStr {
		t.Error("Existing annotation was overwritten")
	}
	if pod.Labels["existing-label"] != existingValueStr {
		t.Error("Existing label was overwritten")
	}

	// Verify new values are added
	if pod.Annotations["new-annotation"] != newValueStr {
		t.Error("Failed to add new annotation")
	}
	if pod.Labels["new-label"] != newValueStr {
		t.Error("Failed to add new label")
	}
}

func TestInitContainerAnnotationFormat(t *testing.T) {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "test-pod",
			Namespace:   "default",
			Annotations: make(map[string]string),
		},
		Spec: corev1.PodSpec{
			InitContainers: []corev1.Container{
				{
					Name:  "init-container-1",
					Image: "busybox:latest",
				},
				{
					Name:  "init-container-2",
					Image: "alpine:latest",
				},
			},
			Containers: []corev1.Container{
				{
					Name:  "main-container",
					Image: "nginx:latest",
				},
			},
		},
	}

	// Simulate the webhook annotation logic for init containers
	for i := range pod.Spec.InitContainers {
		annotationKey := pod.Spec.InitContainers[i].Name + ".initContainer.imageshift.dev/original"
		pod.Annotations[annotationKey] = pod.Spec.InitContainers[i].Image
	}

	// Verify annotations are correctly keyed by init container name
	expected := map[string]string{
		"init-container-1.initContainer.imageshift.dev/original": "busybox:latest",
		"init-container-2.initContainer.imageshift.dev/original": "alpine:latest",
	}

	for key, value := range expected {
		if pod.Annotations[key] != value {
			t.Errorf("Expected annotation %s=%s, got %s", key, value, pod.Annotations[key])
		}
	}
}

func TestContainerAnnotationFormat(t *testing.T) {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "test-pod",
			Namespace:   "default",
			Annotations: make(map[string]string),
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "container-1",
					Image: "nginx:latest",
				},
				{
					Name:  "container-2",
					Image: "redis:latest",
				},
			},
		},
	}

	// Simulate the webhook annotation logic for containers
	for i := range pod.Spec.Containers {
		annotationKey := pod.Spec.Containers[i].Name + ".container.imageshift.dev/original"
		pod.Annotations[annotationKey] = pod.Spec.Containers[i].Image
	}

	// Verify annotations are correctly keyed by container name
	expected := map[string]string{
		"container-1.container.imageshift.dev/original": "nginx:latest",
		"container-2.container.imageshift.dev/original": "redis:latest",
	}

	for key, value := range expected {
		if pod.Annotations[key] != value {
			t.Errorf("Expected annotation %s=%s, got %s", key, value, pod.Annotations[key])
		}
	}
}

func TestMutatedLabel(t *testing.T) {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
			Labels:    make(map[string]string),
		},
	}

	// Simulate setting the mutated label
	hasChanged := true
	if hasChanged {
		pod.Labels["imageshift.dev/mutated"] = mutatedLabelValue
	}

	if pod.Labels["imageshift.dev/mutated"] != mutatedLabelValue {
		t.Error("Expected mutated label to be set")
	}
}
