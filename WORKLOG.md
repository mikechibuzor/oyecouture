# WORKLOG — append-only (newest on top)

One short entry per work session or meaningful chunk. Format:
`## <ISO datetime> · <Codex|Claude> · <one-line summary>` then bullets: **Did / State / Next / Blockers**.

---

## 2026-08-11T15:50:14+01:00 · Codex · Admin settings deployed
- **Did:** Added DB-backed storefront settings using the existing `settings` table. Public `GET /site-settings` now powers the announcement bar and storefront WhatsApp number; authenticated `GET/PATCH /admin/settings` powers the new `/admin/settings` page. Admin Settings currently manages announcement text and WhatsApp number; contact owner email remains env-backed.
- **State:** Committed/pushed/deployed as `d400339`. Verified locally with `go test ./...`, Node 22 `npm run typecheck`, Node 22 `NUXT_IGNORE_LOCK=1 npm run build`, `git diff --check`, local endpoint smoke tests, and live staging checks: `/api/healthz` 200, `/api/site-settings` 200, `/api/admin/settings` 401 without a session.
- **Next:** User visually tests `/admin/settings`; then continue admin phase 2 with home-content controls.
- **Blockers:** None for this slice.

## 2026-08-11T14:38:16+01:00 · Codex · Admin money inputs switched to naira locally
- **Did:** Changed admin product and delivery money inputs to display/edit naira while keeping API/database contracts in kobo. Product drafts now use `priceNaira` and convert to `priceKobo` only in `useAdminProducts`; delivery drafts use `feeNaira` and convert to `feeKobo` only in the save handler.
- **State:** User tested locally and approved. Committed/pushed/deployed as `b4f2ab4`; live health verified 200 after deploy.
- **Next:** Continue admin phase 2 settings/home-content controls.
- **Blockers:** None for this slice.

## 2026-08-11T10:16:25+01:00 · Codex · Admin delivery management deployed
- **Did:** Added authenticated admin delivery-zone management: `GET /admin/delivery-zones`, `PATCH /admin/delivery-zones/{name}`, a Nuxt `/admin/delivery` page, `useAdminDeliveryZones`, sidebar navigation, fee/active controls, save loading, expired-session handling, and top-right success notifications. Public checkout still only sees active zones. Also tightened `FeeForZone` so real DB errors no longer look like an invalid delivery choice.
- **State:** Committed/pushed as `6bb208a` and deployed to staging. Verified locally with `GOCACHE=/private/tmp/oyecouture-go-cache go test ./...`, Node 22 `npm run typecheck`, and Node 22 `NUXT_IGNORE_LOCK=1 npm run build`. Verified live health 200 and `/api/admin/delivery-zones` returns 401 without an admin session.
- **Next:** Visually test `/admin/delivery` with the user's admin session, then continue phase 2 into settings/home-content controls.
- **Blockers:** Same launch blockers: user-owned end-to-end Paystack/order/email test, netcup withdrawal, Hostinger renewal decision, credential rotation, real suit photography.

## 2026-08-11T09:00:00+01:00 · Codex · Latest main deployed to staging
- **Did:** Deployed latest `origin/main` to Hostinger staging with `ssh oyecouture 'cd ~/oyecouture && git pull --ff-only && ./deploy/deploy.sh'`. This put contact form + throttle, product shipping/returns disclosures, disclosure marker fix, and real shipping/returns policy pages live.
- **State:** Staging now runs `a592859` (`feat(care): populate shipping and returns from the real policy copy`). Deploy rebuilt/restarted API and web containers successfully and reclaimed disk, leaving 41G free of 48G.
- **Verification:** Live `https://82-29-184-79.sslip.io/api/healthz` returned 200 with API/database ok. Live `OPTIONS /api/contact` from the staging origin returned 204 and the expected CORS headers.
- **Next:** User will run the end-to-end Paystack/order/email test. Codex continues admin phase 2, starting with visual admin product verification and then delivery/settings/content controls.
- **Blockers:** Netcup withdrawal still time-sensitive (~13 Aug) unless already sent; Hostinger renewal economics; credential rotation; no suit photography.

