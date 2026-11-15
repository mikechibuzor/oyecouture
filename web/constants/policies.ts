import type { CarePageContent } from '~/types/policy'
import { list, paragraph } from '~/types/policy'

// The client's four policy documents, received 15 Nov 2026, set as the content
// of /terms, /legal, /shipping and /returns.
//
// The wording is the client's and is reproduced as written, with one exception:
// the documents typeset the name as "Oye Couture", and it is normalised here to
// "Oyecouture" to match the rest of the site. Nothing else is paraphrased -
// these are the terms the shop is held to, so tightening the prose is not ours
// to do.
//
// Section numbers from the documents are deliberately dropped. The pages render
// each section as a heading rather than a numbered clause, and the Privacy
// Policy's own numbering skips from 7 to 9, which would have to be either
// reproduced as an error or silently corrected. Neither is necessary when the
// headings carry no numbers at all.

export const CONTACT = {
  email: 'hello@oyecouturebespoke.com',
  phone: '+234 815 219 6451',
  website: 'https://oyecouturebespoke.com/',
} as const

// The figures that appear in more than one place: the product page's delivery
// summary quotes them, and three of the four policy documents restate them.
// They live here once so a change to the promise cannot land on one page and
// not another - which is exactly what happened before, when the product page
// said "after work is completed" and the policy said "after dispatch".
export const DELIVERY_TERMS = {
  production: 'three (3) weeks or less',
  domestic: '1–3 working days after dispatch',
  international: '5–7 days after dispatch',
  expediteSurcharge: '20%',
} as const

const contactBlocks = (closing?: string) => [
  paragraph(`Email: ${CONTACT.email}`),
  paragraph(`Phone/WhatsApp: ${CONTACT.phone}`),
  paragraph(`Website: ${CONTACT.website}`),
  ...(closing ? [paragraph(closing)] : []),
]

export const TERMS_PAGE: CarePageContent = {
  breadcrumb: 'Terms',
  heading: 'Terms & Conditions',
  intro:
    'Welcome to Oyecouture Bespoke. These Terms and Conditions govern purchases '
    + 'made through our website and other official sales channels. By placing an '
    + 'order, you agree to the following terms.',
  sections: [
    {
      title: 'Products and Services',
      blocks: [
        paragraph(
          'Oyecouture offers bespoke, made-to-measure and ready-to-wear fashion '
          + 'products. Product descriptions, images, prices and specifications are '
          + 'provided to guide customers. Minor variations in colour, texture and '
          + 'appearance may occur due to fabric characteristics, lighting and '
          + 'screen settings.',
        ),
      ],
    },
    {
      title: 'Orders and Payments',
      blocks: [
        list(
          'All orders require 100% payment upfront.',
          'An order is considered confirmed upon successful receipt and verification of full payment.',
          'Production or order processing will not commence until payment has been confirmed.',
          'Customers are responsible for ensuring that their order details, measurements, '
          + 'delivery information and product selections are accurate.',
          'Any additional customisation or changes requested after order confirmation '
          + 'may attract additional charges.',
        ),
      ],
    },
    {
      title: 'Order Cancellation',
      blocks: [
        list(
          'Customers may request cancellation and a refund within 24 hours of successful payment.',
          'Upon expiration of the 24-hour period, the order shall be deemed to have '
          + 'entered production for the purposes of our cancellation policy.',
          'Cancellation or refund requests made after this period may be declined or '
          + 'subject to reasonable charges for costs already incurred.',
          'This provision does not affect any refund or cancellation rights available '
          + 'to customers under applicable consumer protection laws.',
        ),
      ],
    },
    {
      title: 'Bespoke Orders and Measurements',
      blocks: [
        paragraph(
          'Bespoke garments are produced according to individual customer measurements '
          + 'and approved specifications. Customers are responsible for ensuring that '
          + 'measurements and design preferences supplied by them are accurate.',
        ),
        paragraph(
          'Oyecouture shall not be responsible for fitting issues arising solely from '
          + 'inaccurate customer-supplied measurements. However, where an issue results '
          + 'from a manufacturing error, appropriate corrective measures will be considered.',
        ),
      ],
    },
    {
      title: 'Production and Delivery',
      blocks: [
        paragraph(
          `Our standard production timeline is ${DELIVERY_TERMS.production}, subject to `
          + 'the nature and complexity of the order.',
        ),
        paragraph(
          `Domestic delivery is estimated at ${DELIVERY_TERMS.domestic}, while `
          + `international delivery is estimated at ${DELIVERY_TERMS.international}.`,
        ),
        paragraph(
          'Expedited production may be available at an additional charge of '
          + `${DELIVERY_TERMS.expediteSurcharge} of the item's price, subject to availability.`,
        ),
      ],
    },
    {
      title: 'Intellectual Property',
      blocks: [
        paragraph(
          'All brand names, logos, photographs, product images, designs and other '
          + 'creative materials displayed on this website belong to or are lawfully '
          + 'used by Oyecouture. Unauthorised reproduction, distribution or commercial '
          + 'use is prohibited.',
        ),
      ],
    },
    {
      title: 'Customer Complaints',
      blocks: [
        paragraph(
          'Customers may contact us through our official communication channels for '
          + 'enquiries, complaints or assistance relating to their orders.',
        ),
      ],
    },
    {
      title: 'Governing Law',
      blocks: [
        paragraph(
          'These Terms and Conditions shall be governed by the laws of the Federal '
          + 'Republic of Nigeria, including applicable consumer protection legislation.',
        ),
      ],
    },
    {
      title: 'Contact Us',
      blocks: [paragraph('For enquiries, please contact:'), ...contactBlocks()],
    },
  ],
}

