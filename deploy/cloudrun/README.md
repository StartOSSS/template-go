# Cloud Run deployment

CI executes `skaffold run` with the Cloud Run deployer targeting `${GCP_PROJECT}` and
`${GCP_REGION}`. Ensure the following secrets are set in GitHub Actions:

- `GCP_SA` – service account with deploy + Secret Manager permissions
- `GAR_REPOSITORY` – Artifact Registry repository used for pushing `todo-app`,
  `todo-integration`, and `todo-k6`
- `SERVICE_NAME` – Cloud Run service receiving the deployment

Terraform provisions the database password + `DATABASE_URL` inside Secret Manager. Every
`skaffold run` (preview or production) is followed by a `gcloud run services update` that
mounts the `DATABASE_URL` secret as an environment variable so the container never stores
credentials in git or CI.

Metrics/Logs/Traces automatically flow to Cloud Monitoring because OTLP exporters inherit
the Cloud Run metadata server endpoint.

## Workflows

- `Deploy Cloud Run` runs on pushes to `main` (or manually) and refreshes the production
  revision with the latest commit.
- `Preview Deploy` listens for PR comments. Comment `/deploy` on a pull request to run the
  Skaffold `preview` profile (Cloud Build + Cloud Run) via `make preview-deploy`, tag the
  resulting revision with a short branch slug, and surface the URL in both the workflow
  summary and PR comments so features can be reviewed independently.
