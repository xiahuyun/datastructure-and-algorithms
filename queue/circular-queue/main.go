package main

import "fmt"

type queue struct {
	arrs [7]int
	head int
	tail int
}

func (q *queue) enqueue(data int) error {
	// queue is full
	if (q.tail+1)%len(q.arrs) == q.head {
		return fmt.Errorf("queue is full, can not enqueue the data [%d]", data)
	}

	q.arrs[q.tail] = data
	if (q.tail + 1) == len(q.arrs) {
		q.tail = 0
	} else {
		q.tail++
	}

	return nil
}

func (q *queue) dequeue() (error, int) {
	// queue is empty
	if q.head == q.tail {
		return fmt.Errorf("empty queue"), 0
	}

	result := q.arrs[q.head]
	if (q.head + 1) == len(q.arrs) {
		q.head = 0
	} else {
		q.head++
	}

	return nil, result
}

func main() {
	q := &queue{}

	fmt.Println(q.enqueue(1))
	fmt.Println(q.enqueue(2))
	fmt.Println(q.enqueue(3))
	fmt.Println(q.enqueue(4))
	fmt.Println(q.enqueue(5))
	fmt.Println(q.enqueue(6))
	fmt.Println(q.enqueue(7))
	fmt.Println(q.enqueue(8))

	fmt.Println(q.dequeue())
	fmt.Println(q.dequeue())
	fmt.Println(q.dequeue())
	fmt.Println(q.dequeue())
	fmt.Println(q.dequeue())
	fmt.Println(q.dequeue())
	fmt.Println(q.dequeue())
	fmt.Println(q.dequeue())

	fmt.Println(q.head, q.tail)

	fmt.Println(q.enqueue(1))
	fmt.Println(q.enqueue(2))
	fmt.Println(q.enqueue(3))
	fmt.Println(q.enqueue(4))

	fmt.Println(q.head, q.tail, q.arrs)
}
