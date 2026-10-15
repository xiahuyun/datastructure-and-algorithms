package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// ============================================================================
// 场景 1：GOMAXPROCS(1) + 纯 CPU 循环 → 直观对比：抢占开启 vs 抢占关闭
// ============================================================================
// 运行方式：
//
//	go run main.go                                        # 异步抢占开启（默认）
//	GODEBUG=asyncpreemptoff=1 go run main.go               # 异步抢占关闭
//
// 观察点：
//   - 抢占开启：B 能被定期抢占到，心跳间隔接近 500ms
//   - 抢占关闭：纯 for { i++ } 永不释放 CPU，B 彻底饿死！
//
// 原理：
//
//	Go 调度器是协作式的。goroutine 只有在以下情况才会主动让出 CPU：
//	  ① 函数调用（可能触发栈检查/调度检查点）
//	  ② channel / select 阻塞
//	  ③ 系统调用（网络、文件等）
//	  ④ time.Sleep / runtime.Gosched()
//
//	关键：哪怕是 time.Since() 这种标准库函数，底层也会调用 runtime 的
//	nanotime()，形成调度检查点。所以要做"永不让出 CPU"的循环，必须
//	循环体里只有纯算术操作，**连函数调用都不能有**。
//
//	Go 1.14+ 通过异步信号抢占（~10ms 发一次 SIGURG）硬打断死循环。
//	GODEBUG=asyncpreemptoff=1 可以关闭它，让我们看到 Go 1.14 以前的世界。
func demoStarvation() {
	fmt.Println("========== 场景 1：GOMAXPROCS(1) + 纯 CPU 循环 ==========")
	fmt.Println("运行命令对比：")
	fmt.Println("  默认（抢占开启）:   go run main.go")
	fmt.Println("  关闭异步抢占:       GODEBUG=asyncpreemptoff=1 go run main.go")
	fmt.Println()

	runtime.GOMAXPROCS(1)
	fmt.Printf("GOMAXPROCS = %d\n", runtime.GOMAXPROCS(0))

	// 预先测好：这台机器 100 亿次纯 i++ ≈ 2.66 秒
	// 用 200 亿次 → 约 5.3 秒，确保 A 跑的时间足够让 B 打印多条心跳
	const PURE_LOOP_ITERATIONS uint64 = 20_000_000_000

	// 用一个 flag 等 A 启动后再启动 B（确保 A 先拿到 CPU）
	aStarted := make(chan struct{})

	var wg sync.WaitGroup

	// goroutine A：真正的"永不释放 CPU"的纯循环
	// 循环体里 **只有 i++ 和 i < N 比较**，没有函数调用、没有 channel、
	// 没有 time.Now、没有 fmt.Println —— 这才是协作式调度下的"死循环怪物"
	wg.Add(1)
	go func() {
		defer wg.Done()

		fmt.Printf("[A] 启动！将要跑 %d 次纯算术循环（约 5 秒）\n",
			PURE_LOOP_ITERATIONS)
		close(aStarted) // 通知 B 可以启动了
		var i uint64 = 0
		// 注意：循环体里除了 i++ 和 i < N，什么都不能有
		for i < PURE_LOOP_ITERATIONS {
			i++
		}
		fmt.Printf("[A] 循环结束！i=%d\n", i)

	}()

	// 等 A 启动信号
	<-aStarted
	println("[B] 收到 A 启动信号，马上启动心跳 ticker")
	// goroutine B：每 500ms 打一条心跳，共 12 条（12 × 500ms = 6s > 5.3s）
	// 理论上 B 应该在 A 跑完之前就把 12 条心跳打完
	wg.Add(1)
	go func() {
		defer wg.Done()
		start := time.Now()
		fmt.Printf("[B] 启动 @ %v，期望每 500ms 打一条心跳，共 12 条\n",
			start.Format("15:04:05.000"))
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		count := 0
		for range ticker.C {
			count++
			elapsed := time.Since(start).Round(time.Millisecond)
			fmt.Printf("[B] 心跳 %2d @ %v (+%v)\n",
				count, time.Now().Format("15:04:05.000"), elapsed)
			if count >= 12 {
				fmt.Println("[B] 12 条心跳打完，退出")
				return
			}
		}
	}()

	mainStart := time.Now()
	wg.Wait()
	fmt.Printf("\n总耗时: %v\n", time.Since(mainStart).Round(time.Millisecond))
	fmt.Println()

	fmt.Println("────────────────────── 输出解读 ──────────────────────")
	fmt.Println("理论预期：12 条心跳 × 500ms = 6 秒内全部打印完毕")
	fmt.Println()
	fmt.Println("✅ 异步抢占开启（默认）：")
	fmt.Println("   B 能被定期抢占到，心跳间隔接近 500ms，全部在 6s 内完成")
	fmt.Println()
	fmt.Println("❌ 异步抢占关闭（GODEBUG=asyncpreemptoff=1）：")
	fmt.Println("   B 一条心跳都出不来！直到 A 跑完 5 秒后，B 的 ticker 事件")
	fmt.Println("   被一次性回放（tick 事件被缓冲），瞬间打完 12 条心跳。")
	fmt.Println("   这就是 Go 1.14 以前的典型问题：CPU 密集循环彻底饿死他人。")
	fmt.Println()
	fmt.Println("💡 赶紧跑一遍看看吧：")
	fmt.Println("   GODEBUG=asyncpreemptoff=1 go run main.go\n")
}

