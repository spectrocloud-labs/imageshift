---
sidebar_position: 2
---

# Quick Start

Get ImageShift up and running in minutes with this quick start guide.

## 1. Install ImageShift

Using Helm (recommended):

```bash
helm install imageshift oci://ghcr.io/wcrum/imageshift/chart \
  --version <version> \
  -n imageshift-system \
  --create-namespace
```

Or using Kustomize:

```bash
kubectl apply -k https://github.com/wcrum/imageshift/config/default
```

## 2. Verify Installation

Check that the ImageShift pods are running:

```bash
kubectl get pods -n imageshift-system
```

You should see the controller manager pod in a `Running` state:

```
NAME                                           READY   STATUS    RESTARTS   AGE
imageshift-controller-manager-xxx-yyy          1/1     Running   0          30s
```

## 3. Create an Imageshift Resource

Create a file named `imageshift.yaml`:

```yaml
apiVersion: imageshift.dev/v1
kind: Imageshift
metadata:
  name: imageshift
spec:
  mappings:
    swap:
      - registry: docker.io
        target: registry.internal.example.com/dockerhub
```

Apply it to your cluster:

```bash
kubectl apply -f imageshift.yaml
```

## 4. Enable a Namespace

Label any namespace where you want image swapping to occur:

```bash
kubectl label namespace default imageshift.dev=enabled
```

## 5. Test the Mutation

Deploy a pod with a Docker Hub image:

```bash
kubectl run nginx --image=nginx:latest -n default
```

Verify the image was swapped:

```bash
kubectl get pod nginx -n default -o jsonpath='{.spec.containers[0].image}'
```

The output should show the redirected image:

```
registry.internal.example.com/dockerhub/library/nginx:latest
```

## 6. Clean Up (Optional)

To remove the test pod:

```bash
kubectl delete pod nginx -n default
```

To disable image swapping for a namespace:

```bash
kubectl label namespace default imageshift.dev-
```

## Next Steps

- [Helm Installation](./helm) - Full Helm installation options and values
- [Kustomize Installation](./kustomize) - Kustomize-based installation and customization
- [CRD Reference](../reference/crd) - Complete Imageshift specification
- [Configuration Examples](../reference/configuration) - Real-world configuration examples
