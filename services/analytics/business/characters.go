package business

import (
	"context"
	"fmt"
	"reflect"
	"tenkhours/pkg/auth"
	"tenkhours/pkg/db"
	"tenkhours/pkg/utils"
	"tenkhours/services/analytics/graph/model"

	analyticsRepo "tenkhours/services/analytics/repo"
	coreRepo "tenkhours/services/core/repo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CharactersHandler struct {
	SnapshotsRepo       *analyticsRepo.SnapshotsRepo
	CharactersRepo      *coreRepo.CharactersRepo
	ProfilesRepo        *coreRepo.ProfilesRepo
	CapturedRecordsRepo *analyticsRepo.CapturedRecordsRepo
}

func NewCharactersHandler(snapshotsRepo *analyticsRepo.SnapshotsRepo, charactersRepo *coreRepo.CharactersRepo, profilesRepo *coreRepo.ProfilesRepo, capturedRepo *analyticsRepo.CapturedRecordsRepo) *CharactersHandler {
	return &CharactersHandler{
		SnapshotsRepo:       snapshotsRepo,
		CharactersRepo:      charactersRepo,
		ProfilesRepo:        profilesRepo,
		CapturedRecordsRepo: capturedRepo,
	}
}

func (r *CharactersHandler) GetSnapshots(ctx context.Context, characterID *primitive.ObjectID, filter *model.Filter) ([]analyticsRepo.Snapshot, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coreRepo.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	pineline := mongo.Pipeline{}
	matchStage := bson.D{}

	if characterID != nil {
		character, err := r.CharactersRepo.GetCharacterByID(*characterID)
		if err != nil {
			return nil, fmt.Errorf("failed to get character by ID")
		}

		if character.ProfileID != profile.ID {
			return nil, auth.ErrorPermissionDenied
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
	snapshots, err := r.SnapshotsRepo.GetSnapshots(pineline)
	if err != nil {
		return nil, err
	}

	return snapshots, nil
}

func (r *CharactersHandler) CreateNewSnapshot(ctx context.Context, characterID primitive.ObjectID, description *string) (*analyticsRepo.Snapshot, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coreRepo.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	character, err := r.CharactersRepo.GetCharacterByID(characterID)
	if err != nil {
		return nil, err
	}

	if character.ProfileID != profile.ID {
		return nil, auth.ErrorPermissionDenied
	}

	if profile.AvailableSnapshots <= 0 {
		return nil, fmt.Errorf("no available snapshots")
	}

	// Set CharacterID and ProfileID to Nil to omit them
	character.ProfileID = primitive.NilObjectID
	character.ID = primitive.NilObjectID

	// Compare with the latest snapshot
	latestSnapshot, err := r.SnapshotsRepo.GetLatestSnapshotByCharacterID(characterID)
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
		_, err = r.SnapshotsRepo.CreateSnapshot(snapshot)
		if err != nil {
			return nil, fmt.Errorf("failed to create snapshot: %v", err)
		}

		// Decrement available snapshots
		profile.AvailableSnapshots--
		_, err = r.ProfilesRepo.UpdateProfile(&profile)
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

func (r *CharactersHandler) CreateCapturedRecord(ctx context.Context, characterID primitive.ObjectID) (*model.CapturedRecord, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coreRepo.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	character, err := r.CharactersRepo.GetCharacterByID(characterID)
	if err != nil {
		return nil, err
	}

	if character.ProfileID != profile.ID {
		return nil, auth.ErrorPermissionDenied
	}

	// Extract time from the custom metrics
	metricRecord := model.CapturedRecordCustomMetric{}
	metricRecords := make([]model.CapturedRecordCustomMetric, 0)
	for _, metric := range character.CustomMetrics {
		metricRecord.ID = metric.ID
		metricRecord.Time = metric.Time
		metricRecords = append(metricRecords, metricRecord)
	}

	// Make a new captured record
	capturedRecord := &model.CapturedRecord{
		ID:               primitive.NewObjectID(),
		Timestamp:        utils.Now(),
		TotalFocusedTime: character.TotalFocusedTime,
		CustomMetrics:    metricRecords,
		Metadata: model.CapturedRecordMetadata{
			ProfileID:   profile.ID,
			CharacterID: characterID,
		},
	}

	// Save the captured record
	err = r.CapturedRecordsRepo.CreateCapturedRecord(capturedRecord)
	if err != nil {
		return nil, err
	}

	return capturedRecord, nil
}

func (r *CharactersHandler) GetAnalyticResults(ctx context.Context, characterID *primitive.ObjectID, filter *model.Filter) (map[string]interface{}, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coreRepo.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	pipeline := mongo.Pipeline{}
	matchStage := bson.D{}

	if characterID != nil {
		character, err := r.CharactersRepo.GetCharacterByID(*characterID)
		if err != nil {
			return nil, fmt.Errorf("failed to get character by ID")
		}

		if character.ProfileID != profile.ID {
			return nil, auth.ErrorPermissionDenied
		}

		matchStage = append(matchStage, bson.E{Key: "metadata.character_id", Value: *characterID})
	} else {
		matchStage = append(matchStage, bson.E{Key: "metadata.profile_id", Value: profile.ID})
	}

	pipeline = append(pipeline, bson.D{{Key: "$match", Value: matchStage}})

	capturedRecords, err := r.CapturedRecordsRepo.GetCapturedRecords(pipeline)
	if err != nil {
		return nil, err
	}

	return processCapturedRecords(capturedRecords), nil
}
