package metrics

var CreateCustomMetricVariable = func(characterId interface{}) map[string]interface{} {
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

var UpdateCustomMetricVariable = map[string]interface{}{
	"name":        "Update metric name",
	"description": "This is the updated custom metric description",
	"style": map[string]interface{}{
		"color": "#123456",
		"icon":  "update_icon.png",
	},
}
