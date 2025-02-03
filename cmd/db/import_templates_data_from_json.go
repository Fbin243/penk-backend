package db

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"tenkhours/pkg/db/base"
	"tenkhours/services/core/entity"

	mongodb "tenkhours/pkg/db/mongo"

	mongorepo "tenkhours/services/core/repo/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Style struct {
	Icon  string `json:"icon"`
	Color string `json:"color"`
}

type Property struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Unit  string `json:"unit"`
	Type  string `json:"type"`
}

type Metric struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Style       Style      `json:"style"`
	Properties  []Property `json:"properties"`
}

type Template struct {
	Name        string   `json:"name"`
	Emoji       string   `json:"emoji"`
	Description string   `json:"description"`
	Color       string   `json:"color"`
	Metrics     []Metric `json:"metrics"`
}

func readTemplatesFromJSON() error {
	jsonFile, err := os.Open("cmd/db/templates.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	categoriesMap := map[string][]Template{}
	jsonParser := json.NewDecoder(jsonFile)
	jsonParser.Decode(&categoriesMap)

	templateRepo := mongorepo.NewTemplateRepo(mongodb.GetDBManager().DB)
	templateCategoryRepo := mongorepo.NewTemplateCategoryRepo(mongodb.GetDBManager().DB)

	templates := []entity.Template{}

	// Loop through the categories and insert them into the database
	for cateName, cateTemplates := range categoriesMap {
		category := entity.TemplateCategory{
			BaseEntity:  &base.BaseEntity{},
			Name:        cateName,
			Description: fmt.Sprintf("This category contains templates related to %s", cateName),
		}

		// Insert the category into the database
		_, err = templateCategoryRepo.InsertOne(context.Background(), &category)
		if err != nil {
			return err
		}

		// Loop through the templates and insert them into the slide
		for _, template := range cateTemplates {
			repoTemplate := mapToRepoTemplate(template)
			repoTemplate.CategoryID = category.ID
			templates = append(templates, repoTemplate)
		}
	}

	// Insert the templates into the database
	_, err = templateRepo.InsertMany(context.Background(), templates)
	if err != nil {
		return err
	}

	return nil
}

func mapToRepoTemplate(t Template) entity.Template {
	template := entity.Template{
		BaseEntity:  &base.BaseEntity{},
		Name:        t.Name,
		Description: t.Description,
		CategoryID:  primitive.NewObjectID().Hex(),
		Style: entity.TemplateStyle{
			Color: t.Color,
			Icon:  t.Emoji,
		},
		Metrics: []entity.TemplateMetric{},
	}

	for _, metric := range t.Metrics {
		properties := []entity.TemplateProperty{}
		for _, property := range metric.Properties {
			properties = append(properties, entity.TemplateProperty{
				Name:  property.Name,
				Value: property.Value,
				Unit:  property.Unit,
				Type:  entity.MetricPropertyType(property.Type),
			})
		}

		template.Metrics = append(template.Metrics, entity.TemplateMetric{
			Name:        metric.Name,
			Description: metric.Description,
			Style: entity.MetricStyle{
				Color: metric.Style.Color,
				Icon:  metric.Style.Icon,
			},
			Properties: properties,
		})
	}

	return template
}