## 2026-08-04T23:30:00+01:00 · Claude · Review of Codex's work + contact form throttle
- **Reviewed** the 12 commits since handoff plus the contact endpoint. **Codex's work is good.** All tests pass, `gofmt`/`vet` clean, staging healthy at `a773f80`. Verified: all 9 admin product routes behind `RequireSession`; the Cloudinary secret is server-side only and never reaches the browser; image upload caps at 12MB, **sniffs content type rather than trusting the header**, and allowlists jpeg/png/webp; the dev CORS relaxation is correctly gated behind `appEnv == "production"` — confirmed live that staging gives no CORS header to a foreign origin *or* to localhost. Codex also improved on my webhook fix, adding `APIError`/`ErrVerifyRejected` so Paystack 4xx/404 are terminal while 5xx/timeouts retry. My earlier WhatsApp finding was fixed correctly (`NUXT_PUBLIC_WHATSAPP_NUMBER` now reaches the web container; live in the SSR payload).
- **One real finding, fixed (`86d89a3`):** `POST /contact` was public, unauthenticated and unthrottled while sending an email per submission. Now 5 per hour per IP, charging only messages that actually send. Extracted `internal/ratelimit` and migrated admin sign-in onto it so the client-address handling exists once.
- **Two process notes:** Codex shipped 12 commits before writing any worklog entry (since backfilled) — a usage-limit cutoff in that window would have lost the context. And Codex recorded that browser verification was blocked by the in-app browser not being signed in to admin, so **the entire product management UI has been typechecked and built but never visually verified**.
- **Baton correction:** the `Uncommitted local work` line described the contact endpoint as uncommitted; it is committed as `76e8848`. Fixed.
- **Next:** unchanged and now overdue — **run one order end to end**. It has been the stated next action across two handoffs. Also: deploy the contact endpoint *with* the throttle, and visually verify the product admin.
- **Blockers:** ⏰ netcup (€68.64, ~13 Aug) · Hostinger renewal · three credentials to rotate · no suit photography.

## 2026-08-04T21:35:00+01:00 · Codex · Contact form wired to real email endpoint
- **Did:** Added public `POST /contact` in the Go API, with server-side validation for name/email/topic/rating/message and SendByte-backed delivery to `OWNER_NOTIFY_EMAIL` (defaulting to `iwuagwuchibuzor@gmail.com` for now). Wired the Contact page through a lean Nuxt composable instead of keeping submit logic in the template. Follow-up fixed dev CORS so `http://127.0.0.1:3010` can call `:8080` while production remains locked to the configured site origin.
- **State:** Local API restarted on `:8080`; storefront dev server remains on `http://127.0.0.1:3010`. Smoke-tested `/contact` route with an invalid payload and got the expected `invalid_submission` response without sending email. Browser-origin preflight now returns `Access-Control-Allow-Origin: http://127.0.0.1:3010`; a valid contact POST returned `201`, though one prior SendByte sandbox call timed out before succeeding on retry.
- **Verification:** `GOCACHE=/private/tmp/oyecouture-go-cache go test ./...`, `PATH="$HOME/.nvm/versions/node/v22.12.0/bin:$PATH" npm run typecheck`, and `PATH="$HOME/.nvm/versions/node/v22.12.0/bin:$PATH" NUXT_IGNORE_LOCK=1 npm run build` pass. The shell's default Node `20.12.1` is too old for Nuxt 3.21's parser stack; use Node 22 for frontend checks.
- **Next:** Commit/push/deploy this contact endpoint when the user approves; optionally send one real contact-form test email to confirm SendByte delivery into the owner's inbox.
- **Blockers:** Local/prod email sending still depends on valid `EMAIL_PROVIDER`, `SENDBYTE_API_KEY`, and `EMAIL_FROM`; production already has live email config per the current baton.

## 2026-08-03T02:10:00+01:00 · Claude · Handing the baton to Codex
- **Did:** SendByte **production** key installed on the server (`sk_live_`); the local `.env` deliberately keeps the sandbox key so local runs cannot send real mail. Verified the live key is nowhere in git. User seeded their own admin account on staging and confirmed they can sign in.
- **State:** Staging is fully working at **https://82-29-184-79.sslip.io** — 30 active products (22 real + 8 placeholder suits), 91 images all verified loading from Cloudinary, ghost admin reachable at `/admin`. **Sandbox for money (`sk_test_` Paystack), LIVE for email (`sk_live_` SendByte).** 5 test orders in the staging DB to purge before go-live.
- **Baton correction:** the `Current status` block at the top of `HANDOFF-STATE.md` had gone stale — it still read "Last worked by: Codex · 2026-08-02" and described SendByte as local-only, while every section below it was current. Rewritten, and the superseded "next action" heading demoted so only one is live.
- **Next (for Codex):** run **one order end to end** — order → Paystack sandbox card → webhook → `paid` → email in a real inbox. Every piece is in place; the joins have never been exercised together. Then client price review, then ghost admin phase 2 (product create/edit + Cloudinary upload).
- **Blockers:** ⏰ netcup withdrawal (€68.64, ~13 Aug, user sending it) · ⏰ Hostinger renews at $19.49/mo which exceeds the ₦25k retainer — move to 12-month or cancel ~28 days from 2 Aug · no suit photography · three credentials pasted in chat need rotating.

