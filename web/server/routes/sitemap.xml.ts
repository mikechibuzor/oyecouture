import { SITE } from '~/constants/seo'

// The sitemap is how the product pages get discovered.
//
// Nothing on the storefront links to every piece at once, so without this a
// crawler only reaches whatever happens to be on the homepage. Product URLs are
// read from the API at request time rather than hardcoded, so a piece added
// through the admin is listed without a redeploy.

interface SitemapProduct {
  slug: string
}

interface ProductListResponse {
  products?: SitemapProduct[]
}

// Static routes worth indexing. The checkout, payment, confirmation, tracking
// and wishlist pages are deliberately absent: they are either private to one
// customer or meaningless without an order behind them.
const STATIC_PATHS: readonly string[] = [
  '/',
  '/suits',
  '/trad-attires',
  '/casuals/men',
  '/casuals/women',
  '/about',
  '/contact',
  '/size-guide',
  '/measurements',
  '/shipping',
  '/returns',
  '/terms',
  '/legal',
]

const urlEntry = (path: string, priority: string): string =>
  `  <url>\n`
  + `    <loc>${SITE.url}${path}</loc>\n`
  + `    <priority>${priority}</priority>\n`
  + `  </url>`

export default defineEventHandler(async (event) => {
  const paths = [...STATIC_PATHS]

  // A failed product fetch must not take the sitemap down with it: a sitemap
  // listing only the static pages is far better than a 500, which Google reads
  // as the site having no sitemap at all.
  try {
    const apiBase = useRuntimeConfig(event).apiInternalBase
    const response = await $fetch<ProductListResponse>(`${apiBase}/products`)

    for (const product of response.products ?? []) {
      if (product.slug) {
        paths.push(`/product/${product.slug}`)
      }
    }
  } catch (fetchError) {
    console.error('sitemap: product fetch failed', fetchError)
  }

  const body
    = '<?xml version="1.0" encoding="UTF-8"?>\n'
      + '<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">\n'
      + paths.map((path) => urlEntry(path, path === '/' ? '1.0' : '0.8')).join('\n')
      + '\n</urlset>'

  setHeader(event, 'content-type', 'application/xml; charset=utf-8')
  setHeader(event, 'cache-control', 'public, max-age=3600')

  return body
})
