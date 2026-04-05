package main

import "fmt"

var stack = []int{}

func push(data int) {
	stack = append(stack, data)
}

func pop() (error, int) {
	length := len(stack)
	if length <= 0 {
		return fmt.Errorf("empty stack"), 0
	}

	result := stack[length-1]
	fmt.Println(stack)
	stack = stack[:length-1]
	fmt.Println(stack)
	return nil, result
}

func main() {
	push(1)
	push(2)
	push(3)

	fmt.Println(pop())
	fmt.Println(pop())
	fmt.Println(pop())
	fmt.Println(pop())
}
