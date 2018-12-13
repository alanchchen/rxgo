package rx

type Any interface{}
type DataChannel chan Any
type StopChannel chan struct{}

type observable func(DataChannel, StopChannel)

type Filter interface {
	Filter(item Any) bool
}

type FilterFunc func(item Any) bool

func (f FilterFunc) Filter(item Any) bool {
	return f(item)
}

type Transformer interface {
	Transform(item Any) Any
}

type TransformerFunc func(item Any) Any

func (f TransformerFunc) Transform(item Any) Any {
	return f(item)
}
