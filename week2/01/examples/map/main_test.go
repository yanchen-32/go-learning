package main

import (
	"errors"
	"reflect"
	"testing"
)

func sampleInventory(t *testing.T) map[string]Product {
	t.Helper()
	inventory := make(map[string]Product)
	for _, product := range []Product{
		{SKU: "book", Name: "Book", Stock: 10},
		{SKU: "pen", Name: "Pen", Stock: 0},
		{SKU: "cup", Name: "Cup", Stock: 2},
	} {
		if err := AddProduct(inventory, product); err != nil {
			t.Fatal(err)
		}
	}
	return inventory
}

func TestAddAndFindProduct(t *testing.T) {
	inventory := sampleInventory(t)
	if product, ok := FindProduct(inventory, "pen"); !ok || product.Stock != 0 {
		t.Fatalf("zero-stock product = %+v, exists = %v", product, ok)
	}
	if _, ok := FindProduct(inventory, "missing"); ok {
		t.Fatal("missing key should not exist")
	}
	if err := AddProduct(inventory, Product{SKU: "book", Name: "Replacement", Stock: 99}); !errors.Is(err, ErrDuplicateSKU) {
		t.Fatalf("duplicate SKU error = %v", err)
	}
	if product, _ := FindProduct(inventory, "book"); product != (Product{SKU: "book", Name: "Book", Stock: 10}) {
		t.Fatalf("duplicate add changed original product: %+v", product)
	}
	for _, product := range []Product{{Name: "No SKU"}, {SKU: "no-name"}, {SKU: "negative", Name: "Invalid", Stock: -1}} {
		if err := AddProduct(inventory, product); !errors.Is(err, ErrInvalidProduct) {
			t.Fatalf("AddProduct(%+v) error = %v", product, err)
		}
	}
	if err := AddProduct(nil, Product{SKU: "book", Name: "Book"}); !errors.Is(err, ErrNilInventory) {
		t.Fatalf("nil map error = %v", err)
	}
}

func TestStockChanges(t *testing.T) {
	inventory := sampleInventory(t)
	if err := AddStock(inventory, "pen", 5); err != nil {
		t.Fatal(err)
	}
	if product, _ := FindProduct(inventory, "pen"); product.Stock != 5 {
		t.Fatalf("stock after add = %d, want 5", product.Stock)
	}
	if err := RemoveStock(inventory, "pen", 5); err != nil {
		t.Fatal(err)
	}
	if product, _ := FindProduct(inventory, "pen"); product.Stock != 0 {
		t.Fatalf("stock after removal = %d, want 0", product.Stock)
	}
}

func TestStockErrorsKeepInventory(t *testing.T) {
	for _, tt := range []struct {
		name   string
		change func(map[string]Product, string, int) error
		sku    string
		amount int
		want   error
	}{
		{"insufficient", RemoveStock, "book", 11, ErrInsufficientStock},
		{"zero stock", RemoveStock, "pen", 1, ErrInsufficientStock},
		{"add missing", AddStock, "missing", 1, ErrProductNotFound},
		{"remove missing", RemoveStock, "missing", 1, ErrProductNotFound},
		{"add zero", AddStock, "book", 0, ErrInvalidAmount},
		{"add negative", AddStock, "book", -1, ErrInvalidAmount},
		{"remove zero", RemoveStock, "book", 0, ErrInvalidAmount},
		{"remove negative", RemoveStock, "book", -1, ErrInvalidAmount},
	} {
		t.Run(tt.name, func(t *testing.T) {
			inventory := sampleInventory(t)
			before := sampleInventory(t)
			if err := tt.change(inventory, tt.sku, tt.amount); !errors.Is(err, tt.want) {
				t.Fatalf("error = %v, want %v", err, tt.want)
			}
			if !reflect.DeepEqual(inventory, before) {
				t.Fatalf("failed operation changed inventory: %v", inventory)
			}
		})
	}
}

func TestTotalsLowStockAndRemove(t *testing.T) {
	inventory := sampleInventory(t)
	if got := TotalStock(inventory); got != 12 {
		t.Fatalf("total = %d, want 12", got)
	}
	// 比较完整商品集合，不依赖 map 遍历顺序。
	got := make(map[string]Product)
	for _, product := range LowStock(inventory, 2) {
		got[product.SKU] = product
	}
	want := map[string]Product{"pen": {SKU: "pen", Name: "Pen", Stock: 0}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("LowStock(2) = %v, want %v", got, want)
	}
	if err := RemoveProduct(inventory, "book"); err != nil {
		t.Fatal(err)
	}
	if _, ok := FindProduct(inventory, "book"); ok || TotalStock(inventory) != 2 {
		t.Fatal("product was not removed")
	}
	before := map[string]Product{}
	for sku, product := range inventory {
		before[sku] = product
	}
	if err := RemoveProduct(inventory, "missing"); !errors.Is(err, ErrProductNotFound) {
		t.Fatalf("remove missing error = %v", err)
	}
	if !reflect.DeepEqual(inventory, before) {
		t.Fatal("failed removal changed inventory")
	}
	if TotalStock(nil) != 0 || len(LowStock(nil, 10)) != 0 {
		t.Fatal("nil inventory should have no stock")
	}
}
