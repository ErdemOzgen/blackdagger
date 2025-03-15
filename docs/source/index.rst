Blackdagger
======================================
.. raw:: html

   <div>
      <div class="github-star-button">
      <iframe src="https://ghbtns.com/github-btn.html?user=ErdemOzgen&repo=blackdagger&type=star&count=true&size=large" frameborder="0" scrolling="0" width="160px" height="30px"></iframe>
      </div>
   </div>

Blackdagger is a powerful Cron alternative that comes with a Web UI. It allows you to define dependencies between commands as a `Directed Acyclic Graph (DAG) <https://en.wikipedia.org/wiki/Directed_acyclic_graph>`_ in a declarative :ref:`YAML Format`. Additionally, Blackdagger natively supports running Docker containers, making HTTP requests, and executing commands over SSH. Blackdagger was designed to be easy to use, self-contained, and require no coding, making it ideal for small projects.

Purpose
--------
Legacy systems often have complex and implicit dependencies between jobs. When there are hundreds of cron jobs on a server, it can be difficult to keep track of these dependencies and to determine which job to rerun if one fails. It can also be a hassle to SSH into a server to view logs and manually rerun shell scripts one by one. 

Blackdagger aims to solve these problems by allowing you to explicitly visualize and manage pipeline dependencies as a DAG, and by providing a web UI for checking dependencies, execution status, and logs and for rerunning or stopping jobs with a simple mouse click.

How It Works
------------
Blackdagger simplifies workflow management by operating as a standalone command-line tool, leveraging the local file system for data storage — eliminating the need for database management systems or cloud services. It enables the definition of DAGs in an intuitive, declarative YAML format, ensuring that existing programs can be seamlessly integrated without any modifications.

Why Not Use an Existing Workflow Scheduler Like Airflow?
---------
While there are numerous workflow schedulers like Airflow available, these often necessitate the authoring of DAGs through programming languages such as Python. For legacy systems with extensive job configurations, incorporating code in languages like Perl or Shell Script can already be a complex endeavor. Introducing an additional layer with such tools can further complicate maintainability. In contrast, Blackdagger is crafted for simplicity and usability, requiring no coding skills. This makes it a perfect fit for smaller projects looking for a straightforward, self-sufficient workflow management solution.


.. toctree::
   :maxdepth: 2
   :caption: About Framework
   :hidden:

   About the Framework <framework-intro>
   Components <components/index>

.. toctree::
   :caption: First Steps
   :hidden:

   installation
   quickstart

.. toctree::
   :caption: Interfaces
   :hidden:

   cli
   web_interface

.. toctree::
   :caption: Writing YAMLs
   :hidden:

   yaml_format
   examples
   advanced
   executors
   scheduler
   email

.. toctree::
   :caption: Configurations
   :hidden:

   config
   auth
   api_token

.. toctree::
   :caption: Docker
   :hidden:  

   docker-compose
   docker

.. toctree::
   :caption: REST API
   :hidden:    

   rest

.. toctree::
   :caption: FAQ and Contribution
   :hidden:  

   faq
   contrib