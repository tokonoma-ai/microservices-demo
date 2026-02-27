#!/usr/bin/env bash
# Build dev images for all buildable services and load them into the kind cluster.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="${SCRIPT_DIR}/.."
KIND_CLUSTER="${KIND_CLUSTER:-qw}"

echo "============================================================================="
echo "Building dev images for Sock Shop"
echo "============================================================================="

### Carts (Java 17) ###########################################################

echo
echo "==> Building carts JAR (Java 17)..."
docker run --rm \
  -v "${REPO_ROOT}/carts":/build \
  -w /build \
  maven:3-eclipse-temurin-17 \
  mvn package -DskipTests -q

echo "==> Building carts Docker image..."
docker build -t weaveworksdemos/carts:dev "${REPO_ROOT}/carts"

### Orders (Java 17) ###########################################################

echo
echo "==> Building orders JAR (Java 17)..."
docker run --rm \
  -v "${REPO_ROOT}/orders":/build \
  -w /build \
  maven:3-eclipse-temurin-17 \
  mvn package -DskipTests -q

echo "==> Building orders Docker image..."
docker build -t weaveworksdemos/orders:dev "${REPO_ROOT}/orders"

### Shipping (Java 17) #########################################################

echo
echo "==> Building shipping JAR (Java 17)..."
docker run --rm \
  -v "${REPO_ROOT}/shipping":/build \
  -w /build \
  maven:3-eclipse-temurin-17 \
  mvn package -DskipTests -q

echo "==> Building shipping Docker image..."
docker build -t weaveworksdemos/shipping:dev "${REPO_ROOT}/shipping"

### Queue-master (Java 17) #####################################################

echo
echo "==> Building queue-master JAR (Java 17)..."
docker run --rm \
  -v "${REPO_ROOT}/queue-master":/build \
  -w /build \
  maven:3-eclipse-temurin-17 \
  mvn package -DskipTests -q

echo "==> Building queue-master Docker image..."
docker build -t weaveworksdemos/queue-master:dev "${REPO_ROOT}/queue-master"

### Catalogue (Go 1.24, multi-stage) ##########################################

echo
echo "==> Building catalogue Docker image (multi-stage)..."
docker build -t weaveworksdemos/catalogue:dev \
  -f "${REPO_ROOT}/catalogue/docker/catalogue/Dockerfile" \
  "${REPO_ROOT}/catalogue"

### Payment (Go 1.24, multi-stage) #############################################

echo
echo "==> Building payment Docker image (multi-stage)..."
docker build -t weaveworksdemos/payment:dev \
  -f "${REPO_ROOT}/payment/docker/payment/Dockerfile" \
  "${REPO_ROOT}/payment"

### User (Go 1.24, multi-stage) ################################################

echo
echo "==> Building user Docker image (multi-stage)..."
docker build -t weaveworksdemos/user:dev \
  -f "${REPO_ROOT}/user/docker/user/Dockerfile" \
  "${REPO_ROOT}/user"

### Front-end (Node.js 20) #####################################################

echo
echo "==> Building front-end Docker image..."
docker build -t weaveworksdemos/front-end:dev "${REPO_ROOT}/front-end"

### Load into kind #############################################################

echo
echo "==> Loading images into kind cluster '${KIND_CLUSTER}'..."
kind load docker-image weaveworksdemos/carts:dev        --name "${KIND_CLUSTER}"
kind load docker-image weaveworksdemos/orders:dev        --name "${KIND_CLUSTER}"
kind load docker-image weaveworksdemos/shipping:dev      --name "${KIND_CLUSTER}"
kind load docker-image weaveworksdemos/queue-master:dev  --name "${KIND_CLUSTER}"
kind load docker-image weaveworksdemos/catalogue:dev     --name "${KIND_CLUSTER}"
kind load docker-image weaveworksdemos/payment:dev       --name "${KIND_CLUSTER}"
kind load docker-image weaveworksdemos/user:dev          --name "${KIND_CLUSTER}"
kind load docker-image weaveworksdemos/front-end:dev     --name "${KIND_CLUSTER}"

echo
echo "============================================================================="
echo "Build complete. Images loaded into kind cluster '${KIND_CLUSTER}'."
echo "Run ./bin/deploy.sh to deploy."
echo "============================================================================="
