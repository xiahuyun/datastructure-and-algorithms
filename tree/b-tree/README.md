# B-Tree (Go)

这个项目在 `tree/b-tree` 下实现了一个 B-Tree，使用教材常见的“`m` 阶”语义：

- 每个节点最多 `m` 个子节点
- 每个节点最多 `m-1` 个关键字
- 非根节点至少 `ceil(m/2)` 个子节点

## 已实现能力

- `New(order int)` / `NewWithFreeList(order, freelist)`
- `ReplaceOrInsert(item)`：插入或替换同 key 项
- `DebugString()`：按层打印树结构（便于学习和排错）
- `Get(item)`：查找
- `Delete(item)` / `DeleteMin()` / `DeleteMax()`：删除
- `Min()` / `Max()` / `Len()`
- `Ascend*` / `Descend*`：顺序、逆序与范围遍历
- `Clone()` / `Clear(addNodesToFreelist bool)`

## 运行测试

```bash
cd tree/b-tree
go test ./...
```

## 运行示例

```bash
cd tree/b-tree
go run ./example
```

## 逐步演示模式

```bash
cd tree/b-tree
go run ./example --step
```

会在每一步插入/删除后打印当前树结构，便于对照教材理解分裂、借位、合并过程。

## 插入实现说明

`ReplaceOrInsert` 采用“迭代下探 + 路径栈回溯分裂”的方式：

1. 从根节点一路下探到叶子，记录经过的父节点与子节点下标
2. 在叶子插入新 key
3. 若节点溢出（key 数超过 `m-1`），按中位数分裂并向父节点提升
4. 持续向上处理，直到根节点或不再溢出

## 删除实现说明

`Delete` / `DeleteMin` / `DeleteMax` 采用非递归 top-down 迭代方式：

1. 从根向下迭代
2. 下探前确保目标子节点不是最小关键字数（必要时借位或合并）
3. 到叶子后删除；若在内部节点命中，使用前驱/后继替换并继续向下删除
