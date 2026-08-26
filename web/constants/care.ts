import { PRODUCT_POLICIES } from '~/constants/products'
import type { ProductPolicy } from '~/types/product'

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
    body: line.text,
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
  shipping: carePageFromPolicy(policyBySlug('shipping')),
  returns: carePageFromPolicy(policyBySlug('returns')),
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
