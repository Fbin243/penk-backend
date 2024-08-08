package users

var UserQuery = `query { 
	user {
        autoSnapshot
        availableSnapshots
        createdAt
        currentCharacterID
        email
        firebaseUID
        id
        imageURL
        name
        updatedAt
        characters {
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
    }
}`

var UpdateAccountQuery = `mutation UpdateAccount($name: String, $imageURL: String, $currentCharacterID: String, $autoSnapshot: Boolean) {
	updateAccount(
		input: {
			name: $name
			imageURL: $imageURL
			currentCharacterID: $currentCharacterID
			autoSnapshot: $autoSnapshot
		}
	) {
		autoSnapshot
		availableSnapshots
		createdAt
		currentCharacterID
		email
		firebaseUID
		id
		imageURL
		name
		updatedAt
	}
}`
