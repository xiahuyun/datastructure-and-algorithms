package main

import "fmt"

// Graph 用邻接表表示无向无权图。
type Graph map[string][]string

// BFS 在无权图上做广度优先搜索，返回从 start 到 target 的最短路径
// （边数最少的路径）、边数距离，以及是否可达。
//
// BFS 保证在「每条边权重相同（无权）」时第一次访问到某个节点时，
// 走过的边数就是最短距离。
func BFS(g Graph, start, target string) ([]string, int, bool) {
	if start == target {
		return []string{start}, 0, true
	}

	parent := map[string]string{start: ""} // 记录每个节点由谁访问到，用于重建路径
	dist := map[string]int{start: 0}       // 起点到各节点的最短边数

	queue := []string{start}
	head := 0 // 队头下标，避免切片头部删除带来的 O(n) 开销

	for head < len(queue) {
		cur := queue[head]
		head++

		if cur == target {
			return rebuildPath(parent, target), dist[cur], true
		}

		for _, next := range g[cur] {
			if _, seen := parent[next]; seen {
				continue // 已经访问过，说明存在更短或同样的路径，跳过
			}
			parent[next] = cur
			dist[next] = dist[cur] + 1
			queue = append(queue, next)
		}
	}

	return nil, 0, false
}

// rebuildPath 从 parent 回溯，重建 start -> ... -> target 的路径。
func rebuildPath(parent map[string]string, target string) []string {
	path := []string{}
	for cur := target; cur != ""; cur = parent[cur] {
		path = append([]string{cur}, path...)
	}
	return path
}

func main() {
	// 无向无权图：
	//
	//     A -- B
	//     |    |
	//     C -- D -- E
	//          |    |
	//          F -- G
	g := Graph{
		"A": {"B", "C"},
		"B": {"A", "D"},
		"C": {"A", "D"},
		"D": {"B", "C", "E", "F"},
		"E": {"D", "G"},
		"F": {"D", "G"},
		"G": {"E", "F"},
	}

	start, target := "A", "G"
	path, dist, ok := BFSV2(g, start, target)
	if !ok {
		fmt.Printf("%s 到 %s 不可达\n", start, target)
		return
	}

	fmt.Printf("从 %s 到 %s 的最短路径: %v\n", start, target, path)
	fmt.Printf("最短边数: %d\n", dist)
}

func BFSV2(g Graph, start, target string) ([]string, int, bool) {
	if start == target {
		return []string{start}, 0, true
	}

	visited := map[string]bool{start: true} // 记录访问过的节点，避免重复访问
	queue := []string{start}
	dist := map[string]int{start: 0}       // 记录每个节点的路径长度
	parent := map[string]string{start: ""} // 记录每个节点的父节点，用于重建路径

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if cur == target {
			return rebuildPathV2(parent, target), dist[cur], true
		}

		for _, next := range g[cur] {
			if _, seen := visited[next]; seen {
				continue // 已经访问过，跳过
			}
			visited[next] = true
			queue = append(queue, next)
			dist[next] = dist[cur] + 1
			parent[next] = cur
		}
	}
	return nil, 0, false
}

func rebuildPathV2(parent map[string]string, target string) []string {
	path := []string{}
	for cur := target; cur != ""; cur = parent[cur] {
		path = append([]string{cur}, path...)
	}
	return path
}
