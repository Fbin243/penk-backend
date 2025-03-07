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
        categories {
            id
            name
            description
            style {
                color
                icon
            }
            metrics {
                id
                name
                value
                unit
            }
        }
        metrics {
            id
            name
            value
            unit
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
		name
		tags
		profileID
	}
}`
