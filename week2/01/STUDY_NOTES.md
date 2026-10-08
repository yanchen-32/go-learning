# 2026-10-08 学习总结：interface 与 JSON

## 今日学习内容

今天重点学习了 Go 的 `interface` 和标准库 `encoding/json`，并运行了本目录中的两个案例：

```bash
go run ./examples/interface
go run ./examples/json
```

## interface

### 核心理解

- interface 描述的是一组行为，而不是具体数据。
- Go 通过方法集隐式实现 interface，不需要写 `implements`。
- 调用方可以依赖一个小 interface，而不必依赖具体实现类型。
- 同一个 interface 可以有多种实现，例如控制台通知器和内存通知器。
- 依赖 interface 后，测试可以传入内存实现或 fake 实现，避免真正调用外部服务。

### Value receiver 与 pointer receiver

- value receiver 的方法既可以由值调用，也可以由指针调用。
- 如果 interface 要求的方法只实现在 pointer receiver 上，通常只有该类型的指针满足 interface。
- `MemoryNotifier.Notify` 需要向 slice 中追加消息，所以使用 pointer receiver，确保修改保留在原对象中。

### 编译期接口检查

案例使用以下形式确认类型满足 interface：

```go
var _ Notifier = ConsoleNotifier{}
var _ Notifier = (*MemoryNotifier)(nil)
```

如果实现缺少方法或方法签名不匹配，程序会在编译阶段报错。

### 设计原则

- interface 应尽量小，只包含调用方真正需要的行为。
- interface 通常定义在使用它的一方，而不是为了抽象而提前创建大接口。
- 接收 interface，便于替换实现；创建对象时通常仍返回明确的具体类型。

## JSON

### 核心理解

- `json.Unmarshal` 把 JSON 字节解析到 Go 值中。
- `json.Marshal` 或 `json.MarshalIndent` 把 Go 值编码为 JSON。
- `Unmarshal` 需要修改目标值，因此通常传入指针。
- 需要参与 JSON 编解码的 struct 字段通常必须导出，即字段名以大写字母开头。
- struct tag 用来控制 JSON 字段名称和编码行为。

示例：

```go
type Config struct {
	AppName  string          `json:"app_name"`
	Port     int             `json:"port"`
	Features map[string]bool `json:"features,omitempty"`
}
```

### 常用 JSON tag

- `json:"app_name"`：指定 JSON 字段名。
- `json:"name,omitempty"`：字段为零值时省略。
- `json:"-"`：不参与 JSON 编解码。

### 两层校验

JSON 处理需要区分两种错误：

1. 语法或类型错误：由 `json.Unmarshal` 返回，例如 JSON 缺少括号或字段类型不匹配。
2. 业务规则错误：JSON 虽然能成功解析，但内容不一定有效，例如应用名称为空或端口超出范围。

因此解析流程应当是：

```text
JSON 字节 → 语法解析 → struct → 业务校验 → 可用配置
```

### 错误包装

案例使用 `%w` 给底层错误增加上下文：

```go
fmt.Errorf("decode config: %w", err)
```

这样既能说明错误发生在哪个操作，也可以继续使用 `errors.Is` 或 `errors.As` 检查底层错误。

## 今日易错点

- interface 变量保存的是动态类型和动态值，不能只凭表面上的 `nil` 判断所有情况。
- pointer receiver 会影响一个类型是否满足 interface。
- JSON 字段未导出时，即使写了 tag，`encoding/json` 也不能正常设置该字段。
- JSON 解析成功不等于业务数据合法，必须单独做业务校验。
- map 编码为 JSON 对象时，不应依赖键的业务遍历顺序。
- 底层解析函数应返回错误，由调用方决定打印、退出或继续处理。

## 后续练习

- 独立实现至少两种 interface 实现，并使用内存实现编写测试。
- 测试 pointer receiver 对 interface method set 的影响。
- 为 JSON 配置解析补充合法数据、非法 JSON、错误字段类型和非法业务值测试。
- 继续完成 slice、map，并最终完成 `practice` 中的订单统计器综合练习。
