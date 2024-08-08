package common

import "fmt"

func ErrNotFound(key ContextKey) error {
	return fmt.Errorf("%s is not found", key)
}
