package shared

import (
	"errors"
	"sync"
)

type RingBuffer[T any] struct {
	buffer []T
	size   int
	start  int
	end    int
	full   bool
	mu     sync.Mutex
}

func NewRingBuffer[T any](capacity int) *RingBuffer[T] {
	return &RingBuffer[T]{
		buffer: make([]T, capacity),
		size:   capacity,
	}
}

// Push 寫入一筆資料，滿了會覆蓋最舊資料
func (r *RingBuffer[T]) Push(value T) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.buffer[r.end] = value
	r.end = (r.end + 1) % r.size

	if r.full {
		r.start = (r.start + 1) % r.size
	} else if r.end == r.start {
		r.full = true
	}
}

// Pop 取出最舊資料
func (r *RingBuffer[T]) Pop() (T, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var zero T
	if r.IsEmpty() {
		return zero, errors.New("buffer is empty")
	}
	val := r.buffer[r.start]
	r.buffer[r.start] = zero // 清除被取出的資料
	r.start = (r.start + 1) % r.size
	r.full = false
	return val, nil
}

// GetAll 回傳所有元素（依照順序）
func (r *RingBuffer[T]) GetAll() []T {
	r.mu.Lock()
	defer r.mu.Unlock()

	var out []T
	if r.IsEmpty() {
		return out
	}

	var count int
	if r.full {
		count = r.size
	} else if r.end >= r.start {
		count = r.end - r.start
	} else {
		count = r.size - r.start + r.end
	}

	for i := 0; i < count; i++ {
		idx := (r.start + i) % r.size
		out = append(out, r.buffer[idx])
	}
	return out
}

func (r *RingBuffer[T]) IsEmpty() bool {
	return !r.full && r.start == r.end
}

func (r *RingBuffer[T]) IsFull() bool {
	return r.full
}

func (r *RingBuffer[T]) Capacity() int {
	return r.size
}
