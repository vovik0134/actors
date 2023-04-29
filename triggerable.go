package triggerable

import (
	"context"
	"time"
)

const NoRetryTimeout = -1 * time.Second

func (p *Triggerable) NotifyWhenTriggered(actions chan<- *Action) {
	triggered := make(chan struct{}, 1)
	triggerFunc := func(ctx context.Context) {
		select {
		case <-ctx.Done():
		case triggered <- struct{}{}:
		}
	}

	go p.notifyFunc(p.ctx, triggerFunc)

	action := &Action{Run: p.runFunc, RetryAfterTimeout: p.retryAfterTimeout}

	go func() {
		for {
			select {
			case <-p.ctx.Done():
				return
			case <-triggered:
				actions <- action
			}
		}
	}()
}

func New(ctx context.Context, opts ...Option) *Triggerable {
	p := &Triggerable{
		ctx:               ctx,
		runFunc:           func(ctx context.Context) error { return nil },
		notifyFunc:        func(context.Context, TriggerFunc) {},
		retryAfterTimeout: NoRetryTimeout,
	}

	for _, o := range opts {
		o(p)
	}

	return p
}

func WithRunFunc(runFunc RunFunc) Option {
	return func(p *Triggerable) {
		p.runFunc = runFunc
	}
}

func WithNotifyFunc(notifyFunc NotifyFunc) Option {
	return func(p *Triggerable) {
		p.notifyFunc = notifyFunc
	}
}

func WithRetryAfterTimeout(retryAfterTimeout time.Duration) Option {
	return func(p *Triggerable) {
		p.retryAfterTimeout = retryAfterTimeout
	}
}

type Option func(p *Triggerable)

type Triggerable struct {
	ctx               context.Context
	runFunc           RunFunc
	notifyFunc        NotifyFunc
	retryAfterTimeout time.Duration
}
