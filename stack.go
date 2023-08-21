package main

import (
	"fmt"
	"sync"
)

type Stack[T any] struct {
	data []T
	len  int
	mu   sync.Mutex
}

func (s *Stack[T]) Push(el T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.len < len(s.data) {
		// no need to do s.len + 1 here, since we are increment later
		s.data[s.len] = el
	} else {
		s.data = append(
			s.data,
			el,
		)
	}
	s.len += 1
}

func (s *Stack[T]) Pop() T {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.data) == 0 {
		panic(fmt.Errorf("tried to pop from empty stack: %+v", s))
	}
	s.len -= 1
	return s.data[s.len]
}

func (s *Stack[T]) Peek() T {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.data[s.len - 1]
}

func (s *Stack[T]) IsEmpty() bool {
	return s.len == 0
}

func (s *Stack[T]) Size() int {
	return s.len
}

func (s *Stack[T]) Get(i int) T {
	return s.data[i]
}
