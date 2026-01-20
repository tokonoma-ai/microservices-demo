#!/bin/bash

set -e

NAMESPACE="sock-shop"
MONITORING_NAMESPACE="monitoring"
JAEGER_NAMESPACE="jaeger"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}Setting up port forwards...${NC}"
echo ""

# Function to get pod name (wait for pod to be ready)
get_pod() {
    local namespace=$1
    local selector=$2
    local max_wait=${3:-60}
    local elapsed=0
    
    while [ ${elapsed} -lt ${max_wait} ]; do
        local pod=$(kubectl get pod -n ${namespace} -l ${selector} -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)
        if [ -n "${pod}" ]; then
            local status=$(kubectl get pod -n ${namespace} ${pod} -o jsonpath='{.status.phase}' 2>/dev/null)
            if [ "${status}" = "Running" ]; then
                echo "${pod}"
                return 0
            fi
        fi
        sleep 2
        elapsed=$((elapsed + 2))
    done
    return 1
}

# Cleanup function
cleanup() {
    echo ""
    echo -e "${YELLOW}Stopping port forwards...${NC}"
    if [ -f /tmp/sock-shop-port-forwards.pid ]; then
        local pids=$(cat /tmp/sock-shop-port-forwards.pid)
        for pid in ${pids}; do
            [ -n "${pid}" ] && kill ${pid} 2>/dev/null || true
        done
        rm /tmp/sock-shop-port-forwards.pid
    fi
    exit
}

trap cleanup INT TERM

# Port forward Sock Shop frontend
echo -e "${BLUE}Waiting for frontend pod...${NC}"
kubectl wait --for=condition=ready pod -l name=front-end -n ${NAMESPACE} --timeout=120s || true
FRONTEND_POD=$(get_pod ${NAMESPACE} "name=front-end" 10)
if [ -n "${FRONTEND_POD}" ]; then
    echo -e "${GREEN}Forwarding Sock Shop frontend (http://localhost:8080)${NC}"
    kubectl port-forward -n ${NAMESPACE} ${FRONTEND_POD} 8080:8079 > /tmp/frontend-port-forward.log 2>&1 &
    FRONTEND_PID=$!
    echo "  Pod: ${FRONTEND_POD}, PID: ${FRONTEND_PID}"
    sleep 1
else
    echo -e "${RED}Error: Frontend pod not found or not ready${NC}"
    FRONTEND_PID=""
fi

# Port forward Grafana (to port 3001 as requested)
echo -e "${BLUE}Waiting for Grafana pod...${NC}"
kubectl wait --for=condition=ready pod -l app=grafana -n ${MONITORING_NAMESPACE} --timeout=120s || true
GRAFANA_POD=$(get_pod ${MONITORING_NAMESPACE} "app=grafana,component=core" 10)
if [ -n "${GRAFANA_POD}" ]; then
    echo -e "${GREEN}Forwarding Grafana (http://localhost:3001)${NC}"
    kubectl port-forward -n ${MONITORING_NAMESPACE} ${GRAFANA_POD} 3001:3000 > /tmp/grafana-port-forward.log 2>&1 &
    GRAFANA_PID=$!
    echo "  Pod: ${GRAFANA_POD}, PID: ${GRAFANA_PID}"
    sleep 1
else
    echo -e "${RED}Error: Grafana pod not found or not ready${NC}"
    GRAFANA_PID=""
fi

# Port forward Jaeger UI
echo -e "${BLUE}Waiting for Jaeger pod...${NC}"
kubectl wait --for=condition=ready pod -l jaeger-infra=jaeger-pod -n ${JAEGER_NAMESPACE} --timeout=120s || true
JAEGER_POD=$(get_pod ${JAEGER_NAMESPACE} "jaeger-infra=jaeger-pod" 10)
if [ -n "${JAEGER_POD}" ]; then
    echo -e "${GREEN}Forwarding Jaeger UI (http://localhost:16686)${NC}"
    kubectl port-forward -n ${JAEGER_NAMESPACE} ${JAEGER_POD} 16686:16686 > /tmp/jaeger-port-forward.log 2>&1 &
    JAEGER_PID=$!
    echo "  Pod: ${JAEGER_POD}, PID: ${JAEGER_PID}"
    sleep 1
else
    echo -e "${RED}Error: Jaeger pod not found or not ready${NC}"
    JAEGER_PID=""
fi

# Save PIDs to file for easy cleanup
echo "${FRONTEND_PID} ${GRAFANA_PID} ${JAEGER_PID}" > /tmp/sock-shop-port-forwards.pid

echo ""
echo -e "${GREEN}Port forwards are running!${NC}"
echo ""
echo "Access points:"
echo "  - Sock Shop:    http://localhost:8080"
echo "  - Grafana:      http://localhost:3001 (admin/admin)"
echo "  - Jaeger UI:    http://localhost:16686"
echo ""
echo "Logs are in /tmp/*-port-forward.log"
echo "Press Ctrl+C to stop all port forwards"
echo ""

# Wait for all port forwards
wait
