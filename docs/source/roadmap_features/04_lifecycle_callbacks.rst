Feature 04: Lifecycle Callback Improvements
===========================================

Lifecycle callbacks were clarified for handlerOn usage with webhook executor.

Full Callback Workflow
----------------------

.. code-block:: yaml

  name: handler-on-webhook-callback

  env:
    - CALLBACK_URL: https://example.org/hooks/workflow-status

  handlerOn:
    success:
      executor:
        type: webhook
        config:
          url: "${CALLBACK_URL}"
          method: "POST"
          timeout: 10
          headers:
            Content-Type: "application/json"
          body: '{"status":"COMPLETED"}'
          successStatusCodes: [200, 202]
          silent: true
    failure:
      executor:
        type: webhook
        config:
          url: "${CALLBACK_URL}"
          method: "POST"
          timeout: 10
          headers:
            Content-Type: "application/json"
          body: '{"status":"FAILED"}'
          successStatusCodes: [200, 202]
          silent: true

  steps:
    - name: check_identity
      command: bash
      script: |
        id
