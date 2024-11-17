package business

import (
    "context"
    "tenkhours/services/core/repo"
)

type TemplatesBusiness struct {
    templateRepo *repo.TemplateRepo
}

func NewTemplatesBusiness(templateRepo *repo.TemplateRepo) *TemplatesBusiness {
    return &TemplatesBusiness{
        templateRepo: templateRepo,
    }
}

func (biz *TemplatesBusiness) GetTemplates(ctx context.Context) ([]repo.Template, error) {
    return biz.templateRepo.GetTemplates()
}

func (biz *TemplatesBusiness) GetTemplateByID(ctx context.Context, id string) (*repo.Template, error) {
    templates, err := biz.templateRepo.GetTemplates()
    if err != nil {
        return nil, err
    }

    for _, template := range templates {
        if template.ID == id {
            return &template, nil
        }
    }
    return nil, nil
}