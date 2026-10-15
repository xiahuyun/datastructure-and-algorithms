// https://leetcode.cn/problems/binary-tree-inorder-traversal/description/

package main

import (
	"fmt"
	"strings"
)

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// buildTree 按 LeetCode 的层序遍历格式构建二叉树，nil 表示空节点。
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
// inorderRecursive 递归解法：中序定义（左 -> 根 -> 右）的直译。
// 时间 O(n)，空间 O(高度) 递归栈，树退化为链表时最坏 O(n)。

	func inorderRecursive(root *TreeNode) []int {
		var res []int
		var walk func(*TreeNode)
		walk = func(node *TreeNode) {
			if node == nil {
				return
			}
			walk(node.Left)             // 左
			res = append(res, node.Val) // 根
			walk(node.Right)            // 右
		}
		walk(root)
		return res
	}
*/

var result = make([]int, 0)

func inorderRecursive(root *TreeNode) []int {
	inorder(root)
	return result
}

func inorder(root *TreeNode) {
	if root == nil {
		return
	}
	inorder(root.Left)
	result = append(result, root.Val)
	inorder(root.Right)
}

// inorderIterative 迭代解法：显式栈模拟递归。
// 核心不变式：栈里保存"左子树尚未处理完"的祖先链。
// 每轮三步：一路向左压栈到底 -> 弹出栈顶（其左子树已完成，轮到"根"）->
// 访问后转向右子树。循环条件 cur != nil || len(stack) > 0：
// 还有未进入的子树，或还有未回溯的祖先。
// 时间 O(n)，空间 O(高度)。
func inorderIterative(root *TreeNode) []int {
	var res []int
	var stack []*TreeNode
	cur := root
	for cur != nil || len(stack) > 0 {
		for cur != nil { // 一路向左，沿途全部压栈
			stack = append(stack, cur)
			cur = cur.Left
		}
		cur = stack[len(stack)-1] // 弹出：左子树已处理完，轮到"根"
		stack = stack[:len(stack)-1]
		res = append(res, cur.Val) // 访问
		cur = cur.Right            // 转向右子树
	}
	return res
}

// inorderMorris Morris 解法：O(1) 额外空间（既不用栈也不用递归）。
// 难点：没有栈，走完左子树后怎么"回来"？
// 答案：借中序前驱（左子树的最右节点）的空右指针当"回头路"——
// 进入左子树前，先把前驱的 Right 指向自己（临时线索）；
// 顺着线索恰好能回到自己。再次到达某节点且前驱指向自己，
// 说明左子树已遍历完：撤掉线索、访问、转向右子树。
// 每条边最多走 2 次（一次找前驱、一次回溯），时间仍 O(n)。
// 遍历过程中树被临时修改，结束后恢复原状（main_test.go 有专门断言）。
func inorderMorris(root *TreeNode) []int {
	var res []int
	cur := root
	for cur != nil {
		if cur.Left == nil {
			res = append(res, cur.Val) // 无左子树：直接访问
			cur = cur.Right
			continue
		}
		pred := cur.Left // 找中序前驱：左子树的最右节点
		for pred.Right != nil && pred.Right != cur {
			pred = pred.Right
		}
		if pred.Right == nil {
			pred.Right = cur // 第一次到达：建线索，进入左子树
			cur = cur.Left
		} else {
			pred.Right = nil           // 第二次到达（前驱指向自己）：撤线索
			res = append(res, cur.Val) // 左子树完成，访问
			cur = cur.Right
		}
	}
	return res
}

func main() {
	impls := []struct {
		name string
		fn   func(*TreeNode) []int
	}{
		{"递归", inorderRecursive},
		{"迭代", inorderIterative},
		{"Morris", inorderMorris},
	}
	demos := []struct {
		desc string
		vals []any
	}{
		{"官方示例1", []any{1, nil, 2, 3}},
		{"满二叉树(BST)", []any{4, 2, 6, 1, 3, 5, 7}},
		{"空树", nil},
	}
	fmt.Println("演示（三种实现依次为：递归 / 迭代 / Morris）")
	for _, d := range demos {
		tree := buildTree(d.vals)
		results := make([]string, len(impls))
		for i, impl := range impls {
			results[i] = fmt.Sprint(impl.fn(tree))
		}
		fmt.Printf("%s %v -> %s\n", d.desc, d.vals, strings.Join(results, " / "))
	}
}
