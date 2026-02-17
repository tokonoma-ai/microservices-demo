#!/usr/bin/env bash
# Build dev images for carts, orders, catalogue and load them into the kind cluster.
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

### Orders (Java 8) ############################################################

echo
echo "==> Building orders JAR (Java 8)..."
docker run --rm \
  -v "${REPO_ROOT}/orders":/build \
  -w /build \
  maven:3-eclipse-temurin-8 \
  mvn package -DskipTests -q

echo "==> Building orders Docker image..."
docker build -t weaveworksdemos/orders:dev "${REPO_ROOT}/orders"

### Catalogue (Go 1.24, multi-stage) ##########################################

echo
echo "==> Building catalogue Docker image (multi-stage)..."
docker build -t weaveworksdemos/catalogue:dev \
  -f "${REPO_ROOT}/catalogue/docker/catalogue/Dockerfile" \
  "${REPO_ROOT}/catalogue"

### Load into kind #############################################################

echo
echo "==> Loading images into kind cluster '${KIND_CLUSTER}'..."
kind load docker-image weaveworksdemos/carts:dev     --name "${KIND_CLUSTER}"
kind load docker-image weaveworksdemos/orders:dev     --name "${KIND_CLUSTER}"
kind load docker-image weaveworksdemos/catalogue:dev  --name "${KIND_CLUSTER}"

echo
echo "============================================================================="
echo "Build complete. Images loaded into kind cluster '${KIND_CLUSTER}'."
echo "Run bin/run-dev.sh to deploy."
echo "============================================================================="
