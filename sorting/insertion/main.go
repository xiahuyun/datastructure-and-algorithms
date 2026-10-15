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
	var arr = []int{4, 5, 6, 3, 2, 1}
	fmt.Println(InsertionSortV2(arr))
}

func InsertionSortV2(arr []int) []int {
	for i := 0; i < len(arr); i++ {
		for j := i; j > 0; j-- {
			if arr[j] > arr[j-1] {
				break
			}
			arr[j], arr[j-1] = arr[j-1], arr[j]
		}
	}
	return arr
}
