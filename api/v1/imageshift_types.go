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
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ImageshiftSpec defines the desired state of Imageshift.
type ImageshiftSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +default:value="docker.io"
	Default  string            `json:"default,omitempty"`
	Mappings ImageshiftMapping `json:"mappings,omitempty"`

	// +default:value="imageshift.dev"
	NamespaceSelector string `json:"namespaceSelector"`
}

type ImageshiftMapping struct {
	Swap      []ImageshiftSwap      `json:"swap,omitempty"`
	ExactSwap []ImageshiftExactSwap `json:"exactSwap,omitempty"`
	RegexSwap []ImageshiftRegexSwap `json:"regexSwap,omitempty"`
}

type ImageshiftSwap struct {
	Registry string `json:"registry"`
	Target   string `json:"target"`
}

type ImageshiftExactSwap struct {
	Reference string `json:"reference"`
	Target    string `json:"target"`
}

type ImageshiftRegexSwap struct {
	Expression string `json:"expression"`
	Target     string `json:"target"`
}

// ImageshiftStatus defines the observed state of Imageshift.
type ImageshiftStatus struct {
	// Conditions represent the latest available observations of the Imageshift's state
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// LastReconciled is the timestamp of the last successful reconciliation
	// +optional
	LastReconciled metav1.Time `json:"lastReconciled,omitempty"`

	// ConfigValid indicates whether the current configuration is valid
	// +optional
	ConfigValid bool `json:"configValid,omitempty"`

	// MutatedPodCount tracks the number of pods that have been mutated
	// +optional
	MutatedPodCount int64 `json:"mutatedPodCount,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster

// Imageshift is the Schema for the imageshifts API.
type Imageshift struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ImageshiftSpec   `json:"spec,omitempty"`
	Status ImageshiftStatus `json:"status,omitempty"`
}

func (i *Imageshift) SwapImage(image string) string {
	ref, _ := name.ParseReference(image, name.WithDefaultRegistry(i.Spec.Default))

	registry := ref.Context().RegistryStr()

	// if registry == default registry return image
	var newImage string

	if registry == i.Spec.Default {
		newImage = ref.Name()
	}

	for _, swap := range i.Spec.Mappings.Swap {
		if swap.Registry == registry {
			identifier := ref.Identifier()
			switch len(strings.Split(identifier, ":")) {
			case 1:
				newImage = fmt.Sprintf("%s/%s:%s", swap.Target, ref.Context().RepositoryStr(), ref.Identifier())
			case 2:
				tag := strings.Split(strings.Split(image, "@")[0], ":")[1]
				newImage = fmt.Sprintf("%s/%s:%s@%s", swap.Target, ref.Context().RepositoryStr(), tag, ref.Identifier())
			default:
				newImage = fmt.Sprintf("%s/%s", swap.Target, ref.Context().RepositoryStr())
			}
			break
		}
	}

	for _, exactSwap := range i.Spec.Mappings.ExactSwap {
		if exactSwap.Reference == ref.String() {
			newImage = exactSwap.Target
			break
		}
	}

	for _, regexSwap := range i.Spec.Mappings.RegexSwap {
		re := regexp.MustCompile(regexSwap.Expression)

		match := re.FindStringSubmatch(ref.String())
		if match != nil {
			newRepository := strings.Replace(ref.String(), match[0], "", -1)

			newImage = regexSwap.Target

			if len(match) > 1 {
				for m := 1; m < len(match); m++ {
					newImage = strings.Replace(newImage, "$"+strconv.Itoa(m), match[m], -1)
					newImage = fmt.Sprintf("%s%s", newImage, newRepository)
				}
			}
			break
		}
	}

	return newImage
}

// +kubebuilder:object:root=true

// ImageshiftList contains a list of Imageshift.
type ImageshiftList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Imageshift `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Imageshift{}, &ImageshiftList{})
}
