package base

import "context"

// Interface of base repository to be embedded to concrete repository interfaces
type IBaseRepo[M any] interface {
	CountAll(ctx context.Context) (int64, error)
	InsertMany(ctx context.Context, ms []M) ([]M, error)
	InsertOne(ctx context.Context, m *M) (*M, error)
	FindAll(ctx context.Context) ([]M, error)
	FindByID(ctx context.Context, id string) (*M, error)
	UpdateByID(ctx context.Context, id string, m *M) (*M, error)
	DeleteByID(ctx context.Context, id string) (*M, error)
}
