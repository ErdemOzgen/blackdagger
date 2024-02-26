Installation
============

.. contents::
    :local:

You can install blackdagger quickly using Homebrew or by downloading the latest binary from the Releases page on GitHub.



Via Bash script
---------------

.. code-block:: bash

   curl -L https://raw.githubusercontent.com/ErdemOzgen/blackdagger/main/scripts/downloader.sh | bash

Via Docker
----------

.. code-block:: bash

   docker run \
   --rm \
   -p 8080:8080 \
   -v $HOME/.blackdagger/dags:/home/blackdagger/.blackdagger/dags \
   -v $HOME/.blackdagger/data:/home/blackdagger/.blackdagger/data \
   -v $HOME/.blackdagger/logs:/home/blackdagger/.blackdagger/logs \
   ErdemOzgen/blackdagger:latest

Via GitHub Release Page
-----------------------

Download the latest binary from the `Releases page <https://github.com/ErdemOzgen/blackdagger/releases>`_ and place it in your ``$PATH`` (e.g. ``/usr/local/bin``).