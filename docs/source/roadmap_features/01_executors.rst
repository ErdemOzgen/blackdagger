Feature 01: New Executors
=========================

This update adds three executors for external integrations and infrastructure workflows.

Webhook Executor
----------------

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

Terraform Executor
------------------

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

Ansible Executor
----------------

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
