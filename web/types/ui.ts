// One choice in a UiSelect. `value` is what the model receives; `label` is what
// the reader sees. Generic so each call site keeps its own union type rather
// than widening to string.
export interface SelectOption<TValue extends string = string> {
  readonly value: TValue
  readonly label: string
  readonly disabled?: boolean
}
