# 94. 二叉树的中序遍历（Binary Tree Inorder Traversal）

> https://leetcode.cn/problems/binary-tree-inorder-traversal/description/

## 题目描述

给定一个二叉树的根节点 `root`，返回它的**中序遍历**（左 -> 根 -> 右）。

```
输入: [1, null, 2, 3]        1
                              \
                               2
                              /
                             3
输出: [1, 3, 2]
```

约束：树中节点数在 `[0, 100]` 内。**进阶（follow-up）**：递归算法很简单，你可以通过迭代算法完成吗？

## 中序遍历为什么重要

对**二叉搜索树**做中序遍历恰好得到升序序列——这是 BST 最有用的性质。
本目录演示的满二叉树 `[4,2,6,1,3,5,7]` 就是 BST，中序输出 `[1,2,3,4,5,6,7]`。

## 递归解法（inorderRecursive）

中序定义的直译：`walk(左)` -> `访问根` -> `walk(右)`。

## 迭代解法（inorderIterative）

用显式栈模拟递归，核心是栈的不变式：**栈里保存"左子树尚未处理完"的祖先链**。
每轮三步：

1. 从当前节点一路向左，沿途全部压栈；
2. 弹出栈顶——它的左子树已处理完，轮到"根"，访问它；
3. 转向右子树，回到第 1 步。

## Morris 解法（inorderMorris）：O(1) 空间

递归栈和显式栈都需要 O(高度) 空间。Morris 不用任何栈：
**借中序前驱（左子树的最右节点）的空右指针当"回头路"**。

- 进入左子树前，把前驱的 `Right` 指向自己（临时线索）；
- 左子树遍历完会顺着线索恰好回到自己；
- 再次到达某节点且前驱指向自己：说明左子树完成——撤线索、访问、转右。

每条边最多走 2 次（一次找前驱、一次回溯），时间仍是 O(n)。
遍历过程中树被**临时修改**，结束后恢复原状——
`main_test.go` 的 `TestMorrisRestoresTree` 专门断言了这一点（这是 Morris 最易写错的地方）。

## 三方案对比

| 方案 | 时间 | 额外空间 | 一句话评价 |
|---|---|---|---|
| 递归 | O(n) | O(高度) 栈 | 定义直译，最易写 |
| 迭代 | O(n) | O(高度) 栈 | follow-up 要求，栈模板可复用到前/后序 |
| Morris | O(n) | O(1) | 空间最优，但要小心线索的建立与撤除 |

## 运行与测试

```bash
go run main.go   # 演示输出（同一棵树分别跑三种实现）
go test -v       # 单元测试
```

测试覆盖的用例（三种实现跑同一组）：

- LeetCode 五个官方示例（含空树、单节点、只有左/右孩子）
- 满二叉树（BST 中序升序）
- 左斜链（每个节点都要建立与撤除线索）
- 100 节点的右斜退化链（题目边界 + 递归深度 100 的栈安全）
- Morris 遍历后树结构完全恢复的专门断言

## 相关题目

- [144. 二叉树的前序遍历](https://leetcode.cn/problems/binary-tree-preorder-traversal/) / [145. 后序遍历](https://leetcode.cn/problems/binary-tree-postorder-traversal/)：同一套栈模板，只改访问时机
- [173. 二叉搜索树迭代器](https://leetcode.cn/problems/binary-search-tree-iterator/)：本质是"迭代中序的暂停版"
- [98. 验证二叉搜索树](https://leetcode.cn/problems/validate-binary-search-tree/)：利用中序有序性
- [538. 把二叉搜索树转换为累加树](https://leetcode.cn/problems/convert-bst-to-greater-tree/)：反序中序（右 -> 根 -> 左）
