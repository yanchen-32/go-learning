# Week 2 / 01：slice、map、interface 与 JSON

本目录先提供四个彼此独立、可以直接运行的小案例，再保留一个不含答案的综合练习区。

## 目录说明

| 路径 | 内容 |
| --- | --- |
| `examples/slice/main.go` | slice 的创建、追加、遍历、筛选、复制和删除。 |
| `examples/map/main.go` | map 的创建、写入、comma-ok 查询、删除和稳定输出。 |
| `examples/interface/main.go` | interface 的隐式实现、不同实现和 pointer receiver。 |
| `examples/json/main.go` | JSON tag、解码、业务校验和编码。 |
| `STUDY_NOTES.md` | 2026-10-08 的 interface 与 JSON 学习总结。 |
| `practice/README.md` | 综合盲写题目；练习代码和测试由学习者自行创建。 |

## 运行案例

在本目录执行：

```bash
go run ./examples/slice
go run ./examples/map
go run ./examples/interface
go run ./examples/json
```

统一检查：

```bash
gofmt -w .
go test ./...
go vet ./...
```

建议依次阅读和运行四个案例，确认能解释每一处语法后，再打开 `practice/README.md` 开始综合练习。
