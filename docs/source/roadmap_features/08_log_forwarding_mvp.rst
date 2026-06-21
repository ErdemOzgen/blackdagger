Feature 08: Centralized Log Forwarding MVP
==========================================

The MVP introduces asynchronous forwarding, retry behavior, and monitor endpoints.

Core Capabilities
-----------------

- Async log forwarding pipeline.
- HTTP sink implementation.
- Buffering and retry-oriented behavior.
- Local monitor endpoints for health and metrics.
- Metrics in JSON and Prometheus text formats.

Full config.yaml Example
------------------------

.. code-block:: yaml

  # config.yaml example for centralized log forwarding MVP
  # monitor endpoints:
  # - http://127.0.0.1:8091/log-forwarding/health
  # - http://127.0.0.1:8091/log-forwarding/metrics (JSON)
  # - http://127.0.0.1:8091/log-forwarding/metrics/prometheus
  logForwarding:
    enabled: true
    sinkType: http
    httpURL: https://logs.example.com/ingest
    timeoutSec: 5
    includeStepOutput: false
    queueSize: 256
    maxRetries: 3
    initialBackoffMS: 100
    maxBackoffMS: 2000
    monitorEnabled: true
    monitorHost: 127.0.0.1
    monitorPort: 8091
    monitorBasePath: /log-forwarding
    headers:
      Authorization: "Bearer ${LOG_FORWARDING_TOKEN}"

Full Metrics Scrape Workflow
----------------------------

.. code-block:: yaml

  name: scrape-log-forwarding-metrics

  steps:
    - name: scrape_prometheus_metrics
      command: curl
      args:
        - -fsS
        - http://127.0.0.1:8091/log-forwarding/metrics/prometheus
