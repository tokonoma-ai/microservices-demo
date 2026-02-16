#!/usr/bin/env bash
# Inject failures for Tokonoma demo scenarios.
#
# Usage:
#   ./bin/demo-inject-errors.sh scenario1   # Print a user ID for targeted search
#   ./bin/demo-inject-errors.sh scenario2   # Kill carts-db to trigger errors
set -euo pipefail

NAMESPACE="sock-shop"
RESTORE_DELAY="${RESTORE_DELAY:-60}"

usage() {
  echo "Usage: $0 <scenario1|scenario2>"
  echo
  echo "  scenario1  - Look up a valid user ID for targeted transaction trace"
  echo "  scenario2  - Kill carts-db pod to trigger cart errors (restores after ${RESTORE_DELAY}s)"
  exit 1
}

scenario1() {
  echo "==> Scenario 1: Targeted Transaction Trace"
  echo "    Fetching user IDs from the user service..."

  local user_pod
  user_pod=$(kubectl get pods -n "${NAMESPACE}" -l name=user -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)

  if [[ -z "${user_pod}" ]]; then
    echo "ERROR: No user pod found. Is Sock Shop deployed?"
    exit 1
  fi

  # Port-forward to user service temporarily
  kubectl port-forward -n "${NAMESPACE}" svc/user 8085:80 >/dev/null 2>&1 &
  local pf_pid=$!
  sleep 2

  local users_json
  users_json=$(curl -s http://localhost:8085/customers 2>/dev/null || true)
  kill "${pf_pid}" 2>/dev/null || true
  wait "${pf_pid}" 2>/dev/null || true

  if [[ -z "${users_json}" || "${users_json}" == "null" ]]; then
    echo "    Could not fetch users. Using default test user ID."
    echo
    echo "    User ID for demo: 57a98d98e4b00679b4a830ae"
  else
    local user_id
    user_id=$(echo "${users_json}" | python3 -c "
import sys, json
data = json.load(sys.stdin)
customers = data if isinstance(data, list) else data.get('_embedded', {}).get('customer', data.get('customers', []))
if customers:
    print(customers[0].get('id', customers[0].get('_links', {}).get('self', {}).get('href', '').split('/')[-1]))
" 2>/dev/null || echo "57a98d98e4b00679b4a830ae")

    echo
    echo "    User ID for demo: ${user_id}"
  fi

  echo
  echo "    Use this in your demo prompt:"
  echo '    "Search the sockshop index for user ID <user_id> and find any errors in their transactions"'
}

scenario2() {
  echo "==> Scenario 2: Broad Incident Investigation"
  echo "    Killing carts-db pod to simulate database failure..."

  kubectl delete pod -l name=carts-db -n "${NAMESPACE}" --wait=false
  echo "    carts-db pod deleted. Cart operations will now fail."
  echo
  echo "    The load test will continue hitting the carts service, generating ERROR logs."
  echo "    Wait 30-60 seconds for errors to accumulate in Quickwit, then start the demo."
  echo
  echo "    Demo prompt:"
  echo '    "Users are reporting that adding to cart is failing. Investigate the sockshop logs."'
  echo

  if [[ "${RESTORE_DELAY}" -gt 0 ]]; then
    echo "    carts-db will recover automatically (Kubernetes will restart the pod)."
    echo "    To keep it down longer, run: kubectl scale deployment carts-db -n ${NAMESPACE} --replicas=0"
    echo "    To restore: kubectl scale deployment carts-db -n ${NAMESPACE} --replicas=1"
  fi
}

case "${1:-}" in
  scenario1) scenario1 ;;
  scenario2) scenario2 ;;
  *) usage ;;
esac
