package mongorepo

import (
	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/mongo"
)

type TemplateCategoryRepo struct {
	*mongodb.BaseRepo[entity.TemplateTopic, TemplateTopic]
}

func NewTemplateTopicRepo(db *mongo.Database) *TemplateCategoryRepo {
	return &TemplateCategoryRepo{
		mongodb.NewBaseRepo(
			db.Collection(mongodb.TemplateTopicsCollection),
			&mongodb.Mapper[entity.TemplateTopic, TemplateTopic]{},
		),
	}
}
