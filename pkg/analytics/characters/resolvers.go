package characters

import (
	"context"
	"fmt"
	"reflect"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db"
	"tenkhours/pkg/db/analyticsdb"
	"tenkhours/pkg/db/coredb"
	"tenkhours/pkg/utils"

	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CharactersResolver struct {
	SnapshotsRepo  *analyticsdb.SnapshotsRepo
	CharactersRepo *coredb.CharactersRepo
	UsersRepo      *coredb.UsersRepo
}

func NewCharactersResolver(snapshotsRepo *analyticsdb.SnapshotsRepo, charactersRepo *coredb.CharactersRepo, usersRepo *coredb.UsersRepo) *CharactersResolver {
	return &CharactersResolver{
		SnapshotsRepo:  snapshotsRepo,
		CharactersRepo: charactersRepo,
		UsersRepo:      usersRepo,
	}
}

func (r *CharactersResolver) GetSnapshotsByUserID(params graphql.ResolveParams) (interface{}, error) {
	user, ok := params.Context.Value(auth.UserKey).(coredb.User)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	snapshots, err := r.SnapshotsRepo.GetSnapshotsByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	return snapshots, nil
}

func (r *CharactersResolver) GetSnapshotsByCharacterID(params graphql.ResolveParams) (interface{}, error) {
	user, ok := params.Context.Value(auth.UserKey).(coredb.User)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	characterID := params.Args["characterID"].(string)
	characterOID, err := primitive.ObjectIDFromHex(characterID)
	if err != nil {
		return nil, err
	}

	character, err := r.CharactersRepo.GetCharacterByID(characterOID)
	if err != nil {
		return nil, fmt.Errorf("failed to get character by ID")
	}

	if character.UserID != user.ID {
		return nil, auth.ErrorPermissionDenied
	}

	snapshots, err := r.SnapshotsRepo.GetSnapshotsByCharacterID(characterOID)
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshots by character ID")
	}

	return snapshots, nil
}

func (r *CharactersResolver) CreateNewSnapshot(params graphql.ResolveParams) (interface{}, error) {
	user, ok := params.Context.Value(auth.UserKey).(coredb.User)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	characterID := params.Args["characterID"].(string)
	characterOID, err := primitive.ObjectIDFromHex(characterID)
	if err != nil {
		return nil, err
	}

	character, err := r.CharactersRepo.GetCharacterByID(characterOID)
	if err != nil {
		return nil, err
	}

	if character.UserID != user.ID {
		return nil, auth.ErrorPermissionDenied
	}

	if user.AvailableSnapshots <= 0 {
		return nil, fmt.Errorf("no available snapshots")
	}

	// Set CharacterID and UserID to Nil to omit them
	character.UserID = primitive.NilObjectID
	character.ID = primitive.NilObjectID

	// Compare with the latest snapshot
	latestSnapshot, err := r.SnapshotsRepo.GetLatestSnapshotByCharacterID(characterOID)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, err
		}
	} else if reflect.DeepEqual(latestSnapshot.Character, *character) {
		return nil, fmt.Errorf("no changes detected")
	}

	snapshot := &analyticsdb.Snapshot{
		ID:        primitive.NewObjectID(),
		Timestamp: utils.Now(),
		Metadata: analyticsdb.Metadata{
			UserID:      user.ID,
			CharacterID: characterOID,
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
		user.AvailableSnapshots--
		_, err = r.UsersRepo.UpdateUser(&user)
		if err != nil {
			return nil, fmt.Errorf("failed to update user: %v", err)
		}

		return *snapshot, nil
	}

	return session.WithTransaction(context.TODO(), callback)
}
