package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// 雪花算法（Snowflake）：Twitter 开源的分布式全局唯一 ID 生成算法。
// 把一个 64 位整数切分成若干段，每段表达不同含义：
//
//	0 - 0000000000 0000000000 0000000000 0000000000 0 - 00000 - 00000 - 000000000000
//	|1位符号位|      41位毫秒时间戳(相对epoch)      | 5位机房 | 5位机器 | 12位序列号 |
//
// 因为时间戳放在高位，所以生成的 ID 整体呈「趋势递增」，对 MySQL InnoDB
// 的 B+ 树索引非常友好，适合分库分表场景下的主键。

const (
	// epoch 是自定义的起始时间戳，ID 里存的「时间戳」是相对它的偏移，而非绝对时间。
	// 41 位相对时间戳大约可用 69 年。
	epoch int64 = 1704067200000 // 2024-01-01 00:00:00 UTC

	workerIDBits     = 5  // 机器 ID 占的位数
	datacenterIDBits = 5  // 机房 ID 占的位数
	sequenceBits     = 12 // 序列号占的位数（同一毫秒内能区分 2^12=4096 个）
)

const (
	maxWorkerID     = -1 ^ (-1 << workerIDBits)     // 31
	maxDatacenterID = -1 ^ (-1 << datacenterIDBits) // 31
	maxSequence     = -1 ^ (-1 << sequenceBits)     // 4095

	workerIDShift     = sequenceBits                                   // 12
	datacenterIDShift = sequenceBits + workerIDBits                    // 17
	timestampShift    = sequenceBits + workerIDBits + datacenterIDBits // 22
)

// Snowflake 是一个雪花算法发号器。
// 一个实例对应一个「机房+机器」组合，只要 (datacenterID, workerID) 不重复，
// 任意机器任意时刻生成的 ID 就全局唯一。
type Snowflake struct {
	datacenterID int64
	workerID     int64

	// state 把「上一次的毫秒时间戳」和「该毫秒内的序列号」打包成单个 64 位整数：
	//   高 41 位 = lastTimestamp，低 12 位 = sequence。
	// 把它们合成一个数，就能用一次 CAS(Compare-And-Swap) 原子地同时更新两者，
	// 从而做到「无锁」生成 ID，避免普通互斥锁在高并发下的竞争开销。
	state uint64
}

// NewSnowflake 创建一个发号器实例，datacenterID 与 workerID 需在合法范围内。
func NewSnowflake(datacenterID, workerID int64) (*Snowflake, error) {
	if datacenterID < 0 || datacenterID > maxDatacenterID {
		return nil, fmt.Errorf("datacenterID 必须在 [0, %d] 之间，得到 %d", maxDatacenterID, datacenterID)
	}
	if workerID < 0 || workerID > maxWorkerID {
		return nil, fmt.Errorf("workerID 必须在 [0, %d] 之间，得到 %d", maxWorkerID, workerID)
	}
	return &Snowflake{datacenterID: datacenterID, workerID: workerID}, nil
}

// NextID 生成下一个全局唯一 ID。
func (s *Snowflake) NextID() (int64, error) {
	for {
		// 1. 读出当前状态，拆出「上一次的毫秒时间戳」和「当前序列号」。
		state := atomic.LoadUint64(&s.state)
		lastTimestamp := int64(state >> sequenceBits)
		sequence := state & maxSequence

		now := time.Now().UnixMilli()

		// 2. 处理时钟回拨：当前时间比上次还早，说明系统时钟被回拨了。
		//    这是雪花算法唯一的隐患，这里选择直接拒绝生成，避免吐出重复 ID。
		if now < lastTimestamp {
			return 0, fmt.Errorf("时钟回拨 %d 毫秒，拒绝生成 ID", lastTimestamp-now)
		}

		var newTimestamp int64
		var newSequence uint64
		switch {
		case now == lastTimestamp:
			// 同一毫秒：序列号 +1。加满(maxSequence)后归 0，表明这一毫秒的号发完了，
			// 需要等到下一毫秒再发。
			newSequence = (sequence + 1) & maxSequence
			if newSequence == 0 {
				newTimestamp = waitNextMillis(lastTimestamp)
			} else {
				newTimestamp = lastTimestamp
			}
		default: // now > lastTimestamp，进入新的一毫秒，序列号从头开始
			newTimestamp = now
			newSequence = 0
		}

		// 3. 用 CAS 原子地抢占「新时间戳 + 新序列号」这个状态。
		newState := (uint64(newTimestamp) << sequenceBits) | newSequence
		if atomic.CompareAndSwapUint64(&s.state, state, newState) {
			// 抢占成功，把各段拼成一个 64 位 ID。
			return ((newTimestamp - epoch) << timestampShift) |
				(s.datacenterID << datacenterIDShift) |
				(s.workerID << workerIDShift) |
				int64(newSequence), nil
		}
		// 抢占失败：有别的 goroutine 抢先更新了 state，重新读、重新算。
	}
}

// waitNextMillis 忙等到下一毫秒。用于「当前毫秒序列号用尽」的情况。
func waitNextMillis(lastTimestamp int64) int64 {
	for {
		now := time.Now().UnixMilli()
		if now > lastTimestamp {
			return now
		}
	}
}

// Parse 把一个 ID 反解回它的各个组成部分，方便演示和排查。
func (s *Snowflake) Parse(id int64) (generatedAt time.Time, datacenterID, workerID, sequence int64) {
	sequence = id & maxSequence
	workerID = (id >> workerIDShift) & maxWorkerID
	datacenterID = (id >> datacenterIDShift) & maxDatacenterID
	generatedAt = time.UnixMilli((id >> timestampShift) + epoch)
	return
}

func main() {
	sf, err := NewSnowflake(1, 2) // 机房 1、机器 2
	if err != nil {
		panic(err)
	}

	// 顺序生成演示：验证 ID 单调递增。
	fmt.Println("== 顺序生成 ==")
	last := int64(0)
	for i := 0; i < 5; i++ {
		id, err := sf.NextID()
		if err != nil {
			panic(err)
		}
		fmt.Printf("  ID = %d（比上一个更大:%v）\n", id, id > last)
		last = id
	}

	// 拆解一个 ID，看它的结构。
	fmt.Println("\n== ID 结构拆解 ==")
	ts, dc, w, seq := sf.Parse(last)
	fmt.Printf("  生成时间: %s | 机房: %d | 机器: %d | 序列号: %d\n",
		ts.Format("2006-01-02 15:04:05.000"), dc, w, seq)

	// 并发生成 + 唯一性校验。
	fmt.Println("\n== 并发唯一性校验 ==")
	const total = 100_000
	const workers = 8

	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		seen = make(map[int64]struct{}, total)
		dup  int
	)
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < total/workers; i++ {
				id, err := sf.NextID()
				if err != nil {
					fmt.Println("生成失败:", err)
					return
				}
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

	fmt.Printf("  %d 个 goroutine 并发生成 %d 个 ID，重复数: %d\n", workers, total, dup)
}