# HANDOFF-STATE.md — the baton

**Read this first. Keep it true at all times.** This is the single source of truth for where we are and what to do next. Update it continuously (a usage limit may cut you off without warning).

---

## Current status
- **Phase:** 1 (MVP). First 40% milestone **paid**; build authorized.
- **Milestone:** A — Foundations & catalogue (see `05-BUILD-ROADMAP.md`).
- **Step:** A.4 in progress: typed storefront shell foundations are being implemented.
- **Last worked by:** Codex · 2026-07-29.
- **Repo / branch:** `/Users/chibuzoriwuagwu/Downloads/share/oyecouture` · `main`.
- **Last implementation commits:** `b42bb60` (product HTTP endpoints), `36f40bc` (product read service), `bc20818` (catalogue seed), `2e66fbc` (database migrations), `576ee77` (design tokens/Jost).

## ▶ Next action (the ONE thing to do next)
Complete Milestone A step A.4 only: build and verify the Nuxt storefront shell from `02-DESIGN-SPEC.md` (`AnnouncementBar`, `NavBar` with dropdowns/mobile menu, and `Footer`) using the existing SCSS tokens. Do not start A.5 storefront pages in the same step.

## Done & verified
- Handoff docs `00`–`06`, `catalogue/`, `design-source/`, `reference-docs/` are complete and current.
- Git repository initialized on `main`; commits authored as `mikechibuzor` with no co-author trailers.
- Repository and generated service resources are consistently named `oyecouture`.
- Root `.env.example` matches `06-SETUP-AND-ENV.md`; Compose project is explicitly named `oyecouture` and runs PostgreSQL 16.
- Nuxt `3.21.10` scaffolded in `web/` with strict TypeScript, Node `22.12.0`, shared types/constants, and a typed API-health composable.
- Frontend styling uses modular SCSS with design-system tokens, reusable breakpoint mixins, global rules, page partials, and the `sass-embedded` compiler.
- Approved typography and color foundations are aligned: self-hosted Jost Variable stands in for Gilmer, the type scale reaches a 60px display step, and champagne/gold/taupe tokens match the design spec.
- Collaboration rules require one-line intent comments above functions and computed properties, with concise walkthrough comments at meaningful transitions inside multi-step functions.
- Embedded versioned migrations create `products`, `product_images`, `settings`, and `delivery_zones` with money, section, image-position, and uniqueness constraints; API startup applies pending versions before serving.
- A deterministic validated generator seeds the 27 canonical products from `catalogue/product-data.js`, mapping naira to integer kobo and preserving RTW/bespoke `has_sizes` behavior and catalogue sort order.
- Public product endpoints return active products only: `GET /products` supports validated `section` and literal `q` search filters, while `GET /products/{slug}` provides detail lookup with stable JSON errors and ordered image arrays.
- Go API scaffolded in `api/` with config, PostgreSQL pool, CORS, graceful shutdown, and `GET /healthz`.
- Verified: `go test ./...`, Go build, `npm run typecheck`, and `npm run build` all pass.
- Runtime verified: `/healthz` returns HTTP 200 with API/database `ok`; Nuxt SSR renders web/API/PostgreSQL as connected.
- **Claude review sign-off (2026-07-29):** A.1–A.3 independently reviewed and re-verified (`go vet`/`build`/`test` all pass). Approved with no changes. Claude's three earlier design-token findings (display font, champagne `#e0cf9c`, type scale) were already resolved by Codex in `576ee77`; no fix commit needed. Golden invariants confirmed (integer `priceKobo`, hidden `R2Key`, validated inputs, active-only reads, RTW/bespoke split).

## In progress (exact state)
- Local PostgreSQL, product-enabled Go API, and Nuxt dev services are running. A.4 shell implementation is in progress; A.5 remains untouched.

## How to run / verify current state
```bash
cd /Users/chibuzoriwuagwu/Downloads/share/oyecouture
nvm use
docker compose up -d postgres

cd api && go run ./cmd/server
cd ../web && npm install && npm run dev
```
- Seed integrity: `cd catalogue && npm run check:seed`
- Web: `http://localhost:3000`
- API health: `http://localhost:8080/healthz`
- Verify API: `curl -sS http://localhost:8080/healthz`
- Verify products: `curl -sS 'http://localhost:8080/products?section=casuals-men&q=linen'`

## Blockers / waiting on the user (client deliverables)
- Nomba credentials (sandbox first), WhatsApp business number, delivery-fee amounts per zone, product photography, customer-care page copy, domain DNS access. None block A.4, but they gate later steps.

## Decisions / deviations from the numbered docs
- Health proof uses `/healthz` rather than `/health`; it reports both API and PostgreSQL state.
- Nuxt is pinned to `3.21.10`; Node is pinned to `22.12.0` because this Nuxt release requires Node `^20.19.0 || >=22.12.0`.
- Frontend domain values, API paths, configuration defaults, and reusable copy live in imported types/constants. Prefer `as const` literal unions; use enums where runtime enum semantics improve clarity.
- SCSS structure follows the maintainable organization used by `business-webapp-v4`, while using modern Sass `@use` modules instead of deprecated `@import`.
- Jost Variable is self-hosted through `@fontsource-variable/jost` as the licensed free Gilmer stand-in; keep the family behind shared tokens so licensed Gilmer can replace it before launch.
- Database migrations use embedded `golang-migrate` v4.18.3, pinned to preserve the established Go 1.23 baseline.
- The canonical source contains 27 products (17 bespoke, 10 RTW), correcting the stale exact count of 28 in the handoff catalogue prose; approximate “~28” scope references remain valid.
- `product_images`, `settings`, and `delivery_zones` start empty because R2 keys, delivery fees, and client settings have not been supplied; schema constraints are ready for those later inputs.
- Product responses use camelCase fields, integer `priceKobo`, `images: []` for missing photography, and `{products,count}` / `{product}` envelopes for strict frontend contracts.
- Storefront queries reject unsupported sections, malformed slugs, and searches over 100 characters; `%`, `_`, and backslashes are escaped as literal search characters.
- Product image URLs are derived server-side from `R2_PUBLIC_BASE`, keeping storage keys out of public JSON.

## Reminders (non-negotiable)
Money in kobo (int64, no floats) · `paid` only after server-side verify · guest checkout only in phase 1 · mobile-first · ghost admin stays unlinked/unmentioned to client · commits as `mikechibuzor`, no `Co-Authored-By`.
