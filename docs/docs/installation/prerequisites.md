---
sidebar_position: 1
---

# Prerequisites

Before installing ImageShift, ensure your environment meets the following requirements.

## Kubernetes Version

- **Kubernetes 1.19+** is required
- The mutating admission webhook API must be enabled (enabled by default)

Verify your Kubernetes version:

```bash
kubectl version
```

## kubectl

You need `kubectl` configured to communicate with your cluster:

```bash
kubectl cluster-info
```

## cert-manager

ImageShift uses cert-manager to manage TLS certificates for the webhook endpoints.

### Installing cert-manager

If cert-manager is not already installed, install version 1.16.3 or later:

```bash
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.16.3/cert-manager.yaml
```

Verify cert-manager is running:

```bash
kubectl get pods -n cert-manager
```

You should see three pods running:
- `cert-manager`
- `cert-manager-cainjector`
- `cert-manager-webhook`

Wait for all pods to be ready before proceeding with ImageShift installation.

## Helm (Optional)

If you prefer to install ImageShift using Helm:

- **Helm 3.x** is required

Verify your Helm version:

```bash
helm version
```

## Cluster Permissions

You need cluster-admin permissions to:
- Create Custom Resource Definitions (CRDs)
- Create ClusterRoles and ClusterRoleBindings
- Create MutatingWebhookConfigurations and ValidatingWebhookConfigurations
- Create resources in the `imageshift-system` namespace

## Network Requirements

The Kubernetes API server must be able to reach the ImageShift webhook service. This is typically automatic in most cluster configurations.

If you're using network policies, ensure traffic is allowed from the API server to the `imageshift-system` namespace on port 443.

## Next Steps

Once you've verified all prerequisites, proceed to the [Quick Start](./quick-start) guide.
