## Tokonoma Demo Setup

Demonstrates log investigation using Tokonoma (toko-mcp) against live Sock Shop services.

### Prerequisites

- `docker`, `kubectl`, `helm` v3+
- `kind` (for local clusters) or `awscli` + `eksctl` (for EKS)
- The [2o](https://github.com/tokonoma-ai/2o) platform repo cloned alongside this repo

### 1. Platform (one-time)

```bash
cd ../2o/platform
./bin/setup          # creates kind cluster "qw"
./bin/deploy         # deploys Quickwit, OTel collector, MCP server
./bin/port-forward   # waits for readiness, starts port-forwards
```

### 2. Build and Deploy Sock Shop

All services are built from source. Specify your target cluster:

```bash
# kind (local)
./bin/build --kind && ./bin/deploy --kind

# EKS
./bin/build --eks && ./bin/deploy --eks
```

### 3. Demo Setup (load generators)

```bash
cd ../load-generators
./bin/demo   # build checkout-injector, deploy load-test + CronJob, trigger first failure
```

Background load (load-test) hits front-end. Checkout-fail-injector runs every 15 min and once immediately at startup.

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
| `bin/build --kind\|--eks` | Build dev images for all services |
| `bin/deploy --kind\|--eks` | Deploy Sock Shop to the target cluster |
| `bin/demo-ready` | Verify cluster is ready for the demo |
| `bin/port-forward` | Port-forward front-end to localhost:8080 |

See [DEVELOP.md](./DEVELOP.md) for detailed build/deploy documentation.

# Sock Shop : A Microservice Demo Application

The application is the user-facing part of an online shop that sells socks. It is intended to aid the demonstration and testing of microservice and cloud native technologies.

It is built using [Spring Boot](http://projects.spring.io/spring-boot/), [Go kit](http://gokit.io) and [Node.js](https://nodejs.org/) and is packaged in Docker containers.

You can read more about the [application design](./internal-docs/design.md).
