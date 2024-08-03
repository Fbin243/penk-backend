package query

var CreateCustomMetric = `mutation CreateCustomMetric($characterID: String!, $name: String, $style: MetricStyleInput, $description: String) {
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
