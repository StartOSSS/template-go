# Contributor guide

## Flow

1. Update specs (`docs/spec`, `proto/`) before writing code.
2. Implement application changes inside `internal/` or `cmd/`.
3. Update Terraform/Skaffold/Helm manifests when touching infrastructure.
4. Run the full test matrix locally via `make all` before opening a PR.

## Commit & PR policy

- Conventional commits (e.g. `feat: add todo filter`)
- PR labels drive release automation: `release:patch` (default), `release:minor`, `release:major`
- Every PR must update documentation and tests alongside code changes

## Testing matrix

| Command | Purpose |
|---------|---------|
| `make lint` | Buf, gofmt, golangci-lint, terraform fmt/check, shellcheck, scorecard dry-run |
| `make test` | Go unit tests + container-structure-tests |
| `make integration-test` | Deploy via Skaffold/Helm and run the integration Helm test hook |
| `make load-test` | Execute Grafana k6 Helm test hook |
| `make security-scan` | osv-scanner + syft + grype + gitleaks |
| `make plan` | Terraform init/plan (backend disabled) for dev/preprod/prod |

## SSDF alignment

The checklist below maps to [NIST SP 800-218](https://nvlpubs.nist.gov/nistpubs/SpecialPublications/NIST.SP.800-218.pdf):

- PO.1.1 / PO.1.2 – Spec-first docs + protobufs ensure requirements traceability
- PO.5.1 – Makefile and Skaffold encode reproducible builds (12-factor aligned)
- PW.4.1 – Security static analysis (linters + scanners) in CI
- PW.8.2 – Tests (unit/integration/load) automated per PR via chained workflows
- RV.1.2 – Dependabot + osv-scanner guard dependencies
- RV.3.1 – Scorecard, SLSA provenance, and signed container attestations

Document any deviations or compensating controls inside this file when enhancing the
template.
