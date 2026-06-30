# Amprobe Kubernetes Deployment

This directory contains base Kubernetes manifests for deploying the Amprobe Server.

## Prerequisites

- Kubernetes cluster (1.25+)
- kubectl configured
- Container image built and pushed to a registry accessible by the cluster

## Quick Start

```bash
# 1. Create namespace and secrets
kubectl apply -f deploy/k8s/namespace.yaml
kubectl apply -f deploy/k8s/secret.yaml

# 2. Create ConfigMap and PVCs
kubectl apply -f deploy/k8s/configmap.yaml
kubectl apply -f deploy/k8s/pvc.yaml

# 3. Deploy the Server
kubectl apply -f deploy/k8s/deployment.yaml
kubectl apply -f deploy/k8s/service.yaml

# 4. Verify
kubectl get pods -n amprobe
kubectl logs -n amprobe -l app=amprobe-server
```

## Important Security Notes

1. **Secrets**: `secret.yaml` contains placeholder values. In production, use:
   - External Secret Management (e.g., Vault, Sealed Secrets, AWS Secrets Manager)
   - `kubectl create secret` with values generated from a secure random source

2. **Signing Key**: `AMPROBE_AUTH_SIGNING_KEY` must be ≥32 bytes of random data.
   ```bash
   openssl rand -hex 32
   ```

3. **ConfigMap**: The base config disables TLS on the control tunnel. Enable it in production by:
   - Setting `Control.TLS.Enable = true`
   - Mounting certificate files into `/etc/amprobe/control-certs`

## Health Probes

The deployment configures `livenessProbe` (`/health`) and `readinessProbe` (`/ready`).
The readiness probe checks DB connectivity and tunnel state when dependencies are injected.

## Storage

Two PVCs are created:
- `amprobe-data`: SQLite database and install ID file
- `amprobe-logs`: Application logs

For production, consider migrating to PostgreSQL or ClickHouse by updating the ConfigMap.
