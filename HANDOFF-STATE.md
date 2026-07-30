# HANDOFF-STATE.md — the baton

**Read this first. Keep it true at all times.** This is the single source of truth for where we are and what to do next. Update it continuously (a usage limit may cut you off without warning).

---

## Current status
- **Phase:** 1 (MVP). First 40% milestone **paid**; build authorized.
- **Milestone:** A — Foundations & catalogue (see `05-BUILD-ROADMAP.md`).
- **Step:** A.6 complete: the approved instant product-search overlay is implemented and verified. Milestone A is complete.
- **Last worked by:** Codex · 2026-07-30.
- **Repo / branch:** `/Users/chibuzoriwuagwu/Downloads/share/oyecouture` · `main`.
- **Last implementation commits:** `562faf8` (reference-inspired button animations), `d743037` (instant search overlay), `803149b` (single-column mobile catalogue), `5bcfaa8` (storefront microinteractions), `9cc1333` (approved-design alignment), `8736afc` (product detail), `ed4f3d9` (shop/category pages).

## ▶ Next action (the ONE thing to do next)
**Wait for user review and explicit authorization before starting Milestone B step 8.** The next roadmap step is cart state only. Do not begin checkout or payments until the cart runs and the user authorizes the next step.

### ✅ Nav dropdown hover+click fixed by Claude (`7848780`)
The dropdown now opens on **hover and click**. Bug fixed: hover opened it, then clicking the same trigger closed it — and because the pointer was still inside, no further `mouseenter` could reopen it. `useStorefrontNavigation` now tracks `pointerItemId`; `toggleDropdown` refuses to close a menu the pointer is currently holding open, while keyboard and touch activation (no pointer) still toggle normally. Verified live: hover opens → click keeps open → pointer-leave closes → click-without-pointer opens → click again closes.
**Verification lesson for both agents:** asserting dropdown state via `document.querySelector('.desktop-navigation-dropdown')` is wrong three ways — the panel uses `v-show` (always in the DOM), there are **two** panels (Casuals + Collection) so `querySelector` returns the wrong one, and Vue flushes DOM updates asynchronously so a synchronous read after `.click()` is stale. Assert on `getComputedStyle(document.getElementById('collection-navigation-menu')).display`, and read in a separate round-trip after a click.

### ✅ M1 motion foundation — DONE by Claude (`d2544dc`), claim released
Shipped: `design-system/_motion.scss` tokens, the global `prefers-reduced-motion` block, site-wide `:active` tap feedback, the three existing transitions retuned to tokens, and the **marquee WCAG 2.2.2 fix**. Codex is clear to work anywhere again. **From now on, all timing/easing must come from `_motion.scss` tokens — no raw `ms` or `cubic-bezier()` in component files** (currently zero remain outside the token file; keep it that way). M2 transactional motion (cart drawer, 4-state action button, wishlist, S3 WhatsApp confirm) ships *with* the cart in Milestone B — build it from `07-MOTION-AND-MICROINTERACTIONS.md` §4.4–4.5 rather than retrofitting.

