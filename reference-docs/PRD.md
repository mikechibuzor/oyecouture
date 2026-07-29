# Oyecouture — Product Requirements Document (PRD/FRD)

**Version:** 1.1 · **Date:** 19 July 2026 · **Status:** Internal (engineering copy — includes ghost-admin scope not shared with client)

---

## 1. Overview

Oyecouture is a premium Nigerian bespoke fashion house (suits, traditional attires, casuals) launching its first e-commerce storefront. The visual design is complete (7 Claude-designed pages + cart drawer). Phase 1 delivers a working store on the client's existing domain: catalogue, search, cart, a WhatsApp-first confirmation flow for bespoke garments, online payment (Nomba) via checkout or secure payment link, and order confirmations by email and WhatsApp.

**Brand:** Oyecouture Bespoke — "Your clothing speaks long before you do"
**Catalogue at launch:** ~28 products across 4 sections (Suits, Trad Attires, Casuals Men, Casuals Women), ₦165,000–₦950,000
**Domain:** already owned by client
**Client-visible scope:** storefront + checkout only. **Internal scope:** ghost admin (see §7).

## 2. Product model — the key insight

The catalogue splits into two fulfilment types, already encoded in the design (`hasSizes` flag):

| Type | Examples | Order flow |
|---|---|---|
| **Ready-to-wear** (`hasSizes: true`) | Shirts, polos, chinos, knitwear | Pick size S–XL → cart → checkout → **pay now** → ship |
| **Bespoke** (`hasSizes: false`) | Suits, agbada, aso-oke, kaftans | Cart → **"Confirm order on WhatsApp"** (order saved, no payment) → measurements & fabric agreed on WhatsApp → client shares the order's **secure payment link** → customer pays → production |

**Payment only happens after measurements are known** for bespoke items. Customers who already know their measurements can submit them in a structured form before the WhatsApp hand-off, so the client receives them with the order.

**Mixed carts** (bespoke + RTW together): the whole order takes the WhatsApp-confirm route — one conversation, one payment link. (Assumption; revisit if it confuses users.)

## 3. Users

- **The Groom / Ceremony buyer** — high-ticket trad (agbada ₦950k, wedding ensembles); needs trust signals, WhatsApp contact, fittings
- **The Executive** — suits; values speed, polish, "book a fitting"
- **The Casual shopper** — ready-to-wear under ₦100k; wants a normal add-to-cart → pay → deliver flow
- **Diaspora buyer** — "worldwide delivery" banner; pays by card (Nomba supports NGN & USD cards), delivery quoted via WhatsApp
- **Store owner (client)** — confirms bespoke orders on WhatsApp, sends payment links, receives order emails; later manages catalogue via admin (phase 2 reveal)

## 4. Phase 1 — MVP functional requirements

### 4.1 Storefront (from the Claude design, made real)
- **FR-1** Home/Shop page: announcement bar, hero, category tiles, "All Pieces" grid with category filter chips (All / Suits / Aso-Oke / Casuals / Kaftans), sort (Featured / Price ↑↓ / Newest), product cards with image carousel + tag badges (New / Bespoke)
- **FR-2** Category pages ×4: Suits, Trad Attires, Casuals Men, Casuals Women
- **FR-3** Product detail: gallery + thumbnails, breadcrumbs, price, description, size selector (RTW only), Add to Bag, bespoke WhatsApp note
- **FR-4** **Product search**: nav search icon opens search overlay; instant results by name/category/description as you type (Postgres full-text/trigram — no external search service needed at this catalogue size)
- **FR-5** About page (4 sections per design)
- **FR-6** Footer utility pages: Size Guide, Shipping & Delivery, Returns & Exchanges, Contact Us, Book a Fitting (form → email + WhatsApp link), Legal/Privacy, Terms — *content supplied by client, we template*
- **FR-7** Responsive (mobile-first — Nigerian traffic is majority mobile), SEO meta/OG tags per product, sitemap

### 4.2 Cart & orders
- **FR-8** Cart drawer per design: line items with thumbnail, qty stepper, remove, subtotal, WhatsApp note; persisted in localStorage (guest)
- **FR-9** Cart CTA branches:
  - all-RTW cart → **Checkout** (pay now)
  - any bespoke item → **"Confirm order on WhatsApp"** per design
