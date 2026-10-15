# datastructure-and-algorithms 项目目录说明

这个仓库按数据结构和算法主题拆分，包含 Go 语言实现示例。

## 顶层目录用途

| 目录 | 用途 | 关键内容 |
| --- | --- | --- |
| `searching` | 查找算法示例 | `go.mod`, `main.go` |
| `array` | 数组类题目示例 | `color-type/main.go` |
| `sorting` | 多种排序算法实现 | `bubble/`, `selection/`, `insertion/`, `quick-sort/`, `merge-sort/`, `heap-sort/` |
| `stack` | 栈结构实现与说明 | `README.md`, `array-stack/main.go`, `linked-stack/main.go` |
| `queue` | 队列结构实现 | `circular-queue/main.go`, `linked-queue/main.go`, `priority-queue/main.go` |
| `linked-list` | 链表题目与实现 | `reverse-linked-list/`, `add-two-numbers/`, `add-two-numbers-ii/` |
| `recursion` | 递归示例 | `main.go` |
| `tree` | 树结构实现 | `b-tree/`, `symmetric-tree/`, `binary-tree-inorder-traversal/`, `lowest-common-ancestor-of-a-binary-tree/` |
| `graph` | 图遍历算法示例（BFS/DFS） | BFS: `main.go`; DFS: `dfs/` |
| `backtracking` | 回溯算法示例 | `eight-queens/main.go`, `eight-queens/README.md` |
| `id-generator` | 全局唯一 ID 生成算法示例 | `snowflake/`, `segment/` |
| `goroutine-mpx1` | Go 调度器（GMP）学习示例与可视化动画 | `main.go`, `scheduler-animation.html`, `gmp-animation.html`, `README.md` |

## sorting 子目录细分

| 子目录 | 用途 | 关键内容 |
| --- | --- | --- |
| `sorting/bubble` | 冒泡排序 | `main.go` |
| `sorting/selection` | 选择排序 | `main.go` |
| `sorting/insertion` | 插入排序 | `go.mod`, `main.go` |
| `sorting/quick-sort` | 快速排序 | `go.mod`, `main.go` |
| `sorting/merge-sort` | 归并排序 | `go.mod`, `main.go` |
| `sorting/heap-sort` | 堆排序 | `go.mod`, `main.go` |

## 结构类子目录细分

| 子目录 | 用途 | 关键内容 |
| --- | --- | --- |
| `stack/array-stack` | 顺序栈（数组实现） | `main.go` |
| `stack/linked-stack` | 链式栈（链表实现） | `main.go` |
| `queue/circular-queue` | 循环队列 | `main.go` |
| `queue/linked-queue` | 链式队列 | `main.go` |
| `queue/priority-queue` | 优先队列（container/heap） | `main.go` |
| `linked-list/reverse-linked-list` | 反转链表题解与实现 | `README.md`, `main.go` |
| `linked-list/add-two-numbers` | 2. 两数相加（链表竖式加法） | `README.md`, `main.go` |
| `linked-list/add-two-numbers-ii` | 445. 两数相加 II（反转/栈/递归三解法） | `README.md`, `main.go` |
| `tree/b-tree` | 教材语义 B-Tree（m 阶） | `go.mod`, `btree.go`, `btree_test.go`, `example/main.go` |
| `tree/symmetric-tree` | 101. 对称二叉树（递归/迭代两解法） | `README.md`, `main.go` |
| `tree/binary-tree-inorder-traversal` | 94. 二叉树的中序遍历（递归/迭代/Morris） | `README.md`, `main.go` |
| `tree/lowest-common-ancestor-of-a-binary-tree` | 236. 二叉树的最近公共祖先（递归/路径两解法） | `README.md`, `main.go` |

## tree/b-tree 子目录说明

| 文件/目录 | 用途 |
| --- | --- |
| `tree/b-tree/btree.go` | B-Tree 核心实现，接口风格类似 google/btree |
| `tree/b-tree/btree_test.go` | 插入、删除、范围遍历、克隆等测试 |
| `tree/b-tree/example/main.go` | 普通演示与 `--step` 逐步演示模式 |
| `tree/b-tree/README.md` | 该子项目的运行方式与实现说明 |

## 其他子目录说明

| 子目录 | 用途 | 关键内容 |
| --- | --- | --- |
| `array/color-type` | 75. 颜色分类（荷兰国旗问题，原地分区） | `go.mod`, `main.go` |
| `goroutine-mpx1` | `runtime.GOMAXPROCS(1)` 下的多协程调度行为实验，附带调度可视化动画 | `main.go`, `scheduler-animation.html`, `gmp-animation.html`, `README.md` |

## 占位目录说明

| 目录 | 说明 |
| --- | --- |
| `linked-list/circular-linked-list` | 当前为预留目录，待补充环形链表示例 |
| `queue/blocking-queue` | 当前为预留目录，待补充阻塞队列示例 |
