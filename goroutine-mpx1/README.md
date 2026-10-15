# goroutine-mpx1 · GOMAXPROCS(1) 下的 Go 调度实验

用可运行的最小示例，观察 `runtime.GOMAXPROCS(1)`（只有一个 P / 逻辑处理器）时多个 goroutine 是如何被调度的，并配有可视化动画帮助理解 GMP 模型的运行过程。

## 目录内容

| 文件 | 说明 |
| --- | --- |
| `main.go` | 实验代码，包含 4 个场景函数 |
| `scheduler-animation.html` | 泳道时间轴风格的调度动画（OS 线程 M / main / A / B 四条泳道） |
| `gmp-animation.html` | GMP 盒子物理移动风格动画（Global Queue / runnext / local queue / running slot） |
| `go.mod` | 模块定义 |

## 运行方式

```bash
cd goroutine-mpx1

# 异步抢占开启（Go 1.14+ 默认行为）
go run main.go

# 关闭异步抢占，回到 Go 1.14 之前的"协作式调度"世界
GODEBUG=asyncpreemptoff=1 go run main.go
```

`main.go` 顶部 `main()` 里默认只启用一个场景，切换方式是把 `demoStarvation()`、`demoGosched()`、`demoChannelBlocking()`、`demoEventLoop()` 按需取消注释。

## 场景说明

| 场景 | 函数 | 观察重点 |
| --- | --- | --- |
| 抢占开启 vs 关闭 | `demoStarvation` | 纯 CPU 循环（无任何调度检查点）时，心跳 goroutine 是否被饿死 |
| 主动让出 | `demoGosched` | `runtime.Gosched()` 显式让出 CPU 后多协程如何交替 |
| channel 阻塞 | `demoChannelBlocking` | 无缓冲 channel 的收发阻塞如何自然触发调度 |
| 单线程事件循环 | `demoEventLoop` | `GOMAXPROCS(1)` + `select` 实现串行处理，无需加锁 |

## 核心结论

- Go 的调度是**协作式**的：goroutine 只在函数调用、channel/select 阻塞、系统调用、`time.Sleep`、`runtime.Gosched()` 这些**调度检查点**才会让出 CPU。
- `demoStarvation` 里 A 的循环体只有 `i++` 和比较，**不含任何函数调用**，因此在关闭异步抢占时会独占唯一的 P 长达数秒，B 的 ticker 心跳被彻底饿死（直到 A 结束才被一次性回放）。
- Go 1.14+ 引入了**异步信号抢占**（约每 10ms 发一次 `SIGURG`），能硬打断死循环，让 B 每隔约 500ms 正常打出心跳。
- `GODEBUG=asyncpreemptoff=1` 可关闭异步抢占，用来对比两种调度行为。

## 动画说明

两个 HTML 文件均可直接用浏览器打开，无需依赖或构建步骤：

- `scheduler-animation.html`：以时间轴形式展示每个 goroutine 在不同时刻的状态，并标注 `SIGURG` 抢占点。
- `gmp-animation.html`：以盒子移动的形式展示 G 在 `runnext`、local queue、Global Queue、running slot 之间的流转，支持「异步抢占开启 / 关闭」两种模式切换。

例如：

```bash
open gmp-animation.html
```
