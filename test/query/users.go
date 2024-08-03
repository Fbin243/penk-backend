package query

var User = `query { 
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
	}
}`

var UpdateAccount = `mutation UpdateAccount($name: String, $imageURL: String, $currentCharacterID: String, $autoSnapshot: Boolean) {
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
