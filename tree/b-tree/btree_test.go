package btree

import (
	"math/rand"
	"slices"
	"testing"
)

type Int int

func (a Int) Less(b Item) bool {
	return a < b.(Int)
}

func TestReplaceOrInsertAndGet(t *testing.T) {
	tr := New(5)
	for i := 1; i <= 100; i++ {
		if old := tr.ReplaceOrInsert(Int(i)); old != nil {
			t.Fatalf("unexpected replaced item: %v", old)
		}
	}
	if tr.Len() != 100 {
		t.Fatalf("expected len 100, got %d", tr.Len())
	}

	old := tr.ReplaceOrInsert(Int(42))
	if old == nil || old.(Int) != 42 {
		t.Fatalf("expected replaced 42, got %v", old)
	}
	if tr.Len() != 100 {
		t.Fatalf("replace should not change len, got %d", tr.Len())
	}

	got := tr.Get(Int(77))
	if got == nil || got.(Int) != 77 {
		t.Fatalf("expected 77, got %v", got)
	}
	if tr.Get(Int(1000)) != nil {
		t.Fatalf("unexpected item for missing key")
	}
}

func TestDeleteDeleteMinDeleteMax(t *testing.T) {
	tr := New(5)
	for i := 1; i <= 40; i++ {
		tr.ReplaceOrInsert(Int(i))
	}

	for _, v := range []Int{7, 9, 14, 28, 33} {
		removed := tr.Delete(v)
		if removed == nil || removed.(Int) != v {
			t.Fatalf("expected delete %d, got %v", v, removed)
		}
	}
	if tr.Delete(Int(999)) != nil {
		t.Fatalf("delete of missing key should return nil")
	}

	min := tr.DeleteMin()
	max := tr.DeleteMax()
	if min == nil || min.(Int) != 1 {
		t.Fatalf("expected min 1, got %v", min)
	}
	if max == nil || max.(Int) != 40 {
		t.Fatalf("expected max 40, got %v", max)
	}
}

func TestAscendingAndDescending(t *testing.T) {
	tr := New(5)
	for i := 1; i <= 10; i++ {
		tr.ReplaceOrInsert(Int(i))
	}

	var asc []int
	tr.Ascend(func(item Item) bool {
		asc = append(asc, int(item.(Int)))
		return true
	})
	if !slices.Equal(asc, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}) {
		t.Fatalf("unexpected ascend result: %v", asc)
	}

	var ascRange []int
	tr.AscendRange(Int(4), Int(8), func(item Item) bool {
		ascRange = append(ascRange, int(item.(Int)))
		return true
	})
	if !slices.Equal(ascRange, []int{4, 5, 6, 7}) {
		t.Fatalf("unexpected ascend range: %v", ascRange)
	}

	var desc []int
	tr.Descend(func(item Item) bool {
		desc = append(desc, int(item.(Int)))
		return true
	})
	if !slices.Equal(desc, []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}) {
		t.Fatalf("unexpected descend result: %v", desc)
	}

	var descRange []int
	tr.DescendRange(Int(8), Int(3), func(item Item) bool {
		descRange = append(descRange, int(item.(Int)))
		return true
	})
	if !slices.Equal(descRange, []int{8, 7, 6, 5, 4}) {
		t.Fatalf("unexpected descend range: %v", descRange)
	}
}

func TestSplitKeepsMedianAsParentForOrder3(t *testing.T) {
	tr := New(3)
	tr.ReplaceOrInsert(Int(20))
	tr.ReplaceOrInsert(Int(10))
	tr.ReplaceOrInsert(Int(5))

	var asc []int
	tr.Ascend(func(item Item) bool {
		asc = append(asc, int(item.(Int)))
		return true
	})
	if !slices.Equal(asc, []int{5, 10, 20}) {
		t.Fatalf("unexpected sorted result after split: %v", asc)
	}

	if root := tr.root; root == nil || len(root.items) != 1 || int(root.items[0].(Int)) != 10 {
		t.Fatalf("expected root median to be 10, got %+v", root)
	}

	got := tr.DebugString()
	want := "L0: [10]\nL1: [5] [20]"
	if got != want {
		t.Fatalf("unexpected debug view:\nwant:\n%s\ngot:\n%s", want, got)
	}
}

func TestCloneAndClear(t *testing.T) {
	tr := New(5)
	for i := 1; i <= 20; i++ {
		tr.ReplaceOrInsert(Int(i))
	}
	cp := tr.Clone()
	if cp.Len() != tr.Len() {
		t.Fatalf("clone len mismatch: clone=%d tree=%d", cp.Len(), tr.Len())
	}

	cp.Delete(Int(10))
	if tr.Get(Int(10)) == nil {
		t.Fatalf("clone should be deep copy; original changed unexpectedly")
	}

	tr.Clear(true)
	if tr.Len() != 0 || tr.Min() != nil || tr.Max() != nil {
		t.Fatalf("clear did not reset tree")
	}
}

func TestRandomDeleteAgainstSortedSlice(t *testing.T) {
	tr := New(6)
	var ref []int

	for i := 1; i <= 200; i++ {
		tr.ReplaceOrInsert(Int(i))
		ref = append(ref, i)
	}

	r := rand.New(rand.NewSource(42))
	for i := 0; i < 120; i++ {
		idx := r.Intn(len(ref))
		key := ref[idx]

		removed := tr.Delete(Int(key))
		if removed == nil || int(removed.(Int)) != key {
			t.Fatalf("expected delete %d, got %v", key, removed)
		}

		ref = append(ref[:idx], ref[idx+1:]...)
	}

	var got []int
	tr.Ascend(func(item Item) bool {
		got = append(got, int(item.(Int)))
		return true
	})
	if !slices.Equal(got, ref) {
		t.Fatalf("tree/order mismatch after random deletes")
	}
}
