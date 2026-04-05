package main

import "fmt"

func SelectSort(arr []int) []int {
	n := len(arr)
	min := -1
	for i := 0; i < n; i++ {
		min = i
		for j := i + 1; j < n; j++ {
			if arr[j] < arr[min] {
				min = j
			}
		}
		arr[i], arr[min] = arr[min], arr[i]
	}
	return arr
}

func main() {
	var arr = []int{3, 6, 2, 1, 5}
	fmt.Println(SelectSort(arr))
}
