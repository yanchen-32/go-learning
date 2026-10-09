package main

import "fmt"

func main() {
	input := []byte(`[
		{"id":"o-1","customer":"Alice","status":"paid","items":["book","pen"],"amount":3500},
		{"id":"o-2","customer":"Bob","status":"pending","items":["keyboard"],"amount":12000},
		{"id":"o-3","customer":"Alice","status":"cancelled","items":["mouse"],"amount":5000}
	]`)

	orders, err := ParseOrders(input)
	if err != nil {
		fmt.Println("订单解析或校验失败:", err)
		return
	}
	fmt.Printf("已解析并校验 %d 个订单\n", len(orders))
	fmt.Println("Alice 的订单:", FilterByCustomer(orders, "Alice"))
	fmt.Println("已支付订单:", FilterOrders(orders, StatusFilter{Status: StatusPaid}))
	fmt.Println("金额至少 5000 分的订单:", FilterOrders(orders, MinAmountFilter{Minimum: 5000}))

	output, err := EncodeStatistics(Summarize(orders))
	if err != nil {
		fmt.Println("统计结果编码失败:", err)
		return
	}
	fmt.Println("统计结果（金额单位：分）:")
	fmt.Println(string(output))
}
