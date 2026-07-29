# 02 · Design Spec

The approved design is the source of truth. The 7 files in `design-source/` are self-contained — **open them in a browser to see the exact target** (all images are embedded as base64; no server needed). This document is the written extraction so you can rebuild faithfully without re-parsing 2.3 MB bundles.

> The `.dc.html` files are bundled/compiled design exports: a small HTML shell plus one large inline `<script>` that renders the page and holds base64 images. They are **reference only** — do not try to reuse their compiled JS. Rebuild the markup cleanly in Nuxt from this spec.

## Brand & tone

Oyecouture Bespoke — luxury tailoring. Tagline motifs: *"Tailor your ambition with effortless precision"*, *"Your clothing speaks long before you do"*, *"Be exclusive"*. Visual language: high-contrast charcoal + cream, champagne-gold accents, big confident uppercase display type, generous whitespace, editorial product photography. Feels like a fashion house, not a mass retailer.

## Color tokens (extracted from the bundle)

| Token | Hex | Use |
|---|---|---|
| `charcoal` | `#1a1a1a` | announcement bar, hero bg, footer bg, primary buttons, headings, body text on light |
| `cream` | `#faf9f5` | page background |
| `cream-2` | `#f7f3ee` | alternate section background |
| `sand` | `#ece7e1` / `#e6e1da` / `#e8e2d6` | product-image placeholders, soft cards |
| `champagne` | `#e0cf9c` | primary CTA fill (e.g. "Discover our collection"), hero accent |
| `gold` | `#a98a3d` | active nav link, eyebrow labels, links, "Add to bag" secondary |
| `gold-bright` | `#caa44a` | hover/interactive gold |
| `taupe` | `#8a827a` / `#a59c92` | muted/secondary text |
| `brown-ink` | `#4a443e` / `#3a352f` | dark secondary text |
| `white` | `#ffffff` | nav bar, cards |
| accent jewel tones | `#5a2230` (maroon), `#1f3a5f` (navy), `#2b2f3a` | occasional tag/category accents — verify per use |

Define these as CSS variables / Tailwind theme colors. Support the whole thing on mobile first.

## Typography

- **Display/headings:** `Gilmer` (weights: Heavy, Bold) — geometric sans, **UPPERCASE**, tight tracking, very large on hero.
- **Body:** `Gilmer` (Regular, Medium), sentence case, comfortable line-height.
- **Eyebrows/labels/nav:** uppercase, **wide letter-spacing** (~0.15–0.25em), small size, gold or taupe.
- ⚠️ **Gilmer is a commercial font (Fontfabric).** It is embedded in the design bundle. For the real build you need a proper license, OR substitute a close free geometric sans. Good free stand-ins: **Jost**, **Poppins**, or **Sora**. **Flag this to the client/user as a licensing decision before launch.** Until resolved, develop with Jost as the placeholder and keep the font behind one CSS variable so it's a one-line swap.

## Global layout

### Announcement bar (top, full-width, charcoal)
Centered, uppercase, gold text, small, wide tracking:
`COMPLIMENTARY TAILORING CONSULTATION · WORLDWIDE DELIVERY · BE EXCLUSIVE`
(Text must be editable via the ghost admin `settings`.)

### Header / nav (white, sticky)
- Left: **Oyecouture** logo — circular emblem mark + wordmark "Oyecouture" with "BESPOKE" spaced beneath.
- Center links: `HOME`, `TRAD ATTIRES`, `SUITS`, `CASUALS ▾` (dropdown: Men / Women), `COLLECTION ▾` (dropdown). Active link is gold.
- Right icons: **search**, **account**, **wishlist (heart)**, **cart (bag)**. Cart opens the drawer; search opens the search overlay. Account/wishlist are present in the design but not functional in Phase 1 (render, defer behavior).
- Mobile: collapses to a hamburger; icons remain accessible.

### Footer (charcoal)
Four columns + a brand blurb:
- **STAY CONNECTED** — blurb "Oyecouture Bespoke — where confidence, simplicity, and class become clothing." + social icons: Instagram, X, TikTok, **WhatsApp**.
- **SHOP** — Trad Attires, Suits, Casuals.
- **ABOUT** — Our Story, The Oyecouture Symbol, Craftsmanship, Mission and Vision.
- **CUSTOMER CARE** — Size Guide, Shipping & Delivery, Returns & Exchanges, Book a Fitting, Contact Us.
- Bottom bar: `© 2026 Oyecouture. All rights reserved.` + `LEGAL & PRIVACY` + `TERMS`.

## Page-by-page

