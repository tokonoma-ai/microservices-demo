#!/usr/bin/env bash
set -euo pipefail

NAMESPACE="sock-shop"

echo "============================================================================="
echo "Waiting for services and starting port-forwards"
echo "============================================================================="
echo

echo "==> Waiting for front-end..."
kubectl rollout status deployment/front-end -n "${NAMESPACE}" --timeout=120s

echo
echo "============================================================================="
echo "Starting port-forwards (Ctrl+C to stop all)"
echo "============================================================================="
echo
echo "Access Points:"
echo "  Sock Shop:  http://localhost:8080"
echo

cleanup() {
  echo
  echo "==> Stopping port-forwards..."
  kill $(jobs -p) 2>/dev/null || true
  exit 0
}
trap cleanup SIGINT SIGTERM

kubectl port-forward -n "${NAMESPACE}" svc/front-end 8080:80 &

echo "==> Port-forwards running. Press Ctrl+C to stop."
echo

wait
