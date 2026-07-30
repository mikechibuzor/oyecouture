#!/usr/bin/env bash
# Build and (re)start the production stack. Safe to re-run — this is also the
# redeploy path for every subsequent release.
#
#   ./deploy/deploy.sh
#
# Expects .env.production in the repo root (see .env.production.example).
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_ROOT"

ENV_FILE=".env.production"
COMPOSE_FILE="compose.prod.yaml"

if [[ ! -f "$ENV_FILE" ]]; then
  echo "error: $ENV_FILE not found. Copy .env.production.example and fill it in." >&2
  exit 1
fi

# Fail early with a readable message rather than a confusing compose error.
for required in SITE_HOST ACME_EMAIL POSTGRES_PASSWORD; do
  if ! grep -qE "^${required}=.+" "$ENV_FILE"; then
    echo "error: $required is not set in $ENV_FILE" >&2
    exit 1
  fi
done

echo "==> Building images"
docker compose --env-file "$ENV_FILE" -f "$COMPOSE_FILE" build

echo "==> Starting stack"
docker compose --env-file "$ENV_FILE" -f "$COMPOSE_FILE" up -d --remove-orphans

echo "==> Waiting for the API to report healthy"
for attempt in $(seq 1 30); do
  if docker compose --env-file "$ENV_FILE" -f "$COMPOSE_FILE" exec -T api /server --version >/dev/null 2>&1 \
    || curl -fsS "http://localhost/api/healthz" >/dev/null 2>&1; then
    echo "    API is responding"
    break
  fi
  if [[ "$attempt" -eq 30 ]]; then
    echo "warning: API did not respond in time; check: docker compose -f $COMPOSE_FILE logs api" >&2
  fi
  sleep 2
done

SITE_HOST_VALUE="$(grep -E '^SITE_HOST=' "$ENV_FILE" | cut -d= -f2-)"

echo
echo "==> Deployed. https://${SITE_HOST_VALUE}"
echo "    Health:  https://${SITE_HOST_VALUE}/api/healthz"
echo "    Logs:    docker compose -f ${COMPOSE_FILE} logs -f"
