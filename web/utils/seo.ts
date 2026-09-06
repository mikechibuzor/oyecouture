import { SITE } from '~/constants/seo'
import type { Product } from '~/types/product'

// Structured data tells Google what this business and its products are, rather
// than leaving it to infer them from prose. It is what produces the shop's own
// panel on a brand search, and the price and availability shown beneath a
// product result.

const postalAddress = {
  '@type': 'PostalAddress',
  'streetAddress': SITE.street,
  'addressLocality': SITE.city,
  'addressRegion': SITE.region,
  'addressCountry': SITE.country,
} as const

// ClothingStore rather than the generic LocalBusiness: the more specific type
// is what makes the result eligible for the retail treatment.
export const buildBusinessSchema = (): Record<string, unknown> => ({
  '@context': 'https://schema.org',
  '@type': 'ClothingStore',
  '@id': `${SITE.url}/#business`,
  'name': SITE.name,
  'url': SITE.url,
  'image': `${SITE.url}${SITE.shareImage}`,
  'address': postalAddress,
  'areaServed': SITE.countryName,
  'description':
    'Bespoke tailoring house in Lagos making custom agbada, kaftans, aso-oke, '
    + 'suits and casual wear to measure.',
})

export const buildWebsiteSchema = (): Record<string, unknown> => ({
  '@context': 'https://schema.org',
  '@type': 'WebSite',
  '@id': `${SITE.url}/#website`,
  'name': SITE.name,
  'url': SITE.url,
  'publisher': { '@id': `${SITE.url}/#business` },
})

// Prices are kobo integers throughout the codebase; schema.org wants a decimal
// string in major units, so the conversion happens here rather than being
// mistaken for a formatting concern.
const nairaFromKobo = (amountKobo: number): string => (amountKobo / 100).toFixed(2)

export const buildProductSchema = (product: Product): Record<string, unknown> => ({
  '@context': 'https://schema.org',
  '@type': 'Product',
  'name': product.name,
  'description': product.description,
  'sku': product.slug,
  'image': product.images?.map((image) => image.url) ?? [],
  'brand': { '@type': 'Brand', 'name': SITE.name },
  'offers': {
    '@type': 'Offer',
    'url': `${SITE.url}/product/${product.slug}`,
    'priceCurrency': 'NGN',
    'price': nairaFromKobo(product.priceKobo),
    'availability': 'https://schema.org/InStock',
    'seller': { '@id': `${SITE.url}/#business` },
  },
})

// Nuxt's useHead takes the JSON-LD as script content; serialising here keeps the
// call sites to one line each.
export const jsonLdScript = (
  schema: Record<string, unknown>,
): { type: string, innerHTML: string } => ({
  type: 'application/ld+json',
  innerHTML: JSON.stringify(schema),
})
