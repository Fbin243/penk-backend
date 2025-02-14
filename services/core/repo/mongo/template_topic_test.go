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

func NewTemplateTopic() *entity.TemplateTopic {
	return &entity.TemplateTopic{
		BaseEntity: &base.BaseEntity{
			ID:        primitive.NewObjectID().Hex(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:        "example topic",
		Description: "example description",
	}
}

func assertTemplateTopic(t *testing.T, expected, actual *entity.TemplateTopic) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Description, actual.Description)
}

func TestGetTemplateCategories(t *testing.T) {
	categoryMap := make(map[string]*entity.TemplateTopic)
	for i := 0; i < 3; i++ {
		category := NewTemplateTopic()
		categoryMap[category.ID] = category
		createdTopic, err := templateTopicRepo.InsertOne(context.Background(), category)
		defer cleanUpTopic(createdTopic.ID)
		assert.Nil(t, err)
		assert.Equal(t, createdTopic, category)
	}

	topics, err := templateTopicRepo.FindAll(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, len(topics), 3)
	for _, c := range topics {
		assertTemplateTopic(t, categoryMap[c.ID], &c)
	}
}

func cleanUpTopic(id string) {
	templateTopicRepo.DeleteByID(context.Background(), id)
}
