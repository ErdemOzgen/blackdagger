Feature 06: Symlink and Extension Handling for SPEC UI/API
===========================================================

DAG store behavior was improved for symlinked roots/files and extension resolution.

Delivered Improvements
----------------------

- Symlink-aware DAG root resolution in DAG store.
- Recursive DAG discovery uses resolved roots consistently.
- SPEC lookup preserves real extension handling for both .yaml and .yml.
- Dedicated regression tests for symlinked roots/files and .yml access.

Files
-----

- internal/persistence/local/dag_store.go
- internal/persistence/local/dag_store_test.go
