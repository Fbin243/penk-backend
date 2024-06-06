package character

type CustomMetricData struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type CharacterData struct {
	ID               string             `json:"id"`
	Name             string             `json:"name"`
	Tags             []string           `json:"tags"`
	TotalFocusedTime float64            `json:"total_focused_time"`
	CustomMetrics    []CustomMetricData `json:"custom_metrics"`
}

func createCustomMetric(id string, name string, value string) CustomMetricData {
	return CustomMetricData{
		ID:    id,
		Name:  name,
		Value: value,
	}
}

func createCharacterArrays() []*CharacterData {
	dummyCharacters := make([]*CharacterData, 0)

	dummyCharacters = append(dummyCharacters, &CharacterData{
		ID:               "1",
		Name:             "Learn Go",
		Tags:             []string{"#developer", "#blockchain"},
		TotalFocusedTime: 48.7,

		CustomMetrics: []CustomMetricData{createCustomMetric("C01", "Title", "Intern"), createCustomMetric("C02", "Project", "2")},
	})

	dummyCharacters = append(dummyCharacters, &CharacterData{
		ID:               "2",
		Name:             "Learn React",
		Tags:             []string{"#developer", "#React"},
		TotalFocusedTime: 22,
		CustomMetrics:    []CustomMetricData{createCustomMetric("C03", "Title", "Junior"), createCustomMetric("C04", "Project", "2")},
	})

	dummyCharacters = append(dummyCharacters, &CharacterData{
		ID:               "3",
		Name:             "Learn C",
		Tags:             []string{"#developer", "#C++"},
		TotalFocusedTime: 22,
		CustomMetrics:    []CustomMetricData{createCustomMetric("C05", "Title", "Senior"), createCustomMetric("C06", "IDK", "rafcefsdvsd")},
	})

	return dummyCharacters
}
