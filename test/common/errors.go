package common

import "fmt"

func ErrNotFoundInContext(key string) error {
	return fmt.Errorf("%s is not found", key)
}
