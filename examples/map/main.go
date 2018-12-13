package main

import (
	"fmt"
	"time"

	rx "github.com/alanchchen/rxgo"
)

func main() {
	source := rx.Interval(1 * time.Second)

	source.Map(func(i int) string {
		return time.Now().Format("2006-01-02T15:04:05")
	}).Subscribe(func(s string) string {
		fmt.Println(s)
		return s
	})

	<-time.After(10 * time.Second)
}
