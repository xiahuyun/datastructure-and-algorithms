# datastructure-and-algorithms 项目目录说明

这个仓库按数据结构和算法主题拆分，包含 Go 语言实现示例。

## 顶层目录用途

| 目录 | 用途 | 关键内容 |
| --- | --- | --- |
| `searching` | 查找算法示例 | `go.mod`, `main.go` |
| `sorting` | 多种排序算法实现 | `bubble-sort/`, `selection/`, `insertion/`, `quick-sort/`, `merge-sort/`, `heap-sort/` |
| `stack` | 栈结构实现与说明 | `README.md`, `array-stack/main.go`, `linked-stack/main.go` |
| `queue` | 队列结构实现 | `circular-queue/main.go`, `linked-queue/main.go` |
| `linked-list` | 链表题目与实现 | `reverse-linked-list/README.md`, `reverse-linked-list/main.go` |
| `recursion` | 递归示例 | `main.go` |

## sorting 子目录细分

| 子目录 | 用途 | 关键内容 |
| --- | --- | --- |
| `sorting/bubble-sort` | 冒泡排序 | `main.go` |
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
| `linked-list/reverse-linked-list` | 反转链表题解与实现 | `README.md`, `main.go` |

## 占位目录说明

| 目录 | 说明 |
| --- | --- |
| `linked-list/circular-linked-list` | 当前为预留目录，待补充环形链表示例 |
| `queue/blocking-queue` | 当前为预留目录，待补充阻塞队列示例 |
