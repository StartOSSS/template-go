# GEMINI Agent Playbook

Gemini often performs exploratory refactors or cross-cutting improvements. Follow these
rules to keep the template healthy.

## Repo fundamentals

- **Packages**: `cmd/server` (entrypoint), `internal/*` (config, db, observability, todo),
  `proto/` (API contracts), `charts/` + `deploy/` + `skaffold.yaml` (deploy), and
  `terraform/**` (infra). Respect this layering; do not import from `cmd/` inside
  `internal/` packages.
- **Secrets**: Terraform builds Cloud SQL + Secret Manager. Passwords are generated via
  `random_password` and stored as `<env>-todo-database-password`; `DATABASE_URL` secrets are
  `<env>-todo-database-url`. Cloud Run services are updated via `--set-secrets` after each
  Skaffold run—never bake secrets into code or YAML.
- **Preview deploys**: `scripts/preview-deploy.sh` slugifies branches, runs Skaffold, then
  mounts the secret via `gcloud run services update`. Keep it POSIX-compatible and
  idempotent.

## Workflow expectations

1. **Gather context** – read existing docs (`README.md`, `docs/best-practices.md`,
   `deploy/cloudrun/README.md`, `terraform/README.md`).
2. **Spec-first** – adjust `docs/spec` + `proto/` before editing Go services.
3. **Decide test plan** – figure out which of `make lint/test/integration-test/load-test/
   security-scan/preview-deploy/plan` you must run.
4. **Implement** – prefer declarative files (YAML/Terraform) and shared helpers (Makefile)
   over bespoke scripts.
5. **Document** – any new command, env var, or workflow step must be captured in the docs.

## Coding & testing standards

- Go code uses `gofmt`, `goimports`, and is organized into small packages. Compose services
  via interfaces so tests can stub dependencies without third-party mock frameworks.
- Prefer table-driven tests with clear Arrange/Act/Assert sections. Keep fixtures in the
  same package or `testdata` directories.
- Integration hooks (`make integration-test`, `make load-test`) rely on Skaffold + Helm
  tests; update them whenever the behavior/contract changes.
- For Terraform, keep resources idempotent and ensure `make plan` stays green for all tfvars
  (dev/preprod/prod). Use modules from terraform-google-modules whenever possible.

## Checklist before pushing

- [ ] Specs + proto regenerated if API surface changed.
- [ ] `make lint` and `make test` run locally (attach logs when reporting issues).
- [ ] Additional targets executed as needed (`integration-test`, `load-test`,
      `security-scan`, `preview-deploy`, `plan`).
- [ ] Documentation updated (README + relevant docs under `docs/` or `deploy/`).
- [ ] Terraform plans inspected for dev/preprod/prod; ensure generated secrets remain in SM
      and Cloud Run mounts them via Skaffold/GitHub Actions.
- [ ] CI workflows adjusted when new automation is introduced (step summary + PR comments).

If any item is blocked (e.g., gcloud unavailable), note it explicitly in the PR summary.
