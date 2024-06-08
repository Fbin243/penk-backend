package character

type MetricsType interface {
}

type CustomMetricData struct {
	ID          string `json:"id"`
	CharacterID string `json:"character_id"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Value       string `json:"value"`
}

type CharacterData struct {
	ID               string             `json:"id"`
	UserID           string             `json:"user_id"`
	Name             string             `json:"name"`
	Tags             []string           `json:"tags"`
	TotalFocusedTime float64            `json:"total_focused_time"`
	CustomMetrics    []CustomMetricData `json:"custom_metrics"`
}

func createCustomMetric(id string, characterID string, metricType string, name string, value string) CustomMetricData {
	return CustomMetricData{
		ID:          id,
		CharacterID: characterID,
		Type:        metricType,
		Name:        name,
		Value:       value,
	}
}

func createCharacterArrays() []*CharacterData {
	dummyCharacters := make([]*CharacterData, 0)

	dummyCharacters = append(dummyCharacters, &CharacterData{
		ID:               "1",
		UserID:           "22f83",
		Name:             "Learn Go",
		Tags:             []string{"#developer", "#blockchain"},
		TotalFocusedTime: 48.7,

		CustomMetrics: []CustomMetricData{createCustomMetric("C01", "1", "string", "Title", "Intern"), createCustomMetric("C02", "1", "int", "Project", "2")},
	})

	dummyCharacters = append(dummyCharacters, &CharacterData{
		ID:               "2",
		UserID:           "13h867",
		Name:             "Learn React",
		Tags:             []string{"#developer", "#React"},
		TotalFocusedTime: 22,
		CustomMetrics:    []CustomMetricData{createCustomMetric("C03", "2", "string", "Title", "Junior"), createCustomMetric("C04", "2", "int", "Project", "3")},
	})

	dummyCharacters = append(dummyCharacters, &CharacterData{
		ID:               "3",
		UserID:           "ee35h6",
		Name:             "Learn C",
		Tags:             []string{"#developer", "#C++"},
		TotalFocusedTime: 22,
		CustomMetrics:    []CustomMetricData{createCustomMetric("C05", "3", "string", "Title", "Senior"), createCustomMetric("C06", "3", "string", "Have projects?", "yes")},
	})

	return dummyCharacters
}
