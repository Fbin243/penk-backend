package business

import (
	"context"
	"fmt"
	"time"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/errors"
	"tenkhours/services/currency/entity"

	rdb "tenkhours/pkg/db/redis"
)

type RewardBusiness struct {
	RewardRepo IRewardRepo
}

func NewRewardBusiness(RewardRepo IRewardRepo) *RewardBusiness {
	return &RewardBusiness{
		RewardRepo: RewardRepo,
	}
}

func (biz *RewardBusiness) GetRewardByProfileID(ctx context.Context) (*entity.Reward, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	reward, err := biz.RewardRepo.GetRewardByProfileID(ctx, authSession.ProfileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reward - err: %v", err)
	}

	return reward, nil
}

func (biz *RewardBusiness) ClaimDailyReward(ctx context.Context) (*entity.Reward, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	reward, err := biz.RewardRepo.GetRewardByProfileID(ctx, authSession.ProfileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reward - err: %v", err)
	}

	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")

	now := time.Now().In(loc)

	today7AM := time.Date(now.Year(), now.Month(), now.Day(), 7, 0, 0, 0, loc)
	yesterday7AM := today7AM.Add(-24 * time.Hour)

	// If reward exists, calculate streak
	if reward.StreakCount != 0 {
		lastClaim := reward.ClaimedAt.In(loc)

		if lastClaim.After(today7AM) {
			return nil, fmt.Errorf("already claimed today")
		}
		if lastClaim.After(yesterday7AM) {
			reward.StreakCount += 1
			if reward.StreakCount > 7 {
				reward.StreakCount = 1 // Reset after 7-day completion
			}
		} else {
			// Missed a day -> Reset streak
			reward.StreakCount = 1
		}

	} else {
		reward.StreakCount = 1
	}

	// Determine fish reward based on streak
	var fishReward int
	switch {
	case reward.StreakCount <= 2:
		fishReward = 1
	case reward.StreakCount >= 3 && reward.StreakCount <= 6:
		fishReward = 2
	case reward.StreakCount == 7:
		fishReward = 3
	}

	var updatedReward *entity.Reward

	updatedReward, err = biz.RewardRepo.UpdateReward(ctx, authSession.ProfileID, reward.StreakCount, int32(fishReward))
	if err != nil {
		return nil, fmt.Errorf("failed to update reward - err: %v", err)
	}

	return updatedReward, nil
}
