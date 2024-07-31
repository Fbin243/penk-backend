package pineline

import (
	"context"
)

type Stage func(ctx context.Context) error

func Pineline(stages ...Stage) Stage {
	return func(ctx context.Context) error {
		for _, stage := range stages {
			err := stage(ctx)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
