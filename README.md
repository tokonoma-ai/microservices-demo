# Sock Shop Microservices Demo

A microservices demo application for an online sock shop.

## Prerequisites

- Kubernetes cluster (e.g., Docker Desktop, minikube, kind)
- `kubectl` configured to connect to your cluster

## Deploy

```bash
cd deploy/kubernetes
./deploy.sh
```

This deploys:
- Sock Shop services (namespace: `sock-shop`)
- Prometheus + Grafana (namespace: `monitoring`)
- Jaeger tracing (namespace: `jaeger`)

## Port Forward

```bash
cd deploy/kubernetes
./port-forward.sh
```

Access points:
- Sock Shop: http://localhost:8080
- Grafana: http://localhost:3001 (admin/admin)
- Jaeger UI: http://localhost:16686

Press `Ctrl+C` to stop all port forwards.

## Important

Pods may show as started but need to be `1/1 Running` before the catalogue appears on the demo site. Check status with:

```bash
kubectl get pods -n sock-shop
```

Wait until all pods show `1/1` in the READY column.
