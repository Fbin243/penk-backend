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

func NewTemplate() *entity.Template {
	metric := entity.TemplateMetric{
		Name:  "example metric",
		Value: 100,
		Unit:  "example unit",
	}

	category := entity.TemplateCategory{
		Name:        "example category",
		Description: "example description",
		Style: entity.CategoryStyle{
			Color: "#000000",
			Icon:  "example",
		},
		Metrics: []entity.TemplateMetric{
			metric,
			metric,
			metric,
		},
	}

	return &entity.Template{
		BaseEntity: &base.BaseEntity{
			ID:        primitive.NewObjectID().Hex(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:        "example template",
		Description: "example description",
		TopicID:     primitive.NewObjectID().Hex(),
		Style: entity.TemplateStyle{
			Color: "#000000",
			Icon:  "example",
		},
		Categories: []entity.TemplateCategory{
			category,
			category,
			category,
		},
	}
}

func assertTemplate(t *testing.T, expected, actual *entity.Template) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Description, actual.Description)
	assert.Equal(t, expected.TopicID, actual.TopicID)
	assert.Equal(t, expected.Style, actual.Style)
	assert.Equal(t, expected.Categories, actual.Categories)
}

// func TestGetTemplates(t *testing.T) {
// 	templateMap := make(map[string]*entity.Template)
// 	for i := 0; i < 3; i++ {
// 		template := NewTemplate()
// 		templateMap[template.ID] = template
// 		createdTemplate, err := templateRepo.InsertOne(context.Background(), template)
// 		defer cleanUpTemplate(createdTemplate.ID)
// 		assert.Nil(t, err)
// 		assert.Equal(t, *createdTemplate, *template)
// 	}

// 	templates, err := templateRepo.FindAll(context.Background())
// 	assert.Nil(t, err)
// 	assert.Equal(t, len(templates), 3)
// 	for _, template := range templates {
// 		assertTemplate(t, templateMap[template.ID], &template)
// 	}
// }

func cleanUpTemplate(id string) {
	templateRepo.DeleteByID(context.Background(), id)
}
