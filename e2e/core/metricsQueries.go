package core

var CreateCustomMetricQuery = `
mutation CreateCustomMetric($characterID: ObjectID!, $name: String, $style: MetricStyleInput, $description: String) {
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
mutation UpdateCustomMetric ($id: ObjectID!, $characterID: ObjectID!, $name: String, $style: MetricStyleInput, $description: String, $properties: 
[MetricPropertyInput!]) {
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
}
`

var DeleteCustomMetricQuery = `
mutation DeleteCustomMetric ($id: ObjectID!, $characterID: ObjectID!) {
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
mutation ResetCustomMetric ($id: ObjectID!, $characterID: ObjectID!) {
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
