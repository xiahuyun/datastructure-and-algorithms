package main

import "fmt"

type queue struct {
	head, tail *node
}

type node struct {
	data int
	next *node
	pre  *node
}

func (q *queue) enqueue(data int) {
	n := &node{data: data}
	if q.tail == nil && q.head == q.tail {
		q.head, q.tail = n, n
		return
	}

	q.tail.next = n
	n.pre = q.tail
	q.tail = n
}

func (q *queue) dequeue() (error, int) {
	if q.head == nil && q.head == q.tail {
		return fmt.Errorf("empty queue"), 0
	}

	var result int

	// only one element left
	if q.head == q.tail {
		result = q.head.data
		q.head, q.tail = nil, nil
		return nil, result
	}

	result = q.head.data
	newhead := q.head.next
	newhead.pre = nil
	q.head.next = nil
	q.head = newhead

	return nil, result
}

func main() {
	q := &queue{}
	q.enqueue(1)
	q.enqueue(2)
	q.enqueue(3)

	fmt.Println(q.dequeue())
	fmt.Println(q.dequeue())
	fmt.Println(q.dequeue())

	fmt.Println(q.dequeue())
	fmt.Println(q.dequeue())
}
