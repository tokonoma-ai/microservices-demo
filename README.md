## Tokonoma Demo Setup

Demonstrates log investigation using Tokonoma (toko-mcp) against live Sock Shop services.

### Prerequisites

- `docker`, `kubectl`, `helm` v3+
- `kind` (for local clusters), `awscli` + `eksctl` (for EKS), or `gcloud` (for GKE)
- The [2o](https://github.com/tokonoma-ai/2o) platform repo cloned alongside this repo

### 1. Platform (one-time)

```bash
cd ../2o/platform
./bin/setup          # creates kind cluster "tokonoma"
./bin/deploy --kind  # deploys Quickwit, OTel collector, MCP server
./bin/port-forward   # waits for readiness, starts port-forwards
```

### 2. Build and Deploy Sock Shop

All services are built from source. Specify your target cluster:

```bash
# kind (local)
./bin/build --kind && ./bin/deploy --kind

# EKS
./bin/build --eks && ./bin/deploy --eks

# GKE
./bin/build --gke && ./bin/deploy --gke
```

### 3. Demo Setup (load generators)

```bash
cd ../load-generators
./bin/demo   # build demodatagen, deploy, trigger first checkout failure
```

The demodatagen load generator runs weighted-random user journeys continuously and triggers injection scenarios (checkout failures, payment errors) every ~30 minutes.

### 4. Verify Readiness

```bash
./bin/demo-ready
```

### 5. Demo Prompt

"Checkout failed for user 57a98d98e4b00679b4a830af in the last 15 minutes. Find the error in the sockshop logs."

To trigger another failure: `cd ../load-generators && ./bin/demo` (or re-run the CronJob manually).

### Scripts

| Script | Purpose |
|--------|---------|
| `bin/build --kind\|--eks\|--gke` | Build dev images for all services |
| `bin/deploy --kind\|--eks\|--gke` | Deploy Sock Shop to the target cluster |
| `bin/demo-ready` | Verify cluster is ready for the demo |
| `bin/port-forward` | Port-forward front-end to localhost:8080 |

See [DEVELOP.md](./DEVELOP.md) for detailed build/deploy documentation.

# Sock Shop : A Microservice Demo Application

The application is the user-facing part of an online shop that sells socks. It is intended to aid the demonstration and testing of microservice and cloud native technologies.

It is built using [Spring Boot](http://projects.spring.io/spring-boot/) (Java 17), [Go kit](http://gokit.io) (Go 1.24) and [Node.js](https://nodejs.org/) (Node 20), and is packaged in Docker containers.

You can read more about the [application design](./internal-docs/design.md).
