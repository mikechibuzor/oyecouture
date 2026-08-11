import { API_ROUTES } from '~/constants/api'
import { CATEGORY_TILES, PRODUCT_COPY } from '~/constants/products'
import type { HomeContent, HomeContentDraft } from '~/types/settings'

export const defaultHomeContent = (): HomeContent => ({
  hero: {
    headingLead: PRODUCT_COPY.home.heroHeading.lead,
    headingTail: PRODUCT_COPY.home.heroHeading.tail,
    description: PRODUCT_COPY.home.heroDescription,
    actionLabel: PRODUCT_COPY.home.heroAction,
    actionTo: '/#collection',
  },
  categories: {
    eyebrow: PRODUCT_COPY.home.categoryEyebrow,
    heading: PRODUCT_COPY.categoryTilesHeading,
    tiles: CATEGORY_TILES.map((tile) => ({
      id: tile.heading.toLowerCase().replace(/\s+/g, '-'),
      label: tile.label,
      heading: tile.heading,
      imagePath: tile.imagePath,
      imageAlt: tile.imageAlt,
      to: tile.to,
    })),
  },
})

export const toHomeContentDraft = (homeContent: HomeContent): HomeContentDraft => ({
  hero: { ...homeContent.hero },
  categories: {
    eyebrow: homeContent.categories.eyebrow,
    heading: homeContent.categories.heading,
    tiles: homeContent.categories.tiles.map((tile) => ({ ...tile })),
  },
})

export const useHomeContent = () => {
  const fallbackContent = defaultHomeContent()

  return useFetch<HomeContent>(API_ROUTES.homeContent, {
    baseURL: useApiBaseUrl(),
    default: () => fallbackContent,
    key: 'home-content',
    retry: 0,
  })
}