export const PRIVACY_PAGE: CarePageContent = {
  breadcrumb: 'Privacy',
  heading: 'Privacy Policy',
  intro:
    'Oyecouture Bespoke respects your privacy and is committed to protecting the '
    + 'personal information you share with us when using our website, purchasing our '
    + 'products or accessing our tailoring services.',
  sections: [
    {
      title: 'Our Commitment',
      blocks: [
        paragraph(
          'This Policy explains how we collect, use and protect your information in '
          + 'accordance with the Nigeria Data Protection Act 2023 and applicable data '
          + 'protection requirements.',
        ),
      ],
    },
    {
      title: 'Information We Collect',
      blocks: [
        paragraph('Depending on your interaction with us, we may collect:'),
        list(
          'Your name, email address and telephone number.',
          'Billing and delivery addresses.',
          'Order history and transaction references.',
          'Body measurements and fitting preferences for bespoke services.',
          'Design, fabric and product preferences.',
          'Website usage and technical information, including device and browser information.',
        ),
      ],
    },
    {
      title: 'How We Use Your Information',
      blocks: [
        paragraph('We collect and use personal information to:'),
        list(
          'Process and fulfil customer orders.',
          'Provide bespoke tailoring and measurement services.',
          'Arrange deliveries and communicate order updates.',
          'Respond to enquiries and complaints.',
          'Improve our products, services and website.',
          'Maintain business and transaction records.',
          'Comply with applicable legal obligations.',
        ),
      ],
    },
    {
      title: 'Sharing of Information',
      blocks: [
        paragraph(
          'We may share relevant customer information with trusted third-party service '
          + 'providers where necessary, including:',
        ),
        list(
          'Payment processors.',
          'Courier and delivery companies.',
          'Website hosting and technical service providers.',
          'Professional advisers and regulatory authorities where legally required.',
        ),
        paragraph('We do not sell customers’ personal information to third parties.'),
      ],
    },
    {
      title: 'Protection and Retention of Information',
      blocks: [
        paragraph(
          'We take reasonable technical and organisational measures to protect customer '
          + 'information against unauthorised access, disclosure, alteration or loss.',
        ),
        paragraph(
          'Personal information is retained only for as long as reasonably necessary to '
          + 'fulfil the purpose for which it was collected, comply with legal obligations '
          + 'or resolve disputes.',
        ),
      ],
    },
    {
      title: 'Customer Photographs and Marketing',
      blocks: [
        paragraph(
          'We may wish to showcase completed garments, customer testimonials or '
          + 'photographs for promotional purposes.',
        ),
        paragraph(
          'Where such materials identify a customer, we will obtain the necessary '
          + 'permission or establish another lawful basis before using them for '
          + 'marketing or publication.',
        ),
        paragraph(
          'Customers are not required to consent to promotional use of their photographs '
          + 'as a condition of purchasing our products or services.',
        ),
      ],
    },
    {
      title: 'Cookies',
      blocks: [
        paragraph(
          'Our website may use cookies and similar technologies to improve browsing '
          + 'experience, maintain website functionality and understand website usage.',
        ),
        paragraph(
          'Where required, consent will be obtained before non-essential cookies are '
          + 'activated. Customers may also manage cookies through their browser settings.',
        ),
      ],
    },
    {
      title: 'Your Privacy Rights',
      blocks: [
        paragraph('Subject to applicable law, you have the right to:'),
        list(
          'Request access to your personal information.',
          'Request correction of inaccurate information.',
          'Request deletion of your information in appropriate circumstances.',
          'Object to certain forms of processing.',
          'Withdraw consent where processing is based on consent.',
          'Request a copy of your information where data portability applies.',
          'Lodge a complaint with the Nigeria Data Protection Commission.',
        ),
      ],
    },
    {
      title: 'Changes to this Policy',
      blocks: [
        paragraph(
          'We may update this Privacy Policy from time to time to reflect changes in our '
          + 'operations, technology or applicable laws. Updated versions will be '
          + 'published on our website.',
        ),
      ],
    },
    {
      title: 'Contact Us',
      blocks: [
        paragraph('For privacy-related enquiries or requests, please contact:'),
        ...contactBlocks(),
      ],
    },
  ],
  footnote:
    'Oyecouture Bespoke is committed to treating your personal information with care, '
    + 'confidentiality and respect.',
}

