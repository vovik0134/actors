package triggerable

import (
	"context"
	"fmt"
)

func (p *triggerableImpl) EnqueueActionWhenTriggered(actions chan<- action) {
	triggered := make(chan struct{}, 1)

	triggerFunc := func(ctx context.Context) {
		p.logger.Debug(ctx, fmt.Sprintf("action %q triggered", p.action.Name()))

		select {
		case <-ctx.Done():
		case triggered <- struct{}{}:
		}
	}

	go p.notifyFunc(p.ctx, triggerFunc)

	go func() {
		for {
			select {
			case <-p.ctx.Done():
				return
			case <-triggered:
				actions <- p.action
			}
		}
	}()
}

func New(
	ctx context.Context,
	logger logger,
	action action,
	notifyFunc func(ctx context.Context, triggerFunc func(context.Context)),
) *triggerableImpl {
	return &triggerableImpl{
		ctx:        ctx,
		logger:     logger,
		action:     action,
		notifyFunc: notifyFunc,
	}
}

type triggerableImpl struct {
	ctx        context.Context
	logger     logger
	action     action
	notifyFunc func(ctx context.Context, triggerFunc func(context.Context))
}
