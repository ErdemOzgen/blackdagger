Feature 09: CLI and Configuration Integration
=============================================

Module wiring and startup behavior were updated for log forwarding initialization and monitor startup.

Configuration Surface
---------------------

The following config knobs are integrated end-to-end:

- enabled
- sinkType
- httpURL
- timeoutSec
- includeStepOutput
- queueSize
- maxRetries
- initialBackoffMS
- maxBackoffMS
- monitorEnabled
- monitorHost
- monitorPort
- monitorBasePath
- headers

Related Files
-------------

- cmd/modules.go
- cmd/modules_test.go
- docs/source/config.rst
