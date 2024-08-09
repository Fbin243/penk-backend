package characters

var CreateCharacterQuery = `
mutation CreateCharacter($name: String!, $gender: Boolean, $avatar: String, $tags: [String]) {
	createCharacter(
		input: { 
			name: $name, 
			gender: $gender, 
			avatar: $avatar, 
			tags: $tags 
		}
	) {
		avatar
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
mutation UpdateCharacter($id: ID!, $avatar: String, $gender: Boolean, $name: String, $tags: [String]) {
	updateCharacter(
		id: $id
		input: {
			avatar: $avatar,
			gender: $gender,
			name: $name
			tags: $tags,
		}
	) {
		id
		avatar
		gender
		name
		tags
	}
}`

var DeleteCharacterQuery = `
mutation DeleteCharacter($id: ID!) {
	deleteCharacter(id: $id) {
		avatar
		gender
		id
		limitedMetricNumber
		name
		tags
		totalFocusedTime
		userID
	}
}`
