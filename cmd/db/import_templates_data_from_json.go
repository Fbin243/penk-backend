package db

import (
	"encoding/json"
	"fmt"
	"os"
	"tenkhours/pkg/db"
	"tenkhours/services/core/repo"

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

	templatesRepo := repo.NewTemplatesRepo(db.GetDBManager().DB)
	templateCategoriesRepo := repo.NewTemplateCategoriesRepo(db.GetDBManager().DB)

	templates := []repo.Template{}

	// Loop through the categories and insert them into the database
	for cateName, cateTemplates := range categoriesMap {
		category := repo.TemplateCategory{
			BaseModel:   &db.BaseModel{},
			Name:        cateName,
			Description: fmt.Sprintf("This category contains templates related to %s", cateName),
		}

		// Insert the category into the database
		_, err = templateCategoriesRepo.InsertOne(&category)
		if err != nil {
			return err
		}

		// Loop throught the templates and insert them into the slide
		for _, template := range cateTemplates {
			repoTemplate := mapToRepoTemplate(template)
			repoTemplate.CategoryID = category.ID
			templates = append(templates, repoTemplate)
		}
	}

	// Insert the templates into the database
	_, err = templatesRepo.InsertMany(templates)
	if err != nil {
		return err
	}

	return nil
}

func mapToRepoTemplate(t Template) repo.Template {
	template := repo.Template{
		BaseModel:   &db.BaseModel{},
		Name:        t.Name,
		Description: t.Description,
		CategoryID:  primitive.NewObjectID(),
		Style: repo.TemplateStyle{
			Color: t.Color,
			Icon:  t.Emoji,
		},
		Metrics: []repo.TemplateMetric{},
	}

	for _, metric := range t.Metrics {
		properties := []repo.TemplateProperty{}
		for _, property := range metric.Properties {
			properties = append(properties, repo.TemplateProperty{
				Name:  property.Name,
				Value: property.Value,
				Unit:  property.Unit,
				Type:  repo.MetricPropertyType(property.Type),
			})
		}

		template.Metrics = append(template.Metrics, repo.TemplateMetric{
			Name:        metric.Name,
			Description: metric.Description,
			Style: repo.MetricStyle{
				Color: metric.Style.Color,
				Icon:  metric.Style.Icon,
			},
			Properties: properties,
		})
	}

	return template
}