## Done & verified
- Handoff docs `00`–`06`, `catalogue/`, `design-source/`, `reference-docs/` are complete and current.
- Git repository initialized on `main`; commits authored as `mikechibuzor` with no co-author trailers.
- Repository and generated service resources are consistently named `oyecouture`.
- Root `.env.example` matches `06-SETUP-AND-ENV.md`; Compose project is explicitly named `oyecouture` and runs PostgreSQL 16.
- Nuxt `3.21.10` scaffolded in `web/` with strict TypeScript, Node `22.23.1`, shared types/constants, and a typed API-health composable.
- Frontend styling uses modular SCSS with design-system tokens, reusable breakpoint mixins, global rules, page partials, and the `sass-embedded` compiler.
- Approved typography and color foundations are aligned: the extracted Gilmer demo fonts map all four design weights, the tracking scale matches the source, and champagne/gold/taupe tokens match the design spec. A proper Gilmer webfont licence remains required before launch.
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
- A.6 product search is complete: a full-screen on-brand overlay uses typed shared copy/config/contracts, six approved top terms, six editorial recommendations, and race-safe debounced `GET /products?q=` requests without adding a frontend dependency.
- A.6 accessibility is verified: the input autofocuses, body scrolling locks, Tab stays in the dialog, Escape closes it, desktop focus returns to Search, and mobile focus returns to the visible menu trigger after its menu closes.
- A.6 responsive browser verification passes at 1440px desktop and 390px mobile: four-column desktop and one-column mobile results, three live `linen` matches, accurate singular/plural headings, empty and clear states, `/product/linen` navigation, and no console warnings or errors.
- Mobile catalogue cards now render one full-width column below 640px while retaining the existing three-column tablet and four-column desktop grids.
- A.6 verification passes on Node `22.23.1`: `npm run typecheck`, production `npm run build`, fresh dev-server `vue-tsc` with zero errors, API-backed browser interactions, and responsive visual inspection.
- Storefront interaction refinement is complete: Add to Bag uses a shared typed fizzy-particle treatment, the home hero CTA uses the requested two-circle collision fill, and the mobile menu uses a reusable SVG line-to-X morph with a traced ring. All three use existing motion tokens, composited properties, reduced-motion fallbacks, and no new dependency.
- The interaction refinement passes strict Nuxt typecheck, production build, HTTP `200` checks for Home and `/product/linen`, and in-app mobile verification of the menu open/close state and accessible label. Add to Bag remains intentionally disabled until Milestone B step 8; its hover treatment is available for visual review without adding cart behavior early.
- A.5 design-alignment pass now matches the approved Shop source: exact 48px/700 hero typography and line break, 300-weight support copy, duplicated five-image staggered rail, 40px/600 collection heading, source-width content, champagne category accents, horizontal arrows, 3:4 cards, black/champagne tags, 26px translucent carousel controls, source-weight card text, and dashed image placeholders. The user-approved `Explore Oyecouture / Shop By Category` heading and Add to Bag icon remain.
- The design-alignment pass was verified at 1440px desktop and 390px mobile, including carousel interaction, placeholder state, no horizontal overflow, strict Nuxt typecheck, production build, fresh dev-server startup with zero `vue-tsc` errors, and HTTP 200 SSR responses for Home and Suits.
- **Claude review sign-off (2026-07-29):** A.1–A.3 independently reviewed and re-verified (`go vet`/`build`/`test` all pass). Approved with no changes. Claude's three earlier design-token findings (display font, champagne `#e0cf9c`, type scale) were already resolved by Codex in `576ee77`; no fix commit needed. Golden invariants confirmed (integer `priceKobo`, hidden `R2Key`, validated inputs, active-only reads, RTW/bespoke split).
- **Claude review sign-off — A.4 (2026-07-29):** Storefront shell reviewed and approved. Design fidelity verified against `design-source/` (announcement copy exact; nav Home/Trad/Suits/Casuals▾/Collection▾; Casuals→Men/Women; the COLLECTION names "Dear Groom / Wei-Xin / OYE — Man of Legacy" confirmed present in all 7 design files — extracted, not invented; footer Shop/About/Customer Care + IG/X/TikTok/WhatsApp + legal links match). Code quality strong: nav state in a composable, full a11y (aria-controls/expanded, Escape, click-outside, close-on-route-change), lucide tree-shaken icons, search/account/wishlist/cart correctly emit events for future features. One required A.5 follow-up flagged: remove `:external` per link as pages land (see Decisions). Note: Claude's local `npm run typecheck` was inconclusive due to a missing `oxc-parser` native dep corrupting a generated `.nuxt` file (env artifact, not a source defect); Codex's own typecheck/build/`go test` pass.

