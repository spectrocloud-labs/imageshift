---
sidebar_position: 2
---

# Configuration Examples

Real-world configuration examples for common ImageShift use cases.

## Air-Gapped Environment

For disconnected networks where all images must come from an internal registry:

```yaml
apiVersion: imageshift.dev/v1
kind: Imageshift
metadata:
  name: imageshift
spec:
  default: docker.io
  mappings:
    swap:
      # Docker Hub
      - registry: docker.io
        target: airgap-registry.internal/dockerhub

      # GitHub Container Registry
      - registry: ghcr.io
        target: airgap-registry.internal/ghcr

      # Google Container Registry
      - registry: gcr.io
        target: airgap-registry.internal/gcr

      # Quay.io
      - registry: quay.io
        target: airgap-registry.internal/quay

      # Kubernetes Registry
      - registry: registry.k8s.io
        target: airgap-registry.internal/k8s
```

## Multi-Region Deployment

Use regional registries to minimize latency and cross-region data transfer:

### US-West Region

```yaml
apiVersion: imageshift.dev/v1
kind: Imageshift
metadata:
  name: imageshift
spec:
  mappings:
    swap:
      - registry: docker.io
        target: us-west-2.registry.example.com/dockerhub
      - registry: ghcr.io
        target: us-west-2.registry.example.com/ghcr
```

### EU-Central Region

```yaml
apiVersion: imageshift.dev/v1
kind: Imageshift
metadata:
  name: imageshift
spec:
  mappings:
    swap:
      - registry: docker.io
        target: eu-central-1.registry.example.com/dockerhub
      - registry: ghcr.io
        target: eu-central-1.registry.example.com/ghcr
```

## AWS ECR Redirection

Redirect ECR images from one account or region to another using regex:

```yaml
apiVersion: imageshift.dev/v1
kind: Imageshift
metadata:
  name: imageshift
spec:
  mappings:
    regexSwap:
      # Redirect all ECR images to a central account
      - expression: "^(\\d+)\\.dkr\\.ecr\\.([a-z0-9-]+)\\.amazonaws\\.com/(.*)$"
        target: "999999999999.dkr.ecr.us-east-1.amazonaws.com/mirror/$1/$2/$3"

      # Or redirect to a different region
      - expression: "^(\\d+)\\.dkr\\.ecr\\.(us-west-2)\\.amazonaws\\.com/(.*)$"
        target: "$1.dkr.ecr.us-east-1.amazonaws.com/$3"
```

## Google Container Registry (GCR) Redirection

Redirect GCR images to a different project or Artifact Registry:

```yaml
apiVersion: imageshift.dev/v1
kind: Imageshift
metadata:
  name: imageshift
spec:
  mappings:
    regexSwap:
      # Redirect gcr.io to Artifact Registry
      - expression: "^gcr\\.io/([^/]+)/(.*)$"
        target: "us-docker.pkg.dev/$1/gcr-mirror/$2"

      # Redirect regional GCR to Artifact Registry
      - expression: "^([a-z]+)\\.gcr\\.io/([^/]+)/(.*)$"
        target: "$1-docker.pkg.dev/$2/gcr-mirror/$3"
```

## Pinned Image Versions

Use exact swaps to enforce specific image versions:

```yaml
apiVersion: imageshift.dev/v1
kind: Imageshift
metadata:
  name: imageshift
spec:
  mappings:
    exactSwap:
      # Pin nginx to a specific approved version
      - reference: nginx:latest
        target: approved-registry.example.com/nginx:1.25.3-alpine
      - reference: nginx
        target: approved-registry.example.com/nginx:1.25.3-alpine

      # Pin redis to a specific version
      - reference: redis:latest
        target: approved-registry.example.com/redis:7.2.4
      - reference: redis
        target: approved-registry.example.com/redis:7.2.4

      # Pin postgres
      - reference: postgres:latest
        target: approved-registry.example.com/postgres:16.1
      - reference: postgres:16
        target: approved-registry.example.com/postgres:16.1
```

