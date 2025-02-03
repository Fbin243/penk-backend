package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/errors"
	"tenkhours/services/core/entity"

	rdb "tenkhours/pkg/db/redis"

	"go.mongodb.org/mongo-driver/mongo"
)

type TemplateBusiness struct {
	TemplateRepo         ITemplateRepo
	TemplateCategoryRepo ITemplateCategoryRepo
}

func NewTemplateBusiness(templateRepo ITemplateRepo, templateCategoryRepo ITemplateCategoryRepo) *TemplateBusiness {
	return &TemplateBusiness{
		TemplateRepo:         templateRepo,
		TemplateCategoryRepo: templateCategoryRepo,
	}
}

// GetTemplates returns all templates
func (biz *TemplateBusiness) GetTemplates(ctx context.Context) ([]entity.Template, error) {
	_, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.Unauthorized()
	}

	return biz.TemplateRepo.FindAll(ctx)
}

// Get template category by ID
func (biz *TemplateBusiness) GetTemplateCategory(ctx context.Context, id string) (*entity.TemplateCategory, error) {
	_, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.Unauthorized()
	}

	category, err := biz.TemplateCategoryRepo.FindByID(ctx, id)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return category, nil
}
