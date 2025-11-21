# Repository Agent Guide

Use this document (plus `CLAUDE.md`, `GEMINI.md`, and `copilot-instructions.md`) as the
authoritative reference when automating changes to this repo.

## Mission

Deliver production-grade Golang gRPC/Connect services for Cloud Run. Every change must keep
the repo usable as a template: pristine docs, green CI, reproducible infra, and minimal
manual steps.

## Repository layout

- `cmd/server`: main entrypoint. Keep it thin—wire configs, DI, and HTTP servers here.
- `internal/<layer>`: business logic (config, db, observability, todo service). Add new
  subpackages instead of bloating existing ones. Export interfaces for collaborators.
- `proto/`, `docs/spec/`: API-first workflow. Update specs + protobufs before code.
- `charts/`, `deploy/`, `skaffold.yaml`: deployment stack (Minikube + Cloud Run).
- `terraform/**`: infrastructure modules + environment configs.
- `scripts/`: repeatable automation (bootstrap, preview deploy). Prefer shell + Makefile
  targets over ad-hoc instructions.

## Implementation workflow

1. **Plan** – read open docs + `docs/best-practices.md`. Identify which envs/files change.
2. **Update specs first** – touch `proto/` + `docs/spec/` before internal Go code.
3. **Code** – keep packages layered (`cmd` -> `internal` -> leaf packages). Use Go
   interfaces for external dependencies (DB, messaging) to simplify tests/mocks.
4. **Secrets** – never hardcode credentials. Use Terraform + Secret Manager (see
   `terraform/modules/app`) and mount them in Cloud Run via Skaffold workflows.
5. **Docs** – update README, `docs/`, or workflow docs for any user-facing change.
6. **Automation** – prefer declarative config (YAML/TOML/terraform) over bespoke scripts.

## Coding standards

- Run `gofmt`/`goimports` (or `make fmt`) on all Go code. Keep files under `cmd/` focused on
  wiring; business logic lives in `internal/` packages with clear interfaces.
- Use table-driven tests. Place `_test.go` files next to the code under test. Stub
  dependencies via interfaces (see `internal/db.Pool`). If a mock is missing, add an
  interface + hand-written stub instead of introducing a heavy mocking framework.
- Keep Makefile targets idempotent and composable. If you add commands, document them in
  `README.md` and reference them from relevant workflows.

## Testing & QA

Run relevant targets before pushing:

- `make lint` – Buf, gofmt, golangci-lint, terraform fmt, shellcheck, hadolint.
- `make test` – Go unit tests + container structure tests.
- `make integration-test`, `make load-test` – Helm hooks when touching runtime paths.
- `make security-scan` – osv-scanner, syft, grype, gitleaks.
- `make plan` – Terraform init/plan for dev/preprod/prod (backend disabled).

CI mirrors those checks plus preview deploys, Cloud Run releases, and Terraform planning.
Never merge with red checks. When adding workflows, ensure they align with SLSA/SSDF
guidance already referenced in `docs/`.

## Infra & deployment notes

- Terraform auto-generates database credentials and stores them in Secret Manager
  (`<env>-todo-database-url` + `<env>-todo-database-password`). Cloud Run services are
  updated via `gcloud run services update --set-secrets DATABASE_URL=...` after Skaffold
  deploys. Never commit plaintext secrets.
- Preview envs come from `make preview-deploy` / `/deploy`. Ensure `scripts/preview-deploy.sh`
  stays cross-platform and reuses `gcloud` auth flows.
- Document any new workflow/command in README + `deploy/cloudrun/README.md` (if relevant).

## Review checklist for agents

- [ ] Specs + proto updated if API behavior changed.
- [ ] Go code formatted, tests added/updated, and targeted `make` commands executed.
- [ ] Docs/README/workflow docs updated for user-facing or automation changes.
- [ ] Terraform or Skaffold changes validated with `make plan` / `skaffold run` as
      appropriate. Secrets handled via SM (no plaintext in git).
- [ ] CI implications considered: new workflows include permissions, summary output, and
      PR comments where needed (e.g., Terraform plan).

See `CLAUDE.md`, `GEMINI.md`, and `copilot-instructions.md` for agent-specific callouts.
