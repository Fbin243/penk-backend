package core

var ProfileQuery = `
query { 
	profile {
        createdAt
        currentCharacterID
        email
        firebaseUID
        id
        imageURL
        name
        updatedAt
        characters {
            gender
            id
            name
            tags
            profileID
            categories {
                description
                id
                name
                style {
                    color
                    icon
                }
                metrics {
                     id
                    name
                    unit
                    value
                }
            }
            metrics {
                id
                name
                unit
                value
            }
        }
    }
}`

var UpdateProfileQuery = `
mutation UpdateProfile($input: ProfileInput!) {
	updateProfile(input: $input) {
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

var DeleteProfileQuery = `
mutation DeleteProfile {
    deleteProfile {
        id
        createdAt
        updatedAt
        name
        email
        firebaseUID
        imageURL
        currentCharacterID
    }
}`
