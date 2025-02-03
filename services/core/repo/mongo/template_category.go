package mongorepo

import (
	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/mongo"
)

type TemplateCategoryRepo struct {
	*mongodb.BaseRepo[entity.TemplateCategory, TemplateCategory]
}

func NewTemplateCategoryRepo(db *mongo.Database) *TemplateCategoryRepo {
	return &TemplateCategoryRepo{
		mongodb.NewBaseRepo(
			db.Collection(mongodb.TemplateCategoriesCollection),
			&mongodb.Mapper[entity.TemplateCategory, TemplateCategory]{},
		),
	}
}
