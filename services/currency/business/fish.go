package business

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db"
	"tenkhours/pkg/errors"
	coreRepo "tenkhours/services/core/repo"
	"tenkhours/services/currency/graph/model"
	"tenkhours/services/currency/repo"
	config "tenkhours/services/currency/utils"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FishBusiness struct {
	FishRepo       *repo.FishRepo
	ProfilesRepo   *coreRepo.ProfilesRepo
	CharactersRepo *coreRepo.CharactersRepo
	RedisClient    *redis.Client
}

type CatchFishResult struct {
	fishType string
	number   int32
}

func NewFishBusiness(FishRepo *repo.FishRepo, CharactersRepo *coreRepo.CharactersRepo, ProfilesRepo *coreRepo.ProfilesRepo, redisClient *redis.Client) *FishBusiness {
	return &FishBusiness{
		FishRepo:       FishRepo,
		CharactersRepo: CharactersRepo,
		ProfilesRepo:   ProfilesRepo,
		RedisClient:    redisClient,
	}
}

func (biz *FishBusiness) GetFishByProfileID(ctx context.Context) (*repo.Fish, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coreRepo.Profile)
	if !ok {
		return nil, errors.ErrorUnauthorized
	}

	fish, err := biz.FishRepo.GetFishByProfileID(profile.ID)

	// Check if there is no fish data or failed to get fish
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, fmt.Errorf("failed to get fish data: %v", err)
	}

	// Check if fish document doesn't exist
	if err == mongo.ErrNoDocuments {
		newFish := &repo.Fish{
			ID:        primitive.NewObjectID(),
			ProfileID: profile.ID,
			Gold:      0,
			Normal:    0,
		}

		// Create new fish document
		fish, err = biz.FishRepo.CreateFish(newFish)
		if err != nil {
			return nil, fmt.Errorf("failed to create fish data: %v", err)
		}
	}

	return fish, nil
}

// CatchFish used to catch the fish for a user
// This will be called in client when they finish their tracking sessions
func (biz *FishBusiness) CatchFish(ctx context.Context, profileID primitive.ObjectID) (*CatchFishResult, error) {
	//Create a random number generator base on fish caught rate
	rand.Seed(time.Now().UnixNano())

	// Load fish configurations
	fishConfigPath := os.Getenv("FISH_CONFIG_PATH")
	if fishConfigPath == "" {
		return nil, fmt.Errorf("FISH_CONFIG_PATH environment variable not set")
	}
	fishConfigs, err := config.LoadFishConfigs(fishConfigPath)

	// Create a random float number from 0.0 to 1.0
	randomNumber := rand.Float64()

	// this will check for all the case in config file
	var selectedFishConfig *config.FishConfig

	for _, cfg := range fishConfigs {
		if randomNumber <= cfg.Rate {
			selectedFishConfig = &cfg
			break
		}
		randomNumber -= cfg.Rate
	}

	// Retrieve current fish data from Redis
	fishKey := db.GetFishKey(profileID.Hex())

	fishDataJSON, err := biz.RedisClient.Get(ctx, fishKey).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("fish data not found in Redis")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve fish data from Redis: %v", err)
	}

	//Check for no fish caught
	if selectedFishConfig == nil {
		return &CatchFishResult{
			fishType: "none",
			number:   0,
		}, nil
	}

	// Parse existing fish data using encoding/json
	fish := &repo.Fish{}
	err = json.Unmarshal([]byte(fishDataJSON), fish)
	if err != nil {
		return nil, fmt.Errorf("failed to parse fish data: %v", err)
	}

	count := int32(0)

	// Update fish counts
	switch selectedFishConfig.Type {
	case "normal":
		fish.Normal += int32(selectedFishConfig.Number)
		count = int32(selectedFishConfig.Number)
	case "gold":
		fish.Gold += int32(selectedFishConfig.Number)
		count = int32(selectedFishConfig.Number)
	default:
		return nil, fmt.Errorf("unknown fish type")
	}

	// Save updated fish data to Redis
	updatedFishDataJSON, err := json.Marshal(fish)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize updated fish data: %v", err)
	}

	err = biz.RedisClient.Set(ctx, fishKey, updatedFishDataJSON, 24*time.Hour).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to save updated fish data to Redis: %v", err)
	}

	return &CatchFishResult{
		fishType: selectedFishConfig.Type,
		number:   count,
	}, nil
}

