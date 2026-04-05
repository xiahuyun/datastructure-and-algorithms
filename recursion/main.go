// https://leetcode.cn/problems/climbing-stairs/

package main

import "fmt"

// 1 2 3 4 5 ...
// 1 2 3 5 8 ...
func climbStairs3(n int) int {
	if n == 1 {
		return 1
	} else if n == 2 {
		return 2
	}

	pre, now := 1, 2
	temp := -1
	for i := 2; i < n; i++ {
		temp = now
		now = pre + now
		pre = temp
	}

	return now
}

var stairs = make(map[int]int)

func climbStairs2(n int) int {
	if n == 1 {
		return 1
	} else if n == 2 {
		return 2
	}

	if result, ok := stairs[n]; ok {
		return result
	}

	stairs[n] = climbStairs2(n-1) + climbStairs2(n-2)
	return stairs[n]
}

func climbStairs1(n int) int {
	if n == 1 {
		return 1
	} else if n == 2 {
		return 2
	}

	return climbStairs1(n-1) + climbStairs1(n-2)
}

func main() {
	fmt.Println(climbStairs1(10))
	fmt.Println(climbStairs2(10))
	fmt.Println(climbStairs3(10))
}
