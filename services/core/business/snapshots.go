package business

import (
	"context"
	"reflect"
	"tenkhours/pkg/auth"
	"tenkhours/pkg/db"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/graph/model"
	"tenkhours/services/core/repo"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SnapshotsBusiness struct {
	ProfilesRepo   *repo.ProfilesRepo
	CharactersRepo *repo.CharactersRepo
	SnapshotsRepo  *repo.SnapshotsRepo
}

func NewSnapshotsBusiness(profilesRepo *repo.ProfilesRepo, charactersRepo *repo.CharactersRepo, snapshotsRepo *repo.SnapshotsRepo) *SnapshotsBusiness {
	return &SnapshotsBusiness{
		ProfilesRepo:   profilesRepo,
		CharactersRepo: charactersRepo,
		SnapshotsRepo:  snapshotsRepo,
	}
}

func (biz *SnapshotsBusiness) GetSnapshots(ctx context.Context, characterID *primitive.ObjectID, filter *model.DateTimeFilter) ([]repo.Snapshot, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(db.AuthSession)
	if !ok {
		return nil, errors.Unauthorized()
	}

	pineline := mongo.Pipeline{}
	matchStage := bson.D{}

	if characterID != nil {
		character, err := biz.CharactersRepo.FindByID(*characterID)
		if err != nil {
			return nil, err
		}

		if character.ProfileID != authSession.ProfileID {
			return nil, errors.PermissionDenied()
		}

		matchStage = append(matchStage, bson.E{Key: "metadata.character_id", Value: *characterID})
	} else {
		matchStage = append(matchStage, bson.E{Key: "metadata.profile_id", Value: authSession.ProfileID})
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

func (biz *SnapshotsBusiness) CreateNewSnapshot(ctx context.Context, characterID primitive.ObjectID, description *string) (*repo.Snapshot, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(db.AuthSession)
	if !ok {
		return nil, errors.Unauthorized()
	}

	character, err := biz.CharactersRepo.FindByID(characterID)
	if err != nil {
		return nil, err
	}

	if character.ProfileID != authSession.ProfileID {
		return nil, errors.PermissionDenied()
	}

	profile, err := biz.ProfilesRepo.FindByID(authSession.ProfileID)
	if profile.AvailableSnapshots <= 0 {
		return nil, errors.NewError(errors.ErrCodeLimitSnapshot, "no available snapshots")
	}

	// Create the snapshot character without unnecessary fields
	snapshotCharacter := repo.SnapshotCharacter{
		Name:             character.Name,
		Gender:           character.Gender,
		Tags:             character.Tags,
		TotalFocusedTime: character.TotalFocusedTime,
		CustomMetrics: lo.Map(character.CustomMetrics, func(metric repo.CustomMetric, _ int) repo.SnapshotMetric {
			return repo.SnapshotMetric{
				ID:          metric.ID,
				Name:        metric.Name,
				Description: metric.Description,
				Time:        metric.Time,
				Style:       metric.Style,
			}
		}),
	}

	// Compare with the latest snapshot
	latestSnapshot, err := biz.SnapshotsRepo.GetLatestSnapshotByCharacterID(characterID)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	} else if reflect.DeepEqual(latestSnapshot.Character, snapshotCharacter) {
		return nil, errors.NewError(errors.ErrCodeDuplicateSnapshot, "duplicate snapshot")
	}

	snapshot := &repo.Snapshot{
		ID:        primitive.NewObjectID(),
		Timestamp: utils.Now(),
		Metadata: repo.Metadata{
			ProfileID:   authSession.ProfileID,
			CharacterID: characterID,
		},
		Character: snapshotCharacter,
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
			return nil, err
		}

		profile.AvailableSnapshots--
		_, err = biz.ProfilesRepo.UpdateProfile(profile)
		if err != nil {
			return nil, err
		}

		return *snapshot, nil
	}

	result, err := session.WithTransaction(context.TODO(), callback)
	if err != nil {
		return nil, err
	}

	newSnapshot, ok := result.(repo.Snapshot)
	if !ok {
		return nil, errors.NewError(errors.ErrCodeInternalServer, "failed to create snapshot")
	}

	return &newSnapshot, nil
}
