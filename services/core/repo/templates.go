package repo

import (
	"tenkhours/pkg/db"

	"go.mongodb.org/mongo-driver/mongo"
)

type TemplatesRepo struct {
	*db.BaseRepo[Template]
}

func NewTemplatesRepo(mongodb *mongo.Database) *TemplatesRepo {
	return &TemplatesRepo{db.NewBaseRepo[Template](mongodb.Collection(db.TemplatesCollection))}
}