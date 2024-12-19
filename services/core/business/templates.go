package business

import (
	"context"
	"tenkhours/pkg/auth"
	"tenkhours/pkg/errors"
	"tenkhours/services/core/repo"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TemplatesBusiness struct {
	TemplatesRepo           *repo.TemplatesRepo
	TemplatesCategoriesRepo *repo.TemplateCategoriesRepo
}

func NewTemplatesBusiness(templatesRepo *repo.TemplatesRepo, templateCategoriesRepo *repo.TemplateCategoriesRepo) *TemplatesBusiness {
	return &TemplatesBusiness{
		TemplatesRepo:           templatesRepo,
		TemplatesCategoriesRepo: templateCategoriesRepo,
	}
}

// GetTemplates returns all templates
func (biz *TemplatesBusiness) GetTemplates(ctx context.Context) ([]repo.Template, error) {
	_, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, errors.ErrorUnauthorized
	}

	return biz.TemplatesRepo.FindAll()
}

// Get template category by ID
func (biz *TemplatesBusiness) GetTemplateCategory(ctx context.Context, id primitive.ObjectID) (*repo.TemplateCategory, error) {
	_, ok := ctx.Value(auth.ProfileKey).(repo.Profile)
	if !ok {
		return nil, errors.ErrorUnauthorized
	}

	category, err := biz.TemplatesCategoriesRepo.FindByID(id)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return category, nil
}
