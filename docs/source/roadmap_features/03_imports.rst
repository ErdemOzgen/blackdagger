Feature 03: DAG Imports
=======================

The imports feature enables modular workflows with nested import handling and validation.

main_workflow.yaml
------------------

.. code-block:: yaml

  name: import-example
  imports:
    - ./common_steps
    - ./notify_steps

  steps:
    - name: run_pipeline
      command: echo "running pipeline"
      depends:
        - validate_data

common_steps.yaml
-----------------

.. code-block:: yaml

  steps:
    - name: prepare_data
      command: echo "preparing data"
    - name: validate_data
      command: echo "validating data"
      depends:
        - prepare_data

notify_steps.yaml
-----------------

.. code-block:: yaml

  steps:
    - name: notify
      command: echo "notifying"
      depends:
        - run_pipeline

Validation Behavior
-------------------

- Nested imports are supported.
- Circular import chains are rejected.
- Duplicate merge paths are rejected.
