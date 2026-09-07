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

// stageConstraints are the two CHECK constraints that must both list every
// fulfilment stage. They are separate tables written in separate migrations,
// which is exactly how they drifted apart.
var stageConstraints = map[string]*regexp.Regexp{
	"orders_fulfillment_stage_valid": regexp.MustCompile(
		`(?s)CONSTRAINT orders_fulfillment_stage_valid CHECK \((.*?)\)\s*[,;]`,
	),
	"order_fulfillment_events_stage_valid": regexp.MustCompile(
		`(?s)CONSTRAINT order_fulfillment_events_stage_valid CHECK \((.*?)\)\s*[,;]`,
	),
}

// TestStagesAreAllowedByEveryStageConstraint closes the gap that made moving an
// order to out-for-delivery impossible in production.
//
// Migration 000018 added the stage and widened the constraint on orders, but
// not the one on order_fulfillment_events, written two migrations earlier. A
// transition updates the order and inserts the event in one transaction, so the
// event's CHECK failed, the transaction rolled back, and the admin saw only
// "order tracking could not be loaded".
//
// The email-type test above is the same lesson learned once already. Testing
// only the constraint that broke last time is how the next one gets through, so
// this asserts every stage against every stage constraint.
func TestStagesAreAllowedByEveryStageConstraint(t *testing.T) {
	t.Parallel()

	// Restated rather than imported: the tracking package cannot be imported
	// here without a cycle, and requiring a deliberate edit in both places is
	// what makes the omission visible.
	stages := []string{
		"payment_confirmed",
		"measurements_approved",
		"preparing_order",
		"in_production",
		"quality_check",
		"ready_for_dispatch",
		"dispatched",
		"out_for_delivery",
		"delivered",
	}

	names, err := fs.Glob(migrationFiles, "*.up.sql")
	if err != nil {
		t.Fatalf("glob migrations: %v", err)
	}
	sort.Strings(names)

	for constraintName, pattern := range stageConstraints {
		// Later migrations redefine the constraint, so the last match wins the
		// same way it does when they are applied in order.
		allowed := ""
		for _, name := range names {
			content, err := fs.ReadFile(migrationFiles, name)
			if err != nil {
				t.Fatalf("read %s: %v", name, err)
			}
			if match := pattern.FindStringSubmatch(string(content)); match != nil {
				allowed = match[1]
			}
		}
		if allowed == "" {
			t.Fatalf("no %s constraint found in the migrations", constraintName)
		}

		for _, stage := range stages {
			if !strings.Contains(allowed, "'"+stage+"'") {
				t.Fatalf(
					"stage %q is not permitted by %s; add it in a new migration "+
						"or the transition will roll back with an opaque internal error",
					stage, constraintName,
				)
			}
		}
	}
}
