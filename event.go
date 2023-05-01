package triggerable

import "context"

func (e *eventImpl) Fired(ctx context.Context) <-chan struct{} {
	fired := make(chan struct{})

	go func() {
		e.eventFunc(ctx, func() {
			select {
			case <-ctx.Done():
				return
			case fired <- struct{}{}:
				return
			}
		})
	}()

	return fired
}

func Event(eventFunc func(ctx context.Context, fire func())) *eventImpl {
	return &eventImpl{eventFunc: eventFunc}
}

type eventImpl struct {
	eventFunc func(ctx context.Context, fire func())
}

type event interface {
	Fired(ctx context.Context) <-chan struct{}
}
