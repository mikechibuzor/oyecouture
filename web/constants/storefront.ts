import type {
  StorefrontAction,
  StorefrontFooterGroup,
  StorefrontNavigationItem,
  StorefrontSocialLink,
} from '~/types/storefront'
import {
  StorefrontActionId,
  StorefrontNavigationItemId,
  StorefrontRoute,
  StorefrontSocialPlatform,
} from '~/types/storefront'

export const STOREFRONT_COPY = {
  announcement:
    'COMPLIMENTARY TAILORING CONSULTATION · WORLDWIDE DELIVERY · BE EXCLUSIVE',
  announcementLabel: 'Store announcement',
  brandAlt: 'Oyecouture Bespoke',
  brandLogoPath: '/images/brand/oyecouture-logo.png',
  navigationLabel: 'Primary navigation',
  mobileMenuId: 'mobile-navigation-menu',
  openMenu: 'Open navigation menu',
  closeMenu: 'Close navigation menu',
  footer: {
    heading: 'Stay Connected',
    description:
      'Oyecouture Bespoke, where confidence, simplicity, and class become clothing.',
    copyright: '© 2026 Oyecouture. All rights reserved.',
    legalNavigationLabel: 'Legal navigation',
    socialNavigationLabel: 'Oyecouture social media',
  },
} as const

export const BESPOKE_SERVICE_COPY = {
  eyebrow: 'The Bespoke Service',
  heading: 'Prefer a piece made to a single measure, yours?',
  action: 'Book an Appointment',
} as const

export const STOREFRONT_NAVIGATION: readonly StorefrontNavigationItem[] = [
  {
    id: StorefrontNavigationItemId.Home,
    label: 'Home',
    to: StorefrontRoute.Home,
  },
  {
    id: StorefrontNavigationItemId.TradAttires,
    label: 'Trad Attires',
    to: StorefrontRoute.TradAttires,
  },
  {
    id: StorefrontNavigationItemId.Suits,
    label: 'Suits',
    to: StorefrontRoute.Suits,
  },
  {
    id: StorefrontNavigationItemId.Casuals,
    label: 'Casuals',
    menuId: 'casuals-navigation-menu',
    children: [
      { label: 'Men', to: StorefrontRoute.CasualsMen },
      { label: 'Women', to: StorefrontRoute.CasualsWomen },
    ],
  },
]

export const DESKTOP_HEADER_ACTIONS = [
  StorefrontActionId.Search,
  StorefrontActionId.Wishlist,
  StorefrontActionId.Cart,
] as const satisfies readonly StorefrontActionId[]

export const STOREFRONT_ACTION_BY_ID = {
  [StorefrontActionId.Search]: {
    id: StorefrontActionId.Search,
    label: 'Search',
  },
  [StorefrontActionId.Account]: {
    id: StorefrontActionId.Account,
    label: 'Account',
  },
  [StorefrontActionId.Wishlist]: {
    id: StorefrontActionId.Wishlist,
    label: 'Wishlist',
  },
  [StorefrontActionId.Cart]: {
    id: StorefrontActionId.Cart,
    label: 'Bag',
  },
} as const satisfies Readonly<Record<StorefrontActionId, StorefrontAction>>

export const MOBILE_HEADER_ACTIONS = [
  StorefrontActionId.Wishlist,
  StorefrontActionId.Cart,
] as const satisfies readonly StorefrontActionId[]

export const MOBILE_MENU_ACTIONS = [
  StorefrontActionId.Search,
] as const satisfies readonly StorefrontActionId[]

export const STOREFRONT_FOOTER_GROUPS = [
  {
    heading: 'Shop',
    links: [
      { label: 'Trad Attires', to: StorefrontRoute.TradAttires },
      { label: 'Suits', to: StorefrontRoute.Suits },
      { label: 'Casuals', to: StorefrontRoute.CasualsMen },
    ],
  },
  {
    heading: 'About',
    links: [
      { label: 'Our Story', to: StorefrontRoute.OurStory },
      { label: 'The Oyecouture Symbol', to: StorefrontRoute.Symbol },
      { label: 'Craftsmanship', to: StorefrontRoute.Craftsmanship },
      { label: 'Mission and Vision', to: StorefrontRoute.Mission },
    ],
  },
  {
    heading: 'Customer Care',
    links: [
      { label: 'Size Guide', to: StorefrontRoute.SizeGuide },
      { label: 'Shipping & Delivery', to: StorefrontRoute.Shipping },
      { label: 'Returns & Exchanges', to: StorefrontRoute.Returns },
      { label: 'Book a Fitting', to: StorefrontRoute.BookFitting },
      { label: 'Contact Us', to: StorefrontRoute.Contact },
    ],
  },
] as const satisfies readonly StorefrontFooterGroup[]

export const STOREFRONT_SOCIAL_LINKS = [
  {
    platform: StorefrontSocialPlatform.Instagram,
    label: 'Oyecouture on Instagram',
    href: 'https://www.instagram.com/_oyecouture?igsi=OGwyYWtsY2VxbG4z&utm_source=qr',
  },
  {
    platform: StorefrontSocialPlatform.TikTok,
    label: 'Oyecouture on TikTok',
    href: 'https://www.tiktok.com/@oye_oyecouture?_r=1&_t=ZS-99CydxzaFeM',
  },
  {
    platform: StorefrontSocialPlatform.X,
    label: 'Oyecouture on X',
    href: 'https://x.com/_oyecouture?s=11',
  },
] as const satisfies readonly StorefrontSocialLink[]

export const INTERNAL_STOREFRONT_ROUTES = [
  StorefrontRoute.Home,
  StorefrontRoute.TradAttires,
  StorefrontRoute.Suits,
  StorefrontRoute.CasualsMen,
  StorefrontRoute.CasualsWomen,
  StorefrontRoute.Collection,
  StorefrontRoute.Wishlist,
] as const

const internalStorefrontRouteSet: ReadonlySet<StorefrontRoute> = new Set(
  INTERNAL_STOREFRONT_ROUTES,
)

// Reports whether a completed route should use Nuxt client-side navigation.
export const isInternalStorefrontRoute = (storefrontRoute: StorefrontRoute): boolean =>
  internalStorefrontRouteSet.has(storefrontRoute)
