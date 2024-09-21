package core

var CreateCharacterQuery = `
mutation CreateCharacter($name: String!, $gender: Boolean, $tags: [String]) {
	createCharacter(
		input: { 
			name: $name, 
			gender: $gender, 
			tags: $tags 
		}
	) {
        gender
        id
        limitedMetricNumber
        name
        tags
        totalFocusedTime
        userID
        customMetrics {
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
}`

var UpdateCharacterQuery = `
mutation UpdateCharacter($id: ID!, $gender: Boolean, $name: String, $tags: [String]) {
	updateCharacter(
		id: $id
		input: {
			gender: $gender,
			name: $name
			tags: $tags,
		}
	) {
		id
		gender
		name
		tags
	}
}`

var DeleteCharacterQuery = `
mutation DeleteCharacter($id: ID!) {
	deleteCharacter(id: $id) {
		gender
		id
		limitedMetricNumber
		name
		tags
		totalFocusedTime
		userID
	}
}`
