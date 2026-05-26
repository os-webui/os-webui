package lazy

type funcOption[T any] struct {
	f func(T)
}

func (fdo *funcOption[T]) apply(do T) {
	fdo.f(do)
}
func newFuncOption[T any](f func(T)) *funcOption[T] {
	return &funcOption[T]{
		f: f,
	}
}

type Option[T any] interface {
	apply(T)
}
