.. _Examples:

Examples
========

.. contents::
    :local:

Hello World
------------

.. code-block:: yaml

  name: hello world
  steps:
    - name: s1
      command: echo hello world
    - name: s2
      command: echo done!
      depends:
        - s1


Conditional Steps
------------------

.. code-block:: yaml

  params: foo
  steps:
    - name: step1
      command: echo start
    - name: foo
      command: echo foo
      depends:
        - step1
      preconditions:
        - condition: "$1"
          expected: foo
    - name: bar
      command: echo bar
      depends:
        - step1
      preconditions:
        - condition: "$1"
          expected: bar

.. image:: https://raw.githubusercontent.com/ErdemOzgen/blackdagger/main/examples/images/conditional.png


File Output
------------

.. code-block:: yaml

  steps:
    - name: write hello to '/tmp/hello.txt'
      command: echo hello
      stdout: /tmp/hello.txt

Passing Output to Next Step
---------------------------

.. code-block:: yaml

  steps:
    - name: pass 'hello'
      command: echo hello
      output: OUT1
    - name: output 'hello world'
      command: bash
      script: |
        echo $OUT1 world
      depends:
        - pass 'hello'

Running a Docker Container
--------------------------

.. code-block:: yaml

  steps:
    - name: deno_hello_world
      executor: 
        type: docker
        config:
          image: "denoland/deno:1.10.3"
          host:
            autoRemove: true
      command: run https://examples.deno.land/hello-world.ts

See :ref:`docker executor` for more details.

Sending HTTP Requests
---------------------

.. code-block:: yaml

  steps:
    - name: get fake json data
      executor: http
      command: GET https://jsonplaceholder.typicode.com/comments
      script: |
        {
          "timeout": 10,
          "headers": {},
          "query": {
            "postId": "1"
          },
          "body": ""
        }

Sending a Webhook
-----------------

.. code-block:: yaml

  steps:
    - name: trigger_deployment_webhook
      executor:
        type: webhook
        config:
          url: "https://example.org/hooks/deploy"
          method: "POST"
          timeout: 10
          headers:
            Authorization: "Bearer ${WEBHOOK_TOKEN}"
            Content-Type: "application/json"
          query:
            env: "prod"
          body: '{"service":"api","version":"1.2.3"}'
          successStatusCodes: [200, 202]
          silent: true

Handler Callbacks with Webhooks
-------------------------------

Use lifecycle hooks to notify external systems when a workflow succeeds or
fails.

.. code-block:: yaml

  env:
    - CALLBACK_URL: https://example.org/hooks/workflow-status

  handlerOn:
    success:
      executor:
        type: webhook
        config:
          url: "${CALLBACK_URL}"
          method: "POST"
          headers:
            Content-Type: "application/json"
          body: '{"status":"COMPLETED"}'
          silent: true
    failure:
      executor:
        type: webhook
        config:
          url: "${CALLBACK_URL}"
          method: "POST"
          headers:
            Content-Type: "application/json"
          body: '{"status":"FAILED"}'
          silent: true

  steps:
    - name: check_identity
      command: bash
      script: |
        id

Running Terraform
-----------------

.. code-block:: yaml

  steps:
    - name: terraform_apply
      executor:
        type: terraform
        config:
          binary: terraform
          workingDir: ./infra
          init: true
          initArgs:
            - -upgrade
          subcommand: apply
          varFiles:
            - env/prod.tfvars
          vars:
            image_tag: "1.2.3"
            region: "us-east-1"
          targets:
            - module.app
          autoApprove: true
          env:
            TF_IN_AUTOMATION: "true"

Running Ansible Playbook
------------------------

.. code-block:: yaml

  steps:
    - name: run_ansible_playbook
      executor:
        type: ansible
        config:
          binary: ansible-playbook
          playbook: deploy/site.yml
          inventory: inventory/prod.ini
          tags:
            - deploy
            - app
          extraVars:
            image_tag: "1.2.3"
            environment: "prod"
          become: true
          forks: 20
          diff: true
          check: false
          env:
            ANSIBLE_STDOUT_CALLBACK: yaml

Querying JSON Data with jq
----------------------------

.. code-block:: yaml

  steps:
    - name: run query
      executor: jq
      command: '{(.id): .["10"].b}'
      script: |
        {"id": "sample", "10": {"b": 42}}

Expected Output:

.. code-block:: json

    {
        "sample": 42
    }


Formatting JSON Data with jq
----------------------------

.. code-block:: yaml

  steps:
    - name: format json
      executor: jq
      script: |
        {"id": "sample", "10": {"b": 42}}

Expected Output:

.. code-block:: json

    {
        "10": {
            "b": 42
        },
        "id": "sample"
    }


Outputting Raw Values with jq
-----------------------------

.. code-block:: yaml

  steps:
    - name: output raw value
      executor:
        type: jq
        config:
          raw: true
      command: '.id'
      script: |
        {"id": "sample", "10": {"b": 42}}

Expected Output:

.. code-block:: sh

    sample


Sending Email Notifications
---------------------------

.. image:: https://raw.githubusercontent.com/ErdemOzgen/blackdagger/main/examples/images/email.png

.. code-block:: yaml

  steps:
    - name: Sending Email on Finish or Error
      command: echo "hello world"

  mailOn:
    failure: true
    success: true

  smtp:
    host: "smtp.foo.bar"
    port: "587"
    username: "<username>"
    password: "<password>"
  errorMail:
    from: "foo@bar.com"
    to: "foo@bar.com"
    prefix: "[Error]"
    attachLogs: true
  infoMail:
    from: "foo@bar.com"
    to: "foo@bar.com"
    prefix: "[Info]"
    attachLogs: true


Sending Email
-------------

.. code-block:: yaml

  smtp:
    host: "smtp.foo.bar"
    port: "587"
    username: "<username>"
    password: "<password>"

  steps:
    - name: step1
      executor:
        type: mail
        config:
          to: <to address>
          from: <from address>
          subject: "Sample Email"
          message: |
            Hello world


Customizing Signal Handling on Stop
-----------------------------------

.. code-block:: yaml

  steps:
    - name: step1
      command: bash
      script: |
        for s in {1..64}; do trap "echo trap $s" $s; done
        sleep 60
      signalOnStop: "SIGINT"

Importing DAG Files
-------------------

Use ``imports`` to build modular workflows from reusable YAML files.

.. code-block:: yaml

  # main_workflow.yaml
  name: import-example
  imports:
    - ./common_steps
    - ./notify_steps

  steps:
    - name: run_pipeline
      command: echo "running pipeline"
      depends:
        - validate_data

.. code-block:: yaml

  # common_steps.yaml
  steps:
    - name: prepare_data
      command: echo "preparing data"
    - name: validate_data
      command: echo "validating data"
      depends:
        - prepare_data

.. code-block:: yaml

  # notify_steps.yaml
  steps:
    - name: notify
      command: echo "notifying"
      depends:
        - run_pipeline

Scraping Log Forwarding Metrics
-------------------------------

Use this example to scrape Prometheus-formatted log forwarding metrics from
the local monitor endpoint.

.. code-block:: yaml

  name: scrape-log-forwarding-metrics

  steps:
    - name: scrape_prometheus_metrics
      command: curl
      args:
        - -fsS
        - http://127.0.0.1:8091/log-forwarding/metrics/prometheus

