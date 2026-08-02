# 06 · Setup & Environment

## Accounts / credentials needed (chase these early)

| Service | Purpose | Who provides | Status |
|---|---|---|---|
| **Paystack** | active payments (test + live keys) | us | ✅ sandbox key supplied |
| **Nomba** | reserved payment provider/fallback (`clientId`, `clientSecret`, `accountId`) | client suggested Nomba; client to provide account | ⬜ pending |
| **Cloudflare R2** | image storage (S3-compatible keys) | us (free tier) | ⬜ pending |
| **SendByte** | transactional email (API key + verified sender domain for live delivery) | us | ✅ sandbox key supplied |
| **Oracle Cloud** | Always Free ARM VM (app + Postgres) | us (free tier) | ⬜ pending |
| **Neon** (optional) | free Postgres if not self-hosting on the VM | us | ⬜ optional |
| **WhatsApp business number** | the wa.me hand-off target | client | ⬜ pending |
| **Domain DNS access** | point the domain at the VM, HTTPS | client (owns domain) | ⬜ pending |
| **Sentry / GA** | errors + analytics | us (free) | ⬜ pending |

## Example environment variables (`.env.example`)

```bash
# Core
APP_ENV=development
API_BASE_URL=http://localhost:8080
WEB_BASE_URL=http://localhost:3000
PUBLIC_SITE_URL=http://localhost:3000        # used in emails, OG tags, pay links

# Postgres
DATABASE_URL=postgres://oye:oye@localhost:5432/oyecouture?sslmode=disable

# Payments — Paystack (active)
PAYMENT_PROVIDER=paystack                     # paystack | nomba (ghost-admin toggle later)
PAYSTACK_SECRET_KEY=
PAYSTACK_PUBLIC_KEY=

# Payments — Nomba (reserved)
NOMBA_CLIENT_ID=
NOMBA_CLIENT_SECRET=
NOMBA_ACCOUNT_ID=
NOMBA_ENV=sandbox                             # sandbox | production
NOMBA_WEBHOOK_SECRET=

# Cloudflare R2 (S3-compatible)
R2_ACCOUNT_ID=
R2_ACCESS_KEY_ID=
R2_SECRET_ACCESS_KEY=
R2_BUCKET=oyecouture-media
R2_PUBLIC_BASE=                               # public CDN base for images

# Email (SendByte)
EMAIL_PROVIDER=sendbyte
SENDBYTE_API_KEY=
EMAIL_FROM="Oyecouture <orders@oyecouture-domain>"
OWNER_NOTIFY_EMAIL=                           # where new-order alerts go

# WhatsApp (Phase 1 = deep link only; no API)
WHATSAPP_NUMBER=234XXXXXXXXXX                  # E.164 without '+', used in wa.me links

# Observability
SENTRY_DSN=
GA_MEASUREMENT_ID=
```

## Local dev

```bash
# 1. Postgres (docker is easiest)
docker run -d --name oye-pg -e POSTGRES_USER=oye -e POSTGRES_PASSWORD=oye -e POSTGRES_DB=oyecouture -p 5432:5432 postgres:16

# 2. API (Go)
cd api && cp ../.env.example .env   # fill values
go run ./cmd/server                 # migrations + seed from catalogue/product-data.js on first run

# 3. Web (Nuxt)
cd web && npm install && npm run dev # http://localhost:3000
```

To view the **original design** without the server: open any file in `../design-source/` directly in a browser (see `design-source/HOW-TO-VIEW.md`).

## Deploy (free tier)

- Provision the **Oracle Always Free ARM VM** (Ubuntu). Install Postgres (or use Neon), Go, Node, Caddy/Nginx.
- Build the Nuxt app (`npm run build`, Node server output) and the Go binary; run both under systemd. Caddy for HTTPS + reverse proxy (`/` → Nuxt, `/api` → Go).
- Point the client domain's DNS at the VM; Caddy auto-provisions HTTPS.
- Set the same env vars on the server (prod payment env, real SendByte sender, R2 keys).
- Keep hosting on our accounts for now (client declined running costs); migrate to paid managed Postgres + a bigger box when traffic warrants, quoted then.

## Gotchas
- **Gilmer font is commercial** — dev with Jost/Poppins behind a CSS var; resolve licensing before launch.
- **Never** mark an order `paid` on the payment redirect — only after server-side verify / verified webhook.
- Webhooks must be idempotent + signature-verified.
- Money in **kobo**. Prices in `product-data.js` are naira ×100.
- The **`/admin` area stays unlinked** and unmentioned to the client (ghost admin).