## 2026-08-03T01:40:00+01:00 · Claude · Real catalogue live on Cloudinary
- **Did:** Uploaded the client's photography (709MB → 102MB, 92 images) to Cloudinary and replaced the design-mockup catalogue with 22 real products — 11 aso-oke looks and 11 casual two-pieces (migration `000010`). Placeholders are **deactivated, not deleted**, so order history keeps its links. Suits stay on placeholders (no suit shoot yet). Also uploaded the 58 legacy mockup images, because `R2_PUBLIC_BASE` is a single value and pointing it at Cloudinary broke the remaining suits otherwise.
- **State:** 30 active products, 91 image rows, **every image verified loading**. Credits at 0.4% of 25/month. `R2_PUBLIC_BASE=https://res.cloudinary.com/exgithrw/image/upload/f_auto,q_auto`. Category tiles for Trads/Casuals now use the client's own shots (`campaign-4`, `casual-14`); Suits still the mockup mannequin.
- **Three things worth not repeating:** (1) `product_images_position_valid` capped position at **3** — a mockup-era assumption — while real looks run to 6 frames; migration 000010 lifts it to 12 first. (2) Resetting a dirty migration while the API is **crash-looping** does not work — the old container re-applies and re-dirties it within seconds; stop the api container first. (3) `.product-card` carried `reveal-on-scroll`, which ties opacity to scroll progress; a tall card low in the viewport never completes its range and renders at ~0.96 opacity the whole time it is visible, so identical garments looked different depending on grid position. Removed from cards, kept on headings. Verified opacity 1 at seven scroll positions.
- **Measured delivery sizes (do not re-estimate):** casuals ~99KB, aso-oke 125KB at w_600 / 408KB at w_1200. Realistic ceiling ≈ **16,000 page views/month**, not the 29–40k first projected.
- **Next:** SendByte **production** key still to go into `SENDBYTE_API_KEY` (DKIM + merged SPF are both propagated, so the domain should verify) — note that real order emails start sending once it is in. Then client review of the 22 dummy prices, and ghost admin phase 2 product/image management.
- **Blockers:** ⏰ **netcup withdrawal drafted but UNSENT — €68.64, window closes ~13 Aug.** Hostinger renews at $19.49/mo unless moved to a 12-month term. No suit photography. Cloudinary API secret was pasted in chat — regenerate.

## 2026-08-02T20:30:30+01:00 · Codex · Polished paid-order email templates
- **Did:** Reworked the order confirmation and payment receipt templates with a premium image-free Oyecouture wordmark, warmer framed layout, preheaders, summary/status blocks, cleaner totals, and a fulfilment note. Removed reliance on a hosted logo so SendByte/local previews no longer show a broken image.
- **State:** Verified `go test ./...` and `git diff --check`; new test asserts paid-order emails do not include external `<img>` tags and do include the brand wordmark.
- **Next:** Restart API, then create a new paid test order to view the updated templates in SendByte. Existing sent orders keep their original template because outbox idempotency prevents duplicate sends.
- **Blockers:** Live inbox delivery still needs SendByte live key + verified sender/domain.

