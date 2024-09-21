package core

import "strconv"

func castType(valueType string, value string) any {
	var newValue interface{}

	switch valueType {
	case "Text":
		newValue = value

	case "Number":
		newValue, _ = strconv.ParseFloat(value, 64)

	default:
		newValue = value
	}

	return newValue
}

func convertListToSlice(list interface{}) []string {
	var slice []string

	for _, item := range list.([]interface{}) {
		slice = append(slice, item.(string))
	}

	return slice
}
