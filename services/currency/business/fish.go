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

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func (biz *FishBusiness) GetFishByProfileID(ctx context.Context, fishType string) (*model.Fish, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coreRepo.Profile)
	if !ok {
		return nil, errors.ErrorUnauthorized
	}

	repoFish, err := biz.FishRepo.GetFishByProfileID(profile.ID, fishType)
	if err != nil {
		return nil, fmt.Errorf("failed to find fish repo: %v", err)
	}

	// Convert type repo.fish to type model.fish
	tempNum := int(repoFish.Numbers) //create a pointer to the int value directly
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

	if *input.Type != "gold" && *input.Type != "normal" { //checktype of input
		return nil, fmt.Errorf("Wrong type of fish, accepted gold or normal fish")
	}

	// Check that if the fish repo is existed or not
	existingFish, err := biz.FishRepo.GetFishByProfileID(profile.ID, *input.Type)
	if err == nil && existingFish != nil {
		return nil, fmt.Errorf("fish with type %s already exists for this profile", *input.Type)
	} else if err != nil && err != mongo.ErrNoDocuments {
		return nil, fmt.Errorf("failed to check existing fish: %v", err)
	}

	// Create new fish
	newFish := &repo.Fish{
		ProfileID: profile.ID,
		Numbers:   0,
		Type:      *input.Type,
	}

	insertedFish, err := biz.FishRepo.CreateFish(newFish)
	if err != nil {
		return nil, fmt.Errorf("failed to create fish: %v", err)
	}

	tmpNumber := int(insertedFish.Numbers) //create a pointer to the int value directly
	tmpType := insertedFish.Type

	modelFish := &model.Fish{
		ID:        insertedFish.ID.Hex(),
		ProfileID: insertedFish.ProfileID.Hex(),
		Numbers:   &tmpNumber,
		Type:      &tmpType,
	}

	return modelFish, nil
}

func (biz *FishBusiness) UpdateFish(ctx context.Context, input model.FishInput) (*model.Fish, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coreRepo.Profile)
	if !ok {
		return nil, errors.ErrorUnauthorized
	}

	// Dereference input.Number to get the int value
	tmpNewFish := int32(*input.Number)

	// Check if the number is negative
	if tmpNewFish < 0 {
		return nil, fmt.Errorf("Number must be non-negative")
	}

	updatedFish := &repo.Fish{
		ProfileID: profile.ID,
		Numbers:   tmpNewFish,
		Type:      *input.Type,
	}

	insertedFish, err := biz.FishRepo.UpdateFish(updatedFish, profile.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update fish: %v", err)
	}

	tmpNumber := int(insertedFish.Numbers)
	tmpType := insertedFish.Type

	modelFish := &model.Fish{
		ID:        insertedFish.ID.Hex(),
		ProfileID: insertedFish.ProfileID.Hex(),
		Numbers:   &tmpNumber,
		Type:      &tmpType,
	}

	return modelFish, nil
}

// catchfish, This will be called in client when they finish their tracking
func (biz *FishBusiness) CatchFish(ctx context.Context) (*model.Fish, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coreRepo.Profile)
	if !ok {
		return nil, errors.ErrorUnauthorized
	}

	rand.Seed(time.Now().UnixNano()) // Seed the random number generator

	// Determine fish type based on probability
	var fishType = "none"
	var fishFlag int32 //Check for condition

	randomNumber := rand.Float64()
	fmt.Println("hahduadfhsudfhusdfhusdfhufdshdfsuo", randomNumber)

	switch {
	case randomNumber <= 0.3015: // 30.15% chance of 1 Fish
		fishType, fishFlag = "normal", 1
	case randomNumber <= 0.3225: // 2.10% chance of 2 Fish
		fishType, fishFlag = "normal", 2
	case randomNumber <= 0.3278: // 0.0053% chance of 1 Golden Fish
		fishType, fishFlag = "gold", 1
	default:
		return nil, fmt.Errorf("Unlucky, next time :)))")
	}

	// Update fish count using the helper function
	insertedFish, err := biz.updateFishCount(profile.ID, fishType, fishFlag)
	if err != nil {
		return nil, err
	}

	// Prepare response model
	tmpNumber := int(insertedFish.Numbers)
	tmpType := insertedFish.Type

	modelFish := &model.Fish{
		ID:        insertedFish.ID.Hex(),
		ProfileID: insertedFish.ProfileID.Hex(),
		Numbers:   &tmpNumber,
		Type:      &tmpType,
	}

	return modelFish, nil
}

// helper function for catching fish
func (biz *FishBusiness) updateFishCount(profileID primitive.ObjectID, fishType string, additionalCount int32) (*repo.Fish, error) {
	// Get the current fish record by profileID and fishType
	repoFish, err := biz.FishRepo.GetFishByProfileID(profileID, fishType)
	if err != nil {
		return nil, fmt.Errorf("failed to find fish package: %v", err)
	}

	updatedFish := &repo.Fish{
		ProfileID: repoFish.ProfileID,
		Numbers:   repoFish.Numbers + additionalCount,
		Type:      repoFish.Type,
	}

	insertedFish, err := biz.FishRepo.UpdateFish(updatedFish, profileID)
	if err != nil {
		return nil, fmt.Errorf("failed to update fish: %v", err)
	}

	return insertedFish, nil
}

