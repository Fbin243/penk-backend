package core

var CreateMetricPropertyQuery = `
mutation CreateMetricProperty ($characterID: ID!, $metricID: ID!, $input: MetricPropertyInput!) {
    createMetricProperty(
        characterID: $characterID
        metricID: $metricID
        input: $input
    ) {
        id
        name
        type
        unit
        value
    }
}`

var UpdateMetricPropertyQuery = `
mutation UpdateMetricProperty ($id: ID!, $characterID: ID!, $metricID: ID!, $input: MetricPropertyInput!) {
    updateMetricProperty(
        id: $id
        characterID: $characterID
        metricID: $metricID
        input: $input
    ) {
        id
        name
        type
        unit
        value
    }
}`

var DeleteMetricPropertyQuery = `
mutation DeleteMetricProperty($id: ID!, $characterID: ID!, $metricID: ID!) {
    deleteMetricProperty (id: $id, characterID: $characterID, metricID: $metricID)  {
        id
        name
        type
        unit
        value
    }
}`
