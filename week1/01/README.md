# Go 基础 L3 尝试：最小 Task

这个目录是独立 Go module，不属于 DataAgentX。练习分成两份：

- `task/`：第一遍实现，包含较完整的校验、展示方法和测试。
- `secondpass/`：关闭资料后应从空文件重写的核心结构。
- `cmd/first/`、`cmd/second/`：两遍实现各自的可运行入口。

## 目录与文件说明

| 文件 | 作用 |
| --- | --- |
| `go.mod` | 声明这个练习的 module 路径和使用的 Go 版本。 |
| `README.md` | 说明练习目标、运行方式、核心知识点和第二遍重写记录。 |
| `task/task.go` | 第一遍的完整实现，定义 `Task`、constructor、value receiver、pointer receiver 和错误值。 |
| `task/task_test.go` | 测试第一遍的创建、重命名、完成状态及错误分支。 |
| `cmd/first/main.go` | 第一遍的命令行入口，演示调用 package 并在调用方处理 `error`。 |
| `secondpass/task.go` | 第二遍的空白练习文件；请关闭资料后在这里独立实现核心结构。 |
| `secondpass/task_test.go` | 第二遍的验收测试；重写时可以只保留此文件检查结果。 |
| `cmd/second/main.go` | 第二遍的命令行入口，验证重写版本能够被正常调用。 |

## 运行

```bash
go version
go env
go run ./cmd/first
go run ./cmd/second
go test ./...
```

## 建议练习流程（50～60 分钟）

1. 用 5 分钟确认环境并阅读题目。
2. 用 20～25 分钟自行完成第一遍；这一遍允许查 Go 官方资料。
3. 关闭资料，暂时移开第一遍代码，用 15～20 分钟从空文件写第二遍。
4. 用 10 分钟运行程序、补测试并填写下方记录。

第二遍如果需要重新查看第一遍或资料才能写完，就把等级记录为 L2；能独立写出核心结构、解释取舍并调通，才把这次尝试记为 L3。

## 五个问题

### module 与 package 的区别

module 是一组有共同版本边界的 Go package，由根目录的 `go.mod` 定义，并使用其中的 module path 作为导入路径前缀。package 是同一目录中一起编译的一组 Go 源文件，是组织和复用代码的基本单位。一个 module 通常包含多个 package。

### struct 的作用

`struct` 把一组相关字段组合成一个自定义类型。这里的 `Task` 把 `ID`、`Title`、`Done` 组合起来，使数据和操作这些数据的方法形成一个清晰的模型。

### value receiver 与 pointer receiver 的区别

value receiver 得到值的一份副本；方法内对 receiver 字段的修改不会改变原对象，适合只读、较小的值。pointer receiver 得到对象地址，可修改原对象，也避免复制整个值。为同一类型设计方法时通常保持 receiver 风格一致，除非有明确理由区分。

### 为什么修改状态的方法通常需要 pointer receiver

如果 `MarkDone` 使用 value receiver，它修改的只是副本的 `Done`，调用结束后原 `Task` 仍未完成。使用 `*Task` 才能让修改保留在调用者持有的对象中。

### error 在哪里创建、由谁处理

本例在 `task` package 中用 `errors.New` 创建可复用的错误值，在 constructor 和 `Rename` 中返回错误。调用方（`main` 或测试）决定如何处理：记录、退出、显示给用户，或用 `errors.Is` 断言错误类型。底层函数负责提供上下文明确的错误，调用方负责制定处理策略。

## 第二遍独立重写记录

已提交第二遍代码并完成三次验收；最终技术验收通过。

- 日期：2026-09-26
- 用时：未提供
- 是否关闭资料且未查看第一遍：否；验收过程中查看并使用了编译错误定位提示
- 完成内容：`struct`、constructor function、method、value receiver、pointer receiver 和 `error` 均已正确使用
- 曾卡住的位置：package/API 名称不匹配，以及 `IsDone()` 的返回类型；现已全部修正
- 格式检查：通过；`gofmt -d secondpass/task.go` 无输出
- 测试结果：通过；`go test ./...`、`go vet ./...` 和 `go run ./cmd/second` 均成功
- 自评：L2
- 判断依据：代码已达到最低完成标准，但第二遍经过两轮外部反馈才完成，不符合“关闭资料独立写出”的 L3 条件；本次如实保留为 L2，之后可另做一次全新盲写再挑战 L3

本次实现已完成。若要挑战 L3，请另建空白练习，在不查看本实现和外部提示的情况下独立完成，再运行同样的三项命令验收。
