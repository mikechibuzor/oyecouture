# WORKLOG — append-only (newest on top)

One short entry per work session or meaningful chunk. Format:
`## <ISO datetime> · <Codex|Claude> · <one-line summary>` then bullets: **Did / State / Next / Blockers**.

---

## 2026-10-06 · Claude · Shop availability toggle, and the bespoke heading asks the question

- **Did:** (1) Changed the bespoke band heading to "Would you like to have a bespoke piece?"
  (`9829d20`). The band already carries "Book an Appointment" as its action, so the heading only
  had to pose the question; the old line also repeated the measure wording the size guide and
  measurements pages already make. The About page passes its own heading and is untouched.
  (2) Built the shop availability hold (`9990012`), ported from Fantasy by HOG: an owner-only
  toggle plus editable headline and message in Settings, a gate on every visitor-facing API
  route, a holding page on the brand's own tokens, and migration 23 seeding the three keys off.
- **State:** Both committed; the copy change is **pushed**, the hold is **committed but not
  pushed and not deployed** — migration 23 has not run anywhere. `go build`, `go vet`, the full
  Go suite, `nuxt typecheck` and a production Nuxt build all pass on Node 22. `holdGate` is a
  plain function with five unit tests: open, held, admin bypass, failed read (stays open), and
  held-despite-error (stays closed).
- **Decision worth knowing:** Chibuzor chose "fully dark". I held everything a person can reach
  but left `POST /webhooks/paystack` open, because it is Paystack's server rather than a visitor
  — blocking it would charge a customer and leave their order unpaid. Flagged to him; one line
  in `server.go` if he wants it blocked too.
- **Next:** Run it against a live stack and a browser — nothing here has been exercised beyond
  compilation. Toggle the hold, confirm a visitor is held, an admin is not, and the webhook still
  answers.
- **Blockers:** No local Postgres was running, so no end-to-end check was possible. A deploy here
  restarts the containers and needs Chibuzor's go-ahead.

## 2026-09-15T18:43:05+01:00 · Claude · HOG staging added to the running Caddy, no code change here

- **Did:** At Chibuzor's request, added `hog-staging.82-29-184-79.sslip.io` (the Fantasy by HOG staging site) to `oyecouture-caddy`. It used HOG's `deploy/staging/caddy-live.sh`: a `docker network connect hog-edge` and a graceful reload. Only these notes were pushed to this repo; they are docs-only, so the auto-deploy does not rebuild.
- **State:**
  - Caddy was not restarted.
  - The storefront returned 200, `/api/healthz` ok and `www` 301, both before and after.
  - Let's Encrypt issued the HOG certificate.
  - The running Caddy config is this repo's Caddyfile plus HOG's block.
- **Next:** Ship the permanent version with the next real Oyecouture release. See HANDOFF-STATE.
- **Blockers:** None for Oyecouture. A reboot or a Caddy restart drops HOG staging until `caddy-live.sh` is re-run.

## 2026-09-06T21:55:00+01:00 · Claude · Storefront was invisible to search; metadata built and deployed

- **Did:** Diagnosed why the shop ranked below its own Instagram and Facebook for its own name, and fixed it. Added `constants/seo.ts`, `composables/useStorefrontSeo.ts` (title, description, canonical, Open Graph, Twitter), `utils/seo.ts` (ClothingStore, WebSite and per-product Offer schemas), and `server/routes/sitemap.xml.ts`. Wired the homepage, all four category pages and the product detail page, moved the care pages' metadata into `StorefrontCarePage`, and converted about/contact/size-guide off bare `useSeoMeta` so they gain canonical and link previews.
- **State:** The Google result read **"No information is available for this page"** — not a ranking problem. The homepage served 295KB of HTML with **no `<title>` and no meta description**, as did every category and product page. `sitemap.xml` 404'd, so the 27 product pages were reachable only from whatever the homepage linked. No canonical, and no Open Graph, so links shared to WhatsApp previewed as bare URLs. **No visible copy changed** — the searchable nouns (kaftan, agbada, aso-oke, Lagos) live in the metadata, so the hero keeps the brand's voice. Verified against the **production build**, not the dev server: all thirteen page types return title + canonical + og:image, the homepage carries two schemas, the product page carries its Offer, and the sitemap returns 40 URLs (13 static, 27 live products). Full Go suite passes (20 packages).
- **Next:** **Google Search Console needs the owner's account** — verify the domain and submit `/sitemap.xml`; that is the fastest remaining lever and no agent can do it. Then finish the audit-log detail below.
- **Blockers:** Expectation to set with the client: brand terms (`oyecouture`, `oyecouturebespoke`) are winnable within weeks, but head terms like `kaftan` or `tailoring` are not reachable by code — those belong to global retailers. The titles target local and long-tail intent instead ("bespoke agbada Lagos", "custom kaftan Nigeria"). **Local dev was broken for both agents and now has a cause:** Node v20.12.1 cannot `require()` an ESM module, which `oxc-walker` needs for `oxc-parser`; Docker uses `node:22-alpine`, which is why production built fine and no Mac could. Use Node 22 locally. Launch blockers unchanged: Terms and Privacy are still stubs while the site takes real cards, and production still carries the developer's WhatsApp number.

