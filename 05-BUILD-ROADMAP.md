# 05 · Build Roadmap

Cadence ~daily. The hard commitment is a **working demo on a live link by 31 July 2026**, launch by **end of August 2026**. Build in this order — each step leaves something runnable.

## Milestone A — Foundations & catalogue (Week 1, by ~26 Jul)
1. Scaffold repo: Nuxt 3 (`web/`) + Go API (`api/`) + Postgres + `.env.example`. Get both dev servers talking.
2. Migrations for `products`, `product_images`, `settings`, `delivery_zones` (see 04-DATA-MODEL.md). **Seed from `catalogue/product-data.js`.**
3. Go endpoints: `GET /products` (with `?section=` + `?q=` search), `GET /products/:slug`.
4. Nuxt storefront shell from **02-DESIGN-SPEC**: AnnouncementBar, NavBar (+ dropdowns + mobile menu), Footer, design tokens (colors + Jost placeholder font behind one CSS var).
5. Home/Shop page: hero, category tiles, filter chips + sort, ProductCard grid wired to the API. Category pages (Suits, Trad Attires, Casuals Men/Women). Product detail page with gallery, size selector (RTW only) + bespoke note.
6. Search overlay (instant results from `?q=`).
7. Deploy to the Oracle free-tier VM as **staging on a live link** — even this early, so the client sees momentum.

## Milestone B — Cart, orders, payments, emails → THE DEMO (Week 2, by 31 Jul)
8. `useCart` composable (localStorage) + CartDrawer per design, with the **branching CTA**: all-RTW → "CHECKOUT"; any bespoke → "CONFIRM ORDER ON WHATSAPP".
9. Order creation API + `orders`/`order_items`/`measurements` tables + order-number + `pay_token` generation.
10. **Bespoke flow:** confirm → create `pending_confirmation` order → optional MeasurementForm → open WhatsApp deep link (`wa.me/<number>?text=Order OYE-… …`).
11. **RTW checkout:** contact + address + delivery zone/fee → create `awaiting_payment` order.
12. **Payments:** the `PaymentProvider` interface + **Nomba adapter (sandbox)**; `InitCheckout` → redirect; webhook + `VerifyTransaction`; mark `paid` only after verify. `/pay/<token>` page for bespoke.
13. **Emails (SendByte):** customer order confirmation + payment receipt after verified payment; owner notifications can follow. River worker/outbox retry can be layered in.
14. Order-confirmation page (with WhatsApp CTA for bespoke).
15. Redeploy staging → **this is the demo delivered on 31 Jul** (storefront + cart + WhatsApp confirm + sandbox checkout).

## Milestone C — Hardening & the ghost admin (Weeks 3–4)
16. **Paystack fallback** adapter + health-check failover + provider toggle in settings.
17. Payment **reconciliation** job (re-verify pending; alert owner on mismatch). Webhook idempotency.
18. **Ghost admin** (`/admin`, unlinked, auth): product CRUD + R2 image upload + activate/reorder; order list/detail with status transitions, measurement view, **copy payment link**, resend email; settings (delivery fees, WhatsApp number, announcement text, provider toggle).
19. Customer-care pages (Size Guide, Shipping, Returns, Contact, Book-a-Fitting form → email+WhatsApp, Terms, Privacy) — template now, fill client copy when supplied.
20. SEO (sitemap, OG tags, indexable product pages), analytics, Sentry. Mobile QA. Lighthouse ≥ 90.
21. Live payment tests with real Nomba credentials (small real amounts, then refund).

## Milestone D — Content & launch (Weeks 5–6, target end of Aug)
22. Load real product photography (R2) + real page copy from the client; flip `active=true` as photos land.
23. Resolve the **Gilmer font licensing** decision (license it or keep the free substitute).
24. UAT with the client; fix-list; final deploy to the **client domain** with HTTPS. Hand over the 2-week support window.

## Cross-cutting reminders
- Money in **kobo**; never floats. `paid` only after server-side verify. Mobile-first. Guest-only. Idempotent webhooks.
- Keep the admin **unlinked and unmentioned** to the client (ghost).
- Chase the **open items** early (Nomba creds, WhatsApp number, delivery fees, photos, care-page copy) — the timeline slips on missing content, not on code. See 01-PROJECT-BRIEF.
