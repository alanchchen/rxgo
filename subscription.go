package rx

import (
	"sync"
)

type Subscription interface {
	Unsubscribe()
}

type subscription struct {
	cleanUp func()
}

func newSubscription(o observable, subscribeFunc func(item Any) Any, handlerFunc Any) Subscription {
	resultCh := make(DataChannel)
	stopCh := make(StopChannel)
	subsMu := sync.Mutex{}
	subs := make(map[Subscription]Subscription)

	go o(resultCh, stopCh)
	go func() {
		defer func() {
			for sub := range subs {
				sub.Unsubscribe()
			}
		}()

		for {
			select {
			case item, ok := <-resultCh:
				if !ok {
					return
				}

				if ob, ok := item.(Observable); ok {
					sub := ob.Subscribe(handlerFunc)
					subsMu.Lock()
					subs[sub] = sub
					subsMu.Unlock()
				} else {
					subscribeFunc(item)
				}
			case <-stopCh:
				return
			}
		}
	}()

	return &subscription{
		cleanUp: func() {
			close(stopCh)
		},
	}
}

func (sub *subscription) Unsubscribe() {
	if sub.cleanUp != nil {
		sub.cleanUp()
	}
}
