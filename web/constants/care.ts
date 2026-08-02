// Customer-care pages. The design links to these but never designed them —
// every link is href="#" there — and the PRD assigns the copy to the client.
// These render the real page furniture with clearly-marked placeholder body
// copy, so the structure is demo-ready and swapping in real text is a
// content change rather than a build.
export const CARE_PLACEHOLDER_NOTICE
  = 'This page is ready for your wording — send it over and we will drop it in.'

export const CARE_PAGES = {
  sizeGuide: {
    breadcrumb: 'Size Guide',
    heading: 'Size Guide',
    intro:
      'Ready-to-wear pieces come in S, M, L and XL. Every made-to-measure piece is cut to your own measurements, taken with us on WhatsApp or at a fitting.',
    sections: [
      { title: 'Ready-to-wear', body: 'Standard sizing chart with chest, waist and length measurements for each size.' },
      { title: 'Made-to-measure', body: 'How we take your measurements, and what we need from you before we cut.' },
      { title: 'Between sizes?', body: 'Guidance on choosing when you sit between two ready-to-wear sizes.' },
    ],
  },
  shipping: {
    breadcrumb: 'Shipping & Delivery',
    heading: 'Shipping & Delivery',
    intro:
      'We deliver across Nigeria and worldwide. Made-to-measure timelines depend on the piece and are confirmed with you directly.',
    sections: [
      { title: 'Within Lagos', body: 'Delivery window and fee for Lagos addresses.' },
      { title: 'Rest of Nigeria', body: 'Delivery window and fee for addresses outside Lagos.' },
      { title: 'International', body: 'Worldwide delivery is quoted per order and arranged with us directly.' },
    ],
  },
  returns: {
    breadcrumb: 'Returns & Exchanges',
    heading: 'Returns & Exchanges',
    intro:
      'Ready-to-wear pieces can be returned or exchanged. Made-to-measure pieces are cut to one person and are handled case by case.',
    sections: [
      { title: 'Ready-to-wear', body: 'Return window, condition requirements and how to start a return.' },
      { title: 'Made-to-measure', body: 'How alterations and fit adjustments are handled after delivery.' },
      { title: 'Faults', body: 'What to do if a piece arrives damaged or incorrect.' },
    ],
  },
  contact: {
    breadcrumb: 'Contact Us',
    heading: 'Contact Us',
    intro:
      'The fastest way to reach us is on WhatsApp — that is where fittings, fabric and delivery are arranged.',
    sections: [
      { title: 'WhatsApp', body: 'Message us for orders, measurements and delivery.' },
      { title: 'Email', body: 'For press, wholesale and partnership enquiries.' },
      { title: 'Studio', body: 'Studio address and opening hours for fittings.' },
    ],
  },
  legal: {
    breadcrumb: 'Legal & Privacy',
    heading: 'Legal & Privacy',
    intro: 'How we handle your information when you order from us.',
    sections: [
      { title: 'What we collect', body: 'The details we take at checkout and why we need them.' },
      { title: 'How we use it', body: 'How your details are used to fulfil and deliver your order.' },
      { title: 'Your rights', body: 'How to request a copy of your data or ask us to delete it.' },
    ],
  },
  terms: {
    breadcrumb: 'Terms',
    heading: 'Terms & Conditions',
    intro: 'The terms that apply when you place an order with Oyecouture Bespoke.',
    sections: [
      { title: 'Orders', body: 'How an order is formed, confirmed and paid for.' },
      { title: 'Made-to-measure', body: 'Terms specific to pieces cut to a single measure.' },
      { title: 'Cancellations', body: 'When an order can be cancelled and what happens to payment.' },
    ],
  },
} as const
