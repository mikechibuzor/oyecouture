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

# The password is interpolated into DATABASE_URL, so URL-reserved characters
# corrupt the connection string. Caught here because the failure otherwise
# appears only as an API crash loop after a full build.
POSTGRES_PASSWORD_VALUE="$(grep -E '^POSTGRES_PASSWORD=' "$ENV_FILE" | head -1 | cut -d= -f2-)"
if printf '%s' "$POSTGRES_PASSWORD_VALUE" | grep -q '[/@:?#]'; then
  echo "error: POSTGRES_PASSWORD contains a URL-reserved character (/ @ : ? #)." >&2
  echo "       It is interpolated into DATABASE_URL and will break it." >&2
  echo "       Generate a safe one with: openssl rand -hex 32" >&2
  exit 1
fi

echo "==> Building images"
docker compose --env-file "$ENV_FILE" -f "$COMPOSE_FILE" build

echo "==> Starting stack"
docker compose --env-file "$ENV_FILE" -f "$COMPOSE_FILE" up -d --remove-orphans

echo "==> Waiting for the API to report healthy"
# Pinned to loopback through the real hostname. Requesting http://localhost
# instead looks correct but proves nothing: Caddy answers plain HTTP with a 308
# to HTTPS, and curl counts an unfollowed redirect as success, so the check
# exited 0 without ever reaching the API. Using the configured hostname keeps
# TLS SNI matching the certificate Caddy holds.
SITE_HOST_VALUE="$(grep -E '^SITE_HOST=' "$ENV_FILE" | cut -d= -f2-)"
API_HEALTHY=0
for attempt in $(seq 1 30); do
  # The body is asserted, not just the status: a 200 carrying a failed database
  # check is exactly the rollout this gate exists to stop.
  if curl -fsS --resolve "${SITE_HOST_VALUE}:443:127.0.0.1" \
      "https://${SITE_HOST_VALUE}/api/healthz" 2>/dev/null | grep -q '"status":"ok"'; then
    echo "    API is healthy"
    API_HEALTHY=1
    break
  fi
  sleep 2
done

if [[ "$API_HEALTHY" -ne 1 ]]; then
  echo "error: the API did not report healthy after 60s." >&2
  echo "       Inspect with: docker compose -f $COMPOSE_FILE logs api" >&2
  exit 1
fi

# Every build leaves cache and a superseded image behind. Left alone that grows
# by roughly half a gigabyte per deploy and eventually fills the disk, which
# surfaces months later as a build failing for no obvious reason. Capped by
# size rather than age: age-based retention keeps everything during a busy
# week of deploys, which is exactly when it grows fastest. 4GB keeps
# incremental builds fast while bounding the worst case.
echo "==> Reclaiming disk from old builds"
docker builder prune --force --max-used-space 4GB >/dev/null 2>&1 \
  || docker builder prune --force --filter until=168h >/dev/null 2>&1 \
  || true
docker image prune --force >/dev/null 2>&1 || true
echo "    $(df -h / | awk 'NR==2{print $4" free of "$2}')"

SITE_HOST_VALUE="$(grep -E '^SITE_HOST=' "$ENV_FILE" | cut -d= -f2-)"

echo
echo "==> Deployed. https://${SITE_HOST_VALUE}"
echo "    Health:  https://${SITE_HOST_VALUE}/api/healthz"
echo "    Logs:    docker compose -f ${COMPOSE_FILE} logs -f"
