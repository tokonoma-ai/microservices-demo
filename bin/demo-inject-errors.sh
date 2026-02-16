#!/usr/bin/env bash
# Trigger checkout failure injection for Tokonoma demo.
#
# Usage:
#   ./bin/demo-inject-errors.sh trigger   # Run checkout-fail injector immediately
#   ./bin/demo-inject-errors.sh scenario2 # Kill carts-db for cart-failure demo (optional)
set -euo pipefail

LOADTEST_NS="loadtest"
NAMESPACE="sock-shop"

usage() {
  echo "Usage: $0 <trigger|scenario2>"
  echo
  echo "  trigger   - Run checkout-fail injector now (creates one failure)"
  echo "  scenario2 - Kill carts-db to trigger cart errors (optional, different demo)"
  exit 1
}

trigger() {
  echo "==> Triggering checkout failure now..."
  kubectl delete job checkout-fail-now -n "${LOADTEST_NS}" 2>/dev/null || true
  kubectl create job checkout-fail-now --from=cronjob/checkout-fail-injector -n "${LOADTEST_NS}"
  echo "Job started. Check logs: kubectl logs job/checkout-fail-now -n loadtest -f"
}

scenario2() {
  echo "==> Killing carts-db (cart failure demo)..."
  kubectl delete pod -l name=carts-db -n "${NAMESPACE}" --wait=false
  echo "carts-db deleted. Cart ops will fail. Demo prompt:"
  echo '"Users are reporting that adding to cart is failing. Investigate the sockshop logs."'
}

case "${1:-}" in
  trigger) trigger ;;
  scenario2) scenario2 ;;
  *) usage ;;
esac
