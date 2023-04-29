package examples

import (
	"context"
	"math"
	"testing"
	"time"
	"triggerable"
)

func TestPeriodic(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	interval := 100 * time.Millisecond
	expectedCalls := 5
	actualCalls := 0

	runFunc := func(ctx context.Context) error {
		if actualCalls == expectedCalls {
			cancelFunc()
			return nil
		}
		actualCalls++

		return nil
	}

	notifyFunc := func(ctx context.Context, triggerFunc triggerable.TriggerFunc) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(interval):
				triggerFunc(ctx)
			}
		}
	}

	periodic := triggerable.New(
		ctx,
		triggerable.WithRunFunc(runFunc),
		triggerable.WithNotifyFunc(notifyFunc),
	)

	loop := triggerable.Loop(periodic)

	startedAt := time.Now()
	if err := loop.Run(ctx); err != nil {
		t.Fatalf("loop failed with unexpected error: %s", err)
	}

	expectedDuration := time.Duration(expectedCalls) * interval
	actualDuration := time.Since(startedAt)

	if math.Abs(float64(actualDuration-expectedDuration)) >= float64(expectedDuration)*0.25 {
		t.Fatalf("run take too long. expected %s, but actual is %s", expectedDuration.String(), actualDuration.String())
	}

	if expectedCalls != actualCalls {
		t.Fatalf("expected %d calls, but called %d times", expectedCalls, actualCalls)
	}
}
