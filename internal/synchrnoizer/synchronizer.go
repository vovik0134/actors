package synchrnoizer

import (
	"actors/internal"
	"context"
	"fmt"
	"time"
)

func (s *Synchronizer) Synchronize(ctx context.Context) error {
	for {
		iterationTimeout := time.After(time.Minute)

		select {
		case <-ctx.Done():
			fmt.Println("synchronization stopped")
			return nil
		case action := <-s.actions:
			if err := action(ctx); err != nil {
				fmt.Println("error occurred running action", err)
			}
		case <-iterationTimeout:
		}
	}
}

func New(actors ...internal.Actor) *Synchronizer {
	actions := make(chan internal.ActFunc, 1)
	for _, a := range actors {
		a.Notify(actions)
	}

	return &Synchronizer{actions: actions}
}

type Synchronizer struct {
	actions <-chan internal.ActFunc
}
