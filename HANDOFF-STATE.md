# HANDOFF-STATE.md — the baton

**Read this first. Keep it true at all times.** This is the single source of truth for where we are and what to do next. Update it continuously (a usage limit may cut you off without warning).

---

## Current status
- **Phase:** 1 (MVP). First 40% milestone **paid**; build authorized.
- **Milestone:** A — Foundations & catalogue (see `05-BUILD-ROADMAP.md`).
- **Step:** A.5 complete: the API-backed shop, category, and product detail pages are implemented and verified.
- **Last worked by:** Codex · 2026-07-29.
- **Repo / branch:** `/Users/chibuzoriwuagwu/Downloads/share/oyecouture` · `main`.
- **Last implementation commits:** `7c7bca4` (browser-QA polish), `8736afc` (product detail), `ed4f3d9` (shop/category pages), `82d45fa` (catalogue components), `68417f9` (typed catalogue foundation).

## ▶ Next action (the ONE thing to do next)
**Wait for user review and explicit authorization before starting A.6.** The next roadmap step is A.6 only: the instant search overlay backed by `GET /products?q=`. Do not begin the cart; cart state starts at Milestone B step 8.

## Done & verified
- Handoff docs `00`–`06`, `catalogue/`, `design-source/`, `reference-docs/` are complete and current.
- Git repository initialized on `main`; commits authored as `mikechibuzor` with no co-author trailers.
- Repository and generated service resources are consistently named `oyecouture`.
- Root `.env.example` matches `06-SETUP-AND-ENV.md`; Compose project is explicitly named `oyecouture` and runs PostgreSQL 16.
- Nuxt `3.21.10` scaffolded in `web/` with strict TypeScript, Node `22.23.1`, shared types/constants, and a typed API-health composable.
- Frontend styling uses modular SCSS with design-system tokens, reusable breakpoint mixins, global rules, page partials, and the `sass-embedded` compiler.
- Approved typography and color foundations are aligned: self-hosted Jost Variable stands in for Gilmer, the type scale reaches a 60px display step, and champagne/gold/taupe tokens match the design spec.
- Collaboration rules require one-line intent comments above functions and computed properties, with concise walkthrough comments at meaningful transitions inside multi-step functions.
- Embedded versioned migrations create `products`, `product_images`, `settings`, and `delivery_zones` with money, section, image-position, and uniqueness constraints; API startup applies pending versions before serving.
- A deterministic validated generator seeds the 27 canonical products from `catalogue/product-data.js`, mapping naira to integer kobo and preserving RTW/bespoke `has_sizes` behavior and catalogue sort order.
- Public product endpoints return active products only: `GET /products` supports validated `section` and literal `q` search filters, while `GET /products/{slug}` provides detail lookup with stable JSON errors and ordered image arrays.
- Go API scaffolded in `api/` with config, PostgreSQL pool, CORS, graceful shutdown, and `GET /healthz`.
- Verified: `go test ./...`, Go build, `npm run typecheck`, and `npm run build` all pass.
- Runtime verified: `/healthz` returns HTTP 200 with API/database `ok`; Nuxt SSR renders web/API/PostgreSQL as connected.
- A.4 storefront shell is complete: responsive announcement bar, approved Oyecouture logo, desktop dropdown navigation, accessible mobile menu, typed action/icon maps, charcoal footer, and modular component/layout SCSS.
- A.4 browser verification passes at desktop and 370px mobile widths with no horizontal overflow, working menu open/close behavior, 17 hydrated Lucide icons, all services connected, and no new console errors or hydration warnings.
- A.5 storefront pages are complete: immersive Home/Shop hero, approved category imagery, three category tiles, filter chips, price/newest sorting, API-backed ProductCard grid, Suits/Trad/Casuals Men/Casuals Women routes, and one slug-based product detail route.
- A.5 preserves the purchase split: RTW detail pages expose a typed `S/M/L/XL` selector with `M` selected by default; bespoke pages render no size selector and show the approved WhatsApp sizing/fabric/delivery note.
- A.5 uses API product images first and approved WebP imagery extracted from `design-source/` as temporary fallbacks; missing client photography renders the approved sand placeholder with product name.
- A.5 browser verification passes at 1280px desktop and 370px mobile: 27 cards, 3 category tiles, correct 8-item Suits route, working Aso-Oke filter and descending-price sort, working RTW gallery/size controls, correct bespoke branch, zero horizontal overflow, and no console or hydration errors.
- A.5 verification passes: `npm run typecheck`, `npm run build`, `npm run check:seed`, `go test ./...`, live API queries, and responsive in-app browser checks.
- **Claude review sign-off (2026-07-29):** A.1–A.3 independently reviewed and re-verified (`go vet`/`build`/`test` all pass). Approved with no changes. Claude's three earlier design-token findings (display font, champagne `#e0cf9c`, type scale) were already resolved by Codex in `576ee77`; no fix commit needed. Golden invariants confirmed (integer `priceKobo`, hidden `R2Key`, validated inputs, active-only reads, RTW/bespoke split).
- **Claude review sign-off — A.4 (2026-07-29):** Storefront shell reviewed and approved. Design fidelity verified against `design-source/` (announcement copy exact; nav Home/Trad/Suits/Casuals▾/Collection▾; Casuals→Men/Women; the COLLECTION names "Dear Groom / Wei-Xin / OYE — Man of Legacy" confirmed present in all 7 design files — extracted, not invented; footer Shop/About/Customer Care + IG/X/TikTok/WhatsApp + legal links match). Code quality strong: nav state in a composable, full a11y (aria-controls/expanded, Escape, click-outside, close-on-route-change), lucide tree-shaken icons, search/account/wishlist/cart correctly emit events for future features. One required A.5 follow-up flagged: remove `:external` per link as pages land (see Decisions). Note: Claude's local `npm run typecheck` was inconclusive due to a missing `oxc-parser` native dep corrupting a generated `.nuxt` file (env artifact, not a source defect); Codex's own typecheck/build/`go test` pass.

