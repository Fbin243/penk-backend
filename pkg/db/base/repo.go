package base

import "context"

// Interface of base repository to be embedded to concrete repository interfaces
type IBaseRepo[M any] interface {
	// Count(ctx context.Context, filter any) (int, error)
	// Exists(ctx context.Context, filter any) (bool, error)

	InsertOne(ctx context.Context, m *M) (*M, error)
	InsertMany(ctx context.Context, ms []M) ([]M, error)

	FindByID(ctx context.Context, id string) (*M, error)
	FindByIDs(ctx context.Context, ids []string) ([]M, error)
	// FindOne(ctx context.Context, filter any) (*M, error)
	// FindMany(ctx context.Context, filter any) ([]M, error)

	FindAndUpdateByID(ctx context.Context, id string, m *M) (*M, error)
	FindAndUpdateByIDs(ctx context.Context, m []M) ([]M, error)
	UpdateByID(ctx context.Context, id string, m *M) error
	// UpdateOne(ctx context.Context, filter, update any) error
	// UpdateMany(ctx context.Context, filter, update any) error

	FindAndDeleteByID(ctx context.Context, id string) (*M, error)
	DeleteByID(ctx context.Context, id string) error
	// DeleteOne(ctx context.Context, filter any) error
	// DeleteMany(ctx context.Context, filter any) error

	// AggregateQuery(ctx context.Context, pipeline any) ([]M, error)
	// AggregateCount(ctx context.Context, pipeline any) (int, error)
}
