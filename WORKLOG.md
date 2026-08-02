# WORKLOG — append-only (newest on top)

One short entry per work session or meaningful chunk. Format:
`## <ISO datetime> · <Codex|Claude> · <one-line summary>` then bullets: **Did / State / Next / Blockers**.

---

## 2026-08-02T14:35:40+01:00 · Codex · Ghost admin 1.2 — VBAN-style orders filters
- **Did:** Recovered Claude's in-progress ghost-admin orders work, confirmed the real repo is `/Users/chibuzoriwuagwu/Downloads/share/oyecouture`, inspected the VBAN reference filter, and committed `966392e` with order metrics, date-range API support, a single Filter button with status/type/from/to categories, applied filter chips, admin sidebar/metric cards, and shared select polish.
- **State:** `go test ./...`, `npx vue-tsc --noEmit -p tsconfig.json`, and `git diff --check` pass. `npm run typecheck`/`npm run build` still hit the known local Nuxt `oxc-walker`/`oxc-parser` parser-resolution issue after escalation, so they were not treated as source failures.
- **Next:** Build order detail/actions: authenticated detail endpoint/page, line items, contact/address, measurements, totals, copy `/pay/<token>`, status transitions, and one-tap WhatsApp reply. Keep `paid` server-verified only.
- **Blockers:** Nomba sandbox keys, Resend API key, real WhatsApp number, and netcup server/deploy remain unchanged.

## 2026-08-02T00:00:00+01:00 · Claude · Ghost admin 1.1 — authentication
- **Did:** Migration `000006` creating `admin_users` + `admin_sessions`; `internal/admin` (auth service, store, handler, `RequireSession`); `cmd/admin` CLI to seed or reset the owner account; routes `POST /admin/login`, `POST /admin/logout`, `GET /admin/me` wired in `server.go` with a startup goroutine sweeping expired sessions (`22065c9`).
- **State:** Verified live end to end — 401 without a cookie, identical error for unknown email and wrong password (with a dummy bcrypt compare so timing does not leak which accounts exist), HttpOnly cookie set, `/admin/me` 200 with it, logout 204 with the row deleted server-side, lockout after 8 failed attempts. Password is bcrypt cost 12 in the DB; **nothing secret in `.env`**. `go test ./...` passes.
- **Two bugs caught by that verification, both now regression-tested** in `internal/admin/handler_test.go`: the throttle keyed on `RemoteAddr` *including the ephemeral source port*, so every attempt landed in its own bucket and the limiter never fired; and it charged *successful* logins against the same budget, so the owner could lock themselves out simply by signing in and out. `clientAddress` now strips the port and honours the first hop of `X-Forwarded-For`; only failures are charged.
- **Next:** `GET /admin/orders` behind `RequireSession`, then the `/admin` UI on a `noindex` layout with no storefront link. Then order detail + measurements, `/pay/<token>`, status transitions, WhatsApp reply.
- **Blockers:** unchanged — Nomba Client ID, Resend key, server for deploy (netcup parked, ~13 Aug).

## 2026-07-30T19:05:00+01:00 · Claude · Orders, confirm + checkout (Milestone B steps 10–11)
- **Did:** `POST /orders` with **server-side pricing** — the storefront sends only slug/size/quantity and every price is re-read from the database (`33aaeae`). `/confirm` page: contact + optional measurements → saves the order → WhatsApp opens carrying the reference (`60861c6`). `/checkout` for ready-to-wear with **DB-driven delivery zones** and server-resolved fees (`bad7d67`), plus `GET /delivery-zones`.
- **State:** Verified against the live API and database. **A client cannot set a price** — sending `priceKobo` is rejected (`invalid_body`, unknown field), and an invented delivery zone is rejected too. Order kind is derived server-side: any unsized item makes the whole order `bespoke`/`pending_confirmation`; ready-to-wear is `awaiting_payment`. Confirmed end to end: `OYE-1008` (bespoke, WhatsApp message carries the order no) and `OYE-1011` (ready-to-wear, ₦78,000 + ₦5,000 Lagos = ₦83,000, address + zone persisted). **Every route now returns 200 — no 404s left.**
- **Next:** Payments (Nomba adapter + `/pay/<token>` page) and order emails — both blocked on credentials. Home still shows 27 products vs the design's curated 13 (user decision).
- **Blockers:** Nomba sandbox keys, Resend API key, real WhatsApp number. Server for deploy (netcup billing parked, deadline 13 Aug).

**⚠️ Nuxt trap worth knowing:** calling `useRuntimeConfig()` (or any Nuxt composable) **inside an async event handler** throws — the Nuxt instance is gone by then, and it surfaces as a vague "something failed". Resolve config during setup and close over it. That is why `useOrders()`/`useDeliveryZones()` capture the base URL at setup rather than per-call.

**⚠️ Delivery fees are PLACEHOLDERS** (Lagos ₦5,000 / Rest of Nigeria ₦10,000 / International ₦0) — the client has not supplied real pricing. They live in `delivery_zones`, so correcting them is a data change, not a deploy.

