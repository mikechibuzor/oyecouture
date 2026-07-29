# 07 · Motion & Microinteraction Design

**Status:** Approved design direction. Read before adding any animation to the storefront.
**Applies to:** all of `web/`. Pairs with `02-DESIGN-SPEC.md` (static design) — this doc governs everything that *moves*.

---

## 0. The one-paragraph brief

Oyecouture sells ₦950,000 hand-embroidered agbada and ₦78,000 linen shirts to a majority-mobile Nigerian audience. The motion language must read as **tailoring, not tech**: weighted, quiet, and confident. Almost every surface is still; movement is an accent used as sparingly as the champagne gold in the palette. We buy our polish with **three signature moments** and a lot of invisible feedback — not with animation everywhere.

---

## 1. Five principles

**P1 — Cloth, not rubber.**
Motion falls and settles like fabric: decelerate into place, never bounce back. **Spring, elastic, and overshoot easings are banned** on this project (one narrow exception: §6.3 wishlist). Bounce reads as fast fashion; weighted deceleration reads as bespoke.

**P2 — The garment never moves.**
Animate the *frame* and the *chrome* — labels, borders, overlays, buttons. The product photograph itself only scales slowly inside a fixed frame (max `1.03`). Never rotate, skew, bounce, or distort merchandise. The clothes are the client's craft; we don't fling them around.

**P3 — Bespoke is slower than ready-to-wear.**
A deliberate tempo split that mirrors the core product rule. Bespoke surfaces (agbada, suits, the WhatsApp confirm) use the `--duration-ceremony` tier and more generous reveals. RTW is brisk and transactional. **The tempo tells you which flow you're in before you read a word.**

**P4 — Motion must answer a question.**
Every animation answers exactly one of: *Did my action register?* · *Where did this come from / go?* · *What just changed?* If it answers none of those, cut it. Decorative-only motion does not ship.

**P5 — Fast where money moves.**
Browsing may breathe. Cart, checkout, and payment must feel instant — nothing on a transactional path exceeds `--duration-slow`. Anxiety is the enemy at checkout; delight is fine on the homepage.

---

## 2. Motion tokens

Create `web/assets/scss/design-system/_motion.scss` and `@use` it from `main.scss` (before `globals`), matching the existing token-file pattern.

```scss
:root {
  // ---- Duration ----------------------------------------------------------
  --duration-instant: 90ms;    // tap/press feedback, colour-only shifts
  --duration-quick: 160ms;     // hover states, chevrons, focus rings, chips
  --duration-base: 240ms;      // dropdowns, small reveals, tooltips, exits
  --duration-slow: 380ms;      // cart drawer, mobile menu, gallery crossfade
  --duration-editorial: 620ms; // scroll reveals, hero headline, section entrances
  --duration-ceremony: 800ms;  // bespoke-only moments (see P3)

  // ---- Easing ------------------------------------------------------------
  // Signature curve. Long, quiet tail — the "settle". Use for anything entering.
  --ease-settle: cubic-bezier(0.16, 1, 0.3, 1);
  // Finer settle for small elements (chips, chevrons, icons).
  --ease-fine: cubic-bezier(0.33, 1, 0.68, 1);
  // Exits: accelerate away, always faster than the matching entrance.
  --ease-exit: cubic-bezier(0.4, 0, 1, 1);
  // Neutral, for colour/opacity-only changes with no spatial meaning.
  --ease-standard: cubic-bezier(0.4, 0, 0.2, 1);

  // ---- Displacement ------------------------------------------------------
  // Luxury = short travel. Reveals never exceed --lift-lg.
  --lift-sm: 0.25rem;  // 4px  — hover nudges, inline validation
  --lift-md: 0.5rem;   // 8px  — dropdowns, menu panels
  --lift-lg: 1rem;     // 16px — scroll reveals, section entrances

  // ---- Choreography ------------------------------------------------------
  --stagger-tight: 40ms;  // grid cards, cart line items
  --stagger-loose: 80ms;  // hero headline lines, section blocks
  --zoom-image: 1.03;     // the ONLY scale allowed on product photography
}
```

**Rules of use**
- Entering → `--ease-settle`. Leaving → `--ease-exit`, and **exit duration ≈ 0.6× the entrance** (in ≈ 240ms out ≈ 160ms). Things should leave faster than they arrive.
- Never write a raw `ms` value or `cubic-bezier()` in a component file. If you need a value the scale doesn't have, add a token here and say why in the baton.
- Total stagger across a group must not exceed **~300ms**. Stagger the first ~6 items, then reveal the remainder together.

