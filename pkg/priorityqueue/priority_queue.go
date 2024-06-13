package priorityqueue

import (
	"container/heap"
	"errors"
)

var ErrOutOfCapacity = errors.New("out of capacity")

type PriorityQueue[T any] struct {
	items    []T
	capacity int
	less     func(a, b T) bool
}

func NewPriorityQueue[T any](capacity int, less func(a, b T) bool) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		items:    make([]T, 0, capacity),
		capacity: capacity,
		less:     less,
	}
}

func (pq *PriorityQueue[T]) Len() int {
	return len(pq.items)
}

func (pq *PriorityQueue[T]) Less(i, j int) bool {
	return pq.less(pq.items[i], pq.items[j])
}

func (pq *PriorityQueue[T]) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
}

func (pq *PriorityQueue[T]) Push(x interface{}) {
	pq.items = append(pq.items, x.(T))
}

func (pq *PriorityQueue[T]) Pop() interface{} {
	old := pq.items
	n := len(old)
	item := old[n-1]
	pq.items = old[0 : n-1]
	return item
}

func (pq *PriorityQueue[T]) Enqueue(item T) error {
	if pq.Len() >= pq.capacity {
		return ErrOutOfCapacity
	}
	heap.Push(pq, item)
	return nil
}

func (pq *PriorityQueue[T]) Dequeue() (T, error) {
	if pq.Len() == 0 {
		var zero T
		return zero, errors.New("queue is empty")
	}
	item := heap.Pop(pq).(T)
	return item, nil
}
