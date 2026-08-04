export interface ContactFeedbackRequest {
  readonly name: string
  readonly email: string
  readonly topic: string
  readonly rating: number
  readonly message: string
}

export interface ContactFeedbackResponse {
  readonly received: boolean
}
