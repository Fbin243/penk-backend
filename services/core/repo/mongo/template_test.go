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
	property := entity.TemplateProperty{
		Name:  "example property",
		Type:  entity.MetricPropertyTypeNumber,
		Value: "100",
		Unit:  "example unit",
	}

	metric := entity.TemplateMetric{
		Name:        "example metric",
		Description: "example description",
		Style: entity.MetricStyle{
			Color: "#000000",
			Icon:  "example",
		},
		Properties: []entity.TemplateProperty{
			property,
			property,
			property,
		},
	}

	return &entity.Template{
		BaseEntity: &base.BaseEntity{
			ID:        primitive.NewObjectID().Hex(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:        "example",
		Description: "example description",
		CategoryID:  primitive.NewObjectID().Hex(),
		Style: entity.TemplateStyle{
			Color: "#000000",
			Icon:  "example",
		},
		Metrics: []entity.TemplateMetric{
			metric,
			metric,
			metric,
		},
	}
}

func assertTemplate(t *testing.T, expected, actual *entity.Template) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Description, actual.Description)
	assert.Equal(t, expected.CategoryID, actual.CategoryID)
	assert.Equal(t, expected.Style, actual.Style)
	assert.Equal(t, expected.Metrics, actual.Metrics)
}

func TestGetTemplates(t *testing.T) {
	templateMap := make(map[string]*entity.Template)
	for i := 0; i < 3; i++ {
		template := NewTemplate()
		templateMap[template.ID] = template
		createdTemplate, err := templateRepo.InsertOne(context.Background(), template)
		defer cleanUpTemplate(createdTemplate.ID)
		assert.Nil(t, err)
		assert.Equal(t, *createdTemplate, *template)
	}

	templates, err := templateRepo.FindAll(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, len(templates), 3)
	for _, template := range templates {
		assertTemplate(t, templateMap[template.ID], &template)
	}
}

func cleanUpTemplate(id string) {
	templateRepo.DeleteByID(context.Background(), id)
}
