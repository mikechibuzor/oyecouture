export interface SizeChartRow {
  readonly size: string
  readonly ukSize: string
  readonly primary: string
  readonly secondary: string
}

export interface SizeChart {
  readonly id: string
  readonly title: string
  readonly primaryLabel: string
  readonly secondaryLabel: string
  readonly rows: readonly SizeChartRow[]
}

export interface MeasurementStep {
  readonly id: string
  readonly title: string
  readonly imagePath: string
  readonly imageAlt: string
  readonly instruction: string
  readonly allowance?: string
}

export interface MeasurementField {
  readonly id: string
  readonly label: string
}

export interface MeasurementWorksheet {
  readonly id: 'suit' | 'native'
  readonly label: string
  readonly fields: readonly MeasurementField[]
}

export const SIZE_CHARTS: readonly SizeChart[] = [
  {
    id: 'tops',
    title: 'Shirts, tops and jackets',
    primaryLabel: 'Chest',
    secondaryLabel: 'Waist',
    rows: [
      { size: 'S', ukSize: '34 / 36', primary: '91–94', secondary: '79–83' },
      { size: 'M', ukSize: '36 / 38', primary: '95–99', secondary: '84–88' },
      { size: 'L', ukSize: '38 / 40', primary: '100–105', secondary: '88–92' },
      { size: 'XL', ukSize: '42 / 44', primary: '106–111', secondary: '93–98' },
      { size: 'XXL', ukSize: '46 / 48', primary: '112–118', secondary: '99–105' },
      { size: 'XXXL', ukSize: '50 / 52', primary: '119–126', secondary: '106–113' },
    ],
  },
  {
    id: 'bottoms',
    title: 'Shorts and trousers',
    primaryLabel: 'Waist',
    secondaryLabel: 'Hip',
    rows: [
      { size: 'S', ukSize: '30 / 32', primary: '79–83', secondary: '82–86' },
      { size: 'M', ukSize: '32 / 34', primary: '84–88', secondary: '87–91' },
      { size: 'L', ukSize: '34 / 36', primary: '88–92', secondary: '92–96' },
      { size: 'XL', ukSize: '36 / 38', primary: '93–98', secondary: '97–112' },
      { size: 'XXL', ukSize: '40 / 42', primary: '99–105', secondary: '113–119' },
      { size: 'XXXL', ukSize: '42 / 44', primary: '106–113', secondary: '120–127' },
    ],
  },
]

export const MEASUREMENT_STEPS: readonly MeasurementStep[] = [
  {
    id: 'back',
    title: 'Back',
    imagePath: '/images/measurements/back.webp',
    imageAlt: 'Tape measure placed across the back of a tailoring mannequin',
    instruction: 'Measure from one shoulder bone across the back to the other shoulder bone.',
  },
  {
    id: 'sleeve',
    title: 'Sleeve',
    imagePath: '/images/measurements/sleeve.webp',
    imageAlt: 'Sleeve length measured from shoulder to wrist',
    instruction: 'Measure from the shoulder bone down to the wrist, keeping the arm relaxed.',
  },
  {
    id: 'chest',
    title: 'Chest',
    imagePath: '/images/measurements/chest.webp',
    imageAlt: 'Chest circumference measured on a tailoring mannequin',
    instruction: 'Wrap the tape around the fullest part of the chest and keep it level.',
    allowance: 'Add 3 inches to the exact body measurement.',
  },
  {
    id: 'stomach',
    title: 'Stomach',
    imagePath: '/images/measurements/stomach.webp',
    imageAlt: 'Stomach circumference measured on a tailoring mannequin',
    instruction: 'Wrap the tape around the fullest part of the stomach without pulling it tight.',
    allowance: 'Add 2 inches to the exact body measurement.',
  },
  {
    id: 'top-length',
    title: 'Top length',
    imagePath: '/images/measurements/top-length.webp',
    imageAlt: 'Jacket and native top length measured vertically',
    instruction: 'Measure from the base of the neck down to the preferred jacket or native top length.',
  },
  {
    id: 'biceps',
    title: 'Biceps / round sleeve',
    imagePath: '/images/measurements/biceps.webp',
    imageAlt: 'Upper arm circumference measured with a tape',
    instruction: 'Wrap the tape around the fullest part of the upper arm.',
    allowance: 'Add 1.5 inches to the exact body measurement.',
  },
  {
    id: 'neck',
    title: 'Neck',
    imagePath: '/images/measurements/neck.webp',
    imageAlt: 'Neck circumference measured at the collar line',
    instruction: 'Measure around the base of the neck where a shirt collar naturally sits.',
  },
  {
    id: 'waist',
    title: 'Waist',
    imagePath: '/images/measurements/waist.webp',
    imageAlt: 'Waist circumference measured on a tailoring mannequin',
    instruction: 'Wrap the tape around the natural waistline while standing comfortably.',
  },
  {
    id: 'thigh',
    title: 'Thigh',
    imagePath: '/images/measurements/thigh.webp',
    imageAlt: 'Thigh circumference measured on a tailoring mannequin',
    instruction: 'Measure around the fullest part of one thigh.',
    allowance: 'Add 2 inches to the exact body measurement.',
  },
  {
    id: 'trouser-length',
    title: 'Trouser length',
    imagePath: '/images/measurements/trouser-length.webp',
    imageAlt: 'Diagram showing outseam trouser length',
    instruction: 'Measure from the waistline down the outside of the leg to the preferred hem.',
  },
  {
    id: 'calf',
    title: 'Calf',
    imagePath: '/images/measurements/calf.webp',
    imageAlt: 'Diagram showing where to measure the calf',
    instruction: 'Wrap the tape around the fullest part of the calf.',
    allowance: 'Add 1 inch to the exact body measurement.',
  },
  {
    id: 'ankle',
    title: 'Ankle',
    imagePath: '/images/measurements/ankle.webp',
    imageAlt: 'Ankle circumference measured with a tape',
    instruction: 'Measure around the ankle at the point where the trouser opening will sit.',
  },
  {
    id: 'cap',
    title: 'Cap',
    imagePath: '/images/measurements/cap.webp',
    imageAlt: 'Head circumference measured for a traditional cap',
    instruction: 'Wrap the tape around the head where the cap will sit, just above the ears.',
  },
]

const SUIT_FIELDS: readonly MeasurementField[] = [
  { id: 'back', label: 'Back / shoulder' },
  { id: 'sleeve', label: 'Sleeve' },
  { id: 'round-sleeve', label: 'Round sleeve' },
  { id: 'chest', label: 'Chest' },
  { id: 'stomach', label: 'Stomach' },
  { id: 'jacket-length', label: 'Jacket length' },
  { id: 'trouser-length', label: 'Trouser length' },
  { id: 'waist', label: 'Waist' },
  { id: 'thigh', label: 'Thigh' },
  { id: 'calf', label: 'Calf' },
  { id: 'ankle', label: 'Ankle' },
  { id: 'neck', label: 'Neck' },
]

export const MEASUREMENT_WORKSHEETS: readonly MeasurementWorksheet[] = [
  { id: 'suit', label: 'Suit / tuxedo', fields: SUIT_FIELDS },
  {
    id: 'native',
    label: 'Native / kaftan / agbada',
    fields: [
      ...SUIT_FIELDS.map((field) =>
        field.id === 'jacket-length'
          ? { id: 'top-length', label: 'Top length' }
          : field),
      { id: 'cap', label: 'Cap' },
    ],
  },
]