---

## 3. The three signature moments

Everything else in this document is invisible plumbing. These are the moments a customer would describe to a friend. Get these right; keep the rest quiet.

**S1 — The hero headline reveal (home).**
"TAILOR YOUR AMBITION / WITH EFFORTLESS / PRECISION" reveals line by line via a clip/mask wipe upward, `--duration-editorial`, `--stagger-loose`, `--ease-settle`. Subcopy then CTA follow. Reads like a garment being unfolded. This is the first impression — it earns the extra 600ms.

**S2 — Card → product image morph.**
Tapping a product card morphs the card's image into the detail page's hero image using the **View Transitions API** (§5.3). The single highest-leverage motion on the site: it makes an MPA feel native, and it is ~15 lines of CSS. It also visually confirms *"this is the same garment"*, which matters when someone is about to spend ₦950k.

**S3 — The bespoke WhatsApp confirm.**
The money moment. When a bespoke cart's CTA is pressed: button settles into a loading state, then the drawer hands off with a calm, deliberate `--duration-ceremony` transition into the confirmation. No sparkle, no confetti — **dignity and reassurance**. The customer is trusting a stranger with a large sum; motion here should feel like a tailor nodding, not a game.

---

## 4. Component-by-component specification

### 4.1 Global shell

| Element | Trigger | Motion | Tokens |
|---|---|---|---|
| Header | scroll past hero | border + subtle shadow fade in; backdrop blur already present | `--duration-base`, `--ease-standard` |
| Nav link | hover / focus | colour → gold **+** underline wipes in from left (`transform: scaleX()`, `transform-origin: left`) | `--duration-quick`, `--ease-fine` |
| Nav chevron | dropdown open | rotate 180° *(already built — retune to tokens)* | `--duration-quick` |
| Desktop dropdown | hover / click | fade + `--lift-md` down; exit fade only | in `--duration-base` / out `--duration-quick` |
| Mobile menu panel | tap | slide + fade from `--lift-md`; children stagger | `--duration-slow`, `--stagger-tight` |
| **All icon buttons** | `:active` | `scale(0.94)` | `--duration-instant` |
| Announcement bar | — | **static.** Do not auto-rotate messages (see §8) | — |

> **Tap feedback is not optional.** Touch devices have no hover, so `:active` press states are the *only* confirmation a tap registered. Every button, chip, card, and link needs one.

### 4.2 Product grid & catalogue

| Element | Trigger | Motion | Tokens |
|---|---|---|---|
| Card image | hover (desktop) | crossfade to second image — the fashion-industry standard | `--duration-slow`, `--ease-standard` |
| Card image | load | blur-up from LQIP placeholder (§7.3) | `--duration-slow` |
| Card frame | hover | action buttons fade in; title → gold | `--duration-quick` |
| Wishlist heart | click | fill + one contained pulse `1 → 1.15 → 1` | `--duration-base` |
| Filter chip | click | background/colour fill | `--duration-quick`, `--ease-fine` |
| Grid re-filter | filter/sort change | grid region fades out `--duration-quick`, new set fades in + `--lift-sm` with `--stagger-tight` | see §7.2 — **do not FLIP** |
| Piece count | filter change | crossfade the number only | `--duration-quick` |
| Loading state | fetch pending | calm opacity pulse `0.6 ↔ 1` over 1.2s — **not** a sweeping shimmer | `--ease-standard` |
| Category tile | hover | image `--zoom-image`, gradient deepens, arrow nudges `--lift-sm` right *(partly built)* | `--duration-slow` |

### 4.3 Product detail

| Element | Trigger | Motion | Tokens |
|---|---|---|---|
| Hero image | arrival from grid | **S2 view-transition morph** | §5.3 |
| Gallery main | thumbnail click | crossfade | `--duration-slow` |
| Thumbnail | hover | border → ink; image `--zoom-image` | `--duration-quick` |
| Size selector | select | fill + border; **no bounce** | `--duration-quick`, `--ease-fine` |
| Add to Bag | press | 4-state machine (§4.5) | — |
| Bespoke note | mount (bespoke only) | fade + `--lift-sm`, delayed ~120ms after price | `--duration-ceremony` ← **P3** |

### 4.4 Cart drawer *(Milestone B — build with this spec in hand)*