## 2026-09-02T18:40:00+01:00 · Claude · All four client email designs built and live

- **Did:** Built the remaining three emails to the client's designs after 01 shipped earlier: 02 Shipped, 04 Delivered, and 03 Out for delivery. Added migration `000017` for the handover record and `000018` for the out-for-delivery stage plus rider fields, threaded both through the tracking service, repository, email payload and admin UI, and added tests for the stage rules, the rendered design blocks, and the progress bar.
- **State:** Production on `cbded52`, migration `18:false`, all three rider columns present, stage constraint accepts `out_for_delivery`, health green. Every email was **rendered and inspected**, not merely compiled. `out_for_delivery` is optional — `dispatched` moves either to it or straight to `delivered` — because forcing a courier parcel through a local-rider stage would park international orders somewhere nobody can observe. Handover and rider fields are refused by the API on any transition but their own, and their panels are omitted rather than rendered empty. One near-miss caught before shipping: the public timeline does `FULFILLMENT_PATHS.indexOf(currentStage)`, which returns `-1` for an off-path stage and renders **every** step incomplete, so a paid and shipped order would have looked untouched; `out_for_delivery` now reads as `dispatched` there.
- **Next:** The user's queue: super admin/roles and permissions, then the customer-initiated payment escalation and refund flow. **Build self-service password change first** — `admin_users.role` exists but nothing enforces it on any route, and without a password-change route every account the developer creates keeps a password the developer knows, which makes an audit trail meaningless before it is written.
- **Blockers:** Deliberately omitted from the designs and still open: the **Leave a review** panel (no review system; a test keeps it absent), product photo thumbnails (`email.Item` carries no image URL), and **Email preferences / Unsubscribe** (no system, and a dead link is worse than none). Launch blockers unchanged and needing no code: **Terms and Privacy are still stub sentences while the site takes real cards**, and production still carries the developer's WhatsApp number rather than the client's. Local `npm run typecheck` cannot run in the scratchpad clone (npm optional-dependency native binding bug); the Docker build typechecks properly and is the gate.

## 2026-09-02T12:10:00+01:00 · Claude · A swallowed payment, and the client's launch feedback

- **Did:** Recovered order `OYE-1002`, which was paid and stuck at `awaiting_payment`. Fixed the guard that caused it, made the confirmation page wait for async payments to settle, and cleared the client's post-launch feedback: added a favicon, made section terms searchable, removed duplicated storefront chrome from the confirmation and receipt pages, renamed "the atelier" to "Oyecouture" in customer copy, tightened product-card spacing, and fixed the amount-paid figure.
- **State:** All deployed and verified live. The money bug: Paystack charged ₦101.53 for a ₦100 order because the account passed the transaction fee to the customer on `bank_transfer`; `verificationMatchesOrder` demanded exact equality, rejected the webhook as a mismatch, and because mismatches are terminal it answered 200 so Paystack never retried. **Every bank-transfer order would have failed identically** — card matched exactly, which hid it. The guard now rejects only underpayment. Recovery went through `POST /payments/{payToken}/verify`, not a hand-edited row, so `OYE-1002` is `paid`/`payment_confirmed` with both emails sent and a system-attributed event. Search: Suits 0→7, Casuals 0→10. Confirmation page verified in-browser against the client's design; the `AMOUNT PAID` figure had been rendering at **11px instead of 24px** because a bare `p` label rule outscored `.confirmation-summary-amount` — two properties already carried `!important` to escape the same collision. Node 22 typecheck and the full Go suite pass.
- **Next:** In the user's stated order: super admin/roles/permissions, then a customer-initiated payment escalation and refund flow. **Build self-service password change first** — `admin_users.role` exists but nothing enforces it on any route, and with no password-change route every account the developer creates keeps a password the developer knows, which defeats the point of the roles work.
- **Blockers:** The four email templates are **not** built; a payment placed today still sends the old emails. They need an `out_for_delivery` stage, rider name/contact/window, signature capture, a review system and unsubscribe handling. Two `!important` margin overrides remain in `_confirmation.scss` and `_receipt.scss` with the same nested-`p` specificity smell; cosmetically fine, worth scoping when next in those files.

