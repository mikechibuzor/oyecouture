package payments

import (
	"errors"
	"testing"
)

// TestVerificationMatchesOrder covers the money guard. Every case here is a way
// an order could be marked paid for the wrong amount, which is the worst thing
// this system can get wrong: the piece ships and the money never arrived.
func TestVerificationMatchesOrder(t *testing.T) {
	t.Parallel()

	const total int64 = 8_300_000

	matching := func() (Attempt, Order, Verification) {
		return Attempt{AmountKobo: total, Currency: CurrencyNGN},
			Order{TotalKobo: total},
			Verification{AmountKobo: total, Currency: CurrencyNGN}
	}

	testCases := []struct {
		name    string
		mutate  func(*Attempt, *Order, *Verification)
		wantErr error
	}{
		{
			name:   "exact match is accepted",
			mutate: func(*Attempt, *Order, *Verification) {},
		},
		{
			name: "underpayment is rejected",
			mutate: func(_ *Attempt, _ *Order, verification *Verification) {
				verification.AmountKobo = total - 100
			},
			wantErr: ErrPaymentAmountMismatch,
		},
		{
			name: "paying one kobo short is rejected",
			mutate: func(_ *Attempt, _ *Order, verification *Verification) {
				verification.AmountKobo = total - 1
			},
			wantErr: ErrPaymentAmountMismatch,
		},
		{
			name: "overpayment is rejected rather than silently accepted",
			mutate: func(_ *Attempt, _ *Order, verification *Verification) {
				verification.AmountKobo = total + 100
			},
			wantErr: ErrPaymentAmountMismatch,
		},
		{
			name: "zero payment is rejected",
			mutate: func(_ *Attempt, _ *Order, verification *Verification) {
				verification.AmountKobo = 0
			},
			wantErr: ErrPaymentAmountMismatch,
		},
		{
			name: "an attempt recorded for the wrong amount is rejected",
			mutate: func(attempt *Attempt, _ *Order, _ *Verification) {
				attempt.AmountKobo = total - 500
			},
			wantErr: ErrPaymentAmountMismatch,
		},
		{
			name: "an order repriced after the attempt is rejected",
			mutate: func(_ *Attempt, order *Order, _ *Verification) {
				order.TotalKobo = total + 500
			},
			wantErr: ErrPaymentAmountMismatch,
		},
		{
			name: "a foreign currency is rejected even when the number matches",
			mutate: func(attempt *Attempt, _ *Order, verification *Verification) {
				attempt.Currency = "USD"
				verification.Currency = "USD"
			},
			wantErr: ErrPaymentCurrency,
		},
		{
			name: "a currency the attempt did not use is rejected",
			mutate: func(_ *Attempt, _ *Order, verification *Verification) {
				verification.Currency = "USD"
			},
			wantErr: ErrPaymentCurrency,
		},
		{
			name: "a missing currency is rejected",
			mutate: func(attempt *Attempt, _ *Order, verification *Verification) {
				attempt.Currency = ""
				verification.Currency = ""
			},
			wantErr: ErrPaymentCurrency,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			attempt, order, verification := matching()
			testCase.mutate(&attempt, &order, &verification)

			err := verificationMatchesOrder(attempt, order, verification)

			if !errors.Is(err, testCase.wantErr) {
				t.Fatalf("verificationMatchesOrder() = %v, want %v", err, testCase.wantErr)
			}
		})
	}
}
