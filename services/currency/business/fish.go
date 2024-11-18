package business

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/errors"
	coreRepo "tenkhours/services/core/repo"
	"tenkhours/services/currency/graph/model"
	"tenkhours/services/currency/repo"
	config "tenkhours/services/currency/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FishBusiness struct {
	FishRepo       *repo.FishRepo
	ProfilesRepo   *coreRepo.ProfilesRepo
	CharactersRepo *coreRepo.CharactersRepo
}

func NewFishBusiness(FishRepo *repo.FishRepo, CharactersRepo *coreRepo.CharactersRepo, ProfilesRepo *coreRepo.ProfilesRepo) *FishBusiness {
	return &FishBusiness{
		FishRepo:       FishRepo,
		CharactersRepo: CharactersRepo,
		ProfilesRepo:   ProfilesRepo,
	}
}

func (biz *FishBusiness) GetFishByProfileID(ctx context.Context) (*model.Fish, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coreRepo.Profile)
	if !ok {
		return nil, errors.ErrorUnauthorized
	}

	repoFish, err := biz.FishRepo.GetFishByProfileID(profile.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find fish repo: %v", err)
	}

	return biz.fishRepoToModel(repoFish), nil
}

// catchfish, This will be called in client when they finish their tracking
func (biz *FishBusiness) CatchFish(ctx context.Context) (*model.Fish, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coreRepo.Profile)
	if !ok {
		return nil, errors.ErrorUnauthorized
	}

	rand.Seed(time.Now().UnixNano())

	// Load fish configurations
	fishConfigs, err := config.LoadFishConfigs("fish_config.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to load fish configurations: %v", err)
	}

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
		return nil, fmt.Errorf("unlucky, next time :)))")
	}

	fish, err := biz.FishRepo.GetFishByProfileID(profile.ID)
	if err != nil {
		// If no fish document exists, create a new one
		fish = &repo.Fish{
			ID:        primitive.NewObjectID(),
			ProfileID: profile.ID,
			Counts:    repo.FishCounts{}, // Initialize empty counts
		}
	}

	switch fishType {
	case "normal":
		fish.Counts.Normal++
	case "gold":
		fish.Counts.Gold++
	}

	if _, err = biz.FishRepo.UpdateFish(fish, profile.ID); err != nil {
		return nil, fmt.Errorf("failed to update fish counts: %v", err)
	}

	return biz.fishRepoToModel(fish), nil
}

// Helper function to convert repo.Fish to model.Fish
func (biz *FishBusiness) fishRepoToModel(repoFish *repo.Fish) *model.Fish {
	gold := int(repoFish.Counts.Gold)
	normal := int(repoFish.Counts.Normal)

	counts := model.FishCounts{ // Create FishCounts
		Gold:   &gold,
		Normal: &normal,
	}
	return &model.Fish{
		ID:        repoFish.ID.Hex(),
		ProfileID: repoFish.ProfileID.Hex(),
		Counts:    &counts,
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

	cost := int32(0)
	switch fishType {
	case "normal":
		cost = 3
		if fish.Counts.Normal < cost {
			return false, fmt.Errorf("not enough normal fish to trade")
		}
		fish.Counts.Normal -= cost
	case "gold":
		cost = 1
		if fish.Counts.Gold < cost {
			return false, fmt.Errorf("not enough gold fish to trade")
		}
		fish.Counts.Gold -= cost
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

	character.LimitedMetricNumber++
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

	cost := int32(0)
	switch fishType {
	case "normal":
		cost = 3
		if fish.Counts.Normal < cost {
			return false, fmt.Errorf("not enough normal fish to trade")
		}
		fish.Counts.Normal -= cost
	case "gold":
		cost = 1
		if fish.Counts.Gold < cost {
			return false, fmt.Errorf("not enough gold fish to trade")
		}
		fish.Counts.Gold -= cost
	default:
		return false, fmt.Errorf("invalid fish type: %s", fishType)
	}

	profileData, err := biz.ProfilesRepo.GetProfileByFirebaseUID(profile.FirebaseUID)
	if err != nil {
		return false, fmt.Errorf("failed to get profile: %v", err)
	}

	profileData.AvailableSnapshots++
	if _, err := biz.ProfilesRepo.UpdateProfile(profileData); err != nil {
		return false, fmt.Errorf("failed to update available snapshots: %v", err)
	}

	if _, err := biz.FishRepo.UpdateFish(fish, profile.ID); err != nil {
		return false, fmt.Errorf("failed to update fish: %v", err)
	}

	return true, nil
}
