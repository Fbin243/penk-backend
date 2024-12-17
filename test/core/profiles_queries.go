package core

var ProfileQuery = `query { 
	profile {
        autoSnapshot
        limitedCharacterNumber
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
            gender
            id
            limitedMetricNumber
            name
            tags
            totalFocusedTime
            profileID
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

var UpdateProfileQuery = `mutation UpdateProfile($input: ProfileInput!) {
	updateProfile(input: $input) {
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
        availableSnapshots
        autoSnapshot
        characters {
            id
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
    }
}
`
