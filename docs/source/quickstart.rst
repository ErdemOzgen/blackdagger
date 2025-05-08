Quick Start Guide
=================

.. contents::
    :local:

1. Launch the Web UI
---------------------

Start the server with the command ``blackdagger start-all`` or ``blackdagger server`` and browse to http://127.0.0.1:8080 to explore the Web UI.

*Note: The server will be started on port* ``8080`` *by default. You can change the port by passing* ``--port`` *option. See* :ref:`Host and Port Configuration` *for more details.*

2. Create a New DAG
-------------------

Navigate to the DAG List page by clicking the menu in the left panel of the Web UI. Then create a DAG by clicking the ``New`` button at the top of the page. Enter a name for the DAG (e.g. ``example-dag``) and click the ``Create`` button.

*Note: DAG (YAML) files will be placed in ~/.config/blackdagger/dags by default. See* :ref:`Configuration Options` *for more details.*

3. Edit the DAG
---------------

Go to the ``SPEC`` Tab and hit the ``Edit`` button. Copy & Paste the following YAML code into the editor.

.. code-block:: yaml

    schedule: "* * * * *" # Run the DAG every minute
    steps:
      - name: s1
        command: echo Hello blackdagger
      - name: s2
        command: echo done!
        depends:
          - s1

4. Execute the DAG
-------------------

You can execute the example by pressing the `Start` button. You can see "Hello blackdagger" in the log page in the Web UI.