//	THOSE FUNCTIONS BELOW ARE FOR TRADING

func (biz *FishBusiness) UnlockMetricsWithNormalFish(ctx context.Context, characterID string) (bool, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coreRepo.Profile)
	if !ok {
		return false, errors.ErrorUnauthorized
	}

	repoFish, err := biz.FishRepo.GetFishByProfileID(profile.ID, "normal")
	if err != nil {
		return false, fmt.Errorf("failed to find fish package: %v", err)
	}

	if repoFish.Numbers < 3 { // Check if the number of fish is enough to trade
		return false, fmt.Errorf("There's no fish for you to trade")
	}

	characterOID, err := primitive.ObjectIDFromHex(characterID)
	if err != nil {
		return false, err
	}
	character, err := biz.CharactersRepo.GetCharacterByID(characterOID)
	if err != nil {
		return false, fmt.Errorf("failed to find character: %v", err)
	}

	// Update character metrics
	character.LimitedMetricNumber += 1
	if _, err := biz.CharactersRepo.UpdateCharacter(character); err != nil {
		return false, fmt.Errorf("failed to update metrics limited: %v", err)
	}

	// Decrease fish count
	repoFish.Numbers -= 3 //just a sample at this phase
	if _, err := biz.FishRepo.UpdateFish(repoFish, profile.ID); err != nil {
		return false, fmt.Errorf("failed to update fish: %v", err)
	}

	return true, nil
}

func (biz *FishBusiness) UnlockMetricsWithGoldFish(ctx context.Context, characterID string) (bool, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coreRepo.Profile)
	if !ok {
		return false, errors.ErrorUnauthorized
	}

	repoFish, err := biz.FishRepo.GetFishByProfileID(profile.ID, "gold")
	if err != nil {
		return false, fmt.Errorf("failed to find fish package: %v", err)
	}

	if repoFish.Numbers < 1 { // Check if the number of fish is enough to trade
		return false, fmt.Errorf("There's no fish for you to trade")
	}

	characterOID, err := primitive.ObjectIDFromHex(characterID)
	if err != nil {
		return false, err
	}
	character, err := biz.CharactersRepo.GetCharacterByID(characterOID)
	if err != nil {
		return false, fmt.Errorf("failed to find character: %v", err)
	}

	// Update character metrics
	character.LimitedMetricNumber += 1
	if _, err := biz.CharactersRepo.UpdateCharacter(character); err != nil {
		return false, fmt.Errorf("failed to update metrics limited: %v", err)
	}

	// Decrease fish count
	repoFish.Numbers -= 1 //just a sample at this phase
	if _, err := biz.FishRepo.UpdateFish(repoFish, profile.ID); err != nil {
		return false, fmt.Errorf("failed to update fish: %v", err)
	}

	return true, nil
}

func (biz *FishBusiness) BuySnapshotsWithNormalFish(ctx context.Context) (bool, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coreRepo.Profile)
	if !ok {
		return false, errors.ErrorUnauthorized
	}

	repoFish, err := biz.FishRepo.GetFishByProfileID(profile.ID, "normal")
	if err != nil {
		return false, fmt.Errorf("failed to find fish package: %v", err)
	}

	if repoFish.Numbers < 3 { // Check if the number of fish is enough to trade
		return false, fmt.Errorf("There's no fish for you to trade")
	}

	repoProfile, err := biz.ProfilesRepo.GetProfileByFirebaseUID(profile.FirebaseUID)
	if err != nil {
		return false, fmt.Errorf("failed get profile: %v", err)
	}

	// Update profile available snapshots
	repoProfile.AvailableSnapshots += 1
	if _, err := biz.ProfilesRepo.UpdateProfile(repoProfile); err != nil {
		return false, fmt.Errorf("failed to update available snapshots: %v", err)
	}

	// Decrease fish count
	repoFish.Numbers -= 3 //just a sample at this phase
	if _, err := biz.FishRepo.UpdateFish(repoFish, profile.ID); err != nil {
		return false, fmt.Errorf("failed to update fish: %v", err)
	}

	return true, nil
}

func (biz *FishBusiness) BuySnapshotsWithGoldFish(ctx context.Context) (bool, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coreRepo.Profile)
	if !ok {
		return false, errors.ErrorUnauthorized
	}

	repoFish, err := biz.FishRepo.GetFishByProfileID(profile.ID, "gold")
	if err != nil {
		return false, fmt.Errorf("failed to find fish package: %v", err)
	}

	if repoFish.Numbers < 1 { // Check if the number of fish is enough to trade
		return false, fmt.Errorf("There's no fish for you to trade")
	}

	repoProfile, err := biz.ProfilesRepo.GetProfileByFirebaseUID(profile.FirebaseUID)
	if err != nil {
		return false, fmt.Errorf("failed get profile: %v", err)
	}

	// Update profile available snapshots
	repoProfile.AvailableSnapshots += 1
	if _, err := biz.ProfilesRepo.UpdateProfile(repoProfile); err != nil {
		return false, fmt.Errorf("failed to update available snapshots: %v", err)
	}

	// Decrease fish count
	repoFish.Numbers -= 1 //just a sample at this phase
	if _, err := biz.FishRepo.UpdateFish(repoFish, profile.ID); err != nil {
		return false, fmt.Errorf("failed to update fish: %v", err)
	}

	return true, nil
}
