.. _Configuration Options:

Configuration Options
=====================

.. contents::
    :local:

.. _Environment Variables:

Environment Variables
----------------------

The following environment variables can be used to configure the BLACKDAGGER. Default values are provided in the parentheses:

- ``BLACKDAGGER_HOST`` (``127.0.0.1``): The host to bind the server to.
- ``BLACKDAGGER_PORT`` (``8080``): The port to bind the server to.
- ``BLACKDAGGER_DAGS`` (``$BLACKDAGGER_HOME/dags``): The directory containing the DAGs.
- ``BLACKDAGGER_IS_BASICAUTH`` (``0``): Set to 1 to enable basic authentication.
- ``BLACKDAGGER_BASICAUTH_USERNAME`` (``""``): The username to use for basic authentication.
- ``BLACKDAGGER_BASICAUTH_PASSWORD`` (``""``): The password to use for basic authentication.
- ``BLACKDAGGER_LOG_DIR`` (``$BLACKDAGGER_HOME/logs``): The directory where logs will be stored.
- ``BLACKDAGGER_DATA_DIR`` (``$BLACKDAGGER_HOME/data``): The directory where application data will be stored.
- ``BLACKDAGGER_SUSPEND_FLAGS_DIR`` (``$BLACKDAGGER_HOME/suspend``): The directory containing DAG suspend flags.
- ``BLACKDAGGER_ADMIN_LOG_DIR`` (``$BLACKDAGGER_HOME/logs/admin``): The directory where admin logs will be stored.
- ``BLACKDAGGER_BASE_CONFIG`` (``$BLACKDAGGER_HOME/config.yaml``): The path to the base configuration file.
- ``BLACKDAGGER_NAVBAR_COLOR`` (``""``): The color to use for the navigation bar. E.g., ``red`` or ``#ff0000``.
- ``BLACKDAGGER_NAVBAR_TITLE`` (``BLACKDAGGER``): The title to display in the navigation bar. E.g., ``BLACKDAGGER - PROD`` or ``BLACKDAGGER - DEV``
- ``BLACKDAGGER_WORK_DIR``: The working directory for DAGs. If not set, the default value is DAG location. Also you can set the working directory for each DAG steps in the DAG configuration file. For more information, see :ref:`specifying working dir`.
- ``BLACKDAGGER_WORK_DIR``: The working directory for DAGs. If not set, the default value is DAG location. Also you can set the working directory for each DAG steps in the DAG configuration file. For more information, see :ref:`specifying working dir`.
- ``BLACKDAGGER_CERT_FILE``: The path to the SSL certificate file.
- ``BLACKDAGGER_KEY_FILE`` : The path to the SSL key file.

Note: If ``BLACKDAGGER_HOME`` environment variable is not set, the default value is ``$HOME/.blackdagger``.

Config File
--------------

You can create ``admin.yaml`` file in the ``$BLACKDAGGER_HOME`` directory (default: ``$HOME/.blackdagger/``) to override the default configuration values. The following configuration options are available:

.. code-block:: yaml

    host: <hostname for web UI address>                          # default: 127.0.0.1
    port: <port number for web UI address>                       # default: 8080

    # path to the DAGs directory
    dags: <the location of DAG configuration files>              # default: ${BLACKDAGGER_HOME}/dags
    
    # Web UI Color & Title
    navbarColor: <ui header color>                               # header color for web UI (e.g. "#ff0000")
    navbarTitle: <ui title text>                                 # header title for web UI (e.g. "PROD")
    
    # Basic Auth
    isBasicAuth: <true|false>                                    # enables basic auth
    basicAuthUsername: <username for basic auth of web UI>       # basic auth user
    basicAuthPassword: <password for basic auth of web UI>       # basic auth password

    # API Token
    isAuthToken: <true|false>                                    # enables API token
    authToken: <token for API access>                            # API token

    # Base Config
    baseConfig: <base DAG config path>                           # default: ${BLACKDAGGER_HOME}/config.yaml

    # Working Directory
    workDir: <working directory for DAGs>                        # default: DAG location

    # SSL Configuration
    tls:
        certFile: <path to SSL certificate file>
        keyFile: <path to SSL key file>

.. _Host and Port Configuration:

Server's Host and Port Configuration
-------------------------------------

To specify the host and port for running the BLACKDAGGER server, there are a couple of ways to do it.

The first way is to specify the ``BLACKDAGGER_HOST`` and ``BLACKDAGGER_PORT`` environment variables. For example, you could run the following command:

.. code-block:: sh

    BLACKDAGGER_PORT=8000 blackdagger server

The second way is to use the ``--host`` and ``--port`` options when running the ``blackdagger server`` command. For example:

.. code-block:: sh

    blackdagger server --port=8000

See :ref:`Environment Variables` for more information.
