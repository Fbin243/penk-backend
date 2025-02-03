package mongorepo_test

import (
	"context"
	"testing"
	"time"

	"tenkhours/pkg/db/base"
	"tenkhours/services/core/entity"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewTemplateCategory() *entity.TemplateCategory {
	return &entity.TemplateCategory{
		BaseEntity: &base.BaseEntity{
			ID:        primitive.NewObjectID().Hex(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:        "example",
		Description: "example description",
	}
}

func assertTemplateCategory(t *testing.T, expected, actual *entity.TemplateCategory) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Description, actual.Description)
}

func TestGetTemplateCategories(t *testing.T) {
	categoryMap := make(map[string]*entity.TemplateCategory)
	for i := 0; i < 3; i++ {
		category := NewTemplateCategory()
		categoryMap[category.ID] = category
		createdCategory, err := templateCategoryRepo.InsertOne(context.Background(), category)
		defer cleanUpCategory(createdCategory.ID)
		assert.Nil(t, err)
		assert.Equal(t, createdCategory, category)
	}

	categories, err := templateCategoryRepo.FindAll(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, len(categories), 3)
	for _, c := range categories {
		assertTemplateCategory(t, categoryMap[c.ID], &c)
	}
}

func cleanUpCategory(id string) {
	templateCategoryRepo.DeleteByID(context.Background(), id)
}
