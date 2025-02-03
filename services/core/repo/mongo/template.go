package mongorepo

import (
	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/mongo"
)

type TemplateRepo struct {
	*mongodb.BaseRepo[entity.Template, Template]
}

func NewTemplateRepo(db *mongo.Database) *TemplateRepo {
	return &TemplateRepo{mongodb.NewBaseRepo(
		db.Collection(mongodb.TemplatesCollection),
		&mongodb.Mapper[entity.Template, Template]{},
	)}
}
