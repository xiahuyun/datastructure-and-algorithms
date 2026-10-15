# id-generator/segment — 数据库号段模式（Leaf-segment）

分库分表下需要全局唯一 ID，除了雪花算法（`id-generator/snowflake`），另一条主流思路是**号段模式**，代表实现是美团 Leaf-segment 和滴滴 TinyID。

## 核心思想：批量取号

与其「每生成一个 ID 都访问一次数据库」，不如「一次取一批（号段），缓存在内存里慢慢发」。

数据库里维护一张表：

```sql
CREATE TABLE id_alloc (
  biz_tag  varchar(128) PRIMARY KEY, -- 业务标识
  max_id   bigint     NOT NULL,      -- 当前已分配的最大 ID
  step     int        NOT NULL       -- 每次取号的步长
);
```

取号段时只执行一条 SQL：

```sql
UPDATE id_alloc SET max_id = max_id + step WHERE biz_tag = ?
```

单条 `UPDATE` 在数据库里是原子的（行锁），所以并发安全天然有保证。取回的那段号缓存到本地内存，业务请求 ID 时直接从内存分配，发完再去取下一段。于是 DB 压力从「每请求 1 次」降到「每 step 次 1 次」。

## 运行方式

```bash
cd id-generator/segment
go run main.go
```

输出大致如下：

```
== 批量取号演示 ==
  步长(step)=1000，顺序发前 3 个 ID:
    ID = 1
    ID = 2
    ID = 3
  发了 3 个 ID，命中数据库次数 = 1（1000 个号一次取完）

== 多业务独立 ==
  user  继续 = 4
  order 起始 = 1（独立序号）

== 并发唯一性校验 ==
  8 个 goroutine 并发生成 10000 个 ID，重复数 = 0
  命中数据库次数 = 10（共 10000 个 ID、步长 1000，理论应为 10）
```

## 并发模型

- **数据库侧**：靠单条 `UPDATE` 的原子性（行锁）保证不同进程/机器取到不重复的号段。
- **本地侧**：号段在内存里发放，用一把 `sync.Mutex` 串行化（临界区极短，只是两条 int64 自增）。真正访问数据库的地方只有「号段发完」那一刻。

## 一个缺点：临界点卡顿

当号段用完、需要去数据库取新段的那一刻，所有请求都会被阻塞等待，导致那一下响应时间骤增（TP999 变差）。美团 Leaf 用**双 buffer** 解决：当前号段快发完时，就异步预取下一段到备用 buffer，切过去时无需等待。

## 与雪花算法的对比

| 维度 | 号段模式（segment） | 雪花算法（snowflake） |
| --- | --- | --- |
| 依赖 | 数据库（用于取号段） | 无外部组件 |
| ID 形态 | 连续整数 | 64 位时间戳+机器+序列号 |
| 有序性 | 严格单调递增 | 趋势递增（有间隙） |
| 单点风险 | 依赖 DB 可用性 | 依赖机器时钟（回拨） |
| 典型实现 | 美团 Leaf-segment、滴滴 TinyID | 百度 UidGenerator、美团 Leaf-snowflake |