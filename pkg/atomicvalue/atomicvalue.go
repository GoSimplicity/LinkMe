package atomicvalue

import (
	"sync/atomic"
)

type AtomicValue[T any] struct {
	value atomic.Value
}

func (a *AtomicValue[T]) Load() T {
	val := a.value.Load()
	if val == nil {
		var zero T
		return zero
	}
	return val.(T)
}

func (a *AtomicValue[T]) Store(val T) {
	a.value.Store(val)
}
