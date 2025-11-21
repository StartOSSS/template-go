# Copilot Instructions

These prompts guide GitHub Copilot (or any in-editor assistant) when working inside this
repo.

## Default workflow

1. **Read context** – scan `README.md`, `docs/best-practices.md`, and the file you’re
   editing. Suggesting changes without understanding the layered architecture (`cmd/`,
   `internal/`, `proto/`, `terraform`, `charts`, `deploy`) is not acceptable.
2. **Spec-first mindset** – when adding/changing APIs, update `docs/spec/*` and
   `proto/` definitions before implementing handlers.
3. **Reference helpers** – prefer calling existing functions (config loaders, DB helpers,
   observability utilities) and existing Makefile targets instead of writing new ad-hoc
   scripts.

## Coding rules for Copilot

- **Formatting**: all Go code must be `gofmt`/`goimports` clean. Suggest table-driven tests
  with descriptive names. Keep functions small and packages cohesive.
- **Testing**: whenever adding logic, propose `_test.go` updates near the implementation.
  Use the standard library `testing` package and interface-based stubs. Remind the user to
  run `make lint` + `make test` (and `make integration-test` / `make load-test` /
  `make security-scan` / `make plan` when appropriate) before pushing.
- **Docs**: whenever a command, workflow, or variable changes, prompt the user to update the
  relevant doc (README, `docs/best-practices.md`, `deploy/cloudrun/README.md`,
  `terraform/README.md`).
- **Secrets**: never suggest hardcoding credentials. Refer to Terraform’s random password +
  Secret Manager pattern and the Cloud Run `--set-secrets` step.

## Review checklist Copilot should remind about

- Specs/protos updated (if API change).
- Go code layered correctly and covered by tests.
- Terraform `make plan` verified for dev/preprod/prod; secrets stay in SM.
- Skaffold/GitHub workflows adjusted + documented if behavior changes.
- Preview deploy script still slugifies branches, verifies `gcloud`/`skaffold`, mounts
  secrets, and prints the URL.

When suggesting changes, include the relevant commands (`make lint`, `make test`, etc.) in
the response so contributors remember to run them before committing.
