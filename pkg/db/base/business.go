package base

import (
	"context"
)

type IBaseBusiness[Entity, EntityInput, Filter, OrderBy any] interface {
	Count(ctx context.Context, filter *Filter) (int, error)
	Get(ctx context.Context, filter *Filter, orderBy *OrderBy, limit, offset *int) ([]Entity, error)
	Upsert(ctx context.Context, input *EntityInput) (*Entity, error)
	Delete(ctx context.Context, id string) (*Entity, error)
}
