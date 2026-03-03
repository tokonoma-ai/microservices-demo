# Developing and Deploying Sock Shop

This repo is the single source of truth for deploying Sock Shop to a **kind** cluster.
It contains both the Kubernetes manifests and the source code for services that are
actively modified.

## What's in this repo

| Directory | Contents |
|-----------|----------|
| `deploy/kubernetes/manifests/` | Base Kustomize manifests (all services) |
| `deploy/kubernetes/kustomization.yaml` | Kustomize root — patches all 8 app services to `:dev` |
| `carts/`, `catalogue/`, `orders/`, `payment/`, `shipping/`, `queue-master/`, `user/`, `front-end/` | Source code for services built locally |
| `bin/` | Build, deploy, and utility scripts |

## Services and build status

All 8 app services are built from source.

| Service | Image when stock | Built locally? | Notes |
|---------|-----------------|---------------|-------|
| **carts** | `weaveworksdemos/carts:0.4.8` | Yes → `:dev` | Java 17 |
| **orders** | `weaveworksdemos/orders:0.4.7` | Yes → `:dev` | Java 17 |
| **catalogue** | `weaveworksdemos/catalogue:0.3.5` | Yes → `:dev` | Go 1.24 |
| **payment** | `weaveworksdemos/payment:0.4.3` | Yes → `:dev` | Go 1.24 |
| **front-end** | `weaveworksdemos/front-end:0.3.12` | Yes → `:dev` | Node.js 20 |
| **shipping** | `weaveworksdemos/shipping:0.4.8` | Yes → `:dev` | Java 17 |
| **queue-master** | `weaveworksdemos/queue-master:0.3.1` | Yes → `:dev` | Java 17 |
| **user** | `weaveworksdemos/user:0.4.7` | Yes → `:dev` | Go 1.24 |

Databases, RabbitMQ, Redis all use stock images and have no source code here.

## Build and deploy

### 1. Build all services and load into kind

```bash
./bin/build-dev
```

This builds all 8 services and loads the `:dev` images into the kind cluster:
- **carts** (Java 17, Maven)
- **orders** (Java 17, Maven)
- **shipping** (Java 17, Maven)
- **queue-master** (Java 17, Maven)
- **catalogue** (Go 1.24, multi-stage Docker)
- **payment** (Go 1.24, multi-stage Docker)
- **user** (Go 1.24, multi-stage Docker)
- **front-end** (Node.js 20, Docker)

### 2. Deploy

```bash
./bin/deploy
```

### 3. Or do both steps: build then deploy

```bash
./bin/build-dev && ./bin/deploy
```

## Building a single service

From the repo root:

**Java services (carts, orders, shipping, queue-master):**
```bash
cd carts  # or orders, shipping, queue-master
docker run --rm -v "$(pwd)":/build -w /build maven:3-eclipse-temurin-17 mvn package -DskipTests -q
docker build -t weaveworksdemos/carts:dev .
kind load docker-image weaveworksdemos/carts:dev --name qw
kubectl rollout restart deployment/carts -n sock-shop
```

Use `maven:3-eclipse-temurin-17` for all Java services (carts, orders, shipping, queue-master).

**Go services (catalogue, payment, user):**
```bash
cd catalogue  # or payment, user
docker build -t weaveworksdemos/catalogue:dev -f docker/catalogue/Dockerfile .
kind load docker-image weaveworksdemos/catalogue:dev --name qw
kubectl rollout restart deployment/catalogue -n sock-shop
```

**Node.js service (front-end):**
```bash
cd front-end
docker build -t weaveworksdemos/front-end:dev .
kind load docker-image weaveworksdemos/front-end:dev --name qw
kubectl rollout restart deployment/front-end -n sock-shop
```

## Tokonoma demo setup

After deploying sock-shop, deploy load generators and trigger the demo from the
[load-generators](https://github.com/tokonoma-ai/load-generators) repo:

```bash
cd ../load-generators
./bin/demo
```

Verify with `./bin/demo-ready` (from this repo).

## Utility scripts

| Script | Purpose |
|--------|---------|
| `bin/deploy` | Deploy sock-shop |
| `bin/build-dev` | Build all 8 dev services and load into kind |
| `bin/port-forward` | Port-forward frontend (:8080) |
| `bin/demo-ready` | Verify demo readiness |

## Cluster layout

| Namespace | Contents | Managed by |
|-----------|----------|------------|
| `sock-shop` | All Sock Shop services | This repo (`bin/deploy`) |
| `loadtest` | Load generators + checkout-fail-injector | [load-generators](https://github.com/tokonoma-ai/load-generators) repo |
| `qw` | Quickwit, Prometheus, Grafana, OTel, MCP server | Platform scripts (separate) |
