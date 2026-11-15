import { PRIVACY_PAGE, RETURNS_PAGE, SHIPPING_PAGE, TERMS_PAGE } from '~/constants/policies'
import type { CarePageContent } from '~/types/policy'

// The four care pages now render the client's own policy documents rather than
// placeholder summaries.
//
// Shipping and returns used to be derived from PRODUCT_POLICIES so the delivery
// promise could not differ between the product page and the policy page. That
// no longer works: a product page cannot carry a full legal document. The
// figures those two share - production time, delivery estimates, the expedited
// surcharge - are instead held once in DELIVERY_TERMS and read by both, which
// keeps the guarantee the old arrangement was there for.
export const CARE_PAGES: Readonly<Record<string, CarePageContent>> = {
  shipping: SHIPPING_PAGE,
  returns: RETURNS_PAGE,
  legal: PRIVACY_PAGE,
  terms: TERMS_PAGE,
}
