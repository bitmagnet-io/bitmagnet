package runner

import (
	"context"
	"sync"
)

func OnceShutdowner(s Shutdowner) Shutdowner {
	var (
		once sync.Once
		err  error
	)

	return func(ctx context.Context) error {
		once.Do(func() {
			err = s.Call(ctx)
		})

		return err
	}
}
