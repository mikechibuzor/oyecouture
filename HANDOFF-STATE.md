# HANDOFF-STATE.md — the baton

**Read this first. Keep it true at all times.** This is the single source of truth for where we are and what to do next. Update it continuously (a usage limit may cut you off without warning).

---

## Current status
- **Phase:** 1 (MVP). First 40% milestone **paid**; build authorized.
- **Milestone:** A — Foundations & catalogue (see `05-BUILD-ROADMAP.md`).
- **Step:** A.1 complete and verified. Stop boundary reached; A.2 has not started.
- **Last worked by:** Codex · 2026-07-29.
- **Repo / branch:** `/Users/chibuzoriwuagwu/Downloads/share/oyecouture` · `main`.
- **Last implementation commits:** `e09120a` (SCSS), `21d6f8f` (Nuxt), `822419f` (Go API), `6c5f6a6` (root dev foundation).

## ▶ Next action (the ONE thing to do next)
Wait for user authorization, then begin Milestone A step A.2 only: add migrations for `products`, `product_images`, `settings`, and `delivery_zones`, then seed from `catalogue/product-data.js`. Do not start A.3 endpoints in the same step.

## Done & verified
- Handoff docs `00`–`06`, `catalogue/`, `design-source/`, `reference-docs/` are complete and current.
- Git repository initialized on `main`; commits authored as `mikechibuzor` with no co-author trailers.
- Repository and generated service resources are consistently named `oyecouture`.
- Root `.env.example` matches `06-SETUP-AND-ENV.md`; Compose project is explicitly named `oyecouture` and runs PostgreSQL 16.
- Nuxt `3.21.10` scaffolded in `web/` with strict TypeScript, Node `22.12.0`, shared types/constants, and a typed API-health composable.
- Frontend styling uses modular SCSS with design-system tokens, reusable breakpoint mixins, global rules, page partials, and the `sass-embedded` compiler.
- Go API scaffolded in `api/` with config, PostgreSQL pool, CORS, graceful shutdown, and `GET /healthz`.
- Verified: `go test ./...`, Go build, `npm run typecheck`, and `npm run build` all pass.
- Runtime verified: `/healthz` returns HTTP 200 with API/database `ok`; Nuxt SSR renders web/API/PostgreSQL as connected.

## In progress (exact state)
- Local PostgreSQL, Go API, and Nuxt dev services are running for review. No A.2 work is in progress.

## How to run / verify current state
```bash
cd /Users/chibuzoriwuagwu/Downloads/share/oyecouture
nvm use
docker compose up -d postgres

cd api && go run ./cmd/server
cd ../web && npm install && npm run dev
```
- Web: `http://localhost:3000`
- API health: `http://localhost:8080/healthz`
- Verify API: `curl -sS http://localhost:8080/healthz`

## Blockers / waiting on the user (client deliverables)
- Nomba credentials (sandbox first), WhatsApp business number, delivery-fee amounts per zone, product photography, customer-care page copy, domain DNS access. None block A.2, but they gate later steps.

## Decisions / deviations from the numbered docs
- Health proof uses `/healthz` rather than `/health`; it reports both API and PostgreSQL state.
- Nuxt is pinned to `3.21.10`; Node is pinned to `22.12.0` because this Nuxt release requires Node `^20.19.0 || >=22.12.0`.
- Frontend domain values, API paths, configuration defaults, and reusable copy live in imported types/constants. Prefer `as const` literal unions; use enums where runtime enum semantics improve clarity.
- SCSS structure follows the maintainable organization used by `business-webapp-v4`, while using modern Sass `@use` modules instead of deprecated `@import`.

## Reminders (non-negotiable)
Money in kobo (int64, no floats) · `paid` only after server-side verify · guest checkout only in phase 1 · mobile-first · ghost admin stays unlinked/unmentioned to client · commits as `mikechibuzor`, no `Co-Authored-By`.
