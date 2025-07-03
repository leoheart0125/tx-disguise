package shared

import (
	"testing"
)

func TestRingBuffer_PushPop(t *testing.T) {
	rb := NewRingBuffer[int](3)

	// Test empty buffer
	if !rb.IsEmpty() {
		t.Error("Buffer should be empty initially")
	}

	// Push elements
	rb.Push(1)
	rb.Push(2)
	if rb.IsEmpty() {
		t.Error("Buffer should not be empty after push")
	}

	// Pop elements
	v, err := rb.Pop()
	if err != nil || v != 1 {
		t.Errorf("Expected 1, got %v, err: %v", v, err)
	}
	v, err = rb.Pop()
	if err != nil || v != 2 {
		t.Errorf("Expected 2, got %v, err: %v", v, err)
	}

	// Pop from empty
	_, err = rb.Pop()
	if err == nil {
		t.Error("Expected error when popping from empty buffer")
	}
}

func TestRingBuffer_Overwrite(t *testing.T) {
	rb := NewRingBuffer[int](2)
	rb.Push(1)
	rb.Push(2)
	if !rb.IsFull() {
		t.Error("Buffer should be full")
	}
	// Overwrite oldest
	rb.Push(3)
	if !rb.IsFull() {
		t.Error("Buffer should still be full after overwrite")
	}
	vals := rb.GetAll()
	if len(vals) != 2 || vals[0] != 2 || vals[1] != 3 {
		t.Errorf("Expected [2 3], got %v", vals)
	}
}

func TestRingBuffer_GetAll(t *testing.T) {
	rb := NewRingBuffer[string](3)
	rb.Push("a")
	rb.Push("b")
	vals := rb.GetAll()
	if len(vals) != 2 || vals[0] != "a" || vals[1] != "b" {
		t.Errorf("Expected [a b], got %v", vals)
	}

	rb.Push("c")
	rb.Push("d") // Overwrites "a"
	vals = rb.GetAll()
	if len(vals) != 3 || vals[0] != "b" || vals[1] != "c" || vals[2] != "d" {
		t.Errorf("Expected [b c d], got %v", vals)
	}
}

func TestRingBuffer_Capacity(t *testing.T) {
	rb := NewRingBuffer[float64](5)
	if rb.Capacity() != 5 {
		t.Errorf("Expected capacity 5, got %d", rb.Capacity())
	}
}
