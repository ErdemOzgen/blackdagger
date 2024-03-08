Using Docker Compose
===================================

To automate workflows based on cron expressions, it is necessary to run both the ui server and scheduler process. Here is an example `docker-compose.yml` setup for running blackdagger using Docker Compose.

.. code-block:: yaml

    version: "3.9"
    services:

      # init container updates permission
      init:
        image: "ErdemOzgen/blackdagger:latest"
        user: root
        volumes:
          - blackdagger:/home/blackdagger/.blackdagger
        command: chown -R blackdagger /home/blackdagger/.blackdagger/

      # ui server process
      server:
        image: "ErdemOzgen/blackdagger:latest"
        environment:
          - blackdagger_PORT=8080
          - blackdagger_DAGS=/home/blackdagger/.blackdagger/dags
        restart: unless-stopped
        ports:
          - "8080:8080"
        volumes:
          - blackdagger:/home/blackdagger/.blackdagger
          - ./dags/:/home/blackdagger/.blackdagger/dags
        depends_on:
          - init

      # scheduler process
      scheduler:
        image: "erdemozgen/blackdagger:latest"
        environment:
          - blackdagger_DAGS=/home/blackdagger/.blackdagger/dags
        restart: unless-stopped
        volumes:
          - blackdagger:/home/blackdagger/.blackdagger
          - ./dags/:/home/blackdagger/.blackdagger/dags
        command: blackdagger scheduler
        depends_on:
          - init

    volumes:
      blackdagger: {}
