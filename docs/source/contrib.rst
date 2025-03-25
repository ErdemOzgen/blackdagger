Contribution Guide
===================

We welcome contributions of any size and skill level. If you have an idea for a new feature or have found a bug, please open an issue on the GitHub repository.

Prerequisite
-------------

* `Go version 1.19 or later. <https://go.dev/doc/install>`_
* Latest version of `Node.js <https://nodejs.org/en/download/>`_.
* `yarn <https://yarnpkg.com/>`_ package manager.

Setting up your local environment
----------------------------------

#. Clone the repository to your local machine.
#. Navigate to the root directory of the cloned repository and build the frontend project by running the following command:

   .. code-block:: sh

      make build-ui

#. Run the following command to start the `Blackdagger` application:

   .. code-block:: sh

      go run main.go

#. Now you can change the source code and build the binary by running the following command:

   .. code-block:: sh

      make build

#. Run the following command to start the `Blackdagger` application:

   .. code-block:: sh

      ./bin/blackdagger

Running Tests
-------------

   Run the following command to run the tests:

   .. code-block:: sh

      make test

Code Structure
---------------

- ``cmd``: Contains the main application entry point.
- ``config``: Contains simple configurations for network purposes.
- ``docs``: Contains the documentation for the project.
- ``examples``: Contains the example workflows.
- ``schemas``: Contains code that determines YAML format.
- ``scripts``: Contains scripts for download and installation of Blackdagger itself and extra utilities.
- ``ui``: Contains frontend code for the Web UI.
- ``internal``: Contains the internal code for the project.
  - ``agent``: Contains the code for runnning the workflows.
  - ``client``: Contains the actions that can be taken on the client side.
  - ``config``: Contains the code for loading the configuration.
  - ``constants``: Constants used in the project.
  - ``dag``: Contains the code for parsing the workflow definition.
  - ``frontend``: Contains the frontend setting and the connections between frontend and backend.
  - ``logger``: Contains the code for logging functions.
  - ``mailer``: Contains the code for mail functions.
  - ``persistence``: Contains the code for storage and databases.
  - ``scheduler``: Contains the code for scheduler which schedules and runs the processes in a DAG.
  - ``sock``: Contains the code for interacting with the socket.
  - ``test``: Contains the main tests for ``internal``.
  - ``util``: Contains the code for utility functions.

Setting up your local environment for front end development
-------------------------------------------------------------

#. Clone the repository to your local machine.
#. Navigate to the root directory of the cloned repository and build the frontend project by running the following command:

   .. code-block:: sh

      make build-ui

#. Run the following command to start the `blackdagger` application:

   .. code-block:: sh

      go run main.go server

#. Navigate to ``ui`` directory and run the following command to install the dependencies:

   .. code-block:: sh

      yarn install
      yarn start

#. Open the browser and navigate to http://localhost:8081.

#. Make changes to the source code and refresh the browser to see the changes.

Branches
---------

* ``main``: The main branch where the source code always reflects a production-ready state.
