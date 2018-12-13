package main

import (
	"fmt"
	"time"

	rx "github.com/alanchchen/rxgo"
)

func main() {
	intervalSource := rx.Interval(1 * time.Second)

	oddSource := rx.From(intervalSource).Filter(func(i int) bool {
		return i%2 != 0
	})

	evenSource := rx.From(intervalSource).Filter(func(i int) bool {
		return i%2 == 0
	})

	oddSub := oddSource.Subscribe(func(i int) int {
		fmt.Println("oddSource", i)
		return 0
	})
	defer oddSub.Unsubscribe()

	evenSub := evenSource.Subscribe(func(i int) int {
		fmt.Println("evenSource", i)
		return 0
	})
	defer evenSub.Unsubscribe()

	<-time.After(10 * time.Second)
}
