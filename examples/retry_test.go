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

	testErr := fmt.Errorf("test error")

	runFunc := func(ctx context.Context) error {
		return testErr
	}

	eventFunc := func(ctx context.Context, trigger func(context.Context)) {
		// triggers only once
		trigger(ctx)
	}

	expectedRetries := 5
	actualRetries := 0

	retryOnErrorFunc := func(ctx context.Context, err error) bool {
		if err != testErr {
			return false
		}

		if actualRetries == expectedRetries {
			cancelFunc()
			return false
		}
		time.Sleep(100 * time.Millisecond)
		actualRetries++

		return true
	}

	action := triggerable.Action(runFunc, triggerable.WithName("retryable"), triggerable.WithRetryOnError(retryOnErrorFunc))
	retryableTrigger := triggerable.Trigger(ctx, action, eventFunc)

	loop := triggerable.Loop(logger, retryableTrigger)

	if err := loop.Run(ctx); err != nil {
		t.Fatalf("loop failed with unexpected error: %s", err)
	}

	if expectedRetries != actualRetries {
		t.Fatalf("expected %d calls, but called %d times", expectedRetries, actualRetries)
	}
}
