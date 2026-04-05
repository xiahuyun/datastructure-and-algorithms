package main

import "fmt"

func BinarySearch(arr []int, value int) bool {
	low, high := 0, len(arr)
	for low <= high {
		mid := low + (high-low)/2
		if arr[mid] > value {
			high = mid - 1
		} else if arr[mid] < value {
			low = mid + 1
		} else {
			return true
		}
	}

	return false
}

func main() {
	var arr = []int{8, 11, 19, 23, 27, 33, 45, 55, 67, 98}
	fmt.Println(BinarySearch(arr, 20))
}
