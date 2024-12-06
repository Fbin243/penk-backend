package core

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "tenkhours/services/core/business"
    "tenkhours/services/core/repo"
)

func TestTemplates(t *testing.T) {
    templateRepo := repo.NewTemplateRepo("../data/templates.csv")
    templatesBiz := business.NewTemplatesBusiness(templateRepo)
    
    t.Run("GetTemplates", func(t *testing.T) {
        templates, err := templatesBiz.GetTemplates(context.Background())
        assert.NoError(t, err)
        assert.NotEmpty(t, templates)
        
        template := templates[0]
        assert.Equal(t, "1", template.ID)
        assert.Equal(t, "Frontend Developer", template.Name)
        assert.NotEmpty(t, template.Metrics)
    })

    t.Run("GetTemplateByID", func(t *testing.T) {
        template, err := templatesBiz.GetTemplateByID(context.Background(), "1")
        assert.NoError(t, err)
        assert.NotNil(t, template)
        assert.Equal(t, "Frontend Developer", template.Name)
    })
}