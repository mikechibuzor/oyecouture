#!/usr/bin/env bash
# Poll origin/main and deploy when it moves. Run from cron on the server.
#
# This exists because the two obvious alternatives are both unavailable here:
# GitHub Actions is blocked by a failed payment authorization on the account,
# and deploying by hand needs a laptop that can reach the box, which the
# developer's ISP cannot. Having the server pull for itself sidesteps both, and
# also closes the original gap where work sat pushed-but-undeployed for days
# because nobody remembered to run the deploy.
#
# Install (the copy is deliberate, see below):
#   cp ~/oyecouture/deploy/autodeploy.sh ~/autodeploy.sh && chmod +x ~/autodeploy.sh
#   crontab -e   →   */5 * * * * /home/deploy/autodeploy.sh >> /home/deploy/autodeploy.log 2>&1
#
# Cron must run the COPY at ~/autodeploy.sh, never the one inside the repo.
# This script pulls, so running it from the repo would let git rewrite the file
# while bash is still reading it — bash reads scripts incrementally, so an
# updated script resumes at a byte offset that now points mid-line. Re-install
# the copy whenever this file changes.
set -euo pipefail

REPO="${HOME}/oyecouture"
LOCK="/tmp/oyecouture-autodeploy.lock"

log() { printf '%s  %s\n' "$(date -u '+%Y-%m-%dT%H:%M:%SZ')" "$*"; }

# A deploy builds two images on a single core and can outlast the 5 minute cron
# interval. Without this the next tick would start a second build against the
# same working copy and Docker daemon. Non-blocking: skip this tick rather than
# queue, since the next one is only minutes away.
exec 9>"$LOCK"
if ! flock -n 9; then
  exit 0
fi

cd "$REPO"

git fetch --quiet origin main

LOCAL="$(git rev-parse HEAD)"
REMOTE="$(git rev-parse origin/main)"

if [[ "$LOCAL" == "$REMOTE" ]]; then
  exit 0
fi

log "origin/main moved: ${LOCAL:0:7} -> ${REMOTE:0:7}"

# Whether anything the running site is built from actually changed. Handoff docs
# change on every chunk of work in this project and rebuild nothing, so treating
# them as a release would mean several minutes of pointless image building on a
# 1 vCPU box many times a day.
if git diff --name-only "$LOCAL" "$REMOTE" | grep -qvE '(\.md|\.gitignore)$'; then
  NEEDS_RELEASE=1
else
  NEEDS_RELEASE=0
fi

# --ff-only so a diverged working copy on the server stops here with a readable
# error instead of silently creating a merge commit on the box.
if ! git merge --ff-only origin/main --quiet; then
  log "ERROR: fast-forward failed. The server working copy has diverged from origin/main."
  log "       Resolve by hand: cd $REPO && git status"
  exit 1
fi

if [[ "$NEEDS_RELEASE" -eq 0 ]]; then
  log "documentation only, repository updated without rebuilding"
  exit 0
fi

log "deploying $(git log --oneline -1)"

if ./deploy/deploy.sh; then
  log "deploy OK: $(git rev-parse --short HEAD)"
else
  # deploy.sh exits nonzero when the API never reports healthy, so this is a
  # real failed rollout and not just noisy output.
  log "DEPLOY FAILED for $(git rev-parse --short HEAD)"
  log "       Inspect: cd $REPO && docker compose -f compose.prod.yaml --env-file .env.production logs api"
  exit 1
fi
