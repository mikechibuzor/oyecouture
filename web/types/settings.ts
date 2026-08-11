export interface SiteSettings {
  readonly announcement: string
  readonly whatsappNumber: string
}

export interface SiteSettingsDraft {
  announcement: string
  whatsappNumber: string
}

export interface HomeContent {
  readonly hero: HomeHeroContent
  readonly categories: HomeCategoryContent
}

export interface HomeHeroContent {
  readonly headingLead: string
  readonly headingTail: string
  readonly description: string
  readonly actionLabel: string
  readonly actionTo: string
}

export interface HomeCategoryContent {
  readonly eyebrow: string
  readonly heading: string
  readonly tiles: readonly HomeCategoryTile[]
}

export interface HomeCategoryTile {
  readonly id: string
  readonly label: string
  readonly heading: string
  readonly imagePath: string
  readonly imageAlt: string
  readonly to: string
}

export interface HomeContentDraft {
  hero: {
    headingLead: string
    headingTail: string
    description: string
    actionLabel: string
    actionTo: string
  }
  categories: {
    eyebrow: string
    heading: string
    tiles: HomeCategoryTileDraft[]
  }
}

export interface HomeCategoryTileDraft {
  id: string
  label: string
  heading: string
  imagePath: string
  imageAlt: string
  to: string
}
