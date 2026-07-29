# WORKLOG — append-only (newest on top)

One short entry per work session or meaningful chunk. Format:
`## <ISO datetime> · <Codex|Claude> · <one-line summary>` then bullets: **Did / State / Next / Blockers**.

---

## 2026-07-29T23:30:00+01:00 · Claude · Motion & microinteraction system researched and specified
- **Did:** Extensive research pass (2026 luxury/fashion ecommerce motion, View Transitions API, CSS scroll-driven animations, motion-library comparison, duration/easing standards, WCAG motion criteria, INP/compositor performance, mobile touch + LQIP patterns). Wrote **`07-MOTION-AND-MICROINTERACTIONS.md`**: 5 principles, a full `_motion.scss` token scale (duration/easing/displacement/stagger), 3 signature moments, a component-by-component spec (shell, grid, detail, cart, buttons, search, forms), a CSS-first 5-tier implementation stack, accessibility + performance rules, a banned list, and phased delivery mapped to `05-BUILD-ROADMAP.md`. Indexed it in `00-README-START-HERE.md`.
- **State:** Spec only — **no code written**. Current motion in the build is minimal (chevron 160ms, mobile menu 150ms, tile scale 1.025/300ms, hero marquee 30s). No motion library installed, and the spec says keep it that way unless justified.
- **Next:** M1 foundation from §9 (motion tokens + global reduced-motion block + `:active` tap feedback + retune existing transitions + **fix the marquee WCAG 2.2.2 issue**) — small and low-risk, pairs well with A.6. Then M2 transactional motion ships *with* the cart in Milestone B.
- **Blockers:** ⚠️ **WCAG 2.2.2 finding (launch blocker):** the infinite hero marquee needs a pause/stop mechanism; `prefers-reduced-motion` alone does not satisfy the criterion. Detail in §6.2.

## 2026-07-29T21:10:00+01:00 · Claude · Typography fidelity fix — self-hosted Gilmer + letter-spacing scale
- **Did:** User flagged font/letter-spacing/weight drift vs the Claude design. Diagnosed: letter-spacing was absent site-wide (only a `letter-spacing:0` on the debug page), and the face was Jost (a stand-in), not the design's Gilmer. Extracted the design's exact type system from `design-source` (tracking .01–.32em by role) and the four embedded **Gilmer** OTFs → converted to WOFF2 in `web/public/fonts/`. Added `_fonts.scss` (@font-face, Gilmer weight ranges 300-400 / 500-600 / 700 / 800), Gilmer family tokens + a `--tracking-*` scale in `_typography.scss`, and applied tracking across announcement bar, nav, footer, hero, category, product detail, and catalogue per the design. Removed the now-unused Jost css import from `nuxt.config.ts`.
- **State:** SCSS compiles clean (sass-embedded), fonts serve 200 from the dev server, letter-spacing declarations 1→27. Committed `26ebd64` as `mikechibuzor` (no co-author). **Visual side-by-side NOT yet done** — the browser preview tool was unavailable this session; values are extracted directly from the design so fidelity should be high, but a rendered check vs `design-source/` is still pending.
- **Next:** When the browser tool is back, render `http://localhost:3000/` vs the design and fine-tune any tracking/weight that looks off; then resume the authorized A.6 (search overlay). Services are already running (Postgres/API/web).
- **Blockers:** ⚠️ **Gilmer licensing** — Gilmer is commercial (Fontfabric). User chose to self-host it now for an exact match; a webfont licence is required before production launch (see Decisions). `@fontsource-variable/jost` dep is now unused and can be pruned.

## 2026-07-29T19:44:31+01:00 · Codex · Milestone A.5 storefront pages completed
- **Did:** Added strict product enums/types/constants, API catalogue composables, approved fallback imagery, reusable ProductCard/filter/sort/gallery/size/bespoke components, the Home/Shop experience, four category routes, and slug-based product details; converted completed shell links to SPA navigation.
- **State:** `npm run typecheck`, production build, catalogue seed validation, and `go test ./...` pass. Desktop/mobile browser QA confirms 27 products, correct filtering/sorting, working RTW size/gallery interactions, the bespoke no-size WhatsApp-note branch, zero overflow, and no console/hydration errors.
- **Next:** Wait for user review and explicit authorization, then implement A.6 search overlay only. Cart begins at Milestone B step 8.
- **Blockers:** None for A.6. Client photography remains outstanding, so approved design-source imagery and sand placeholders are active fallbacks.

## 2026-07-29T17:05:00+01:00 · Claude · Reviewed A.4; approved; A.5 authorized by user
- **Did:** Reviewed the storefront shell against `02-DESIGN-SPEC.md` and `design-source/`. Verified the COLLECTION names ("Dear Groom / Wei-Xin / OYE — Man of Legacy") appear in all 7 design files (extracted, not invented), announcement copy is exact, and nav/footer/social/legal structure matches. Confirmed a11y and typed contracts. Recorded sign-off in `HANDOFF-STATE.md`.
- **State:** Approved — no changes needed. Only `web/pages/index.vue` exists, so nav links are intentionally `:external` (documented). Local web typecheck was inconclusive here (missing `oxc-parser` native dep corrupts a generated `.nuxt` file — env artifact, not a source defect); Codex's typecheck/build/`go test` pass.
- **Next:** **A.5 AUTHORIZED by user.** Build the storefront pages on the shell; **remove `:external` from each nav/footer link as its page lands** (required, see baton Decisions). Do not start A.6 (cart) in the same step. Handing the baton to Codex.
- **Blockers:** None for A.5 structure. Product photography + care-page copy gate final content.

