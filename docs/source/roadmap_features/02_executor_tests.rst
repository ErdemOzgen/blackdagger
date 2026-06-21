Feature 02: Executor Test Coverage
=================================

Executor coverage was expanded with targeted regression tests.

Covered Areas
-------------

- Webhook config decoding.
- Environment expansion behavior.
- Success and failure behavior.
- Custom success status code handling.
- Required URL validation.
- Terraform argument building.
- Terraform required subcommand validation.
- Ansible argument building.
- Ansible required playbook validation.
- HTTP executor callback-oriented env expansion and URL validation.

Test Files
----------

- internal/dag/executor/webhook_test.go
- internal/dag/executor/terraform_test.go
- internal/dag/executor/ansible_test.go
- internal/dag/executor/http_test.go
