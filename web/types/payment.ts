export interface PaymentOrderItem {
  readonly name: string
  readonly size?: string
  readonly quantity: number
  readonly lineKobo: number
}

export interface PaymentOrder {
  readonly orderNo: string
  readonly status: string
  readonly kind: string
  readonly customerName: string
  readonly email: string
  readonly phone: string
  readonly address: string
  readonly deliveryZone: string
  readonly subtotalKobo: number
  readonly deliveryFeeKobo: number
  readonly totalKobo: number
  readonly items: readonly PaymentOrderItem[]
  readonly createdAt: string
}

export interface PaymentSummary {
  readonly order: PaymentOrder
  readonly payable: boolean
  readonly message: string
}

export interface InitializePaymentResponse {
  readonly payment: {
    readonly authorizationUrl: string
    readonly accessCode: string
    readonly reference: string
  }
}

export interface VerifyPaymentResponse {
  readonly result: {
    readonly order: PaymentOrder
    readonly reference: string
    readonly status: string
    readonly paid: boolean
    readonly message: string
    // How they paid: card, bank_transfer, ussd. Absent until payment succeeds.
    readonly channel?: string
  }
}
