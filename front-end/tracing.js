'use strict';

/**
 * OpenTelemetry: Zipkin export + B3 propagation (multi-header X-B3-*), aligned with Java/Go services.
 * Must be loaded before express, request, or other instrumented modules (see server.js).
 */
if (process.env.OTEL_SDK_DISABLED === 'true') {
  module.exports = { started: false };
  return;
}

var NodeSDK = require('@opentelemetry/sdk-node').NodeSDK;
var getNodeAutoInstrumentations =
  require('@opentelemetry/auto-instrumentations-node').getNodeAutoInstrumentations;
var ZipkinExporter = require('@opentelemetry/exporter-zipkin').ZipkinExporter;
var B3Propagator = require('@opentelemetry/propagator-b3').B3Propagator;
var B3InjectEncoding = require('@opentelemetry/propagator-b3').B3InjectEncoding;

function zipkinEndpoint() {
  if (process.env.OTEL_EXPORTER_ZIPKIN_ENDPOINT) {
    return process.env.OTEL_EXPORTER_ZIPKIN_ENDPOINT;
  }
  var z = process.env.ZIPKIN;
  if (z) {
    return z.replace(/\/api\/v1\/spans\/?$/i, '/api/v2/spans');
  }
  return 'http://127.0.0.1:9411/api/v2/spans';
}

var sdk = new NodeSDK({
  serviceName: process.env.OTEL_SERVICE_NAME || 'front-end',
  traceExporter: new ZipkinExporter({
    url: zipkinEndpoint(),
  }),
  textMapPropagator: new B3Propagator({
    injectEncoding: B3InjectEncoding.MULTI_HEADER,
  }),
  instrumentations: [
    getNodeAutoInstrumentations({
      '@opentelemetry/instrumentation-fs': { enabled: false },
      '@opentelemetry/instrumentation-dns': { enabled: false },
    }),
  ],
});

sdk.start();

function shutdown() {
  sdk
    .shutdown()
    .catch(function () {})
    .finally(function () {
      process.exit(0);
    });
}

process.once('SIGTERM', shutdown);
process.once('SIGINT', shutdown);

module.exports = { started: true };
