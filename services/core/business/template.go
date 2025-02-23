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

func NewTemplateBusiness(templateRepo ITemplateRepo, templateTopicRepo ITemplateCategoryRepo) *TemplateBusiness {
	return &TemplateBusiness{
		TemplateRepo:         templateRepo,
		TemplateCategoryRepo: templateTopicRepo,
	}
}

// GetTemplates returns all templates
func (biz *TemplateBusiness) GetTemplates(ctx context.Context) ([]entity.Template, error) {
	_, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	return biz.TemplateRepo.FindAll(ctx)
}

// Get template category by ID
func (biz *TemplateBusiness) GetTemplateCategory(ctx context.Context, id string) (*entity.TemplateTopic, error) {
	_, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	category, err := biz.TemplateCategoryRepo.FindByID(ctx, id)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return category, nil
}
