// 跳表（skip list）示例：基于 github.com/huandu/skiplist 演示
// 1. 基本增删改查与遍历
// 2. 区间查询（Find + Next）
// 3. 并发读写的封装（huandu/skiplist 本身非并发安全，需加锁）

package main

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/huandu/skiplist"
)

// ConcurrentSkipList 在 huandu/skiplist 之上加读写锁，提供并发安全的有序 map。
type ConcurrentSkipList struct {
	mu   sync.RWMutex
	list *skiplist.SkipList
}

// NewConcurrentSkipList 创建一个 key 为 int 的并发安全跳表。
func NewConcurrentSkipList() *ConcurrentSkipList {
	return &ConcurrentSkipList{list: skiplist.New(skiplist.Int)}
}

// Set 写入 key -> value。
func (c *ConcurrentSkipList) Set(key int, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.list.Set(key, value)
}

// Get 读取 key 对应的 value。
func (c *ConcurrentSkipList) Get(key int) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.list.GetValue(key)
}

// Remove 删除 key。
func (c *ConcurrentSkipList) Remove(key int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.list.Remove(key)
}

// Len 返回元素个数。
func (c *ConcurrentSkipList) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.list.Len()
}

// Range 返回 key 落在 [start, end] 闭区间内的所有 value（按 key 升序）。
func (c *ConcurrentSkipList) Range(start, end int) []interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var result []interface{}
	for e := c.list.Find(start); e != nil; e = e.Next() {
		key := e.Key().(int)
		if key > end {
			break
		}
		result = append(result, e.Value)
	}
	return result
}

func main() {
	demoBasic()
	demoRange()
	demoConcurrent()
}

// demoBasic 演示基本操作：Set / Get / GetValue / 遍历 / Remove。
func demoBasic() {
	fmt.Println("== 1. 基本操作 ==")

	list := skiplist.New(skiplist.Int)

	list.Set(20, "hello")
	list.Set(10, "world")
	list.Set(40, 3.14)
	list.Set(40, 1000) // 覆盖已有 key

	// Get 返回元素节点，节点上的 Value 是导出字段。
	e := list.Get(10)
	fmt.Printf("Get(10) -> value=%v\n", e.Value)

	// GetValue 直接拿值，类似 map。
	if v, ok := list.GetValue(20); ok {
		fmt.Printf("GetValue(20) -> %v\n", v)
	}

	// 从头到尾遍历。
	fmt.Print("遍历全部: ")
	for e := list.Front(); e != nil; e = e.Next() {
		fmt.Printf("{%v:%v} ", e.Key(), e.Value)
	}
	fmt.Println()

	// 删除并拿到被删元素。
	removed := list.Remove(40)
	fmt.Printf("Remove(40) -> removed=%v, 剩余 Len=%d\n", removed.Value, list.Len())
	fmt.Println()
}

// demoRange 演示区间查询：先 Find 定位下界，再 Next 向后遍历直到越过上界。
func demoRange() {
	fmt.Println("== 2. 区间查询 ==")

	list := skiplist.New(skiplist.Int)
	for i := 0; i <= 100; i += 10 {
		list.Set(i, fmt.Sprintf("score-%d", i))
	}

	// 查询 [25, 70] 闭区间内的元素。
	fmt.Println("区间 [25, 70] 内的元素:")
	for e := list.Find(25); e != nil; e = e.Next() {
		key := e.Key().(int)
		if key > 70 {
			break
		}
		fmt.Printf("  {%v: %v}\n", key, e.Value)
	}
	fmt.Println()
}

// demoConcurrent 演示并发读写：多个 goroutine 同时读写封装后的 ConcurrentSkipList。
func demoConcurrent() {
	fmt.Println("== 3. 并发读写 ==")

	cl := NewConcurrentSkipList()

	const (
		writers = 8
		readers = 4
		ops     = 1000
	)

	var wg sync.WaitGroup

	// 并发写：每个 writer 负责一段不重叠的 key 区间，保证最终结果可预期。
	const perWriter = ops / writers // 1000/8 = 125，最终共 1000 个不同 key
	for w := 0; w < writers; w++ {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			for i := 0; i < perWriter; i++ {
				cl.Set(start+i, start)
			}
		}(w * perWriter)
	}

	// 并发读。
	for r := 0; r < readers; r++ {
		wg.Add(1)
		go func(seed int) {
			defer wg.Done()
			rnd := rand.New(rand.NewSource(int64(seed + 100)))
			for i := 0; i < ops; i++ {
				_, _ = cl.Get(rnd.Intn(1000))
			}
		}(r)
	}

	wg.Wait()
	fmt.Printf("并发读写完成，最终元素数 Len = %d\n", cl.Len())
}