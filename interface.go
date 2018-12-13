package rx

type Observable interface {
	SwitchMap(fn Any) Observable
	FlatMap(fn Any) Observable
	Map(fn Any) Observable
	Filter(fn Any) Observable
	Take(n int) Observable
	Skip(n int) Observable
	Subscribe(fn Any) Subscription
}
