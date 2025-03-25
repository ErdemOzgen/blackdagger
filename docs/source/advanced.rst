Advanced Examples
========

.. contents::
    :local:

Running sub workflows
~~~~~~~~~~~~~~~~~~~~~~~~
Organize complex workflows using sub workflow:

.. code-block:: yaml

  steps:
    - name: sub workflow
      run: sub_workflow
      params: "FOO=BAR"

The result of the sub workflow will be available from the standard output of the sub workflow in JSON format.

Example:

.. code-block:: json

  {
    "name": "sub_workflow"
    "params": "FOO=BAR",
    "outputs": {
      "RESULT": "ok",
    }
  }

You can access the output of the sub workflow using the `output` field:

.. code-block:: yaml

  steps:
    - name: sub workflow
      run: sub_workflow
      params: "FOO=BAR"
      output: SUB_RESULT

    - name: use sub workflow output
      command: echo $SUB_RESULT
      depends:
        - sub workflow

Command Substitution
~~~~~~~~~~~~~~~~~
Use command output in configurations:

.. code-block:: yaml

  env:
    TODAY: "`date '+%Y%m%d'`"
  steps:
    - name: use date
      command: "echo hello, today is ${TODAY}"

Lifecycle Hooks
~~~~~~~~~~~~~
React to DAG state changes:

.. code-block:: yaml

  handlerOn:
    success:
      command: echo "succeeded!"
    cancel:
      command: echo "cancelled!"
    failure:
      command: echo "failed!"
    exit:
      command: echo "exited!"
  steps:
    - name: main task
      command: echo hello

Repeat Steps
~~~~~~~~~~
Execute steps periodically:

.. code-block:: yaml

  steps:
    - name: repeating task
      command: main.sh
      repeatPolicy:
        repeat: true
        intervalSec: 60