package mongorepo_test

import (
	"context"
	"testing"

	"tenkhours/pkg/db/base"
	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/entity"

	"github.com/stretchr/testify/assert"
)

func NewCategory() *entity.Category {
	return &entity.Category{
		BaseEntity: &base.BaseEntity{
			ID:        mongodb.GenObjectID(),
			CreatedAt: utils.Now(),
			UpdatedAt: utils.Now(),
		},
		CharacterID: mongodb.GenObjectID(),
		Name:        "Category name",
		Description: "Category description",
		Style: entity.CategoryStyle{
			Color: "Category color",
			Icon:  "Category icon",
		},
	}
}

func assertCategory(t *testing.T, expected, actual entity.Category) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.CreatedAt, actual.CreatedAt)
	assert.Equal(t, expected.UpdatedAt, actual.UpdatedAt)
	assert.Equal(t, expected.CharacterID, actual.CharacterID)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Description, actual.Description)
	assert.Equal(t, expected.Style.Color, actual.Style.Color)
	assert.Equal(t, expected.Style.Icon, actual.Style.Icon)
}

func TestCategoryRepo(t *testing.T) {
	characterID := mongodb.GenObjectID()

	// Create category
	for range 3 {
		category := NewCategory()
		category.CharacterID = characterID
		createdCategory, err := categoryRepo.InsertOne(context.Background(), category)
		assert.Nil(t, err)
		assertCategory(t, *category, *createdCategory)
	}

	// Find all categories of a character
	categories, err := categoryRepo.Find(context.Background(), entity.CategoryPipeline{
		Filter: &entity.CategoryFilter{
			CharacterID: &characterID,
		},
	})
	assert.Nil(t, err)
	assert.Len(t, categories, 3)

	categoryID := categories[0].ID
	// Check exist
	err = categoryRepo.Exist(context.Background(), characterID, categoryID)
	assert.Nil(t, err)

	// Count all categories of a character
	count, err := categoryRepo.CountByCharacterID(context.Background(), characterID)
	assert.Nil(t, err)
	assert.Equal(t, int64(3), count)

	// Delete all categories of a character
	err = categoryRepo.DeleteByCharacterID(context.Background(), characterID)
	assert.Nil(t, err)
	count, err = categoryRepo.CountByCharacterID(context.Background(), characterID)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), count)

	// Check exist
	err = categoryRepo.Exist(context.Background(), characterID, categoryID)
	assert.NotNil(t, err)
	assert.Equal(t, errors.ErrMongoNotFound, err)

	// Create two categories for two characters
	characterIDs := []string{mongodb.GenObjectID(), mongodb.GenObjectID()}
	for i := range 2 {
		category := NewCategory()
		category.CharacterID = characterIDs[i]
		_, err := categoryRepo.InsertOne(context.Background(), category)
		assert.Nil(t, err)
	}

	// Delete categories of two characters
	err = categoryRepo.DeleteByCharacterIDs(context.Background(), characterIDs)
	assert.Nil(t, err)
	// count, err = categoryRepo.CountAll(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, int64(0), count)
}
