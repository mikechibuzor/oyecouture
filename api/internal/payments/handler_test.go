package payments

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"testing"
)

func TestValidPaystackSignature(t *testing.T) {
	t.Parallel()

	body := []byte(`{"event":"charge.success","data":{"reference":"OYE-1001-test"}}`)
	secret := "test-secret"
	mac := hmac.New(sha512.New, []byte(secret))
	mac.Write(body)
	signature := hex.EncodeToString(mac.Sum(nil))

	if !validPaystackSignature(body, signature, secret) {
		t.Fatal("validPaystackSignature() rejected a valid signature")
	}
	if validPaystackSignature(body, signature, "wrong-secret") {
		t.Fatal("validPaystackSignature() accepted a signature with the wrong secret")
	}
}
