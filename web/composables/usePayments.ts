import { API_ROUTES } from '~/constants/api'
import type {
  InitializePaymentResponse,
  PaymentSummary,
  VerifyPaymentResponse,
} from '~/types/payment'

export const usePayments = () => {
  const apiBaseUrl = useApiBaseUrl()

  const fetchPayment = (payToken: string): Promise<PaymentSummary> =>
    $fetch<PaymentSummary>(API_ROUTES.payment(payToken), {
      baseURL: apiBaseUrl,
    })

  const initializePayment = (payToken: string): Promise<InitializePaymentResponse> =>
    $fetch<InitializePaymentResponse>(API_ROUTES.paymentInitialize(payToken), {
      baseURL: apiBaseUrl,
      method: 'POST',
    })

  const verifyPayment = (
    payToken: string,
    reference: string,
  ): Promise<VerifyPaymentResponse> =>
    $fetch<VerifyPaymentResponse>(API_ROUTES.paymentVerify(payToken), {
      baseURL: apiBaseUrl,
      method: 'POST',
      body: { reference },
    })

  return { fetchPayment, initializePayment, verifyPayment }
}
