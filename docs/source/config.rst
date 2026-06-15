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
- ``BLACKDAGGER_DAGS`` (``$HOME/.config/blackdagger/dags``): The directory containing the DAGs.
- ``BLACKDAGGER_IS_BASICAUTH`` (``0``): Set to 1 to enable basic authentication.
- ``BLACKDAGGER_BASICAUTH_USERNAME`` (``""``): The username to use for basic authentication.
- ``BLACKDAGGER_BASICAUTH_PASSWORD`` (``""``): The password to use for basic authentication.
- ``BLACKDAGGER_LOG_DIR`` (``$HOME/.local/share/blackdagger/logs``): The directory where logs will be stored.
- ``BLACKDAGGER_DATA_DIR`` (``$HOME/.local/share/blackdagger/history``): The directory where application data will be stored.
- ``BLACKDAGGER_SUSPEND_FLAGS_DIR`` (``$HOME/.config/blackdagger/suspend``): The directory containing DAG suspend flags.
- ``BLACKDAGGER_ADMIN_LOG_DIR`` (``$HOME/.local/share/admin``): The directory where admin logs will be stored.
- ``BLACKDAGGER_BASE_CONFIG`` (``$HOME/.config/blackdagger/base.yaml``): The path to the base configuration file.
- ``BLACKDAGGER_NAVBAR_COLOR`` (``""``): The color to use for the navigation bar. E.g., ``red`` or ``#ff0000``.
- ``BLACKDAGGER_NAVBAR_TITLE`` (``BLACKDAGGER``): The title to display in the navigation bar. E.g., ``BLACKDAGGER - PROD`` or ``BLACKDAGGER - DEV``
- ``BLACKDAGGER_WORK_DIR``: The working directory for DAGs. If not set, the default value is DAG location. Also you can set the working directory for each DAG steps in the DAG configuration file. For more information, see :ref:`specifying working dir`.
- ``BLACKDAGGER_WORK_DIR``: The working directory for DAGs. If not set, the default value is DAG location. Also you can set the working directory for each DAG steps in the DAG configuration file. For more information, see :ref:`specifying working dir`.
- ``BLACKDAGGER_CERT_FILE``: The path to the SSL certificate file.
- ``BLACKDAGGER_KEY_FILE`` : The path to the SSL key file.
- ``BLACKDAGGER_SKIP_INITIAL_DAG_PULLS``: Set to 1 to skip the automatic pull of default DAG YAML files during startup.
- ``BLACKDAGGER_DAG_REPO``: The prefix for the DAG repository. This is used to pull the categorized DAG YAML files from the repository. The default value is ``https://github.com/ErdemOzgen/blackdagger-``.
- ``BLACKDAGGER_LOG_FORWARDING_ENABLED``: Set to 1 to enable centralized log forwarding.
- ``BLACKDAGGER_LOG_FORWARDING_SINK_TYPE``: Log forwarding sink type. Current MVP supports ``http``.
- ``BLACKDAGGER_LOG_FORWARDING_HTTP_URL``: HTTP endpoint URL used to forward log lines.
- ``BLACKDAGGER_LOG_FORWARDING_TIMEOUT_SEC``: HTTP timeout in seconds for log forwarding requests.
- ``BLACKDAGGER_LOG_FORWARDING_INCLUDE_STEP_OUTPUT``: Set to 1 to include step stdout and stderr files in forwarding.
- ``BLACKDAGGER_LOG_FORWARDING_QUEUE_SIZE``: Async queue size for buffered log forwarding.
- ``BLACKDAGGER_LOG_FORWARDING_MAX_RETRIES``: Maximum retry attempts per forwarded log record.
- ``BLACKDAGGER_LOG_FORWARDING_INITIAL_BACKOFF_MS``: Initial retry backoff in milliseconds.
- ``BLACKDAGGER_LOG_FORWARDING_MAX_BACKOFF_MS``: Maximum retry backoff in milliseconds.
- ``BLACKDAGGER_LOG_FORWARDING_MONITOR_ENABLED``: Set to 1 to expose a local health and metrics endpoint for log forwarding.
- ``BLACKDAGGER_LOG_FORWARDING_MONITOR_HOST``: Host address used by the local monitor endpoint.
- ``BLACKDAGGER_LOG_FORWARDING_MONITOR_PORT``: Port number used by the local monitor endpoint.
- ``BLACKDAGGER_LOG_FORWARDING_MONITOR_BASE_PATH``: Base path for monitor endpoints (for example, ``/log-forwarding``).