| Element | Motion | Tokens |
|---|---|---|
| Scrim | fade to dim | `--duration-base` |
| Panel | slide from right, settling | in `--duration-slow` / out `--duration-base`, `--ease-settle` |
| Line items | stagger in | `--stagger-tight` |
| Remove item | height collapse + fade | `--duration-base`, `--ease-exit` |
| Subtotal | value change | crossfade the number (never obscure the price — §8) | `--duration-quick` |
| Cart icon count | item added | number swap + contained scale pulse | `--duration-base` |
| RTW CTA → checkout | firm press + loading state | §4.5 |
| **Bespoke CTA → WhatsApp** | **S3** | `--duration-ceremony` |

Accessibility: focus trap, `Escape` closes, focus returns to the trigger, scrim click closes — reuse the dismissal pattern already in `useStorefrontNavigation`.

### 4.5 The four-state action button *(one shared pattern)*

Used by Add to Bag, Checkout, Confirm on WhatsApp, and every form submit:

1. **Idle** → 2. **Press** `scale(0.98)`, `--duration-instant` → 3. **Loading** label crossfades to a quiet indicator; button disabled; width **locked** (no reflow) → 4. **Success** brief checkmark + confirmation label for ~1.4s, then either return to idle or hand off.

Failure state replaces success with an inline message that fades + lifts `--lift-sm`. **Never shake** — shake is punitive and off-brand.

### 4.6 Search overlay *(A.6 — next up)*

Overlay fade `--duration-base`; input autofocused; results fade + `--lift-sm` with `--stagger-tight` as debounced results land; empty state crossfades. Keep it fast — search is a transactional path (P5).

### 4.7 Forms, checkout & measurement

- Field focus: border colour + ring, `--duration-quick`. Focus rings appear **instantly** — never animate a focus indicator in.
- Inline validation: message fades + lifts `--lift-sm`.
- Measurement form (bespoke): sections reveal at `--duration-base`; this is a long form — keep it calm and quiet, no per-field theatrics.
- **Payment redirect:** a full-screen calm state — "Taking you to secure payment…" with a quiet indicator. This is the highest-anxiety second in the entire funnel; motion's only job is reassurance.

---

## 5. Implementation stack

Deliberately **CSS-first**. Nigerian mobile traffic is majority mid-range Android; every KB of JS is a tax on the people most likely to buy.

**Tier 1 — CSS transitions & keyframes.** ~90% of this document. Zero bytes of JS.

**Tier 2 — Vue `<Transition>` / `<TransitionGroup>`.** Enter/leave for drawers, overlays, list items. Already in use for the mobile menu.

**Tier 3 — CSS scroll-driven animations** for reveal-on-scroll. `animation-timeline: view()` with named ranges — no `IntersectionObserver`, no JS, runs off the main thread. ~84% global support in 2026 (Chrome/Edge 115+, Safari 18+, Firefox 132+), so **always wrap in `@supports`**:

```scss
@supports (animation-timeline: view()) {
  .reveal {
    animation: reveal-rise linear both;
    animation-timeline: view();
    animation-range: entry 10% cover 35%;
  }
}
@keyframes reveal-rise {
  from { opacity: 0; transform: translateY(var(--lift-lg)); }
  to   { opacity: 1; transform: none; }
}
```
Unsupported browsers simply see the final state — a perfect progressive-enhancement fallback. **Never** hide content behind a JS-dependent reveal.

**Tier 4 — View Transitions API** for page navigation and **S2**. Cross-document transitions are production-ready cross-browser as of 2026 (Chrome 126+, Safari 18.2+). Nuxt supports this natively:

```ts
// nuxt.config.ts
app: { viewTransition: true }
```
Set to `true` (not `'always'`) — Nuxt then **automatically skips transitions for `prefers-reduced-motion: reduce`**, which does part of §6 for us. For S2, put a per-product `view-transition-name` on the card image and the same name on the detail hero.

**Tier 5 — a motion library: NOT INSTALLED, and that's a decision.**
Motion for Vue (~5kb) is the pick *if* we ever need spring physics or complex sequencing. Today nothing in this spec requires it. **Do not add a motion dependency without recording the justification in the baton.**

---

## 6. Accessibility — non-negotiable

**6.1 Reduced motion.** Ship this global block:

```scss
@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
    scroll-behavior: auto !important;
  }
}
```
**Reduce movement, never remove feedback.** A reduced-motion user must still see the button change to its loading state, the chip turn active, the heart fill. Kill the *travel*, keep the *information*. Scroll reveals must resolve to fully visible — never leave content stuck at `opacity: 0`.

