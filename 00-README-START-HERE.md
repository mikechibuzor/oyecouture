# Oyecouture — Build Handoff (START HERE)

**You are picking up a real, paid client project.** The first milestone payment has landed and the build is starting now. This folder contains everything you need to understand what has been decided and to start writing code — no prior conversation required.

## What this project is (30 seconds)

**Oyecouture Bespoke** is a premium Nigerian fashion house (bespoke suits, traditional attire, casuals). We are building their **first e-commerce storefront** on a domain they already own. The visual design is already done (7 pages, provided as self-contained HTML in `design-source/`). We are turning that design into a working store with search, cart, online payments (Nomba), a WhatsApp-first flow for bespoke items, and email + WhatsApp order confirmations.

**Target: a working demo on a live link by 31 July 2026; launch by end of August 2026.** Build cadence is roughly daily.

## The one thing you must not get wrong

The catalogue has **two kinds of products**, and they behave completely differently at checkout. This is the heart of the whole build:

| Type | Flag in data | Examples | Flow |
|---|---|---|---|
| **Ready-to-wear (RTW)** | `hasSizes: true` | shirts, polos, chinos, knitwear | pick size → cart → **checkout → pay now** → ship |
| **Bespoke** | `hasSizes: false` | suits, agbada, aso-oke, kaftans | cart → **"Confirm order on WhatsApp"** (order saved, NO payment yet) → measurements/fabric/delivery agreed on WhatsApp → owner shares the order's **secure payment link** → customer pays → production |

Bespoke customers **pay only after their measurements are known.** Get this branching right and the rest follows.

## Read these in order

1. **[01-PROJECT-BRIEF.md](01-PROJECT-BRIEF.md)** — full context: client, scope, phase boundaries, what's explicitly in/out, the "ghost admin" (read this — it's a deliberate, confidential decision).
2. **[02-DESIGN-SPEC.md](02-DESIGN-SPEC.md)** — the visual system extracted from the design: colors, fonts, every page, every component. Pair it with the live files in `design-source/`.
3. **[03-ARCHITECTURE.md](03-ARCHITECTURE.md)** — stack (Nuxt 3 + Go + Postgres), the payment-provider abstraction, folder layout, hosting.
4. **[04-DATA-MODEL.md](04-DATA-MODEL.md)** — tables, order state machine, money handling.
5. **[05-BUILD-ROADMAP.md](05-BUILD-ROADMAP.md)** — the actual build order, step by step, toward the 31 July demo.
6. **[06-SETUP-AND-ENV.md](06-SETUP-AND-ENV.md)** — accounts, env vars, local dev, free-tier deploy.
7. **[07-MOTION-AND-MICROINTERACTIONS.md](07-MOTION-AND-MICROINTERACTIONS.md)** — the motion system: tokens, per-component animation spec, accessibility and performance rules. **Read before adding any animation.**

## What's in the other folders

- **`design-source/`** — the 7 approved design pages as self-contained `.dc.html` files. **Open any of them directly in a browser** to see the exact target (images are embedded; no server needed). The product page is `Oyecouture Product.dc.html?id=<slug>` (e.g. `?id=agbada`, `?id=linen`). See [design-source/HOW-TO-VIEW.md](design-source/HOW-TO-VIEW.md).
- **`catalogue/`** — `product-data.js` is the **real, canonical catalogue** used by the design (27 products, prices, the `hasSizes` flag). `CATALOGUE.md` is the human-readable version. Seed the database from this.
- **`reference-docs/`** — the internal PRD/FRD and the client quotation, verbatim, if you want the primary source.

## Ground rules carried over from planning

- **Money is stored in kobo (integer), never floats.** Prices in `product-data.js` are already in naira (e.g. `850000` = ₦850,000) — multiply by 100 for storage.
- **Never trust a payment redirect** — a Nomba order is only `paid` after a server-side verify / webhook.
- **Mobile-first.** Nigerian traffic is majority mobile; the design must be perfect on a phone.
- **Guest checkout only** in phase 1. No accounts, no login, no saved profiles yet.
- Branch/commit convention: work under `mikechibuzor`; **do not add Co-Authored-By trailers.** Ticket-style branch names if a tracker is used.
