package repo_test

import (
	"tenkhours/pkg/db"
	"tenkhours/services/core/repo"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetTemplates(t *testing.T) {
	property := repo.TemplateProperty{
		Name:  "example property",
		Type:  repo.MetricPropertyTypeNumber,
		Value: "100",
		Unit:  "example unit",
	}

	metric := repo.TemplateMetric{
		Name:        "example metric",
		Description: "example description",
		Style: repo.MetricStyle{
			Color: "#000000",
			Icon:  "example",
		},
		Properties: []repo.TemplateProperty{
			property,
			property,
			property,
		},
	}

	template := &repo.Template{
		BaseModel:   &db.BaseModel{},
		Name:        "example",
		Description: "example description",
		CategoryID:  primitive.NewObjectID(),
		Style: repo.TemplateStyle{
			Color: "#000000",
			Icon:  "example",
		},
		Metrics: []repo.TemplateMetric{
			metric,
			metric,
			metric,
		},
	}

	for i := 0; i < 3; i++ {
		createdTemplate, err := templatesRepo.InsertOne(template)
		defer cleanUpTemplate(createdTemplate.ID)
		assert.Nil(t, err)
		assert.Equal(t, *createdTemplate, *template)
	}

	templates, err := templatesRepo.FindAll()
	assert.Nil(t, err)
	assert.Equal(t, len(templates), 3)
}

func cleanUpTemplate(id primitive.ObjectID) {
	templatesRepo.DeleteByID(id)
}