// Get Fish data from Redis to store in db
func (biz *FishBusiness) UpdateFishFromRedis(fish *repo.Fish, profileID primitive.ObjectID) (bool, error) {
	currentFish, err := biz.FishRepo.GetFishByProfileID(profileID)
	if err != nil {
		return false, fmt.Errorf("failed to get fish data from DB: %v", err)
	}

	//Accumulate the values of gold and normal
	if fish.Gold != 0 {
		currentFish.Gold += fish.Gold
	}
	if fish.Normal != 0 {
		currentFish.Normal += fish.Normal
	}

	_, err = biz.FishRepo.UpdateFish(currentFish, profileID)
	if err != nil {
		return false, fmt.Errorf("failed to save updated fish data to DB: %v", err)
	}

	return true, nil
}

//Those functions below are for trading

// Use fish to trade metrics
func (biz *FishBusiness) UnlockMetrics(ctx context.Context, fishType model.FishType, characterID string) (bool, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coreRepo.Profile)
	if !ok {
		return false, errors.ErrorUnauthorized
	}

	fish, err := biz.FishRepo.GetFishByProfileID(profile.ID)
	if err != nil {
		return false, fmt.Errorf("failed to find fish: %v", err)
	}

	exchangeConfigPath := os.Getenv("EXCHANGE_CONFIG_PATH") //load config
	if exchangeConfigPath == "" {
		return false, fmt.Errorf("EXCHANGE_CONFIG_PATH environment variable not set")
	}

	exchangeConfigs, err := config.LoadExchangeConfigs(exchangeConfigPath)

	cost := 0
	increase := 0
	foundConfig := false

	for _, cfg := range exchangeConfigs {
		if cfg.ItemType == "metric" && cfg.FishType == string(fishType) {
			cost = cfg.Number
			increase = cfg.Increase
			foundConfig = true
			break
		}
	}

	if !foundConfig {
		return false, fmt.Errorf("invalid exchange configuration for metric and %s fish", fishType)
	}

	switch fishType {
	case model.FishTypeNormal:
		if fish.Normal < int32(cost) {
			return false, fmt.Errorf("not enough normal fish to trade")
		}
		fish.Normal -= int32(cost)
	case model.FishTypeGold:
		if fish.Gold < int32(cost) {
			return false, fmt.Errorf("not enough gold fish to trade")
		}
		fish.Gold -= int32(cost)
	default:
		return false, fmt.Errorf("invalid fish type: %s", fishType)
	}

	characterOID, err := primitive.ObjectIDFromHex(characterID)
	if err != nil {
		return false, err
	}
	character, err := biz.CharactersRepo.GetCharacterByID(characterOID)
	if err != nil {
		return false, fmt.Errorf("failed to find character: %v", err)
	}

	//Increment based on config
	character.LimitedMetricNumber += int32(increase)

	if _, err := biz.CharactersRepo.UpdateCharacter(character); err != nil {
		return false, fmt.Errorf("failed to update metrics limited: %v", err)
	}

	if _, err := biz.FishRepo.UpdateFish(fish, profile.ID); err != nil {
		return false, fmt.Errorf("failed to update fish: %v", err)
	}

	return true, nil
}

