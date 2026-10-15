// https://leetcode.cn/problems/add-two-numbers/description/?envType=study-plan-v2&envId=top-100-liked

package main

import (
	"fmt"
	"strconv"
	"strings"
)

type ListNode struct {
	Val  int
	Next *ListNode
}

// buildList 根据切片构建链表，切片顺序即链表顺序（个位在前），方便构造测试数据。
func buildList(vals []int) *ListNode {
	dummy := &ListNode{}
	cur := dummy
	for _, v := range vals {
		cur.Next = &ListNode{Val: v}
		cur = cur.Next
	}
	return dummy.Next
}

// listToSlice 将链表转回切片，便于断言和打印。
func listToSlice(head *ListNode) []int {
	var vals []int
	for ; head != nil; head = head.Next {
		vals = append(vals, head.Val)
	}
	return vals
}

// numString 把逆序存储的数字切片还原成可读的十进制字符串，如 [2,4,3] -> "342"。
func numString(vals []int) string {
	var b strings.Builder
	for i := len(vals) - 1; i >= 0; i-- {
		b.WriteString(strconv.Itoa(vals[i]))
	}
	return b.String()
}

// addTwoNumbers 用链表模拟竖式加法。
// 逆序存储意味着头节点就是个位，两指针同步前移，逐位相加：
// 每轮 sum = carry + l1.Val + l2.Val，本位写 sum%10，进位更新为 sum/10。
// 循环条件是三者的"或"：任一链表未走完、或还有进位，就要继续——
// 最高位的进位会凭空多出一个节点（如 99 + 1 = 100）。
// 哑节点 dummy 统一了结果头节点的构建，免去对第一个节点的特殊处理。
// 时间复杂度 O(max(m,n))，额外空间 O(1)（结果链表不计）。
func addTwoNumbers(l1, l2 *ListNode) *ListNode {
	dummy := &ListNode{}
	cur := dummy
	carry := 0
	for l1 != nil || l2 != nil || carry != 0 {
		sum := carry
		if l1 != nil {
			sum += l1.Val
			l1 = l1.Next
		}
		if l2 != nil {
			sum += l2.Val
			l2 = l2.Next
		}
		carry = sum / 10
		cur.Next = &ListNode{Val: sum % 10}
		cur = cur.Next
	}
	return dummy.Next
}

func main() {
	demos := []struct {
		l1, l2 []int
	}{
		{[]int{2, 4, 3}, []int{5, 6, 4}},
		{[]int{0}, []int{0}},
		{[]int{9, 9}, []int{1}},
		{[]int{9, 9, 9, 9, 9, 9, 9}, []int{9, 9, 9, 9}},
	}
	for _, d := range demos {
		sum := listToSlice(addTwoNumbers(buildList(d.l1), buildList(d.l2)))
		fmt.Printf("%s + %s = %s\n", numString(d.l1), numString(d.l2), numString(sum))
	}
}
