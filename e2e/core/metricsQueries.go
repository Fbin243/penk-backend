package core

var CreateCustomMetricQuery = `
mutation CreateCustomMetric($characterID: ID!, $name: String, $style: MetricStyleInput, $description: String) {
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

var UpdateCustomMetricQuery = `
mutation UpdateCustomMetric ($id: ID!, $characterID: ID!, $name: String, $style: MetricStyleInput, $description: String, $properties: 
[MetricPropertyInput]) {
     updateCustomMetric(
        id: $id
        characterID: $characterID
        input: {
            properties: $properties
            name: $name
            description: $description
            style: $style
        }
    ) {
        description
        id
        limitedPropertyNumber
        name
        time
        properties {
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
}
`

var DeleteCustomMetricQuery = `
mutation DeleteCustomMetric ($id: ID!, $characterID: ID!) {
    deleteCustomMetric(id: $id, characterID: $characterID) {
        description
        id
        limitedPropertyNumber
        name
        time
        properties {
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

var ResetCustomMetricQuery = `
mutation ResetCustomMetric ($id: ID!, $characterID: ID!) {
    resetCustomMetric(id: $id, characterID: $characterID) {
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