- **FR-10** Bespoke confirm flow: capture contact (name, email, phone) + optional structured measurement form (neck, chest, waist, hips, sleeve, inseam, length, height, weight, notes) → order saved as `pending_confirmation` → deep link opens WhatsApp with prefilled message (`Order OYE-1042 — The Emir Kaftan ×1 — I'd like to confirm my order`)
- **FR-11** RTW checkout (guest only): contact, delivery address, delivery zone → fee, order summary → pay
- **FR-12** Delivery fees: configurable zones — Lagos / Rest of Nigeria / International ("delivery arranged via WhatsApp after payment")

### 4.3 Payments (provider-agnostic by design)
- **FR-13** Go `PaymentProvider` interface: `InitCheckout`, `VerifyTransaction`, `HandleWebhook`. Adapters: **Nomba (primary)**, **Paystack (fallback)**. Provider selected by config + health check; automatic failover if primary init fails
- **FR-14** Two entry points, one engine:
  - RTW checkout → Nomba hosted checkout (`POST /v1/checkout/order` → redirect)
  - Bespoke → every order gets a **secure payment page** `oyecouture.com/pay/<token>`; client shares it on WhatsApp once measurements are agreed; page shows order summary + delivery fee and launches the same hosted checkout
- **FR-15** Webhook + server-side verify before any order is marked paid (never trust redirect alone); webhooks idempotent, signatures verified
- **FR-16** Order states: `pending_confirmation (bespoke) → awaiting_payment → paid → in_progress → shipped → delivered / cancelled / refunded` (RTW orders start at `awaiting_payment`)
- **FR-17** Payment reconciliation job (River or cron): re-verify pending payments; alert owner on mismatch

### 4.4 Notifications (email + WhatsApp at every step)
- **FR-18** On bespoke order confirmation: customer email ("we've received your order, continue on WhatsApp") + owner email with items, contact and measurements (if provided) + owner wa.me quick-reply link
- **FR-19** On payment: customer receipt email + owner "paid" email — for **every** order
- **FR-20** WhatsApp side in phase 1: conversation-thread based — the bespoke flow *opens* the thread from the customer's side (free, no API), and the client confirms/updates there; all prefilled links carry the order number. *Automated outbound WhatsApp messages require the WhatsApp Business API (Meta per-conversation billing) — scoped as the phase-3 upgrade, the site is built so it can be slotted in*

### 4.5 Non-functional
- **NFR-1** Stack: **Nuxt 3 (SSR) + Go API + PostgreSQL**; River for async jobs (emails, payment verification); images on **Cloudflare R2** (zero egress; Go backend generates WebP size variants at upload)
- **NFR-2** Lighthouse ≥ 90 performance on product pages; images lazy-loaded, WebP
- **NFR-3** All money stored in kobo (int64); no floats
- **NFR-4** Deploy: **₦0/month at launch** — Oracle Cloud Always Free ARM VM running Nuxt + Go + PostgreSQL (self-hosted on the VM, or Neon free tier); always-on so no cold starts; R2 + Resend free tiers; HTTPS on client domain. Client declined running costs, so keep hosting on Uzor & Co. free-tier accounts for now; migrate to paid (~₦15–30k/mo) only when traffic outgrows the free tiers, quoted at that point. Avoid Render free Postgres (expires after 90 days).
- **NFR-5** Basic analytics (Plausible or GA4) + error tracking (Sentry)

## 5. Data model (phase 1)

```
products(id, slug, name, section, category, price_kobo, description, tag, has_sizes, active, sort_order)
product_images(id, product_id, url, position)
orders(id, order_no, type, status, customer_name, email, phone, address, delivery_zone,
       delivery_fee_kobo, subtotal_kobo, total_kobo, pay_token, payment_provider, payment_ref, created_at)
order_items(id, order_id, product_id, name_snapshot, price_kobo_snapshot, size, qty)
measurements(id, order_id, neck, chest, waist, hips, sleeve, inseam, length, height, weight, notes)
payment_events(id, order_id, provider, event_type, payload_json, created_at)
delivery_zones(id, name, fee_kobo, active)
admin_users(id, email, password_hash, role)        -- ghost admin
settings(key, value)                                -- announcement bar, WhatsApp number, provider toggle
```

