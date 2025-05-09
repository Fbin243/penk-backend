package utils

import (
	"encoding/json"
)

func PrettyJSON(data any) string {
	prettyJSON, _ := json.MarshalIndent(data, "", "  ")
	return string(prettyJSON)
}
