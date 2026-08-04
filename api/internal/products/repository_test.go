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

// TestBuildAdminListQueryIncludesArchivedProducts verifies the owner-facing base query.
func TestBuildAdminListQueryIncludesArchivedProducts(t *testing.T) {
	t.Parallel()

	query, countQuery, queryArguments := buildAdminListQuery(AdminListFilter{
		Limit:  25,
		Offset: 50,
	})

	if strings.Contains(query, "product.active = TRUE") {
		t.Fatalf("admin query should not force active-only products: %q", query)
	}
	if !strings.Contains(query, "ORDER BY product.active DESC") {
		t.Fatalf("expected active products first, got %q", query)
	}
	if countQuery != "SELECT COUNT(*) FROM products AS product\nWHERE TRUE" {
		t.Fatalf("unexpected count query %q", countQuery)
	}
	if len(queryArguments) != 2 || queryArguments[0] != 25 || queryArguments[1] != 50 {
		t.Fatalf("unexpected paging arguments: %#v", queryArguments)
	}
}

// TestBuildAdminListQueryUsesParameterizedFilters verifies safe admin filters.
func TestBuildAdminListQueryUsesParameterizedFilters(t *testing.T) {
	t.Parallel()

	active := false
	query, countQuery, queryArguments := buildAdminListQuery(AdminListFilter{
		Section: SectionTrad,
		Active:  &active,
		Query:   `agbada_%`,
		Limit:   10,
		Offset:  20,
	})

	for _, expectedQueryFragment := range []string{
		"product.section = $1",
		"product.active = $2",
		"product.slug ILIKE $3 ESCAPE '\\'",
		"LIMIT $4 OFFSET $5",
	} {
		if !strings.Contains(query, expectedQueryFragment) {
			t.Fatalf("expected query to contain %q, got %q", expectedQueryFragment, query)
		}
	}
	for _, expectedCountFragment := range []string{
		"product.section = $1",
		"product.active = $2",
		"product.description ILIKE $3 ESCAPE '\\'",
	} {
		if !strings.Contains(countQuery, expectedCountFragment) {
			t.Fatalf("expected count query to contain %q, got %q", expectedCountFragment, countQuery)
		}
	}
	if len(queryArguments) != 5 {
		t.Fatalf("expected 5 query arguments, got %d", len(queryArguments))
	}
	if queryArguments[0] != SectionTrad || queryArguments[1] != false ||
		queryArguments[2] != `%agbada\_\%%` || queryArguments[3] != 10 ||
		queryArguments[4] != 20 {
		t.Fatalf("unexpected query arguments: %#v", queryArguments)
	}
}
