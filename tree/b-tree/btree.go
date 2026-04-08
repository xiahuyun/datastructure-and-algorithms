package btree

import (
	"fmt"
	"sort"
	"strings"
)

// Item is an element stored in the B-Tree.
// Less defines a strict ordering for items.
type Item interface {
	Less(than Item) bool
}

// ItemIterator is used by traversal functions.
// Return false to stop iterating early.
type ItemIterator func(item Item) bool

// FreeList is kept for API compatibility with google/btree style constructors.
// This implementation does not reuse nodes through freelist.
type FreeList struct{}

// NewFreeList returns a freelist placeholder for API compatibility.
func NewFreeList(_ int) *FreeList {
	return &FreeList{}
}

// BTree is a classic B-Tree implementation.
//
// This implementation uses textbook "order m" semantics:
// - each node has at most m children
// - each node has at most m-1 items
// - each non-root node has at least ceil(m/2) children
type BTree struct {
	order  int
	root   *node
	length int
}

type node struct {
	leaf     bool
	items    []Item
	children []*node
}

// New creates a new B-Tree with textbook order m.
// order must be >= 3.
func New(order int) *BTree {
	if order < 3 {
		panic("btree: order must be at least 3")
	}
	return &BTree{
		order: order,
		root: &node{
			leaf: true,
		},
	}
}

// NewWithFreeList creates a new B-Tree.
// The freelist parameter is accepted for API compatibility.
func NewWithFreeList(order int, _ *FreeList) *BTree {
	return New(order)
}

func (t *BTree) minItems() int {
	return (t.order+1)/2 - 1
}

func (t *BTree) maxItems() int {
	return t.order - 1
}

// Len returns number of items in the tree.
func (t *BTree) Len() int {
	return t.length
}

// Clear removes all items from the tree.
// addNodesToFreelist is accepted for compatibility.
func (t *BTree) Clear(_ bool) {
	t.root = &node{leaf: true}
	t.length = 0
}

// Clone returns a deep copy of the tree structure.
func (t *BTree) Clone() *BTree {
	cloned := &BTree{
		order:  t.order,
		length: t.length,
	}
	if t.root == nil {
		cloned.root = &node{leaf: true}
		return cloned
	}
	cloned.root = cloneNode(t.root)
	return cloned
}

func cloneNode(n *node) *node {
	if n == nil {
		return nil
	}
	cp := &node{
		leaf:  n.leaf,
		items: append([]Item(nil), n.items...),
	}
	if !n.leaf {
		cp.children = make([]*node, len(n.children))
		for i := range n.children {
			cp.children[i] = cloneNode(n.children[i])
		}
	}
	return cp
}

// Get returns the equal item if found, otherwise nil.
func (t *BTree) Get(item Item) Item {
	if t.root == nil || t.length == 0 || item == nil {
		return nil
	}
	cur := t.root
	for {
		i, found := cur.find(item)
		if found {
			return cur.items[i]
		}
		if cur.leaf {
			return nil
		}
		cur = cur.children[i]
	}
}

// Min returns the minimum item in the tree.
func (t *BTree) Min() Item {
	if t.root == nil || t.length == 0 {
		return nil
	}
	n := t.root
	for !n.leaf {
		n = n.children[0]
	}
	return n.items[0]
}

// Max returns the maximum item in the tree.
func (t *BTree) Max() Item {
	if t.root == nil || t.length == 0 {
		return nil
	}
	n := t.root
	for !n.leaf {
		n = n.children[len(n.children)-1]
	}
	return n.items[len(n.items)-1]
}