// ============================================================================
// 场景 2：GOMAXPROCS(1) + 通过 runtime.Gosched() 主动让出 → 协程能协作
// ============================================================================
func demoGosched() {
	fmt.Println("========== 场景 2：显式 runtime.Gosched() 让多个协程协作 ==========")
	runtime.GOMAXPROCS(1)

	var wg sync.WaitGroup

	// goroutine A：每次循环后主动 Gosched()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			println("[A] 迭代", i)
			runtime.Gosched() // 主动让出 CPU，让 B 有机会运行
		}
	}()

	// goroutine B：每次循环后也 Gosched()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			println("[B] 迭代", i)
			runtime.Gosched()
		}
	}()

	wg.Wait()
	fmt.Println()
}

// ============================================================================
// 场景 3：GOMAXPROCS(1) + channel 阻塞 → 调度器自然切换
// ============================================================================
// channel 发送/接收、定时器、网络 IO 等阻塞操作都会触发调度器，
// 让当前 goroutine 进入 waiting 状态，调度别的可运行 goroutine。
func demoChannelBlocking() {
	fmt.Println("========== 场景 3：channel 阻塞自然触发调度 ==========")
	runtime.GOMAXPROCS(1)

	ch := make(chan int)

	// goroutine A：生产者
	go func() {
		println("[A] 启动生产者 goroutine")
		for i := 0; i < 5; i++ {
			println("[A] 发送", i)
			ch <- i // 阻塞在这里，等待 B 接收
		}
		close(ch)
	}()

	// goroutine B：消费者
	println("[B] 启动消费者 goroutine")
	for v := range ch {
		println("[B] 收到", v)
	}
	println()
}

// ============================================================================
// 场景 4：GOMAXPROCS(1) 的真实用途 —— 节省内存 + 单线程安全模型
// ============================================================================
// - 每个 OS 线程默认占用 2MB 栈空间，GOMAXPROCS(1) 可以把线程数压到 1
// - 某些场景（比如需要串行处理、或为了简化并发模型）只允许一个线程
// - Go 的 select 循环 + GOMAXPROCS(1) 常常用来实现 "单线程事件循环"
func demoEventLoop() {
	fmt.Println("========== 场景 4：GOMAXPROCS(1) 作为单线程事件循环 ==========")
	runtime.GOMAXPROCS(1)

	ch1 := make(chan string)
	ch2 := make(chan string)
	done := make(chan struct{})

	// 多个生产者 goroutine（可以并行提交任务到 channel）
	go func() {
		for i := 0; i < 3; i++ {
			ch1 <- fmt.Sprintf("来自 ch1 的消息 #%d", i)
			time.Sleep(100 * time.Millisecond)
		}
	}()
	go func() {
		for i := 0; i < 3; i++ {
			ch2 <- fmt.Sprintf("来自 ch2 的消息 #%d", i)
			time.Sleep(80 * time.Millisecond)
		}
	}()
	go func() {
		time.Sleep(500 * time.Millisecond)
		close(done)
	}()

	// 单线程事件循环：所有消息在单个 OS 线程上串行处理
	// 处理逻辑不需要加锁，天然线程安全
loop:
	for {
		select {
		case msg := <-ch1:
			println("[事件循环] 收到", msg)
		case msg := <-ch2:
			println("[事件循环] 收到", msg)
		case <-done:
			break loop
		}
	}
	println()
}

func main() {
	fmt.Printf("默认 GOMAXPROCS = %d (CPU 核心数)\n", runtime.NumCPU())
	fmt.Println()

	/*
		demoStarvation()
		demoGosched()
		demoChannelBlocking()
	*/

	demoEventLoop()

	//demoPrint()
}

func demoPrint() {
	runtime.GOMAXPROCS(1)
	for i := 0; i < 5; i++ {
		go func(n int) { println(n) }(i)
	}
	time.Sleep(100 * time.Millisecond)
}
