# 101. 对称二叉树（Symmetric Tree）

> https://leetcode.cn/problems/symmetric-tree/description/

## 题目描述

给你一个二叉树的根节点 `root`，检查它是否**轴对称**。

示例 1（对称）：

```
        1
       / \
      2   2
     / \ / \
    3  4 4  3
```

示例 2（不对称）：

```
        1
       / \
      2   2
       \    \
        3    3
```

**进阶（follow-up）**：可以运用递归和迭代两种方法解决本题吗？

## 核心洞察：把"对称"翻译成"镜像"

整棵树对称 ⟺ 根的**左右子树互为镜像**。两棵树 A、B 互为镜像的递归定义：

1. 两棵都为空，或都非空且 `A.Val == B.Val`；
2. **交叉比较**：A 的左子树与 B 的右子树互为镜像，A 的右子树与 B 的左子树互为镜像。

"左对右、右对左"的交叉是精髓：沿对称轴翻转后，左就变成了右。
如果顺手写成"左对左、右对右"的同向比较，那就变成了 [100. 相同的树](https://leetcode.cn/problems/same-tree/)。

## 递归解法（isSymmetricRecursive）

`isSymmetric(root) = isMirror(root, root)`：从根出发把左右子树当成两棵树来比较，
空树天然对称（`isMirror(nil, nil) = true`）。
每个节点最多访问一次。时间 O(n)；空间 O(高度) 递归栈，树退化为链表时最坏 O(n)。

## 迭代解法（isSymmetricIterative）

用队列替代递归栈：队列中**成对**存放待比较的节点 `(l, r)`，
每次出队一对，恰一边为空或值不等立即返回 false；
否则把 `(l.Left, r.Right)` 与 `(l.Right, r.Left)` 两对继续入队。
子问题拆分与递归完全一致。时间 O(n)；空间 O(n) 队列。

## 复杂度对比

| 方案 | 时间 | 空间 |
|---|---|---|
| 递归 | O(n) | O(高度)，最坏 O(n) |
| 迭代 | O(n) | O(n) |

## 运行与测试

```bash
go run main.go   # 演示输出（同一棵树分别跑两种实现）
go test -v       # 单元测试
```

测试覆盖的用例（两种实现跑同一组）：

- LeetCode 两个官方示例
- 空树、单节点（边界）
- 根的左右值不同
- 镜像位置值不同 / 结构不同（考察交叉比较是否正确）
- 结构非满但恰好镜像（避免只看形状不看位置）
- 511 节点的镜像满二叉树 + 破坏一个叶子后的对照（题目约束节点数 ≤ 1000）

## 相关题目

- [100. 相同的树](https://leetcode.cn/problems/same-tree/)：把交叉比较改成同向比较
- [226. 翻转二叉树](https://leetcode.cn/problems/invert-binary-tree/)：树对称 ⟺ `invert(root) == root`
- [572. 另一棵树的子树](https://leetcode.cn/problems/subtree-of-another-tree/)