export const SHIPPING_PAGE: CarePageContent = {
  breadcrumb: 'Shipping',
  heading: 'Shipping & Delivery',
  intro:
    'At Oyecouture, we are committed to ensuring that every order is carefully '
    + 'prepared, packaged and delivered to our customers.',
  sections: [
    {
      title: 'Order Processing and Production',
      blocks: [
        paragraph('All orders require full payment before processing or production.'),
        paragraph(
          'Our standard production timeline for bespoke and made-to-order garments is '
          + `${DELIVERY_TERMS.production}, following confirmation of payment, measurements `
          + 'and approved specifications.',
        ),
        paragraph(
          'Production timelines may vary depending on the complexity of the design, '
          + 'fabric availability and other relevant production requirements.',
        ),
      ],
    },
    {
      title: 'Domestic Delivery',
      blocks: [
        paragraph(
          `Deliveries within Nigeria are estimated at ${DELIVERY_TERMS.domestic}, `
          + 'depending on the destination and courier arrangements.',
        ),
      ],
    },
    {
      title: 'International Delivery',
      blocks: [
        paragraph('We offer international delivery to available destinations.'),
        paragraph(
          `International deliveries are estimated at ${DELIVERY_TERMS.international}, `
          + 'subject to customs clearance, courier operations and destination-country '
          + 'requirements.',
        ),
        paragraph(
          'Any applicable customs duties, import taxes or charges imposed by the '
          + 'destination country shall be borne by the customer unless otherwise stated.',
        ),
      ],
    },
    {
      title: 'Expedited Production',
      blocks: [
        paragraph(
          'Customers who require urgent orders may request expedited production, '
          + 'subject to availability.',
        ),
        paragraph(
          `An additional charge of ${DELIVERY_TERMS.expediteSurcharge} of the item’s `
          + 'price applies to expedited production.',
        ),
        paragraph(
          'Customers are encouraged to communicate their required delivery dates before '
          + 'placing urgent orders.',
        ),
      ],
    },
    {
      title: 'Delivery Information',
      blocks: [
        paragraph(
          'Customers are responsible for providing accurate delivery addresses, telephone '
          + 'numbers and other necessary delivery information.',
        ),
        paragraph(
          'Please ensure that someone is available to receive your package at the '
          + 'designated address.',
        ),
        paragraph(
          'Additional delivery charges may apply where a failed delivery results from '
          + 'incorrect delivery information or the unavailability of the recipient.',
        ),
      ],
    },
    {
      title: 'Tracking and Delivery Issues',
      blocks: [
        paragraph(
          'Where tracking is available, relevant tracking information will be '
          + 'communicated to the customer after dispatch.',
        ),
        paragraph(
          'Customers are encouraged to inspect their packages upon receipt and promptly '
          + 'report missing items, visible damage or delivery discrepancies.',
        ),
        paragraph(
          'Where an order is delayed, lost or damaged in transit, please contact us so '
          + 'that we can investigate and take appropriate steps to resolve the matter.',
        ),
      ],
    },
    {
      title: 'Packaging',
      blocks: [
        paragraph(
          'Orders are carefully packaged to preserve the quality and presentation of our '
          + 'products. Our packaging may include garment bags, dust protection and '
          + 'garment care information, depending on the order.',
        ),
      ],
    },
    {
      title: 'Contact Us',
      blocks: [
        paragraph('For shipping enquiries or delivery assistance:'),
        ...contactBlocks(),
      ],
    },
  ],
  footnote: 'Thank you for choosing Oyecouture Bespoke.',
}

