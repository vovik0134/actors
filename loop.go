package triggerable

import (
	"context"
	"time"
)

func (l *LoopImpl) Run(ctx context.Context) error {
	actions := make(chan *Action, len(l.notifiables))
	for _, n := range l.notifiables {
		n.NotifyWhenTriggered(actions)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case action := <-actions:
			if err := action.Run(ctx); err != nil {
				if action.RetryAfterTimeout != NoRetryTimeout {
					l.scheduleRetry(ctx, actions, action)
				}
			}
		}
	}
}

func Loop(notifiables ...Notifiable) *LoopImpl {
	return &LoopImpl{notifiables: notifiables}
}

type LoopImpl struct {
	notifiables []Notifiable
}

func (l *LoopImpl) scheduleRetry(ctx context.Context, actions chan<- *Action, action *Action) {
	go func() {
		select {
		case <-ctx.Done():
		case <-time.After(action.RetryAfterTimeout):
			actions <- action
		}
	}()
}
