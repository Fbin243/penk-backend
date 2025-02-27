package rpc_test

import (
	"testing"

	"tenkhours/pkg/db/base"
	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/utils"
	"tenkhours/proto/pb/core"
	"tenkhours/services/core/entity"

	"github.com/jinzhu/copier"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

var category = entity.Category{
	ID:          mongodb.GenObjectID(),
	Name:        "Example name",
	Description: "Example desc",
	Style: entity.CategoryStyle{
		Color: "red",
		Icon:  "icon.png",
	},
}

var metric = entity.Metric{
	ID:         mongodb.GenObjectID(),
	CategoryID: &category.ID,
	Name:       "Example name",
	Value:      1.0,
	Unit:       "unit",
}

var entityCharacter = &entity.Character{
	BaseEntity: &base.BaseEntity{
		ID:        mongodb.GenObjectID(),
		CreatedAt: utils.Now(),
		UpdatedAt: utils.Now(),
	},
	ProfileID: mongodb.GenObjectID(),
	Name:      "Example",
	Tags:      []string{"#tag1", "#tag2"},
	Gender:    false,
	Categories: []entity.Category{
		category, category, category,
	},
	Metrics: []entity.Metric{
		metric, metric, metric,
	},
	Vision: entity.Vision{
		Name:        "Example name",
		Description: "Example desc",
	},
}

func TestMapCharacter(t *testing.T) {
	rpcCharacter := &core.Character{}
	copier.Copy(rpcCharacter, entityCharacter)
	rpcCharacter.CreatedAt = entityCharacter.CreatedAt.Unix()
	rpcCharacter.UpdatedAt = entityCharacter.UpdatedAt.Unix()
	rpcCharacterCategories := lo.Map(rpcCharacter.Categories, func(c *core.Category, _ int) entity.Category {
		return entity.Category{
			ID:          c.Id,
			Name:        c.Name,
			Description: c.Description,
			Style: entity.CategoryStyle{
				Color: c.Style.Color,
				Icon:  c.Style.Icon,
			},
		}
	})
	rpcCharacterMetrics := lo.Map(rpcCharacter.Metrics, func(m *core.Metric, _ int) entity.Metric {
		return entity.Metric{
			ID:         m.Id,
			CategoryID: m.CategoryId,
			Name:       m.Name,
			Value:      float64(m.Value),
			Unit:       m.Unit,
		}
	})

	assert.Equal(t, entityCharacter.ID, rpcCharacter.Id)
	assert.Equal(t, entityCharacter.CreatedAt.Unix(), rpcCharacter.CreatedAt)
	assert.Equal(t, entityCharacter.UpdatedAt.Unix(), rpcCharacter.UpdatedAt)
	assert.Equal(t, entityCharacter.ProfileID, rpcCharacter.ProfileId)
	assert.Equal(t, entityCharacter.Name, rpcCharacter.Name)
	assert.Equal(t, entityCharacter.Tags, rpcCharacter.Tags)
	assert.Equal(t, entityCharacter.Gender, rpcCharacter.Gender)
	assert.Equal(t, entityCharacter.Categories, rpcCharacterCategories)
	assert.Equal(t, entityCharacter.Metrics, rpcCharacterMetrics)
	assert.Equal(t, entityCharacter.Vision.Name, rpcCharacter.Vision.Name)
	assert.Equal(t, entityCharacter.Vision.Description, *rpcCharacter.Vision.Description)
}

var rpcCategoryInput = &core.CategoryInput{
	Id:          lo.ToPtr(mongodb.GenObjectID()),
	Name:        "Example name",
	Description: lo.ToPtr("Example desc"),
	Style: &core.CategoryStyleInput{
		Color: "red",
		Icon:  "icon.png",
	},
}

var rpcMetricInput = &core.MetricInput{
	Id:         lo.ToPtr(mongodb.GenObjectID()),
	CategoryId: lo.ToPtr(mongodb.GenObjectID()),
	Name:       "Example name",
	Value:      1.0,
	Unit:       "unit",
}

var rpcCharacterInput = &core.CharacterInput{
	Id:     lo.ToPtr(mongodb.GenObjectID()),
	Name:   "Example",
	Gender: false,
	Tags:   []string{"#tag1", "#tag2"},
	Categories: []*core.CategoryInput{
		rpcCategoryInput, rpcCategoryInput, rpcCategoryInput,
	},
	Metrics: []*core.MetricInput{
		rpcMetricInput, rpcMetricInput, rpcMetricInput,
	},
	Vision: &core.VisionInput{
		Name:        "Example name",
		Description: lo.ToPtr("Example desc"),
	},
}

func TestMapCharacterInput(t *testing.T) {
	entityCharacterInput := &entity.CharacterInput{}
	copier.Copy(entityCharacterInput, rpcCharacterInput)

	assert.Equal(t, entityCharacterInput.ID, rpcCharacterInput.Id)
	assert.Equal(t, entityCharacterInput.Name, rpcCharacterInput.Name)
	assert.Equal(t, entityCharacterInput.Tags, rpcCharacterInput.Tags)
	assert.Equal(t, entityCharacterInput.Gender, rpcCharacterInput.Gender)
	assert.Equal(t, entityCharacterInput.Vision.Name, rpcCharacterInput.Vision.Name)
	assert.Equal(t, entityCharacterInput.Vision.Description, rpcCharacterInput.Vision.Description)
	assert.Len(t, entityCharacterInput.Categories, len(rpcCharacterInput.Categories))
	assert.Len(t, entityCharacterInput.Metrics, len(rpcCharacterInput.Metrics))
	assert.Equal(t, entityCharacterInput.Categories[0].ID, rpcCharacterInput.Categories[0].Id)
	assert.Equal(t, entityCharacterInput.Categories[0].Name, rpcCharacterInput.Categories[0].Name)
	assert.Equal(t, entityCharacterInput.Categories[0].Description, rpcCharacterInput.Categories[0].Description)
	assert.Equal(t, entityCharacterInput.Categories[0].Style.Color, rpcCharacterInput.Categories[0].Style.Color)
	assert.Equal(t, entityCharacterInput.Categories[0].Style.Icon, rpcCharacterInput.Categories[0].Style.Icon)
	assert.Equal(t, entityCharacterInput.Metrics[0].ID, rpcCharacterInput.Metrics[0].Id)
	assert.Equal(t, entityCharacterInput.Metrics[0].CategoryID, rpcCharacterInput.Metrics[0].CategoryId)
	assert.Equal(t, entityCharacterInput.Metrics[0].Name, rpcCharacterInput.Metrics[0].Name)
	assert.Equal(t, entityCharacterInput.Metrics[0].Value, float64(rpcCharacterInput.Metrics[0].Value))
	assert.Equal(t, entityCharacterInput.Metrics[0].Unit, rpcCharacterInput.Metrics[0].Unit)
}
