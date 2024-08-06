package metrics

var CreateCustomMetricQuery = `mutation CreateCustomMetric($characterID: ID!, $name: String, $style: MetricStyleInput, $description: String) {
	createCustomMetric(
		characterID: $characterID
		input: {
			name: $name,
			style: $style,
			description: $description
		}
	) {
		description
		id
		limitedPropertyNumber
		name
		time
		properties {
			id
			name
			type
			unit
			value
		}
		style {
			color
			icon
		}
	}
}`
