.. _Remote Instance Configuration:

Remote Instance Management
===========

.. contents::
    :local:

Blackdagger UI can be configured to connect to remote instances, allowing management of DAGs across different environments from a single interface.

How to configure
----------------
Create ``config.yaml`` in ``$HOME/.config/blackdagger/`` to configure remote instances. Example configuration:

.. code-block:: yaml

    # Remote Instance Configuration
    remoteNodes:
    - name: "dev"                                # name of the remote instance
      apiBaseUrl: "http://localhost:8080/api/v1" # Base API URL of the remote instance it must end with /api/v1

      # Authentication settings for the remote instance
      # Basic authentication
      isBasicAuth: true              # Enable basic auth (optional)
      basicAuthUsername: "username"     # Basic auth username (optional)
      basicAuthPassword: "password"    # Basic auth password (optional)

      # api token authentication
      isAuthToken: true              # Enable API token (optional)
      authToken: "your-secret-token" # API token value (optional)

      # TLS settings
      skipTLSVerify: false           # Skip TLS verification (optional)

Using Remote Instances
-----------------
Once configured, remote instances can be selected from the dropdown menu in the top right corner of the UI. This allows you to:

- Switch between different environments
- View and manage DAGs on remote instances
- Monitor execution status across instances

The UI will maintain all functionality while operating on the selected remote instance.