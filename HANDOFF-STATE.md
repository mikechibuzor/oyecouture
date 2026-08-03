# HANDOFF-STATE.md — the baton

**Read this first. Keep it true at all times.** This is the single source of truth for where we are and what to do next. Update it continuously (a usage limit may cut you off without warning).

---

## Current status
- **Phase:** 1 (MVP). First 40% milestone **paid**; build authorized.
- **Milestone:** Storefront complete · ghost admin phase 1 complete · **staging is deployed and live** · the client's real catalogue is in.
- **Last worked by:** Claude · 2026-08-03.
- **Repo:** `/Users/chibuzoriwuagwu/Downloads/share/oyecouture` · `main` · **now has a GitHub remote: private repo `mikechibuzor/oyecouture`. Push to `origin`.**
- **Staging:** **https://82-29-184-79.sslip.io** — real Let's Encrypt cert, all five containers healthy.
- **Server access:** `ssh oyecouture` (config on the user's Mac, key `~/.ssh/oyecouture`). Deploy with `cd ~/oyecouture && git pull && ./deploy/deploy.sh`.
- **Admin:** https://82-29-184-79.sslip.io/admin — the user seeded their own account on 3 Aug; **neither agent knows the password, and that is correct.** Re-seed or reset with `docker compose -f compose.prod.yaml --env-file .env.production exec -it api /admin`.

## ▶ NEXT ACTION — run one order end to end (nobody has done this yet)

Every component is now in place for the first time: a public HTTPS webhook URL, a verified sending domain, a live email key, and real products. **The full chain has never been exercised once:**

```
place order → Paystack sandbox card → webhook fires → order flips to paid → confirmation + receipt email arrives
```

Each part is tested in isolation; the joins between them are not. This is the highest-risk untested path in the project and the thing a client demo depends on. **Do this before starting phase 2.** Use the user's own email address — see the credentials note below, emails are real now.

After that, in order: **(1)** client review of the 22 placeholder prices and the product groupings (casuals 7/8/9 and 23–26 are the least certain), **(2)** ghost admin **phase 2** — product create/edit with Cloudinary image upload, then editable home content, **(3)** chase the client for **suit photography** (suits are ₦275k–₦850k and still show a mockup mannequin).

⏰ **Two clocks:** netcup withdrawal (€68.64, closes ~13 Aug — the user said they would send it) and the **Hostinger renewal at $19.49/mo (~₦26,900), which exceeds the whole ₦25k maintenance retainer** — move to the 12-month term (~₦9,646/mo) or cancel, ~28 days from 2 Aug.

⚠️ **Staging is sandbox for money but LIVE for email.** Paystack is `sk_test_`; SendByte is `sk_live_`. Any order placed on staging sends a genuine email to the address entered. There are already **5 test orders** in the staging database that must be purged before go-live.

## 🔑 CREDENTIALS TO REGENERATE (Claude, 3 Aug)

Both work; both were pasted into a chat transcript and should be rotated once setup settles:
1. **SendByte production key** (`sk_live_…`) — now live on the server in `.env.production`. The local `.env` deliberately keeps the **sandbox** key so local runs cannot send real mail.
2. **Cloudinary API secret** — in the local gitignored `.env`. The key also had to be moved from *media user* to *admin* role or uploads 403.
3. (Older, still outstanding) **Nomba private key** — pasted in chat during sandbox setup.

⚠️ **Order emails are now REAL.** Any order placed on staging sends a genuine confirmation and receipt to the address entered. Use the user's own address for test orders.

**▶ Never run end to end yet:** order → Paystack sandbox card → webhook → status `paid` → SendByte email landing in a real inbox. Every part is now in place (public HTTPS webhook URL, verified sending domain, live key) but the whole chain has not been exercised once. Do this before any client demo.

## 🖼 PRODUCT PHOTOGRAPHY + CLOUDINARY (Claude, 3 Aug)

**Cloudinary is the image backend** (user's decision; it was also the original PRD plan). Credentials in the gitignored `.env` as `CLOUDINARY_CLOUD_NAME=exgithrw`, `CLOUDINARY_API_KEY`, `CLOUDINARY_API_SECRET`. ⚠️ **The API secret was pasted into chat — regenerate it in the Cloudinary console once setup is done.** The key had to be moved from *media user* to *admin* role; a media-user key returns `403 missing permissions (actions=["create"])` on upload.

**Serving needs no code change.** `buildPublicImageURL` is `base + "/" + key`, so pointing `R2_PUBLIC_BASE` at `https://res.cloudinary.com/exgithrw/image/upload/f_auto,q_auto` and storing the Cloudinary public_id in `product_images.r2_key` just works. (`r2_key`/`R2_PUBLIC_BASE` are now misnomers — they hold a Cloudinary public id and delivery base.)

**92 images uploaded** from the client's two Drive folders, resized 2400px/q88 first: **709 MB → 102 MB**. Public ids are `oyecouture/casual-NN`, `oyecouture/asoke-look-NN-N`, `oyecouture/campaign-N`. Uses **0.4% of the 25 free monthly credits**.

**Measured delivery sizes (do not re-estimate — these are real):** casuals ~99 KB at w_1200; aso-oke **125 KB at w_600, 254 KB at w_900, 408 KB at w_1200** — textile detail compresses poorly. Realistic ceiling ≈ **16,000 page views/month**, not the 29–40k first projected. Keep grids at w_400–600 and consider `q_auto:eco` for thumbnails.

**The client's real inventory, grouped and priced (dummy prices, awaiting client correction):**
- **Aso-oke → `trad`, aso-oke filter, 12 looks:** full agbada (L1, L2, L4, L5) ₦780,000 · aso-oke two-piece (L8, L9, L11, L12) ₦520,000 · embroidered buba+sokoto (L3, L6) ₦295,000. **L7 is a couple shot** (wedding ensemble, ₦860,000?) and **L10 shows two models** (L9 + L1 together — probably lifestyle, not a product). Both need the user's call.
- **Casuals → 11 products from 33 photos:** long-sleeve two-piece ₦115,000 · short-sleeve two-piece ₦98,000 · women's two-piece ₦105,000.
- **7 campaign images** (`Group *`, `October *`) are multi-model shots — lookbook/banners, not products.

**▶ DECIDED: the 27 existing products are design-mockup placeholders, not real inventory.** `trad` and `casuals-*` are to be **replaced** by the real photographed products; `suits` stays on placeholders until suit photography arrives.

## 🚀 STAGING SERVER IS LIVE (Claude, 2 Aug) — Hostinger, netcup abandoned

**Hosting moved from netcup to Hostinger.** ⏰ **netcup order `54429992` must still be formally withdrawn — €68.64 unpaid, window closes ~13 Aug.** Nothing depends on it any more; not withdrawing just buys a 12-month liability.

- **Server:** Hostinger KVM 1 · `82.29.184.79` · Ubuntu 24.04.4 LTS · 3.9 GB RAM · 1 vCPU · 48 GB NVMe.
- **Billing:** 1-month term, $10.74 paid. ⚠️ **Renews at $19.49/mo (~₦26,900) — more than the entire ₦25k maintenance retainer.** Switch to the 12-month term (~₦9,646/mo) or cancel before it rolls. Set a reminder ~28 days out.
- **Access:** `ssh oyecouture` (config written on the user's Mac). Dedicated key `~/.ssh/oyecouture`, **no passphrase** — created because the only existing key was labelled `chibuzor@ovalfi.com`, i.e. the work key, and client production must not share a credential with the employer.
- **Hardened:** `deploy` user with passwordless sudo + docker group; 2 GB swap (4 GB RAM on 1 vCPU, and the Nuxt build runs `vue-tsc`); `ufw` 22/80/443 only; **key-only SSH — password auth verified refused**. Deploy-user login was verified *before* passwords were disabled, deliberately.
- **Code:** now on GitHub — **private** repo `mikechibuzor/oyecouture`, 90 commits. Previously the repo existed only on the user's laptop. Server pulls via a **read-only deploy key** (`~/.ssh/github_deploy` on the VPS).
- **Config:** `~/oyecouture/.env.production`, mode 600, sandbox Paystack + SendByte keys, freshly generated Postgres password, `WHATSAPP_NUMBER=2348094939264` (the user's own number, for testing), `SITE_HOST=82-29-184-79.sslip.io`.

**Three bugs the first real deploy exposed — none of which local builds or tests could have caught:**
1. **`456bbdb`** — `api/Dockerfile` pinned `golang:1.23-alpine` while `go.mod` requires `go >= 1.25.0`. `GOTOOLCHAIN` is `local` in that image, so it failed instead of fetching a toolchain. Keep the Dockerfile base and the `go` directive in step.
2. **`02cbc66`** — the runbook told you to generate `POSTGRES_PASSWORD` with `openssl rand -base64 32`, but compose interpolates it into `DATABASE_URL` and base64 emits `/`, which terminates the authority section. The API crash-looped on `invalid port after host` **only after a full image build**. Generator is now `openssl rand -hex 32`, and `deploy.sh` rejects `/ @ : ? #` up front.
3. **`7660408`** — the ghost admin had **no account and no way to create one** on a deployed box: the seeding CLI needs Go, the distroless runtime image has none, and the runbook had no seeding step. The `admin` binary now ships in the image — `docker compose exec -it api /admin`.

**Live staging URL: https://82-29-184-79.sslip.io** — real Let's Encrypt cert, HTTP→HTTPS 308. Verified from the public internet: `/`, `/suits`, `/admin`, `/api/healthz`, `/api/products`, `/api/delivery-zones` all 200; **Postgres (5432) and the API (8080) are not reachable from outside** — only Caddy is. Fresh DB at migration v9: 27 products, 3 delivery zones, 0 orders.

**▶ STAGING → PRODUCTION is a config swap, not a rebuild.** When going live: point the client's domain A record at `82.29.184.79`, set `SITE_HOST` to it, swap in **live** Paystack keys, re-register the webhook in the **live** Paystack dashboard, then `./deploy/deploy.sh`. Caddy issues the new certificate itself.
⚠️ **Purge test data before go-live.** There are ~40 orders, 10 marked paid, ₦1,543,000 of fake revenue. Truncate `orders`, `order_items`, `measurements`, `payment_attempts`, `outbound_emails`; **keep** `products`, `product_images`, `delivery_zones`, `admin_users`.

## 🔍 PAYMENTS REVIEW (Claude, 2 Aug) — passed, with two items to decide

Reviewed `internal/payments` end to end against the running stack. **The money path is sound.** Verified, not assumed:
- **Amount and currency are checked server-side** against *both* the stored attempt and the order total (`repository.go`), so neither a tampered attempt nor a repriced order can confirm.
- **Row-locked** (`FOR UPDATE OF pa, o`) and idempotent — an already-`paid` order commits without reprocessing.
- **Cross-order references are blocked**: the callback binds the reference to the order behind the `payToken`.
- **The webhook never trusts its own payload for money** — it takes only the reference and re-verifies against Paystack's API. Signature is HMAC-SHA512 over the **raw body**, compared with `hmac.Equal`, checked *before* parsing. Forged and unsigned webhooks both 401 (verified live).
- **`paid` cannot be set by hand.** `PATCH /admin/orders/{no}/status` rejects `paid` and illegal transitions with 400, and 401s unauthenticated — enforced in the API, not just hidden in the UI (verified live).
- Secret key is never logged, `.env` is gitignored and was never committed, `.env.example` holds empty placeholders.
- Public `GET /payments/{payToken}` excludes email, phone, pay token and order id via `json:"-"`.

**FIXED (`d29c830`):** the amount/currency guard had **no test** — the service tests stub the repository, so the single most important money check was never exercised. Extracted it as `verificationMatchesOrder` and covered underpayment, overpayment, zero, tampered attempt, repriced order and three currency cases. Confirmed real by removing the guard: six cases fail without it.

**✅ ITEM 1 FIXED (`6685057` + `beeef97`) — webhook retry storm.** Terminal failures (`ErrReferenceNotFound`, `ErrReferenceMismatch`, `ErrPaymentAmountMismatch`, `ErrPaymentCurrency`, `ErrOrderNotPayable`, `ErrOrderNotFound`) are now acknowledged with **200** and logged, because redelivery can never resolve them. **Anything unrecognised still returns 500 deliberately** — a DB blip or an unreachable Paystack must be retried, and wrongly acknowledging one would drop a real payment. Covered by `webhook_retry_test.go` (14 cases + signature gate); mutation-checked — breaking the classification fails 6 subtests.
**The first attempt (`6685057`) did not actually fix the case that fires.** Replaying a signed webhook against staging showed it still returning 500: when Paystack itself answers **400** for a reference, the adapter flattened every non-2xx into one opaque string, so it read as transient. `beeef97` keeps the status code (`APIError.Retryable()`) and gives `status:false` its own sentinel (`ErrVerifyRejected`) — **4xx and explicit rejections are acknowledged, 5xx still retries.** Verified live: signed+unknown → **200**, unsigned → 401, bad signature → 401. Lesson: unit tests passed on the first attempt; only replaying against the real provider found the gap.

**▶ ONE ITEM STILL FOR THE USER TO DECIDE (not a bug, do not silently "fix"):**
2. **`payment_attempts.provider_payload` stores Paystack's entire verify response.** That includes the customer's `email`, `phone`, `ip_address`, and an `authorization` object with `bin`, `last4`, `exp_month`, `exp_year`, `bank`, `account_name`, `receiver_bank_account_number` and a reusable `authorization_code`. No full PAN, so not PCI cardholder data, but it widens the blast radius of a database leak and is customer PII under Nigerian data-protection rules with no retention policy. Suggested: persist a trimmed payload (status, amount, currency, reference, paid_at, gateway_response, channel, last4/card_type for support) instead of the whole blob.

## Ghost admin phase 1 — history and detail (superseded as the next action; see NEXT ACTION at the top)
**Remember what the ghost admin is:** a full admin layer built in phase 1 but **not billed and not revealed** to the client, who declined it. Pre-priced at ₦300k as an add-on for when they ask. It must be genuinely usable, not a stub — that is the whole point.

**✅ 1.1 AUTH — DONE (`22065c9`).** Session-cookie auth over `admin_users` + `admin_sessions` (migration `000006`). Seed or reset the owner account with:
```
cd api && go run ./cmd/admin
```
Nothing secret lives in `.env`; the password is bcrypt cost 12 in the DB. `POST /admin/login`, `POST /admin/logout`, `GET /admin/me` (behind `RequireSession`). 12h sessions, HttpOnly cookie `oyecouture_admin_session`, `Secure` only when `APP_ENV=production`. Expired sessions swept by a startup goroutine.
Verified live end to end: 401 without a cookie · identical error for unknown email and wrong password · cookie set HttpOnly · `/admin/me` 200 with it · logout 204 and the row deleted server-side · lockout after 8 failed attempts.
**Two bugs found and fixed during that verification, worth not reintroducing:** (a) the limiter keyed on `RemoteAddr` *including the ephemeral port*, so every attempt got its own bucket and it never fired — `clientAddress` now strips the port and takes the first hop of `X-Forwarded-For`; (b) it charged *successful* logins against the same budget, so the owner could lock themselves out just by signing in and out — only failures count now. Both are covered by `internal/admin/handler_test.go`.

**✅ 1.2 ORDERS LIST — DONE (`93c564c`, `fddbe22`, `966392e`, `a8a372c`, `9619edc`, `ac24935`, `2be9f3c`).** `GET /admin/orders` is behind `RequireSession`; `/admin` uses the noindex admin layout and remains unlinked from the storefront. The order desk now has a fixed admin sidebar, summary metric cards, paginated table, search, and a single VBAN-style **Filter** button with category submenus for status, type, from date, and to date. Date filters are passed as `YYYY-MM-DD`; the API treats `to` as inclusive through end-of-day, drops invalid dates, and keeps a backwards range usable by dropping the upper bound. `GET /admin/orders/metrics` powers the summary strip; paid revenue counts only paid/in-progress/shipped/delivered orders. Shared `UiSelect` replaced native selects for the admin metric period and storefront/checkout select polish. Pagination is now numbered (`1 … current±1 … last`), the filter group sits left of the table, search is constrained to half width on desktop, the filter popover clamps to the viewport, metric cards are slightly taller, table columns align with headers, rows open their order detail page, and the admin surfaces have quieter radius/hover/press microinteractions using existing motion tokens.

**✅ 1.3 ORDER DETAIL + ACTIONS — DONE (`d701dfd`).** Added authenticated `GET /admin/orders/{orderNo}` and `PATCH /admin/orders/{orderNo}/status`, plus `/admin/orders/[orderNo]`. Detail shows line items, contact/address, delivery, measurements, totals, payment link copy, WhatsApp reply, and safe status actions. Admin actions deliberately cannot set `paid`; `paid` remains reserved for server-side payment verification/webhook. Allowed manual flow is confirmation/awaiting-payment cancellation plus fulfilment transitions after payment (`paid → in_progress → shipped → delivered`, with cancellation/refund from paid where useful).

**▶ 2.1 PRODUCTS + IMAGES — next up.** Build product create/edit/read admin APIs and UI, including active state, section/filter fields, price in kobo, has-sizes behavior, and image management. Important: `product_images` is currently DB-backed but the storefront still has local fallback maps; the ghost-admin story needs images to flow DB/R2 → API → frontend without deploys. Avoid real R2 uploads until credentials are supplied; design the interface and persistence boundary so it can accept uploaded image metadata later.

**Then 2.2 editable home/content settings:** hero headline both lines, sub-copy, CTA label, sliding gallery upload/reorder, announcement bar, category tiles, delivery fees, WhatsApp number.

## Earlier next-action — demo sprint, split between agents (30 Jul)
📅 **NO FIXED DEMO DATE (updated 2 Aug).** The client has told the user to demo **when the build is ready** — 31 Jul and 3 Aug are both dead as deadlines. Do not optimise for speed or ship thin work to hit a date; build it properly and demo when it is genuinely worth showing. The 40% / ₦180,000 milestone still triggers on demo delivery, so don't drift either. **The only hard date left is the netcup withdrawal window, ~13 Aug** (see PARKED item). Roadmap **step 7 (deploy staging to a live link) was skipped** — we jumped from A.6 to step 8. Two parallel tracks, claimed to avoid collision:

**✅ CLAUDE — deploy artifacts DONE (`d965adb`), server pending.** Production stack committed: `api/Dockerfile` (distroless), `web/Dockerfile`, `compose.prod.yaml` (caddy + web + api + postgres + nightly pg_dump), `deploy/Caddyfile` (auto-TLS, `/api/*` → Go, `/*` → Nuxt), `deploy/deploy.sh`, `deploy/README.md` (full server runbook), `.env.production.example`. Postgres publishes no ports. **Nuxt now resolves the API differently for SSR vs browser** — `useApiBaseUrl()` returns `apiInternalBase` (`http://api:8080`) on the server and `public.apiBase` (`/api`) in the browser; all three API composables use it. Same-origin in production means **no CORS**. Blocked only on the server existing: netcup order `54429992` is being re-done on monthly billing (the 12-month term was billed €68.64 upfront; 14-day withdrawal window closes ~13 Aug).

**🔒 CLAUDE owns (do not touch):** (a) **DB-driven product images** — `product_images` is empty and every image is currently a hardcoded frontend map (`PRODUCT_FALLBACK_IMAGES` in `constants/products.ts`), so no image can change without a deploy, which breaks the ghost-admin story. Seeding rows + `R2_PUBLIC_BASE` so images flow DB→API→frontend. (b) **Deploy artifacts + staging on a live link** (roadmap step 7): Dockerfiles, prod compose, Caddy auto-HTTPS, deploy script.

**✅ CLAUDE — cart DONE (`68b168a`), steps 8–9 complete. Codex: please review.** Guest cart with localStorage persistence + validator, 430px drawer matching the design, nav badge, and the RTW/bespoke split enforced (any bespoke line → prefilled `wa.me` deep link; RTW-only → `/checkout`). 14 edge cases verified live. **⚠️ Known gap: `/checkout` does not exist yet, so the RTW CTA points at a 404 — that is step 11.** Step 10 (persist the order server-side + measurement form) also still open. Claude continues on 10/11 unless reassigned.

**🔒 (previously Codex, now reassigned — left for history):** **Milestone B steps 8–10** — cart state, the cart drawer per `02-DESIGN-SPEC.md`, and the bespoke **"Confirm order on WhatsApp"** flow + measurement form. Build the M2 transactional motion in as you go from `07-MOTION-AND-MICROINTERACTIONS.md` §4.4–4.5 (drawer choreography, 4-state action button, cart count) — do not retrofit. Use the WhatsApp number from env; do not hardcode it.

### ⏰ PARKED — netcup billing (deadline 13 Aug 2026)
Staging deploy is blocked only on the server existing. netcup order **54429992** was placed on a **12-month term** and invoiced **€68.64 upfront**; user wants monthly/no-minimum instead and has asked netcup to amend it. netcup replied that they will attend to the request. **The invoice is deliberately unpaid — do not pay it, and nothing is provisioned until it is.** If netcup has not resolved this by ~**10 Aug**, the formal withdrawal email must go out before the **14-day window closes ~13 Aug**, or the 12-month term binds. Surface this to the user if it goes quiet.

**Blocked on launch credentials/config:** public HTTPS/staging URL for real Paystack webhook delivery, Paystack live keys before production charges, SendByte verified sender/domain before live email delivery, real WhatsApp business number, R2 credentials/product photography/content. Nomba remains reserved for later provider fallback and does not block the active Paystack path.

### ✅ Nav dropdown hover+click fixed by Claude (`7848780`)
The dropdown now opens on **hover and click**. Bug fixed: hover opened it, then clicking the same trigger closed it — and because the pointer was still inside, no further `mouseenter` could reopen it. `useStorefrontNavigation` now tracks `pointerItemId`; `toggleDropdown` refuses to close a menu the pointer is currently holding open, while keyboard and touch activation (no pointer) still toggle normally. Verified live: hover opens → click keeps open → pointer-leave closes → click-without-pointer opens → click again closes.
**Verification lesson for both agents:** asserting dropdown state via `document.querySelector('.desktop-navigation-dropdown')` is wrong three ways — the panel uses `v-show` (always in the DOM), there are **two** panels (Casuals + Collection) so `querySelector` returns the wrong one, and Vue flushes DOM updates asynchronously so a synchronous read after `.click()` is stale. Assert on `getComputedStyle(document.getElementById('collection-navigation-menu')).display`, and read in a separate round-trip after a click.

### ✅ M1 motion foundation — DONE by Claude (`d2544dc`), claim released
Shipped: `design-system/_motion.scss` tokens, the global `prefers-reduced-motion` block, site-wide `:active` tap feedback, the three existing transitions retuned to tokens, and the **marquee WCAG 2.2.2 fix**. Codex is clear to work anywhere again. **From now on, all timing/easing must come from `_motion.scss` tokens — no raw `ms` or `cubic-bezier()` in component files** (currently zero remain outside the token file; keep it that way). M2 transactional motion (cart drawer, 4-state action button, wishlist, S3 WhatsApp confirm) ships *with* the cart in Milestone B — build it from `07-MOTION-AND-MICROINTERACTIONS.md` §4.4–4.5 rather than retrofitting.

### 🔍 DESIGN AUDIT (Claude, 30 Jul) — findings + status
Compared every page against `design-source/` by extracting rendered text and computed styles from both sides.

**FIXED (`3678100`):**
- Category pages ×4 restored to the approved design: light `#f4f0ec` band, `Home / X` breadcrumb, left-aligned **40px/700 sentence-case** H1 (was a dark hero with an uppercase title + invented eyebrow), the design's real descriptions, and per-section chips (Suits `Three-Piece/Two-Piece/Blazers`, Trad `Kaftans/Agbada/Aso-Oke`, Casuals none). Removed the invented `The Collection / X Collection` heading and the sort control (design has neither on category pages).
- **`StorefrontBespokeServiceBand`** added to the layout — it appears above the footer on *every* design page and was missing entirely. Symbol SVG, champagne eyebrow (`--tracking-section`), 34px/500 `#f3eee6` heading, champagne 47px button.
- Cart WhatsApp copy now matches the design's voice ("Hi Oyecouture! I'd like to order the following:" … "Please confirm availability, sizing and delivery details.") with `x1` quantities.

**Grid + spacing corrected (`later commit`):** measured at a **1440px** reference — the design uses **different column counts per page**: Home is **4-up** (313px columns), category pages are **3-up** (425px). Both share a **22px gap** and **56px side gutters**. Ours was 4-up everywhere with tighter gutters. Category `h1` is **56px at desktop**, not 40px — the 40px figure recorded earlier was the *mobile* size, measured before the browser was widened. **Lesson for both agents: always widen the viewport before measuring the design; the default preview width is ~537px and silently yields mobile values.**

**Steps 10–11 delivered (30 Jul):** orders persist with server-side pricing; `/confirm` (bespoke) and `/checkout` (ready-to-wear) both live; delivery zones DB-driven. **All routes return 200 — the `/checkout` 404 is closed.**

**STILL OUTSTANDING:**
1. **`/about` 404s — the entire About page is unbuilt.** Design has hero "YOUR CLOTHING SPEAKS / LONG BEFORE YOU DO", tab nav, and four sections (01 Our Story, 02 The Oyecouture Symbol incl. the CONFIDENCE/SIMPLICITY/CLASS/YOU grid, 03 Craftsmanship, 04 Mission & Vision). All copy is in the design; the footer links to it from four places.
2. **Home shows 27 products; the design shows a curated 13** (monarch, sovereign, regent, agbada, heritage, emir, senator, classic, linen, polo, heir, fila, chino). Needs a product decision — full catalogue or curated.
3. Sort labels use `Price: Low to High`; design uses `Price · Low to High` (middot), and the design has no "Sort by" label.
4. `/checkout` still 404s (step 11).

**VERIFIED CORRECT:** announcement bar · nav + dropdowns · footer (every column/link/legal line) · product detail structure and copy · product card structure.

**KEPT (previously user-approved deviations):** `Explore Oyecouture / Shop By Category` on Home; the Add-to-Bag icon.

**Bespoke rule confirmed:** WhatsApp routing keys off `hasSizes === false` (17 products: all suits + trad), **not** `tag === 'Bespoke'` (only 3). Verified the design behaves the same — `emir` is untagged yet its cart CTA is "CONFIRM ORDER ON WHATSAPP". Using the tag would send 14 unsized products to card checkout with no measurements.

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
- Ghost admin phase 1 is complete: authenticated list/detail APIs, admin orders desk, sidebar, metrics strip, numbered pagination, search, VBAN-style single Filter button with status/type/from/to categories, date-range API support, applied filter chips, order detail page, payment-link copy, WhatsApp reply, and safe manual status actions that cannot mark an order `paid`.
- Paystack payment foundation is complete for sandbox: ready-to-wear checkout redirects to `/pay/<payToken>`, payment attempts are stored per order, Paystack initialize returns hosted checkout URLs, callback verification checks provider status/amount/currency/reference server-side, webhook signatures use Paystack HMAC SHA512 over the raw body, and order status flips to `paid` only through verified successful payment confirmation.
- Paystack initialize is now duplicate-charge guarded: unpaid orders reuse an existing initialized checkout URL/reference, and `payment_attempts_one_initialized_attempt_per_order_idx` enforces at most one initialized checkout attempt per order.
- **Claude review sign-off (2026-07-29):** A.1–A.3 independently reviewed and re-verified (`go vet`/`build`/`test` all pass). Approved with no changes. Claude's three earlier design-token findings (display font, champagne `#e0cf9c`, type scale) were already resolved by Codex in `576ee77`; no fix commit needed. Golden invariants confirmed (integer `priceKobo`, hidden `R2Key`, validated inputs, active-only reads, RTW/bespoke split).
- **Claude review sign-off — A.4 (2026-07-29):** Storefront shell reviewed and approved. Design fidelity verified against `design-source/` (announcement copy exact; nav Home/Trad/Suits/Casuals▾/Collection▾; Casuals→Men/Women; the COLLECTION names "Dear Groom / Wei-Xin / OYE — Man of Legacy" confirmed present in all 7 design files — extracted, not invented; footer Shop/About/Customer Care + IG/X/TikTok/WhatsApp + legal links match). Code quality strong: nav state in a composable, full a11y (aria-controls/expanded, Escape, click-outside, close-on-route-change), lucide tree-shaken icons, search/account/wishlist/cart correctly emit events for future features. One required A.5 follow-up flagged: remove `:external` per link as pages land (see Decisions). Note: Claude's local `npm run typecheck` was inconclusive due to a missing `oxc-parser` native dep corrupting a generated `.nuxt` file (env artifact, not a source defect); Codex's own typecheck/build/`go test` pass.

## In progress (exact state)
- Paystack checkout/webhook work is implemented and locally verified but not deployed. SendByte-backed post-payment customer emails are implemented: two separate emails (order confirmation + payment receipt) are sent after server-side successful verification, protected by an `outbound_emails` outbox table and deterministic idempotency keys. Templates now use an image-free Oyecouture wordmark so SendByte/local previews do not show a broken hosted logo. Ghost admin phase 2 product/content management remains next after payment/email verification.

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
- Public HTTPS/staging URL for real Paystack webhook delivery, Paystack live keys before production, SendByte verified sender/domain for live delivery, WhatsApp business number, delivery-fee amounts per zone, product photography, customer-care page copy, domain DNS access. None block local Paystack hosted-checkout testing, but they gate launch/deployment polish.

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
- Paystack is the active payment provider for the first payment phase; Nomba remains reserved for later. Multiple Paystack attempts may exist for one order, but the repository locks the attempt/order row and only one verified successful confirmation can move an awaiting-payment order to `paid`.
- Multiple non-terminal local attempts may exist for audit/recovery, but only one attempt per order may remain `initialized`; repeated Pay clicks return that same Paystack checkout URL/reference until the attempt becomes terminal.
- Browser callback verification must provide a reference belonging to the same order as the `payToken`; webhook confirmation is reference-based because Paystack owns the event source and the signature is validated before processing.

- **Typography (2026-07-29, Claude):** The design's real face is **Gilmer** (commercial, Fontfabric). Per user decision, the four Gilmer weights were extracted from the design bundle and self-hosted as WOFF2 in `web/public/fonts/` (`_fonts.scss`), replacing the Jost stand-in. A `--tracking-*` letter-spacing scale (values taken from the design) now drives all uppercase labels/headings/buttons. ⚠️ **Gilmer needs a proper webfont licence before production launch** — self-hosting it now is for build/demo fidelity; do not launch commercially on the unlicensed font without confirming licensing with the client. The `@fontsource-variable/jost` dependency is now unused and can be pruned. The visual side-by-side against `design-source/Oyecouture Shop.dc.html` was completed in `9cc1333`; desktop/mobile typography and catalogue details were corrected to the source values.

- **Motion system (2026-07-29, Claude):** `07-MOTION-AND-MICROINTERACTIONS.md` is now the authority for anything that moves — **read it before adding any animation.** Key rules: use `_motion.scss` tokens (no raw ms/cubic-bezier), animate only composited properties (`transform`/`opacity`/`filter`/`clip-path`), exits faster than entrances, **no spring/bounce easing**, tap `:active` feedback on everything (mobile has no hover), and nothing on the cart→checkout→payment path exceeds `--duration-slow`. Stack is **CSS-first**: CSS transitions → Vue `<Transition>` → CSS scroll-driven animations (`@supports` guarded) → View Transitions API via Nuxt `app.viewTransition: true`. **No motion library is installed — do not add one without recording the justification here.** ⚠️ **WCAG 2.2.2 launch blocker:** the infinite hero marquee needs a pause/stop mechanism (§6.2).

- **M1 motion implementation notes (2026-07-29, Claude, `d2544dc`):** (a) Press feedback uses the independent `scale` property, **not** `transform: scale()` — this composes with elements that already carry a transform (e.g. `.product-carousel-button`'s `translateY(-50%)`), which `transform: scale()` would have silently clobbered. Same reason `.category-tile img` hover now uses `scale: var(--zoom-image)`. (b) **Design deviation — hero marquee:** on `(hover: none)` devices the auto-scroll is dropped and the rail becomes a swipeable overflow strip; on `(hover: hover)` it animates and pauses on `:hover`/`:focus-within`. Rationale: touch devices cannot hover to pause, so hover-pause alone would leave the majority-mobile audience non-compliant with WCAG 2.2.2; this also saves battery on mid-range Android. Visually the rail is unchanged — only its motion behaviour differs by input type. Flag to the user if they want the auto-scroll back on touch (it would then need a visible pause control).
- **Local env caveat resolved for current checkout:** the earlier Nuxt `oxc-parser`/`oxc-walker` optional parser issue was fixed by keeping portable `oxc-parser` in web dev dependencies and removing the machine-specific Darwin binding. On 2026-08-02, `npx vue-tsc --noEmit -p tsconfig.json` and `npm run build` both pass on Node `22.23.1`.

## Reminders (non-negotiable)
Money in kobo (int64, no floats) · `paid` only after server-side verify · guest checkout only in phase 1 · mobile-first · ghost admin stays unlinked/unmentioned to client · commits as `mikechibuzor`, no `Co-Authored-By`.
