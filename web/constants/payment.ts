export const PAYMENT_COPY = {
  breadcrumb: 'Payment',
  heading: 'Secure payment',
  description: 'Review your order and continue to Paystack to complete payment.',
  loading: 'Loading payment details...',
  loadFailed: 'This payment link could not be loaded.',
  notPayable: 'This order is not ready for payment.',
  alreadyPaid: 'This order has already been paid.',
  payAction: 'Pay securely',
  starting: 'Opening Paystack...',
  verifying: 'Confirming payment...',
  verified: 'Payment confirmed',
  pending: 'Payment is not complete yet',
  failed: 'The payment could not be verified.',
  // Bank transfers and USSD settle after Paystack sends the customer back, so
  // the first check routinely finds nothing. Saying so plainly beats the old
  // behaviour of declaring failure on a payment that lands seconds later.
  settling: 'Waiting for your bank',
  settlingBody:
    'Transfers and USSD payments take a moment to clear. This page is watching for it and will update on its own.',
  settlingSlow: 'Still waiting for your bank',
  settlingSlowBody:
    'This is taking longer than usual. Your money is safe. If it has left your account, the order will confirm itself and we will email you the receipt.',
  checkAgain: 'Check again',
  order: 'Order',
  customer: 'Customer',
  items: 'Items',
  total: 'Total',
  retry: 'Try again',
  continueShopping: 'Continue shopping',
} as const

// Copy for the post-payment confirmation page, per the client's approved design.
export const CONFIRMATION_COPY = {
  eyebrow: 'Payment confirmed',
  greeting: (firstName: string): string => `Thank you, ${firstName}`,
  intro:
    'Your order is confirmed and now with Oyecouture. Below is everything that was paid for, along with where each piece goes next.',
  orderReference: 'Order reference',
  orderDate: 'Order date',
  paymentMethod: 'Payment method',
  amountPaid: 'Amount paid',
  yourPieces: 'Your pieces',
  itemCount: (count: number): string => `${count} ${count === 1 ? 'item' : 'items'}`,
  whatHappensNext: 'What happens next',
  paymentSummary: 'Payment summary',
  subtotal: 'Subtotal',
  tailoring: 'Tailoring & finishing',
  tailoringValue: 'Included',
  delivery: 'Delivery',
  deliveryComplimentary: 'Complimentary',
  chargeNote: 'Charged once. Your card details are never stored by Oyecouture.',
  deliveryDetails: 'Delivery details',
  recipient: 'Recipient',
  address: 'Address',
  contact: 'Contact',
  estimatedDelivery: 'Estimated delivery',
  downloadReceipt: 'Download receipt',
  contactAtelier: 'Contact Oyecouture',
  footerNote:
    'Every piece is finished by hand, so small variations in cloth and stitch are part of the work.',
  footerQuestion: 'Questions about sizing or delivery? Oyecouture answers within a day.',
} as const

// Receipt copy. Wording follows the client's approved email template so the
// printed receipt and the emailed confirmation read as the same document.
export const RECEIPT_COPY = {
  title: 'Receipt',
  banner: 'Complimentary tailoring consultation · Worldwide delivery',
  eyebrow: 'Payment received',
  heading: 'Receipt',
  printAction: 'Print or save as PDF',
  // Said plainly because the browser's own dialog is the thing that saves the
  // file, and a customer who expects a download otherwise thinks it failed.
  printHint: 'Choose "Save as PDF" as the destination in the print dialog to keep a copy.',
  back: 'Back to your order',
  orderDetails: 'Order details',
  orderReference: 'Order reference',
  orderDate: 'Order date',
  payment: 'Payment',
  yourPieces: 'Your pieces',
  subtotal: 'Subtotal',
  tailoring: 'Tailoring consultation',
  complimentary: 'Complimentary',
  delivery: 'Delivery',
  totalPaid: 'Total paid',
  deliveringTo: 'Delivering to',
  signOff: 'We are glad to be dressing you.',
  brandName: 'Oyecouture Bespoke',
  returnsNote:
    'Ready-to-wear pieces may be exchanged within 7 days of delivery in original condition. Made-to-measure and bespoke commissions are cut to your measurements and are final sale, though we will adjust the fit at no charge.',
  address: 'Oyecouture Bespoke · 6 Force Road, Tafawa Balewa Square (TBS), Lagos, Nigeria',
  receiptFor: (orderNo: string): string => `Receipt for order ${orderNo}.`,
} as const

