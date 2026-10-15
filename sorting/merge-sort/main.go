package main

import "fmt"

func merge(left []int, right []int) []int {
	n, m := len(left), len(right)
	var result = make([]int, n+m)

	i, j := 0, 0
	k := 0
	for i < n && j < m {
		if left[i] > right[j] {
			result[k] = right[j]
			j++
			k++
		} else {
			result[k] = left[i]
			i++
			k++
		}
	}

	if i < n {
		for i < n {
			result[k] = left[i]
			i++
			k++
		}
	} else if j < m {
		for j < m {
			result[k] = right[j]
			k++
			j++
		}
	}

	return result
}

func MergeSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}

	mid := len(arr) / 2

	left := MergeSort(arr[:mid])
	right := MergeSort(arr[mid:])

	return merge(left, right)
}

func mergeOverwrite(arr []int, left, mid, right int) {
	n1 := mid - left + 1
	n2 := right - mid

	leftArr := make([]int, n1)
	rightArr := make([]int, n2)

	for i := 0; i < n1; i++ {
		leftArr[i] = arr[left+i]
	}
	for j := 0; j < n2; j++ {
		rightArr[j] = arr[mid+1+j]
	}

	i, j, k := 0, 0, left
	for i < n1 && j < n2 {
		if leftArr[i] <= rightArr[j] {
			arr[k] = leftArr[i]
			i++
		} else {
			arr[k] = rightArr[j]
			j++
		}
		k++
	}

	for i < n1 {
		arr[k] = leftArr[i]
		i++
		k++
	}

	for j < n2 {
		arr[k] = rightArr[j]
		j++
		k++
	}
}

func MergeSortOverwrite(arr []int, left, right int) {
	if left < right {
		mid := left + (right-left)/2

		MergeSortOverwrite(arr, left, mid)
		MergeSortOverwrite(arr, mid+1, right)

		mergeOverwrite(arr, left, mid, right)
	}
}

func main() {
	var arr = []int{11, 8, 3, 9, 7, 1, 2, 5}
	result := MergeSortV2(arr, 0, len(arr)-1)
	fmt.Println(result, arr)

	//MergeSortOverwrite(arr, 0, len(arr)-1)
	//fmt.Println(arr)
}

// 归并排序
// mergeSort(p, q) = merge(mergeSort(p, mid), mergeSort(mid+1, q))
// if p >= q, return

func MergeSortV2(arr []int, left, right int) []int {
	if left >= right {
		return arr[left : right+1]
	}

	mid := left + (right-left)/2
	leftSort := MergeSortV2(arr, left, mid)
	rightSort := MergeSortV2(arr, mid+1, right)
	return mergeV2(leftSort, rightSort)
}

func mergeV2(left, right []int) []int {
	n, m := len(left), len(right)
	var result = make([]int, n+m)

	i, j := 0, 0
	for i < n && j < m {
		if right[j] < left[i] {
			result[i+j] = right[j]
			j++
		} else {
			result[i+j] = left[i]
			i++
		}
	}

	if i < n {
		for ; i < n; i++ {
			result[i+m] = left[i]
		}
	} else if j < m {
		for ; j < m; j++ {
			result[j+n] = right[j]
		}
	}

	return result
}
