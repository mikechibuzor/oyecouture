export const CONFIRM_COPY = {
  breadcrumb: 'Confirm order',
  heading: 'Confirm your order',
  description:
    'Leave your details and we will pick this up with you on WhatsApp — sizing, fabric and delivery are agreed there before any payment is taken.',
  contactHeading: 'Your details',
  name: 'Full name',
  email: 'Email address',
  phone: 'Phone number',
  phoneHint: 'Include your country code, e.g. +234',
  deliveryHeading: 'Delivery',
  address: 'Delivery address',
  zone: 'Delivery area',
  zonePlaceholder: 'Select a delivery area',
  zoneHint: 'Delivery fees for bespoke pieces are confirmed on WhatsApp.',
  zoneRequired: 'Please choose a delivery area.',
  measurementsHeading: 'Already know your measurements?',
  measurementsHint:
    'Optional — share them now and we can skip straight to fabric and delivery.',
  measurementsToggle: 'Add my measurements',
  notes: 'Anything else we should know',
  summaryHeading: 'Your order',
  subtotal: 'Subtotal',
  submit: 'Confirm order',
  submitting: 'Confirming…',
  emptyBag: 'Your bag is empty.',
  emptyBagAction: 'Browse the collection',
  successEyebrow: 'Order received',
  successHeading: 'We have your order',
  successBody:
    'Quote this reference when we speak. Continue on WhatsApp and we will confirm sizing, fabric and delivery.',
  successAction: 'Continue on WhatsApp',
  errorFallback: 'We could not place that order. Please try again.',
} as const

// cm unless stated; the design's measurement set.
export const MEASUREMENT_FIELDS = [
  { key: 'neckCm', label: 'Neck (cm)' },
  { key: 'chestCm', label: 'Chest (cm)' },
  { key: 'waistCm', label: 'Waist (cm)' },
  { key: 'hipsCm', label: 'Hips (cm)' },
  { key: 'sleeveCm', label: 'Sleeve (cm)' },
  { key: 'inseamCm', label: 'Inseam (cm)' },
  { key: 'lengthCm', label: 'Length (cm)' },
  { key: 'heightCm', label: 'Height (cm)' },
  { key: 'weightKg', label: 'Weight (kg)' },
] as const
