package handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

func TestExecPassesDeadlineIntoJob(t *testing.T) {
	t.Parallel()
	done := make(chan error, 1)
	h := New("test", func(ctx context.Context, _ model.QueueJob) error {
		<-ctx.Done()
		done <- ctx.Err()
		return ctx.Err()
	}, JobTimeout(20*time.Millisecond))
	err := Exec(context.Background(), h, model.QueueJob{})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline error, got %v", err)
	}
	select {
	case received := <-done:
		if !errors.Is(received, context.DeadlineExceeded) {
			t.Fatalf("handler got %v", received)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out handler did not observe context cancellation")
	}
}

func TestExecReturnsHandlerError(t *testing.T) {
	t.Parallel()
	sentinel := errors.New("temporarily unavailable")
	h := New("test", func(context.Context, model.QueueJob) error { return sentinel })
	if err := Exec(context.Background(), h, model.QueueJob{}); !errors.Is(err, sentinel) {
		t.Fatalf("handler error not retained for retry: %v", err)
	}
}
