package main

import (
	"fmt"
	"time"

	rx "github.com/alanchchen/rxgo"
)

func main() {
	intervalSource := rx.Interval(1 * time.Second)

	sub1 := intervalSource.Subscribe(func(i int) int {
		fmt.Println("sub1", i)
		return i
	})
	defer sub1.Unsubscribe()

	<-time.After(1 * time.Second)

	sub2 := intervalSource.Subscribe(func(i int) int {
		fmt.Println("sub2", i)
		return i
	})
	defer sub2.Unsubscribe()

	<-time.After(10 * time.Second)
}
