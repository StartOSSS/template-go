#!/usr/bin/env bash
set -euo pipefail
psql "$DATABASE_URL" <<'SQL'
INSERT INTO todo_tasks (title, description)
VALUES
    ('Wire OTLP exporters', 'Ensure traces/metrics/logs stream to the collector'),
    ('Write Helm README', 'Document how to reuse the todo-app chart'),
    ('Refresh Terraform state', 'Run make plan to confirm infra drift is zero'),
    ('Polish Grafana dashboards', 'Add request latency + DB utilisation panels'),
    ('Scorecard follow-up', 'Review nightly scorecard alerts in the Security tab'),
    ('Document Helm tests', 'Walk through helm test usage for integration + load'),
    ('Review OTEL config', 'Double-check sampling + resource attributes'),
    ('Tighten IAM bindings', 'Ensure Cloud Run SA has minimum permissions'),
    ('Validate migrations', 'Re-run migrate and confirm it is idempotent'),
    ('Smoke test healthz', 'Hit /healthz and TodoService/HealthCheck'),
    ('Seed demo data', 'Populate todos for quick local validation'),
    ('Exercise CRUD', 'Call Create/List/Update/Delete via Connect client'),
    ('Tune k6 thresholds', 'Align volume test assertions with SLOs'),
    ('Capture DB metrics', 'Ensure pg_stat_statements dashboard renders'),
    ('Verify release tags', 'Ensure label -> semantic version automation works'),
    ('Review Dependabot PRs', 'Update Go modules + Skaffold when needed'),
    ('Refresh secrets', 'Rotate GSM secrets + update workloads'),
    ('Update contributor guide', 'Add new testing steps to docs'),
    ('Expand consumer docs', 'Add Quickstart snippets for Connect clients'),
    ('Demo Cloud Run deploy', 'Ship to staging and confirm telemetry');
SQL
