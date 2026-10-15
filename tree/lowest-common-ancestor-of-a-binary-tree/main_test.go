package main

import "testing"

// impls 汇总两种实现，保证每种实现都跑完全部同一组用例。
var impls = []struct {
	name string
	fn   func(root, p, q *TreeNode) *TreeNode
}{
	{"递归", lowestCommonAncestorRecursive},
	{"路径", lowestCommonAncestorPath},
}

func TestLowestCommonAncestor(t *testing.T) {
	tests := []struct {
		name string
		vals []any
		pVal int
		qVal int
		want int
	}{
		{"示例1: p和q分居根两侧", []any{3, 5, 1, 6, 2, 0, 8, nil, nil, 7, 4}, 5, 1, 3},
		{"示例2: p是q的祖先", []any{3, 5, 1, 6, 2, 0, 8, nil, nil, 7, 4}, 5, 4, 5},
		{"示例3: 根是q的祖先", []any{1, 2}, 1, 2, 1},
		{"同在左子树: 答案在深处", []any{3, 5, 1, 6, 2, 0, 8, nil, nil, 7, 4}, 6, 4, 5},
		{"同在右子树: 答案在深处", []any{3, 5, 1, 6, 2, 0, 8, nil, nil, 7, 4}, 0, 8, 1},
		{"互为兄弟的叶子", []any{3, 5, 1, 6, 2, 0, 8, nil, nil, 7, 4}, 7, 4, 2},
		{"左斜链: 祖先即链上节点", []any{1, 2, nil, 3, nil, 4, nil, 5}, 4, 5, 4},
	}

	for _, impl := range impls {
		for _, tt := range tests {
			t.Run(impl.name+"/"+tt.name, func(t *testing.T) {
				root := buildTree(tt.vals)
				p, q := findNode(root, tt.pVal), findNode(root, tt.qVal)
				if got := impl.fn(root, p, q); got.Val != tt.want {
					t.Fatalf("%s: LCA(%d, %d) = %d, want %d", impl.name, tt.pVal, tt.qVal, got.Val, tt.want)
				}
			})
		}
	}
}

// TestMaxConstraint 覆盖题目边界：节点数最多 10^5。
// 左斜链同时考察最坏情况：递归深度 = 节点数。
// 链上 LCA(最深的 n, 次深的 n-1) = n-1（祖先即邻居）；
// LCA(链首的孩子 2, 链尾 n) = 2（祖先靠近根）。
func TestMaxConstraint(t *testing.T) {
	const n = 100000
	vals := make([]any, 0, 2*n)
	vals = append(vals, 1)
	for v := 2; v <= n; v++ {
		vals = append(vals, v, nil) // 每个节点只有左孩子
	}
	root := buildTree(vals)

	for _, impl := range impls {
		t.Run(impl.name+"/10^5节点左斜链", func(t *testing.T) {
			deep, second := findNode(root, n), findNode(root, n-1)
			if got := impl.fn(root, deep, second); got.Val != n-1 {
				t.Fatalf("%s: LCA(%d, %d) = %d, want %d", impl.name, n, n-1, got.Val, n-1)
			}
			head := findNode(root, 2)
			if got := impl.fn(root, head, deep); got.Val != 2 {
				t.Fatalf("%s: LCA(2, %d) = %d, want 2", impl.name, n, got.Val)
			}
		})
	}
}
