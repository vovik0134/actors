package periodic

import (
	"actors/internal"
	"context"
	"fmt"
	"time"
)

func (a *Actor) Act(_ context.Context) error {
	fmt.Println(fmt.Sprintf("%s: periodic called", time.Now().Format(time.RFC3339Nano)))
	return nil
}

func (a *Actor) Notify(actions chan<- internal.ActFunc) {
	go func() {
		for {
			actions <- a.Act
			time.Sleep(a.period)
		}
	}()
}

func New(period time.Duration) *Actor {
	return &Actor{period: period}
}

type Actor struct {
	period time.Duration
}
