package core

var UpsertCharacterQuery = `
mutation UpsertCharacter($input: CharacterInput!) {
    upsertCharacter(input: $input) {
        id
        createdAt
        updatedAt
        profileID
        name
        gender
        tags
        totalFocusedTime
        limitedMetricNumber
        customMetrics {
            id
            name
            description
            time
            limitedPropertyNumber
            style {
                color
                icon
            }
            properties {
                id
                name
                type
                value
                unit
            }
        }
    }
}`

var DeleteCharacterQuery = `
mutation DeleteCharacter($id: ID!) {
	deleteCharacter(id: $id) {
		gender
		id
		createdAt
		updatedAt
		limitedMetricNumber
		name
		tags
		totalFocusedTime
		profileID
	}
}`
