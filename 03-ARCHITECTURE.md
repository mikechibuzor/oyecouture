# 03 · Architecture

## Stack (decided)

- **Frontend:** Nuxt 3 (Vue 3, SSR on) — SSR for SEO on product pages. TypeScript.
- **Backend:** Go (net/http or a light router like chi; standard `$GOPATH`-style module rooted at `src/` is the team's Go convention). REST/JSON API.
- **Database:** PostgreSQL.
- **Async jobs:** River (Go job queue on Postgres) for emails and payment verification/reconciliation.
- **Images:** Cloudflare R2 (S3-compatible; zero egress). Go backend generates WebP size variants on upload; store keys in DB.
- **Email:** SendByte (transactional email; sandbox logs before live delivery is enabled).
- **Analytics/errors:** Google Analytics (or Plausible) + Sentry.

Rationale for full-custom (no Medusa/CMS): the client's needs will grow and change; we want full control and no third-party commerce framework lock-in. The trade-off — we hand-build the admin — is accepted and is exactly the "ghost admin" we're building anyway.

## Why these choices matter for the code

- **Money in kobo (int64).** `product-data.js` prices are in naira; store `price_kobo = naira * 100`. Never use floats for money.
- **SSR** the storefront/product pages for SEO; the cart/checkout can be client-side.
- **Guest-only** in Phase 1: cart lives in `localStorage`; orders are created server-side at confirm/checkout time with contact details captured then.

## Payment abstraction (build this carefully — it's a headline feature)

Define a provider-agnostic interface in Go and implement two adapters:

```go
type PaymentProvider interface {
    // Create a hosted-checkout order; return a redirect/checkout URL + provider ref.
    InitCheckout(ctx context.Context, o OrderPaymentInput) (CheckoutResult, error)
    // Server-side verify a transaction by reference (source of truth for "paid").
    VerifyTransaction(ctx context.Context, ref string) (PaymentStatus, error)
    // Parse + verify a webhook (signature check), return a normalized event.
    HandleWebhook(ctx context.Context, headers http.Header, body []byte) (WebhookEvent, error)
    Name() string
    Healthy(ctx context.Context) bool
}
```

- **Paystack = active provider** for the first payment phase. It owns hosted checkout, verify, and signed webhooks today.
- **Nomba = reserved provider** behind the same interface for a later fallback/primary switch if needed.
- **Selection:** config flag chooses primary; on init failure or unhealthy primary, fall back automatically. Expose the toggle in the ghost admin `settings`.
- **Golden rule:** an order becomes `paid` **only** after `VerifyTransaction`/verified webhook. Never on the redirect alone. Webhooks must be **idempotent** and **signature-verified**. Record every provider event in `payment_events`.

### Two payment entry points, one engine
1. **RTW checkout** → `InitCheckout` → redirect to hosted checkout → verify → `paid`.
2. **Bespoke** → every order has a **secure pay page** `oyecouture.com/pay/<token>` (opaque token, not the order id). Owner shares it on WhatsApp after measurements are agreed. The page shows the order summary + delivery fee and launches the same hosted checkout.

## Suggested repo layout

```
oyecouture/
  web/                      # Nuxt 3 storefront (+ /admin route group = ghost admin)
    pages/                  # index, [section], product/[slug], cart, checkout, pay/[token], confirm/[order], about, care/*
    pages/admin/            # ghost admin (unlinked)
    components/             # per DESIGN-SPEC component inventory
    composables/            # useCart (localStorage), useSearch, etc.
    assets/                 # fonts (Jost/Gilmer), css tokens
  api/                      # Go service
    cmd/server/main.go
    internal/
      products/  orders/  payments/ (interface + nomba/ + paystack/)
      measurements/  delivery/  settings/  admin/
      email/     (SendByte)  storage/ (R2)  jobs/ (River workers)
    migrations/             # SQL migrations (see 04-DATA-MODEL.md)
  README.md  .env.example
```
(Adjust to taste; the point is the payments package with a clean interface + two adapters, and an unlinked admin area.)

## Hosting — ₦0/month at launch (client declined running costs)

- **App + DB:** Oracle Cloud **Always Free** ARM VM (always-on → no cold starts) running Nuxt + Go + PostgreSQL on the one box. (Alternatively Postgres on **Neon** free tier.) Host it on **our (Uzor & Co.) free-tier accounts** for now.
- **Images:** Cloudflare R2 free (10 GB). **Email:** SendByte sandbox/free tier while available. **Analytics/errors:** free tiers.
- **Avoid** Render's free Postgres — it **expires after 90 days**.
- When traffic outgrows the free tiers, migrate to paid (~₦15–30k/month) and quote it then. Design deploys so this migration is painless (managed Postgres URL via env, stateless app, images already on R2).
- Staging + prod; HTTPS on the client's domain; SEO (sitemap, OG tags, indexable product pages).

## Non-functionals
Lighthouse ≥ 90 on product pages (lazy-load images, WebP). Idempotent + signature-verified webhooks. Reconciliation job re-verifies pending payments and alerts the owner on mismatch. Store all money in kobo.