## In progress (exact state)
- Local PostgreSQL, the product-enabled Go API, and the Node `22.23.1` Nuxt dev server are running. Milestone A is complete; wait for user review and Milestone B step 8 authorization.

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
- Nomba credentials (sandbox first), WhatsApp business number, delivery-fee amounts per zone, product photography, customer-care page copy, domain DNS access. None block cart state, but they gate final content and later checkout/payment/deployment steps.

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
- Search recommendations preserve the approved source order (`agbada`, `monarch`, `senator`, `linen`, `classic`, `wrap`); live queries debounce for 250ms, cancel superseded requests, and leave all matching logic on the validated server endpoint.
- Search results use one column on mobile, two on tablet, and the approved four-column desktop grid. This follows the project's mobile readability decision while retaining the source's desktop dimensions.
- Approved design-source WebPs are versioned local fallbacks while `product_images` is empty. API/R2 media always takes precedence, so client photography can replace the fallbacks without changing page components.
- Nuxt's unused experimental app manifest is disabled to avoid a dev-only Vite alias error. Official Lucide icon data is rendered through a typed local Vue component to keep SSR hydration deterministic on Vue `3.5.40`.
- The requested DEV/CodePen animation references are reproduced as local CSS/Vue mechanics rather than copied wholesale: the collision fill and SVG ring/menu morph preserve their identifying movement, while the Add to Bag fizz is reduced to eight deterministic particles to protect performance and the quiet-luxury motion direction.

- **Typography (2026-07-29, Claude):** The design's real face is **Gilmer** (commercial, Fontfabric). Per user decision, the four Gilmer weights were extracted from the design bundle and self-hosted as WOFF2 in `web/public/fonts/` (`_fonts.scss`), replacing the Jost stand-in. A `--tracking-*` letter-spacing scale (values taken from the design) now drives all uppercase labels/headings/buttons. ⚠️ **Gilmer needs a proper webfont licence before production launch** — self-hosting it now is for build/demo fidelity; do not launch commercially on the unlicensed font without confirming licensing with the client. The `@fontsource-variable/jost` dependency is now unused and can be pruned. The visual side-by-side against `design-source/Oyecouture Shop.dc.html` was completed in `9cc1333`; desktop/mobile typography and catalogue details were corrected to the source values.

- **Motion system (2026-07-29, Claude):** `07-MOTION-AND-MICROINTERACTIONS.md` is now the authority for anything that moves — **read it before adding any animation.** Key rules: use `_motion.scss` tokens (no raw ms/cubic-bezier), animate only composited properties (`transform`/`opacity`/`filter`/`clip-path`), exits faster than entrances, **no spring/bounce easing**, tap `:active` feedback on everything (mobile has no hover), and nothing on the cart→checkout→payment path exceeds `--duration-slow`. Stack is **CSS-first**: CSS transitions → Vue `<Transition>` → CSS scroll-driven animations (`@supports` guarded) → View Transitions API via Nuxt `app.viewTransition: true`. **No motion library is installed — do not add one without recording the justification here.** ⚠️ **WCAG 2.2.2 launch blocker:** the infinite hero marquee needs a pause/stop mechanism (§6.2).

- **M1 motion implementation notes (2026-07-29, Claude, `d2544dc`):** (a) Press feedback uses the independent `scale` property, **not** `transform: scale()` — this composes with elements that already carry a transform (e.g. `.product-carousel-button`'s `translateY(-50%)`), which `transform: scale()` would have silently clobbered. Same reason `.category-tile img` hover now uses `scale: var(--zoom-image)`. (b) **Design deviation — hero marquee:** on `(hover: none)` devices the auto-scroll is dropped and the rail becomes a swipeable overflow strip; on `(hover: hover)` it animates and pauses on `:hover`/`:focus-within`. Rationale: touch devices cannot hover to pause, so hover-pause alone would leave the majority-mobile audience non-compliant with WCAG 2.2.2; this also saves battery on mid-range Android. Visually the rail is unchanged — only its motion behaviour differs by input type. Flag to the user if they want the auto-scroll back on touch (it would then need a visible pause control).
- **Local env caveat (Claude's machine only):** `npm run build` and `npm run typecheck` fail here with an `oxc-parser`/`@nuxt/vite-builder` resolution error. **Verified pre-existing** — reproduced identically with all changes stashed on a clean tree. Not a code defect; Codex's environment builds fine. Claude verifies frontend work via `sass-embedded` compilation + live browser checks instead.

## Reminders (non-negotiable)
Money in kobo (int64, no floats) · `paid` only after server-side verify · guest checkout only in phase 1 · mobile-first · ghost admin stays unlinked/unmentioned to client · commits as `mikechibuzor`, no `Co-Authored-By`.
