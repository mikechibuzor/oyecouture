# How to view the design

These 7 `.dc.html` files are the **approved visual design** — self-contained (all images embedded as base64). Just **double-click any file to open it in a browser.** No server, no internet needed.

| File | Page |
|---|---|
| `Oyecouture Shop.dc.html` | Home / Shop (hero, category tiles, collection grid) |
| `Oyecouture Suits.dc.html` | Suits category |
| `Oyecouture Trad Attires.dc.html` | Trad Attires category |
| `Oyecouture Casuals Men.dc.html` | Casuals · Men category |
| `Oyecouture Casuals Women.dc.html` | Casuals · Women category |
| `Oyecouture Product.dc.html` | Product detail — append `?id=<slug>` |
| `Oyecouture About.dc.html` | About |

**Product detail** needs a product id in the URL, e.g. open `Oyecouture Product.dc.html` then add to the address bar:
- `?id=agbada` → The Royal Agbada (bespoke — no size selector, WhatsApp note)
- `?id=monarch` → The Monarch Three-Piece (bespoke)
- `?id=linen` → The Weekend Linen Shirt (RTW — shows S/M/L/XL size selector)

Compare a bespoke item vs an RTW item to see the two product-page variants the build must produce.

## Important
- These are **compiled/bundled exports** (a big inline `<script>` renders the page). They are **reference only** — do **not** copy their JS into the build. Rebuild the markup cleanly in Nuxt from `../02-DESIGN-SPEC.md`.
- The **cart drawer** and **about** interactions are best seen live in the browser. The cart CTA in the design reads **"CONFIRM ORDER ON WHATSAPP"** for bespoke items — that is the intended bespoke flow.
- Some product images are **placeholders** (sand-colored blocks with the product name) — real photography is a pending client deliverable.
