package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrInvalidAmount     = errors.New("amount must be positive")
	ErrDuplicateSKU      = errors.New("duplicate SKU")
	ErrInvalidProduct    = errors.New("SKU and name are required; stock cannot be negative")
	ErrNilInventory      = errors.New("inventory must be initialized with make")
)

type Product struct {
	SKU   string
	Name  string
	Stock int
}

func AddProduct(inventory map[string]Product, product Product) error {
	if inventory == nil {
		return ErrNilInventory
	}
	if strings.TrimSpace(product.SKU) == "" || strings.TrimSpace(product.Name) == "" || product.Stock < 0 {
		return ErrInvalidProduct
	}
	if _, exists := inventory[product.SKU]; exists {
		return ErrDuplicateSKU
	}
	inventory[product.SKU] = product
	return nil
}

func FindProduct(inventory map[string]Product, sku string) (Product, bool) {
	product, exists := inventory[sku]
	return product, exists
}

func AddStock(inventory map[string]Product, sku string, amount int) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	product, exists := inventory[sku]
	if !exists {
		return ErrProductNotFound
	}
	product.Stock += amount
	inventory[sku] = product
	return nil
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

func RemoveProduct(inventory map[string]Product, sku string) error {
	if _, exists := inventory[sku]; !exists {
		return ErrProductNotFound
	}
	delete(inventory, sku)
	return nil
}

func TotalStock(inventory map[string]Product) int {
	total := 0
	for _, product := range inventory {
		total += product.Stock
	}
	return total
}

// LowStock 返回库存严格低于 threshold 的商品，顺序不作保证。
func LowStock(inventory map[string]Product, threshold int) []Product {
	result := make([]Product, 0)
	for _, product := range inventory {
		if product.Stock < threshold {
			result = append(result, product)
		}
	}
	return result
}

func main() {
	inventory := make(map[string]Product)
	for _, product := range []Product{
		{SKU: "p-1", Name: "Keyboard", Stock: 8},
		{SKU: "p-2", Name: "Mouse", Stock: 0},
	} {
		if err := AddProduct(inventory, product); err != nil {
			fmt.Println("添加失败:", err)
			return
		}
	}
	if err := AddProduct(inventory, Product{SKU: "p-1", Name: "Duplicate", Stock: 99}); err != nil {
		fmt.Println("重复 SKU:", err)
	}

	// comma-ok 可以区分“键不存在”和“键存在但库存恰好为零”。
	product, exists := FindProduct(inventory, "p-2")
	fmt.Printf("p-2 存在=%t，库存=%d\n", exists, product.Stock)
	_, exists = FindProduct(inventory, "missing")
	fmt.Printf("missing 存在=%t\n", exists)
	if err := AddStock(inventory, "p-2", 2); err != nil {
		fmt.Println("增加失败:", err)
		return
	}

	if err := RemoveStock(inventory, "p-1", 3); err != nil {
		fmt.Println("扣减失败:", err)
		return
	}

	if err := RemoveStock(inventory, "p-1", 999); err != nil {
		fmt.Println("库存不足，保持原值:", err)
	}
	fmt.Println("总库存:", TotalStock(inventory))
	// LowStock 本身不保证顺序，展示时按 SKU 排序。
	low := LowStock(inventory, 3)
	sort.Slice(low, func(i, j int) bool { return low[i].SKU < low[j].SKU })
	fmt.Println("库存低于 3 的商品:", low)
	if err := RemoveProduct(inventory, "p-2"); err != nil {
		fmt.Println("删除失败:", err)
		return
	}

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
