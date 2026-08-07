import { PRODUCT_POLICIES } from '~/constants/products'
import type { ProductPolicy } from '~/types/product'

// Customer-care pages. The design links to these but never designed them —
// every link is href="#" there — and the PRD assigns the copy to the client.
// Pages still awaiting the client's wording keep `awaitingCopy: true`, which is
// what shows the placeholder notice; pages with real copy must not carry it.
export const CARE_PLACEHOLDER_NOTICE
  = 'This page is ready for your wording — send it over and we will drop it in.'

// Shipping and returns share their wording with the product page disclosures.
// Two copies of a delivery promise drift, and a customer who reads one figure
// on the product and another on the policy page has been misled.
const carePageFromPolicy = (policy: ProductPolicy) => ({
  breadcrumb: policy.title,
  heading: policy.title,
  intro: policy.intro,
  sections: policy.lines.map((line) => ({
    // Labels carry their own punctuation for inline use; a heading needs it off.
    title: (line.label ?? policy.linesHeading ?? '').replace(/:$/, ''),
    body: line.text.replace(/^—\s*/, ''),
  })),
  footnote: policy.footnote,
})

const policyBySlug = (slug: string): ProductPolicy => {
  const policy = PRODUCT_POLICIES.find((candidate) => candidate.id === slug)
  if (!policy) {
    throw new Error(`missing product policy: ${slug}`)
  }

  return policy
}

export const CARE_PAGES = {

  sizeGuide: {
    awaitingCopy: true,
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
  shipping: carePageFromPolicy(policyBySlug('shipping')),
  returns: carePageFromPolicy(policyBySlug('returns')),
  legal: {
    awaitingCopy: true,
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
    awaitingCopy: true,
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
