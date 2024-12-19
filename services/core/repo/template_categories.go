package repo

import (
	"tenkhours/pkg/db"

	"go.mongodb.org/mongo-driver/mongo"
)

type TemplateCategoriesRepo struct {
	*db.BaseRepo[TemplateCategory]
}

func NewTemplateCategoriesRepo(mongodb *mongo.Database) *TemplateCategoriesRepo {
	return &TemplateCategoriesRepo{db.NewBaseRepo[TemplateCategory](mongodb.Collection(db.TemplateCategoriesCollection))}
}
