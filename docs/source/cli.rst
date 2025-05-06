Command Line Interface
======================

The following commands are available for interacting with Blackdagger:

.. code-block:: sh

  # Runs the DAG
  blackdagger start [--params=<params>] <file>
  
  # Displays the current status of the DAG. If no DAG file is provided, displays the status of all DAGs.
  blackdagger status [file]
  
  # Re-runs the specified DAG run
  blackdagger retry --req=<request-id> <file>
  
  # Stops the DAG execution. If no DAG file is provided, stops all DAGs.
  blackdagger stop [file]
  
  # Restarts the current running DAG. If no DAG file is provided, restarts all DAGs.
  blackdagger restart [file]
  
  # Dry-runs the DAG
  blackdagger dry [--params=<params>] <file>
  
  # Launches both the web UI server and scheduler process
  blackdagger start-all [--host=<host>] [--port=<port>] [--dags=<path to directory>]
  
  # Launches the blackdagger web UI server
  blackdagger server [--host=<host>] [--port=<port>] [--dags=<path to directory>]
  
  # Starts the scheduler process
  blackdagger scheduler [--dags=<path to directory>]

  # Pulls the latest version of the DAGs from the repository
  blackdagger pull [category] [--origin] [--check] [--keep]
  
  # Shows the current binary version
  blackdagger version