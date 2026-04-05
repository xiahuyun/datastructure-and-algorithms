package main

import "fmt"

func InsertionSort(arr []int) []int {
	n := len(arr)
	for i := 1; i < n; i++ {
		value := arr[i]
		j := i - 1
		for ; j >= 0; j-- {
			if value < arr[j] {
				arr[j+1] = arr[j]
			} else {
				break
			}
		}
		arr[j+1] = value
	}
	return arr
}

func main() {
	var arr = []int{3, 6, 2, 1, 5}
	fmt.Println(InsertionSort(arr))
}
