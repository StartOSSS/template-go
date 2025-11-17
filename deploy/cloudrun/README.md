# Cloud Run deployment

CI executes `skaffold run` with the Cloud Run deployer targeting `${GCP_PROJECT}` and
`${GCP_REGION}`. Ensure the following secrets are set in GitHub Actions:

- `GCP_SA` – service account with deploy + Secret Manager permissions
- `DATABASE_URL` – connection string stored in GSM and mounted via Terraform outputs
- `GAR_REPOSITORY` – Artifact Registry repository used for pushing `todo-app`,
  `todo-integration`, and `todo-k6`

Metrics/Logs/Traces automatically flow to Cloud Monitoring because OTLP exporters inherit
the Cloud Run metadata server endpoint.
