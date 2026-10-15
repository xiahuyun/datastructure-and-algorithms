// https://leetcode.cn/problems/symmetric-tree/description/

package main

import "fmt"

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// buildTree 按 LeetCode 的层序遍历格式构建二叉树，nil 表示空节点。
// 例如 [1,2,2,nil,3,nil,3] 表示根为 1，两个孩子都是 2；
// 左孩子 2 的右孩子是 3（左孩子为空），右孩子 2 的右孩子是 3（左孩子为空）。
// 空节点不入队，因此它不再占用后续元素——与 LeetCode 的序列化规则一致。
func buildTree(vals []any) *TreeNode {
	if len(vals) == 0 {
		return nil
	}
	root := &TreeNode{Val: vals[0].(int)}
	queue := []*TreeNode{root}
	i := 1
	for len(queue) > 0 && i < len(vals) {
		node := queue[0]
		queue = queue[1:]
		if v, ok := vals[i].(int); ok {
			node.Left = &TreeNode{Val: v}
			queue = append(queue, node.Left)
		}
		i++
		if i < len(vals) {
			if v, ok := vals[i].(int); ok {
				node.Right = &TreeNode{Val: v}
				queue = append(queue, node.Right)
			}
			i++
		}
	}
	return root
}

/*
// isSymmetricRecursive 递归解法。
// 整棵树对称 ⟺ 根的左右子树互为镜像。镜像的递归定义：
//  1. 两棵空树互为镜像；
//  2. 都非空时值相等，且"左对右、右对左"交叉递归——
//     A 的左子树与 B 的右子树互为镜像，A 的右子树与 B 的左子树互为镜像。
//
// 交叉比较是精髓：沿对称轴翻转后，左就变成了右。
// 顺手写成"左对左、右对右"的同向比较，就变成了 100. 相同的树。
// 时间 O(n)，空间 O(高度) 递归栈，树退化为链表时最坏 O(n)。
func isSymmetricRecursive(root *TreeNode) bool {
	return isMirror(root, root)
}

func isMirror(l, r *TreeNode) bool {
	if l == nil || r == nil {
		return l == r // 两棵都为空才互为镜像
	}
	return l.Val == r.Val && isMirror(l.Left, r.Right) && isMirror(l.Right, r.Left)
}
*/

func isSymmetricRecursive(root *TreeNode) bool {
	if root == nil {
		return true // 空树按定义是对称的，也避免下方解引用 nil
	}
	if root.Left == nil && root.Right == nil {
		return true
	}
	return isMirror(root.Left, root.Right)
}

func isMirror(left, right *TreeNode) bool {
	if left == nil && right == nil {
		return true
	} else if left == nil && right != nil {
		return false
	} else if left != nil && right == nil {
		return false
	} else if left.Val != right.Val {
		return false
	}
	return isMirror(left.Left, right.Right) && isMirror(left.Right, right.Left)
}

// isSymmetricIterative 迭代解法：队列 + 成对比较。
// 队列中始终成对存放"待比较的节点对"，每次出队两个节点：
// 恰一边为空或值不等，立即返回 false；否则把
// (l.左, r.右) 和 (l.右, r.左) 两对继续入队。
// 子问题拆分与递归完全一致，只是用显式队列替代了调用栈。
// 时间 O(n)，空间 O(n)。
func isSymmetricIterative(root *TreeNode) bool {
	queue := []*TreeNode{root, root}
	for len(queue) > 0 {
		l, r := queue[0], queue[1]
		queue = queue[2:]
		if l == nil && r == nil {
			continue
		}
		if l == nil || r == nil || l.Val != r.Val {
			return false
		}
		queue = append(queue, l.Left, r.Right, l.Right, r.Left)
	}
	return true
}

func main() {
	demos := []struct {
		desc string
		vals []any
	}{
		{"示例1: 对称", []any{1, 2, 2, 3, 4, 4, 3}},
		{"示例2: 结构不镜像", []any{1, 2, 2, nil, 3, nil, 3}},
		{"结构非满但恰好镜像", []any{1, 2, 2, 3, nil, nil, 3}},
		{"空树", nil},
	}
	for _, d := range demos {
		tree := buildTree(d.vals)
		fmt.Printf("%-14s %v -> 递归: %v, 迭代: %v\n", d.desc, d.vals, isSymmetricRecursive(tree), isSymmetricIterative(tree))
	}
}
