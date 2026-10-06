package settings

import (
	"net/mail"
	"testing"
)

const configuredFrom = "Oyecouture Bespoke <orders@oyecouturebespoke.com>"

func TestComposeFromKeepsTheConfiguredHeaderWhenNoNameIsSet(t *testing.T) {
	if composed := ComposeFrom("", configuredFrom); composed != configuredFrom {
		t.Fatalf("ComposeFrom() = %q, want the configured header unchanged", composed)
	}
}

func TestComposeFromReplacesOnlyTheDisplayName(t *testing.T) {
	composed := ComposeFrom("Oyecouture Atelier", configuredFrom)

	want := `"Oyecouture Atelier" <orders@oyecouturebespoke.com>`
	if composed != want {
		t.Fatalf("ComposeFrom() = %q, want %q", composed, want)
	}
}

// The address is the one verified with the provider. A studio that types an address into the
// name field must not be able to send from it.
func TestComposeFromCannotChangeTheAddress(t *testing.T) {
	for _, displayName := range []string{
		"someone@elsewhere.com",
		"Oyecouture <someone@elsewhere.com>",
		"Name\r\nBcc: victim@example.com",
	} {
		composed := ComposeFrom(displayName, configuredFrom)

		parsed, err := parseForTest(composed)
		if err != nil {
			t.Fatalf("ComposeFrom(%q) produced an unparseable header %q: %v", displayName, composed, err)
		}
		if parsed != "orders@oyecouturebespoke.com" {
			t.Fatalf("ComposeFrom(%q) sent from %q, want the configured address", displayName, parsed)
		}
	}
}

// A header that cannot be parsed is left exactly as configured rather than rebuilt wrongly.
func TestComposeFromLeavesAnUnparseableHeaderAlone(t *testing.T) {
	const broken = "not-an-address"

	if composed := ComposeFrom("Oyecouture Bespoke", broken); composed != broken {
		t.Fatalf("ComposeFrom() = %q, want %q", composed, broken)
	}
}

// parseForTest reads back the address from a composed header.
func parseForTest(header string) (string, error) {
	parsed, err := mail.ParseAddress(header)
	if err != nil {
		return "", err
	}

	return parsed.Address, nil
}
