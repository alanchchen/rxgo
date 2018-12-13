package rx

import "reflect"

func (o observable) Filter(handlerFunc Any) Observable {
	genericFunc, err := newGenericFunc(
		"Filter", "handlerFunc", handlerFunc,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	filterFunc := func(item Any) bool {
		v := reflect.ValueOf(genericFunc.Call(item))
		if v.Kind() == reflect.Bool {
			return v.Bool()
		}

		return v.Interface() == reflect.Zero(v.Type()).Interface()
	}

	return makeSubject(func(dataDestinationCh DataChannel, stop StopChannel) {
		deliver(o, dataDestinationCh, stop, FilterFunc(filterFunc), nil)
	})
}

func (o observable) Take(n int) Observable {
	index := 0

	return o.Filter(func(v Any) bool {
		defer func() {
			index++
		}()

		return index < n
	})
}

func (o observable) Skip(n int) Observable {
	index := 0

	return o.Filter(func(v Any) bool {
		defer func() {
			index++
		}()

		return index >= n
	})
}
