#!/usr/bin/env bash
# One-command demo setup: build, deploy load + checkout-injector, run injector immediately.
# Run after platform and sock-shop are up.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

echo "==> Tokonoma Demo Setup"
echo

# 1. Build checkout-injector and load-generator-demo, load into kind
echo "==> Building images..."
"${REPO_ROOT}/load-generator-demo/bin/build"

# 2. Load image into kind (ensure available before Job runs)
KIND_CLUSTER="${KIND_CLUSTER:-qw}"
echo "==> Loading checkout-injector into kind..."
kind load docker-image weaveworksdemos/checkout-injector:latest --name "${KIND_CLUSTER}"

# 3. Deploy load-test (background load) + demo manifests (CronJob)
echo "==> Deploying load-test and checkout-injector..."
kubectl apply -f "${REPO_ROOT}/deploy/kubernetes/manifests-loadtest/"
kubectl apply -f "${REPO_ROOT}/deploy/kubernetes/manifests-loadtest-demo/loadtest-demo-dep.yaml"

# 4. Run injector immediately (don't wait 15 min for first CronJob)
echo "==> Triggering checkout failure now..."
kubectl delete job checkout-fail-now -n loadtest 2>/dev/null || true
kubectl create job checkout-fail-now --from=cronjob/checkout-fail-injector -n loadtest

echo
echo "==> Verify: ./bin/demo-ready.sh"
echo "==> Demo prompt: \"Checkout failed for user 57a98d98e4b00679b4a830af. Find the error in the sockshop logs.\""
