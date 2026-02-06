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
| `spec.default` | string | `"docker.io"` | The default registry assumed for images without an explicit registry prefix (e.g., `nginx` becomes `docker.io/library/nginx`) |
| `spec.namespaceSelector` | string | `"imageshift.dev"` | Reserved for future use. Currently ignored - namespaces must use label `imageshift.dev=enabled` |
| `spec.enforceExistingPods` | boolean | `false` | When enabled, the controller will delete pods that have images not matching the swap rules, forcing pod recreation with correct images |
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

## Enforcing Existing Pods

By default, ImageShift only mutates pods at creation time via the admission webhook. Pods that already exist when ImageShift is installed or when rules are changed will not be affected.

To enforce image swap rules on existing pods, enable the `enforceExistingPods` field:

```yaml
apiVersion: imageshift.dev/v1
kind: Imageshift
metadata:
  name: imageshift
spec:
  enforceExistingPods: true
  mappings:
    swap:
      - registry: docker.io
        target: registry.internal.example.com/dockerhub
```

:::warning
**This is a disruptive operation.** When enabled, the controller will **delete** any pods in labeled namespaces that have images not matching the current swap rules. This forces the pods to be recreated by their controllers (Deployment, StatefulSet, etc.) with the correct images applied by the webhook.

Only enable this feature when you understand the impact:
- Pods will be deleted and recreated
- Stateless workloads managed by controllers will recover automatically
- Standalone pods (not managed by a controller) will be permanently deleted
- Brief service interruptions may occur during pod recreation
:::

## Namespace Selection

ImageShift only processes pods in namespaces that have the label `imageshift.dev=enabled`.

To enable a namespace:

```bash
kubectl label namespace my-namespace imageshift.dev=enabled
```

To disable a namespace:

```bash
kubectl label namespace my-namespace imageshift.dev-
```

:::note
The `spec.namespaceSelector` field is reserved for future use. Currently, ImageShift always uses `imageshift.dev` as the label key.
:::

## Pod Annotations and Labels

When ImageShift mutates a pod, it adds the following metadata:

### Labels

| Label | Value | Description |
|-------|-------|-------------|
| `imageshift.dev/mutated` | `true` | Added when at least one container image was swapped |

### Annotations

For each mutated container, ImageShift stores the original image reference:

| Annotation Pattern | Description |
|--------------------|-------------|
| `<container-name>.container.imageshift.dev/original` | Original image for regular containers |
| `<container-name>.initContainer.imageshift.dev/original` | Original image for init containers |

Example for a pod with an nginx container:

```yaml
metadata:
  labels:
    imageshift.dev/mutated: "true"
  annotations:
    nginx.container.imageshift.dev/original: "nginx:latest"
```

You can use these annotations to verify mutations or for debugging:

```bash
kubectl get pod <pod-name> -o jsonpath='{.metadata.annotations}' | jq .
```

## Limitations

- Only one Imageshift resource can exist in the cluster at a time
- The validating webhook rejects creation of additional Imageshift resources
- Init containers are also processed
- Images in Pod templates (Deployments, StatefulSets, etc.) are processed when Pods are created

## Status

The Imageshift resource exposes the following status fields:

| Field | Type | Description |
|-------|------|-------------|
| `conditions` | array | Standard Kubernetes conditions for the resource |
| `lastReconciled` | timestamp | When the resource was last successfully reconciled |
| `configValid` | boolean | Whether the current configuration is valid |
| `mutatedPodCount` | integer | Total number of pods that have been mutated |
