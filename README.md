# template-go

**template-go** is an opinionated Golang backend template optimised for Cloud Run + GitOps
pipelines. It ships with Connect-RPC APIs, Postgres persistence, OTel instrumentation, and
a full developer experience that can be reproduced locally (via Skaffold + Helm +
Minikube) or in CI.

## Highlights

- ✅ Spec-first development via protobuf definitions in `proto/` + docs in `docs/spec/`
- ✅ Connect RPC server implemented in `cmd/server`, backed by Postgres + embedded migrations
- ✅ Health checks (`/healthz` + `TodoService/HealthCheck`) and automatic local seeding when
  running outside GCP
- ✅ Observability: OTLP traces/metrics/logs, Prometheus endpoint, Grafana stack (Grafana,
  Loki, Tempo, Mimir, Pyroscope) bootstrapped through upstream Helm charts
- ✅ GitOps-ready: Skaffold config for Cloud Run deploys, Terraform modules for IAM + secrets,
  GitHub workflows for linting, testing, integration, load, and release automation
- ✅ On-demand preview envs via `/deploy` comments in CI or `make preview-deploy` locally
- ✅ Documented best practices in `docs/best-practices.md` explaining every folder + workflow
- ✅ Database credentials are auto-generated in Terraform, stored in Secret Manager, and mounted into Cloud Run revisions
- ✅ Security & supply-chain posture: Dependabot, OSSF Scorecard, osv-scanner, Syft, Grype,
  Gitleaks, container-structure-tests, and NIST SSDF controls baked into CI

## Getting started

```bash
# 1. bootstrap toolchain & minikube
make bootstrap

# 2. run the app + Postgres + grafana stack locally
make dev

> `make dev` now spins up Postgres plus the Grafana/Tempo/Prometheus trio so you
> can validate OTLP traces + metrics locally. Each component runs with
> single-replica, in-memory storage to keep requirements friendly for
> single-node Minikube clusters.

# 3. run unit tests
make test

# 4. launch integration & k6 volume tests via Helm hooks
make integration-test
make load-test

# 5. confirm Terraform still plans cleanly for every env
make plan
```

The default configuration lives in `.env.example`. Copy it into `.env` and customise secrets
before running the stack. Skaffold pushes application + test images to Artifact Registry
(`GAR_REPOSITORY`) for Cloud Run deploys and loads them straight into Minikube for local
profiles.

## Best practices reference

See `docs/best-practices.md` for the rationale behind each directory, workflow,
and automation choice baked into this template. Reference it when cloning the
repo to keep conventions consistent across services.

## Preview environments

- Comment `/deploy` on any pull request to trigger the **Preview Deploy** workflow. GitHub
  Actions authenticates with Google Cloud, runs the Skaffold `preview` profile via Cloud Build
  + Cloud Run, and tags the resulting revision with a short slug derived from the branch name.
  The workflow captures the Cloud Run URL for that tag, publishes it in the job summary, and
  posts it back on the pull request thread.
- Run `gcloud auth login --brief` (or export `GOOGLE_APPLICATION_CREDENTIALS`) locally and
  call `make preview-deploy` to produce the same tag/URL on demand. The helper script accepts
  overrides for `BRANCH`, `COMMIT_SHA`, `PREVIEW_TAG`, `GCP_PROJECT`, `GCP_REGION`,
  `SERVICE_NAME`, and `SKAFFOLD_DEFAULT_REPO` so forks or sandboxes can point to alternate
  projects/regions.
- `scripts/preview-deploy.sh` bootstraps authentication automatically: it reuses
  `GOOGLE_APPLICATION_CREDENTIALS` in CI and falls back to `gcloud auth login --brief` if no
  active account exists locally. CI workflows must still call
  `google-github-actions/auth@v2` (or similar) before invoking `make preview-deploy` so the
  gcloud SDK has an active account.

## Repository map