export const RETURNS_PAGE: CarePageContent = {
  breadcrumb: 'Returns',
  heading: 'Returns, Refunds & Exchanges',
  intro:
    'At Oyecouture, we take pride in delivering quality garments and ensuring customer '
    + 'satisfaction. Please review our returns and exchange guidelines before placing '
    + 'an order.',
  sections: [
    {
      title: 'Order Cancellation and Refunds',
      blocks: [
        paragraph(
          'Customers may request cancellation and a refund within 24 hours of successful '
          + 'payment.',
        ),
        paragraph(
          'After the expiration of this period, an order shall be deemed to have entered '
          + 'production for the purposes of our cancellation policy. Ordinary cancellation '
          + 'or refund requests arising from a change of mind may no longer be accepted, '
          + 'except where otherwise required by applicable law.',
        ),
      ],
    },
    {
      title: 'Bespoke and Made-to-Measure Orders',
      blocks: [
        paragraph('Due to the customised nature of bespoke garments:'),
        list(
          'Bespoke orders are not ordinarily eligible for returns or exchanges due to a '
          + 'change of mind or personal preference.',
          'Customers must carefully confirm their measurements, designs, fabrics, colours '
          + 'and other specifications before production.',
          'Fitting issues resulting solely from inaccurate measurements supplied or '
          + 'approved by the customer may attract additional alteration charges.',
        ),
      ],
    },
    {
      title: 'Manufacturing Defects',
      blocks: [
        paragraph(
          'Where a garment has a manufacturing defect, poor finishing or materially '
          + 'differs from the approved specifications, customers should notify us '
          + 'promptly, preferably within three (3) days of delivery.',
        ),
        paragraph(
          'Following an assessment, Oyecouture may offer an appropriate remedy, including '
          + 'repair, alteration, replacement or refund, as applicable.',
        ),
        paragraph(
          'The three-day reporting period does not extinguish any statutory consumer rights.',
        ),
      ],
    },
    {
      title: 'Ready-to-Wear Products',
      blocks: [
        paragraph(
          'Requests for discretionary exchanges or store credit for ready-to-wear products '
          + 'should be made within three (3) days of delivery.',
        ),
        paragraph('Returned items must:'),
        list(
          'Be unworn and unused.',
          'Be in their original condition.',
          'Have all original tags and labels attached.',
          'Be free from stains, damage or alterations.',
        ),
        paragraph(
          'Approved discretionary exchanges may be processed through an exchange or store '
          + 'credit. This does not restrict statutory refund rights relating to defective '
          + 'or non-conforming products.',
        ),
      ],
    },
    {
      title: 'How to Request a Return or Exchange',
      blocks: [
        paragraph('To initiate a request:'),
        list(
          'Contact us through our official email or WhatsApp.',
          'Provide your order number and purchase details.',
          'Explain the reason for your request.',
          'Provide photographs where necessary.',
          'Await our instructions before dispatching any item for return.',
        ),
      ],
    },
    {
      title: 'Refund Processing',
      blocks: [
        paragraph(
          'Approved refunds shall ordinarily be processed through the original payment '
          + 'method or another mutually agreed method. Processing and settlement times may '
          + 'depend on the relevant financial institution or payment provider.',
        ),
        paragraph('All refunds required by applicable law shall be honoured.'),
      ],
    },
    {
      title: 'Contact Us',
      blocks: [
        paragraph(`Email: ${CONTACT.email}`),
        paragraph(`Phone/WhatsApp: ${CONTACT.phone}`),
      ],
    },
  ],
  footnote: 'We appreciate your understanding and continued trust in Oyecouture Bespoke.',
}
