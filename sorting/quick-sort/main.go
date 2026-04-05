package main

import (
	"fmt"
	"sort"
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
	QuickSort(arr, 0, len(arr)-1)

	fmt.Println(arr)

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
}
