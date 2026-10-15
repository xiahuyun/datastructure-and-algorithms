package main

import (
	"fmt"
)

func partition(arr []int, left, right int) int {
	pivot := arr[left]
	isRight := true
	for left < right {
		if isRight {
			if arr[right] > pivot {
				right--
			} else {
				arr[left] = arr[right]
				left++
				isRight = false
			}
		} else {
			if arr[left] > pivot {
				arr[right] = arr[left]
				right--
				isRight = true
			} else {
				left++
			}
		}
	}

	if left == right {
		arr[left] = pivot
	}

	return left
}

func QuickSort(arr []int, left, right int) {
	if left < right {
		k := partition(arr, left, right)

		QuickSort(arr, left, k-1)
		QuickSort(arr, k+1, right)
	}
}

type Person struct {
	Name string
	Age  int
}

type ByAge []Person

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }

func main() {
	var arr = []int{11, 8, 3, 9, 7, 1, 2, 5}
	QuickSortV2(arr, 0, len(arr)-1)

	fmt.Println(arr)

	/*
		var arr2 = []int{1, 2, 3, 5, 8, 10}
		QuickSort(arr, 0, len(arr2)-1)

		people := []Person{
			{"Bob", 31},
			{"John", 42},
			{"Michael", 17},
			{"Jenny", 26},
			{"Jenny", 26},
			{"Jenny", 26},
			{"Jenny", 26},
			{"Jenny", 26},
			{"Jenny", 26},
			{"Jenny", 26},
			{"Jenny", 26},
			{"Jenny", 26},
			{"Jenny", 26},
			{"Jenny", 26},
		}

		sort.Sort(ByAge(people))
	*/
}

// QuickSortV2 快速排序v2
// QuickSortV2(p, q) = QuickSortV2(p, k-1) + QuickSortV2(k+1, q)
// if p >= q, return
func QuickSortV2(arr []int, left, right int) {
	if left >= right {
		return
	}

	pivot := partitionV2(arr, left, right)
	QuickSortV2(arr, left, pivot-1)
	QuickSortV2(arr, pivot+1, right)
}

// var arr = []int{11, 8, 3, 9, 7, 1, 2, 5}
// lessTemp = []int{3, 1, 2}
// greaterTemp = []int{9, 7, 8, 11}
func partitionV2(arr []int, left, right int) int {
	var lessTemp = make([]int, 0)
	var greaterTemp = make([]int, 0)
	pivot := arr[right]
	for i := left; i < right+1; i++ {
		if i == right {
			continue
		}
		if arr[i] <= pivot {
			lessTemp = append(lessTemp, arr[i])
		} else {
			greaterTemp = append(greaterTemp, arr[i])
		}
	}

	copy(arr[left:left+len(lessTemp)], lessTemp)
	arr[left+len(lessTemp)] = pivot
	for i := 0; i < len(greaterTemp); i++ {
		arr[left+i+len(lessTemp)+1] = greaterTemp[i]
	}
	return left + len(lessTemp)
}
