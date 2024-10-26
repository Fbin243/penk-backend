package pineline

import (
	"context"
)

type Stage interface {
	Exec(ctx *context.Context) error
}

func Pineline(stages ...Stage) func(ctx *context.Context) error {
	return func(ctx *context.Context) error {
		for _, stage := range stages {
			err := stage.Exec(ctx)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
