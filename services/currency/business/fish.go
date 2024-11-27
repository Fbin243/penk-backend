package business

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"tenkhours/pkg/auth"
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

func NewFishBusiness(FishRepo *repo.FishRepo, CharactersRepo *coreRepo.CharactersRepo, ProfilesRepo *coreRepo.ProfilesRepo, redisClient *redis.Client) *FishBusiness {
	return &FishBusiness{
		FishRepo:       FishRepo,
		CharactersRepo: CharactersRepo,
		ProfilesRepo:   ProfilesRepo,
		RedisClient:    redisClient,
	}
}

func (biz *FishBusiness) GetFishByProfileID(ctx context.Context) (*model.Fish, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coreRepo.Profile)
	if !ok {
		return nil, errors.ErrorUnauthorized
	}

	log.Println("This IS ID", profile.ID)

	fish, err := biz.FishRepo.GetFishByProfileID(profile.ID)
	if err != nil {
		if err == mongo.ErrNoDocuments { // Check if fish document doesn't exist
			newFish := &repo.Fish{
				ID:        primitive.NewObjectID(),
				ProfileID: profile.ID,
				Gold:      0,
				Normal:    0,
			}
			fish, err = biz.FishRepo.CreateFish(newFish) // Create new fish document
			if err != nil {
				return nil, fmt.Errorf("failed to create fish data: %v", err)
			}
		} else {
			return nil, fmt.Errorf("failed to get fish data: %v", err)
		}
	}

	return biz.fishRepoToModel(fish), nil
}

// catchfish, This will be called in client when they finish their tracking
func (biz *FishBusiness) CatchFish(ctx context.Context, profileID primitive.ObjectID) (string, error) {
	rand.Seed(time.Now().UnixNano())

	// Load fish configurations
	fishConfigPath := os.Getenv("FISH_CONFIG_PATH")
	if fishConfigPath == "" {
		return "error", fmt.Errorf("FISH_CONFIG_PATH environment variable not set")
	}
	fishConfigs, err := config.LoadFishConfigs(fishConfigPath)

	fishType := "none"
	randomNumber := rand.Float64()

	for _, cfg := range fishConfigs {
		if randomNumber <= cfg.Rate {
			fishType = cfg.Type
			break
		}
		randomNumber -= cfg.Rate
	}

	if fishType == "none" {
		return "unlucky, next time :)))", nil
	}

	// Retrieve current fish data from Redis
	redisKey := fmt.Sprintf("fish:%s", profileID.Hex())
	if biz.RedisClient == nil {
		return "error", fmt.Errorf("Redis client is not initialized")
	}

	fishDataJSON, err := biz.RedisClient.Get(ctx, redisKey).Result()
	if err == redis.Nil {
		return "error", fmt.Errorf("fish data not found in Redis")
	} else if err != nil {
		return "error", fmt.Errorf("failed to retrieve fish data from Redis: %v", err)
	}

	// Parse existing fish data using encoding/json
	fish := &repo.Fish{}
	err = json.Unmarshal([]byte(fishDataJSON), fish)
	if err != nil {
		return "error", fmt.Errorf("failed to parse fish data: %v", err)
	}

	// Update fish counts
	switch fishType {
	case "normal":
		fish.Normal++
	case "gold":
		fish.Gold++
	}

	// Save updated fish data to Redis
	updatedFishDataJSON, err := json.Marshal(fish)
	if err != nil {
		return "error", fmt.Errorf("failed to serialize updated fish data: %v", err)
	}

	err = biz.RedisClient.Set(ctx, redisKey, updatedFishDataJSON, 4*time.Hour).Err()
	if err != nil {
		return "error", fmt.Errorf("failed to save updated fish data to Redis: %v", err)
	}

	return fishType, nil
}

func (biz *FishBusiness) UpdateFishFromRedis(fish *repo.Fish, profileID primitive.ObjectID) (bool, error) {
	currentFish, err := biz.FishRepo.GetFishByProfileID(profileID)
	if err != nil {
		return false, fmt.Errorf("failed to get fish data from DB: %v", err)
	}

	// Cộng dồn giá trị gold và normal
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

// Helper function to convert repo.Fish to model.Fish
func (biz *FishBusiness) fishRepoToModel(repoFish *repo.Fish) *model.Fish {
	gold := int(repoFish.Gold)
	normal := int(repoFish.Normal)

	return &model.Fish{
		ID:        repoFish.ID.Hex(),
		ProfileID: repoFish.ProfileID.Hex(),
		Gold:      &gold,
		Normal:    &normal,
	}
}

//	THOSE FUNCTIONS BELOW ARE FOR TRADING

func (biz *FishBusiness) UnlockMetrics(ctx context.Context, fishType string, characterID string) (bool, error) {
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
		if cfg.ItemType == "metric" && cfg.FishType == fishType { // Find matching config
			cost = cfg.Number
			increase = cfg.Increase
			foundConfig = true
			break
		}
	}

	if !foundConfig {
		return false, fmt.Errorf("invalid exchange configuration for metric and %s fish", fishType)
	}

	switch fishType { // Deduct fish based on loaded config
	case "normal":
		if fish.Normal < int32(cost) { // Correct type comparison
			return false, fmt.Errorf("not enough normal fish to trade")
		}
		fish.Normal -= int32(cost)
	case "gold":
		if fish.Gold < int32(cost) { // Correct type comparison
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

	character.LimitedMetricNumber += int32(increase) // Increment based on config

	if _, err := biz.CharactersRepo.UpdateCharacter(character); err != nil {
		return false, fmt.Errorf("failed to update metrics limited: %v", err)
	}

	if _, err := biz.FishRepo.UpdateFish(fish, profile.ID); err != nil {
		return false, fmt.Errorf("failed to update fish: %v", err)
	}

	return true, nil
}

func (biz *FishBusiness) BuySnapshots(ctx context.Context, fishType string) (bool, error) {
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
		if cfg.ItemType == "snapshot" && cfg.FishType == fishType {
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
	case "normal":
		if fish.Normal < int32(cost) {
			return false, fmt.Errorf("not enough normal fish to trade")
		}
		fish.Normal -= int32(cost)
	case "gold":
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

func (biz *FishBusiness) OnBoardNewCharacters(ctx context.Context, fishType string) (bool, error) {
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
		if cfg.ItemType == "character" && cfg.FishType == fishType {
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
	case "normal":
		if fish.Normal < int32(cost) {
			return false, fmt.Errorf("not enough normal fish to onboard new characters")
		}
		fish.Normal -= int32(cost)
	case "gold":
		if fish.Gold < int32(cost) {
			return false, fmt.Errorf("not enough gold fish to onboard new characters")
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
