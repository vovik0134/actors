package examples

import (
	"context"
	"fmt"
	"testing"
	"time"
	"triggerable"
)

func TestRetry(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	logger := &Logger{}

	expectedCalls := 5
	actualCalls := 0

	runFunc := func(ctx context.Context) error {
		if actualCalls == expectedCalls {
			cancelFunc()
			return nil
		}
		actualCalls++

		return fmt.Errorf("test error")
	}

	notifyFunc := func(ctx context.Context, triggerFunc func(ctx context.Context)) {
		// triggers only once
		triggerFunc(ctx)
	}

	action := triggerable.Action(runFunc, triggerable.WithName("retryable"), triggerable.WithRetryAfterTimeout(100*time.Millisecond))
	retryable := triggerable.New(ctx, logger, action, notifyFunc)

	loop := triggerable.Loop(logger, retryable)

	if err := loop.Run(ctx); err != nil {
		t.Fatalf("loop failed with unexpected error: %s", err)
	}

	if expectedCalls != actualCalls {
		t.Fatalf("expected %d calls, but called %d times", expectedCalls, actualCalls)
	}
}
