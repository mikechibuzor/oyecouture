package products

import (
	"strings"
	"testing"
)

// TestBuildListQueryUsesActiveProductsAndParameterizedFilters verifies safe filter SQL.
func TestBuildListQueryUsesActiveProductsAndParameterizedFilters(t *testing.T) {
	t.Parallel()

	query, queryArguments := buildListQuery(ListFilter{
		Section: SectionCasualsMen,
		Query:   `50%_\`,
	})

	for _, expectedQueryFragment := range []string{
		"product.active = TRUE",
		"product.section = $1",
		"product.name ILIKE $2 ESCAPE '\\'",
		"GROUP BY product.id",
	} {
		if !strings.Contains(query, expectedQueryFragment) {
			t.Fatalf("expected query to contain %q, got %q", expectedQueryFragment, query)
		}
	}

	if len(queryArguments) != 2 {
		t.Fatalf("expected 2 query arguments, got %d", len(queryArguments))
	}
	if queryArguments[0] != SectionCasualsMen {
		t.Fatalf("expected section argument %q, got %q", SectionCasualsMen, queryArguments[0])
	}
	if queryArguments[1] != `%50\%\_\\%` {
		t.Fatalf("expected escaped search pattern, got %q", queryArguments[1])
	}
}

// TestBuildListQueryWithoutFiltersKeepsAnEmptyArgumentList verifies the base listing query.
func TestBuildListQueryWithoutFiltersKeepsAnEmptyArgumentList(t *testing.T) {
	t.Parallel()

	query, queryArguments := buildListQuery(ListFilter{})

	if !strings.Contains(query, "product.active = TRUE") {
		t.Fatalf("expected active product filter, got %q", query)
	}
	if len(queryArguments) != 0 {
		t.Fatalf("expected no query arguments, got %d", len(queryArguments))
	}
}
