# MAAS Webhook for cert-manager

A Helm chart for deploying the cert-manager MAAS DNS01 challenge solver webhook.

## Installation

### Add the Helm repository

```bash
helm repo add maas-webhook https://kogito-ops.github.io/cert-manager-webhook-maas
helm repo update
```

### Install the chart

```bash
helm install maas-webhook maas-webhook/maas-webhook --namespace cert-manager
```

## Configuration

The following table lists the configurable parameters of the chart and their default values.

| Parameter | Description | Default |
|-----------|-------------|---------|
| `groupName` | The API group name for the webhook | `acme.maas.io` |
| `image.registry` | Container image registry | `ghcr.io` |
| `image.repository` | Container image repository | `kogito-ops/cert-manager-webhook-maas` |
| `image.tag` | Container image tag | `""` (uses chart appVersion) |
| `image.pullPolicy` | Container image pull policy | `IfNotPresent` |
| `certManager.namespace` | cert-manager namespace | `cert-manager` |
| `certManager.serviceAccountName` | cert-manager service account name | `cert-manager` |
| `secretName` | List of secret names containing MAAS API credentials | `["maas-secret"]` |
| `service.type` | Kubernetes service type | `ClusterIP` |
| `service.port` | Kubernetes service port | `443` |
| `resources` | Pod resource requests and limits | `{}` |
| `nodeSelector` | Node labels for pod assignment | `{}` |
| `tolerations` | Tolerations for pod assignment | `[]` |
| `affinity` | Affinity rules for pod assignment | `{}` |

## Example Usage

### 1. Create MAAS API Secret

```bash
kubectl create secret generic maas-secret \
  --from-literal=api-key="consumer-key:token-key:token-secret" \
  --namespace cert-manager
```

### 2. Install the Chart

```bash
helm install maas-webhook maas-webhook/maas-webhook \
  --namespace cert-manager \
  --set groupName=acme.maas.io
```

### 3. Create an Issuer

```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-maas
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: admin@example.com
    privateKeySecretRef:
      name: letsencrypt-maas
    solvers:
    - dns01:
        webhook:
          groupName: acme.maas.io
          solverName: maas
          config:
            secretName: maas-secret
            apiUrl: http://maas.example.com:5240/MAAS
            zoneName: example.com
```

## Uninstallation

```bash
helm uninstall maas-webhook --namespace cert-manager
```