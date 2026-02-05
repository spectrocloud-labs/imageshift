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

package swap

import (
	"testing"

	imageshiftv1 "github.com/spectrocloud-labs/imageshift/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSwapImage_RegistrySwap(t *testing.T) {
	tests := []struct {
		name     string
		config   imageshiftv1.Imageshift
		image    string
		expected string
	}{
		{
			name: "swap gcr.io to private registry",
			config: imageshiftv1.Imageshift{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: imageshiftv1.ImageshiftSpec{
					Default: "docker.io",
					Mappings: imageshiftv1.ImageshiftMapping{
						Swap: []imageshiftv1.ImageshiftSwap{
							{Registry: "gcr.io", Target: "my-registry.example.com"},
						},
					},
				},
			},
			image:    "gcr.io/google-containers/pause:3.1",
			expected: "my-registry.example.com/google-containers/pause:3.1",
		},
		{
			name: "swap quay.io to private registry",
			config: imageshiftv1.Imageshift{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: imageshiftv1.ImageshiftSpec{
					Default: "docker.io",
					Mappings: imageshiftv1.ImageshiftMapping{
						Swap: []imageshiftv1.ImageshiftSwap{
							{Registry: "quay.io", Target: "my-registry.example.com"},
						},
					},
				},
			},
			image:    "quay.io/coreos/etcd:v3.5.0",
			expected: "my-registry.example.com/coreos/etcd:v3.5.0",
		},
		{
			name: "no matching registry - returns empty",
			config: imageshiftv1.Imageshift{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: imageshiftv1.ImageshiftSpec{
					Default: "docker.io",
					Mappings: imageshiftv1.ImageshiftMapping{
						Swap: []imageshiftv1.ImageshiftSwap{
							{Registry: "quay.io", Target: "my-registry.example.com"},
						},
					},
				},
			},
			image:    "gcr.io/google-containers/pause:3.1",
			expected: "",
		},
		{
			name: "swap index.docker.io (normalized docker.io) to private registry",
			config: imageshiftv1.Imageshift{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: imageshiftv1.ImageshiftSpec{
					Default: "docker.io",
					Mappings: imageshiftv1.ImageshiftMapping{
						Swap: []imageshiftv1.ImageshiftSwap{
							// Use the normalized form that go-containerregistry returns
							{Registry: "index.docker.io", Target: "my-registry.example.com"},
						},
					},
				},
			},
			image:    "nginx:latest",
			expected: "my-registry.example.com/library/nginx:latest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SwapImage(tt.config, tt.image)
			if result != tt.expected {
				t.Errorf("SwapImage() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestSwapImage_ExactSwap(t *testing.T) {
	tests := []struct {
		name     string
		config   imageshiftv1.Imageshift
		image    string
		expected string
	}{
		{
			name: "exact match with full reference",
			config: imageshiftv1.Imageshift{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: imageshiftv1.ImageshiftSpec{
					Default: "docker.io",
					Mappings: imageshiftv1.ImageshiftMapping{
						ExactSwap: []imageshiftv1.ImageshiftExactSwap{
							{Reference: "gcr.io/my-project/my-image:v1.0", Target: "my-registry.example.com/my-image:v1.0"},
						},
					},
				},
			},
			image:    "gcr.io/my-project/my-image:v1.0",
			expected: "my-registry.example.com/my-image:v1.0",
		},
		{
			name: "no exact match - returns empty",
			config: imageshiftv1.Imageshift{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: imageshiftv1.ImageshiftSpec{
					Default: "docker.io",
					Mappings: imageshiftv1.ImageshiftMapping{
						ExactSwap: []imageshiftv1.ImageshiftExactSwap{
							{Reference: "gcr.io/project/image:v1.0", Target: "my-registry.example.com/image:v1.0"},
						},
					},
				},
			},
			image:    "gcr.io/project/image:v2.0",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SwapImage(tt.config, tt.image)
			if result != tt.expected {
				t.Errorf("SwapImage() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestSwapImage_RegexSwap(t *testing.T) {
	tests := []struct {
		name     string
		config   imageshiftv1.Imageshift
		image    string
		expected string
	}{
		{
			name: "regex with capture group",
			config: imageshiftv1.Imageshift{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: imageshiftv1.ImageshiftSpec{
					Default: "docker.io",
					Mappings: imageshiftv1.ImageshiftMapping{
						RegexSwap: []imageshiftv1.ImageshiftRegexSwap{
							{Expression: "^gcr.io/(.+)$", Target: "my-registry.example.com/$1"},
						},
					},
				},
			},
			image:    "gcr.io/my-project/my-image:v1.0",
			expected: "my-registry.example.com/my-project/my-image:v1.0",
		},
		{
			name: "invalid regex pattern - skipped gracefully",
			config: imageshiftv1.Imageshift{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: imageshiftv1.ImageshiftSpec{
					Default: "docker.io",
					Mappings: imageshiftv1.ImageshiftMapping{
						RegexSwap: []imageshiftv1.ImageshiftRegexSwap{
							{Expression: "[invalid(regex", Target: "my-registry.example.com"},
						},
					},
				},
			},
			image:    "gcr.io/my-project/my-image:v1.0",
			expected: "",
		},
		{
			name: "regex no match returns empty",
			config: imageshiftv1.Imageshift{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: imageshiftv1.ImageshiftSpec{
					Default: "docker.io",
					Mappings: imageshiftv1.ImageshiftMapping{
						RegexSwap: []imageshiftv1.ImageshiftRegexSwap{
							{Expression: "^quay.io/(.+)$", Target: "my-registry.example.com/$1"},
						},
					},
				},
			},
			image:    "gcr.io/my-project/my-image:v1.0",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SwapImage(tt.config, tt.image)
			if result != tt.expected {
				t.Errorf("SwapImage() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestGetCompiledRegex_Caching(t *testing.T) {
	expression := "^test-pattern-(.+)$"

	// First call should compile and cache
	re1, err := getCompiledRegex(expression)
	if err != nil {
		t.Fatalf("getCompiledRegex() error = %v", err)
	}

	// Second call should return cached version
	re2, err := getCompiledRegex(expression)
	if err != nil {
		t.Fatalf("getCompiledRegex() error = %v", err)
	}

	// Should be the same pointer (from cache)
	if re1 != re2 {
		t.Error("getCompiledRegex() should return cached regex")
	}
}

func TestGetCompiledRegex_InvalidPattern(t *testing.T) {
	_, err := getCompiledRegex("[invalid(pattern")
	if err == nil {
		t.Error("getCompiledRegex() should return error for invalid pattern")
	}
}

func TestSwapImage_CombinedRules(t *testing.T) {
	// Test that exact swap takes precedence when both swap and exact match
	config := imageshiftv1.Imageshift{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: imageshiftv1.ImageshiftSpec{
			Default: "docker.io",
			Mappings: imageshiftv1.ImageshiftMapping{
				Swap: []imageshiftv1.ImageshiftSwap{
					{Registry: "gcr.io", Target: "registry-a.example.com"},
				},
				ExactSwap: []imageshiftv1.ImageshiftExactSwap{
					{Reference: "gcr.io/project/image:v1.0", Target: "registry-b.example.com/special-image:v1.0"},
				},
			},
		},
	}

	// Exact swap should override the registry swap
	result := SwapImage(config, "gcr.io/project/image:v1.0")
	expected := "registry-b.example.com/special-image:v1.0"
	if result != expected {
		t.Errorf("SwapImage() with combined rules = %q, want %q", result, expected)
	}
}

func TestSwapImage_NoMatchReturnsEmpty(t *testing.T) {
	config := imageshiftv1.Imageshift{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: imageshiftv1.ImageshiftSpec{
			Default: "docker.io",
			Mappings: imageshiftv1.ImageshiftMapping{
				Swap: []imageshiftv1.ImageshiftSwap{
					{Registry: "quay.io", Target: "my-registry.example.com"},
				},
			},
		},
	}

	// Image from non-matching registry should return empty
	result := SwapImage(config, "gcr.io/project/image:v1.0")
	if result != "" {
		t.Errorf("SwapImage() = %q, want empty string for non-matching registry", result)
	}
}