**6.2 ⚠️ Existing issue — the hero marquee.** `.home-hero-media-track` runs `30s linear infinite`. WCAG **2.2.2 (Pause, Stop, Hide)** requires a pause/stop mechanism for motion that starts automatically, runs over 5 seconds, and sits alongside other content. The existing `prefers-reduced-motion` handling is good but **does not satisfy 2.2.2 on its own**. Fix by either (a) adding a visible pause control, (b) `animation-play-state: paused` on hover **and** focus-within, or (c) running a finite number of iterations. **Treat as a launch blocker for accessibility.**

**6.3 Other rules.**
- Motion must never be the *only* signal of state — pair with colour **and** text/icon.
- Focus indicators: instant, never animated, never suppressed mid-animation.
- Nothing flashes more than 3×/second.
- The wishlist pulse (§4.2) is the sole sanctioned overshoot — contained, one iteration, ≤`1.15`.

---

## 7. Performance rules

**7.1 Composited properties only.** Animate `transform`, `opacity`, `filter`, `clip-path` — these run on the compositor, off the main thread. **Never animate** `width`, `height`, `top`, `left`, `margin`, or `padding`; they force layout and directly damage **INP**, the strictest Core Web Vital. The one allowed exception is the cart's remove-item height collapse, which is short and off the critical path.

**7.2 Why no FLIP on the filter grid.** FLIP re-ordering is expensive and janky on mid-range Android for little payoff. A fade-out/fade-in of the grid region (§4.2) costs ~2 lines and looks calmer. Revisit only with real device evidence.

**7.3 Images.** Blur-up from an LQIP placeholder, `loading="lazy"` below the fold, explicit `width`/`height` (or `aspect-ratio` — already set) to hold layout and protect **CLS**. Never lazy-load or animate the LCP hero image.

**7.4 Hygiene.** `will-change` only on elements about to animate, removed after. `content-visibility: auto` on long product grids. No JS animation on the interaction path. Re-check Lighthouse ≥ 90 (roadmap item 20) *after* the motion pass — motion is the most common cause of a late Lighthouse regression.

---

## 8. Banned — do not build these

Each of these is either off-brand for quiet luxury, an accessibility problem, or a performance problem:

- ❌ **Spring / bounce / elastic easing** (P1) — reads cheap.
- ❌ **Parallax on product photography** — nausea risk + scroll jank. Editorial/lookbook only, if ever.
- ❌ **Page-load preloader animation** — destroys LCP and dates the site instantly.
- ❌ **Text scramble / typewriter headings** — hurts LCP, screen readers, and patience.
- ❌ **Custom cursors / cursor followers** — desktop-only vanity on a mobile-majority site.
- ❌ **Auto-rotating announcement bar or carousel without pause** — WCAG 2.2.2 (§6.2).
- ❌ **Confetti / celebration effects on purchase** — undignified at a ₦950,000 price point.
- ❌ **Animating a price or total in a way that obscures the number**, ever.
- ❌ **Motion > `--duration-slow` anywhere on cart → checkout → payment** (P5).
- ❌ **Hover-only affordances** — anything reachable only by hover is invisible on a phone.

---

## 9. Phased delivery (mapped to `05-BUILD-ROADMAP.md`)

**M1 — Foundation** *(alongside A.6 search; low risk, do first)*
Add `_motion.scss` tokens + the global reduced-motion block. Retune the existing chevron/menu/tile transitions to tokens. Add `:active` tap feedback site-wide. **Fix the marquee WCAG issue (§6.2).**

**M2 — Transactional motion** *(with Milestone B, the 31 Jul demo)*
Cart drawer choreography, the four-state action button, cart count, wishlist, **S3 bespoke WhatsApp confirm**. These ship with the features themselves — don't retrofit.

**M3 — Signature polish** *(Milestone C, before launch)*
**S1 hero headline**, scroll reveals (Tier 3), **S2 view-transition morph** (Tier 4), image blur-up. Then re-run Lighthouse.

**M4 — Editorial** *(post-launch / Phase 2)*
Scrollytelling for the lookbook and Instagram/blog surfaces, once the client is actually producing shoot content.

---

## 10. Review checklist

Before any motion PR merges:

- [ ] Uses tokens from `_motion.scss` — no raw `ms` or `cubic-bezier()`.
- [ ] Animates only composited properties (§7.1).
- [ ] Exit is faster than entrance.
- [ ] Has a visible `:active` state if it's tappable.
- [ ] Verified under `prefers-reduced-motion: reduce` — feedback survives, travel is gone.
- [ ] Verified on a throttled mid-range mobile profile, not just desktop.
- [ ] Answers one of the P4 questions. If not — delete it.
- [ ] Nothing on the checkout path exceeds `--duration-slow`.
