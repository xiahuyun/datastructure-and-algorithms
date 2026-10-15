package main

import "fmt"

const n = 8 // 八皇后

// abs 返回 x 的绝对值。
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// isValid 判断在第 row 行第 col 列放皇后，是否与前面已放好的皇后冲突。
//
// 用一维数组表示棋盘：queens[i] 存的是第 i 行的皇后放在哪一列。
// 天然保证了「不同行」，所以只需检查「不同列」和「不同对角线」。
func isValid(queens []int, row, col int) bool {
	for i := 0; i < row; i++ {
		// 同列冲突：queens[i] == col
		// 同对角线冲突：行差 == 列差，即 row-i == abs(col-queens[i])
		if queens[i] == col || row-i == abs(col-queens[i]) {
			return false
		}
	}
	return true
}

// solveNQueens 返回 N 皇后的所有解（递归回溯）。
//
// 思路：一行一行地放。当前行尝试每一列，能放就放、递归进入下一行；
// 放不下（或后续走不通）就回到当前行换下一列继续试，这就是「回溯」。
func solveNQueens(n int) [][]int {
	solutions := [][]int{}
	queens := make([]int, n)

	var backtrack func(row int)
	backtrack = func(row int) {
		// 递归出口：n 行都放满，得到一个完整解。
		if row == n {
			sol := make([]int, n)
			copy(sol, queens) // 深拷贝，避免后续回溯覆盖
			solutions = append(solutions, sol)
			return
		}

		// 在当前第 row 行，依次尝试每一列。
		for col := 0; col < n; col++ {
			if isValid(queens, row, col) {
				queens[row] = col  // 放下
				backtrack(row + 1) // 递归到下一行
				// 回溯：无需显式撤销，下一轮循环会直接覆盖 queens[row]
			}
		}
	}

	backtrack(0)
	return solutions
}

// printBoard 打印一个解的棋盘，Q 表示皇后，. 表示空格。
func printBoard(queens []int) {
	for row := 0; row < len(queens); row++ {
		line := make([]byte, len(queens))
		for i := range line {
			line[i] = '.'
		}
		line[queens[row]] = 'Q'
		fmt.Println(string(line))
	}
}

func main() {
	solutions := solveNQueens(n)

	fmt.Printf("八皇后共有 %d 个解\n\n", len(solutions))

	fmt.Println("前 3 个解（每项表示对应行皇后的列号，从 0 开始）:")
	for i := 0; i < 3 && i < len(solutions); i++ {
		fmt.Printf("  解 %d: %v\n", i+1, solutions[i])
	}

	fmt.Println("\n第一个解的棋盘:")
	printBoard(solutions[0])
}