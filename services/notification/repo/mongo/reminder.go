package mongorepo

import (
	mongodb "tenkhours/pkg/db/mongo"
	core_entity "tenkhours/services/core/entity"
	core_mongomodel "tenkhours/services/core/repo/mongo/model"

	"go.mongodb.org/mongo-driver/mongo"
)

type ReminderRepo struct {
	*mongodb.BaseRepo[core_entity.Reminder, core_mongomodel.Reminder]
}

func NewReminderRepo(db *mongo.Database) *ReminderRepo {
	return &ReminderRepo{mongodb.NewBaseRepo[core_entity.Reminder, core_mongomodel.Reminder](
		db.Collection(mongodb.RemindersCollection),
		true,
	)}
}
