---
sidebar_position: 1
---

# Introduction

ImageShift is a Kubernetes mutating admission webhook operator that automatically swaps container images from one registry to another based on configurable mapping rules.

## What is ImageShift?

When pods are created in namespaces labeled with `imageshift.dev=enabled`, the ImageShift mutating webhook intercepts the pod creation and rewrites container image references according to your configured `Imageshift` custom resource.

This happens transparently - your existing Kubernetes manifests don't need to change at all.

## Key Features

- **Zero manifest changes** - Transparently redirect images without modifying deployments
- **Namespace scoped** - Control which namespaces are affected with labels
- **Flexible mapping strategies**:
  - **Registry swap** - Redirect all images from one registry to another
  - **Exact swap** - Replace specific image references
  - **Regex swap** - Use regular expressions for complex matching patterns

## Use Cases

### Air-Gapped Environments

In disconnected networks, redirect all image pulls to your internal mirror registry:

```yaml
spec:
  mappings:
    swap:
      - registry: docker.io
        target: internal-registry.corp/dockerhub
      - registry: gcr.io
        target: internal-registry.corp/gcr
      - registry: ghcr.io
        target: internal-registry.corp/ghcr
```

### Multi-Region Deployments

Use region-specific registries to reduce latency and egress costs:

```yaml
spec:
  mappings:
    swap:
      - registry: docker.io
        target: us-west-2.mirror.example.com/dockerhub
```

### Image Mirroring and Caching

Implement transparent caching by redirecting to a pull-through cache:

```yaml
spec:
  mappings:
    swap:
      - registry: docker.io
        target: cache.example.com/docker
```

### Development and Testing

Test with local registries without changing application manifests:

```yaml
spec:
  mappings:
    swap:
      - registry: docker.io
        target: localhost:5000
```

## How It Works

1. **Install ImageShift** - Deploy the operator and webhooks to your cluster
2. **Create an Imageshift resource** - Define your image mapping rules
3. **Label namespaces** - Add `imageshift.dev=enabled` to target namespaces
4. **Deploy workloads** - Images are automatically redirected based on your rules

## Next Steps

- [Prerequisites](./installation/prerequisites) - Check requirements before installing
- [Quick Start](./installation/quick-start) - Get up and running in minutes
- [CRD Reference](./reference/crd) - Full specification of the Imageshift resource