## 2026-07-30T14:10:00+01:00 · Claude · Design audit + About page
- **Did:** Audited every page against `design-source/` by diffing **computed styles**, not just text. Fixed: category pages ×4 restored to the design (light `#f4f0ec` band, breadcrumb, 40px→56px sentence-case H1, real copy, per-section chips; removed the invented collection heading + sort) `3678100`; grid corrected — **Home is 4-up, category pages are 3-up**, both 22px gap / 56px gutters `559737e`; **footer was fully inverted** — the design's is WHITE with `#4a443e` links, not dark `8e97bf4`; on-card size selector for ready-to-wear, which also lets Add-to-Bag work from the grid `adf90af`; **built the whole About page** `cef6593`.
- **State:** All values verified against the design at 1440px. About page: hero 62px/700, champagne eyebrow `0.42em`, tab pills 39px, gold section labels, 40px/600 headings, 16px/300 `#4a443e` body, stat 34px, mission `0.22em`, principle `0.1em` — every one matches. The band takes per-page overrides so About shows "Step into the world of Oyecouture" + its extra line, with no duplicate band. All four footer About links now resolve.
- **Next:** Remaining audit gaps — Home shows 27 products vs the design's curated 13; sort labels use `:` not `·`; `/checkout` still 404s (step 11); step 10 (persist orders + measurement form) still open.
- **Blockers:** Nomba/Resend credentials; netcup server (billing parked, deadline 13 Aug).

**⚠️ Two measurement lessons, both of which caused real bugs:**
1. **Widen the viewport before measuring the design.** The preview defaults to ~537px, so every reading is silently a *mobile* value. That is how the category H1 got recorded as 40px when it is 56px at desktop.
2. **The design stores text sentence-case and uppercases it in CSS.** Selectors matching on `'OUR STORY'` find nothing; match case-insensitively against `'Our Story'`.

## 2026-07-30T11:40:00+01:00 · Claude · Cart + bespoke WhatsApp confirmation (Milestone B steps 8–9)
- **Did:** Guest cart shipped (`68b168a`). `useCart` (useState, SSR-safe) + `plugins/cart.client.ts` for localStorage with a per-field validator; `CartDrawer` (430px, matches the approved design) + `CartCheckoutAction`; nav bag badge; Add-to-Bag wired on grid and detail. **The RTW/bespoke split is enforced in the cart:** any bespoke line routes the whole bag to a prefilled `wa.me` deep link; RTW-only routes to `/checkout`. Line key is `slug::size`, so the same shirt in M and L are two lines but the same shirt twice merges to qty 2.
- **State:** 14 cases verified live — persistence across reload, merge-vs-split by size, mixed-cart routing, RTW-only routing, quantity clamp (max 10, + disables), corrupt-JSON recovery, invalid-line rejection **with storage self-cleaning on load**, empty state, Escape close, body-scroll lock/restore, and a sized product from the grid routing to the detail page instead of guessing a size. SCSS compiles, `go build`/`go test` pass, zero console errors. Found and fixed two real bugs mid-test: a missing `WHATSAPP_NUMBER` produced a broken `wa.me/?text=` link (CTA now disables when unconfigured), and sanitised data was not written back to storage.
- **Next:** `/checkout` route does not exist yet — the RTW CTA currently points at a 404. That is step 11. Step 10 (persisting the order server-side + measurement form) also still open.
- **Blockers:** None for step 10/11. Nomba + Resend credentials still gate payments and emails.

## 2026-07-30T09:05:00+01:00 · Claude · Fixed nav dropdown hover/click conflict
- **Did:** Fixed the reported wart where hovering opened the menu and the next click closed it (with the pointer still inside, no `mouseenter` could reopen it). `useStorefrontNavigation` now tracks `pointerItemId`; `toggleDropdown` will not close a menu the pointer is holding open, while keyboard/touch activation still toggles. Committed `7848780`.
- **State:** Verified live through all five paths: hover opens → click keeps it open → pointer-leave closes → click with no pointer opens → click again closes. SCSS compiles; no console errors.
- **Next:** Milestone B step 8 (cart state) still awaits user authorization. M2 transactional motion ships with the cart.
- **Blockers:** None. Recorded a verification lesson in `HANDOFF-STATE.md`: dropdown state must be asserted via computed `display` on the specific panel id, not `querySelector` presence (`v-show` + two panels + async Vue flush all mislead).

## 2026-07-30T00:20:00+01:00 · Claude · M1 motion foundation shipped
- **Did:** Implemented `07-MOTION-AND-MICROINTERACTIONS.md` §9 M1 (`d2544dc`): new `design-system/_motion.scss` (duration/easing/displacement/stagger/press tokens) wired into `main.scss`; global `prefers-reduced-motion` block in `_globals.scss` that kills travel but preserves state feedback; site-wide `:active` press feedback plus smooth colour/border transitions on buttons and links; retuned all three pre-existing transitions (nav chevron, mobile menu, category tile) to tokens with exits now faster than entrances; **fixed the hero marquee WCAG 2.2.2 violation**.
- **State:** Verified live in-browser — every motion token resolves, marquee `animationPlayState` goes `running → paused → running` on hover/unhover, and all four guard rules (hover-pause, touch fallback, global reduced-motion, tap feedback) confirmed present in the cascade. SCSS compiles clean; home + product SSR return 200; zero console errors; no horizontal overflow at 375px. **Zero raw `ms`/`cubic-bezier()` remain outside the token file.**
- **Next:** M2 transactional motion ships *with* the cart in Milestone B (§4.4–4.5: cart drawer choreography, the 4-state action button, cart count, wishlist, S3 bespoke WhatsApp confirm) — build it in, don't retrofit. M3 signature polish (hero reveal, scroll reveals, View Transitions morph, image blur-up) before launch.
- **Blockers:** None. Two notes recorded in `HANDOFF-STATE.md` Decisions: the marquee's per-input-type behaviour is a deliberate design deviation, and Claude's local `build`/`typecheck` env failure is pre-existing (proven against a stashed clean tree), so Claude verifies via sass compile + browser.

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
