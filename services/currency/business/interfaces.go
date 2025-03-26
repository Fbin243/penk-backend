package business

import (
	"context"

	"tenkhours/pkg/db/base"
	"tenkhours/services/currency/entity"
)

type IRewardBusiness interface {
	GetRewardByProfileID(ctx context.Context) (*entity.Reward, error)
	ClaimDailyReward(ctx context.Context) (*entity.Reward, error)
}

type IRewardRepo interface {
	base.IBaseRepo[entity.Reward]
	GetRewardByProfileID(ctx context.Context, profileID string) (*entity.Reward, error)
	UpdateReward(ctx context.Context, profileID string, streakCount, fishCount int32) (*entity.Reward, error)
}

type ICoreClient interface{}