// ReplaceOrInsert inserts item into the tree.
// If an equal item already exists, it is replaced and the old item is returned.
func (t *BTree) ReplaceOrInsert(item Item) Item {
	if item == nil {
		return nil
	}
	if t.root == nil {
		t.root = &node{leaf: true}
	}

	type pathElem struct {
		parent   *node
		childIdx int
	}

	cur := t.root
	var path []pathElem

	for {
		idx, found := cur.find(item)
		if found {
			old := cur.items[idx]
			cur.items[idx] = item
			return old
		}
		if cur.leaf {
			cur.items = insertItemAt(cur.items, idx, item)
			break
		}
		path = append(path, pathElem{
			parent:   cur,
			childIdx: idx,
		})
		cur = cur.children[idx]
	}

	t.length++

	// Bubble splits from leaf to root using recorded path.
	for len(cur.items) > t.maxItems() {
		promoted, right := t.splitOverflowNode(cur)

		if len(path) == 0 {
			t.root = &node{
				leaf:     false,
				items:    []Item{promoted},
				children: []*node{cur, right},
			}
			return nil
		}

		top := path[len(path)-1]
		path = path[:len(path)-1]

		parent := top.parent
		idx := top.childIdx
		parent.items = insertItemAt(parent.items, idx, promoted)
		parent.children = insertChildAt(parent.children, idx+1, right)
		cur = parent
	}
	return nil
}

// splitOverflowNode splits a node that temporarily holds maxItems()+1 keys.
// It promotes the median item and returns a new right sibling.
func (t *BTree) splitOverflowNode(n *node) (promoted Item, right *node) {
	mid := len(n.items) / 2
	promoted = n.items[mid]

	right = &node{
		leaf:  n.leaf,
		items: append([]Item(nil), n.items[mid+1:]...),
	}
	if !n.leaf {
		right.children = append([]*node(nil), n.children[mid+1:]...)
	}

	n.items = n.items[:mid]
	if !n.leaf {
		n.children = n.children[:mid+1]
	}
	return promoted, right
}

// DebugString returns a level-by-level tree view, useful for learning and debugging.
// Example:
// L0: [10]
// L1: [5] [20]
func (t *BTree) DebugString() string {
	if t.root == nil || t.length == 0 {
		return "<empty>"
	}

	var b strings.Builder
	level := 0
	queue := []*node{t.root}

	for len(queue) > 0 {
		if level > 0 {
			b.WriteByte('\n')
		}
		fmt.Fprintf(&b, "L%d:", level)

		next := make([]*node, 0, len(queue)*t.order)
		for _, n := range queue {
			b.WriteByte(' ')
			b.WriteString(formatItems(n.items))
			if !n.leaf {
				next = append(next, n.children...)
			}
		}
		queue = next
		level++
	}

	return b.String()
}

func formatItems(items []Item) string {
	var b strings.Builder
	b.WriteByte('[')
	for i, item := range items {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%v", item)
	}
	b.WriteByte(']')
	return b.String()
}

// Delete removes item equal to key and returns removed item.
func (t *BTree) Delete(item Item) Item {
	if item == nil || t.root == nil || t.length == 0 {
		return nil
	}

	cur := t.root
	target := item
	var removed Item

	for {
		idx, found := cur.find(target)
		if found {
			if removed == nil {
				removed = cur.items[idx]
			}
			if cur.leaf {
				cur.items, _ = removeItemAt(cur.items, idx)
				break
			}

			left := cur.children[idx]
			right := cur.children[idx+1]

			if len(left.items) > t.minItems() {
				pred := t.maxItem(left)
				cur.items[idx] = pred
				target = pred
				cur = left
				continue
			}
			if len(right.items) > t.minItems() {
				succ := t.minItem(right)
				cur.items[idx] = succ
				target = succ
				cur = right
				continue
			}

			t.mergeChildren(cur, idx)
			cur = cur.children[idx]
			continue
		}

		if cur.leaf {
			return nil
		}

		childIdx := idx
		if len(cur.children[childIdx].items) == t.minItems() {
			childIdx = t.fill(cur, childIdx)
		}
		cur = cur.children[childIdx]
	}

	t.length--
	t.shrinkRoot()
	return removed
}

// DeleteMin removes and returns minimum item.
func (t *BTree) DeleteMin() Item {
	if t.root == nil || t.length == 0 {
		return nil
	}

	cur := t.root
	for !cur.leaf {
		childIdx := 0
		if len(cur.children[childIdx].items) == t.minItems() {
			childIdx = t.fill(cur, childIdx)
		}
		cur = cur.children[childIdx]
	}

	if len(cur.items) == 0 {
		return nil
	}
	var removed Item
	cur.items, removed = removeItemAt(cur.items, 0)
	if removed != nil {
		t.length--
		t.shrinkRoot()
	}
	return removed
}

