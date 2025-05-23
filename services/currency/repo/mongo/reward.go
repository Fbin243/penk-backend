package mongorepo

import (
	"context"
	"fmt"
	"time"

	"tenkhours/services/currency/entity"

	"tenkhours/pkg/db/base"
	mongodb "tenkhours/pkg/db/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RewardRepo struct {
	*mongodb.BaseRepo[entity.Reward, Reward]
}

func NewRewardRepo(db *mongo.Database) *RewardRepo {
	return &RewardRepo{mongodb.NewBaseRepo[entity.Reward, Reward](
		db.Collection(mongodb.RewardCollection),
		true,
	)}
}

func (r *RewardRepo) GetRewardByProfileID(ctx context.Context, profileID string) (*entity.Reward, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var reward entity.Reward
	filter := bson.M{"profile_id": mongodb.ToObjectID(profileID)}

	err := r.Collection.FindOne(ctx, filter).Decode(&reward)
	if err == nil {
		return &reward, nil
	}
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	newReward := entity.Reward{
		BaseEntity:  &base.BaseEntity{},
		ProfileID:   profileID,
		StreakCount: 0,
	}

	_, err = r.InsertOne(ctx, &newReward)
	if err != nil {
		return nil, fmt.Errorf("failed to create new reward: %w", err)
	}

	return &newReward, nil
}

func (r *RewardRepo) UpdateReward(ctx context.Context, profileID string, streakCount, fishCount int32) (*entity.Reward, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"profile_id": mongodb.ToObjectID(profileID)}
	update := bson.M{
		"$set": bson.M{
			"streak_count": streakCount,
			"claimed_at":   time.Now(),
		},
		"$inc": bson.M{
			"fish_count": fishCount,
		},
	}

	var updatedReward entity.Reward
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	err := r.Collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedReward)
	if err != nil {
		return nil, fmt.Errorf("failed to update daily reward: %w", err)
	}

	return &updatedReward, nil
}

func (r *RewardRepo) DeleteReward(ctx context.Context, profileID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"profile_id": mongodb.ToObjectID(profileID)}

	var deletedReward entity.Reward
	err := r.Collection.FindOneAndDelete(ctx, filter).Decode(&deletedReward)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return fmt.Errorf("failed to delete reward: %w", err)
	}

	return nil
}
