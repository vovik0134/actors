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

	notifyFunc := func(ctx context.Context, triggerFunc triggerable.TriggerFunc) {
		// triggers only once
		triggerFunc(ctx)
	}

	retryable := triggerable.New(
		ctx,
		triggerable.WithRunFunc(runFunc),
		triggerable.WithNotifyFunc(notifyFunc),
		triggerable.WithRetryAfterTimeout(100*time.Millisecond),
	)

	loop := triggerable.Loop(retryable)

	if err := loop.Run(ctx); err != nil {
		t.Fatalf("loop failed with unexpected error: %s", err)
	}

	if expectedCalls != actualCalls {
		t.Fatalf("expected %d calls, but called %d times", expectedCalls, actualCalls)
	}
}
