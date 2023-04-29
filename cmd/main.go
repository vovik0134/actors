package main

import (
	"actors/internal/eventual"
	"actors/internal/periodic"
	"actors/internal/synchrnoizer"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	events := make(chan struct{}, 1)

	p := periodic.New(time.Second)
	e := eventual.New(events)
	s := synchrnoizer.New(p, e)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-stop
		fmt.Println(fmt.Sprintf("%s: got stop signal", time.Now().Format(time.RFC3339Nano)))
		cancel()
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			time.Sleep(time.Duration(rand.Int31n(10)) * time.Second)
			events <- struct{}{}
		}
	}()

	_ = s.Synchronize(ctx)
}