// DeleteMax removes and returns maximum item.
func (t *BTree) DeleteMax() Item {
	if t.root == nil || t.length == 0 {
		return nil
	}

	cur := t.root
	for !cur.leaf {
		childIdx := len(cur.children) - 1
		if len(cur.children[childIdx].items) == t.minItems() {
			childIdx = t.fill(cur, childIdx)
		}
		cur = cur.children[childIdx]
	}

	if len(cur.items) == 0 {
		return nil
	}
	last := len(cur.items) - 1
	var removed Item
	cur.items, removed = removeItemAt(cur.items, last)
	if removed != nil {
		t.length--
		t.shrinkRoot()
	}
	return removed
}

func (t *BTree) shrinkRoot() {
	if t.root == nil {
		t.root = &node{leaf: true}
		return
	}
	if len(t.root.items) == 0 && !t.root.leaf {
		t.root = t.root.children[0]
	}
	if len(t.root.items) == 0 && t.length == 0 {
		t.root = &node{leaf: true}
	}
}

func (t *BTree) fill(parent *node, idx int) int {
	if idx > 0 && len(parent.children[idx-1].items) > t.minItems() {
		t.borrowFromPrev(parent, idx)
		return idx
	}
	if idx < len(parent.children)-1 && len(parent.children[idx+1].items) > t.minItems() {
		t.borrowFromNext(parent, idx)
		return idx
	}

	if idx < len(parent.children)-1 {
		t.mergeChildren(parent, idx)
		return idx
	}
	t.mergeChildren(parent, idx-1)
	return idx - 1
}

func (t *BTree) borrowFromPrev(parent *node, idx int) {
	child := parent.children[idx]
	leftSibling := parent.children[idx-1]

	child.items = insertItemAt(child.items, 0, parent.items[idx-1])
	parent.items[idx-1] = leftSibling.items[len(leftSibling.items)-1]
	leftSibling.items = leftSibling.items[:len(leftSibling.items)-1]

	if !child.leaf {
		move := leftSibling.children[len(leftSibling.children)-1]
		leftSibling.children = leftSibling.children[:len(leftSibling.children)-1]
		child.children = insertChildAt(child.children, 0, move)
	}
}

func (t *BTree) borrowFromNext(parent *node, idx int) {
	child := parent.children[idx]
	rightSibling := parent.children[idx+1]

	child.items = append(child.items, parent.items[idx])
	parent.items[idx] = rightSibling.items[0]
	rightSibling.items = rightSibling.items[1:]

	if !child.leaf {
		move := rightSibling.children[0]
		rightSibling.children = rightSibling.children[1:]
		child.children = append(child.children, move)
	}
}

func (t *BTree) mergeChildren(parent *node, idx int) {
	left := parent.children[idx]
	right := parent.children[idx+1]

	left.items = append(left.items, parent.items[idx])
	left.items = append(left.items, right.items...)
	if !left.leaf {
		left.children = append(left.children, right.children...)
	}

	parent.items, _ = removeItemAt(parent.items, idx)
	parent.children, _ = removeChildAt(parent.children, idx+1)
}

func (t *BTree) minItem(n *node) Item {
	cur := n
	for !cur.leaf {
		cur = cur.children[0]
	}
	return cur.items[0]
}

func (t *BTree) maxItem(n *node) Item {
	cur := n
	for !cur.leaf {
		cur = cur.children[len(cur.children)-1]
	}
	return cur.items[len(cur.items)-1]
}

// Ascend calls iterator for every item in ascending order.
func (t *BTree) Ascend(iterator ItemIterator) {
	if t.root == nil || t.length == 0 || iterator == nil {
		return
	}
	t.scanAscend(t.root, iterator)
}

// AscendRange calls iterator for items in [greaterOrEqual, lessThan).
func (t *BTree) AscendRange(greaterOrEqual, lessThan Item, iterator ItemIterator) {
	if t.root == nil || t.length == 0 || iterator == nil {
		return
	}
	t.scanAscend(t.root, func(item Item) bool {
		if greaterOrEqual != nil && item.Less(greaterOrEqual) {
			return true
		}
		if lessThan != nil && !item.Less(lessThan) {
			return false
		}
		return iterator(item)
	})
}