| Path | Purpose |
|------|---------|
| `cmd/server` | Connect RPC HTTP server + health endpoint |
| `internal/` | Config, DB, observability, and todo service layers |
| `charts/todo-app` | Reusable Helm chart powering app + integration/load Helm tests |
| `charts/observability-dashboards` | Dashboards surfaced inside Grafana |
| `deploy/helm-values` | Overrides for upstream Helm charts (Postgres, Grafana stack) |
| `integration/` | Black-box Connect client container executed as a Helm test |
| `loadtest/` | Grafana k6 plan + container definition |
| `terraform/modules` | Baseline Cloud Run infrastructure module |
| `terraform/components/*` | Component definitions + per-environment tfvars |
| `docs/` | Specs, consumer onboarding, contributor handbook |

## Google AIP alignment

- `docs/spec/aip-compliance.md` tracks how the template implements the generally applicable
  Google AIPs (1, 2, 3, 8, 9, 100, 111, 200, 205) and provides a checklist for new APIs.
- Every spec under `docs/spec/` now carries a `Version:` header with an ISO-8601 tag per
  [AIP-3](https://google.aip.dev/3). Update the date when the protobuf contract or behavior
  changes.
- Pull requests proposing beta/GA-ready changes must copy the latest AIP checklist into the
  description and link to design reviews in accordance with [AIP-100](https://google.aip.dev/100).
- Temporary exceptions require a `TODO(aip-200)` note plus an issue link so they do not set
  precedent (AIP-200) and beta blockers are tracked in the spec per [AIP-205](https://google.aip.dev/205).

## Observability

- Prometheus scrape endpoint on `${METRICS_ADDR}` (default `:9090`)
- Grafana dashboards for API + Postgres metrics seeded via
  `charts/observability-dashboards/dashboards`
- Tempo + Loki + Pyroscope stacks available locally; Cloud Run deploys export to Cloud
  Monitoring / Cloud Trace / Cloud Profiler instead (see `deploy/cloudrun/README.md`).

## GitHub workflows

| Workflow | Description |
|----------|-------------|
| `Lint` | Buf, Go, Terraform, Shell, Docker, and Bash linting + scorecard checks |
| `Test` | Go unit tests + container-structure-tests and chained E2E Helm runs |
| `E2E` | Minikube + Skaffold deploy, Helm integration + load hooks, Grafana smoke test |
| `Security` | osv-scanner, Syft SBOM, Grype image scan, Gitleaks |
| `Deploy Cloud Run` | Builds images with Skaffold and deploys to Cloud Run using the `GCP_SA` secret |
| `Preview Deploy` | Comment `/deploy` on a PR to run the Skaffold preview profile (Cloud Build + Cloud Run) and return a tagged URL |
| `Release` | Semantic tagging automation triggered after PR merge based on labels |
| `Scorecard` | OSSF scorecard nightly run (results uploaded to the Security tab and exposed
  via the badge below) |
| `Terraform` | Validates fmt + plans all three environments whenever `terraform/**` changes |

Badges:

- ![lint](https://github.com/example/template-go/actions/workflows/lint.yml/badge.svg)
- ![test](https://github.com/example/template-go/actions/workflows/test.yml/badge.svg)
- ![e2e](https://github.com/example/template-go/actions/workflows/e2e.yml/badge.svg)
- ![scorecard](https://api.securityscorecards.dev/projects/github.com/example/template-go/badge)

> ℹ️ The Scorecard workflow uploads SARIF + JSON results back into GitHub's Security tab and
> powers the badge above via `api.securityscorecards.dev`, so the status is always sourced
> from the latest nightly run without needing to commit generated artifacts.

## Versioning & release process

`release.yml` promotes every merge into `main` by tagging the commit. PR labels
`release:major|minor|patch` control the bump (patch by default). Container images are
published to Artifact Registry (`GAR_REPOSITORY`) and Terraform wires the service accounts /
permissions required by Cloud Run.

## Security posture

- SLSA-aligned build provenance surfaced via GitHub workflows + attestations
- NIST SSDF controls mapped in `docs/contributor/guide.md`
- Automated scanning: osv-scanner, Syft SBOM, Grype vulnerability detection, Gitleaks secret scanning
- Scorecard + Dependabot keep dependencies up-to-date

## Contributing

See `docs/contributor/guide.md` for branching strategy, commit conventions, testing matrix,
and how to run `helm test` hooks for integration/load scenarios. PRs must be spec-aligned:
update protobuf definitions + docs first, regenerate artifacts, then implement services.
