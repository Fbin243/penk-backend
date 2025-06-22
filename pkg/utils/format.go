package utils

import (
	"encoding/json"
)

func PrettyJSON(data interface{}) string {
	prettyJSON, _ := json.MarshalIndent(data, "", "  ")
	return string(prettyJSON)
}
