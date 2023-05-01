package triggerable

import (
	"context"
)

func (p *triggerImpl) Triggered(ctx context.Context) chan action {
	fired := p.event.Fired(ctx)
	actions := make(chan action)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-fired:
				select {
				case actions <- p.action:
				}
			}
		}
	}()

	return actions
}

func Trigger(event event, action action) *triggerImpl {
	return &triggerImpl{event: event, action: action}
}

type triggerImpl struct {
	event  event
	action action
}

type trigger interface {
	Triggered(ctx context.Context) chan action
}
