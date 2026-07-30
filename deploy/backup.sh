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

mkdir -p "$BACKUP_DIR"

# Dump to a temporary file first so a failed run never leaves a truncated
# backup that looks valid.
pg_dump -h postgres -U "$POSTGRES_USER" -d "$POSTGRES_DB" \
  | gzip > "${BACKUP_FILE}.partial"

mv "${BACKUP_FILE}.partial" "$BACKUP_FILE"
echo "backup written: ${BACKUP_FILE}"

find "$BACKUP_DIR" -name 'oyecouture-*.sql.gz' -mtime "+${RETENTION_DAYS}" -delete