## 2026-08-02T23:00:00+01:00 · Claude · Staging live on Hostinger
- **Did:** Provisioned and hardened Hostinger KVM 1 (`82.29.184.79`, Ubuntu 24.04.4, 3.9GB/1vCPU/48GB): `deploy` user, 2GB swap, Docker, `ufw` 22/80/443, key-only SSH. Created a **dedicated** `~/.ssh/oyecouture` key — the user's only existing key was the work one (`chibuzor@ovalfi.com`), and client production must not share a credential with the employer. Pushed the repo to **private** GitHub `mikechibuzor/oyecouture` (it previously existed only on the laptop); server pulls via a read-only deploy key. Deployed the full stack.
- **State:** **https://82-29-184-79.sslip.io** live with a real Let's Encrypt cert and HTTP→HTTPS 308. All five containers up. Verified from the public internet: storefront, `/admin`, and every API route 200; **5432 and 8080 unreachable from outside**. Fresh DB at v9 — 27 products, 3 zones, 0 orders. Sandbox Paystack/SendByte keys, `WHATSAPP_NUMBER=2348094939264` (user's own, for testing).
- **Three bugs only a real deploy could surface, all fixed:** Dockerfile on Go 1.23 vs `go.mod` 1.25 (`456bbdb`); the runbook's `openssl rand -base64 32` producing a `/` that broke `DATABASE_URL` and crash-looped the API after a full build (`02cbc66`, now hex + a fail-fast preflight); and the ghost admin having **no account and no way to create one** on a server, since the CLI needs Go and distroless has none (`7660408`, binary now ships in the image).
- **Next:** User runs `docker compose exec -it api /admin` to set their own admin password (deliberately not set by me). Then register `POST /webhooks/paystack` against the staging HTTPS URL — the last untested part of the payment path. Then ghost admin phase 2.
- **Blockers:** ⏰ **netcup withdrawal still not sent — €68.64, closes ~13 Aug.** Hostinger renews at $19.49/mo (~₦26,900) after the first month unless switched to a 12-month term. Live Paystack keys, SendByte verified sender, Gilmer licence for launch.

## 2026-08-02T19:55:06+01:00 · Codex · SendByte paid-order emails
- **Did:** Implemented SendByte transactional email after verified Paystack payment: two separate customer emails (order confirmation + payment receipt), `outbound_emails` outbox/audit table, deterministic idempotency keys, stale-pending retry recovery, local HTML/text templates with Oyecouture branding, and payment-service hooks for callback/webhook settlement.
- **State:** SendByte sandbox key is only in ignored `.env`; committed env files have placeholders. Verified `go test ./...`, `npx vue-tsc --noEmit -p tsconfig.json`, `npm run build` on Node `22.23.1`, and `git diff --check`. No real/sandbox email was sent in this session.
- **Next:** Optionally trigger one sandbox payment verification with explicit approval to confirm both emails appear in SendByte logs, then continue ghost admin phase 2 product/content management.
- **Blockers:** Public HTTPS/staging URL for real Paystack webhooks, Paystack live keys before production, SendByte verified sender/domain for live delivery, real WhatsApp number, R2 credentials/photos/content/DNS.

## 2026-08-02T17:20:00+01:00 · Claude · Review of Codex's ghost admin + Paystack work
- **Did:** Reviewed `internal/payments` and the ghost admin against the running stack rather than reading only. Confirmed live: forged/unsigned webhooks 401; `PATCH .../status` refuses `paid` and illegal transitions with 400 and 401s unauthenticated, with the DB row unchanged; admin list/detail/metrics 200; unknown pay token 404. Confirmed by reading: amount **and** currency checked against both attempt and order total, `FOR UPDATE OF pa, o` locking, idempotent re-confirmation, cross-order reference binding, webhook signature over the raw body with `hmac.Equal` before parsing, secret never logged, `.env` never committed, PII excluded from the public payment JSON.
- **Fixed (`d29c830`):** the amount/currency guard had no test — the service tests stub the repository, so it was never exercised. Extracted `verificationMatchesOrder` as a pure function and covered 10 cases (underpayment, one kobo short, overpayment, zero, tampered attempt, repriced order, three currency cases). Mutation-checked: removing the guard fails six of them.
- **State:** `gofmt`, `go vet`, `go test ./...` and `vue-tsc --noEmit` all clean. Migrations at version 8, not dirty. Codex's work is good — the money path is the strongest code in the repo.
- **Next:** Two decisions recorded in `HANDOFF-STATE.md` for the user, neither a silent fix: (1) terminal webhook errors return 500 and so retry forever — acknowledge with 200 and alert instead, keeping 500 for transient failures; (2) `provider_payload` stores Paystack's full response including customer email/phone/IP and card metadata with a reusable `authorization_code` — consider persisting a trimmed payload. Neither blocks phase 2.
- **Blockers:** unchanged — public HTTPS URL for real webhook delivery, remaining launch credentials/content.