### 1. Home / Shop (`Oyecouture Shop.dc.html`)
- **Hero** (charcoal, full-bleed): huge uppercase headline "TAILOR YOUR AMBITION WITH EFFORTLESS PRECISION"; subcopy "Invest in pieces that reflect your ambition brilliantly. Bespoke tailoring that speaks before you do."; champagne CTA button "DISCOVER OUR COLLECTION →".
- **Category tiles** — three large image cards labelled SHOP → **Suits**, **Trads**, **Casuals**, each with a → arrow, linking to the category pages.
- **The Collection** — eyebrow "THE COLLECTION", heading "All Pieces". Filter chips: `ALL`, `SUITS`, `ASO-OKE`, `CASUALS`, `KAFTANS`. Right side: piece count ("13 pieces") + sort dropdown ("Featured", plus Price ↑/↓, Newest). Then the product grid (see Product Card).
- Footer.

### 2. Category pages (`Suits`, `Trad Attires`, `Casuals Men`, `Casuals Women`)
Same shell (announcement, nav, footer). A section header + the product grid filtered to that section. Use the same Product Card and grid. Section→products mapping is in `product-data.js` (`section` field: `suits`, `trad`, `casuals-men`, `casuals-women`).

### 3. Product detail (`Oyecouture Product.dc.html?id=<slug>`)
- Breadcrumbs: Home / <Section> / <Product name>.
- **Left:** large image gallery with prev/next arrows + 3 thumbnails below; a tag badge overlay (e.g. BESPOKE, NEW) when present.
- **Right:** category eyebrow (gold), product name (large), price (₦, `toLocaleString`), description paragraph.
  - **If `hasSizes: true`:** a `SIZE` selector — S / M / L / XL pills, one selected by default. Then **ADD TO BAG** (charcoal button).
  - **If `hasSizes: false` (bespoke):** **no** size selector. **ADD TO BAG** (charcoal). Below the button, the note: *"Every piece is finished by hand. Sizing, fabric & delivery details are confirmed with us on WhatsApp after checkout."*
- Sold-out/placeholder images render as sand-colored blocks with the product name — most catalogue items currently have placeholder images (client owes real photography).

### 4. Cart drawer (slides in from right — see the real screenshot the user provided)
- Header: **"Your Bag · <count>"** + close X.
- Line items: thumbnail, category label (uppercase gray), product name (bold), a **qty stepper** (− n +), line price (₦), and a per-line remove X. Size shown for RTW items.
- Footer of drawer: **Subtotal** row. If the cart contains any bespoke item, show the note: *"This order will be confirmed with us on WhatsApp — sizing, fabric & delivery details are finalised there."* and the primary CTA is a dark button with gold text + WhatsApp icon: **"CONFIRM ORDER ON WHATSAPP"**. For an all-RTW cart, the CTA is **"CHECKOUT"** and goes to the pay-now flow.

### 5. About (`Oyecouture About.dc.html`)
- Hero (charcoal) with faint large emblem watermark: eyebrow "ABOUT", headline "YOUR CLOTHING SPEAKS LONG BEFORE YOU DO", intro paragraph, and section-jump chips: `OUR STORY`, `THE OYECOUTURE SYMBOL`, `CRAFTSMANSHIP`, `MISSION & VISION`.
- Then numbered editorial sections: `01 — Our Story` ("A language of confidence, worn into every room"), `02 — The Oyecouture Symbol`, `03 — Craftsmanship`, `04 — Mission & Vision`, each with copy + an image.

### 6. Checkout / order-confirmation / search overlay — NOT in the design
These screens were never designed. Build them in the same design language (charcoal/cream/champagne, Gilmer/Jost, uppercase labels, generous spacing). Keep them minimal and on-brand. This is expected scope — flag to the user that they follow the design system rather than a prior mockup.

## Components to build (inventory)

`AnnouncementBar`, `NavBar` (+ dropdowns, mobile menu), `SearchOverlay`, `Footer`, `CategoryTile`, `ProductCard` (image carousel w/ dots, wishlist heart, tag badge, name, price, View Product + Add to Bag), `FilterChips`, `SortDropdown`, `ProductGallery` (+ thumbnails), `SizeSelector`, `BespokeNote`, `CartDrawer` (+ line item, qty stepper), `MeasurementForm` (bespoke), `CheckoutForm` (RTW), `PayLinkPage` (`/pay/<token>`), `OrderConfirmation`, plus the ghost-admin screens.

## Product Card (used on Home + all category pages)
Image area (carousel with small dot indicators; wishlist heart top-right; optional tag badge top-left e.g. "BESPOKE"/"NEW"), then: category label (uppercase gray), product name (bold), price (₦), and two actions: **VIEW PRODUCT** (outline) + **ADD TO BAG** (solid charcoal). Slot naming in the design follows `oye-p-<id>-<1|2|3>` for the three images per product.
