package orders

import (
	"strings"
	"testing"
	"time"
)

// TestPeriodStartResolvesEachWindow pins the summary-strip windows. "Today"
// runs from local midnight, not 24 hours back, so an order placed this morning
// still counts at 11pm.
func TestPeriodStartResolvesEachWindow(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.August, 2, 15, 30, 0, 0, time.UTC)

	testCases := []struct {
		name   string
		period MetricsPeriod
		want   *time.Time
	}{
		{
			name:   "today starts at local midnight",
			period: PeriodToday,
			want:   timePointer(time.Date(2026, time.August, 2, 0, 0, 0, 0, time.UTC)),
		},
		{
			name:   "last week looks back seven days",
			period: PeriodLastWeek,
			want:   timePointer(time.Date(2026, time.July, 26, 15, 30, 0, 0, time.UTC)),
		},
		{
			name:   "fourteen days",
			period: PeriodLast14Days,
			want:   timePointer(time.Date(2026, time.July, 19, 15, 30, 0, 0, time.UTC)),
		},
		{
			name:   "ninety days crosses months",
			period: PeriodLast90Days,
			want:   timePointer(time.Date(2026, time.May, 4, 15, 30, 0, 0, time.UTC)),
		},
		{name: "all time has no lower bound", period: PeriodAllTime, want: nil},
		{name: "unknown period falls back to all time", period: "last-decade", want: nil},
		{name: "empty period falls back to all time", period: "", want: nil},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got := periodStart(testCase.period, now)

			if testCase.want == nil {
				if got != nil {
					t.Fatalf("periodStart() = %v, want nil", got)
				}
				return
			}

			if got == nil || !got.Equal(*testCase.want) {
				t.Fatalf("periodStart() = %v, want %v", got, *testCase.want)
			}
		})
	}
}

func timePointer(value time.Time) *time.Time { return &value }

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

// TestNormalizeFilterHandlesDateBounds covers the two date rules: "to" names a
// day the owner wants included, and a backwards range must not silently return
// nothing.
func TestNormalizeFilterHandlesDateBounds(t *testing.T) {
	t.Parallel()

	t.Run("to runs to the end of the named day", func(t *testing.T) {
		t.Parallel()

		to := time.Date(2026, time.August, 2, 0, 0, 0, 0, time.Local)
		normalized := normalizeFilter(ListFilter{To: &to})

		want := time.Date(2026, time.August, 3, 0, 0, 0, 0, time.Local)
		if normalized.To == nil || !normalized.To.Equal(want) {
			t.Fatalf("To = %v, want %v", normalized.To, want)
		}
	})

	t.Run("backwards range drops the upper bound", func(t *testing.T) {
		t.Parallel()

		from := time.Date(2026, time.August, 10, 0, 0, 0, 0, time.Local)
		to := time.Date(2026, time.August, 1, 0, 0, 0, 0, time.Local)
		normalized := normalizeFilter(ListFilter{From: &from, To: &to})

		if normalized.To != nil {
			t.Fatalf("To = %v, want nil for a backwards range", normalized.To)
		}
		if normalized.From == nil || !normalized.From.Equal(from) {
			t.Fatalf("From = %v, want %v", normalized.From, from)
		}
	})
}

// TestParseFilterDate treats anything unparseable as absent, so a hand-edited
// URL still returns a list.
func TestParseFilterDate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		given   string
		wantNil bool
	}{
		{name: "valid date", given: "2026-08-02"},
		{name: "empty is absent", given: "", wantNil: true},
		{name: "wrong format is absent", given: "02/08/2026", wantNil: true},
		{name: "nonsense is absent", given: "yesterday", wantNil: true},
		{name: "impossible date is absent", given: "2026-13-45", wantNil: true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got := ParseFilterDate(testCase.given)

			if testCase.wantNil && got != nil {
				t.Fatalf("ParseFilterDate(%q) = %v, want nil", testCase.given, got)
			}
			if !testCase.wantNil && got == nil {
				t.Fatalf("ParseFilterDate(%q) = nil, want a date", testCase.given)
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