## 2026-09-01T15:20:00+01:00 · Claude · Production domain live and the money path proven with real money

- **Did:** Moved the site onto the client's own domain and to live payments. Pointed DNS (A `@` → `82.29.184.79`, `www` follows by CNAME) leaving MX and SPF untouched; set `SITE_HOST=oyecouturebespoke.com` and `OWNER_NOTIFY_EMAIL=frontdeskoyecouture@gmail.com`; installed the live Paystack keys server-side without their values passing through the transcript. Added a `www` redirect to the Caddyfile. Purged 36 development orders plus 114 cascaded child rows and reset `order_number_seq` so the first real order was `OYE-1001`, backing the database up first.
- **State:** **The end-to-end chain has now been exercised once, with real money, and it works.** `OYE-1001` for ₦100 on live keys: payment attempt `success`, and `order_fulfillment_events` records `payment_confirmed` at 14:07:17 with `created_by_admin_id IS NULL`, which only the signature-verified webhook writes. `order_confirmation` and `payment_receipt` both `sent` on the first attempt a second later, no API errors. Certificate is Let's Encrypt with a matching CN, valid to 30 Nov; `www` 301s to the apex; HTTP 308s to HTTPS. Autodeploy shipped `b1cd9bf` on its own with no human involvement, confirming the poll loop works unattended.
- **Next:** Restore `adire-linen-set` to ₦105,000 — it is live at ₦100 right now. Then the client's WhatsApp number, real Terms and Privacy copy, and the remaining launch blockers.
- **Blockers:** ⚠️ **Caddyfile edits do not take effect on deploy.** `compose.prod.yaml` bind-mounts the single file, so a pull replacing it leaves the container on the old inode; the container must be force-recreated. `deploy.sh` does not do this and should be taught to. ⚠️ **Real cards now charge real money and Paystack test cards no longer work** — live-mode webhook signatures do not match test-mode ones.

## 2026-09-01T13:25:00+01:00 · Claude · Automatic deploys, and a health gate that never worked

- **Did:** Made deployment automatic without GitHub Actions. Added `deploy/autodeploy.sh`, installed as `~/autodeploy.sh` on the server under a 5-minute cron: it fetches `origin/main`, fast-forwards, and runs `deploy/deploy.sh` when non-documentation files changed, guarded by `flock`. Fixed `deploy/deploy.sh`'s health check, which had never verified anything. Wrote `.github/workflows/deploy.yml` (tests gating deploy, asserted health and migration state) and parked it on `workflow_dispatch` after establishing Actions cannot run on this account. Pushed Codex's stranded `2aee56d`/`df7cbe3`, deployed, and took the first off-server database dump since 25 Aug (`~/oyecouture-20260901T093424Z.sql.gz`, 13 tables, 36 orders, gzip verified).
- **State:** Server on `012264d`, in sync with origin, all five containers up, migration `16:false`, `/api/healthz` ok, and `/`, `/suits`, `/casuals/men`, `/product/monarch`, `/admin/login` all 200. Autodeploy verified twice: a real deploy that built both images and reported `deploy OK: 012264d`, then a second run that correctly no-opped. **The old health gate was a no-op on every deploy this project has ever done** — `curl -fsS http://localhost/api/healthz` gets Caddy's 308 to HTTPS, and curl treats an unfollowed redirect as success, so it exited 0 without reaching the API. It now pins `SITE_HOST` to loopback so TLS SNI matches, asserts `"status":"ok"`, and fails the deploy instead of warning.
- **Next:** The end-to-end order test, still never run: order → Paystack **test** card → webhook → `paid` → email. Must happen on the test keys currently installed; the live keys reject test-signed webhooks. Only after it passes should live keys go in.
- **Blockers:** GitHub Actions is blocked account-wide by a failed card authorization (GitHub billing → Payment information); autodeploy makes this non-urgent. ⚠️ **Cron runs the copy at `~/autodeploy.sh`, not the repo's — re-copy it whenever `deploy/autodeploy.sh` changes.** Airtel Nigeria still cannot route to Hostinger; MTN can. That no longer blocks deploys, since the server pulls for itself, but it does block SSH and admin access from an Airtel connection.

