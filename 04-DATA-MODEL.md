# 04 · Data Model

Postgres. All money in **kobo (int64)**. Seed `products`/`product_images` from `catalogue/product-data.js`.

## Tables

```sql
-- Catalogue
products(
  id            uuid pk,
  slug          text unique,           -- matches product-data.js id (e.g. 'monarch','agbada','linen')
  name          text,
  section       text,                  -- 'suits' | 'trad' | 'casuals-men' | 'casuals-women'
  category      text,                  -- e.g. 'Three-Piece','Agbada','Shirts'
  price_kobo    bigint,                -- naira * 100
  description   text,
  tag           text,                  -- '', 'New', 'Bespoke'
  has_sizes     boolean,               -- TRUE = RTW (size selector, pay-now); FALSE = bespoke (WhatsApp-first)
  active        boolean default true,  -- ghost admin can hide placeholder items until photos arrive
  sort_order    int default 0,
  created_at    timestamptz default now()
)

product_images(
  id          uuid pk,
  product_id  uuid fk -> products,
  r2_key      text,                    -- key in Cloudflare R2
  position    int                      -- 1..3 (design uses up to 3 images)
)

-- Orders
orders(
  id                uuid pk,
  order_no          text unique,       -- human ref e.g. 'OYE-1042'
  type              text,              -- 'rtw' | 'bespoke' | 'mixed'
  status            text,              -- see state machine below
  customer_name     text,
  email             text,
  phone             text,              -- WhatsApp-capable number
  address           text,
  delivery_zone     text,              -- 'lagos' | 'nigeria' | 'international'
  delivery_fee_kobo bigint,
  subtotal_kobo     bigint,
  total_kobo        bigint,
  pay_token         text unique,       -- opaque token for /pay/<token>
  payment_provider  text,              -- 'nomba' | 'paystack'
  payment_ref       text,              -- provider transaction reference
  created_at        timestamptz default now()
)

order_items(
  id                 uuid pk,
  order_id           uuid fk -> orders,
  product_id         uuid fk -> products,
  name_snapshot      text,             -- capture name at purchase time
  price_kobo_snapshot bigint,          -- capture price at purchase time
  size               text,             -- 'S'|'M'|'L'|'XL' or null (bespoke)
  qty                int
)

-- Bespoke measurements (optional, one row per order when provided)
measurements(
  id        uuid pk,
  order_id  uuid fk -> orders,
  neck      numeric, chest numeric, waist numeric, hips numeric,
  sleeve    numeric, inseam numeric, length numeric,
  height    numeric, weight numeric,
  notes     text
)

-- Payments audit (every provider event; enables idempotency + reconciliation)
payment_events(
  id          uuid pk,
  order_id    uuid fk -> orders,
  provider    text,
  event_type  text,                    -- 'init','webhook','verify'
  payload_json jsonb,
  created_at  timestamptz default now()
)

-- Config (ghost admin editable)
delivery_zones(id uuid pk, name text, fee_kobo bigint, active boolean)
settings(key text pk, value text)      -- 'announcement_text','whatsapp_number','payment_provider', etc.
admin_users(id uuid pk, email text unique, password_hash text, role text)  -- ghost admin auth
```

## Order state machine

```
                       ┌─────────────── RTW ───────────────┐
                       ▼                                    │
   (create) ──▶ awaiting_payment ──pay+verify──▶ paid ──▶ in_progress ──▶ shipped ──▶ delivered
                       ▲                            ▲
   BESPOKE:            │                            │
   (create) ──▶ pending_confirmation ──owner shares pay link, customer pays+verify──┘
                       │
   any ──▶ cancelled / refunded
```

- **RTW** orders start at `awaiting_payment`.
- **Bespoke** orders start at `pending_confirmation` (no payment yet); after WhatsApp agreement the owner triggers the pay link; on verified payment they go to `paid` → `in_progress`.
- Transition to `paid` happens **only** via `VerifyTransaction` or a verified webhook — never on redirect.

## Seeding from product-data.js

`catalogue/product-data.js` exports `products[]` (28 items) with `{id, name, cat, section, price, tag, hasSizes, description}` and helper `sections`, `sizeOptions = ['S','M','L','XL']`, `fmt()` (₦ formatter). Map: `slug=id`, `category=cat`, `price_kobo=price*100`, `has_sizes=hasSizes`. Images are placeholders for most items — insert `product_images` rows as real photos arrive; use `active=false` to hide items lacking photos at launch if desired.
