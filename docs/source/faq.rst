FAQ
===

How Long Will the History Data be Stored?
------------------------------------------

By default, the execution history data is retained for 30 days. However, you can customize this setting by modifying the `histRetentionDays` field in a YAML file.

How to Use Specific Host and Port or `blackdagger server`?
-----------------------------------------------------

To configure the host and port for `blackdagger server`, you can set the environment variables `BLACKDAGGER_HOST` and `BLACKDAGGER_PORT`. Refer to the :ref:`Configuration Options` for more details.

How to Specify the DAGs Directory for `blackdagger server` and `blackdagger scheduler`?
--------------------------------------------------------------------------

You can customize the directory used to store DAG files by setting the environment variable `BLACKDAGGER_DAGS`. See :ref:`Configuration Options` for more information.

How Can I Retry a DAG from a Specific Task?
--------------------------------------------

If you want to retry a DAG from a specific task, you can set the status of that task to `failed` by clicking the step in the Web UI. When you rerun the DAG, it will execute the failed task and any subsequent tasks.

Can I Use Symlinked DAG Files and Still Edit SPEC in the Web UI?
-----------------------------------------------------------------

Yes. Blackdagger supports DAG files that are symlinked and can also work when
the configured DAG directory is a symlink.

If SPEC cannot be shown or saved:

- Ensure the process user has read and write permissions for the symlink
	target files.
- Ensure files use ``.yaml`` or ``.yml`` extensions.
- As a fallback, set ``BLACKDAGGER_DAGS`` to the real target directory path
	instead of the symlink path.

Can I Run the Same DAG in Parallel with Different Parameters?
-------------------------------------------------------------

Not natively in the current runtime. A DAG is guarded by a single local socket,
so only one active run of the same DAG is allowed at a time.

Workarounds:

- Use separate DAG files for each parameter set.
- Use one orchestrator DAG that fans out work in parallel inside a step
	(for example with ``xargs -P`` or GNU parallel).

How Does It Track Running Processes Without DBMS?
-------------------------------------------------

`blackdagger` uses Unix sockets to communicate with running processes.
