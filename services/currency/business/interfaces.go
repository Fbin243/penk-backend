package business

import (
	"context"

	"tenkhours/pkg/db/base"
	"tenkhours/services/currency/entity"
)

type ICurrencyBusiness interface {
	GetFish(ctx context.Context, profileID string) (*entity.Fish, error)
	CreateFish(ctx context.Context, profileID string) (*entity.Fish, error)
	CatchFish(ctx context.Context) (*entity.CatchFishResult, error)
	UpdateFish(ctx context.Context, fish *entity.Fish) (*entity.Fish, error)
	DeleteFish(ctx context.Context, profileID string) (*entity.Fish, error)
}

type IFishRepo interface {
	base.IBaseRepo[entity.Fish]
	GetFishByProfileID(ctx context.Context, profileID string) (*entity.Fish, error)
	UpdateFishByProfileID(ctx context.Context, profileID string, fish *entity.Fish) (*entity.Fish, error)
	DeleteFishByProfileID(ctx context.Context, profileID string) (*entity.Fish, error)
}

type ICoreClient interface{}
