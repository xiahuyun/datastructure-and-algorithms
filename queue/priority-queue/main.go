package main

import (
	"container/heap"
	"fmt"
)

// ============ 一、用 container/heap 做堆排序 ============

// IntHeap 是小顶堆（Less 返回 h[i] < h[j]）。
// 想改成大顶堆，把 Less 里的比较号反过来即可。
type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

// Push 往堆里塞一个元素（container/heap 会在末尾追加后自动上浮）。
func (h *IntHeap) Push(x interface{}) { *h = append(*h, x.(int)) }

// Pop 弹出堆顶元素（container/heap 会把末尾元素换到堆顶再下沉）。
func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// HeapSort 借助小顶堆实现升序排序。
func HeapSort(nums []int) []int {
	h := IntHeap(nums)
	heap.Init(&h) // 建堆，O(n)

	sorted := make([]int, 0, len(nums))
	for h.Len() > 0 {
		sorted = append(sorted, heap.Pop(&h).(int)) // 反复弹堆顶，O(n log n)
	}
	return sorted
}

// ============ 二、优先队列 ============

// Item 优先队列里的元素：priority 越小优先级越高（小顶堆）。
type Item struct {
	value    string // 业务值
	priority int    // 优先级，决定出队顺序
	index    int    // 该元素在堆中的下标，供 heap.Fix 使用
}

// PriorityQueue 实现 heap.Interface 的优先队列（小顶堆）。
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // 避免内存泄漏
	item.index = -1
	*pq = old[:n-1]
	return item
}

func main() {
	// 1. 堆排序
	nums := []int{3, 1, 4, 1, 5, 9, 2, 6}
	fmt.Println("堆排序结果:", HeapSort(nums))

	// 2. 优先队列：按优先级（数字小者优先）出队
	items := []*Item{
		{value: "吃饭", priority: 3},
		{value: "写代码", priority: 1},
		{value: "睡觉", priority: 4},
		{value: "健身", priority: 2},
	}
	pq := make(PriorityQueue, 0)
	for _, it := range items {
		heap.Push(&pq, it)
	}

	// 演示动态更新优先级：把「吃饭」提到最高优先级（0）
	items[0].priority = 0
	heap.Fix(&pq, items[0].index)

	// 演示运行中临时插入元素
	heap.Push(&pq, &Item{value: "喝水", priority: 2})

	fmt.Print("优先队列出队顺序: ")
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		fmt.Printf("%s(%d) ", item.value, item.priority)
	}
	fmt.Println()
}