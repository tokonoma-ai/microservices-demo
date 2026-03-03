# Installing sock-shop on Kubernetes (kind)

This repo’s instructions target a **kind** cluster. From the repository root:

## Deploy (Kustomize)

```bash
./bin/build-dev   # build all services, load into kind
./bin/deploy      # deploy with locally-built :dev images
```

All 8 app services use locally-built `:dev` images. Build with `build-dev` first, then deploy. See [DEVELOP.md](../../DEVELOP.md) for details.

## Carts and MongoDB version

The carts service uses **Spring Boot 2.7**, **Java 17**, and the **MongoDB Java driver 4.x** (via Spring Data MongoDB 3.x). It talks to MongoDB over the modern wire protocol, so carts-db in the manifests uses **mongo:7** (driver 4.x supports MongoDB 5–8). Do not use legacy mongo:4.4 with this stack; the driver no longer uses OP_QUERY.

## Kubernetes manifests

There are 2 sets of manifests for deploying Sock Shop on Kubernetes: one in the [manifests directory](manifests/), and complete-demo.yaml. The complete-demo.yaml is a single file manifest
made by concatenating all the manifests from the manifests directory, so please regenerate it when changing files in the manifests directory.

## Monitoring

All monitoring is performed by prometheus. All services expose a `/metrics` endpoint. All services have a Prometheus Histogram called `request_duration_seconds`, which is automatically appended to create the metrics `_count`, `_sum` and `_bucket`.

The manifests for the monitoring are spread across the [manifests-monitoring](./manifests-monitoring) and [manifests-alerting](./manifests-alerting/) directories.

To use them, please run `kubectl create -f <path to directory>`.

### What's Included?

* Sock-shop grafana dashboards
* Alertmanager with 500 alert connected to slack
* Prometheus with config to scrape all k8s pods, connected to local alertmanager.

### Ports

Grafana will be exposed on the NodePort `31300` and Prometheus is exposed on `31090`. If running on a real cluster, the easiest way to connect to these ports is by port forwarding in a ssh command:
```
ssh -i $KEY -L 3000:$NODE_IN_CLUSTER:31300 -L 9090:$NODE_IN_CLUSTER:31090 ubuntu@$BASTION_IP
```
Where all the pertinent information should be entered. Grafana and Prometheus will be available on `http://localhost:3000` or `:9090`.

If on Minikube, you can connect via the VM IP address and the NodePort.

## Build from source

The application services (carts, catalogue, front-end, orders, payment, queue-master, shipping, user) live in separate repositories under the [microservices-demo](https://github.com/microservices-demo) GitHub organization. This repo contains only deployment config and a few auxiliaries (openapi, healthcheck).

To build and deploy using your own images:

1. Clone each service repo (e.g. `microservices-demo/carts`, `microservices-demo/front-end`, `microservices-demo/user`, etc.), build the Docker image, and push to your registry (or load into Minikube/kind with `eval $(minikube docker-env)` and build there).
2. Deploy with image overrides: either edit the manifest files to use your image names and tags, or use `kubectl set image` after running `./bin/deploy`, or use Kustomize to override the default `weaveworksdemos/*` images.
