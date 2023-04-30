package examples

import (
	"context"
	"testing"
	"triggerable"
)

func TestEventTriggered(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	logger := &Logger{}

	eventsN := 5
	events := make(chan struct{}, 1)

	// generate events and then close events channel
	go func() {
		for i := 0; i < eventsN; i++ {
			select {
			case <-ctx.Done():
				return
			case events <- struct{}{}:
			}
		}

		close(events)
	}()

	actualCalls := 0
	runFunc := func(ctx context.Context) error {
		if actualCalls == eventsN {
			cancelFunc()
			return nil
		}

		actualCalls++
		return nil
	}

	event := triggerable.Event(func(ctx context.Context, fire func()) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-events:
				fire()
			}
		}
	})

	action := triggerable.Action(runFunc, triggerable.WithName("event"))
	eventTrigger := triggerable.Trigger(action, event)

	loop := triggerable.Loop(logger, eventTrigger)

	if err := loop.Run(ctx); err != nil {
		t.Fatalf("loop failed with unexpected error: %s", err)
	}

	if eventsN != actualCalls {
		t.Fatalf("expected %d calls, but called %d times", eventsN, actualCalls)
	}
}
