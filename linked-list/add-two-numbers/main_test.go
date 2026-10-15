package main

import (
	"slices"
	"testing"
)

func TestAddTwoNumbers(t *testing.T) {
	tests := []struct {
		name string
		l1   []int
		l2   []int
		want []int
	}{
		{"示例1: 342+465=807", []int{2, 4, 3}, []int{5, 6, 4}, []int{7, 0, 8}},
		{"示例2: 0+0=0", []int{0}, []int{0}, []int{0}},
		{"示例3: 9999999+9999=10009998", []int{9, 9, 9, 9, 9, 9, 9}, []int{9, 9, 9, 9}, []int{8, 9, 9, 9, 0, 0, 0, 1}},
		{"长度不等: 99+1=100", []int{9, 9}, []int{1}, []int{0, 0, 1}},
		{"最高位进位: 999+1=1000", []int{9, 9, 9}, []int{1}, []int{0, 0, 0, 1}},
		{"无进位: 3+4=7", []int{3}, []int{4}, []int{7}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := listToSlice(addTwoNumbers(buildList(tt.l1), buildList(tt.l2)))
			if !slices.Equal(got, tt.want) {
				t.Fatalf("addTwoNumbers(%v, %v) = %v, want %v", tt.l1, tt.l2, got, tt.want)
			}
		})
	}
}

// TestMaxConstraint 覆盖题目边界：链表最长 100 个节点（数字可达 10^100）。
// 两个 100 位的全 9 数相加触发全链进位，且结果多出一个节点。
// 这也说明不能把链表转成 int64 相加——必然溢出，本题本质是大数加法。
func TestMaxConstraint(t *testing.T) {
	nines := make([]int, 100)
	for i := range nines {
		nines[i] = 9
	}

	// (10^100 - 1) * 2 = 2*10^100 - 2，即 1 后跟 99 个 9 再一个 8，
	// 逆序存储就是：8 + 99 个 9 + 1，共 101 个节点。
	want := make([]int, 101)
	want[0] = 8
	for i := 1; i <= 99; i++ {
		want[i] = 9
	}
	want[100] = 1

	got := listToSlice(addTwoNumbers(buildList(nines), buildList(nines)))
	if !slices.Equal(got, want) {
		t.Fatalf("max-length addition mismatch:\ngot  %v\nwant %v", got, want)
	}
}
