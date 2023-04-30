package triggerable

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func (l *loopImpl) Run(ctx context.Context) error {
	wg := &sync.WaitGroup{}
	wg.Add(len(l.triggerable))

	for _, t := range l.triggerable {
		go func(t trigger) {
			l.listen(ctx, t.Triggered())
			wg.Done()
		}(t)
	}

	wg.Wait()
	return nil
}

func Loop(logger logger, triggerable ...trigger) *loopImpl {
	return &loopImpl{logger: logger, triggerable: triggerable}
}

type loopImpl struct {
	logger      logger
	triggerable []trigger
}

func (l *loopImpl) listen(ctx context.Context, actions chan action) {
	for {
		select {
		case <-ctx.Done():
			l.logger.Info(ctx, fmt.Sprintf("loop stopped"))
			return
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

func (l *loopImpl) scheduleRetry(ctx context.Context, actions chan action, action action) {
	l.logger.Info(ctx, fmt.Sprintf("scheduling retry for action %q in %s", action.Name(), action.RetryAfterTimeout().String()))

	go func() {
		select {
		case <-ctx.Done():
		case <-time.After(action.RetryAfterTimeout()):
			select {
			case actions <- action:
			}
		}
	}()
}
