package stack

import "container/list"

type Stack struct {
	l *list.List
}

func NewStack() *Stack {
	l := list.New()
	return &Stack{l}
}

func (s *Stack) Push(value interface{}) {
	s.l.PushBack(value)
}

// Pop removes the last element and return
func (s *Stack) Pop() interface{} {
	e := s.l.Back()
	if e != nil {
		s.l.Remove(e)
		return e.Value
	}
	return nil
}

func (s *Stack) Back() interface{} {
	e := s.l.Back()
	if e != nil {
		return e.Value
	}
	return nil
}

func (s *Stack) Len() int {
	return s.l.Len()
}

func (s *Stack) Empty() bool {
	return s.l.Len() == 0
}
