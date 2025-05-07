package mongorepo

import (
	"context"
	"log"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/types"
	"tenkhours/services/core/entity"
	mongomodel "tenkhours/services/core/repo/mongo/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReminderRepo struct {
	*mongodb.BaseRepo[entity.Reminder, mongomodel.Reminder]
}

func NewReminderRepo(db *mongo.Database) *ReminderRepo {
	collection := db.Collection(mongodb.RemindersCollection)

	// Create indexes
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "character_id", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "remind_time", Value: 1},
			},
		},
	}

	_, err := collection.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		log.Println("Error creating indexes for reminders collection:", err)
	}

	return &ReminderRepo{mongodb.NewBaseRepo[entity.Reminder, mongomodel.Reminder](
		collection,
		true,
	)}
}

func (r *ReminderRepo) FindByCharacterID(ctx context.Context, characterID string) ([]entity.Reminder, error) {
	return r.BaseRepo.FindMany(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
}

func (r *ReminderRepo) CountByCharacterID(ctx context.Context, characterID string) (int, error) {
	return r.BaseRepo.Count(ctx, bson.M{"character_id": mongodb.ToObjectID(characterID)})
}

func (r *ReminderRepo) Find(ctx context.Context, p entity.ReminderPipeline) ([]entity.Reminder, error) {
	return r.AggregateQuery(ctx,
		r.addPaginationStage(
			r.addSortStage(
				r.addMatchStage(
					[]bson.M{},
					p.Filter,
				),
				p.OrderBy,
			),
			p.Pagination,
		),
	)
}

func (r *ReminderRepo) CountByFilter(ctx context.Context, filter *entity.ReminderFilter) (int, error) {
	return r.AggregateCount(ctx,
		r.addCountStage(
			r.addMatchStage(
				[]bson.M{},
				filter,
			),
		),
	)
}

func (r *ReminderRepo) addMatchStage(p []bson.M, filter *entity.ReminderFilter) []bson.M {
	if filter == nil {
		return p
	}

	matchStage := bson.M{}
	if filter.CharacterID != nil {
		matchStage["character_id"] = mongodb.ToObjectID(*filter.CharacterID)
	}

	return append(p, bson.M{"$match": matchStage})
}

func (r *ReminderRepo) addSortStage(p []bson.M, orderBy *entity.ReminderOrderBy) []bson.M {
	return p
}

func (r *ReminderRepo) addPaginationStage(p []bson.M, pagination *types.Pagination) []bson.M {
	if pagination == nil {
		return p
	}

	return append(p, mongodb.ToPaginationPineline(pagination)...)
}

func (r *ReminderRepo) addCountStage(p []bson.M) []bson.M {
	return append(p, bson.M{"$count": "count"})
}
