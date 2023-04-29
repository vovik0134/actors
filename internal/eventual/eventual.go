package eventual

import (
	"actors/internal"
	"context"
	"fmt"
	"time"
)

func (e *Actor) Act(_ context.Context) error {
	fmt.Println(fmt.Sprintf("%s: got event", time.Now().Format(time.RFC3339Nano)))
	return nil
}

func (e *Actor) Notify(actions chan<- internal.ActFunc) {
	go func() {
		for {
			<-e.events
			actions <- e.Act
		}
	}()
}

func New(events <-chan struct{}) *Actor {
	return &Actor{events: events}
}

type Actor struct {
	events <-chan struct{}
}
