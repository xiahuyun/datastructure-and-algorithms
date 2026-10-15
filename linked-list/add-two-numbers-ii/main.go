// https://leetcode.cn/problems/add-two-numbers-ii/description/

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

// buildList 根据切片构建链表，切片顺序即链表顺序（最高位在前），方便构造测试数据。
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

// numString 把正序存储的数字切片拼接成可读字符串，如 [7,2,4,3] -> "7243"。
func numString(vals []int) string {
	var b strings.Builder
	for _, v := range vals {
		b.WriteString(strconv.Itoa(v))
	}
	return b.String()
}

// reverseList 迭代反转链表：三指针 pre/cur/next，逐个调转指针方向。
func reverseList(head *ListNode) *ListNode {
	var pre *ListNode
	for cur := head; cur != nil; {
		next := cur.Next
		cur.Next = pre
		pre = cur
		cur = next
	}
	return pre
}

// addTwoNumbersReverse 方案一：反转链表（会修改输入）。
// 正序存储意味着头节点是最高位，顺着链表加是从高位加到低位，进位方向反了。
// 先把两条输入反转成"个位在前"（第 2 题的形态），做完竖式加法后，
// 再把结果反转回正序。额外空间 O(1)，代价是破坏了输入链表。
func addTwoNumbersReverse(l1, l2 *ListNode) *ListNode {
	l1 = reverseList(l1)
	l2 = reverseList(l2)

	// 第 2 题式竖式加法：头节点就是个位，两指针同步前移逐位相加；
	// 循环条件必须带上 carry != 0，让最高位的进位能多出一个节点。
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
	return reverseList(dummy.Next)
}

// addTwoNumbersStack 方案二：栈（不修改输入）。
// 栈（LIFO）天然倒转访问顺序：顺序入栈后，弹出的顺序恰好是低位到高位。
// 两条链表的值分别入栈，循环弹栈逐位相加——弹空的一侧按 0 参与运算，
// 自动对齐长度差；每算出一位就把新节点头插到结果头部，结果天然正序，
// 最高位的进位也走同一条头插逻辑。额外空间 O(m+n)。
func addTwoNumbersStack(l1, l2 *ListNode) *ListNode {
	var s1, s2 []int
	for ; l1 != nil; l1 = l1.Next {
		s1 = append(s1, l1.Val)
	}
	for ; l2 != nil; l2 = l2.Next {
		s2 = append(s2, l2.Val)
	}

	var head *ListNode
	carry := 0
	for len(s1) > 0 || len(s2) > 0 || carry != 0 {
		sum := carry
		if len(s1) > 0 {
			sum += s1[len(s1)-1]
			s1 = s1[:len(s1)-1]
		}
		if len(s2) > 0 {
			sum += s2[len(s2)-1]
			s2 = s2[:len(s2)-1]
		}
		carry = sum / 10
		head = &ListNode{Val: sum % 10, Next: head} // 头插：越晚建成的节点是越高位
	}
	return head
}

// addTwoNumbersRecursive 方案三：递归（不修改输入）。
// 递归先一路走到链表尾，"归"的阶段恰好从低位（尾部）回到高位（头部），
// 进位作为返回值向高层传递——递归本身就是"倒着的遍历"。
// 先比较两条链表长度并保证 l1 不短于 l2，长度差 diff 就是 l1 需要
// 独自前移的步数（等价于给短链表补前导零）。题目约束链表最长 100，
// 递归深度安全。额外空间 O(n) 调用栈。
func addTwoNumbersRecursive(l1, l2 *ListNode) *ListNode {
	n1, n2 := length(l1), length(l2)
	if n1 < n2 {
		l1, l2 = l2, l1 // 保证 l1 不短于 l2
		n1, n2 = n2, n1 // 长度必须一起交换，diff 才是正确的正数
	}
	head, carry := addRecursive(l1, l2, n1-n2)
	if carry != 0 {
		return &ListNode{Val: carry, Next: head}
	}
	return head
}

// length 计算链表长度。
func length(head *ListNode) int {
	n := 0
	for ; head != nil; head = head.Next {
		n++
	}
	return n
}

// addRecursive 对两条按长度差对齐的链表逐位求和。
// diff > 0 表示 l1 还有富余的高位节点，本位 l2 按 0（补前导零）处理；
// diff == 0 后两条链表同步前移。
// 返回值：（已构建好的结果头节点，向上一位传递的进位）。
func addRecursive(l1, l2 *ListNode, diff int) (*ListNode, int) {
	v1, v2 := 0, 0
	if l1 != nil {
		v1 = l1.Val
		l1 = l1.Next
	}
	if diff > 0 {
		diff--
	} else if l2 != nil {
		v2 = l2.Val
		l2 = l2.Next
	}

	var next *ListNode
	carry := 0
	if l1 != nil || l2 != nil {
		next, carry = addRecursive(l1, l2, diff)
	}

	sum := v1 + v2 + carry
	return &ListNode{Val: sum % 10, Next: next}, sum / 10
}

func main() {
	impls := []struct {
		name string
		fn   func(l1, l2 *ListNode) *ListNode
	}{
		{"反转链表", addTwoNumbersReverse},
		{"栈", addTwoNumbersStack},
		{"递归", addTwoNumbersRecursive},
	}
	demos := []struct {
		l1, l2 []int
	}{
		{[]int{7, 2, 4, 3}, []int{5, 6, 4}}, // 官方示例 1：7243 + 564
		{[]int{2, 4, 3}, []int{5, 6, 4}},    // 官方示例 2：243 + 564
		{[]int{9, 9}, []int{1}},             // 最高位进位：99 + 1
	}
	fmt.Println("演示（三种实现依次为：反转链表 / 栈 / 递归）")
	for _, d := range demos {
		results := make([]string, len(impls))
		for i, impl := range impls {
			results[i] = numString(listToSlice(impl.fn(buildList(d.l1), buildList(d.l2))))
		}
		fmt.Printf("%s + %s = %s\n", numString(d.l1), numString(d.l2), strings.Join(results, " / "))
	}
}
