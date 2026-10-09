# Week 2 / 01 练习说明与验收

## 运行

在 `week2/01` 执行：

```bash
go test ./...
go vet ./...
go run ./practice
```

四道主题练习的完整答案及测试在 `examples/` 中，例如运行 `go run ./examples/slice` 查看播放列表的添加、查询、统计、筛选与删除，运行 `go test ./examples/slice` 验证测试。综合练习的实现和测试在 `practice/` 中。

## 实现约定

- slice 使用 `[]Song`，`AddSong` 返回更新后的 slice，`FindSong` 返回歌曲副本及是否存在，`TotalDuration` 返回总秒数。`LongerThan` 严格筛选大于指定秒数的歌曲并返回独立 slice；删除不存在的 ID 返回错误。题目未要求 ID 唯一，重复时查找和删除均处理第一个匹配项。
- map 以 SKU 为 key，使用 `make` 初始化后通过 `AddProduct` 添加，拒绝重复 SKU、空白 SKU、空白名称和负数初始库存。`FindProduct` 使用 comma-ok 区分不存在和零库存。`AddStock` 与 `RemoveStock` 的数量必须大于零，操作失败保持原值。`RemoveProduct` 删除商品，`TotalStock` 统计总库存，`LowStock` 返回库存严格小于阈值的商品 slice，顺序不作保证。
- interface 的业务函数 `SendNotification` 接收 `Notifier`，控制台实现默认打印消息，可注入 `io.Writer` 测试输出；内存实现通过 pointer receiver 保存消息，两种实现都拒绝空白消息。`SendWelcome` 复用发送业务函数，发送错误返回给调用方。
- JSON 配置包含应用名称、环境、端口、功能开关 map 和主机 slice，均通过 tag 指定字段名。应用名称不能全为空白，端口为 1～65535，环境仅允许 `development`、`testing`、`production`。功能开关和主机列表允许省略。`DecodeConfig` 在解码后校验，失败返回零值配置；`EncodeConfig` 在编码前也校验。
- 订单输入必须是 JSON 数组，空数组合法，`null` 不合法。ID 和客户不能全为空白，状态必须合法，金额为非负整数分；ID 按原字符串检测重复。任何错误均返回 `nil` 和错误，不返回部分订单。
- 统计中始终包含三种状态的计数，取消订单也计数，但金额只累加 pending 与 paid。客户按名称精确匹配；最低金额筛选包含等于阈值的订单。筛选保留输入顺序。
- 筛选返回新的外层 slice。订单结构中的 `Items` 仍是浅拷贝，会共享原商品列表的底层数组；如需修改商品列表而不影响原订单，应额外复制 `Items`。

示例入口处理三个订单，状态计数各为 1，非取消订单总金额为 15500 分。

## 删除 slice 对底层数组的影响

slice 答案使用 `copy(songs[index:], songs[index+1:])` 将后续元素前移，清空尾元素后返回缩短的 slice。操作会修改原底层数组；如果另一个 slice 与其共享底层数组，也会观察到变化。删除最坏需要移动 O(n) 个元素。调用方必须接收返回的 slice。

常见的 `append(songs[:index], songs[index+1:]...)` 写法同样会复用并修改底层数组，但不会清空不再使用的尾元素。答案中清空尾元素是为了释放其持有的字符串引用。若要保留原数组，应先分配并复制，再删除。

## interface 的接收与返回

业务函数接收小 interface，让调用方自由提供控制台、内存或 fake 实现。构造函数通常返回具体类型，既能满足接口，又保留查看消息等具体能力。答案中的 `NewMemoryNotifier()` 返回 `*MemoryNotifier`，调用方既能发送通知，也能检查 `Messages`。只有在确实要隐藏实现等情况下，才需要返回 interface。

`var _ Notifier = (*MemoryNotifier)(nil)` 是编译期赋值检查：将一个类型为 `*MemoryNotifier` 的 nil 指针赋给接口，结果由空白标识符丢弃。它不会调用方法；缺少接口要求的方法时编译失败。

## 四道题的测试覆盖

| 目录 | 主要验收场景 |
| --- | --- |
| `examples/slice` | 添加、查找、总时长、首/中/末元素删除、不存在 ID、空列表、单元素、筛选边界及独立底层数组。 |
| `examples/map` | 重复 SKU、不存在键与零库存、增减库存、库存不足及失败不变、删除、总库存、低库存阈值；用集合比较避免依赖 map 顺序。 |
| `examples/interface` | 内存发送次数和内容、控制台输出、空消息不发送、底层发送错误传递；编译期检查两种实现。 |
| `examples/json` | 合法配置、缺失字段、类型错误、非法 JSON、端口上下界、三种合法环境、非法环境、业务校验失败和编码往返。 |

## 十二项验收问题

1. **slice 与数组有什么区别？** 数组长度是类型的一部分，例如 `[3]string`；slice 是包含底层数组指针、长度和容量的描述符，例如 `[]string`。复制数组会复制元素，复制 slice 描述符仍共享底层数组。
2. **len 与 cap 是什么？** `len` 是当前可访问的元素数；`cap` 是从 slice 起始位置到底层数组允许扩展边界的元素数。三下标切片可以限制这个容量边界。
3. **为什么接收 append 的返回值？** 追加会改变 slice 长度，容量不足时还会分配新的底层数组；必须使用返回的描述符，通常写 `s = append(s, value)`。
4. **nil slice 与空 slice 的异同？** 两者长度均可为零，都可以 range、append。nil slice 满足 `s == nil`；非 nil 的空 slice 不满足。标准 JSON 编码分别为 `null` 和 `[]`，`reflect.DeepEqual` 也会区分它们。
5. **如何判断 map key 是否存在？** 使用 `value, ok := m[key]`；`ok` 区分键不存在与值恰好为零。map 案例中库存为 0 的商品仍然存在。
6. **为什么不能依赖 map 遍历顺序？** Go 不保证 map 的迭代顺序。需要稳定展示时应取出 key 后排序；测试应比较键值集合，而不是遍历得到的顺序。
7. **interface 如何隐式实现？** 类型的方法集包含接口要求的全部方法，并且签名匹配，就满足该接口，无需 `implements` 声明。
8. **pointer receiver 如何影响 method set？** 类型 `T` 的方法集包含 receiver 为 `T` 的方法，`*T` 的方法集同时包含 receiver 为 `T` 和 `*T` 的方法。若 `Notify` 只有 pointer receiver，则只有 `*MemoryNotifier` 满足接口。对可寻址值直接调用时的自动取地址不等同于接口赋值时的方法集规则。
9. **为什么 JSON 字段通常必须导出？** `encoding/json` 只能正常读写导出字段；未导出字段即使写 tag 也会被忽略。导出字段也可以用 `json:"-"` 排除。
10. **Marshal 与 Unmarshal 做什么？** `Marshal` 将 Go 值编码为 JSON 字节；`Unmarshal` 将 JSON 字节解码到 Go 值，一般需要传入目标指针。
11. **JSON 语法校验与业务校验有什么区别？** 解码检查 JSON 结构及字段能否装入目标类型，业务校验检查空名称、端口范围、状态和重复 ID 等领域规则。`{"port":0}` 可成功解码为整数，但不代表端口有效。缺失字段可能留下零值，仍需业务校验。
12. **为什么金额使用整数？** 二进制浮点数无法精确表示许多十进制小数。以分为单位的整数可以精确保存和累加日常金额，综合题使用 `int64`；实际业务还应根据允许的金额上限约束整数范围。
