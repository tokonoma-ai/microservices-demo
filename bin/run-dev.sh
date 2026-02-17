#!/usr/bin/env bash
# Deploy Sock Shop using the dev overlay (locally built images for carts, orders, catalogue).
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="${SCRIPT_DIR}/.."

echo "==> Deploying Sock Shop (dev overlay)..."
kubectl apply -k "${REPO_ROOT}/deploy/kubernetes/overlays/dev"

echo "==> Waiting for deployments to roll out..."
for dep in carts catalogue orders front-end payment user shipping queue-master; do
  echo "    Waiting for ${dep}..."
  kubectl rollout status deployment/"${dep}" -n sock-shop --timeout=180s 2>/dev/null || \
    echo "    WARNING: ${dep} did not become ready within 180s"
done

echo
echo "Sock Shop deployed. Pods:"
kubectl get pods -n sock-shop
