package migrations

import (
	"io/fs"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// emailTypeConstraint finds the last definition of outbound_emails_type_valid
// across the migrations, which is the one that ends up on the live table.
var emailTypeConstraint = regexp.MustCompile(
	`(?s)ADD CONSTRAINT outbound_emails_type_valid CHECK \((.*?)\);`,
)

// TestEmailTypesAreAllowedByTheDatabase closes the gap that broke
// out-for-delivery in production: the email type was added as a Go constant,
// the CHECK constraint was not widened to match, and enqueuing one rolled back
// the whole fulfilment transition with an opaque internal error.
//
// The email package cannot be imported here without a cycle, so the constants
// are restated. That is the point: adding a Type* constant should require a
// deliberate edit in both places, and this test names the missing one.
func TestEmailTypesAreAllowedByTheDatabase(t *testing.T) {
	t.Parallel()

	emailTypes := []string{
		"order_confirmation",
		"payment_receipt",
		"fulfillment_update",
		"order_dispatched",
		"order_out_for_delivery",
		"order_delivered",
	}

	names, err := fs.Glob(migrationFiles, "*.up.sql")
	if err != nil {
		t.Fatalf("glob migrations: %v", err)
	}
	sort.Strings(names)

	// Later migrations redefine the constraint, so the last match wins the same
	// way it does when they are applied in order.
	allowed := ""
	for _, name := range names {
		content, err := fs.ReadFile(migrationFiles, name)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		if match := emailTypeConstraint.FindStringSubmatch(string(content)); match != nil {
			allowed = match[1]
		}
	}
	if allowed == "" {
		t.Fatal("no outbound_emails_type_valid constraint found in the migrations")
	}

	for _, emailType := range emailTypes {
		if !strings.Contains(allowed, "'"+emailType+"'") {
			t.Fatalf(
				"email type %q is not permitted by outbound_emails_type_valid; "+
					"add it in a new migration or enqueuing one will roll back its transaction",
				emailType,
			)
		}
	}
}
