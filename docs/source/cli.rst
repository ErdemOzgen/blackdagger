Command Line Interface
======================

The following commands are available for interacting with blackdagger:

.. code-block:: sh

  # Runs the DAG
  blackdagger start [--params=<params>] <file>
  
  # Displays the current status of the DAG
  blackdagger status <file>
  
  # Re-runs the specified DAG run
  blackdagger retry --req=<request-id> <file>
  
  # Stops the DAG execution
  blackdagger stop <file>
  
  # Restarts the current running DAG
  blackdagger restart <file>
  
  # Dry-runs the DAG
  blackdagger dry [--params=<params>] <file>
  
  # Launches both the web UI server and scheduler process
  blackdagger start-all [--host=<host>] [--port=<port>] [--dags=<path to directory>]
  
  # Launches the blackdagger web UI server
  blackdagger server [--host=<host>] [--port=<port>] [--dags=<path to directory>]
  
  # Starts the scheduler process
  blackdagger scheduler [--dags=<path to directory>]
  
  # Shows the current binary version
  blackdagger version