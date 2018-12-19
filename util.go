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

type option struct {
	transformer Transformer
	filter      Filter
}

type deliverOption func(*option)

func withFilter(filter Filter) deliverOption {
	return func(o *option) {
		o.filter = filter
	}
}

func withTransformer(transformer Transformer) deliverOption {
	return func(o *option) {
		o.transformer = transformer
	}
}

func deliver(o observable, dataDestinationCh DataChannel, stop StopChannel, opts ...deliverOption) {
	option := &option{
		filter: FilterFunc(func(item Any) bool {
			return true
		}),
		transformer: TransformerFunc(func(item Any) Any {
			return item
		}),
	}

	for _, opt := range opts {
		opt(option)
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

			if option.filter.Filter(v) {
				dataDestinationCh <- option.transformer.Transform(v)
			}
		case <-stop:
			close(stopCh)
			return
		}
	}
}
