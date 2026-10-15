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
	var arr = []int{4, 5, 6, 3, 2, 1}
	fmt.Println(SelectSortV3(arr))
}

func SelectSortV2(arr []int) []int {
	iter := 0
	for i := 0; i < len(arr)-iter; {
		max := i
		for j := 0; j < len(arr)-iter; j++ {
			if arr[j] > arr[max] {
				max = j
			}
		}
		arr[len(arr)-1-iter], arr[max] = arr[max], arr[len(arr)-1-iter]
		iter++
		i = 0
	}
	return arr
}

func SelectSortV3(arr []int) []int {
	for i := 0; i < len(arr)-1; i++ {
		max := 0
		for j := 0; j < len(arr)-i; j++ {
			if arr[j] > arr[max] {
				max = j
			}
		}
		arr[len(arr)-1-i], arr[max] = arr[max], arr[len(arr)-1-i]
	}
	return arr
}
