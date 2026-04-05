package main

import (
	"fmt"
)

type stack struct {
	head *node
	tail *node
}

type node struct {
	data int
	pre  *node
	next *node
}

func (s *stack) push(data int) {
	n := &node{data: data}
	if s.head == nil && s.head == s.tail {
		s.head, s.tail = n, n
		return
	}

	s.tail.next = n
	n.pre = s.tail
	s.tail = n
}

func (s *stack) pop() (error, int) {
	if s.tail == nil {
		return fmt.Errorf("empty stack"), 0
	}

	var result int

	// only one element left
	if s.head == s.tail {
		result = s.tail.data
		s.head, s.tail = nil, nil
		return nil, result
	}

	node := s.tail
	result = node.data
	s.tail = node.pre
	s.tail.next = nil
	node.pre = nil

	return nil, node.data
}

func main() {
	var s = stack{}
	s.push(1)
	s.push(2)
	s.push(3)

	fmt.Println(s.pop())
	fmt.Println(s.pop())
	fmt.Println(s.pop())

	fmt.Println(s.pop())
	fmt.Println(s.pop())
}
