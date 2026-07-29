# Catalogue (human-readable)

Canonical data: [`product-data.js`](product-data.js) (27 products). Prices in naira; store as kobo (×100). **`Sizes` column = the critical RTW vs bespoke split** — `hasSizes:false` = bespoke (WhatsApp-first, pay-after-measurement); `hasSizes:true` = ready-to-wear (size picker, pay-now).

## Suits — section `suits` (all bespoke, no sizes)
| slug | Name | Category | Price | Tag |
|---|---|---|---|---|
| monarch | The Monarch Three-Piece | Three-Piece | ₦850,000 | Bespoke |
| sovereign | The Sovereign Suit | Two-Piece | ₦620,000 | |
| regent | The Regent Two-Piece | Two-Piece | ₦540,000 | |
| heir | The Heir Blazer | Blazers | ₦380,000 | New |
| ambassador | The Ambassador Three-Piece | Three-Piece | ₦720,000 | New |
| chairman | The Chairman Three-Piece | Three-Piece | ₦680,000 | |
| diplomat | The Diplomat Two-Piece | Two-Piece | ₦490,000 | |
| weekend | The Weekend Linen Blazer | Blazers | ₦275,000 | |

## Trad Attires — section `trad` (all bespoke, no sizes)
| slug | Name | Category | Price | Tag |
|---|---|---|---|---|
| agbada | The Royal Agbada | Agbada | ₦950,000 | Bespoke |
| heritage | Heritage Aso-Oke Set | Aso-Oke | ₦720,000 | |
| fila | Aso-Oke Fila & Wrapper | Aso-Oke | ₦410,000 | New |
| emir | The Emir Kaftan | Kaftans | ₦290,000 | |
| senator | Embroidered Senator Kaftan | Kaftans | ₦240,000 | |
| classic | Classic White Kaftan | Kaftans | ₦165,000 | |
| ceremonial | The Ceremonial Agbada | Agbada | ₦780,000 | |
| executive | The Executive Agbada | Agbada | ₦690,000 | New |
| wedding | Aso-Oke Wedding Ensemble | Aso-Oke | ₦860,000 | Bespoke |

## Casuals · Men — section `casuals-men` (all ready-to-wear, sizes S–XL)
| slug | Name | Category | Price | Tag |
|---|---|---|---|---|
| linen | The Weekend Linen Shirt | Shirts | ₦78,000 | |
| oxford | The Oxford Overshirt | Shirts | ₦86,000 | New |
| polo | The Resort Polo | Shirts | ₦54,000 | |
| chino | Tailored Chino Trouser | Trousers | ₦92,000 | New |
| jogger | The Utility Jogger | Trousers | ₦74,000 | |
| short | Tailored Cotton Short | Trousers | ₦58,000 | |
| crew | The Crewneck Sweater | Knitwear | ₦88,000 | |
| cardigan | The Merino Cardigan | Knitwear | ₦96,000 | New |

## Casuals · Women — section `casuals-women` (all ready-to-wear, sizes S–XL)
| slug | Name | Category | Price | Tag |
|---|---|---|---|---|
| wrap | The Silk Wrap Blouse | Tops | ₦82,000 | New |
| linenshirt | The Linen Shirt Dress | Tops | ₦96,000 | |

## Notes
- Sizes for RTW: `S`, `M`, `L`, `XL`.
- The Shop page's filter chips (`ALL / SUITS / ASO-OKE / CASUALS / KAFTANS`) filter by category/section, not by the four `section` values one-to-one — e.g. "ASO-OKE" and "KAFTANS" filter within `trad`. Wire filters from the `category`/`section` fields.
- Most products currently have **placeholder images** in the design — real photography is a client deliverable. Use `active=false` to hide photo-less items at launch if preferred.
