package main

import "fmt"

func sortColors(nums []int) {
	h, t := 0, len(nums)-1
	for i := range nums {
		if nums[i] != 0 {
			h = i
			break
		}
	}
	n := len(nums)
	for i := range nums {
		if nums[n-i-1] != 2 {
			t = n - i - 1
			break
		}
	}

	for c := h; c <= t; {
		if nums[c] == 2 {
			nums[c], nums[t] = nums[t], nums[c]
			t--
		} else if nums[c] == 1 {
			c++
		} else if nums[c] == 0 && c == h {
			c++
			h++
		} else if nums[c] == 0 && c != h {
			nums[c], nums[h] = nums[h], nums[c]
			c++
			h++
		}
	}
}

func main() {
	nums := []int{2, 0, 2, 1, 1, 0}
	sortColors(nums)
	fmt.Println(nums) // 输出: [0, 0, 1, 1, 2, 2]
}
