package triggerable

import (
	"context"
	"time"
)

type TriggerFunc func(ctx context.Context)
type NotifyFunc func(ctx context.Context, triggerFunc TriggerFunc)

type RunFunc func(ctx context.Context) error

type Action struct {
	Run               RunFunc
	RetryAfterTimeout time.Duration
}

type Notifiable interface {
	NotifyWhenTriggered(chan<- *Action)
}
