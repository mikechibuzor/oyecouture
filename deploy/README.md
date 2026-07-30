# Deploying Oyecouture

One server runs the whole stack. Caddy terminates TLS and routes `/api/*` to the
Go service and everything else to the Nuxt storefront. PostgreSQL is reachable
only on the internal Docker network — it publishes no ports.

```
Internet → Caddy (auto-TLS) ─┬─ /api/* → api  (Go, :8080)
                             └─ /*     → web  (Nuxt SSR, :3000)
                                              └─ postgres (internal only)
```

## First deploy

### 1. Prepare the server (once)

As root on a fresh Ubuntu 24.04 box:

```bash
# A non-root user for day-to-day work
adduser --disabled-password --gecos "" deploy
usermod -aG sudo deploy
mkdir -p /home/deploy/.ssh
cp /root/.ssh/authorized_keys /home/deploy/.ssh/
chown -R deploy:deploy /home/deploy/.ssh
chmod 700 /home/deploy/.ssh && chmod 600 /home/deploy/.ssh/authorized_keys

# Docker
curl -fsSL https://get.docker.com | sh
usermod -aG docker deploy

# Firewall: SSH and web only
ufw allow OpenSSH && ufw allow 80 && ufw allow 443 && ufw --force enable

# Key-only SSH
sed -i 's/^#\?PasswordAuthentication .*/PasswordAuthentication no/' /etc/ssh/sshd_config
sed -i 's/^#\?PermitRootLogin .*/PermitRootLogin prohibit-password/' /etc/ssh/sshd_config
systemctl restart ssh
```

### 2. Get the code and configure

```bash
su - deploy
git clone <repo-url> oyecouture && cd oyecouture
cp .env.production.example .env.production
openssl rand -base64 32   # use for POSTGRES_PASSWORD
nano .env.production
```

`SITE_HOST` must resolve to this server. Before the client delegates DNS, use
`<ip-with-dashes>.sslip.io` (e.g. `203-0-113-10.sslip.io`) — it resolves to the
IP automatically and still earns a real Let's Encrypt certificate.

### 3. Deploy

```bash
./deploy/deploy.sh
```

Caddy requests the certificate on first boot; allow a few seconds. Then check:

```bash
curl -fsS https://$SITE_HOST/api/healthz     # {"status":"ok",...}
curl -fsS -o /dev/null -w '%{http_code}\n' https://$SITE_HOST/
```

## Redeploying

Same command — it rebuilds changed images and restarts only what moved:

```bash
git pull && ./deploy/deploy.sh
```

## Moving to the client's domain

1. Point an `A` record at the server IP.
2. Change `SITE_HOST` in `.env.production`.
3. `./deploy/deploy.sh` — Caddy issues the new certificate automatically.

## Backups

`pg_dump` runs nightly into the `postgres-backups` volume, kept 14 days.

```bash
docker compose -f compose.prod.yaml exec backup ls -lh /backups          # list
docker compose -f compose.prod.yaml exec backup sh /usr/local/bin/backup.sh oye oyecouture   # run now
```

Restore:

```bash
gunzip -c /backups/oyecouture-<timestamp>.sql.gz \
  | docker compose -f compose.prod.yaml exec -T postgres psql -U oye -d oyecouture
```

Copy backups off the server periodically — a dump that only exists on the box it
protects is not a backup.

## Troubleshooting

| Symptom | Check |
|---|---|
| TLS not issued | `docker compose -f compose.prod.yaml logs caddy` — usually `SITE_HOST` does not resolve to this server, or 80/443 are blocked |
| 502 from the storefront | `logs web` — the Nuxt container failed to start |
| API errors | `logs api` — most often `DATABASE_URL` or a failed migration |
| Images missing | `R2_PUBLIC_BASE` must match where the files actually are |
