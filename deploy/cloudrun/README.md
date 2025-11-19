# Cloud Run deployment

CI executes `skaffold run` with the Cloud Run deployer targeting `${GCP_PROJECT}` and
`${GCP_REGION}`. Ensure the following secrets are set in GitHub Actions:

- `GCP_SA` – service account with deploy + Secret Manager permissions
- `DATABASE_URL` – connection string stored in GSM and mounted via Terraform outputs
- `GAR_REPOSITORY` – Artifact Registry repository used for pushing `todo-app`,
  `todo-integration`, and `todo-k6`

Metrics/Logs/Traces automatically flow to Cloud Monitoring because OTLP exporters inherit
the Cloud Run metadata server endpoint.

## Workflows

- `Deploy Cloud Run` runs on pushes to `main` (or manually) and refreshes the production
  revision with the latest commit.
- `Deploy Preview` listens for PR comments. Comment `/deploy` on a pull request to build the
  branch, deploy it as a new revision, and tag that revision with the branch name truncated to
  10 characters. Each tag receives a dedicated Cloud Run URL following the
  `https://<tag>---todo-app-us-central1.a.run.app` format so feature branches can be reviewed
  independently.
