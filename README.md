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

"Checkout failed for user 57a98d98e4b00679b4a830af. Find the error in the sockshop logs."

To trigger another failure: `./bin/demo.sh` (or re-run the CronJob manually).

### Scripts

| Script | Purpose |
|--------|---------|
| `bin/demo.sh` | One-command demo setup (build + deploy + trigger failure) |
| `bin/build-dev.sh` | Build dev images (carts, orders, catalogue), load into kind |
| `bin/run-dev.sh` | Deploy Sock Shop with dev overlay |
| `bin/demo-ready.sh` | Verify cluster is ready for the demo |
| `load-generator-demo/bin/build` | Build checkout-injector image and load into kind |