Note: If ``BLACKDAGGER_HOME`` environment variable is not set, the default value is ``$HOME/.config/blackdagger``.

Config File
--------------

You can create ``config.yaml`` file in the ``$BLACKDAGGER_HOME`` directory (default: ``$HOME/.config/blackdagger``) to override the default configuration values. The following configuration options are available:

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

    # DAG Configuration
    skipInitialDAGPulls: <true|false>                            # skip initial DAG pulls
    dagRepo: <DAG repository prefix>                             # default: https://github.com/ErdemOzgen/blackdagger-

    # Centralized log forwarding (MVP)
    logForwarding:
      enabled: <true|false>                                    # enables central log forwarding
      sinkType: http                                            # supported: http
      httpURL: https://logs.example.com/ingest                  # destination endpoint
      timeoutSec: 5                                             # request timeout in seconds
      includeStepOutput: <true|false>                           # include stdout/stderr files
      queueSize: 256                                             # async buffer size
      maxRetries: 3                                              # retry attempts per log record
      initialBackoffMS: 100                                      # first retry delay (ms)
      maxBackoffMS: 2000                                         # max retry delay (ms)
      monitorEnabled: <true|false>                               # expose local monitor endpoints
      monitorHost: 127.0.0.1                                     # monitor bind host
      monitorPort: 8091                                          # monitor bind port
      monitorBasePath: /log-forwarding                           # monitor base path
      headers:
        Authorization: "Bearer <token>"                      # optional HTTP headers

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

.. _base configuration:

Base Configuration for all DAGs
---------------------------------

Creating a base configuration (default path: ``~/.config/blackdagger/base.yaml``) is a convenient way to organize shared settings among all DAGs. The path to the base configuration file can be configured. See :ref:`Configuration Options` for more details.

Example:

.. code-block:: yaml

    # directory path to save logs from standard output
    logDir: /path/to/stdout-logs/

    # history retention days (default: 30)
    histRetentionDays: 3

    # Email notification settings
    mailOn:
      failure: true
      success: true

    # SMTP server settings
    smtp:
      host: "smtp.foo.bar"
      port: "587"
      username: "<username>"
      password: "<password>"

    # Error mail configuration
    errorMail:
      from: "foo@bar.com"
      to: "foo@bar.com"
      prefix: "[Error]"

    # Info mail configuration
    infoMail:
      from: "foo@bar.com"
      to: "foo@bar.com"
      prefix: "[Info]"

Centralized Log Forwarding (MVP)
--------------------------------

You can forward DAG execution logs to a central HTTP endpoint.

.. code-block:: yaml

    logForwarding:
      enabled: true
      sinkType: http
      httpURL: https://logs.example.com/ingest
      timeoutSec: 5
      includeStepOutput: false
      queueSize: 256
      maxRetries: 3
      initialBackoffMS: 100
      maxBackoffMS: 2000
      monitorEnabled: true
      monitorHost: 127.0.0.1
      monitorPort: 8091
      monitorBasePath: /log-forwarding
      headers:
        Authorization: "Bearer <token>"

When monitor is enabled, the command process exposes these endpoints:

- ``<monitorBasePath>/health``
- ``<monitorBasePath>/metrics`` (JSON)
- ``<monitorBasePath>/metrics/prometheus`` (Prometheus text format)

You can also request Prometheus format from ``<monitorBasePath>/metrics`` using
``?format=prometheus``.

When retries, queue drops, or terminal delivery failures occur, the command logger emits structured events with:

- ``event=log_forwarding_retry``
- ``event=log_forwarding_drop``
- ``event=log_forwarding_failed``

For a runnable configuration example, see ``examples/config_log_forwarding.yaml``.