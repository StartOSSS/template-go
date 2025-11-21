#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cd "${ROOT_DIR}"

: "${GCP_PROJECT:?Set GCP_PROJECT}"
: "${GCP_REGION:?Set GCP_REGION}"
: "${SERVICE_NAME:?Set SERVICE_NAME}"
DEPLOY_ENVIRONMENT=${DEPLOY_ENVIRONMENT:-dev}
DATABASE_SECRET_NAME=${DATABASE_SECRET_NAME:-"${DEPLOY_ENVIRONMENT}-todo-database-url"}

command -v gcloud >/dev/null || { echo "gcloud CLI is required" >&2; exit 1; }
command -v skaffold >/dev/null || { echo "skaffold CLI is required" >&2; exit 1; }
command -v jq >/dev/null || { echo "jq is required to parse Cloud Run responses" >&2; exit 1; }

slugify() {
  local value="$1"
  echo "$value" | tr '[:upper:]' '[:lower:]' | tr -cs 'a-z0-9' '-' |
    sed 's/^-\+//;s/-\+$//' | cut -c1-30
}

branch="${BRANCH:-}"
if [[ -z "$branch" ]]; then
  branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || true)
fi

if [[ -z "$branch" || "$branch" == "HEAD" ]]; then
  branch="${COMMIT_SHA:-preview}"
fi

slug=$(slugify "$branch")
if [[ -z "$slug" ]]; then
  slug=preview
fi

export PREVIEW_TAG=${PREVIEW_TAG:-$slug}
export SKAFFOLD_DEFAULT_REPO=${SKAFFOLD_DEFAULT_REPO:-"us-central1-docker.pkg.dev/${GCP_PROJECT}/template-go"}

echo "Using PREVIEW_TAG=$PREVIEW_TAG"
echo "Using SKAFFOLD_DEFAULT_REPO=$SKAFFOLD_DEFAULT_REPO"

ensure_auth() {
  local active
  active=$(gcloud auth list --filter=status:ACTIVE --format='value(account)' || true)
  if [[ -n "$active" ]]; then
    return
  fi

  if [[ -n "${GOOGLE_APPLICATION_CREDENTIALS:-}" && -f "${GOOGLE_APPLICATION_CREDENTIALS}" ]]; then
    echo "Activating service account from GOOGLE_APPLICATION_CREDENTIALS"
    gcloud auth activate-service-account --key-file "${GOOGLE_APPLICATION_CREDENTIALS}"
  else
    echo "No active gcloud account found; launching login flow"
    gcloud auth login --brief
  fi
}

ensure_auth
gcloud config set project "${GCP_PROJECT}" >/dev/null

echo "Running Skaffold preview deployment"
skaffold run -p preview --default-repo="${SKAFFOLD_DEFAULT_REPO}"

echo "Configuring Cloud Run secret mount for DATABASE_URL via ${DATABASE_SECRET_NAME}"
gcloud run services update "${SERVICE_NAME}" \
  --project "${GCP_PROJECT}" \
  --region "${GCP_REGION}" \
  --set-secrets "DATABASE_URL=${DATABASE_SECRET_NAME}:latest" \
  --quiet

echo "Waiting for preview URL tagged '${PREVIEW_TAG}'"
preview_url=""
for attempt in {1..30}; do
  if service_json=$(gcloud run services describe "${SERVICE_NAME}" \
    --project "${GCP_PROJECT}" \
    --region "${GCP_REGION}" \
    --format json 2>/dev/null); then
    preview_url=$(jq -r --arg TAG "${PREVIEW_TAG}" '.status.traffic[]? | select(.tag==$TAG) | .url' <<<"$service_json" | head -n1)
    if [[ -n "$preview_url" && "$preview_url" != "null" ]]; then
      break
    fi
  fi
  sleep 5
done

if [[ -z "$preview_url" || "$preview_url" == "null" ]]; then
  echo "Preview URL for tag '${PREVIEW_TAG}' was not found" >&2
  exit 1
fi

echo "Preview URL: $preview_url"
