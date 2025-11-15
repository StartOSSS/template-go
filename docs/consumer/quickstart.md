# Consumer onboarding

1. Generate a Connect client using the proto files under `proto/` (use Buf or protoc).
2. Point the client at the Cloud Run URL or local dev endpoint (`http://localhost:8080`).
3. For local testing, authenticate via the optional `X-Demo-User` header. In Cloud Run the
   service relies on IAM + identity aware proxy.
4. Observe metrics and traces using Grafana (local) or Cloud Monitoring (prod).
