package rx

import (
	"reflect"
	"sync"
	"time"
)

// From creates a new Observable from an Observable, a slice or an array.
func From(source Any) Observable {
	typ := reflect.TypeOf(source)

	switch typ.Kind() {
	case reflect.Array, reflect.Slice:
		return makeSubject(func(dataDestinationCh DataChannel, stop StopChannel) {
			defer close(dataDestinationCh)
			arr := reflect.ValueOf(source)

			for i := 0; i < arr.Len(); i++ {
				select {
				case <-stop:
					return
				default:
					dataDestinationCh <- arr.Index(i).Interface()
				}
			}
		})
	case reflect.Func:
		o := reflect.ValueOf(source)
		if fn, ok := o.Interface().(observable); ok {
			return makeSubject(fn)
		}
	}

	panic("unsupported source type")
}

// Range generates a sequence of integral numbers within a specified range.
func Range(from, to int) Observable {
	if from > to {
		from, to = to, from
	}

	return makeSubject(func(dataDestinationCh DataChannel, stop StopChannel) {
		defer close(dataDestinationCh)
		for i := from; i < from+to; i++ {
			select {
			case <-stop:
				return
			default:
				dataDestinationCh <- i
			}
		}
	})
}

// Interval generates a sequence of integral numbers for every given duration.
func Interval(duration time.Duration) Observable {
	timer := time.NewTicker(duration)
	index := 0
	return makeSubject(func(dataDestinationCh DataChannel, stop StopChannel) {
		defer close(dataDestinationCh)
		defer timer.Stop()
		for {
			select {
			case <-stop:
				return
			case <-timer.C:
				dataDestinationCh <- index
				index++
			}
		}
	})
}

func (o observable) Subscribe(handlerFunc Any) Subscription {
	subscribeGenericFunc, err := newGenericFunc(
		"Subscribe", "handlerFunc", handlerFunc,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	subscribeFunc := func(item Any) Any {
		return subscribeGenericFunc.Call(item)
	}

	return newSubscription(o, subscribeFunc, handlerFunc)
}

func makeSubject(source observable) Observable {
	mu := sync.Mutex{}
	children := make(map[DataChannel]DataChannel)
	sendAll := func(data Any) {
		for key := range children {
			key <- data
		}
	}
	closeAll := func() {
		for key := range children {
			close(key)
			mu.Lock()
			delete(children, key)
			mu.Unlock()
		}
	}

	dataSourceCh := make(DataChannel)
	stopSourceCh := make(StopChannel)

	go source(dataSourceCh, stopSourceCh)

	return observable(func(dataDestinationCh DataChannel, stopCh StopChannel) {
		mu.Lock()
		children[dataDestinationCh] = dataDestinationCh
		mu.Unlock()

		for {
			select {
			case data, ok := <-dataSourceCh:
				if ok {
					sendAll(data)
				} else {
					closeAll()
					return
				}
			case <-stopCh:
				mu.Lock()
				delete(children, dataDestinationCh)
				mu.Unlock()
				if len(children) == 0 {
					close(stopSourceCh)
				}
				return
			}
		}
	})
}
