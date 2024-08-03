package variables

var CreateCustomMetric = func(characterId interface{}) map[string]interface{} {
	return map[string]interface{}{
		"characterID": characterId,
		"name":        "Metric name",
		"description": "This is the custom metric description",
		"style": map[string]interface{}{
			"color": "#000000",
			"icon":  "icon.png",
		},
	}
}
