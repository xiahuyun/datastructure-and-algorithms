package main

import "fmt"

// Graph 用邻接表表示无向图。
type Graph map[string][]string

// visit 是访问回调，每当首次访问到一个节点时触发。
type visit func(string)

// DFSRecursive 递归版深度优先搜索。
func DFSRecursive(g Graph, start string, visit visit) {
	visited := map[string]bool{}

	var dfs func(string)
	dfs = func(cur string) {
		visited[cur] = true
		visit(cur)
		// 每次都挑一个还没去过的邻居，一路走到头，再回头换下一条路。
		for _, next := range g[cur] {
			if !visited[next] {
				dfs(next)
			}
		}
	}
	dfs(start)
}

// DFSIterative 迭代版深度优先搜索，用栈显式模拟递归的调用栈。
func DFSIterative(g Graph, start string, visit visit) {
	visited := map[string]bool{}
	stack := []string{start}

	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if visited[cur] {
			continue
		}
		visited[cur] = true
		visit(cur)

		// 逆序入栈，保证弹出顺序与递归版一致。
		for i := len(g[cur]) - 1; i >= 0; i-- {
			if !visited[g[cur][i]] {
				stack = append(stack, g[cur][i])
			}
		}
	}
}

// DFSFindPath 从 start 走到 target，返回一条路径（不保证最短）。
func DFSFindPath(g Graph, start, target string) ([]string, bool) {
	visited := map[string]bool{}
	var path []string

	var dfs func(string) bool
	dfs = func(cur string) bool {
		visited[cur] = true
		path = append(path, cur)

		if cur == target {
			return true
		}
		for _, next := range g[cur] {
			if !visited[next] && dfs(next) {
				return true
			}
		}
		path = path[:len(path)-1] // 此路不通，回退一步（回溯）
		return false
	}

	if dfs(start) {
		return path, true
	}
	return nil, false
}

func main() {
	g := Graph{
		"A": {"B", "C"},
		"B": {"A", "D"},
		"C": {"A", "D"},
		"D": {"B", "C", "E", "F"},
		"E": {"D", "G"},
		"F": {"D", "G"},
		"G": {"E", "F"},
	}

	fmt.Print("递归 DFS 访问顺序: ")
	DFSRecursiveV2(g, "A", func(s string) { fmt.Print(s, " ") })
	fmt.Println()

	fmt.Print("迭代 DFS 访问顺序: ")
	DFSIterative(g, "A", func(s string) { fmt.Print(s, " ") })
	fmt.Println()

	if path, ok := DFSFindPath(g, "A", "G"); ok {
		fmt.Printf("DFS 找到的一条路径（不保证最短）: %v\n", path)
	}
}

func DFSRecursiveV2(g Graph, start string, visit visit) {
	visited := map[string]bool{}

	var dfs func(string)
	dfs = func(cur string) {
		visited[cur] = true
		visit(cur)
		for _, next := range g[cur] {
			if visited[next] {
				continue
			}
			dfs(next)
		}
	}

	dfs(start)
}
