package business

import (
	"context"
	"fmt"
	"tenkhours/pkg/auth"
	"tenkhours/pkg/errors"
	coreRepo "tenkhours/services/core/repo"
	"tenkhours/services/currency/graph/model"
	"tenkhours/services/currency/repo"
)

type FishBusiness struct {
	FishRepo     *repo.FishRepo
	ProfilesRepo *coreRepo.ProfilesRepo
}

func NewFishBusiness(FishRepo *repo.FishRepo) *FishBusiness {
	return &FishBusiness{
		FishRepo: FishRepo,
	}
}

func (biz *FishBusiness) GetFishByProfileID(ctx context.Context) (*model.Fish, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coreRepo.Profile)
	if !ok {
		return nil, errors.ErrorUnauthorized
	}

	repoFish, err := biz.FishRepo.GetFishByProfileID(profile.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find character: %v", err)
	}

	// Convert type repo.fish to type model.fish
	tempNum := int(repoFish.Numbers)
	tempType := repoFish.Type
	modelFish := &model.Fish{
		ID:        repoFish.ID.Hex(),
		ProfileID: repoFish.ProfileID.Hex(),
		Numbers:   &tempNum,
		Type:      &tempType,
	}

	return modelFish, nil
}

func (biz *FishBusiness) CreateFish(ctx context.Context, input model.FishInput) (*model.Fish, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coreRepo.Profile)
	if !ok {
		return nil, errors.ErrorUnauthorized
	}

	newFish := &repo.Fish{
		ProfileID: profile.ID,
		Numbers:   0,
		Type:      *input.Type,
	}

	insertedFish, err := biz.FishRepo.CreateFish(newFish)
	if err != nil {
		return nil, fmt.Errorf("failed to create fish: %v", err)
	}

	modelFish := &model.Fish{
		ID:        insertedFish.ID.Hex(),
		ProfileID: insertedFish.ProfileID.Hex(),
		Numbers:   int(insertedFish.Numbers),
		Type:      insertedFish.Type,
	}

	return modelFish, nil
}
