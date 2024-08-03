package query

var CreateCharacter = `mutation CreateCharacter($name: String!, $gender: Boolean, $avatar: String!, $tags: [String]) {
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
	}
}`

var UpdateCharacter = `mutation UpdateCharacter($id: String!, $avatar: String, $gender: Boolean, $name: String, $tags: [String]) {
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

var DeleteCharacter = `mutation DeleteCharacter($id: String!) {
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
