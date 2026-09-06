// Search metadata for the storefront.
//
// Titles and descriptions carry the literal words customers type - kaftan,
// agbada, aso-oke, bespoke tailoring, Lagos - which the on-page copy
// deliberately does not. The hero reads "Tailor Your Ambition With Effortless
// Precision" because that is the brand's voice; a search engine needs the plain
// noun. Keeping the two separate lets the page stay written the way the client
// wants it while still being findable.

export const SITE = {
  name: 'Oyecouture Bespoke',
  // Used to build absolute canonical and Open Graph URLs, which must be
  // absolute to count.
  url: 'https://oyecouturebespoke.com',
  locale: 'en_NG',
  // Taken from the address already shown on the payment page, so the site and
  // its structured data cannot drift apart.
  street: '6 Force Road, Tafawa Balewa Square (TBS)',
  city: 'Lagos',
  region: 'Lagos',
  country: 'NG',
  countryName: 'Nigeria',
  shareImage: '/images/brand/oyecouture-icon-gold.png',
} as const

// Appended to every page title except the homepage, which leads with the brand
// instead so a search for the name matches the front of the title.
export const TITLE_SUFFIX = `${SITE.name}` as const

export const SEO_COPY = {
  home: {
    title: 'Oyecouture Bespoke | Bespoke Tailoring in Lagos, Nigeria',
    description:
      'Bespoke tailoring house in Lagos. Custom agbada, kaftans, aso-oke, '
      + 'suits and casual wear, made to measure and delivered across Nigeria.',
  },
  suits: {
    title: 'Bespoke Suits & Wedding Suits',
    description:
      'Made-to-measure suits, tuxedos and wedding suits, cut to your '
      + 'measurements in Lagos and delivered across Nigeria.',
  },
  tradAttires: {
    title: 'Agbada, Kaftans & Aso-Oke',
    description:
      'Custom agbada, kaftans and aso-oke for weddings and ceremonies. '
      + 'Traditional Nigerian attire tailored to your measurements.',
  },
  casualsMen: {
    title: "Men's Casual Wear",
    description:
      'Tailored casual shirts, trousers and separates for men, made to '
      + 'measure in Lagos.',
  },
  casualsWomen: {
    title: "Women's Casual Wear",
    description:
      'Tailored casual pieces for women, made to measure in Lagos and '
      + 'delivered across Nigeria.',
  },
  // The care pages - terms, privacy, returns, shipping - are absent on purpose.
  // They all render StorefrontCarePage, which sets their metadata from the same
  // heading and intro the page displays, so listing them here would be a second
  // source of truth for one title.
} as const

// Words appended to a product's own name so a product page can match a search
// for the garment type, not only for the piece's given name. "The Adire Linen
// Set" contains none of the words anyone searches for on its own.
export const PRODUCT_TITLE_SUFFIX = 'Bespoke' as const
