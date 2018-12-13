package rx

import (
	"fmt"
	"reflect"
)

func (o observable) SwitchMap(handlerFunc Any) Observable {
	return o.FlatMap(handlerFunc)
}

func (o observable) FlatMap(handlerFunc Any) Observable {
	genericFunc, err := newGenericFunc(
		"FlatMap", "handlerFunc", handlerFunc,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(Observable))),
	)
	if err != nil {
		panic(err)
	}

	flatMapFunc := func(item Any) Any {
		return genericFunc.Call(item)
	}

	return makeSubject(func(dataDestinationCh DataChannel, stop StopChannel) {
		deliver(o, dataDestinationCh, stop, nil, TransformerFunc(func(item Any) Any {
			mappedObservable, ok := flatMapFunc(item).(observable)
			if !ok {
				panic(fmt.Errorf("Call using %s as type %s", reflect.TypeOf(item), reflect.TypeOf(new(Observable))))
			}

			return flatObservable(mappedObservable, stop)
		}))
	})
}

func (o observable) Map(handlerFunc Any) Observable {
	genericFunc, err := newGenericFunc(
		"Map", "handlerFunc", handlerFunc,
		simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(genericType))),
	)
	if err != nil {
		panic(err)
	}

	mapFunc := func(item Any) Any {
		return genericFunc.Call(item)
	}

	return makeSubject(func(dataDestinationCh DataChannel, stop StopChannel) {
		deliver(o, dataDestinationCh, stop, nil, TransformerFunc(mapFunc))
	})
}
