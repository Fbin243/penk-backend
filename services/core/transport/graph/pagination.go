package graph

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/samber/lo"
)

func paginate[Entity any](ctx context.Context,
	countFunc func() (int, error),
	findFunc func() ([]Entity, error),
) (int, []Entity, error) {
	fields := graphql.CollectAllFields(ctx)

	// Get the total count of items
	var totalCount int
	var err error
	if lo.Contains(fields, "totalCount") {
		totalCount, err = countFunc()
		if err != nil {
			return 0, nil, err
		}
	}

	// Get the items
	var items []Entity
	if lo.Contains(fields, "edges") {
		items, err = findFunc()
		if err != nil {
			return 0, nil, err
		}
	}

	return totalCount, items, nil
}
