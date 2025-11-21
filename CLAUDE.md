# CLAUDE Agent Playbook

Claude is expected to take on larger, multi-step changes. Use this playbook whenever you
modify the template.

## 1. Understand the repo

- **Architecture**: Connect RPC server in `cmd/server`, business logic under `internal/`,
  Terraform + Skaffold for Cloud Run, Helm for Minikube, Grafana-based observability stack.
- **Spec-first**: Update `docs/spec/*` + `proto/` before changing Go code. Regenerate Buf
  outputs via Makefile targets if necessary.
- **Secret flow**: Terraform generates DB passwords (`random_password`) and writes both the
  password + `DATABASE_URL` secrets to Secret Manager. Cloud Run revisions mount
  `DATABASE_URL` via `gcloud run services update --set-secrets ...` after Skaffold deploys.

## 2. Plan the change

- List every command and workflow you will touch. Prefer declarative config (YAML/TOML,
  Terraform, Skaffold) over ad-hoc scripts.
- Identify the environments impacted. `make plan` and the Terraform workflow plan **dev**,
  **preprod**, and **prod** via `terraform/components/app/<env>.tfvars`—use the same files.
- Decide whether preview deploy automation (`scripts/preview-deploy.sh`, `.github/workflows`
  preview job) needs updates.

## 3. Coding standards

- Go code must be `gofmt`/`goimports` clean. Keep `cmd/` thin and push logic into
  `internal/<package>`. Use interfaces to enable hand-written mocks (see `internal/db.Pool`).
- Configuration lives in `internal/config`. Never pull env vars from random places.
- Use table-driven tests. Put `_test.go` files next to implementation files. For external
  deps, create lightweight interfaces and fake implementations.
- Follow the existing logging/observability conventions (OTel spans via middleware,
  Prometheus metrics in `internal/observability`).

## 4. Testing & verification

Before pushing or opening a PR, run the relevant `make` targets locally:

- `make lint`
- `make test`
- `make integration-test` / `make load-test` if runtime behavior changes
- `make security-scan`
- `make plan` (verifies Terraform for dev/preprod/prod)

Capture any preview changes by running `make preview-deploy` (needs `gcloud auth login` or
`GOOGLE_APPLICATION_CREDENTIALS`). Confirm the script updates the Cloud Run secret mount.

## 5. Docs & workflows

- Update `README.md`, `docs/best-practices.md`, `deploy/cloudrun/README.md`, or
  `terraform/README.md` whenever commands, secrets, or workflows change.
- Terraform/Skaffold changes must preserve the automatic Secret Manager + Cloud Run mount
  flow. Mention new secrets or env vars explicitly.
- When touching GitHub Actions, ensure they emit helpful step summaries and PR comments
  (e.g., Terraform plan output).

## 6. Review checklist for Claude

- [ ] Specs + proto updated if the API changed.
- [ ] Go code layered correctly, formatted, and backed by table-driven tests.
- [ ] Docs updated for any user-facing or automation change.
- [ ] Terraform/Skaffold secrets handled via SM (`--set-secrets DATABASE_URL=...`).
- [ ] Relevant `make` targets executed locally; CI failures explained or resolved.
- [ ] Preview + production workflows remain idempotent (Makefile targets safe to re-run).

Hand back changes only after the checklist is satisfied and CI is expected to pass.
