# 01 · Project Brief

## The client & the relationship

- **Client:** Oyecouture Bespoke — a premium Nigerian bespoke fashion house. Menswear-led (suits, agbada, aso-oke, kaftans) with some womenswear casuals.
- This is a **repeat client** (we built a site for them before). They asked to **"start small"** — hence a tightly-scoped Phase 1. They already **own the domain**.
- The client thinks in **naira**, wants a lean first phase, and negotiated the price down. They have chosen **not to pay monthly running costs** for now (see Architecture — we launch on free tiers at ₦0/month).
- Deliver quickly and reliably. Cadence is ~daily; the commercial promise is a live demo by **31 July 2026** and launch by **end of August 2026**.

## Commercials (context for scope decisions, not for the code)

- **Phase 1 fixed price: ₦450,000**, milestones 40% / 40% / 20% (on acceptance / on demo / on launch). First 40% is **paid** — build is a go.
- **Running costs: ₦0/month at launch** (free-tier hosting; client declined ongoing costs).
- Optional maintenance retainer: ₦25,000/month.
- Full scope & pricing verbatim: [reference-docs/QUOTATION.md](reference-docs/QUOTATION.md).

## What Phase 1 includes (build all of this)

1. **Storefront** — Home/Shop, 4 category pages (Suits, Trad Attires, Casuals Men, Casuals Women), About. Filterable + sortable product grid. Product detail pages with image gallery and size selection (RTW only). Fully responsive, mobile-first.
2. **Product search** — instant search from the nav across name/category/description. (Postgres full-text/trigram is enough at ~28 products; no external search service.)
3. **Cart** — drawer per the design; persisted for guests (localStorage). Branches by product type (see below).
4. **Bespoke order flow (WhatsApp-first)** — cart → "Confirm order on WhatsApp"; order saved without payment; optional structured measurement form; deep-link to WhatsApp carrying the order number. Owner later shares a secure payment link.
5. **Ready-to-wear checkout** — contact + delivery address + delivery zone/fee + pay via Nomba now.
6. **Online payments** — Nomba primary, behind a provider-agnostic interface with a Paystack fallback + health-check failover.
7. **Order confirmations** — automated email to customer and owner on every order and every payment; WhatsApp hand-off carries the order number.
8. **Customer-care pages** — Size Guide, Shipping & Delivery, Returns & Exchanges, Contact, Book-a-Fitting (form → email + WhatsApp), Terms, Privacy. (Copy supplied by client; we template.)
9. **Launch** — deploy on the client domain, HTTPS, SEO (indexable product pages, OG tags, sitemap), analytics, error monitoring.

## The two product types (the core rule — repeated because it matters)

- **`hasSizes: true` (RTW)** → size selector S/M/L/XL on the product page → normal cart → **pay-now checkout**.
- **`hasSizes: false` (bespoke)** → no size selector; product page shows the note *"Sizing, fabric & delivery details are confirmed with us on WhatsApp after checkout"* → cart CTA becomes **"Confirm order on WhatsApp"** → order saved as `pending_confirmation`, **no payment taken** → customer optionally submits measurements → WhatsApp opens with a prefilled order-number message → owner agrees details, then shares the order's `/pay/<token>` link → customer pays.
- **Mixed cart** (bespoke + RTW together): treat the whole order as WhatsApp-first (one conversation, one payment link). This is an assumption — revisit if it confuses users.

## The "Ghost Admin" — INTERNAL, confidential, deliberate

**Read carefully.** The client said they do **not** want an admin layer yet ("start small"). We are building one anyway, quietly, because:

- Products, prices, delivery fees, the announcement bar, and the payment-provider toggle must not require code deploys.
- The bespoke flow **needs** an operator surface from day one: to view an order's measurements, **copy its payment link**, and move its status.
- Building it now on a fresh schema costs a fraction of retrofitting it later.

**Rules for the ghost admin:**
- It exists at a route group like `/admin`, behind auth (single admin user), **not linked anywhere** in the storefront.
- It is **not billed** in Phase 1 and **not revealed** to the client until they ask for self-service. When they do, it becomes the pre-priced "Admin & Content Management" add-on (₦300,000: reveal + training + polish + roles/audit log).
- Until then, **we** are the only operators; the client requests catalogue changes through us.
- Build it. Just don't surface it to the client or put it in client-facing materials.

Scope of the ghost admin: product CRUD + image upload (R2 presigned) + activate/deactivate + reorder; order list/detail with status transitions, measurement view, **copy payment link**, resend confirmation email; settings for delivery-zone fees, WhatsApp number, announcement-bar text, and Nomba⇄Paystack toggle.

## Explicitly OUT of Phase 1 (do not build)

Accounts/login, wishlist persistence, saved measurement profiles, discount codes, reviews, multi-currency display, inventory/stock counts (this is a made-to-order house — no stock tracking), blog/lookbook, automated *outbound* WhatsApp messaging (that needs the WhatsApp Business API — Phase 3), and the admin *handover* (the admin is built but not revealed).

The wishlist heart and account icon appear in the design — render them, but they can be visually present/inert or deferred; do not build persistence in Phase 1. Confirm the exact treatment as you go.

## Later phases (so you build Phase 1 without painting us into a corner)

- **Phase 2:** admin handover (reveal the ghost admin), customer accounts + **saved measurement profiles** ("measure once, reorder forever" — this is our differentiator; no local competitor has it), discounts, lookbook/Instagram/blog, back-in-stock notifications, real international delivery rates.
- **Phase 3:** WhatsApp Business API + AI customer-service agent (FAQ, order-status, payment/delivery confirmations, abandoned-cart recovery), size recommendation from height/weight, multi-currency, fitting-calendar booking.

Design the data model and payment layer so these slot in later without a rewrite.

## What the client still owes us (open items — surface these to the user)

Product photography (the design uses placeholders for most items), copy for the customer-care pages, the WhatsApp business number, Nomba account credentials (client suggested Nomba), domain DNS access, delivery-fee amounts per zone, and the refund/returns policy text.
