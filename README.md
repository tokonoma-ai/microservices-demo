[![Build Status](https://travis-ci.org/microservices-demo/microservices-demo.svg?branch=master)](https://travis-ci.org/microservices-demo/microservices-demo)

# DEPRECATED: Sock Shop : A Microservice Demo Application

The application is the user-facing part of an online shop that sells socks. It is intended to aid the demonstration and testing of microservice and cloud native technologies.

It is built using [Spring Boot](http://projects.spring.io/spring-boot/), [Go kit](http://gokit.io) and [Node.js](https://nodejs.org/) and is packaged in Docker containers.

You can read more about the [application design](./internal-docs/design.md).

## Deployment Platforms

The [deploy folder](./deploy/) contains scripts and instructions to provision the application onto your favourite platform.

**Want to change the application code before deploying?** See [DEVELOP.md](./DEVELOP.md) for kind: clone service repos, build images locally, load into kind, and deploy. 

Please let us know if there is a platform that you would like to see supported.

## Bugs, Feature Requests and Contributing

We'd love to see community contributions. We like to keep it simple and use Github issues to track bugs and feature requests and pull requests to manage contributions. See the [contribution information](.github/CONTRIBUTING.md) for more information.

## Screenshot

![Sock Shop frontend](https://github.com/microservices-demo/microservices-demo.github.io/raw/master/assets/sockshop-frontend.png)

## Visualizing the application

Use [Weave Scope](http://weave.works/products/weave-scope/) or [Weave Cloud](http://cloud.weave.works/) to visualize the application once it's running in the selected [target platform](./deploy/).

![Sock Shop in Weave Scope](https://github.com/microservices-demo/microservices-demo.github.io/raw/master/assets/sockshop-scope.png)

## Tokonoma Demo Setup

Demonstrates log investigation using Tokonoma (toko-mcp) against live Sock Shop services.

### Prerequisites

- `docker`, `kind`, `kubectl`, `helm` v3+
- The [2o](https://github.com/khou/2o) platform repo cloned alongside this repo

### 1. Platform (one-time)

```bash
cd ../2o/platform
./bin/setup.sh          # creates kind cluster "qw"
./bin/deploy.sh         # deploys Quickwit, OTel collector, MCP server
./bin/port-forward.sh   # waits for readiness, starts port-forwards
```

### 2. Build and Deploy Sock Shop

```bash
./bin/build-dev.sh      # builds carts, orders, catalogue; loads into kind
./bin/run-dev.sh        # deploys sock-shop with dev overlay
```

### 3. Start Load Test

```bash
kubectl apply -f deploy/kubernetes/manifests-loadtest/
```

### 4. Verify Readiness

```bash
./bin/demo-ready.sh
```

### 5. Run Demo Scenarios

**Scenario 1 -- Targeted Transaction Trace**

```bash
./bin/demo-inject-errors.sh scenario1   # prints a user ID to search for
```

Prompt: "Search the sockshop index for user ID `<user_id>` and find any errors in their transactions"

**Scenario 2 -- Broad Incident Investigation**

```bash
./bin/demo-inject-errors.sh scenario2   # kills carts-db, wait 30-60s for errors
```

Prompt: "Users are reporting that adding to cart is failing. Investigate the sockshop logs."

### Scripts

| Script | Purpose |
|--------|---------|
| `bin/build-dev.sh` | Build dev images (carts, orders, catalogue), load into kind |
| `bin/run-dev.sh` | Deploy Sock Shop with dev overlay |
| `bin/demo-ready.sh` | Verify cluster is ready for the demo |
| `bin/demo-inject-errors.sh` | Inject failures for demo scenarios |