## 6. Explicitly OUT of phase 1
Accounts/login, wishlist persistence, saved measurement profiles, discount codes, reviews, multi-currency display, inventory/stock counts (made-to-order house — no stock tracking yet), blog/lookbook, automated WhatsApp messaging (Business API), admin handover.

## 7. Ghost admin (built in phase 1, **not billed, not revealed**)

Rationale: products, prices, delivery fees and the announcement bar must not require code deploys; building it now on a fresh schema costs a fraction of retrofitting later. The bespoke flow also *needs* an operator surface (payment links, status moves) from day one.

- Separate route group `/admin` behind auth (single admin user), not linked anywhere
- Product CRUD + image upload (R2 presigned upload, variants generated server-side) + activate/deactivate + reorder
- Order list/detail, status transitions, measurement view, **copy payment link** button, resend emails
- Delivery zone fees, WhatsApp number, announcement bar text, payment provider toggle (Nomba ⇄ Paystack)
- **Until reveal:** we are the only operators; client requests changes via us. When client asks for self-service, we quote "Admin & Content Management" as a phase-2 line item (training + handover + any polish), pre-listed in the quotation's add-ons.

## 8. Phase 2 — candidates (quote separately when requested)
1. **Admin handover** — reveal + training + audit log + roles (the pre-built ghost admin, billed at reveal)
2. **Customer accounts** — login (email/OTP), order history, **saved measurement profile** (Indochino-style "measure once, reorder forever"), wishlist persistence
3. **Discounts & promos** — codes, sales pricing
4. **Content** — lookbook/collections page, Instagram feed embed, blog (see quotation notes)
5. **Back-in-stock / new-drop notifications** (email list + WhatsApp broadcast)
6. **International delivery rates** — real zone pricing instead of WhatsApp quote

## 9. Phase 3 — WhatsApp AI agent & growth (nice-to-have backlog)
- **WhatsApp Business API** (via BSP; Meta bills per conversation — utility msgs ≈ <$0.01, marketing ≈ $0.02–0.22) replacing the personal wa.me number
- AI customer-service agent: FAQ, order-status lookup, measurement guidance, payment confirmation, delivery updates
- Automated order-lifecycle WhatsApp messages (order received / payment received / in production / shipped)
- Abandoned-cart recovery + window-shopper re-engagement via WhatsApp (typical e-comm ROI leader)
- Post-delivery follow-up & review collection
- Size-recommendation from height/weight (Hockerty-style algorithm)
- Multi-currency pricing, virtual fitting booking with calendar

## 10. Competitor/market notes (research summary)
- **Deji & Kola** (dejiandkola.com): closest local comp — custom + MTM menswear, global delivery, accounts, occasion-based categories. No visible measurement-profile feature → our phase-2 profile is a differentiator
- **ARIO, HoseDavins, Wyseman**: Nigerian tailors selling online; WhatsApp is the universal measurement/fitting channel — validates confirm-on-WhatsApp as the *primary* bespoke flow, not a workaround
- **Indochino / Hockerty** (global MTM): saved measurement profiles ("measure once"), guided video measurement, algorithmic fit from height/weight, fit guarantee — the phase-2/3 playbook
- **Payments**: Paystack & Flutterwave both 1.5% + ₦100 local (capped ₦2,000); Nomba comparable with cards + transfer + USSD and NGN/USD card support. Provider abstraction (FR-13) makes the choice reversible
- **Fashion e-comm table stakes 2026**: search + size guides (52% of apparel returns are size-related) in phase 1; wishlist, back-in-stock alerts, lookbooks, Instagram shopping mapped into phases 2–3

## 11. Timeline (working ~daily)
- **Week 1 (by 26 Jul):** repo + schema + product/catalogue API + storefront pages wired to real API + prod infra live on domain (staging)
- **Week 2 (by 31 Jul):** cart + bespoke WhatsApp confirm flow + RTW checkout + Nomba sandbox + payment-link page + emails → **workable demo on prod for client**
- **Weeks 3–4:** Paystack fallback + reconciliation, measurement form, search, utility pages, ghost admin, polish, mobile QA, live payment tests
- **Weeks 5–6:** content load (real photos/copy from client), UAT, launch. **Target launch: end of August 2026**

## 12. Open items (client to supply)
Product photography (design has placeholders for most items), copy for utility pages, WhatsApp business number, Nomba account credentials, domain DNS access, delivery fee amounts per zone, refund/returns policy terms.
