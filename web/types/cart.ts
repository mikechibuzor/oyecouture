import type { ProductSection, ProductSize } from '~/types/product'

// How the whole cart must be checked out. A cart containing any bespoke piece
// is confirmed on WhatsApp because measurements must be agreed before money moves.
export enum CartCheckoutRoute {
  Bespoke = 'bespoke',
  ReadyToWear = 'ready-to-wear',
}

// One line in the bag. Price and name are snapshotted at the time of adding so
// the bag stays stable while browsing; the server re-prices from the database
// when the order is created, so a stale snapshot can never set what is charged.
export interface CartLine {
  readonly key: string
  readonly slug: string
  readonly name: string
  readonly section?: ProductSection
  readonly category: string
  readonly priceKobo: number
  readonly hasSizes: boolean
  readonly size: ProductSize | null
  readonly imageUrl: string | null
  readonly quantity: number
}
