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

type Metric struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

type Category struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Style       Style    `json:"style"`
	Metrics     []Metric `json:"metrics"`
}

type Template struct {
	Name        string     `json:"name"`
	Emoji       string     `json:"emoji"`
	Description string     `json:"description"`
	Color       string     `json:"color"`
	Categories  []Category `json:"categories"`
}

func readTemplatesFromJSON() error {
	jsonFile, err := os.Open("cmd/db/templates.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	topicsMap := map[string][]Template{}
	jsonParser := json.NewDecoder(jsonFile)
	jsonParser.Decode(&topicsMap)

	templateRepo := mongorepo.NewTemplateRepo(mongodb.GetDBManager().DB)
	templateTopicRepo := mongorepo.NewTemplateTopicRepo(mongodb.GetDBManager().DB)

	templates := []entity.Template{}

	// Loop through the categories and insert them into the database
	for topicName, topicTemplates := range topicsMap {
		category := entity.TemplateTopic{
			BaseEntity:  &base.BaseEntity{},
			Name:        topicName,
			Description: fmt.Sprintf("This topic contains templates related to %s", topicName),
		}

		// Insert the category into the database
		_, err = templateTopicRepo.InsertOne(context.Background(), &category)
		if err != nil {
			return err
		}

		// Loop through the templates and insert them into the slide
		for _, template := range topicTemplates {
			repoTemplate := mapToRepoTemplate(template)
			repoTemplate.TopicID = category.ID
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
		TopicID:     primitive.NewObjectID().Hex(),
		Style: entity.TemplateStyle{
			Color: t.Color,
			Icon:  t.Emoji,
		},
		Categories: []entity.TemplateCategory{},
	}

	for _, category := range t.Categories {
		metrics := []entity.TemplateMetric{}
		for _, metric := range category.Metrics {
			metrics = append(metrics, entity.TemplateMetric{
				Name:  metric.Name,
				Value: metric.Value,
				Unit:  metric.Unit,
			})
		}

		template.Categories = append(template.Categories, entity.TemplateCategory{
			Name:        category.Name,
			Description: category.Description,
			Style: entity.CategoryStyle{
				Color: category.Style.Color,
				Icon:  category.Style.Icon,
			},
			Metrics: metrics,
		})
	}

	return template
}
