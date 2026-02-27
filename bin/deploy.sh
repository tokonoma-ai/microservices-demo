#!/bin/bash
#
# Deploy Sock Shop to Kubernetes using Kustomize.
# All app services use locally-built :dev images.
#
# Usage:
#   ./bin/deploy.sh
#
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
KUSTOMIZE_ROOT="${REPO_ROOT}/deploy/kubernetes"
NAMESPACE="sock-shop"

echo "Deploying Sock Shop from microservices-demo (${REPO_ROOT})..."
echo ""

echo "Applying with Kustomize (kubectl apply -k ${KUSTOMIZE_ROOT})..."
kubectl apply -k "${KUSTOMIZE_ROOT}"

echo ""
echo "Waiting for deployments to roll out..."
for dep in carts catalogue orders front-end payment user shipping queue-master load-generator; do
  echo "  Waiting for ${dep}..."
  kubectl rollout status deployment/"${dep}" -n "${NAMESPACE}" --timeout=180s 2>/dev/null || \
    echo "  WARNING: ${dep} did not become ready within 180s"
done

echo ""
echo "Deployment complete."
echo ""
echo "Current status:"
kubectl get pods -n "${NAMESPACE}"
echo ""
echo "To access the front-end:"
echo "  NodePort: http://<node-ip>:30001"
echo "  Or: ./bin/port-forward.sh"
echo "  Then open http://localhost:8080"
