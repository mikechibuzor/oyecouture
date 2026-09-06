import { SITE, TITLE_SUFFIX } from '~/constants/seo'

// Values may be reactive: a product page only knows its title once the product
// has loaded, and the tags have to follow it rather than freeze on the
// placeholder.
interface StorefrontSeoInput {
  title: MaybeRefOrGetter<string>
  description: MaybeRefOrGetter<string>
  // Absolute or root-relative image for link previews. Defaults to the brand
  // mark, so a shared link never renders as a blank card.
  image?: MaybeRefOrGetter<string>
  // Set on the homepage, whose title already leads with the brand name.
  bareTitle?: boolean
}

// useStorefrontSeo sets the title, description, canonical and link-preview tags
// for one page.
//
// It exists so canonical and Open Graph are not left to each page to remember.
// They were missing everywhere, which is why a link shared to WhatsApp or
// Instagram previewed as a bare URL with no image or description.
export const useStorefrontSeo = (input: StorefrontSeoInput): void => {
  const route = useRoute()

  // Canonical uses the path without query or hash, so filtered and tracked
  // variants of a page do not compete with it as separate URLs.
  const canonical = computed(() => `${SITE.url}${route.path}`)

  const title = computed(() => {
    const resolved = toValue(input.title)

    return input.bareTitle ? resolved : `${resolved} | ${TITLE_SUFFIX}`
  })

  const description = computed(() => toValue(input.description))

  const shareImage = computed(() => {
    const image = toValue(input.image) ?? SITE.shareImage

    return image.startsWith('http') ? image : `${SITE.url}${image}`
  })

  useSeoMeta({
    title: () => title.value,
    description: () => description.value,
    ogTitle: () => title.value,
    ogDescription: () => description.value,
    ogType: 'website',
    ogUrl: () => canonical.value,
    ogImage: () => shareImage.value,
    ogSiteName: SITE.name,
    ogLocale: SITE.locale,
    twitterCard: 'summary_large_image',
    twitterTitle: () => title.value,
    twitterDescription: () => description.value,
    twitterImage: () => shareImage.value,
  })

  useHead(computed(() => ({
    link: [{ rel: 'canonical', href: canonical.value }],
  })))
}
