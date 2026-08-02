package paystack

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"testing"
)

func TestValidSignature(t *testing.T) {
	t.Parallel()

	body := []byte(`{"event":"charge.success","data":{"reference":"OYE-1001-test"}}`)
	secret := "test-secret"
	mac := hmac.New(sha512.New, []byte(secret))
	mac.Write(body)
	signature := hex.EncodeToString(mac.Sum(nil))

	if !validSignature(body, signature, secret) {
		t.Fatal("validSignature() rejected a valid signature")
	}
	if validSignature(body, signature, "wrong-secret") {
		t.Fatal("validSignature() accepted a signature with the wrong secret")
	}
}
