package examples

import (
	"context"
	"fmt"
	"time"
)

func (*Logger) Debug(_ context.Context, args ...any) {
	fmt.Println(time.Now().Format(logFormat) + ": DEBUG " + fmt.Sprint(args...))
}

func (*Logger) Info(_ context.Context, args ...any) {
	fmt.Println(time.Now().Format(logFormat) + ": INFO  " + fmt.Sprint(args...))
}

type Logger struct {
}

const logFormat = "2006-01-02 15:04:05.000000 -0700 MST"
