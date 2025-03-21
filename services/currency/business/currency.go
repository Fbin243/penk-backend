package business

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"tenkhours/pkg/db/base"
	"tenkhours/services/currency/entity"

	config "tenkhours/services/currency/utils"
)

type CurrencyBusiness struct {
	FishRepo   IFishRepo
	CoreClient ICoreClient
}

func NewCurrencyBusiness(FishRepo IFishRepo, coreClient ICoreClient) *CurrencyBusiness {
	return &CurrencyBusiness{
		FishRepo:   FishRepo,
		CoreClient: coreClient,
	}
}

func (biz *CurrencyBusiness) GetFish(ctx context.Context, profileID string) (*entity.Fish, error) {
	fish, err := biz.FishRepo.GetFishByProfileID(ctx, profileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get fish: %v", err)
	}

	return fish, nil
}

func (biz *CurrencyBusiness) CreateFish(ctx context.Context, profileID string) (*entity.Fish, error) {
	fish := &entity.Fish{
		BaseEntity: &base.BaseEntity{},
		ProfileID:  profileID,
		Normal:     0,
		Gold:       0,
	}

	_, err := biz.FishRepo.InsertOne(ctx, fish)
	if err != nil {
		return nil, fmt.Errorf("failed to create fish: %v", err)
	}

	return fish, nil
}

func (biz *CurrencyBusiness) UpdateFish(ctx context.Context, fish *entity.Fish) (*entity.Fish, error) {
	currentFish, err := biz.FishRepo.GetFishByProfileID(ctx, fish.ProfileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get fish - err2: %v", err)
	}

	currentFish.Gold += fish.Gold
	currentFish.Normal += fish.Normal

	updatedFish, err := biz.FishRepo.UpdateFishByProfileID(ctx, fish.ProfileID, currentFish)
	if err != nil {
		return nil, fmt.Errorf("failed to update fish: %v", err)
	}

	return updatedFish, nil
}

// CatchFish used to catch the fish for a user
// This will be called in client when they finish their tracking sessions
func (biz *CurrencyBusiness) CatchFish(ctx context.Context) (*entity.CatchFishResult, error) {
	// Create a random number generator base on fish caught rate
	rand.Seed(time.Now().UnixNano())

	// Load fish configurations
	fishConfigPath := os.Getenv("FISH_CONFIG_PATH")
	if fishConfigPath == "" {
		return nil, fmt.Errorf("FISH_CONFIG_PATH environment variable not set")
	}
	fishConfigs, err := config.LoadFishConfigs(fishConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load fish config %v", err)
	}

	// Create a random float number from 0.0 to 1.0
	randomNumber := rand.Float64()

	// this will check for all the case in config file
	var selectedFishConfig *config.FishConfig
	cumulativeRate := 0.0

	for _, cfg := range fishConfigs {
		cumulativeRate += cfg.Rate
		if randomNumber <= cumulativeRate {
			selectedFishConfig = &cfg
			break
		}
	}

	// Check for no fish caught
	if selectedFishConfig == nil {
		return &entity.CatchFishResult{
			FishType: entity.FishTypeNone,
			Number:   0,
		}, nil
	}

	count := int32(0)

	// Update fish counts
	switch entity.FishType(selectedFishConfig.Type) {
	case entity.FishTypeNormal:
		count = int32(selectedFishConfig.Number)
	case entity.FishTypeGold:
		count = int32(selectedFishConfig.Number)
	default:
		return nil, fmt.Errorf("unknown fish type: %s", selectedFishConfig.Type)
	}

	return &entity.CatchFishResult{
		FishType: entity.FishType(selectedFishConfig.Type),
		Number:   count,
	}, nil
}

func (biz *CurrencyBusiness) DeleteFish(ctx context.Context, profileID string) (*entity.Fish, error) {
	fish, err := biz.FishRepo.DeleteFishByProfileID(ctx, profileID)
	if err != nil {
		return nil, err
	}

	return fish, nil
}
