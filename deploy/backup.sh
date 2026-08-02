#!/bin/sh
# Nightly database dump with a 14-day retention window.
# Usage: backup.sh <postgres-user> <postgres-db>
set -eu

POSTGRES_USER="$1"
POSTGRES_DB="$2"
BACKUP_DIR="/backups"
RETENTION_DAYS=14
TIMESTAMP="$(date -u +%Y%m%dT%H%M%SZ)"
BACKUP_FILE="${BACKUP_DIR}/oyecouture-${TIMESTAMP}.sql.gz"
WORKING_DUMP="${BACKUP_DIR}/.dump.$$"

mkdir -p "$BACKUP_DIR"
trap 'rm -f "$WORKING_DUMP" "${BACKUP_FILE}.partial"' EXIT

# Dump uncompressed first so pg_dump's own exit status is visible. Piping it
# straight into gzip hides that status: gzip succeeds either way, so a failed
# dump still produced a valid-looking archive of nothing.
if ! pg_dump -h postgres -U "$POSTGRES_USER" -d "$POSTGRES_DB" > "$WORKING_DUMP"; then
  echo "backup failed: pg_dump exited non-zero" >&2
  exit 1
fi

# A dump with no tables is worse than no dump at all: it sits there looking like
# a backup and restores an empty database. That is exactly what a dump taken
# before the API has applied its migrations produces.
if ! grep -q 'CREATE TABLE' "$WORKING_DUMP"; then
  echo "backup rejected: the dump contains no tables" >&2
  exit 1
fi

gzip -c "$WORKING_DUMP" > "${BACKUP_FILE}.partial"
mv "${BACKUP_FILE}.partial" "$BACKUP_FILE"
echo "backup written: ${BACKUP_FILE} ($(wc -c < "$BACKUP_FILE") bytes)"

find "$BACKUP_DIR" -name 'oyecouture-*.sql.gz' -mtime "+${RETENTION_DAYS}" -delete