## 2026-09-01T10:30:00+01:00 · Claude · The 31 Aug "outage" was a routing fault, not a server fault

- **Did:** Retracted the previous entry's outage claim after the developer reached the server from an MTN MiFi. Verified live: SSH/80/443 open, `/api/healthz` ok for API and database, home and admin 200, all five containers up, migration `16:false`, 27 active products, 36 orders, disk 14%, and **29 days of uptime**. Rewrote the false outage section in `HANDOFF-STATE.md`.
- **State:** Nothing was ever wrong with the server. It served continuously throughout, and the client browsed the storefront normally. The developer's ISP could not route to `82.29.184.79`; packets reached Hostinger's Manchester edge and were dropped. A second ISP failed the same way, which made an IP block look plausible. The panel-reported reboot never actually restarted the guest, as the 29-day uptime proves. Server is on `5543292`, exactly in sync with `origin/main`; no application deploy is pending.
- **Next:** Push the 2 laptop-only commits from the `Documents` checkout (`2aee56d`, `df7cbe3`) and these documentation commits. Then the still-unrun end-to-end order test on TEST keys, before any switch to live Paystack keys.
- **Blockers:** None technical. **Methodology note for both agents: every `nc`/`curl`/`ping`/`traceroute` runs on the developer's Mac and inherits their connection, so a failed probe is evidence about the network path and never about the server. When probes fail while the provider panel reports a healthy guest, believe the panel and verify off-network before declaring an outage.**

## 2026-08-31T22:36:21+01:00 · Claude · Staging server unreachable; production-readiness inventory

