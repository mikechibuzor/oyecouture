import { MEASUREMENT_WORKSHEETS } from '~/constants/sizeGuide'
import type { CartLine } from '~/types/cart'
import type { MeasurementProfile, MeasurementUnit } from '~/types/order'
import { ProductSection } from '~/types/product'

const CLIPBOARD_PREFIX = 'OYECOUTURE_MEASUREMENTS_V1:'
const MEASUREMENT_UNIT: MeasurementUnit = 'in'

const SUIT_CATEGORY_TERMS = ['suit', 'tuxedo', 'blazer', 'two-piece', 'three-piece']
const NATIVE_CATEGORY_TERMS = ['agbada', 'aso-oke', 'kaftan', 'buba', 'native']

export interface MeasurementRequirements {
  readonly profile: MeasurementProfile
  readonly fields: readonly { readonly id: string, readonly label: string }[]
}

interface MeasurementClipboardPayload {
  readonly version: 1
  readonly profile: MeasurementProfile
  readonly unit: MeasurementUnit
  readonly values: Record<string, number>
}

const allMeasurementFields = MEASUREMENT_WORKSHEETS
  .flatMap((worksheet) => worksheet.fields)
  .filter((field, fieldIndex, fields) =>
    fields.findIndex((candidate) => candidate.id === field.id) === fieldIndex,
  )

const knownMeasurementIds = new Set(allMeasurementFields.map((field) => field.id))

export const measurementFieldLabels = Object.fromEntries(
  allMeasurementFields.map((field) => [field.id, field.label]),
) as Readonly<Record<string, string>>

export const fieldsForMeasurementProfile = (
  profile: MeasurementProfile,
): MeasurementRequirements['fields'] => {
  if (profile === 'mixed') {
    return allMeasurementFields
  }

  return MEASUREMENT_WORKSHEETS.find((worksheet) => worksheet.id === profile)?.fields ?? []
}

const profileForCartLine = (line: CartLine): MeasurementProfile | null => {
  if (line.hasSizes) {
    return null
  }

  if (line.section === ProductSection.Suits) {
    return 'suit'
  }
  if (line.section === ProductSection.Trad) {
    return 'native'
  }

  const description = `${line.category} ${line.name}`.toLowerCase()
  if (SUIT_CATEGORY_TERMS.some((term) => description.includes(term))) {
    return 'suit'
  }
  if (NATIVE_CATEGORY_TERMS.some((term) => description.includes(term))) {
    return 'native'
  }

  return 'mixed'
}

export const getMeasurementRequirements = (
  lines: readonly CartLine[],
): MeasurementRequirements | null => {
  const profiles = new Set(lines.map(profileForCartLine).filter(Boolean))
  if (profiles.size === 0) {
    return null
  }

  const profile: MeasurementProfile = profiles.size === 1
    ? [...profiles][0] as MeasurementProfile
    : 'mixed'

  return { profile, fields: fieldsForMeasurementProfile(profile) }
}

export const serializeMeasurementClipboard = (
  profile: MeasurementProfile,
  values: Record<string, number>,
): string => `${CLIPBOARD_PREFIX}${JSON.stringify({
  version: 1,
  profile,
  unit: MEASUREMENT_UNIT,
  values,
} satisfies MeasurementClipboardPayload)}`

export const parseMeasurementClipboard = (
  clipboardText: string,
): MeasurementClipboardPayload | null => {
  const trimmedText = clipboardText.trim()
  if (!trimmedText.startsWith(CLIPBOARD_PREFIX)) {
    return null
  }

  try {
    const payload = JSON.parse(trimmedText.slice(CLIPBOARD_PREFIX.length)) as Partial<MeasurementClipboardPayload>
    if (
      payload.version !== 1
      || !['suit', 'native', 'mixed'].includes(payload.profile ?? '')
      || payload.unit !== MEASUREMENT_UNIT
      || typeof payload.values !== 'object'
      || payload.values === null
    ) {
      return null
    }

    const values = Object.fromEntries(
      Object.entries(payload.values).filter(([fieldId, value]) =>
        knownMeasurementIds.has(fieldId)
        && typeof value === 'number'
        && Number.isFinite(value)
        && value > 0,
      ),
    )

    return Object.keys(values).length > 0
      ? { version: 1, profile: payload.profile as MeasurementProfile, unit: MEASUREMENT_UNIT, values }
      : null
  } catch {
    return null
  }
}
