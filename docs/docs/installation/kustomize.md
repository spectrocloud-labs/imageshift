---
sidebar_position: 4
---

# Kustomize Installation

Install ImageShift using Kustomize for GitOps-friendly deployments.

## Quick Install

Install directly from the repository:

```bash
kubectl apply -k https://github.com/wcrum/imageshift/config/default
```

## Install from Source

Clone the repository and install:

```bash
git clone https://github.com/wcrum/imageshift.git
cd imageshift
make install  # Install CRDs
make deploy   # Deploy controller
```

## Kustomize Structure

The ImageShift repository provides the following Kustomize bases:

```
config/
├── crd/              # Custom Resource Definitions
├── rbac/             # RBAC permissions
├── manager/          # Controller deployment
├── webhook/          # Webhook configuration
├── certmanager/      # cert-manager integration
└── default/          # Complete installation
```

## Creating a Custom Overlay

Create your own overlay to customize the installation:

```
my-imageshift/
├── kustomization.yaml
└── patches/
    └── manager-resources.yaml
```

### kustomization.yaml

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: imageshift-system

resources:
  - https://github.com/wcrum/imageshift/config/default

patches:
  - path: patches/manager-resources.yaml
```

### patches/manager-resources.yaml

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: controller-manager
  namespace: imageshift-system
spec:
  replicas: 2
  template:
    spec:
      containers:
        - name: manager
          resources:
            limits:
              cpu: 500m
              memory: 256Mi
            requests:
              cpu: 100m
              memory: 128Mi
```

Apply your overlay:

```bash
kubectl apply -k my-imageshift/
```

## Common Customizations

### Change Namespace

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: custom-namespace

resources:
  - https://github.com/wcrum/imageshift/config/default
```

### Add Image Pull Secrets

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
  - https://github.com/wcrum/imageshift/config/default

patches:
  - target:
      kind: Deployment
      name: controller-manager
    patch: |
      - op: add
        path: /spec/template/spec/imagePullSecrets
        value:
          - name: my-registry-secret
```

### Use a Specific Image Version

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
  - https://github.com/wcrum/imageshift/config/default

images:
  - name: controller
    newName: ghcr.io/wcrum/imageshift
    newTag: v1.0.0
```

## Development Installation

For local development:

```bash
# Create a kind cluster
kind create cluster --name imageshift-dev

# Install CRDs only
make install

# Run the controller locally (outside the cluster)
make run
```

## Uninstall

To remove ImageShift:

```bash
# Remove controller and webhooks
kubectl delete -k https://github.com/wcrum/imageshift/config/default

# Remove CRDs (optional)
kubectl delete -k https://github.com/wcrum/imageshift/config/crd
```

## Next Steps

- [CRD Reference](../reference/crd) - Configure your Imageshift resource
- [Configuration Examples](../reference/configuration) - Real-world configuration examples
