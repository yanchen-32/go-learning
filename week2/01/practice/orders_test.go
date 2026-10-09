package main

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

const sampleJSON = `[
	{"id":"o-1","customer":"Alice","status":"paid","items":["book","pen"],"amount":3500},
	{"id":"o-2","customer":"Bob","status":"pending","items":["keyboard"],"amount":12000},
	{"id":"o-3","customer":"Alice","status":"cancelled","items":["mouse"],"amount":5000},
	{"id":"o-4","customer":"Alice","status":"paid","items":[],"amount":0}
]`

func sampleOrders(t *testing.T) []Order {
	t.Helper()
	orders, err := ParseOrders([]byte(sampleJSON))
	if err != nil {
		t.Fatal(err)
	}
	return orders
}

func TestParseOrders(t *testing.T) {
	orders := sampleOrders(t)
	want := []Order{
		{ID: "o-1", Customer: "Alice", Status: StatusPaid, Items: []string{"book", "pen"}, Amount: 3500},
		{ID: "o-2", Customer: "Bob", Status: StatusPending, Items: []string{"keyboard"}, Amount: 12000},
		{ID: "o-3", Customer: "Alice", Status: StatusCancelled, Items: []string{"mouse"}, Amount: 5000},
		{ID: "o-4", Customer: "Alice", Status: StatusPaid, Items: []string{}, Amount: 0},
	}
	if !reflect.DeepEqual(orders, want) {
		t.Fatalf("ParseOrders() = %+v, want %+v", orders, want)
	}
}

func TestParseOrdersInvalid(t *testing.T) {
	for _, tt := range []struct {
		name  string
		input string
		want  error
	}{
		{"syntax", `[`, nil},
		{"empty input", ``, nil},
		{"object instead of array", `{}`, nil},
		{"null", `null`, ErrExpectedArray},
		{"wrong amount type", `[{"id":"o-1","customer":"Alice","status":"paid","amount":"3500"}]`, nil},
		{"wrong items type", `[{"id":"o-1","customer":"Alice","status":"paid","items":"book","amount":0}]`, nil},
		{"empty ID", `[{"id":"","customer":"Alice","status":"paid","amount":0}]`, ErrEmptyID},
		{"blank ID", `[{"id":"  ","customer":"Alice","status":"paid","amount":0}]`, ErrEmptyID},
		{"missing ID", `[{"customer":"Alice","status":"paid","amount":0}]`, ErrEmptyID},
		{"empty customer", `[{"id":"o-1","customer":"","status":"paid","amount":0}]`, ErrEmptyCustomer},
		{"blank customer", `[{"id":"o-1","customer":" \t","status":"paid","amount":0}]`, ErrEmptyCustomer},
		{"missing customer", `[{"id":"o-1","status":"paid","amount":0}]`, ErrEmptyCustomer},
		{"invalid status", `[{"id":"o-1","customer":"Alice","status":"unknown","amount":0}]`, ErrInvalidStatus},
		{"missing status", `[{"id":"o-1","customer":"Alice","amount":0}]`, ErrInvalidStatus},
		{"negative amount", `[{"id":"o-1","customer":"Alice","status":"paid","amount":-1}]`, ErrNegativeAmount},
		{"duplicate ID", `[
			{"id":"o-1","customer":"Alice","status":"paid","amount":100},
			{"id":"o-1","customer":"Bob","status":"pending","amount":200}
		]`, ErrDuplicateID},
		{"valid then invalid", `[
			{"id":"o-1","customer":"Alice","status":"paid","amount":100},
			{"id":"o-2","customer":"Bob","status":"unknown","amount":200}
		]`, ErrInvalidStatus},
		{"valid then type error", `[
			{"id":"o-1","customer":"Alice","status":"paid","amount":100},
			{"id":"o-2","customer":"Bob","status":"paid","amount":false}
		]`, nil},
	} {
		t.Run(tt.name, func(t *testing.T) {
			orders, err := ParseOrders([]byte(tt.input))
			if err == nil || (tt.want != nil && !errors.Is(err, tt.want)) {
				t.Fatalf("error = %v, want %v", err, tt.want)
			}
			if orders != nil {
				t.Fatalf("error returned partial orders: %+v", orders)
			}
		})
	}
}

