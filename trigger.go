package triggerable

import (
	"context"
)

func (p *triggerImpl) Triggered(ctx context.Context) chan action {
	actions := make(chan action)
	fired := p.event.Fired(ctx)

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

func Trigger(action action, event event) *triggerImpl {
	return &triggerImpl{action: action, event: event}
}

type triggerImpl struct {
	action action
	event  event
}

type trigger interface {
	Triggered(ctx context.Context) chan action
}
