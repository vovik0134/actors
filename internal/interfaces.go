package internal

import "context"

type ActFunc func(ctx context.Context) error

type Actor interface {
	Act(ctx context.Context) error
	Notify(chan<- ActFunc)
}