// Use fish to trade snapshot
func (biz *FishBusiness) BuySnapshots(ctx context.Context, fishType model.FishType) (bool, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coreRepo.Profile)
	if !ok {
		return false, errors.ErrorUnauthorized
	}

	fish, err := biz.FishRepo.GetFishByProfileID(profile.ID)
	if err != nil {
		return false, fmt.Errorf("failed to find fish: %v", err)
	}

	exchangeConfigPath := os.Getenv("EXCHANGE_CONFIG_PATH")
	if exchangeConfigPath == "" {
		return false, fmt.Errorf("EXCHANGE_CONFIG_PATH environment variable not set")
	}

	exchangeConfigs, err := config.LoadExchangeConfigs(exchangeConfigPath)

	cost := 0
	increase := 0
	foundConfig := false

	for _, cfg := range exchangeConfigs {
		if cfg.ItemType == "snapshot" && cfg.FishType == string(fishType) {
			cost = cfg.Number
			increase = cfg.Increase
			foundConfig = true
			break
		}
	}

	if !foundConfig {
		return false, fmt.Errorf("invalid exchange configuration for snapshots and %s fish", fishType)
	}

	switch fishType {
	case model.FishTypeNormal:
		if fish.Normal < int32(cost) {
			return false, fmt.Errorf("not enough normal fish to trade")
		}
		fish.Normal -= int32(cost)
	case model.FishTypeGold:
		if fish.Gold < int32(cost) {
			return false, fmt.Errorf("not enough gold fish to trade")
		}
		fish.Gold -= int32(cost)
	default:
		return false, fmt.Errorf("invalid fish type: %s", fishType)
	}

	profileData, err := biz.ProfilesRepo.GetProfileByFirebaseUID(profile.FirebaseUID)
	if err != nil {
		return false, fmt.Errorf("failed to get profile: %v", err)
	}

	profileData.AvailableSnapshots += int32(increase)

	if _, err := biz.ProfilesRepo.UpdateProfile(profileData); err != nil {
		return false, fmt.Errorf("failed to update available snapshots: %v", err)
	}

	if _, err := biz.FishRepo.UpdateFish(fish, profile.ID); err != nil {
		return false, fmt.Errorf("failed to update fish: %v", err)
	}

	return true, nil
}

// Use fish to unclock new character
func (biz *FishBusiness) UnclockNewCharacters(ctx context.Context, fishType model.FishType) (bool, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coreRepo.Profile)
	if !ok {
		return false, errors.ErrorUnauthorized
	}

	fish, err := biz.FishRepo.GetFishByProfileID(profile.ID)
	if err != nil {
		return false, fmt.Errorf("failed to find fish: %v", err)
	}

	exchangeConfigPath := os.Getenv("EXCHANGE_CONFIG_PATH")
	if exchangeConfigPath == "" {
		return false, fmt.Errorf("EXCHANGE_CONFIG_PATH environment variable not set")
	}

	exchangeConfigs, err := config.LoadExchangeConfigs(exchangeConfigPath)

	cost := 0
	increase := 0
	foundConfig := false

	for _, cfg := range exchangeConfigs {
		if cfg.ItemType == "character" && cfg.FishType == string(fishType) {
			cost = cfg.Number
			increase = cfg.Increase
			foundConfig = true
			break
		}
	}

	if !foundConfig {
		return false, fmt.Errorf("invalid exchange configuration for characters and %s fish", fishType)
	}

	switch fishType {
	case model.FishTypeNormal:
		if fish.Normal < int32(cost) {
			return false, fmt.Errorf("not enough normal fish to trade")
		}
		fish.Normal -= int32(cost)
	case model.FishTypeGold:
		if fish.Gold < int32(cost) {
			return false, fmt.Errorf("not enough gold fish to trade")
		}
		fish.Gold -= int32(cost)
	default:
		return false, fmt.Errorf("invalid fish type: %s", fishType)
	}

	profileData, err := biz.ProfilesRepo.GetProfileByFirebaseUID(profile.FirebaseUID)
	if err != nil {
		return false, fmt.Errorf("failed to get profile: %v", err)
	}

	// increase the limited char number
	profileData.LimitedCharacterNumber += int32(increase)

	if _, err := biz.ProfilesRepo.UpdateProfile(profileData); err != nil {
		return false, fmt.Errorf("failed to update character count: %v", err)
	}

	if _, err := biz.FishRepo.UpdateFish(fish, profile.ID); err != nil {
		return false, fmt.Errorf("failed to update fish: %v", err)
	}

	return true, nil
}