## Combined Mapping Strategies

Use multiple mapping types together for comprehensive coverage:

```yaml
apiVersion: imageshift.dev/v1
kind: Imageshift
metadata:
  name: imageshift
spec:
  default: docker.io
  mappings:
    # First: Registry-level swaps (checked first)
    swap:
      - registry: docker.io
        target: internal.example.com/dockerhub
      - registry: ghcr.io
        target: internal.example.com/ghcr
      - registry: quay.io
        target: internal.example.com/quay

    # Second: Exact matches for pinned versions (can override swap)
    exactSwap:
      - reference: nginx:latest
        target: internal.example.com/approved/nginx:1.25.3
      - reference: redis:latest
        target: internal.example.com/approved/redis:7.2.4

    # Third: Regex for complex patterns (highest priority, checked last)
    regexSwap:
      - expression: "^(\\d+)\\.dkr\\.ecr\\.([a-z0-9-]+)\\.amazonaws\\.com/(.*)$"
        target: "internal.example.com/ecr-mirror/$3"
      - expression: "^gcr\\.io/([^/]+)/(.*)$"
        target: "internal.example.com/gcr-mirror/$1/$2"
```

## Development Environment

Redirect to a local registry for development:

```yaml
apiVersion: imageshift.dev/v1
kind: Imageshift
metadata:
  name: imageshift
spec:
  mappings:
    swap:
      - registry: docker.io
        target: localhost:5000/dockerhub
      - registry: ghcr.io
        target: localhost:5000/ghcr
```

## Pull-Through Cache

Redirect to a pull-through cache proxy:

```yaml
apiVersion: imageshift.dev/v1
kind: Imageshift
metadata:
  name: imageshift
spec:
  mappings:
    swap:
      - registry: docker.io
        target: cache.example.com/docker.io
      - registry: gcr.io
        target: cache.example.com/gcr.io
      - registry: ghcr.io
        target: cache.example.com/ghcr.io
      - registry: quay.io
        target: cache.example.com/quay.io
```

## Namespace-Specific Configuration

While ImageShift applies cluster-wide rules, you can control which namespaces are affected:

```bash
# Enable for production namespaces
kubectl label namespace prod-app imageshift.dev=enabled
kubectl label namespace prod-api imageshift.dev=enabled

# Keep development namespaces using original images
# (don't label them, or explicitly remove the label)
kubectl label namespace dev-app imageshift.dev-
```

## Registry Normalization

ImageShift normalizes registry names to handle aliases. For example:

- `docker.io` is equivalent to `index.docker.io` (Docker Hub's canonical name)
- Images without a registry prefix (e.g., `nginx:latest`) use the `spec.default` registry

This means you can use `docker.io` in your swap rules and it will match images specified as either `docker.io/nginx` or `nginx` (when default is `docker.io`).

## Tips for Writing Regex Patterns

### Escape Special Characters

In YAML, backslashes need to be escaped. Use `\\` for regex escape sequences:

```yaml
# Match digits
expression: "^(\\d+)\\.dkr\\.ecr"  # Correct
expression: "^(\d+)\.dkr\.ecr"     # Won't work as expected
```

### Use Capture Groups

Capture groups `()` can be referenced in the target with `$1`, `$2`, etc.:

```yaml
expression: "^([^/]+)/([^/]+)/(.*)$"
target: "new-registry.com/$1/$2/$3"
```

### Test Patterns

Test your regex patterns before deploying:

```bash
echo "123456789012.dkr.ecr.us-west-2.amazonaws.com/my-app:v1" | \
  sed -E 's/^([0-9]+)\.dkr\.ecr\.([a-z0-9-]+)\.amazonaws\.com\/(.*)$/mirror.example.com\/\1\/\2\/\3/'
```
