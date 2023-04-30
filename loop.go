package triggerable

import (
	"context"
	"fmt"
	"time"
)

func (l *loopImpl) Run(ctx context.Context) error {
	actions := make(chan action, len(l.triggerable))
	for _, n := range l.triggerable {
		n.EnqueueActionWhenTriggered(actions)
	}

	for {
		select {
		case <-ctx.Done():
			l.logger.Info(ctx, fmt.Sprintf("loop stopped"))
			return nil
		case a := <-actions:
			l.logger.Debug(ctx, fmt.Sprintf("running action %q", a.Name()))
			if err := a.Run(ctx); err != nil {
				l.logger.Info(ctx, fmt.Sprintf("action %q failed with error: %s", a.Name(), err))
				if a.RetryAfterTimeout() != NoRetryTimeout {
					l.scheduleRetry(ctx, actions, a)
				}
			}
		}
	}
}

func Loop(logger logger, triggerable ...triggerable) *loopImpl {
	return &loopImpl{logger: logger, triggerable: triggerable}
}

type loopImpl struct {
	logger      logger
	triggerable []triggerable
}

func (l *loopImpl) scheduleRetry(ctx context.Context, actions chan<- action, action action) {
	l.logger.Info(ctx, fmt.Sprintf("scheduling retry for action %q in %s", action.Name(), action.RetryAfterTimeout().String()))

	go func() {
		select {
		case <-ctx.Done():
		case <-time.After(action.RetryAfterTimeout()):
			actions <- action
		}
	}()
}

type triggerable interface {
	EnqueueActionWhenTriggered(chan<- action)
}
