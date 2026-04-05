package main

import "fmt"

func maxIndex(arr []int, left, right int) int {
	n := len(arr) - 1
	if left > n {
		return -1
	}

	if right > n {
		return left
	}

	if arr[left] > arr[right] {
		return left
	}

	return right
}

func maxHeap(arr []int, lagest int) {
	for {
		left := 2 * lagest
		right := 2*lagest + 1
		index := maxIndex(arr, left, right)
		if index == -1 {
			break
		}

		if arr[lagest] < arr[index] {
			arr[lagest], arr[index] = arr[index], arr[lagest]
			lagest = index
		} else {
			break
		}
	}
}

func heapSort(arr []int) {
	n := len(arr) - 1
	for i := n / 2; i >= 1; i-- {
		maxHeap(arr, i)
	}

	i := 0
	for j := n; j >= 1; j-- {
		arr[1], arr[j] = arr[j], arr[1]
		maxHeap(arr[:n-i], 1)
		i++
	}
}

func main() {
	var arr = []int{0, 12, 11, 13, 5, 6, 7}
	fmt.Println(arr)

	heapSort(arr)

	fmt.Println(arr)
}
