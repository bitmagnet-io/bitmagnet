package semaphore

import (
	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
	"github.com/marusama/semaphore/v2"
)

type Semaphore = semaphore.Semaphore

func NewAtomic(size atomic.Reader[int]) (semaphore.Semaphore, func() int) {
	currentSize := size.Get()

	sem := semaphore.New(currentSize)

	unsubscribe := size.Subscribe(func(newSize int) {
		if newSize != currentSize {
			currentSize = newSize
			sem.SetLimit(currentSize)
		}
	})

	return sem, unsubscribe
}
