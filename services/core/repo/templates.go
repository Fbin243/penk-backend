package repo

import (
    "encoding/csv"
    "encoding/json"
    "os"
)

type TemplateRepo struct {
    filePath string
}

type Template struct {
    ID          string           `json:"id"`
    Name        string           `json:"name"`
    Description string           `json:"description"` 
    Icon        string           `json:"icon"`
    Category    string           `json:"category"`
    Metrics     []TemplateMetric `json:"metrics"`
}

type TemplateMetric struct {
    Name        string             `json:"name"`
    Description string             `json:"description"`
    Style       MetricStyle        `json:"style"`
    Properties  []TemplateProperty `json:"properties"`
}

type TemplateProperty struct {
    Name  string `json:"name"`
    Type  string `json:"type"`  
    Value string `json:"value"`
    Unit  string `json:"unit"`
}

func NewTemplateRepo(filePath string) *TemplateRepo {
    return &TemplateRepo{
        filePath: filePath,
    }
}

func (r *TemplateRepo) GetTemplates() ([]Template, error) {
    file, err := os.Open(r.filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records, err := reader.ReadAll()
    if err != nil {
        return nil, err
    }

    templates := make([]Template, 0)
    for _, record := range records[1:] { // Skip header
        var metrics []TemplateMetric
        err = json.Unmarshal([]byte(record[5]), &metrics)
        if err != nil {
            return nil, err
        }

        template := Template{
            ID:          record[0],
            Name:        record[1],
            Description: record[2],
            Icon:        record[3],
            Category:    record[4],
            Metrics:     metrics,
        }
        templates = append(templates, template)
    }

    return templates, nil
}