package rx

import (
	"reflect"
	"sync"
)

func flatObservable(o observable, stop StopChannel) Any {
	dataSourceCh := make(DataChannel)
	stopCh := make(StopChannel)

	go deliver(o, dataSourceCh, stopCh, nil, nil)

	once := sync.Once{}
	s := reflect.Value{}

	for {
		select {
		case i, ok := <-dataSourceCh:
			if !ok {
				return s.Interface()
			}

			once.Do(func() {
				sliceType := reflect.SliceOf(reflect.TypeOf(i))
				s = reflect.MakeSlice(sliceType, 0, 0)
			})

			s = reflect.Append(s, reflect.ValueOf(i))
		case <-stop:
			close(stopCh)
		}
	}
}

func deliver(o observable, dataDestinationCh DataChannel, stop StopChannel, filter Filter, transformer Transformer) {
	if filter == nil {
		filter = FilterFunc(func(item Any) bool {
			return true
		})
	}

	if transformer == nil {
		transformer = TransformerFunc(func(item Any) Any {
			return item
		})
	}

	dataSourceCh := make(DataChannel)
	stopCh := make(StopChannel)

	defer close(dataDestinationCh)

	go o(dataSourceCh, stopCh)
	for {
		select {
		case v, ok := <-dataSourceCh:
			if !ok {
				return
			}

			if filter.Filter(v) {
				dataDestinationCh <- transformer.Transform(v)
			}
		case <-stop:
			close(stopCh)
			return
		}
	}
}
