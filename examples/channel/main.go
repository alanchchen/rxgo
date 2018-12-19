package main

import (
	"fmt"
	"time"

	rx "github.com/alanchchen/rxgo"
)

func main() {
	timerSource := rx.From(time.Tick(1 * time.Second))

	sub := timerSource.Subscribe(func(i time.Time) time.Time {
		fmt.Println(i)
		return i
	})
	defer sub.Unsubscribe()

	<-time.After(10 * time.Second)
}
