package orders

import (
	"strings"
	"testing"
)

// TestNormalizeFilterClampsUntrustedInput verifies a hand-typed query string
// can neither ask for an unbounded read nor error the page out.
func TestNormalizeFilterClampsUntrustedInput(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		given ListFilter
		want  ListFilter
	}{
		{
			name:  "empty filter takes the default page",
			given: ListFilter{},
			want:  ListFilter{Limit: defaultPageSize},
		},
		{
			name:  "unknown status is dropped rather than rejected",
			given: ListFilter{Status: "not_a_status"},
			want:  ListFilter{Limit: defaultPageSize},
		},
		{
			name:  "unknown kind is dropped rather than rejected",
			given: ListFilter{Kind: "haute-couture"},
			want:  ListFilter{Limit: defaultPageSize},
		},
		{
			name:  "known status and kind survive",
			given: ListFilter{Status: StatusPaid, Kind: KindBespoke},
			want:  ListFilter{Status: StatusPaid, Kind: KindBespoke, Limit: defaultPageSize},
		},
		{
			name:  "oversized limit is clamped",
			given: ListFilter{Limit: 99999},
			want:  ListFilter{Limit: maximumPageSize},
		},
		{
			name:  "negative limit falls back to the default",
			given: ListFilter{Limit: -5},
			want:  ListFilter{Limit: defaultPageSize},
		},
		{
			name:  "negative offset is floored",
			given: ListFilter{Offset: -20},
			want:  ListFilter{Limit: defaultPageSize},
		},
		{
			name:  "search is trimmed",
			given: ListFilter{Search: "  OYE-1016 "},
			want:  ListFilter{Search: "OYE-1016", Limit: defaultPageSize},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if got := normalizeFilter(testCase.given); got != testCase.want {
				t.Fatalf("normalizeFilter() = %+v, want %+v", got, testCase.want)
			}
		})
	}
}

// TestNormalizeFilterTruncatesLongSearch keeps an oversized term out of the query.
func TestNormalizeFilterTruncatesLongSearch(t *testing.T) {
	t.Parallel()

	normalized := normalizeFilter(ListFilter{Search: strings.Repeat("a", maximumSearchSize*3)})

	if len(normalized.Search) != maximumSearchSize {
		t.Fatalf("search length = %d, want %d", len(normalized.Search), maximumSearchSize)
	}
}

// TestEscapeLikePattern stops a wildcard in the search box from matching every
// order — searching "%" must find orders containing a literal percent sign.
func TestEscapeLikePattern(t *testing.T) {
	t.Parallel()

	testCases := []struct{ name, given, want string }{
		{name: "percent is literal", given: "100%", want: `100\%`},
		{name: "underscore is literal", given: "a_b", want: `a\_b`},
		{name: "backslash is escaped first", given: `a\b`, want: `a\\b`},
		{name: "plain text is untouched", given: "OYE-1016", want: "OYE-1016"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if got := escapeLikePattern(testCase.given); got != testCase.want {
				t.Fatalf("escapeLikePattern(%q) = %q, want %q", testCase.given, got, testCase.want)
			}
		})
	}
}
