package main

import (
	"errors"
	"fmt"
	"sort"
)

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrInvalidAmount     = errors.New("amount must be positive")
)

type Product struct {
	SKU   string
	Name  string
	Stock int
}

// RemoveStock 取出 map 中的值，修改后再写回。
func RemoveStock(inventory map[string]Product, sku string, amount int) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	product, exists := inventory[sku]
	if !exists {
		return ErrProductNotFound
	}
	if product.Stock < amount {
		return ErrInsufficientStock
	}

	product.Stock -= amount
	inventory[sku] = product
	return nil
}

func main() {
	inventory := make(map[string]Product)
	inventory["p-1"] = Product{SKU: "p-1", Name: "Keyboard", Stock: 8}
	inventory["p-2"] = Product{SKU: "p-2", Name: "Mouse", Stock: 0}

	// comma-ok 可以区分“键不存在”和“键存在但库存恰好为零”。
	product, exists := inventory["p-2"]
	fmt.Printf("p-2 存在=%t，库存=%d\n", exists, product.Stock)

	if err := RemoveStock(inventory, "p-1", 3); err != nil {
		fmt.Println("扣减失败:", err)
	}

	delete(inventory, "p-2")

	// map 遍历顺序不稳定；先排序 key，输出才稳定。
	keys := make([]string, 0, len(inventory))
	for sku := range inventory {
		keys = append(keys, sku)
	}
	sort.Strings(keys)

	for _, sku := range keys {
		fmt.Printf("%s: %+v\n", sku, inventory[sku])
	}
}
