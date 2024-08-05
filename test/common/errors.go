package common

import "fmt"

func ErrNotFoundInContext(key ContextKey) error {
	return fmt.Errorf("%s is not found in context at stage", key)
}
