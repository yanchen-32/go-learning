# Week 2 / 01：slice、map、interface 与 JSON

本目录包含四道主题练习的完整答案、测试和可运行示例，以及订单统计器综合练习及其测试。

## 目录说明

| 路径 | 内容 |
| --- | --- |
| `examples/slice/` | 播放列表：添加、删除、查找、总时长、时长筛选及边界测试。 |
| `examples/map/` | 库存索引：重复 SKU 校验、查询、增减库存、删除、总库存与低库存查询。 |
| `examples/interface/` | 通知服务：小接口、控制台与内存实现、发送业务函数及测试。 |
| `examples/json/` | 配置解析器：JSON tag、环境与端口校验、编解码及测试。 |
| `STUDY_NOTES.md` | 2026-10-08 的 interface 与 JSON 学习总结。 |
| `practice/` | 订单统计器题目、实现、测试及可运行入口。 |
| `EXERCISE_NOTES.md` | 练习约定、设计说明与十二项验收问题的回答。 |

## 运行案例

在本目录执行：

```bash
go run ./examples/slice
go run ./examples/map
go run ./examples/interface
go run ./examples/json
go run ./practice
```

统一检查：

```bash
gofmt -w .
go test ./...
go vet ./...
```

建议依次阅读和运行 `examples/` 中的四个案例，再阅读 `practice/` 的综合实现及测试。独立复习时可以先根据题目重写，再用测试验收。

每个主题的答案在对应目录的 `main.go`，测试在 `main_test.go`。例如只检查播放列表：`go test ./examples/slice`。
