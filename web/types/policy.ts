// The care and policy pages render the client's legal documents, which are
// paragraphs and bullet lists in sequence - a section often reads "here is the
// rule", then a list of conditions, then a qualifying sentence. The previous
// shape held one title and one string, so a list could only be flattened into
// prose, and the qualifier after it had nowhere to go.

export enum CareBlockKind {
  Paragraph = 'paragraph',
  List = 'list',
}

export interface CareParagraph {
  readonly kind: CareBlockKind.Paragraph
  readonly text: string
}

export interface CareList {
  readonly kind: CareBlockKind.List
  readonly items: readonly string[]
}

export type CareBlock = CareParagraph | CareList

export interface CareSection {
  readonly title: string
  readonly blocks: readonly CareBlock[]
}

export interface CarePageContent {
  readonly breadcrumb: string
  readonly heading: string
  readonly intro: string
  readonly sections: readonly CareSection[]
  // Closing line beneath the sections, where the document has one.
  readonly footnote?: string
}

// Shorthand so the policy documents read as content rather than as structure.
export const paragraph = (text: string): CareParagraph => ({
  kind: CareBlockKind.Paragraph,
  text,
})

export const list = (...items: readonly string[]): CareList => ({
  kind: CareBlockKind.List,
  items,
})
