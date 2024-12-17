package business

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/analytics/graph/model"
	analyticsRepo "tenkhours/services/analytics/repo"
	coreRepo "tenkhours/services/core/repo"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CharactersBusiness struct {
	SnapshotsRepo       *analyticsRepo.SnapshotsRepo
	CharactersRepo      *coreRepo.CharactersRepo
	ProfilesRepo        *coreRepo.ProfilesRepo
	CapturedRecordsRepo *analyticsRepo.CapturedRecordsRepo
	RedisClient         *redis.Client
}

func NewCharactersBusiness(snapshotsRepo *analyticsRepo.SnapshotsRepo, charactersRepo *coreRepo.CharactersRepo, profilesRepo *coreRepo.ProfilesRepo, capturedRepo *analyticsRepo.CapturedRecordsRepo, redisClient *redis.Client) *CharactersBusiness {
	return &CharactersBusiness{
		SnapshotsRepo:       snapshotsRepo,
		CharactersRepo:      charactersRepo,
		ProfilesRepo:        profilesRepo,
		CapturedRecordsRepo: capturedRepo,
		RedisClient:         redisClient,
	}
}

func (biz *CharactersBusiness) GetSnapshots(ctx context.Context, characterID *primitive.ObjectID, filter *model.DateTimeFilter) ([]analyticsRepo.Snapshot, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coreRepo.Profile)
	if !ok {
		return nil, errors.ErrorUnauthorized
	}

	pineline := mongo.Pipeline{}
	matchStage := bson.D{}

	if characterID != nil {
		character, err := biz.CharactersRepo.FindByID(*characterID)
		if err != nil {
			return nil, fmt.Errorf("failed to get character by ID")
		}

		if character.ProfileID != profile.ID {
			return nil, errors.ErrorPermissionDenied
		}

		matchStage = append(matchStage, bson.E{Key: "metadata.character_id", Value: *characterID})
	} else {
		matchStage = append(matchStage, bson.E{Key: "metadata.profile_id", Value: profile.ID})
	}

	if filter != nil {
		if filter.Month != nil {
			month := utils.MonthToIntMap[(*filter.Month).String()]
			matchStage = append(matchStage, bson.E{
				Key: "$expr",
				Value: bson.D{{
					Key:   "$eq",
					Value: bson.A{bson.D{{Key: "$month", Value: "$timestamp"}}, month},
				}},
			})
		}

		if filter.Year != nil {
			matchStage = append(matchStage, bson.E{
				Key: "$expr",
				Value: bson.D{{
					Key:   "$eq",
					Value: bson.A{bson.D{{Key: "$year", Value: "$timestamp"}}, *filter.Year},
				}},
			})
		}
	}

	pineline = append(pineline, bson.D{{Key: "$match", Value: matchStage}})
	snapshots, err := biz.SnapshotsRepo.GetSnapshots(pineline)
	if err != nil {
		return nil, err
	}

	return snapshots, nil
}

func (biz *CharactersBusiness) CreateNewSnapshot(ctx context.Context, characterID primitive.ObjectID, description *string) (*analyticsRepo.Snapshot, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coreRepo.Profile)
	if !ok {
		return nil, errors.ErrorUnauthorized
	}

	character, err := biz.CharactersRepo.FindByID(characterID)
	if err != nil {
		return nil, err
	}

	if character.ProfileID != profile.ID {
		return nil, errors.ErrorPermissionDenied
	}

	if profile.AvailableSnapshots <= 0 {
		return nil, fmt.Errorf("no available snapshots")
	}

	// Set CharacterID and ProfileID to Nil to omit them
	character.ProfileID = primitive.NilObjectID
	character.ID = primitive.NilObjectID

	// Compare with the latest snapshot
	latestSnapshot, err := biz.SnapshotsRepo.GetLatestSnapshotByCharacterID(characterID)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	} else if reflect.DeepEqual(latestSnapshot.Character, *character) {
		return nil, fmt.Errorf("no changes detected")
	}

	snapshot := &analyticsRepo.Snapshot{
		ID:        primitive.NewObjectID(),
		Timestamp: utils.Now(),
		Metadata: analyticsRepo.Metadata{
			ProfileID:   profile.ID,
			CharacterID: characterID,
		},
		Character: *character,
	}

	if description != nil {
		snapshot.Description = *description
	}

	session, err := db.GetDBManager().Client.StartSession()
	if err != nil {
		return nil, err
	}

	defer session.EndSession(context.TODO())

	callback := func(ctx mongo.SessionContext) (interface{}, error) {
		_, err = biz.SnapshotsRepo.CreateSnapshot(snapshot)
		if err != nil {
			return nil, fmt.Errorf("failed to create snapshot: %v", err)
		}

		// Decrement available snapshots
		profile.AvailableSnapshots--
		_, err = biz.ProfilesRepo.UpdateProfile(&profile)
		if err != nil {
			return nil, fmt.Errorf("failed to update user profile: %v", err)
		}

		// Restore CharacterID and ProfileID
		snapshot.Character.ProfileID = profile.ID
		snapshot.Character.ID = characterID

		return *snapshot, nil
	}

	result, err := session.WithTransaction(context.TODO(), callback)
	if err != nil {
		return nil, err
	}

	newSnapshot, ok := result.(analyticsRepo.Snapshot)
	if !ok {
		return nil, fmt.Errorf("failed to create snapshot")
	}

	return &newSnapshot, nil
}

// Get analytic results from the captured records from the database
func (biz *CharactersBusiness) GetAnalyticResults(ctx context.Context, characterID *primitive.ObjectID, startTime *time.Time, endTime *time.Time, analyticSections []model.AnalyticSection) (map[string]interface{}, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coreRepo.Profile)
	if !ok {
		return nil, errors.ErrorUnauthorized
	}

	capturedRecordsFilter := CapturedRecordsFilter{
		FilterMethod:        FilterMethodServerSide,
		FilterType:          FilterTypeUser,
		ProfileID:           profile.ID,
		EndTime:             time.Now(),
		CapturedRecordsRepo: biz.CapturedRecordsRepo,
		RedisClient:         biz.RedisClient,
	}

	if characterID != nil {
		character, err := biz.CharactersRepo.FindByID(*characterID)
		if err != nil {
			return nil, fmt.Errorf("failed to get character by ID")
		}

		if character.ProfileID != profile.ID {
			return nil, errors.ErrorPermissionDenied
		}

		capturedRecordsFilter.FilterType = FilterTypeCharacter
		capturedRecordsFilter.CharacterID = *characterID
	}

	if startTime != nil {
		capturedRecordsFilter.StartTime = utils.ResetTimeToBeginningOfDay(*startTime)
	}

	if endTime != nil {
		capturedRecordsFilter.EndTime = utils.ResetTimeToBeginningOfDay(*endTime)
	}

	capturedRecords, err := capturedRecordsFilter.Filter()
	if err != nil {
		return nil, err
	}

	analyticsProcessor := &AnalyticsProcessor{
		AnalyticSections: analyticSections,
		CapturedRecords:  capturedRecords,
		AnalyticResults:  make(map[string]interface{}),
		FilterType:       capturedRecordsFilter.FilterType,
		StartTime:        capturedRecordsFilter.StartTime,
		EndTime:          capturedRecordsFilter.EndTime,
	}

	return analyticsProcessor.ProcessCapturedRecords(), nil
}