// AscendLessThan calls iterator for items < pivot.
func (t *BTree) AscendLessThan(pivot Item, iterator ItemIterator) {
	t.AscendRange(nil, pivot, iterator)
}

// AscendGreaterOrEqual calls iterator for items >= pivot.
func (t *BTree) AscendGreaterOrEqual(pivot Item, iterator ItemIterator) {
	if t.root == nil || t.length == 0 || iterator == nil {
		return
	}
	t.scanAscend(t.root, func(item Item) bool {
		if pivot != nil && item.Less(pivot) {
			return true
		}
		return iterator(item)
	})
}

// Descend calls iterator for every item in descending order.
func (t *BTree) Descend(iterator ItemIterator) {
	if t.root == nil || t.length == 0 || iterator == nil {
		return
	}
	t.scanDescend(t.root, iterator)
}

// DescendRange calls iterator for items in (greaterThan, lessOrEqual], descending.
func (t *BTree) DescendRange(lessOrEqual, greaterThan Item, iterator ItemIterator) {
	if t.root == nil || t.length == 0 || iterator == nil {
		return
	}
	t.scanDescend(t.root, func(item Item) bool {
		if lessOrEqual != nil && lessOrEqual.Less(item) {
			return true
		}
		if greaterThan != nil && !greaterThan.Less(item) {
			return false
		}
		return iterator(item)
	})
}

// DescendLessOrEqual calls iterator for items <= pivot, descending.
func (t *BTree) DescendLessOrEqual(pivot Item, iterator ItemIterator) {
	t.DescendRange(pivot, nil, iterator)
}

// DescendGreaterThan calls iterator for items > pivot, descending.
func (t *BTree) DescendGreaterThan(pivot Item, iterator ItemIterator) {
	if t.root == nil || t.length == 0 || iterator == nil {
		return
	}
	t.scanDescend(t.root, func(item Item) bool {
		if pivot != nil && !pivot.Less(item) {
			return false
		}
		return iterator(item)
	})
}

func (t *BTree) scanAscend(n *node, iterator ItemIterator) bool {
	if n.leaf {
		for _, item := range n.items {
			if !iterator(item) {
				return false
			}
		}
		return true
	}

	for i, item := range n.items {
		if !t.scanAscend(n.children[i], iterator) {
			return false
		}
		if !iterator(item) {
			return false
		}
	}
	return t.scanAscend(n.children[len(n.items)], iterator)
}

func (t *BTree) scanDescend(n *node, iterator ItemIterator) bool {
	if n.leaf {
		for i := len(n.items) - 1; i >= 0; i-- {
			if !iterator(n.items[i]) {
				return false
			}
		}
		return true
	}

	if !t.scanDescend(n.children[len(n.items)], iterator) {
		return false
	}
	for i := len(n.items) - 1; i >= 0; i-- {
		if !iterator(n.items[i]) {
			return false
		}
		if !t.scanDescend(n.children[i], iterator) {
			return false
		}
	}
	return true
}

func (n *node) find(item Item) (int, bool) {
	idx := sort.Search(len(n.items), func(i int) bool {
		return !n.items[i].Less(item)
	})
	if idx < len(n.items) && !item.Less(n.items[idx]) && !n.items[idx].Less(item) {
		return idx, true
	}
	return idx, false
}

func insertItemAt(items []Item, idx int, item Item) []Item {
	items = append(items, nil)
	copy(items[idx+1:], items[idx:])
	items[idx] = item
	return items
}

func removeItemAt(items []Item, idx int) ([]Item, Item) {
	removed := items[idx]
	copy(items[idx:], items[idx+1:])
	last := len(items) - 1
	items[last] = nil
	return items[:last], removed
}

func insertChildAt(children []*node, idx int, child *node) []*node {
	children = append(children, nil)
	copy(children[idx+1:], children[idx:])
	children[idx] = child
	return children
}

func removeChildAt(children []*node, idx int) ([]*node, *node) {
	removed := children[idx]
	copy(children[idx:], children[idx+1:])
	last := len(children) - 1
	children[last] = nil
	return children[:last], removed
}
