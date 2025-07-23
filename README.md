# cert-manager-webhook-maas

A cert-manager ACME DNS01 solver webhook for Canonical MAAS (Metal as a Service).

<div align="center">

[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](LICENSE.md)

</div>

## Overview

This webhook implements an ACME DNS01 challenge solver for cert-manager that
integrates with MAAS DNS. This allows you to use cert-manager with Let's Encrypt
(or other ACME CAs) to automatically provision TLS certificates for domains
managed by MAAS.

## Features

- Integrates with MAAS DNS API for DNS-01 challenges
- Supports wildcard certificates
- Configurable per-issuer with different MAAS instances
- Uses official Canonical gomaasclient library

## Prerequisites

- Kubernetes cluster with cert-manager installed (>= v1.18.0)
- MAAS instance with DNS enabled
- MAAS API credentials

## Installation

### Using Helm (Recommended)

1. Add the Helm repository:

```bash
helm repo add maas-webhook https://kogito-ops.github.io/cert-manager-webhook-maas
helm repo update
```

2. Create a secret containing your MAAS API credentials:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: maas-secret
  namespace: cert-manager
type: Opaque
stringData:
  api-key: "your-consumer-key:your-token-key:your-token-secret"
```

```bash
kubectl apply -f - <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: maas-secret
  namespace: cert-manager
type: Opaque
stringData:
  api-key: "your-consumer-key:your-token-key:your-token-secret"
EOF
```

3. Install the webhook:

```bash
helm install maas-webhook maas-webhook/maas-webhook \
  --namespace cert-manager \
  --set groupName=acme.maas.io
```

### Using kubectl

Alternatively, you can install using the raw manifests:

```bash
kubectl apply -f https://github.com/kogito-ops/cert-manager-webhook-maas/releases/latest/download/maas-webhook.yaml
```

### From source

```bash
git clone https://github.com/kogito-ops/cert-manager-webhook-maas.git
cd cert-manager-webhook-maas
helm install maas-webhook charts/maas-webhook \
  --namespace cert-manager \
  --set groupName=acme.maas.io
```

## Configuration

### Issuer Configuration

Create an Issuer or ClusterIssuer that uses the MAAS webhook:

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
      name: letsencrypt-maas-account
    solvers:
    - dns01:
        webhook:
          groupName: acme.maas.io
          solverName: maas
          config:
            # Name of the secret containing MAAS API credentials
            secretName: maas-secret

            # MAAS API endpoint
            apiUrl: http://maas.example.com:5240/MAAS

            # Optional: DNS zone name
            # If not specified, will be extracted from the domain
            zoneName: example.com

            # Optional: MAAS API version (default: "2.0")
            apiVersion: "2.0"
```

### Certificate Request

Request a certificate using the configured issuer:

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: example-com
  namespace: default
spec:
  secretName: example-com-tls
  issuerRef:
    name: letsencrypt-maas
    kind: ClusterIssuer
  dnsNames:
  - example.com
  - "*.example.com"
```

## MAAS API Credentials

To obtain MAAS API credentials:

1. Log into your MAAS web interface
2. Click on your username in the top right
3. Go to "API keys"
4. Generate a new API key

## Configuration Options

| Parameter    | Description | Default | Required |
| ------------ | ----------- | ------- | -------- |
| `secretName` | Name of the Kubernetes secret containing MAAS API key | - | Yes |
| `apiUrl`     | MAAS API endpoint URL | - | Yes |
| `zoneName`   | DNS zone to use. If not specified, extracted from domain | - | No |
| `apiVersion` | MAAS API version | `2.0` | No |

## Troubleshooting

### Check webhook logs

```bash
kubectl logs -n cert-manager deployment/maas-webhook
```

### Verify DNS record creation

Check if the TXT record was created in MAAS:

```bash
# Using MAAS CLI
maas $PROFILE dnsresources read domain=$DOMAIN name=_acme-challenge.$SUBDOMAIN
```

### Common Issues

1. **Authentication errors**: Verify your API key is correct and the secret is in the correct namespace
2. **DNS record not created**: Check that the zone exists in MAAS and the API URL is correct
3. **Certificate stays in pending**: Check cert-manager logs for detailed error messages

## Development

### Building

```bash
go build -o webhook .
```

### Running tests

```bash
go test ./...
```

### Building Docker image

```bash
docker build -t cert-manager-webhook-maas:latest .
```

## License

This project is licensed under the Apache License 2.0 - see the LICENSE file for details.

## Acknowledgments

- Based on the [cert-manager webhook example](https://github.com/cert-manager/webhook-example)
- Uses the [Canonical gomaasclient](https://github.com/canonical/gomaasclient) library
