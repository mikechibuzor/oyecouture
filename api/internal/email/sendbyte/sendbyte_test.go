package sendbyte

import "testing"

// TestMessageIDExtractsKnownResponseShapes covers the response forms SendByte-style APIs use.
func TestMessageIDExtractsKnownResponseShapes(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		responseBody string
		want         string
	}{
		{name: "top level id", responseBody: `{"id":"email_123"}`, want: "email_123"},
		{name: "top level message id", responseBody: `{"message_id":"email_456"}`, want: "email_456"},
		{name: "camel message id", responseBody: `{"messageId":"email_789"}`, want: "email_789"},
		{name: "nested data id", responseBody: `{"data":{"id":"email_nested"}}`, want: "email_nested"},
		{name: "invalid json", responseBody: `not-json`, want: ""},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got := messageID([]byte(testCase.responseBody))
			if got != testCase.want {
				t.Fatalf("messageID() = %q, want %q", got, testCase.want)
			}
		})
	}
}
