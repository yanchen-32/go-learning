package main

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrEmptyID                = errors.New("wallet id cannot be empty")
	ErrEmptyOwner             = errors.New("wallet owner cannot be empty")
	ErrNegativeInitialBalance = errors.New("initial balance cannot be negative")
	ErrNonPositiveDeposit     = errors.New("deposit amount must be positive")
	ErrNonPositiveWithdraw    = errors.New("withdraw amount must be positive")
	ErrInsufficientBalance    = errors.New("insufficient balance")
)

type Wallet struct {
	id      string
	owner   string
	balance int
}

// NewWallet 创建钱包。失败时返回 nil 和具体错误。
func NewWallet(id, owner string, initialBalance int) (*Wallet, error) {
	id = strings.TrimSpace(id)
	owner = strings.TrimSpace(owner)

	if id == "" {
		return nil, ErrEmptyID
	}
	if owner == "" {
		return nil, ErrEmptyOwner
	}
	if initialBalance < 0 {
		return nil, ErrNegativeInitialBalance
	}

	return &Wallet{
		id:      id,
		owner:   owner,
		balance: initialBalance,
	}, nil
}

// Info 查询钱包信息。value receiver 表示只读，不修改原钱包。
func (w Wallet) Info() string {
	return fmt.Sprintf("钱包编号: %s, 所有者: %s, 当前余额: %d", w.id, w.owner, w.balance)
}

// Deposit 存入金额。pointer receiver 表示必须修改原钱包。
func (w *Wallet) Deposit(amount int) error {
	if amount <= 0 {
		return ErrNonPositiveDeposit
	}

	w.balance += amount
	return nil
}

// Withdraw 支出金额。pointer receiver 表示必须修改原钱包。
func (w *Wallet) Withdraw(amount int) error {
	if amount <= 0 {
		return ErrNonPositiveWithdraw
	}
	if w.balance < amount {
		return ErrInsufficientBalance
	}

	w.balance -= amount
	return nil
}

// IsZero 查询余额是否为零。value receiver 表示只读。
func (w Wallet) IsZero() bool {
	return w.balance == 0
}

func main() {
	w, err := NewWallet("  W001  ", "  Alice  ", 1000)
	if err != nil {
		fmt.Println("创建钱包失败:", err)
		return
	}

	fmt.Println("初始信息:", w.Info())

	if err := w.Deposit(500); err != nil {
		fmt.Println("存入失败:", err)
	} else {
		fmt.Println("存入 500 后:", w.Info())
	}

	if err := w.Withdraw(300); err != nil {
		fmt.Println("支出失败:", err)
	} else {
		fmt.Println("支出 300 后:", w.Info())
	}

	fmt.Println("最终信息:", w.Info())

	err = w.Withdraw(999999)
	if err != nil {
		fmt.Println("余额不足支出失败:", err)
		fmt.Println("失败后信息保持不变:", w.Info())
	}
}