- **Did:** Ran a state check ahead of a production-readiness review. Found `82.29.184.79` completely unreachable (ICMP, 22, 80, 443 all timing out; `traceroute` terminating inside Hostinger's `man01` network; control request to github.com returning 200, so local connectivity is fine). Fast-forwarded the `Downloads/share` checkout from 9 commits behind to `origin/main` at `5543292` and ran the full Go suite, which passes including Codex's new `internal/tracking` package. Inventoried remaining production blockers: `sslip.io` hostname with no real domain, `sk_test_` Paystack, unwritten Legal/Terms copy, three parked navigation routes, mockup suit photography, unlicensed Gilmer font, ~5 test orders plus `OYE-1033` still in the staging database, and 22 unconfirmed placeholder prices.
- **State:** Server down and **not diagnosed from the inside** — the cause is inferred from external probes only. Codex deployed `a0232c8` to this same box earlier today and recorded live health 200, so the outage is hours old rather than gradual. Claude's initial guess of an expired Hostinger term was **wrong and is ruled out**: the panel shows the VPS `srv1875150.hstgr.cloud` **Running** with expiration **2026-09-02**. The VM is powered on and paid, so this is a wedged guest OS, a firewall change, or a network fault. Leading candidate is resource exhaustion from today's deploy on a 1 vCPU / 3.9 GB box whose Nuxt `vue-tsc` build is the known heavy step. Discovered two hazards that were not in the baton: **backups never leave the server** (`deploy/backup.sh` writes to the `postgres-backups` Docker volume on the same VPS, so a reclaim destroys the database and every backup together), and **two unsynchronized checkouts of this repo exist on the Mac** — `Documents/oyecouture` (Codex, holding 2 unpushed commits `2aee56d` and `df7cbe3`) and `Downloads/share/oyecouture` (Claude). Newest local dump remains `oyecouture-20260825T090120Z.sql.gz` from 25 Aug.
- **Next:** User checks the Hostinger panel for the VPS billing/power state. Restore or rebuild, then push the 2 stranded commits, then move backups off the box before anything else ships.
- **Blockers:** Everything server-side is blocked until the VPS is reachable. Nothing can be deployed, verified, or backed up in the meantime.

## 2026-08-31T11:42:24+01:00 · Codex · Product and admin detail hierarchy refined
- **Did:** Separated storefront product summary, order selection, and aftercare; grouped admin orders into overview, fulfilment, commercial, and conditional tailoring sections; segmented product, delivery, settings, and home-content editors with shared section cards and clearer action containment.
- **State:** Implementation is committed locally as `2aee56d`. Node 22 strict typecheck, Nuxt production build, `git diff --check`, and authenticated local browser QA pass across ready-to-wear/bespoke product pages and order/product/delivery/settings/home-content detail views. The local fixture was read only and remains unchanged.
- **Next:** Push and deploy only after fresh explicit approval.
- **Blockers:** None. The workstation disk remains very low on free space, but clearing the generated Nuxt cache provided enough room for validation.

## 2026-08-31T09:42:00+01:00 · Codex · Quantity scope deployed
- **Did:** Pushed `a0232c8` plus its handoff record, deployed through the standard Hostinger script, rebuilt both production images, and recreated API/web.
- **State:** Application images were built from release head `b81d37d`; all five containers are running, PostgreSQL is healthy, migration remains `16:false`, and live health returns 200. A forged quantity-2 bespoke request returned 400 `invalid_items`, and a direct count confirmed it created zero orders.
- **Next:** User verifies ready-to-wear `− / +` controls and fixed bespoke `Qty 1` on staging, then returns to authenticated tracking acceptance.
- **Blockers:** None for quantity. Authenticated tracking and real dispatch-email checks still require the user's staging session and an explicitly selected test order.

## 2026-08-31T03:05:00+01:00 · Codex · Quantity limited by fulfillment type
- **Did:** Kept the existing 1–10 quantity controls for sized/ready-to-wear products while fixing made-to-measure products at one unit across repeat adds, cart restoration, drawer presentation, and server validation. Added backend coverage proving multi-unit bespoke requests fail and multi-unit ready-to-wear requests retain correct totals.
- **State:** Committed locally as `a0232c8`; not pushed or deployed. Full Go tests, Node 22 strict typecheck, Nuxt production build, and `git diff --check` pass. A populated local browser interaction could not run because the local API was unavailable; the temporary web process was stopped.
- **Next:** User approves or requests changes, then commit/push/deploy if explicitly authorized.
- **Blockers:** None in code. Full browser interaction can be repeated when the local API is available.

## 2026-08-31T02:21:00+01:00 · Codex · Customer order tracking deployed
- **Did:** Committed the approved tracking implementation as `8bac779`, pushed `main`, and deployed through the standard Hostinger script. Both production images rebuilt and API/web were recreated; migration `000016` applied cleanly.
- **State:** All five containers are running, PostgreSQL is healthy, migration state is `16:false`, and public health returns 200. A live invalid token returns 404 with `private, no-store` and `noindex, nofollow`; the admin timeline route returns 401 without a session. The running API resolves the correct HTTPS `WEB_BASE_URL`, and SendByte configuration remains present. No Paystack outcome was selected and no real email was sent.
- **Next:** User verifies a populated timeline using their staging admin session, then explicitly selects a user-owned test order before dispatch-email acceptance. Resume the paused production audit afterward.
- **Blockers:** Codex does not have the admin password/session, so authenticated populated-timeline and real inbox acceptance remain user-owned checks.

## 2026-08-31T02:14:00+01:00 · Codex · Tracking release approved and verified locally
- **Did:** Received explicit commit/push/deploy approval, stopped the local review servers, and reran full Go tests, Node 22 strict typecheck, Nuxt production build, and `git diff --check` successfully.
- **State:** Tracking implementation is release-ready and still uncommitted. The disposable review order remains isolated in the developer's local database because deleting database records requires separate approval; no fixture data exists in Git or can reach staging.
- **Next:** Commit as `mikechibuzor`, push `main`, deploy migration `000016`, and run production smoke checks without touching Paystack order `OYE-1033` or sending a real email.
- **Blockers:** None for release. Populated authenticated timeline and dispatch-email acceptance remain user-session/test-order checks after rollout.

## 2026-08-31T01:29:00+01:00 · Codex · Local tracking review environment restored
- **Did:** Reconciled the baton and uncommitted tracking implementation, reran full Go tests and Nuxt typecheck, started the local API/web stack, and added disposable fixture `OYE-QA-TRACKING` with `.invalid` contact data so the populated admin and customer timelines can be reviewed without Paystack or real email.
- **State:** Local API runs on `:8080`, web on `:3010`, and the active local PostgreSQL instance is clean at migration `16:false`. The fixture response contains only approved public fields and returns `private, no-store` plus `noindex, nofollow`; no code was committed, pushed, or deployed.
- **Next:** User reviews both local views and explicitly approves or requests changes. Remove the fixture before commit/release.
- **Blockers:** None. Production rollout and any real dispatch email remain approval-gated.

## 2026-08-30T18:50:00+01:00 · Codex · Email-only customer order tracking implemented
- **Did:** Added migration `000016`, separate bespoke/ready-to-wear fulfillment stage machines, immutable admin-attributed timeline events, optimistic transition locking, secure dispatch-time tracking tokens, ETA/courier metadata, PII-minimized public tracking API/page, authenticated admin timeline/actions, branded milestone email templates, and a PostgreSQL-backed retry worker. Payment verification now creates the first tracking event atomically, while legacy order status remains a compatibility projection. Removed direct shipped/delivered admin actions so fulfillment cannot bypass the audited path.
- **State:** Local only; not committed, pushed, or deployed. Full Go tests, Node 22 typecheck, Nuxt production build, `git diff --check`, migration application, public 200/404 cache and noindex headers, and browser QA pass. A disposable order was advanced through the authenticated admin to Delivered; the public page updated immediately without customer/admin PII, success notification rendered, no horizontal overflow or alert errors appeared, migration state is clean at `16:false`, and all disposable order/admin data was removed.
- **Next:** User reviews locally, then explicitly approves commit/push/deploy. Production smoke checks must cover migration state, authenticated timeline load, invalid token behavior, and a dispatch email using a real user-owned test order. Resume the paused production audit afterward.
- **Blockers:** No implementation blocker. Production email delivery cannot be safely exercised locally because local email config is intentionally absent; the real dispatch email requires the existing production SendByte configuration and an explicitly approved user-owned test order. SMS remains out of scope because there is no approved live SMS provider.

## 2026-08-30T18:20:00+01:00 · Codex · Tracking architecture research and production-audit pause
- **Did:** Paused the production-readiness pass before completing its Paystack TEST transaction, mapped the existing payment, order-status, admin-action, and outbound-email paths, and reviewed current official guidance for customer order-status pages, fulfillment separation, secure object access, transactional outboxes, accessible status updates, and SMS consent.
- **State:** No feature code changed and no payment was submitted. The current system automatically attempts confirmation and receipt emails after verified payment and permits admin transitions from paid through in-progress, shipped, and delivered, but has no status-event history, public tracking page, carrier metadata, transition notifications, or automatic email retry worker. The tailored recommendation is documented in the baton and awaiting user approval.
- **Next:** Implement the approved email-first tracking MVP backend-first, then resume the paused production audit from the remaining payment, email, admin, responsive, accessibility, and failure-path checks.
- **Blockers:** SMS cannot ship in the MVP without a live provider and explicit customer consent; SendByte currently exposes transactional email while SMS remains listed as coming soon.

## 2026-08-30T12:30:02+01:00 · Codex · Measurement-aware checkout deployed
- **Did:** Deployed `origin/main` through the standard production script, rebuilding the API and Nuxt images and recreating both application services.
- **State:** Production runs exact Git head `eb686d6`, containing application feature `e128be8`. All five containers are up, PostgreSQL is healthy, migration state is clean at `15:false`, 42 GB remains free, and independent HTTPS checks return 200 for `/api/healthz` and `/confirm`.
- **Next:** User runs production acceptance for complete, partial, and absent measurement paths, then retries one authenticated admin product image upload.
- **Blockers:** No deployment blocker. Final acceptance still depends on the configured production WhatsApp and Paystack flows.

## 2026-08-29T16:26:00+01:00 · Codex · Measurement-aware bespoke payment routing
- **Did:** Made bespoke payment eligibility server-authoritative. The API derives the expected measurement profile from trusted product section/category data and returns `awaiting_payment` only when every required structured measurement is present under the matching profile. The storefront routes that status to the existing Paystack page, leaves measurement-free orders on the WhatsApp confirmation path, and blocks partial or invalid measurement submissions with explicit guidance.
- **State:** Committed as `e128be8` and pushed to `origin/main` on 2026-08-30; not deployed. Full Go tests, Node 22 typecheck, Nuxt production build, and browser QA pass. A complete native measurement set reached `/pay/{token}` with the correct order total; a partial set stayed on `/confirm`; a measurement-free order created `OYE-1003` and reached the WhatsApp confirmation state.
- **Next:** Deploy after explicit approval, then verify the three paths against production configuration.
- **Blockers:** Local site settings do not contain a WhatsApp number, so the success CTA is hidden locally; production must retain its configured WhatsApp number for the final continuation link.

## 2026-08-27T13:02:13+01:00 · Codex · Production Cloudinary upload configuration restored
- **Did:** Traced admin product image upload failures to the API's explicit `image_upload_unavailable` branch. Confirmed the required Cloudinary cloud name, API key, and API secret were configured locally but missing from production. After explicit user approval, securely synchronized only the Cloudinary settings into `.env.production`, preserved a server-side backup, recreated the API container, and removed temporary credential files without displaying or committing any values.
- **State:** Production `.env.production` and the running API container both report all four Cloudinary variables as set with values redacted. The API container is running without startup errors, and public `/api/healthz` returns HTTP 200 with both API and database healthy. No application-code release was required.
- **Next:** User retries an authenticated product image upload in the production admin to complete the functional acceptance check.
- **Blockers:** None at infrastructure level. The Cloudinary secret remains scheduled for rotation because it appeared in an earlier chat transcript.

## 2026-08-27T11:20:45+01:00 · Codex · Social links and duplicate catalogue cleanup
- **Did:** Added the client-provided Instagram, TikTok, and X destinations to the shared storefront footer with accessible external-link controls. Added reversible migration `000015` to archive the duplicate Diplomat product, retained Chairman Two-Piece, and aligned canonical catalogue/fallback data. A requested two-row collection and Show More control was implemented, then fully reverted when the client withdrew that request; all category filters still render their complete result set.
- **State:** Committed as `119ad8e`, pushed to `origin/main`, and deployed successfully. `GOCACHE=/private/tmp/oyecouture-go-cache go test ./...`, Node 22 typecheck/build, `git diff --check`, and the storefront/catalogue em-dash scan pass. Production is on `119ad8e`, migration state is `15:false`, all five containers are running with PostgreSQL healthy, and live health returns 200. Browser QA confirms all 27 active products render, no Show More control remains, Chairman appears once, Diplomat is absent, and the three social destinations are exact.
- **Next:** User runs production acceptance for the footer links and catalogue, then continues the payment/email and authenticated admin acceptance already tracked in this baton.
- **Blockers:** None.

## 2026-08-26T09:03:59+01:00 · Codex · Product media and universal delivery pricing
- **Did:** Removed synthetic product-card gallery slots so navigation renders only when a product has multiple real or approved fallback images. Made configured delivery-zone fees server-authoritative for bespoke/suit and ready-to-wear orders, added priced options plus delivery/total to `/confirm`, and extended WhatsApp confirmation with delivery and total. Removed em dashes from storefront/canonical catalogue copy, added migration `000014` for existing descriptions, normalized future admin-entered product descriptions, restored the missing shared header frame on Shipping, Returns, Legal, and Terms, and replaced the approximate bespoke-service symbol with the client's approved gold SVG.
- **State:** Committed as `56f339a`, pushed to `origin/main`, and deployed successfully. Full Go tests, Node 22 typecheck, Nuxt production build, `git diff --check`, and the em-dash source scan pass. Browser QA confirmed one-image suit cards have no arrows/dots, multi-image cards retain them, selecting Ibadan updates a mixed bespoke order from ₦1,630,000 to ₦1,633,500 with a ₦3,500 delivery line, the Shipping/Returns headings align exactly with the storefront logo, and the shared bespoke-service band renders the approved vector at the intended size. Staging is on `56f339a`, migration state is `14:false`, all five containers are running, PostgreSQL is healthy, and all checked live routes return 200. The known pre-existing hydration mismatch warning remains unchanged.
- **Next:** User runs production acceptance for gallery behavior, confirmation/WhatsApp, payment/email, and authenticated admin mutations.
- **Blockers:** None for this slice. International remains configured at ₦0 and should be updated in Admin Delivery when the client supplies its price.

## 2026-08-26T07:18:00+01:00 · Codex · Storefront and admin experience hardening
- **Did:** Executed the approved motion/interaction audit: added a reusable focus trap with Escape, scroll locking, and focus restoration for cart/admin dialogs; added an explicit hero pause control; strengthened focus rings and carousel/toast targets; replaced instant/global loading overlays with delayed progress plus geometry-preserving product skeletons; removed fizz/collision spectacle; enabled named product View Transitions and gallery crossfades; preloaded adjacent media and removed fake single-product gallery slots; capped search staggering; added the official `web-vitals` package for CLS/INP/LCP events; and rebuilt mobile admin navigation/tables into compact, usable dialogs/cards. Follow-up changed Wishlist **Add all to bag** into a successful-item move so the wishlist empties after the bag receives everything.
- **State:** Committed as `f36262f`, pushed to `origin/main`, and deployed successfully. Node `22.23.1` `npm run typecheck` and production `npm run build` pass; `git diff --check` is clean. Browser QA passed at 390px for Home, product shared navigation, a five-image Gold Agbada gallery with the next image visible within 80ms, single-image gallery cleanup, mobile admin navigation, product table cards/editor drawer, focus placement, and Wishlist → Add all → empty Wishlist with an incremented bag. All five production containers are running, PostgreSQL is healthy, and live checks returned 200 for health, core storefront, measurements/confirmation, and admin routes. `npm audit --omit=dev` could not reach the npm registry from the sandbox; deployment's dependency install reported two high-severity audit findings for separate review.
- **Next:** User runs production acceptance plus the payment/email tests they own; then resolve client/legal/media launch blockers and triage the two dependency findings.
- **Blockers:** None for staging deployment. Final launch still depends on client/legal/pricing/media/licensing decisions already listed in `HANDOFF-STATE.md`.

## 2026-08-11T22:38:00+01:00 · Codex · Admin URL state and control polish
- **Did:** Made Orders and Products persist admin table state in readable query params (`page`, `limit`, `period`, `q`, `status`, `kind`, `from`, `to`, `section`, `active`) so reload/back/forward preserve the current desk view. Added first-mount query normalization, stale-page clamping when a URL points past the available results, custom admin checkbox styling, stronger focus-visible rings, primary save-button treatment, and number-input spinner removal for admin form fields.
- **State:** Local only; not committed, pushed, or deployed yet. Verified with `git diff --check`, Node 22 `npm run typecheck`, and Node 22 `NUXT_IGNORE_LOCK=1 npm run build`.
- **Next:** User visually tests admin query persistence and control polish locally; then commit/push/deploy the admin modernization bundle if approved.
- **Blockers:** None for this slice.

## 2026-08-11T21:36:00+01:00 · Codex · Admin dashboard modernization pass
- **Did:** Modernized the shared admin dashboard visual system without changing data contracts: dark premium sidebar rail, stronger page headers, elevated metric cards, grouped toolbar/search/filter surfaces, cleaner table headers/row hover/focus affordances, softer status pills, upgraded form field/select focus states, blurred drawer backdrop, smoother drawer/popover entrances, and reduced-motion fallbacks. Also added rows-per-page controls to the paginated Orders and Products desks (`10/25/50/100`), with API `limit`/`offset` kept in sync.
- **State:** Local only; not committed, pushed, or deployed yet. Verified with `git diff --check`, Node 22 `npm run typecheck`, Node 22 `NUXT_IGNORE_LOCK=1 npm run build`, and in-app browser screenshots for `/admin/home-content`, `/admin/products`, and the product edit drawer. Browser-smoked Products rows-per-page from 25 to 10 rows; Orders local DB currently has no rows, so its footer path was compile-verified rather than interaction-verified.
- **Next:** User visually tests core admin routes locally; then decide whether to tune the visual direction further or commit/push/deploy.
- **Blockers:** None for this styling pass.

## 2026-08-11T18:03:48+01:00 · Codex · Home content admin first slice
- **Did:** Added JSON-backed home content using the existing `settings` table: public `GET /home-content`, authenticated `GET/PATCH /admin/home-content`, a new `/admin/home-content` page, active sidebar navigation, `useHomeContent`, and storefront homepage bindings for hero copy, CTA link/label, category section copy, and category tile label/heading/image URL/alt/link fields. Static constants remain the fallback if the endpoint is unavailable or unset.
- **State:** User tested locally and approved. Committed/pushed/deployed as `925ad9b`. Verified with `GOCACHE=/private/tmp/oyecouture-go-cache go test ./...`, Node 22 `npm run typecheck`, Node 22 `NUXT_IGNORE_LOCK=1 npm run build`, `git diff --check`, local smoke checks, and live staging checks: `/api/healthz` 200, `/api/home-content` 200, `/api/admin/home-content` 401 without a session, `/admin/home-content` 200/noindex.
- **Next:** Continue into homepage gallery image upload/reorder when work resumes.
- **Blockers:** None for text/link content. Homepage gallery upload/reorder remains a separate slice.

## 2026-08-11T16:45:45+01:00 · Codex · Delivery areas drawer flow
- **Did:** Added authenticated admin delivery-area creation with `POST /admin/delivery-zones`, duplicate-name handling, and a product-style `/admin/delivery` drawer flow. Delivery areas now render as table rows; clicking a row opens the drawer with fee editing, checkout visibility, and hide/show CTAs. Creating or saving closes the drawer and fires the shared top-right success notification; failures keep the drawer open.
- **State:** Committed/pushed as `c325b92` and deployed in the `925ad9b` staging rollout. Verified locally with `GOCACHE=/private/tmp/oyecouture-go-cache go test ./...`, Node 22 `npm run typecheck`, Node 22 `NUXT_IGNORE_LOCK=1 npm run build`, and `git diff --check`. User tested add/hide/show locally and approved.
- **Next:** Continue into editable home/content controls.
- **Blockers:** None for this slice.

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