## In progress (exact state)
- Local PostgreSQL, the product-enabled Go API, and Nuxt dev server are running. A.5 is complete; wait for user review and A.6 authorization.

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
- Nomba credentials (sandbox first), WhatsApp business number, delivery-fee amounts per zone, product photography, customer-care page copy, domain DNS access. None block A.6, but they gate final content and later purchase/deployment steps.

## Decisions / deviations from the numbered docs
- Health proof uses `/healthz` rather than `/health`; it reports both API and PostgreSQL state.
- Nuxt is pinned to `3.21.10`; Node is pinned to `22.23.1` and the web engine floor is `^20.19.0 || >=22.13.0` for stable `util.styleText` support across the installed build toolchain.
- Frontend domain values, API paths, configuration defaults, and reusable copy live in imported types/constants. Prefer `as const` literal unions; use enums where runtime enum semantics improve clarity.
- SCSS structure follows the maintainable organization used by `business-webapp-v4`, while using modern Sass `@use` modules instead of deprecated `@import`.
- Jost Variable is self-hosted through `@fontsource-variable/jost` as the licensed free Gilmer stand-in; keep the family behind shared tokens so licensed Gilmer can replace it before launch.
- Database migrations use embedded `golang-migrate` v4.18.3, pinned to preserve the established Go 1.23 baseline.
- The canonical source contains 27 products (17 bespoke, 10 RTW), correcting the stale exact count of 28 in the handoff catalogue prose; approximate “~28” scope references remain valid.
- `product_images`, `settings`, and `delivery_zones` start empty because R2 keys, delivery fees, and client settings have not been supplied; schema constraints are ready for those later inputs.
- Product responses use camelCase fields, integer `priceKobo`, `images: []` for missing photography, and `{products,count}` / `{product}` envelopes for strict frontend contracts.
- Storefront queries reject unsupported sections, malformed slugs, and searches over 100 characters; `%`, `_`, and backslashes are escaped as literal search characters.
- Product image URLs are derived server-side from `R2_PUBLIC_BASE`, keeping storage keys out of public JSON.
- The written design spec's charcoal footer takes precedence over the white footer shown in the static shop export.
- Completed A.5 destinations use Nuxt client-side navigation. Future unimplemented destinations remain external anchors until their pages exist, preventing false router warnings without creating placeholders early.
- ProductCard and product-detail `ADD TO BAG` controls are intentionally present but disabled in A.5; bag mutation and persistence begin at Milestone B step 8, preventing an incomplete cart implementation from leaking into this step.
- Approved design-source WebPs are versioned local fallbacks while `product_images` is empty. API/R2 media always takes precedence, so client photography can replace the fallbacks without changing page components.
- Nuxt's unused experimental app manifest is disabled to avoid a dev-only Vite alias error. Official Lucide icon data is rendered through a typed local Vue component to keep SSR hydration deterministic on Vue `3.5.40`.

## Reminders (non-negotiable)
Money in kobo (int64, no floats) · `paid` only after server-side verify · guest checkout only in phase 1 · mobile-first · ghost admin stays unlinked/unmentioned to client · commits as `mikechibuzor`, no `Co-Authored-By`.
