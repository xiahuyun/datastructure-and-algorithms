package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// 数据库号段模式（Leaf-segment / TinyID 的思路）。
//
// 核心思想：不再每生成一个 ID 都访问一次数据库，而是「批量取号」——
// 一次性从数据库取走一段号段（比如 1000 个），缓存在本地内存里慢慢发，
// 发完再去取下一段。数据库的访问频率因此从「每请求 1 次」降到「每 1000 次 1 次」。

// DB 模拟数据库里的 id_alloc 号段表：
//
//	CREATE TABLE id_alloc (
//	  biz_tag  varchar(128) PRIMARY KEY, -- 业务标识
//	  max_id   bigint     NOT NULL,      -- 当前已分配的最大 ID
//	  step     int        NOT NULL       -- 每次取号的步长
//	);
type DB struct {
	mu       sync.Mutex
	maxID    map[string]int64
	allocCnt atomic.Int64 // 统计「取号段」次数，即命中数据库的次数（用于演示）
}

// newDB 创建一个空的号段表。
func newDB() *DB {
	return &DB{maxID: map[string]int64{}}
}

// allocSegment 原子地「取走一段号段」，等价于这条 SQL：
//
//	UPDATE id_alloc SET max_id = max_id + step WHERE biz_tag = ?
//
// 单条 UPDATE 在数据库里是原子的（行锁），天生并发安全，这是整个方案唯一
// 需要「锁」的地方。返回这段号段的 [start, end]（闭区间）。
func (db *DB) allocSegment(bizTag string, step int64) (start, end int64) {
	db.mu.Lock()
	defer db.mu.Unlock()

	cur := db.maxID[bizTag]
	db.maxID[bizTag] = cur + step
	db.allocCnt.Add(1)
	return cur + 1, cur + step
}

// Allocator 每个业务一个实例，从 DB 批量取号段缓存在内存里发放。
type Allocator struct {
	db     *DB
	bizTag string
	step   int64

	mu      sync.Mutex
	current int64 // 下一个要发出去的 ID
	end     int64 // 当前号段的上界（闭区间）
}

func newAllocator(db *DB, bizTag string, step int64) *Allocator {
	return &Allocator{db: db, bizTag: bizTag, step: step, end: -1}
}

// Next 返回下一个 ID。绝大多数调用都在内存里完成，只有号段发完才访问数据库。
func (a *Allocator) Next() int64 {
	a.mu.Lock()
	defer a.mu.Unlock()

	// 号段发完（end 初始为 -1，表示尚未取号），从数据库批量取新的一段。
	if a.current > a.end {
		start, end := a.db.allocSegment(a.bizTag, a.step)
		a.current = start
		a.end = end
	}

	id := a.current
	a.current++
	return id
}

func main() {
	db := newDB()

	// 用户 ID 发放，每段 1000 个。
	userAlloc := newAllocator(db, "user", 1000)

	fmt.Println("== 批量取号演示 ==")
	fmt.Println("  步长(step)=1000，顺序发前 3 个 ID:")
	for i := 0; i < 3; i++ {
		fmt.Printf("    ID = %d\n", userAlloc.Next())
	}
	fmt.Printf("  发了 3 个 ID，命中数据库次数 = %d（1000 个号一次取完）\n", db.allocCnt.Load())

	// 不同业务各自独立序号。
	fmt.Println("\n== 多业务独立 ==")
	orderAlloc := newAllocator(db, "order", 1000)
	fmt.Printf("  user  继续 = %d\n", userAlloc.Next())
	fmt.Printf("  order 起始 = %d（独立序号）\n", orderAlloc.Next())

	// 并发生成 + 唯一性校验，并统计命中数据库的次数。
	fmt.Println("\n== 并发唯一性校验 ==")
	db.allocCnt.Store(0) // 清零，重新统计
	const total = 10_000
	const workers = 8

	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		seen = make(map[int64]struct{}, total)
		dup  int
	)
	task := newAllocator(db, "task", 1000)
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < total/workers; i++ {
				id := task.Next()
				mu.Lock()
				if _, ok := seen[id]; ok {
					dup++
				} else {
					seen[id] = struct{}{}
				}
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	fmt.Printf("  %d 个 goroutine 并发生成 %d 个 ID，重复数 = %d\n", workers, total, dup)
	fmt.Printf("  命中数据库次数 = %d（共 %d 个 ID、步长 1000，理论应为 %d）\n",
		db.allocCnt.Load(), total, total/1000)
}