Feature 05: HTTP Executor Hardening
===================================

HTTP callback reliability was improved for command and handlerOn use cases.

What Changed
------------

- Environment expansion for command-derived method values.
- Environment expansion for URL argument values.
- Environment expansion for query values.
- Explicit URL-required validation with clear error behavior.

Full HTTP Workflow Example
--------------------------

.. code-block:: yaml

  steps:
    - name: get fake json data
      executor:
        type: http
        config:
          timeout: 10
          headers:
          silent: true
          query:
            postId: "1"
          body: ""
      command: GET https://jsonplaceholder.typicode.com/comments
