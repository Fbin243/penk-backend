package core

var CreateCharacterQuery = `
mutation CreateCharacter($name: String!, $gender: Boolean!, $tags: [String!], $customMetrics: [CustomMetricInput!]) {
	createCharacter(
		input: { 
			name: $name, 
			gender: $gender, 
			tags: $tags 
			customMetrics: $customMetrics
		}
	) {
        gender
        id
				createdAt
        updatedAt
        limitedMetricNumber
        name
        tags
        totalFocusedTime
        profileID
        customMetrics {
            description
            name
            style {
                color
                icon
            }
        }
	}
}`

var UpdateCharacterQuery = `
mutation UpdateCharacter($id: ObjectID!, $gender: Boolean!, $name: String!, $tags: [String!]) {
	updateCharacter(
		id: $id
		input: {
			gender: $gender,
			name: $name
			tags: $tags,
		}
	) {
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
mutation DeleteCharacter($id: ObjectID!) {
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
