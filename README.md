## Tokonoma Demo Setup

Demonstrates log investigation using Tokonoma (toko-mcp) against live Sock Shop services.

### Prerequisites

- `docker`, `kind`, `kubectl`, `helm` v3+
- The [2o](https://github.com/tokonoma-ai/2o) platform repo cloned alongside this repo

### 1. Platform (one-time)

```bash
cd ../2o/platform
./bin/setup.sh          # creates kind cluster "qw"
./bin/deploy.sh         # deploys Quickwit, OTel collector, MCP server
./bin/port-forward.sh   # waits for readiness, starts port-forwards
```

### 2. Build and Deploy Sock Shop

```bash
./bin/build-dev.sh      # builds carts, orders, catalogue, etc., loads into kind
./bin/deploy.sh            # deploys sock-shop
```

### 3. Demo Setup (one command)

```bash
./bin/demo.sh   # build checkout-injector, deploy load-test + CronJob, trigger first failure
```

Background load (load-test) hits front-end. Checkout-fail-injector runs every 15 min and once immediately at startup.

### 4. Verify Readiness

```bash
./bin/demo-ready.sh
```

### 5. Demo Prompt

"Checkout failed for user 57a98d98e4b00679b4a830af in the last 15 minutes. Find the error in the sockshop logs."

To trigger another failure: `./bin/demo.sh` (or re-run the CronJob manually).

### Scripts

| Script | Purpose |
|--------|---------|
| `bin/demo.sh` | One-command demo setup (build + deploy + trigger failure) |
| `bin/build-dev.sh` | Build dev images (carts, orders, catalogue), load into kind |
| `bin/deploy.sh` | Deploy Sock Shop |
| `bin/demo-ready.sh` | Verify cluster is ready for the demo |
| `load-generator-demo/bin/build` | Build checkout-injector image and load into kind |

# Sock Shop : A Microservice Demo Application

The application is the user-facing part of an online shop that sells socks. It is intended to aid the demonstration and testing of microservice and cloud native technologies.

It is built using [Spring Boot](http://projects.spring.io/spring-boot/), [Go kit](http://gokit.io) and [Node.js](https://nodejs.org/) and is packaged in Docker containers.

You can read more about the [application design](./internal-docs/design.md).

## Deployment Platforms

The [deploy folder](./deploy/) contains scripts and instructions to provision the application onto your favourite platform.

**Want to change the application code before deploying?** See [DEVELOP.md](./DEVELOP.md) for kind: build images locally, load into kind, and deploy. 


