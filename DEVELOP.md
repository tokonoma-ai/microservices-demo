# Developing and Deploying Sock Shop from Source

This guide covers how to **change application code** and then build and deploy to a **kind** cluster (local only).

## Where the code lives

| What | Where | Used in deploy? |
|------|--------|------------------|
| **Deployment config** (K8s manifests) | This repo (`deploy/kubernetes/manifests/`) | Yes |
| **openapi**, **healthcheck**, **graphs** | This repo (`openapi/`, `healthcheck/`, `graphs/`) | Not in main sock-shop manifests (CI/test) |
| **Main app services** (carts, catalogue, front-end, orders, payment, queue-master, shipping, user, catalogue-db, user-db) | **Separate GitHub repos** | Yes — these are the containers the deploy runs |

So to change the **shop UI, cart, orders, users, catalogue, payment, shipping**, etc., you need the source from the individual service repos.

---

## Next steps: make changes and deploy

### 1. Get the service source code

Clone the repos you want to edit. You can put them next to this repo (e.g. `../carts`, `../front-end`) or in a single folder:

```bash
# From this repo root
./bin/clone-services
```

That script clones all service repos into `../microservices-demo-services/` (or set `SERVICES_DIR`). Or clone by hand:

| Service      | Repo |
|-------------|------|
| carts       | https://github.com/microservices-demo/carts |
| catalogue   | https://github.com/microservices-demo/catalogue |
| front-end   | https://github.com/microservices-demo/front-end |
| orders      | https://github.com/microservices-demo/orders |
| payment     | https://github.com/microservices-demo/payment |
| queue-master| https://github.com/microservices-demo/queue-master |
| shipping    | https://github.com/microservices-demo/shipping |
| user        | https://github.com/microservices-demo/user |

**catalogue-db and user-db:** There are no separate repos (or they were removed). You don’t need them—the deploy uses the pre-built images from Docker Hub (`weaveworksdemos/catalogue-db`, `weaveworksdemos/user-db`). Only the app services in the table need to be cloned when changing code.

### 2. Make your code changes

Edit the service(s) you care about in their respective repos. Each repo has its own Dockerfile and build (Java, Go, Node, etc.).

### 3. Build and load images (kind)

Use tag `dev` and the same image name as the dev overlay (`weaveworksdemos/<service>:dev`). **`docker build -t ...` only creates a local image—it does not upload or push to Docker Hub** (or any registry). For kind you then load that local image into the cluster with `kind load docker-image`.

```bash
export IMAGE_TAG=dev
export IMAGE_PREFIX=weaveworksdemos  
```

From each service repo you changed:

**Java services (carts, orders, payment, shipping, queue-master):** build the JAR first, then Docker:

```bash
cd /path/to/carts   # or orders, payment, shipping, queue-master
mvn -DskipTests package
docker build -t ${IMAGE_PREFIX}/carts:${IMAGE_TAG} .
kind load docker-image ${IMAGE_PREFIX}/carts:${IMAGE_TAG}
```

**One-liner (carts, from this repo root; kind cluster name `qw`):**
```bash
cd carts && mvn -DskipTests package && docker build -t weaveworksdemos/carts:dev . && kind load docker-image weaveworksdemos/carts:dev --name qw && kubectl rollout restart deployment/carts -n sock-shop
```
Use `--name <cluster>` if your kind cluster has a different name.

**Go/Node services (catalogue, front-end, user):** no compile step; just Docker:

```bash
cd /path/to/front-end
docker build -t ${IMAGE_PREFIX}/front-end:${IMAGE_TAG} .
kind load docker-image ${IMAGE_PREFIX}/front-end:${IMAGE_TAG}
```
**One-liner (catalogue, from this repo root; kind cluster name `qw`):**
```bash
cd catalogue && docker build -t weaveworksdemos/catalogue:dev -f docker/catalogue/Dockerfile . && kind load docker-image weaveworksdemos/catalogue:dev --name qw && kubectl rollout restart deployment/catalogue -n sock-shop
```
Repeat for every service you modified (and optionally for all services so versions match).

### 5. Deploy (Kustomize)

**Official images (no build):**
```bash
./bin/deploy
```

**Your built images (after building and `kind load docker-image`):**
```bash
./bin/deploy dev
```

### 6. Verify

```bash
kubectl get pods -n sock-shop
# Open front-end (NodePort 30001 or port-forward)
kubectl port-forward -n sock-shop svc/front-end 8079:80
# Then open http://localhost:8079
```

---

## Summary workflow (kind)

1. **Clone** service repos → `./bin/clone-services` or clone by hand.
2. **Edit** code in the repo(s) you care about.
3. **Build** images: `docker build -t weaveworksdemos/<service>:dev .` (local only; no push). Then **load into kind:** `kind load docker-image weaveworksdemos/<service>:dev`.
4. **Deploy:** `./bin/deploy` (official images) or `./bin/deploy dev` (your images).