## 2026-08-02T16:26:35+01:00 · Codex · Paystack duplicate-charge guard
- **Did:** Added Paystack initialize idempotency: unpaid orders now reuse an existing initialized Paystack checkout URL/reference, and migration `000008` enforces one initialized payment attempt per order with a partial unique index. Concurrent initialize races now fail the losing local attempt and return the surviving checkout URL.
- **State:** Verified `go test ./...`, `go test ./internal/payments`, `npx vue-tsc --noEmit -p tsconfig.json`, `git diff --check`, local migration index creation, and two near-simultaneous initialize calls for `OYE-1068` returning the same reference/URL.
- **Next:** Continue sandbox payment testing with this guard in place; later production hardening can add provider-side expiry cleanup/reconciliation views.
- **Blockers:** Public HTTPS/staging URL for real webhook delivery and remaining launch credentials/content.

## 2026-08-02T16:02:41+01:00 · Codex · Paystack payment foundation
- **Did:** Continued from the payment planning and implemented the Paystack-first payment flow: `payment_attempts` migration, Paystack initialize/verify provider, payment service/repository with idempotent paid settlement, signed `/webhooks/paystack`, public `/payments/{payToken}` APIs, `/pay/{payToken}` and callback pages, and checkout redirect-to-payment for ready-to-wear orders.
- **State:** Sandbox secret lives only in ignored `.env`; committed files contain empty env placeholders only. Verified `go test ./...`, `go test ./internal/payments`, `npx vue-tsc --noEmit -p tsconfig.json`, `npm run build`, and a secret scan outside `.env`. Callback verification now requires the reference to belong to the same order as the pay token; webhook confirmation remains reference-based.
- **Next:** Test a real sandbox card through the browser, then expose the API on public HTTPS and register `POST /webhooks/paystack` in the Paystack dashboard.
- **Blockers:** Public HTTPS/staging URL for real webhook delivery, Paystack public key if inline payment is later desired, Resend API key, real WhatsApp number, client delivery fees/photos/content, and netcup server/deploy.

## 2026-08-02T15:15:48+01:00 · Codex · Admin order desk click polish
- **Did:** Committed `ac24935` for taller metrics and table column alignment; committed `2be9f3c` so clicking or keyboard-opening an order row navigates to the existing detail page while nested links keep native behavior.
- **State:** Verified `npx vue-tsc --noEmit -p tsconfig.json`, `npx sass --no-source-map assets/scss/main.scss /tmp/oyecouture-admin.css`, and `git diff --check`.
- **Next:** Phase 2 product/content management remains next.
- **Blockers:** R2/upload credentials, Nomba sandbox keys, Resend API key, real WhatsApp number, client delivery fees/photos/content, and netcup server/deploy remain unchanged.

## 2026-08-02T14:55:21+01:00 · Codex · Admin filter popover responsiveness
- **Did:** Committed `9619edc` to make the VBAN-style admin filter popover fixed-position and viewport-clamped, with resize/scroll repositioning and internal overflow for smaller screens.
- **State:** Verified `npx vue-tsc --noEmit -p tsconfig.json`, `npx sass --no-source-map assets/scss/main.scss /tmp/oyecouture-admin.css`, and `git diff --check`.
- **Next:** Phase 2 product/content management remains next.
- **Blockers:** R2/upload credentials, Nomba sandbox keys, Resend API key, real WhatsApp number, client delivery fees/photos/content, and netcup server/deploy remain unchanged.

## 2026-08-02T14:50:46+01:00 · Codex · Ghost admin phase 1 complete
- **Did:** Committed `a8a372c` with numbered pagination, left-side filters, half-width search, and quieter admin microinteractions; committed `d701dfd` with authenticated order detail/read API, status update API, detail page, line items, customer/address/measurements, totals, payment-link copy, WhatsApp reply, and guarded manual status actions.
- **State:** Ghost admin phase 1 is complete. Admin actions cannot mark an order `paid`; payment remains reserved for server-side verification/webhook. Verified `go test ./...`, `npx vue-tsc --noEmit -p tsconfig.json`, `npx sass --no-source-map assets/scss/main.scss /tmp/oyecouture-admin.css`, and `git diff --check`.
- **Next:** Phase 2 product/content management. Start with product create/edit + image metadata flow, keeping DB/R2 media boundaries ready while credentials are still missing.
- **Blockers:** R2/upload credentials, Nomba sandbox keys, Resend API key, real WhatsApp number, client delivery fees/photos/content, and netcup server/deploy remain unchanged.

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
