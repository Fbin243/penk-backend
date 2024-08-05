package characters

var CreateCharacterVariable = map[string]interface{}{
	"name":   "Character name",
	"gender": false,
	"avatar": "avatar.png",
	"tags":   []interface{}{"#Tag1", "#Tag2"},
}

var UpdateCharacterVariable = func(idCharacter interface{}) map[string]interface{} {
	return map[string]interface{}{
		"id":     idCharacter,
		"name":   "Update name",
		"gender": true,
		"avatar": "update-avatar.png",
		"tags":   []interface{}{"#update_tag_1", "#update_tag_2"},
	}
}
