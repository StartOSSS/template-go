# template-go best practices

This document captures the conventions that make `template-go` the reference
layout for Cloud Run-focused Go services. Each topic links back to the folder or
workflow that encodes the practice so contributors can trace the "why" as well
as the "how".

## Repository + architecture conventions

- **Layered packages:** Runtime code lives under `cmd/` (entrypoints) and
  `internal/` (domain, adapters, observability). This keeps binaries tiny and
  discourages circular imports while matching Go's standard project layout.
- **Spec-first design:** Every change starts in `docs/spec/*.md` and
  `proto/`. `docs/spec/todo-api.md` spells out RPC contracts and
  non-functional requirements, and Buf (`buf.gen.yaml`, `buf.work.yaml`)
  keeps protobuf generation reproducible.
- **Managed infrastructure code:** Helm charts (`charts/`), Skaffold
  profiles (`skaffold.yaml`), Terraform components (`terraform/components`),
  and Cloud Run docs (`deploy/cloudrun/README.md`) all live next to the app to
  enforce GitOps and enable reproducible environments.

## Development workflow

- **Single Makefile interface:** `Makefile` exports idempotent targets for lint,
  testing, security scans, Terraform planning, and preview deploys so every
  developer + CI job uses the same commands.
- **Tool bootstrap:** `make bootstrap` installs buf, skaffold, minikube, and
  helm to freeze toolchain versions.
- **Pre-commit everywhere:** Install hooks via `pre-commit install`. They run
  `make fmt`, `make lint`, `make test`, and `make plan` on every commit so
  regressions never reach CI.
- **Preview automation:** `scripts/preview-deploy.sh` + `make preview-deploy`
  slugify branches, verify gcloud/skaffold, reuse service-account credentials,
  and surface preview URLs without bespoke scripting.

## Testing + quality gates

- **Multi-stage matrix:** `make lint`, `make test`, `make integration-test`,
  `make load-test`, `make security-scan`, and `make plan` cover formatting,
  unit, integration, load, security, and infrastructure drift with explicit
  commands documented in `docs/contributor/guide.md`.
- **On-cluster validation:** Helm tests (`charts/todo-app/tests`, invoked via
  Skaffold) exercise integration and k6 load paths in the same environment used
  by developers.
- **Container verification:** `structure-test.yaml` plus container-structure-test
  in `make test` ensures runtime images keep required files + permissions.

## Deployment + operations

- **Skaffold profiles for every stage:** `skaffold.yaml` defines `minikube`,
  `cloudrun`, and `preview` profiles. All Cloud Run deploys build via Cloud
  Build (`googleCloudBuild.projectId=${GCP_PROJECT}`) for provenance, and both
  Cloud Run profiles accept `${GCP_PROJECT}`, `${GCP_REGION}`, `${SERVICE_NAME}`
  inputs.
- **Preview environments:** `.github/workflows/preview-deploy.yml` reacts to
  `/deploy` comments, reuses the same script used locally, tags revisions via
  Cloud Run traffic splitting, and posts URLs to PRs.
- **Production deploys:** `.github/workflows/deploy-cloudrun.yml` performs
  authenticated Skaffold runs from the `cloudrun` profile and configures
  Artifact Registry auth via gcloud to keep secrets centralized in GitHub.
- **Terraform discipline:** `.github/workflows/terraform.yml` runs fmt/plan on
  PRs, while `make plan` iterates every environment (`dev`, `preprod`, `prod`)
  with backend disabled to make preview changes safe.

## Observability + reliability

- **End-to-end telemetry:** The server enables OTLP traces/metrics/logs and
  exposes Prometheus scrape endpoints; `charts/observability-dashboards` ships
  curated Grafana dashboards synchronised via Helm.
- **Local Grafana stack:** `make dev` + Skaffold start Postgres, Grafana,
  Tempo, Prometheus, and dashboards so developers can validate telemetry before
  shipping.
- **Health & load checks:** `/healthz`, `TodoService/HealthCheck`, integration
  Helm tests, and k6 load hooks keep regression detection fast.

## Security, supply-chain, and compliance

- **Defense-in-depth scans:** `make security-scan` runs osv-scanner, Syft,
  Grype, and Gitleaks; `.github/workflows/security.yml` mirrors those checks in
  CI and uploads SARIF when supported.
- **Scorecard + Dependabot:** `.github/workflows/scorecard.yml` and
  `dependabot.yml` keep dependencies and governance in line with OSSF guidance.
- **SLSA & SSDF alignment:** `docs/contributor/guide.md` calls out SP 800-218
  checkpoints, while SLSA attestations come from the release workflows. Skaffold's
  Cloud Build integration plus `scripts/build-artifacts.sh` (used by hooks) ensure
  deterministic builds.
- **Secrets management:** Terraform generates Cloud SQL passwords, stores both the
  password + `DATABASE_URL` inside Secret Manager, and Cloud Run revisions mount the
  `DATABASE_URL` secret via `gcloud run services update` after every Skaffold deploy. CI
  relies on GitHub OIDC + the `GCP_SA` secret, keeping credentials out of the repo.

## Documentation + governance

- **Contributor expectations:** `docs/contributor/guide.md` mandates spec-first
  changes, conventional commits, release-label discipline, and "docs + tests with
  every PR".
- **Consumer onboarding:** `docs/consumer/quickstart.md` shows how clients
  integrate with Connect, describe auth headers, and route traffic to previews or
  prod.
- **AIP compliance log:** `docs/spec/aip-compliance.md` tracks adherence to
  Google AIPs (1, 2, 3, 8, 9, 100, 111, 200, 205) and records temporary
  exceptions for auditing.

Use this document when cloning the template into new services so the rationale
behind each directory, workflow, and dependency remains clear.
