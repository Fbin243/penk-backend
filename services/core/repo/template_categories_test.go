package repo_test

import (
	"tenkhours/pkg/db"
	"tenkhours/services/core/repo"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetTemplateCategories(t *testing.T) {
	category := &repo.TemplateCategory{
		BaseModel:   &db.BaseModel{},
		Name:        "example",
		Description: "example description",
	}

	for i := 0; i < 3; i++ {
		createdCategory, err := templateCategoriesRepo.InsertOne(category)
		defer cleanUpCategory(createdCategory.ID)
		assert.Nil(t, err)
		assert.Equal(t, *createdCategory, *category)
	}

	categories, err := templateCategoriesRepo.FindAll()
	assert.Nil(t, err)
	assert.Equal(t, len(categories), 3)
}

func cleanUpCategory(id primitive.ObjectID) {
	templatesRepo.DeleteByID(id)
}
