package priorityqueue

import (
	"container/heap"
	"errors"
)

var (
	ErrOutOfCapacity = errors.New("队列已满")
	ErrEmptyQueue    = errors.New("队列为空")
)

type PriorityQueue[T any] struct {
	items    []T
	capacity int
	less     func(a, b T) bool
}

func NewPriorityQueue[T any](capacity int, less func(a, b T) bool) *PriorityQueue[T] {
	if capacity <= 0 {
		capacity = 1
	}
	pq := &PriorityQueue[T]{
		items:    make([]T, 0, capacity),
		capacity: capacity,
		less:     less,
	}
	heap.Init(pq)
	return pq
}

// Len 返回队列的长度
func (pq *PriorityQueue[T]) Len() int {
	return len(pq.items)
}

// Less 比较两个元素的大小
func (pq *PriorityQueue[T]) Less(i, j int) bool {
	if pq.less == nil {
		return false
	}
	return pq.less(pq.items[i], pq.items[j])
}

// Swap 交换两个元素
func (pq *PriorityQueue[T]) Swap(i, j int) {
	if i >= 0 && i < len(pq.items) && j >= 0 && j < len(pq.items) {
		pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
	}
}

// Push 添加一个元素
func (pq *PriorityQueue[T]) Push(x interface{}) {
	if item, ok := x.(T); ok {
		pq.items = append(pq.items, item)
	}
}

// Pop 移除并返回最小元素
func (pq *PriorityQueue[T]) Pop() interface{} {
	if len(pq.items) == 0 {
		return nil
	}
	old := pq.items
	n := len(old)
	item := old[n-1]
	pq.items = old[0 : n-1]
	return item
}

// Enqueue 添加一个元素
func (pq *PriorityQueue[T]) Enqueue(item T) error {
	if pq.Len() >= pq.capacity {
		return ErrOutOfCapacity
	}
	heap.Push(pq, item)
	return nil
}

// Dequeue 移除并返回最小元素
func (pq *PriorityQueue[T]) Dequeue() (T, error) {
	if pq.Len() == 0 {
		var zero T
		return zero, ErrEmptyQueue
	}
	result := heap.Pop(pq)
	if result == nil {
		var zero T
		return zero, ErrEmptyQueue
	}
	return result.(T), nil
}
