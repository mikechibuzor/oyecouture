// Dates are shown in the owner's own timezone, short enough to scan in a table.
const orderDateFormat = new Intl.DateTimeFormat('en-NG', {
  day: 'numeric',
  month: 'short',
  hour: 'numeric',
  minute: '2-digit',
})

export const formatOrderDate = (isoDate: string): string => {
  const parsed = new Date(isoDate)

  return Number.isNaN(parsed.getTime()) ? 'Not available' : orderDateFormat.format(parsed)
}
