package analytics

import (
	"context"
	"fmt"
	"reflect"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db"
	"tenkhours/pkg/db/analyticsdb"
	"tenkhours/pkg/db/coredb"
	"tenkhours/pkg/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CharactersHandler struct {
	SnapshotsRepo  *analyticsdb.SnapshotsRepo
	CharactersRepo *coredb.CharactersRepo
	ProfilesRepo   *coredb.ProfilesRepo
}

func NewCharactersHandler(snapshotsRepo *analyticsdb.SnapshotsRepo, charactersRepo *coredb.CharactersRepo, profilesRepo *coredb.ProfilesRepo) *CharactersHandler {
	return &CharactersHandler{
		SnapshotsRepo:  snapshotsRepo,
		CharactersRepo: charactersRepo,
		ProfilesRepo:   profilesRepo,
	}
}

func (r *CharactersHandler) GetSnapshotsByProfileID(ctx context.Context) ([]analyticsdb.Snapshot, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coredb.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	snapshots, err := r.SnapshotsRepo.GetSnapshotsByProfileID(profile.ID)
	if err != nil {
		return nil, err
	}

	return snapshots, nil
}

func (r *CharactersHandler) GetSnapshotsByCharacterID(ctx context.Context, characterID primitive.ObjectID) ([]analyticsdb.Snapshot, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coredb.Profile)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	character, err := r.CharactersRepo.GetCharacterByID(characterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get character by ID")
	}

	if character.ProfileID != profile.ID {
		return nil, auth.ErrorPermissionDenied
	}

	snapshots, err := r.SnapshotsRepo.GetSnapshotsByCharacterID(characterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshots by character ID")
	}

	return snapshots, nil
}

func (r *CharactersHandler) CreateNewSnapshot(ctx context.Context, characterID primitive.ObjectID) (*analyticsdb.Snapshot, error) {
	profile, ok := ctx.Value(auth.ProfileKey).(coredb.Profile)
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
	fmt.Printf("latestSnapshot: %v\n", latestSnapshot)
	fmt.Printf("character: %v\n", *character)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	} else if reflect.DeepEqual(latestSnapshot.Character, *character) {
		return nil, fmt.Errorf("no changes detected")
	}

	snapshot := &analyticsdb.Snapshot{
		ID:        primitive.NewObjectID(),
		Timestamp: utils.Now(),
		Metadata: analyticsdb.Metadata{
			ProfileID:   profile.ID,
			CharacterID: characterID,
		},
		Character: *character,
		Asset:     nil,
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

	newSnapshot, ok := result.(analyticsdb.Snapshot)
	if !ok {
		return nil, fmt.Errorf("failed to create snapshot")
	}

	return &newSnapshot, nil
}
