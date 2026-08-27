export enum StorefrontRoute {
  Home = '/',
  TradAttires = '/trad-attires',
  Suits = '/suits',
  CasualsMen = '/casuals/men',
  CasualsWomen = '/casuals/women',
  Collection = '/#collection',
  DearGroom = '/#collection', // TODO: curated collection, no design yet. See audit.
  //  was '/collection/dear-groom',
  WeiXin = '/#collection', // TODO: curated collection, no design yet. See audit.
  //  was '/collection/wei-xin',
  ManOfLegacy = '/#collection', // TODO: curated collection, no design yet. See audit.
  //  was '/collection/man-of-legacy',
  About = '/about',
  OurStory = '/about#our-story',
  Symbol = '/about#symbol',
  Craftsmanship = '/about#craftsmanship',
  Mission = '/about#mission',
  SizeGuide = '/size-guide',
  Measurements = '/measurements',
  Shipping = '/shipping',
  Returns = '/returns',
  BookFitting = '/#bespoke',
  Contact = '/contact',
  Wishlist = '/wishlist',
  Checkout = '/checkout',
  Confirm = '/confirm',
  Legal = '/legal',
  Terms = '/terms',
}

export enum StorefrontNavigationItemId {
  Home = 'home',
  TradAttires = 'trad-attires',
  Suits = 'suits',
  Casuals = 'casuals',
  Collection = 'collection',
}

export enum StorefrontActionId {
  Search = 'search',
  Account = 'account',
  Wishlist = 'wishlist',
  Cart = 'cart',
}

export enum StorefrontSocialPlatform {
  Instagram = 'instagram',
  X = 'x',
  TikTok = 'tiktok',
  WhatsApp = 'whatsapp',
}

export interface StorefrontNavigationChild {
  readonly label: string
  readonly to: StorefrontRoute
}

export interface StorefrontNavigationItem {
  readonly id: StorefrontNavigationItemId
  readonly label: string
  readonly to?: StorefrontRoute
  readonly menuId?: string
  readonly children?: readonly StorefrontNavigationChild[]
}

export interface StorefrontAction {
  readonly id: StorefrontActionId
  readonly label: string
}

export interface StorefrontFooterLink {
  readonly label: string
  readonly to: StorefrontRoute
}

export interface StorefrontFooterGroup {
  readonly heading: string
  readonly links: readonly StorefrontFooterLink[]
}

export interface StorefrontSocialLink {
  readonly platform: StorefrontSocialPlatform
  readonly label: string
  readonly href: string
}
