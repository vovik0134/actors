package triggerable

import (
	"context"
)

func (p *triggerImpl) Triggered() chan action {
	actions := make(chan action)

	go p.eventFunc(p.ctx, func(ctx context.Context) {
		select {
		case <-ctx.Done():
		case actions <- p.action:
		}
	})

	return actions
}

func Trigger(
	ctx context.Context,
	action action,
	eventFunc func(ctx context.Context, trigger func(context.Context)),
) *triggerImpl {
	return &triggerImpl{
		ctx:       ctx,
		action:    action,
		eventFunc: eventFunc,
	}
}

type triggerImpl struct {
	ctx       context.Context
	action    action
	eventFunc func(ctx context.Context, trigger func(context.Context))
}

type trigger interface {
	Triggered() chan action
}
