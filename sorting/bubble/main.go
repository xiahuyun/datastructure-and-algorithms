package main

import "fmt"

func BubbleSort(arr []int) []int {
	n := len(arr)
	for i := 0; i < n; i++ {
		swapped := false
		for j := 0; j < n-i-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
				swapped = true
			}
		}
		if !swapped {
			break
		}
	}
	return arr
}

func main() {
	var arr = []int{4, 5, 6, 3, 2, 1}
	fmt.Println(BubbleSort(arr))
}

func BubbleSortV2(arr []int) []int {
	i := 0
	iter := 0
	for i < len(arr)-iter {
		for j := 0; j < len(arr)-iter; j++ {
			if j+1 < len(arr)-iter && arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
		iter++
		i = 0
	}
	return arr
}