## 2026-07-29T16:41:08+01:00 · Codex · Milestone A.4 storefront shell completed
- **Did:** Built the responsive announcement bar, approved brand logo, desktop dropdown navigation, accessible mobile menu, charcoal footer, typed route/action/icon contracts, and modular SCSS shell; aligned the Node runtime and removed dev-only route, manifest, and icon-hydration noise.
- **State:** `npm run typecheck`, production build, and `go test ./...` pass. Desktop and 370px mobile browser checks show no overflow, working menu interactions, 17 hydrated Lucide icons, no new console warnings/errors, and Nuxt/API/PostgreSQL connected.
- **Next:** Wait for user approval, then begin A.5 storefront pages only.
- **Blockers:** None for A.5 implementation; real product photography and the existing client deliverables still gate final content and later steps.

## 2026-07-29T14:20:00+01:00 · Claude · Reviewed A.1–A.3; independently verified; approved
- **Did:** Full review of the scaffold + catalogue + product read API against `02`–`04`. Re-ran `go vet ./...` (clean), `go build ./...` (clean), `go test ./...` (health + products pass). Read migrations, code-generated seed, and product model/service/handler/routes.
- **State:** Approved — no changes needed. The three design-token findings I raised earlier (Georgia-serif display font, missing champagne `#e0cf9c`, thin/inverted type scale) were **already resolved by Codex in `576ee77`** before I looked, so no fix commit was warranted. Confirmed golden invariants hold: money is integer `priceKobo` (₦850,000 → `85000000`), `R2Key` hidden from JSON, inputs validated pre-DB, active-only reads, stable error envelopes, `has_sizes` RTW/bespoke split intact.
- **Next:** Unchanged — A.4 storefront shell awaits user authorization. Handing the baton back to Codex.
- **Blockers:** None new. Client deliverables still gate later steps.

## 2026-07-29T14:01:01+01:00 · Codex · Milestone A.3 product read API completed
- **Did:** Added typed product models, active-only PostgreSQL repository queries, validated service filters, R2 image URL mapping, JSON handlers, and routes for product listing/search/detail.
- **State:** Go tests, vet, and binary build pass. Live API checks return 27 products (17 bespoke, 10 RTW), correct combined filters/details, literal wildcard search, stable 400/404 errors, CORS, and integer-kobo prices; Nuxt SSR remains connected.
- **Next:** Wait for approval, then begin A.4 storefront shell components only.
- **Blockers:** None for A.4; real product photography remains a later client dependency.

## 2026-07-29T13:45:00+01:00 · Codex · Milestone A.2 migrations and catalogue seed completed
- **Did:** Added embedded up/down migrations for products, images, settings, and delivery zones; generated a validated 27-product seed directly from `product-data.js`; integrated automatic migration before API startup.
- **State:** Go tests and the API binary build pass, seed generation is deterministic, PostgreSQL is clean at migration version 2, and a second startup preserved 27 unique products (17 bespoke, 10 RTW) with integer-kobo prices. API health remains HTTP 200.
- **Next:** Wait for approval, then begin A.3 product read endpoints only.
- **Blockers:** None for A.3; image keys, delivery fees, and client settings remain pending for later steps.

## 2026-07-29T13:22:06+01:00 · Codex · Core design tokens corrected before A.2
- **Did:** Replaced the serif display fallback with self-hosted Jost Variable, normalized the type scale through a 60px display token, and added the approved champagne, bright-gold, and taupe palette values.
- **State:** Nuxt typecheck/build pass; the live stylesheet serves the corrected tokens and local font assets. The function/computed-property walkthrough-comment rule is now durable in `AGENTS.md`.
- **Next:** Wait for approval, then begin A.2 migrations and catalogue seed only.
- **Blockers:** None for A.2; licensed Gilmer remains a pre-launch client decision.

## 2026-07-29T12:33:10+01:00 · Codex · Frontend styling migrated to modular SCSS
- **Did:** Replaced the plain CSS entry with SCSS design-system tokens, breakpoint mixins, globals, and a foundation-page partial; added Nuxt preprocessing and `sass-embedded`.
- **State:** No plain CSS source remains. Nuxt typecheck/build and live SCSS compilation pass; existing page output and service connectivity are preserved.
- **Next:** Wait for approval, then begin A.2 migrations and catalogue seed only.
- **Blockers:** None for A.2.

## 2026-07-29T12:20:25+01:00 · Codex · Milestone A.1 scaffold completed and verified
- **Did:** Initialized and renamed the repository to `oyecouture`; added PostgreSQL 16 Compose, complete `.env.example`, Go API `/healthz`, and a strictly typed Nuxt 3 client that consumes it.
- **State:** PostgreSQL, API, and web dev services are running. Go tests/build and Nuxt typecheck/build pass; runtime health and SSR both report all services connected.
- **Next:** Wait for approval, then begin A.2 migrations and catalogue seed only.
- **Blockers:** None for A.2. Later client credentials/content remain pending as listed in `HANDOFF-STATE.md`.

## 2026-07-29 · Claude · Handoff package + collaboration kit created
- **Did:** Wrote the numbered handoff docs `00`–`06`, `catalogue/` (product-data.js + CATALOGUE.md), copied `design-source/` (7 design files + HOW-TO-VIEW) and `reference-docs/` (PRD + QUOTATION). Added this collaboration kit: `AGENTS.md`, `CLAUDE.md`, `HANDOFF-STATE.md`, `WORKLOG.md`.
- **State:** Planning complete. No code yet. First 40% milestone paid; build authorized.
- **Next:** Milestone A, step A.1 — scaffold the repo (Nuxt 3 + Go + Postgres) and get both dev servers talking; copy the four baton files to the new repo root. See `HANDOFF-STATE.md`.
- **Blockers:** None for A.1. Client deliverables (Nomba creds, WhatsApp number, delivery fees, photos, care-page copy, DNS) gate later steps — chase them.
