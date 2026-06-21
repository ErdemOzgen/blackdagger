Feature 07: Parallel Same-DAG Behavior
======================================

Current runtime behavior does not support multiple active instances of the same DAG concurrently.

Current Limitation
------------------

A DAG is guarded by a single local socket, so only one active run of that DAG is allowed at a time.

Workarounds
-----------

- Use separate DAG files per parameter set.
- Use a parent DAG that fans out concurrent work inside a step.
