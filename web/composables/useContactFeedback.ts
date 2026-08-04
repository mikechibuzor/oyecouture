import { API_ROUTES } from '~/constants/api'
import type { ContactFeedbackRequest, ContactFeedbackResponse } from '~/types/contact'

export const useContactFeedback = () => {
  const apiBaseUrl = useApiBaseUrl()

  const submitContactFeedback = (
    request: ContactFeedbackRequest,
  ): Promise<ContactFeedbackResponse> =>
    $fetch<ContactFeedbackResponse>(API_ROUTES.contact, {
      baseURL: apiBaseUrl,
      method: 'POST',
      body: request,
    })

  return { submitContactFeedback }
}
