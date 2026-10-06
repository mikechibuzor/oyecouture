/** What any visitor can read. Everything here is published on the storefront. */
export interface SiteSettings {
  readonly announcement: string
  readonly whatsappNumber: string
  readonly contactEmail: string
  readonly held: boolean
  readonly holdHeadline: string
  readonly holdMessage: string
}

/**
 * What the settings screen reads and writes. It carries the owner-only fields as well, and is
 * only ever returned by the authenticated admin endpoint.
 */
export interface SiteSettingsDraft {
  announcement: string
  whatsappNumber: string
  contactEmail: string
  notifyEmail: string
  fromName: string
  held: boolean
  holdHeadline: string
  holdMessage: string
}

export type AdminSettings = SiteSettingsDraft

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
