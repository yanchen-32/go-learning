package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	StatusPending   = "pending"
	StatusPaid      = "paid"
	StatusCancelled = "cancelled"
)

var (
	ErrExpectedArray  = errors.New("orders must be a JSON array")
	ErrEmptyID        = errors.New("order ID cannot be empty")
	ErrEmptyCustomer  = errors.New("customer cannot be empty")
	ErrInvalidStatus  = errors.New("invalid order status")
	ErrNegativeAmount = errors.New("order amount cannot be negative")
	ErrDuplicateID    = errors.New("duplicate order ID")
)

type Order struct {
	ID       string   `json:"id"`
	Customer string   `json:"customer"`
	Status   string   `json:"status"`
	Items    []string `json:"items"`
	Amount   int64    `json:"amount"` // 分，避免浮点数表示金额。
}

func (o Order) Validate() error {
	if strings.TrimSpace(o.ID) == "" {
		return ErrEmptyID
	}
	if strings.TrimSpace(o.Customer) == "" {
		return ErrEmptyCustomer
	}
	if o.Amount < 0 {
		return ErrNegativeAmount
	}
	switch o.Status {
	case StatusPending, StatusPaid, StatusCancelled:
		return nil
	default:
		return ErrInvalidStatus
	}
}

// ParseOrders 完整解码并校验后才返回；任一订单错误都返回 nil。
func ParseOrders(data []byte) ([]Order, error) {
	var orders []Order
	if err := json.Unmarshal(data, &orders); err != nil {
		return nil, fmt.Errorf("decode orders: %w", err)
	}
	// JSON null 可解码成 nil slice，但不满足题目的数组输入要求。
	if orders == nil {
		return nil, ErrExpectedArray
	}
	seen := make(map[string]bool)
	for i, order := range orders {
		if err := order.Validate(); err != nil {
			return nil, fmt.Errorf("order[%d]: %w", i, err)
		}
		if seen[order.ID] {
			return nil, fmt.Errorf("order[%d] ID %q: %w", i, order.ID, ErrDuplicateID)
		}
		seen[order.ID] = true
	}
	return orders, nil
}

type Statistics struct {
	StatusCounts map[string]int `json:"status_counts"`
	TotalAmount  int64          `json:"total_amount"`
}

// Summarize 接收已经通过校验的订单，金额不包含 cancelled 订单。
func Summarize(orders []Order) Statistics {
	stats := Statistics{StatusCounts: map[string]int{
		StatusPending: 0, StatusPaid: 0, StatusCancelled: 0,
	}}
	for _, order := range orders {
		stats.StatusCounts[order.Status]++
		if order.Status != StatusCancelled {
			stats.TotalAmount += order.Amount
		}
	}
	return stats
}

// OrderFilter 是调用方所需的最小筛选行为。
type OrderFilter interface {
	Match(Order) bool
}

type StatusFilter struct {
	Status string
}

func (f StatusFilter) Match(order Order) bool {
	return order.Status == f.Status
}

type MinAmountFilter struct {
	Minimum int64
}

func (f MinAmountFilter) Match(order Order) bool {
	return order.Amount >= f.Minimum
}

var _ OrderFilter = StatusFilter{}
var _ OrderFilter = MinAmountFilter{}

// FilterOrders 返回新的外层 slice；Order.Items 仍共享各自底层数组。
func FilterOrders(orders []Order, filter OrderFilter) []Order {
	result := make([]Order, 0)
	for _, order := range orders {
		if filter.Match(order) {
			result = append(result, order)
		}
	}
	return result
}

func FilterByCustomer(orders []Order, customer string) []Order {
	result := make([]Order, 0)
	for _, order := range orders {
		if order.Customer == customer {
			result = append(result, order)
		}
	}
	return result
}

func EncodeStatistics(stats Statistics) ([]byte, error) {
	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode statistics: %w", err)
	}
	return data, nil
}