func TestParseErrorsPreserveCause(t *testing.T) {
	_, err := ParseOrders([]byte(`[`))
	var syntaxError *json.SyntaxError
	if !errors.As(err, &syntaxError) {
		t.Fatalf("error = %v, want wrapped JSON syntax error", err)
	}
	_, err = ParseOrders([]byte(`[{"amount":"wrong"}]`))
	var typeError *json.UnmarshalTypeError
	if !errors.As(err, &typeError) {
		t.Fatalf("error = %v, want wrapped JSON type error", err)
	}
}

func TestSummarizeAndJSON(t *testing.T) {
	stats := Summarize(sampleOrders(t))
	want := Statistics{
		StatusCounts: map[string]int{StatusPending: 1, StatusPaid: 2, StatusCancelled: 1},
		TotalAmount:  15500,
	}
	if !reflect.DeepEqual(stats, want) {
		t.Fatalf("Summarize() = %+v, want %+v", stats, want)
	}
	output, err := EncodeStatistics(stats)
	if err != nil {
		t.Fatal(err)
	}
	// 解码后比较字段，不比较 JSON 字符串中的键顺序。
	var decoded struct {
		StatusCounts map[string]int `json:"status_counts"`
		TotalAmount  int64          `json:"total_amount"`
	}
	if err := json.Unmarshal(output, &decoded); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded.StatusCounts, want.StatusCounts) || decoded.TotalAmount != want.TotalAmount {
		t.Fatalf("decoded statistics = %+v, want %+v", decoded, want)
	}
}

func orderIDs(orders []Order) []string {
	ids := make([]string, 0, len(orders))
	for _, order := range orders {
		ids = append(ids, order.ID)
	}
	return ids
}

func TestFilterByCustomer(t *testing.T) {
	orders := sampleOrders(t)
	filtered := FilterByCustomer(orders, "Alice")
	if got := orderIDs(filtered); !reflect.DeepEqual(got, []string{"o-1", "o-3", "o-4"}) {
		t.Fatalf("Alice order IDs = %v", got)
	}
	filtered[0].Customer = "changed"
	if orders[0].Customer != "Alice" {
		t.Fatal("filtered slice changed original order fields")
	}
	if got := FilterByCustomer(orders, "missing"); len(got) != 0 {
		t.Fatalf("missing customer matches = %+v, want empty", got)
	}
}

func TestInterfaceFilters(t *testing.T) {
	for _, tt := range []struct {
		name   string
		filter OrderFilter
		want   []string
	}{
		{"paid", StatusFilter{Status: StatusPaid}, []string{"o-1", "o-4"}},
		{"minimum inclusive", MinAmountFilter{Minimum: 5000}, []string{"o-2", "o-3"}},
		{"no matching status", StatusFilter{Status: "unknown"}, []string{}},
		{"no matching amount", MinAmountFilter{Minimum: 12001}, []string{}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			orders := sampleOrders(t)
			filtered := FilterOrders(orders, tt.filter)
			if got := orderIDs(filtered); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("filtered IDs = %v, want %v", got, tt.want)
			}
			if len(filtered) > 0 {
				id := filtered[0].ID
				filtered[0].ID = "changed"
				for _, order := range orders {
					if order.ID == "changed" {
						t.Fatalf("filter result for %q shares outer array", id)
					}
				}
			}
		})
	}
}

func TestEmptyOrders(t *testing.T) {
	orders, err := ParseOrders([]byte(`[]`))
	if err != nil || orders == nil || len(orders) != 0 {
		t.Fatalf("empty array = %+v, %v", orders, err)
	}
	want := Statistics{StatusCounts: map[string]int{StatusPending: 0, StatusPaid: 0, StatusCancelled: 0}}
	if got := Summarize(orders); !reflect.DeepEqual(got, want) {
		t.Fatalf("empty statistics = %+v, want %+v", got, want)
	}
	if len(FilterByCustomer(nil, "Alice")) != 0 || len(FilterOrders(nil, StatusFilter{Status: StatusPaid})) != 0 {
		t.Fatal("empty input should have no matches")
	}
}
