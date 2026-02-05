---
sidebar_position: 3
---

# Helm Installation

Install ImageShift using Helm for a customizable deployment.

## Install ImageShift

ImageShift is distributed as an OCI Helm chart. No `helm repo add` is required.

Basic installation:

```bash
helm install imageshift oci://ghcr.io/wcrum/imageshift/chart \
  --version <version> \
  -n imageshift-system \
  --create-namespace
```

With custom values:

```bash
helm install imageshift oci://ghcr.io/wcrum/imageshift/chart \
  --version <version> \
  -n imageshift-system \
  --create-namespace \
  -f values.yaml
```

## Configuration Values

The following table lists the configurable parameters of the ImageShift chart.

### Global Settings

| Parameter | Description | Default |
|-----------|-------------|---------|
| `replicaCount` | Number of controller replicas | `1` |
| `image.repository` | Controller image repository | `ghcr.io/wcrum/imageshift` |
| `image.tag` | Controller image tag | `latest` |
| `image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `imagePullSecrets` | Image pull secrets | `[]` |
| `nameOverride` | Override chart name | `""` |
| `fullnameOverride` | Override full release name | `""` |

### Controller Settings

| Parameter | Description | Default |
|-----------|-------------|---------|
| `controller.resources.limits.cpu` | CPU limit | `500m` |
| `controller.resources.limits.memory` | Memory limit | `128Mi` |
| `controller.resources.requests.cpu` | CPU request | `10m` |
| `controller.resources.requests.memory` | Memory request | `64Mi` |
| `controller.nodeSelector` | Node selector labels | `{}` |
| `controller.tolerations` | Pod tolerations | `[]` |
| `controller.affinity` | Pod affinity rules | `{}` |

### Webhook Settings

| Parameter | Description | Default |
|-----------|-------------|---------|
| `webhook.port` | Webhook server port | `9443` |
| `webhook.certManager.enabled` | Use cert-manager for TLS | `true` |
| `webhook.certManager.issuerRef` | cert-manager issuer reference | `{}` |

### RBAC Settings

| Parameter | Description | Default |
|-----------|-------------|---------|
| `rbac.create` | Create RBAC resources | `true` |
| `serviceAccount.create` | Create service account | `true` |
| `serviceAccount.name` | Service account name | `""` |
| `serviceAccount.annotations` | Service account annotations | `{}` |

## Example values.yaml

```yaml
replicaCount: 2

image:
  repository: ghcr.io/wcrum/imageshift
  tag: <version>  # Replace with desired version from https://github.com/wcrum/imageshift/releases
  pullPolicy: IfNotPresent

controller:
  resources:
    limits:
      cpu: 500m
      memory: 256Mi
    requests:
      cpu: 100m
      memory: 128Mi

  nodeSelector:
    kubernetes.io/os: linux

  tolerations:
    - key: "node-role.kubernetes.io/control-plane"
      operator: "Exists"
      effect: "NoSchedule"

webhook:
  port: 9443
  certManager:
    enabled: true

rbac:
  create: true

serviceAccount:
  create: true
  annotations:
    eks.amazonaws.com/role-arn: arn:aws:iam::123456789012:role/imageshift-role
```

## Upgrade

To upgrade an existing installation:

```bash
helm upgrade imageshift oci://ghcr.io/wcrum/imageshift/chart \
  --version <version> \
  -n imageshift-system
```

## Uninstall

To remove ImageShift:

```bash
helm uninstall imageshift -n imageshift-system
```

Note: This will not remove the CRDs. To remove CRDs:

```bash
kubectl delete crd imageshifts.imageshift.dev
```

## Next Steps

- [CRD Reference](../reference/crd) - Configure your Imageshift resource
- [Configuration Examples](../reference/configuration) - Real-world configuration examples
