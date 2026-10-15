package main

import (
	"slices"
	"testing"
)

// impls 汇总两种实现，保证每种实现都跑完全部同一组用例。
var impls = []struct {
	name string
	fn   func(*TreeNode) bool
}{
	{"递归", isSymmetricRecursive},
	{"迭代", isSymmetricIterative},
}

func TestIsSymmetric(t *testing.T) {
	tests := []struct {
		name string
		vals []any
		want bool
	}{
		{"示例1: 完全对称", []any{1, 2, 2, 3, 4, 4, 3}, true},
		{"示例2: 结构不镜像", []any{1, 2, 2, nil, 3, nil, 3}, false},
		{"空树: 对称", nil, true},
		{"单节点: 对称", []any{1}, true},
		{"根的左右值不同", []any{1, 2, 3}, false},
		{"镜像位置值不同", []any{1, 2, 2, 3, 4, 4, 5}, false},
		{"镜像位置结构不同", []any{1, 2, 2, 2, nil, 2, nil}, false},
		{"结构非满但恰好镜像", []any{1, 2, 2, 3, nil, nil, 3}, true},
	}

	for _, impl := range impls {
		for _, tt := range tests {
			t.Run(impl.name+"/"+tt.name, func(t *testing.T) {
				got := impl.fn(buildTree(tt.vals))
				if got != tt.want {
					t.Fatalf("%s(%v) = %v, want %v", impl.name, tt.vals, got, tt.want)
				}
			})
		}
	}
}

// TestMaxConstraint 覆盖题目边界：节点数最多 1000。
// 深度 9 的满二叉树共 511 个节点（同层节点同值），整树应镜像对称；
// 再把最右下角的叶子（全右路径）改成 99，它与镜像位置（最左下角叶子）
// 不再相等，两种实现都应判为不对称。
func TestMaxConstraint(t *testing.T) {
	symmetricVals := make([]any, 511)
	for i := range symmetricVals {
		d := 1 // 满二叉树按层序编号，节点 i 的深度 = floor(log2(i+1)) + 1
		for n := i + 1; n > 1; n /= 2 {
			d++
		}
		symmetricVals[i] = d
	}
	brokenVals := slices.Clone(symmetricVals)
	brokenVals[len(brokenVals)-1] = 99

	for _, impl := range impls {
		t.Run(impl.name+"/511节点满二叉树", func(t *testing.T) {
			if !impl.fn(buildTree(symmetricVals)) {
				t.Fatalf("%s: 511 节点的镜像满二叉树应判为对称", impl.name)
			}
		})
		t.Run(impl.name+"/破坏一个叶子", func(t *testing.T) {
			if impl.fn(buildTree(brokenVals)) {
				t.Fatalf("%s: 镜像位置的叶子值被破坏后应判为不对称", impl.name)
			}
		})
	}
}
