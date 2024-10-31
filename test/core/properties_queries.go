package core

var CreateMetricPropertyQuery = `
mutation CreateMetricProperty ($characterID: ObjectID!, $metricID: ObjectID!, $input: MetricPropertyInput!) {
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
mutation UpdateMetricProperty ($id: ObjectID!, $characterID: ObjectID!, $metricID: ObjectID!, $input: MetricPropertyInput!) {
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
mutation DeleteMetricProperty($id: ObjectID!, $characterID: ObjectID!, $metricID: ObjectID!) {
    deleteMetricProperty (id: $id, characterID: $characterID, metricID: $metricID)  {
        id
        name
        type
        unit
        value
    }
}`
