// https://leetcode.cn/problems/lowest-common-ancestor-of-a-binary-tree/description/

package main

import "fmt"

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

// findNode 按值查找节点（题目保证节点值互不相同）。
func findNode(root *TreeNode, val int) *TreeNode {
	if root == nil || root.Val == val {
		return root
	}
	if n := findNode(root.Left, val); n != nil {
		return n
	}
	return findNode(root.Right, val)
}

// lowestCommonAncestorRecursive 递归解法（后序遍历，一次扫描）。
//
// 本题灵魂是定义清楚递归函数的"返回值语义"——每个节点向父节点汇报
// 自己子树里 p、q 的寻找结果：
//
//	子树中 p、q 都没有          → nil
//	子树中只找到 p              → p
//	子树中只找到 q              → q
//	子树中 p、q 都在（含已出答案）→ 它们的最近公共祖先
//
// 在每个节点上合并左右两份报告：
//   - 两边都非空：p、q 分居两侧，两条"根→目标"路径在此分叉 → 当前节点就是答案；
//   - 只有一边非空：p、q 同在该侧子树深处，把这份报告原样上传；
//   - 两边都空：上传 nil。
//
// 出口 `root == nil || root == p || root == q` 覆盖了一个容易踩的坑：
// 当 p 是 q 的祖先时（题目约定节点可作为自己的后代），碰到 p 就立即返回、
// 不再下钻。q 的存在性由题目保证，即使它藏在 p 的子树里也无需找到——
// "p 的子树里有没有 q"不影响答案仍是 p。
//
// 时间 O(n)：每个节点至多访问一次；空间 O(高度) 递归栈，斜树最坏 O(n)。
func lowestCommonAncestorRecursive(root, p, q *TreeNode) *TreeNode {
	if root == nil || root == p || root == q {
		return root
	}
	left := lowestCommonAncestorRecursive(root.Left, p, q)
	right := lowestCommonAncestorRecursive(root.Right, p, q)
	if left != nil && right != nil {
		return root // 两份报告在此汇合：路径分叉点
	}
	if left != nil {
		return left // p、q 同在左子树，答案在深处，报告继续上传
	}
	return right // 可能为 nil：两侧都没有找到
}

// findPath 记录从 root 到 target 的路径（含两端），找到返回 true。
// 回溯是关键：某棵子树宣告失败时必须把当前节点弹出路径，
// 否则路径会混入"不在通往 target 上的节点"。
func findPath(root, target *TreeNode, path *[]*TreeNode) bool {
	if root == nil {
		return false
	}
	*path = append(*path, root)
	if root == target ||
		findPath(root.Left, target, path) ||
		findPath(root.Right, target, path) {
		return true
	}
	*path = (*path)[:len(*path)-1] // 撤销选择
	return false
}

// lowestCommonAncestorPath 路径解法：显式存下"根→p"、"根→q"两条路径再比对。
// 两条路径上半段必然重合、在分叉点分开；从头逐位对齐，
// 最后一个相同节点即最近公共祖先。
//
// 它是递归解法的"直觉版"：递归解法把"存两条路径再比对"压缩成
// "边递归边把比对结果往上带"——信息不显式存储，而是藏在返回值里逐层上传；
// 本解法则把两条路径完整摊开在桌面上从头对齐。
//
// 时间 O(n)，空间 O(n)（两条路径 + 回溯栈）。
func lowestCommonAncestorPath(root, p, q *TreeNode) *TreeNode {
	var pathP, pathQ []*TreeNode
	findPath(root, p, &pathP)
	findPath(root, q, &pathQ)
	var lca *TreeNode
	for i := 0; i < len(pathP) && i < len(pathQ); i++ {
		if pathP[i] != pathQ[i] {
			break
		}
		lca = pathP[i]
	}
	return lca
}

func main() {
	demos := []struct {
		desc string
		vals []any
		pVal int
		qVal int
	}{
		{"示例1: 分居根两侧", []any{3, 5, 1, 6, 2, 0, 8, nil, nil, 7, 4}, 5, 1},
		{"示例2: p是q的祖先", []any{3, 5, 1, 6, 2, 0, 8, nil, nil, 7, 4}, 5, 4},
		{"示例3: 根与孩子", []any{1, 2}, 1, 2},
		{"同在左子树深处", []any{3, 5, 1, 6, 2, 0, 8, nil, nil, 7, 4}, 6, 4},
		{"同在右子树深处", []any{3, 5, 1, 6, 2, 0, 8, nil, nil, 7, 4}, 0, 8},
	}
	for _, d := range demos {
		root := buildTree(d.vals)
		p, q := findNode(root, d.pVal), findNode(root, d.qVal)
		fmt.Printf("%-18s p=%-2d q=%-2d -> 递归: %d, 路径: %d\n", d.desc, d.pVal, d.qVal,
			lowestCommonAncestorRecursive(root, p, q).Val,
			lowestCommonAncestorPath(root, p, q).Val)
	}
}
