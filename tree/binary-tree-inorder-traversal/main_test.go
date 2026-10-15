package main

import (
	"slices"
	"testing"
)

// impls 汇总三种实现，保证每种实现都跑完全部同一组用例。
var impls = []struct {
	name string
	fn   func(*TreeNode) []int
}{
	{"递归", inorderRecursive},
	{"迭代", inorderIterative},
	{"Morris", inorderMorris},
}

func TestInorderTraversal(t *testing.T) {
	tests := []struct {
		name string
		vals []any
		want []int
	}{
		{"示例1: 右子树带左孩子", []any{1, nil, 2, 3}, []int{1, 3, 2}},
		{"示例2: 空树", nil, []int{}},
		{"示例3: 单节点", []any{1}, []int{1}},
		{"示例4: 只有左孩子", []any{1, 2}, []int{2, 1}},
		{"示例5: 只有右孩子", []any{1, nil, 2}, []int{1, 2}},
		{"满二叉树(BST): 中序升序", []any{4, 2, 6, 1, 3, 5, 7}, []int{1, 2, 3, 4, 5, 6, 7}},
		{"左斜链: 全是左孩子", []any{4, 3, nil, 2, nil, 1}, []int{1, 2, 3, 4}},
	}

	for _, impl := range impls {
		for _, tt := range tests {
			t.Run(impl.name+"/"+tt.name, func(t *testing.T) {
				got := impl.fn(buildTree(tt.vals))
				if !slices.Equal(got, tt.want) {
					t.Fatalf("%s(%v) = %v, want %v", impl.name, tt.vals, got, tt.want)
				}
			})
		}
	}
}

// TestMaxConstraint 覆盖题目边界：节点数最多 100。
// 构造 100 个节点的右斜退化链（层序序列中节点值与 nil 交替出现），
// 递归深度达到 100 也验证了递归解法的栈安全。
func TestMaxConstraint(t *testing.T) {
	vals := make([]any, 2*100-1)
	for i := range 100 {
		vals[2*i] = i + 1
	}
	want := make([]int, 100)
	for i := range want {
		want[i] = i + 1
	}

	for _, impl := range impls {
		t.Run(impl.name, func(t *testing.T) {
			got := impl.fn(buildTree(vals))
			if !slices.Equal(got, want) {
				t.Fatalf("%s: 100 节点右斜链的中序应为 [1..100]", impl.name)
			}
		})
	}
}

// TestMorrisRestoresTree 验证 Morris 的关键性质：遍历结束后临时线索
// 全部撤除，树结构与遍历前完全一致——这是 Morris 实现最易出错的地方。
func TestMorrisRestoresTree(t *testing.T) {
	trees := [][]any{
		{1, nil, 2, 3},         // 官方示例 1
		{4, 2, 6, 1, 3, 5, 7},  // 满二叉树
		{4, 3, nil, 2, nil, 1}, // 左斜链：每个节点都要建/撤线索
	}
	for _, vals := range trees {
		root := buildTree(vals)
		backup := cloneTree(root)
		inorderMorris(root)
		if !treesEqual(root, backup) {
			t.Fatalf("Morris 遍历后树结构被破坏: %v", vals)
		}
	}
}

// cloneTree 深拷贝一棵树，用于 Morris 的前后对比。
func cloneTree(root *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}
	return &TreeNode{
		Val:   root.Val,
		Left:  cloneTree(root.Left),
		Right: cloneTree(root.Right),
	}
}

// treesEqual 递归比较两棵树的值与结构。
func treesEqual(a, b *TreeNode) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Val == b.Val && treesEqual(a.Left, b.Left) && treesEqual(a.Right, b.Right)
}
