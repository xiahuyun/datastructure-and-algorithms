// https://leetcode.cn/problems/reverse-linked-list/description/?envType=study-plan-v2&envId=top-100-liked

package main

import "fmt"

type ListNode struct {
	Val  int
	Next *ListNode
}

// buildList 根据切片构建链表，方便构造测试数据。
func buildList(vals []int) *ListNode {
	dummy := &ListNode{}
	cur := dummy
	for _, v := range vals {
		cur.Next = &ListNode{Val: v}
		cur = cur.Next
	}
	return dummy.Next
}

func printList(head *ListNode) {
	for cur := head; cur != nil; cur = cur.Next {
		fmt.Printf("%d ", cur.Val)
	}
	fmt.Println()
}

// reverseListV1 迭代（双指针）反转链表。
func reverseListV1(head *ListNode) *ListNode {
	var pre *ListNode = nil
	for cur := head; cur != nil; {
		next := cur.Next
		cur.Next = pre
		pre = cur
		cur = next
	}
	return pre
}

// reverseListV2 递归反转链表。
func reverseListV2(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}
	newHead := reverseListV2(head.Next)
	head.Next.Next = head
	head.Next = nil
	return newHead
}

func main() {
	nums := []int{1, 2, 3, 4, 5}

	fmt.Printf("原始链表:    ")
	printList(buildList(nums))

	fmt.Printf("迭代反转:    ")
	printList(reverseListV1(buildList(nums)))

	fmt.Printf("递归反转:    ")
	printList(reverseListV2(buildList(nums)))

	fmt.Printf("递归反转V3:    ")
	printList(reverseListV3(buildList(nums)))
}

// 递归反转
func reverseListV3(head *ListNode) *ListNode {
	if head == nil {
		return nil
	}

	if head.Next == nil {
		return head
	}
	t := head.Next
	head.Next.Next = head
	head.Next = nil
	return reverseListV3(t)
}
