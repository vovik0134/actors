package triggerable

import (
	"context"
	"fmt"
	"sync"
)

func (l *loopImpl) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg := &sync.WaitGroup{}
	wg.Add(len(l.triggerable))
	var err error

	for _, t := range l.triggerable {
		go func(t trigger) {
			if err = l.listen(ctx, t.Triggered(ctx)); err != nil {
				cancel()
			}
			wg.Done()
		}(t)
	}

	wg.Wait()
	return err
}

func Loop(logger logger, triggerable ...trigger) *loopImpl {
	return &loopImpl{logger: logger, triggerable: triggerable}
}

type loopImpl struct {
	logger      logger
	triggerable []trigger
}

func (l *loopImpl) listen(ctx context.Context, actions chan action) error {
	for {
		select {
		case <-ctx.Done():
			l.logger.Info(ctx, fmt.Sprintf("loop stopped"))
			return nil
		case a := <-actions:
			l.logger.Debug(ctx, fmt.Sprintf("running action %q", a.Name()))

			if err := a.Run(ctx); err != nil {
				l.logger.Info(ctx, fmt.Sprintf("action %q failed with error: %s", a.Name(), err))
				retry, retryFunc := a.RetryOnError(err)
				if !retry {
					return err
				}

				l.logger.Info(ctx, fmt.Sprintf("retrying action %q", a.Name()))

				// retry function can contain time.Sleep call
				// or something blocking,
				// so we running it in separate goroutine
				go func() {
					if retryFunc != nil {
						retryFunc(ctx)
					}

					select {
					case <-ctx.Done():
						return
					case actions <- a:
					}
				}()
			}
		}
	}
}
