package business

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	"tenkhours/pkg/errors"
	"tenkhours/services/currency/entity"

	rdb "tenkhours/pkg/db/redis"

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

	currentFish.Gold = fish.Gold
	currentFish.Normal = fish.Normal

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

// Those functions below are for trading
// Use fish to trade metrics
func (biz *CurrencyBusiness) BuyMetrics(ctx context.Context, fishType entity.FishType, characterID string) (bool, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return false, errors.Unauthorized()
	}

	fish, err := biz.FishRepo.GetFishByProfileID(ctx, authSession.ProfileID)
	if err != nil {
		return false, fmt.Errorf("failed to find fish: %v", err)
	}

	exchangeConfigPath := os.Getenv("EXCHANGE_CONFIG_PATH") // load config
	if exchangeConfigPath == "" {
		return false, fmt.Errorf("EXCHANGE_CONFIG_PATH environment variable not set")
	}

	exchangeConfigs, err := config.LoadExchangeConfigs(exchangeConfigPath)
	if err != nil {
		return false, fmt.Errorf("failed to load exchange config %v", err)
	}

	cost := 0
	// increase := 0
	foundConfig := false

	for _, cfg := range exchangeConfigs {
		if cfg.ItemType == "metric" && cfg.FishType == string(fishType) {
			cost = cfg.Number
			// increase = cfg.Increase
			foundConfig = true
			break
		}
	}

	if !foundConfig {
		return false, fmt.Errorf("invalid exchange configuration for metric and %s fish", fishType)
	}

	switch fishType {
	case entity.FishTypeNormal:
		if fish.Normal < int32(cost) {
			return false, fmt.Errorf("not enough normal fish to trade")
		}
		fish.Normal -= int32(cost)
	case entity.FishTypeGold:
		if fish.Gold < int32(cost) {
			return false, fmt.Errorf("not enough gold fish to trade")
		}
		fish.Gold -= int32(cost)
	default:
		return false, fmt.Errorf("invalid fish type: %s", fishType)
	}

	// // Buy metric
	// err = biz.CoreClient.BuyItem(ctx, lo.ToPtr(authSession.ProfileID), lo.ToPtr(characterID), nil, entity.ItemTypeMetric, int32(increase))
	// if err != nil {
	// 	return false, fmt.Errorf("failed to buy metric: %v", err)
	// }

	if _, err := biz.FishRepo.UpdateFishByProfileID(ctx, authSession.ProfileID, fish); err != nil {
		return false, fmt.Errorf("failed to update fish: %v", err)
	}

	return true, nil
}

// Use fish to trade snapshot
func (biz *CurrencyBusiness) BuySnapshots(ctx context.Context, fishType entity.FishType) (bool, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return false, errors.Unauthorized()
	}

	fish, err := biz.FishRepo.GetFishByProfileID(ctx, authSession.ProfileID)
	if err != nil {
		return false, fmt.Errorf("failed to find fish: %v", err)
	}

	exchangeConfigPath := os.Getenv("EXCHANGE_CONFIG_PATH")
	if exchangeConfigPath == "" {
		return false, fmt.Errorf("EXCHANGE_CONFIG_PATH environment variable not set")
	}

	exchangeConfigs, err := config.LoadExchangeConfigs(exchangeConfigPath)
	if err != nil {
		return false, fmt.Errorf("failed to load exchange config %v", err)
	}

	cost := 0
	// increase := 0
	foundConfig := false

	for _, cfg := range exchangeConfigs {
		if cfg.ItemType == "snapshot" && cfg.FishType == string(fishType) {
			cost = cfg.Number
			// increase = cfg.Increase
			foundConfig = true
			break
		}
	}

	if !foundConfig {
		return false, fmt.Errorf("invalid exchange configuration for snapshots and %s fish", fishType)
	}

	switch fishType {
	case entity.FishTypeNormal:
		if fish.Normal < int32(cost) {
			return false, fmt.Errorf("not enough normal fish to trade")
		}
		fish.Normal -= int32(cost)
	case entity.FishTypeGold:
		if fish.Gold < int32(cost) {
			return false, fmt.Errorf("not enough gold fish to trade")
		}
		fish.Gold -= int32(cost)
	default:
		return false, fmt.Errorf("invalid fish type: %s", fishType)
	}

	// profileData, err := biz.ProfileRepo.FindByID(authSession.ProfileID)
	// if err != nil {
	// 	return false, fmt.Errorf("failed to get profile: %v", err)
	// }

	// profileData.AvailableSnapshots += int32(increase)

	// if _, err := biz.ProfileRepo.UpdateByID(profileData.ID, profileData); err != nil {
	// 	return false, fmt.Errorf("failed to update available snapshots: %v", err)
	// }

	// TODO: @Fbin243 add grpc buy item later (after product team finalize snapshot)

	if _, err := biz.FishRepo.UpdateFishByProfileID(ctx, authSession.ProfileID, fish); err != nil {
		return false, fmt.Errorf("failed to update fish: %v", err)
	}

	return true, nil
}

// Use fish to unclock new character
func (biz *CurrencyBusiness) BuyCharacters(ctx context.Context, fishType entity.FishType) (bool, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return false, errors.Unauthorized()
	}

	fish, err := biz.FishRepo.GetFishByProfileID(ctx, authSession.ProfileID)
	if err != nil {
		return false, fmt.Errorf("failed to find fish: %v", err)
	}

	exchangeConfigPath := os.Getenv("EXCHANGE_CONFIG_PATH")
	if exchangeConfigPath == "" {
		return false, fmt.Errorf("EXCHANGE_CONFIG_PATH environment variable not set")
	}

	exchangeConfigs, err := config.LoadExchangeConfigs(exchangeConfigPath)
	if err != nil {
		return false, fmt.Errorf("failed to load exchange config %v", err)
	}

	cost := 0
	// increase := 0
	foundConfig := false

	for _, cfg := range exchangeConfigs {
		if cfg.ItemType == "character" && cfg.FishType == string(fishType) {
			cost = cfg.Number
			// increase = cfg.Increase
			foundConfig = true
			break
		}
	}

	if !foundConfig {
		return false, fmt.Errorf("invalid exchange configuration for characters and %s fish", fishType)
	}

	switch fishType {
	case entity.FishTypeNormal:
		if fish.Normal < int32(cost) {
			return false, fmt.Errorf("not enough normal fish to trade")
		}
		fish.Normal -= int32(cost)
	case entity.FishTypeGold:
		if fish.Gold < int32(cost) {
			return false, fmt.Errorf("not enough gold fish to trade")
		}
		fish.Gold -= int32(cost)
	default:
		return false, fmt.Errorf("invalid fish type: %s", fishType)
	}

	// Buy character
	// err = biz.CoreClient.BuyItem(ctx, lo.ToPtr(authSession.ProfileID), nil, nil, entity.ItemTypeCharacter, int32(increase))
	// if err != nil {
	// 	return false, fmt.Errorf("failed to buy metric: %v", err)
	// }

	if _, err := biz.FishRepo.UpdateFishByProfileID(ctx, authSession.ProfileID, fish); err != nil {
		return false, fmt.Errorf("failed to update fish: %v", err)
	}

	return true, nil
}
