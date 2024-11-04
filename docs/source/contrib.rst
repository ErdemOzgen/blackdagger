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
- ``docs``: Contains the documentation for the project.
- ``examples``: Contains the example workflows.
- ``models``: Contains code related to data models and structures that represent the data used within the application. 
- ``restapi``: Contains code related to implementing a RESTful API
- ``schemas``: Contains code that determines YAML format.
- ``scripts``: Contains scripts for download and installation of Blackdagger itself and extra utilities.
- ``service``: TBU
- ``ui``: Frontend code for the Web UI.
- ``internal``: Contains the internal code for the project.

  - ``agent``: Contains the code for runnning the workflows.
  - ``config``: Contains the code for loading the configuration.
  - ``constants``: Constants used in the project.
  - ``dag``: Contains the code for parsing the workflow definition.
  - ``engine``: Contains the for running basic functions of Blackdagger.
  - ``errors``: Contains the code for error handling in commands of DAG.
  - ``executor``: Contains the code for different types of executors.
  - ``grep``: TBU
  - ``logger``: Contains the code for logging functions.
  - ``mailer``: Contains the code for mail functions.
  - ``pb``: Contains the code for steps in DAG.
  - ``persistence``: Contains the code for storage and databases.
  - ``proto``: Contains the code for "Protocol Buffers".
  - ``reporter``: Contains the code for reporting the status of the scheduler.
  - ``scheduler``: Contains the code for scheduler which schedules and runs the processes in a DAG.
  - ``sock``: Contains the code for interacting with the socket.
  - ``utils``: Contains the code for utility functions.

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
