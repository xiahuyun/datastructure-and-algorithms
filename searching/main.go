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

func BinarySearchRepeat(arr []int, value int) (bool, int) {
	low, high := 0, len(arr)-1
	for low <= high {
		mid := low + (high-low)/2
		if arr[mid] > value {
			high = mid - 1
		} else if arr[mid] < value {
			low = mid + 1
		} else {
			if mid == len(arr)-1 || arr[mid+1] != value {
				return true, mid
			}
			low = mid + 1
		}
	}
	return false, -1
}

func BinarySearchFirst(arr []int, value int) (bool, int) {
	low, high := 0, len(arr)-1
	for low <= high {
		mid := low + (high-low)/2
		if arr[mid] > value {
			if mid == 0 || arr[mid-1] < value {
				return true, mid
			}
			high = mid - 1
		} else if arr[mid] < value {
			if mid == len(arr)-1 {
				return false, -1
			} else if arr[mid+1] > value {
				return true, mid + 1
			}
			low = mid + 1
		} else {
			if mid == 0 || arr[mid-1] != value {
				return true, mid
			}
			high = mid - 1
		}
	}
	return false, -1
}

func main() {
	var arr = []int{8, 11, 19, 23, 27, 33, 45, 55, 67, 98}
	fmt.Println(BinarySearchFirst(arr, 35))
}
