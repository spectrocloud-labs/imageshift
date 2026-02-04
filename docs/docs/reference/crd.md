---
sidebar_position: 1
---

# CRD Reference

Complete reference for the Imageshift Custom Resource Definition.

## Overview

| Property | Value |
|----------|-------|
| API Version | `imageshift.dev/v1` |
| Kind | `Imageshift` |
| Scope | Cluster |

The Imageshift resource defines image mapping rules that are applied to pods created in namespaces with the appropriate label.

## Spec

### Top-Level Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `spec.default` | string | `"docker.io"` | The default registry for images without an explicit registry |
| `spec.namespaceSelector` | string | `"imageshift.dev"` | The label key used to identify namespaces for image swapping |
| `spec.mappings` | object | | Container for all mapping rules |

### Mappings

The `spec.mappings` object contains three types of mapping rules:

| Field | Type | Description |
|-------|------|-------------|
| `mappings.swap` | array | Registry-level swaps - redirect all images from one registry to another |
| `mappings.exactSwap` | array | Exact image swaps - replace specific image references |
| `mappings.regexSwap` | array | Regex swaps - use regular expressions for complex matching |

#### swap

Registry-level swaps redirect all images from a source registry to a target registry.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `registry` | string | Yes | Source registry to match (e.g., `docker.io`, `gcr.io`) |
| `target` | string | Yes | Target registry to redirect to |

Example:

```yaml
spec:
  mappings:
    swap:
      - registry: docker.io
        target: registry.internal.example.com/dockerhub
      - registry: gcr.io
        target: registry.internal.example.com/gcr
```

With this configuration:
- `nginx:latest` becomes `registry.internal.example.com/dockerhub/library/nginx:latest`
- `gcr.io/google-containers/pause:3.2` becomes `registry.internal.example.com/gcr/google-containers/pause:3.2`

#### exactSwap

Exact swaps replace specific image references with another image.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `reference` | string | Yes | Exact image reference to match |
| `target` | string | Yes | Replacement image reference |

Example:

```yaml
spec:
  mappings:
    exactSwap:
      - reference: nginx:latest
        target: internal-registry.example.com/approved/nginx:1.25
      - reference: redis:7
        target: internal-registry.example.com/approved/redis:7.2.3
```

With this configuration:
- `nginx:latest` becomes `internal-registry.example.com/approved/nginx:1.25`
- `redis:7` becomes `internal-registry.example.com/approved/redis:7.2.3`

#### regexSwap

Regex swaps use regular expressions for complex matching and replacement patterns.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `expression` | string | Yes | Regular expression pattern to match |
| `target` | string | Yes | Replacement string (supports capture groups) |

Example:

```yaml
spec:
  mappings:
    regexSwap:
      - expression: "^(\\d+)\\.dkr\\.ecr\\.(.*)\\.amazonaws\\.com/(.*)$"
        target: "ecr-mirror.internal.example.com/$1/$2/$3"
      - expression: "^gcr\\.io/([^/]+)/(.*)$"
        target: "gcr-mirror.internal.example.com/$1/$2"
```

With this configuration:
- `123456789012.dkr.ecr.us-west-2.amazonaws.com/my-app:v1` becomes `ecr-mirror.internal.example.com/123456789012/us-west-2/my-app:v1`
- `gcr.io/my-project/my-image:latest` becomes `gcr-mirror.internal.example.com/my-project/my-image:latest`

## Processing Order

When an image is being processed, mappings are evaluated in the following order:

1. **swap** - Checked first for registry-level matches
2. **exactSwap** - Checked second for exact matches, can override swap results
3. **regexSwap** - Checked last for pattern matches, can override all previous results

The last matching rule wins (later rules override earlier ones). If no rules match, the image is left unchanged.

## Complete Example

```yaml
apiVersion: imageshift.dev/v1
kind: Imageshift
metadata:
  name: imageshift
spec:
  default: docker.io
  namespaceSelector: imageshift.dev
  mappings:
    # Registry-level swaps (checked first)
    swap:
      - registry: docker.io
        target: registry.internal.example.com/dockerhub
      - registry: ghcr.io
        target: registry.internal.example.com/ghcr
      - registry: quay.io
        target: registry.internal.example.com/quay

    # Exact image swaps (checked second, can override swap)
    exactSwap:
      - reference: nginx:latest
        target: registry.internal.example.com/approved/nginx:1.25-alpine
      - reference: redis:latest
        target: registry.internal.example.com/approved/redis:7.2

    # Regex swaps for ECR (checked last, highest priority)
    regexSwap:
      - expression: "^(\\d+)\\.dkr\\.ecr\\.(.*)\\.amazonaws\\.com/(.*)$"
        target: "registry.internal.example.com/ecr/$1/$2/$3"
```

## Namespace Selection

ImageShift only processes pods in namespaces that have the selector label set to `enabled`.

Default selector label: `imageshift.dev=enabled`

To enable a namespace:

```bash
kubectl label namespace my-namespace imageshift.dev=enabled
```

To disable a namespace:

```bash
kubectl label namespace my-namespace imageshift.dev-
```

To use a custom selector label, set `spec.namespaceSelector`:

```yaml
spec:
  namespaceSelector: custom-label.example.com
```

Then label namespaces with:

```bash
kubectl label namespace my-namespace custom-label.example.com=enabled
```

## Limitations

- Only one Imageshift resource can exist in the cluster at a time
- The validating webhook rejects creation of additional Imageshift resources
- Init containers and ephemeral containers are also processed
- Images in Pod templates (Deployments, StatefulSets, etc.) are processed when Pods are created

## Status

The Imageshift resource does not currently expose status fields. Future versions may include:
- Count of processed pods
- Last reconciliation time
- Error conditions
