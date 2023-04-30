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

	event := triggerable.Event(func(ctx context.Context, fire func()) {
		// triggers only once
		fire()
	})

	expectedRetries := 5
	actualRetries := 0

	retryOnErrorFunc := func(err error) (bool, func(context.Context)) {
		if err != testErr {
			return false, nil
		}

		if actualRetries == expectedRetries {
			return false, nil
		}

		return true, func(ctx context.Context) {
			time.Sleep(100 * time.Millisecond)
			actualRetries++
		}
	}

	action := triggerable.Action(runFunc, triggerable.WithName("retryable"), triggerable.WithRetryOnError(retryOnErrorFunc))
	retryableTrigger := triggerable.Trigger(action, event)

	loop := triggerable.Loop(logger, retryableTrigger)

	if err := loop.Run(ctx); err != testErr {
		t.Fatalf("loop failed with unexpected error: %s", err)
	}

	if expectedRetries != actualRetries {
		t.Fatalf("expected %d calls, but called %d times", expectedRetries, actualRetries)
	}
}
