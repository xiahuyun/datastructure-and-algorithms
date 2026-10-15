package main

import (
	"slices"
	"testing"
)

// impls 汇总三种实现，保证每种实现都跑完全部同一组用例。
var impls = []struct {
	name string
	fn   func(l1, l2 *ListNode) *ListNode
}{
	{"反转链表", addTwoNumbersReverse},
	{"栈", addTwoNumbersStack},
	{"递归", addTwoNumbersRecursive},
}

func TestAddTwoNumbersII(t *testing.T) {
	tests := []struct {
		name string
		l1   []int
		l2   []int
		want []int
	}{
		{"示例1: 7243+564=7807", []int{7, 2, 4, 3}, []int{5, 6, 4}, []int{7, 8, 0, 7}},
		{"示例2: 243+564=807", []int{2, 4, 3}, []int{5, 6, 4}, []int{8, 0, 7}},
		{"长度不等: 5+998=1003", []int{5}, []int{9, 9, 8}, []int{1, 0, 0, 3}},
		{"最高位进位: 99+1=100", []int{9, 9}, []int{1}, []int{1, 0, 0}},
		{"零: 0+0=0", []int{0}, []int{0}, []int{0}},
		{"无进位: 3+4=7", []int{3}, []int{4}, []int{7}},
	}

	for _, impl := range impls {
		for _, tt := range tests {
			t.Run(impl.name+"/"+tt.name, func(t *testing.T) {
				got := listToSlice(impl.fn(buildList(tt.l1), buildList(tt.l2)))
				if !slices.Equal(got, tt.want) {
					t.Fatalf("%s(%v, %v) = %v, want %v", impl.name, tt.l1, tt.l2, got, tt.want)
				}
			})
		}
	}
}

// TestMaxConstraint 覆盖题目边界：链表最长 100 个节点（数字可达 10^100）。
// 两个 100 位的全 9 数相加触发全链进位，且结果多出一个节点。
// 这也再次说明本题是大数加法，不能把链表转成 int64 计算。
func TestMaxConstraint(t *testing.T) {
	nines := make([]int, 100)
	for i := range nines {
		nines[i] = 9
	}

	// (10^100 - 1) * 2 = 2*10^100 - 2，即 1 后跟 99 个 9 再一个 8，
	// 正序存储共 101 个节点。
	want := make([]int, 101)
	want[0] = 1
	for i := 1; i <= 99; i++ {
		want[i] = 9
	}
	want[100] = 8

	for _, impl := range impls {
		t.Run(impl.name, func(t *testing.T) {
			got := listToSlice(impl.fn(buildList(nines), buildList(nines)))
			if !slices.Equal(got, want) {
				t.Fatalf("%s: max-length addition mismatch:\ngot  %v\nwant %v", impl.name, got, want)
			}
		})
	}
}
