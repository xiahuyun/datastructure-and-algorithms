package main

import (
	"btree"
	"fmt"
	"os"
)

type Int int

func (a Int) Less(b btree.Item) bool {
	return a < b.(Int)
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--step" {
		runStepDemo()
		return
	}

	runNormalDemo()
}

func runNormalDemo() {
	tr := btree.New(3)

	for _, v := range []Int{10, 20, 5, 6, 12, 30, 7, 17} {
		tr.ReplaceOrInsert(v)
	}

	fmt.Printf("len=%d min=%v max=%v\n", tr.Len(), tr.Min(), tr.Max())
	fmt.Println("tree:")
	fmt.Println(tr.DebugString())

	fmt.Print("ascend: ")
	tr.Ascend(func(item btree.Item) bool {
		fmt.Printf("%v ", item)
		return true
	})
	fmt.Println()

	removed := tr.Delete(Int(12))
	fmt.Printf("delete 12 => %v\n", removed)

	fmt.Print("descend: ")
	tr.Descend(func(item btree.Item) bool {
		fmt.Printf("%v ", item)
		return true
	})
	fmt.Println()
}

func runStepDemo() {
	tr := btree.New(3)
	inserts := []Int{20, 10, 5, 6, 12, 30, 7, 17}

	fmt.Println("== Step Demo: Insert ==")
	for _, v := range inserts {
		tr.ReplaceOrInsert(v)
		fmt.Printf("insert %-2d =>\n%s\n\n", v, tr.DebugString())
	}

	fmt.Println("== Step Demo: Delete ==")
	for _, v := range []Int{12, 20, 6} {
		removed := tr.Delete(v)
		fmt.Printf("delete %-2d (removed=%v) =>\n%s\n\n", v, removed, tr.DebugString())
	}

	min := tr.DeleteMin()
	fmt.Printf("deleteMin (removed=%v) =>\n%s\n\n", min, tr.DebugString())

	max := tr.DeleteMax()
	fmt.Printf("deleteMax (removed=%v) =>\n%s\n\n", max, tr.DebugString())
}
