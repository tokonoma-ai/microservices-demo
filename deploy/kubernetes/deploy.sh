#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
NAMESPACE="sock-shop"
MONITORING_NAMESPACE="monitoring"
JAEGER_NAMESPACE="jaeger"

echo "Deploying Sock Shop with tracing and monitoring..."

# Create namespaces
echo "Creating namespaces..."
kubectl create namespace ${NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -
kubectl create namespace ${MONITORING_NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -
kubectl create namespace ${JAEGER_NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -

# Deploy Jaeger first (services need it for tracing)
echo "Deploying Jaeger..."
kubectl apply -f ${SCRIPT_DIR}/manifests-jaeger/jaeger.yaml

# Deploy Sock Shop using manifests
echo "Deploying Sock Shop services from manifests..."
kubectl apply -f ${SCRIPT_DIR}/manifests/

# Override services with prepared jaeger-enabled manifests (these have ZIPKIN configured)
echo "Applying tracing-enabled service manifests..."
kubectl apply -f ${SCRIPT_DIR}/manifests-jaeger/user-dep.yaml
kubectl apply -f ${SCRIPT_DIR}/manifests-jaeger/payment-dep.yaml
kubectl apply -f ${SCRIPT_DIR}/manifests-jaeger/catalogue-dep.yaml

# Shipping already has ZIPKIN in the manifest, but tracing is disabled in JAVA_OPTS
# Enable tracing by removing the -Dspring.zipkin.enabled=false flag
echo "Enabling tracing for shipping service..."
kubectl patch deployment shipping -n ${NAMESPACE} --type='json' \
  -p='[{"op": "replace", "path": "/spec/template/spec/containers/0/env", "value": [{"name": "ZIPKIN", "value": "zipkin.jaeger.svc.cluster.local"}, {"name": "JAVA_OPTS", "value": "-Xms64m -Xmx128m -XX:+UseG1GC -Djava.security.egd=file:/dev/urandom"}]}]' 2>/dev/null || \
echo "Note: Shipping service may need manual tracing configuration"

# Deploy Prometheus and Grafana
echo "Deploying monitoring stack (Prometheus + Grafana)..."
kubectl apply -f ${SCRIPT_DIR}/manifests-monitoring/

# Import Grafana dashboards (delete old job first if exists)
echo "Importing Grafana dashboards..."
kubectl delete job grafana-import-dashboards -n ${MONITORING_NAMESPACE} --ignore-not-found=true
kubectl apply -f ${SCRIPT_DIR}/manifests-monitoring/23-grafana-import-dash-batch.yaml

echo ""
echo "Deployment complete!"
echo ""
echo "Services:"
echo "  - Sock Shop frontend: NodePort 30001 (port 8079)"
echo "  - Grafana: NodePort 31300 (port 3000)"
echo "  - Prometheus: NodePort 31090 (port 9090)"
echo "  - Jaeger UI: LoadBalancer service"
echo ""
echo "Use ./port-forward.sh to port forward the UIs locally"
echo ""
echo "Checking service status..."
kubectl get pods -n ${NAMESPACE}
kubectl get pods -n ${MONITORING_NAMESPACE}
kubectl get pods -n ${JAEGER_NAMESPACE}
