package main

import (
	"errors"
	"strings"
	"testing"
)

func TestNewWalletSuccess(t *testing.T) {
	w, err := NewWallet("  W001  ", "  Alice  ", 1050)
	if err != nil {
		t.Fatalf("NewWallet() error = %v", err)
	}
	if w == nil {
		t.Fatal("NewWallet() returned nil wallet")
	}
	if w.id != "W001" {
		t.Errorf("id = %q, want %q", w.id, "W001")
	}
	if w.owner != "Alice" {
		t.Errorf("owner = %q, want %q", w.owner, "Alice")
	}
	if w.balance != 1050 {
		t.Errorf("balance = %d, want %d", w.balance, 1050)
	}
}

func TestNewWalletEmptyID(t *testing.T) {
	w, err := NewWallet("   ", "Alice", 0)
	if err == nil {
		t.Fatal("NewWallet() expected error, got nil")
	}
	if !errors.Is(err, ErrEmptyID) {
		t.Errorf("error = %v, want ErrEmptyID", err)
	}
	if w != nil {
		t.Errorf("wallet = %#v, want nil", w)
	}
}

func TestNewWalletEmptyOwner(t *testing.T) {
	w, err := NewWallet("W001", "   ", 0)
	if err == nil {
		t.Fatal("NewWallet() expected error, got nil")
	}
	if !errors.Is(err, ErrEmptyOwner) {
		t.Errorf("error = %v, want ErrEmptyOwner", err)
	}
	if w != nil {
		t.Errorf("wallet = %#v, want nil", w)
	}
}

func TestNewWalletNegativeInitialBalance(t *testing.T) {
	w, err := NewWallet("W001", "Alice", -1)
	if err == nil {
		t.Fatal("NewWallet() expected error, got nil")
	}
	if !errors.Is(err, ErrNegativeInitialBalance) {
		t.Errorf("error = %v, want ErrNegativeInitialBalance", err)
	}
	if w != nil {
		t.Errorf("wallet = %#v, want nil", w)
	}
}

func TestDepositSuccess(t *testing.T) {
	w, err := NewWallet("W001", "Alice", 1000)
	if err != nil {
		t.Fatalf("NewWallet() error = %v", err)
	}

	if err := w.Deposit(250); err != nil {
		t.Fatalf("Deposit() error = %v", err)
	}
	if w.balance != 1250 {
		t.Errorf("balance = %d, want %d", w.balance, 1250)
	}
}

func TestDepositInvalidKeepsBalance(t *testing.T) {
	for _, amount := range []int{0, -1} {
		w, err := NewWallet("W001", "Alice", 1000)
		if err != nil {
			t.Fatalf("NewWallet() error = %v", err)
		}

		before := w.balance
		err = w.Deposit(amount)

		if err == nil {
			t.Errorf("Deposit(%d) expected error, got nil", amount)
		}
		if !errors.Is(err, ErrNonPositiveDeposit) {
			t.Errorf("Deposit(%d) error = %v, want ErrNonPositiveDeposit", amount, err)
		}
		if w.balance != before {
			t.Errorf("Deposit(%d) changed balance: got %d, want %d", amount, w.balance, before)
		}
	}
}

func TestWithdrawSuccess(t *testing.T) {
	w, err := NewWallet("W001", "Alice", 1000)
	if err != nil {
		t.Fatalf("NewWallet() error = %v", err)
	}

	if err := w.Withdraw(300); err != nil {
		t.Fatalf("Withdraw() error = %v", err)
	}
	if w.balance != 700 {
		t.Errorf("balance = %d, want %d", w.balance, 700)
	}
}

func TestWithdrawInsufficientKeepsBalance(t *testing.T) {
	w, err := NewWallet("W001", "Alice", 500)
	if err != nil {
		t.Fatalf("NewWallet() error = %v", err)
	}

	before := w.balance
	err = w.Withdraw(501)

	if err == nil {
		t.Fatal("Withdraw() expected error, got nil")
	}
	if !errors.Is(err, ErrInsufficientBalance) {
		t.Errorf("error = %v, want ErrInsufficientBalance", err)
	}
	if w.balance != before {
		t.Errorf("balance changed after failed withdraw: got %d, want %d", w.balance, before)
	}
}

func TestWithdrawInvalidKeepsBalance(t *testing.T) {
	for _, amount := range []int{0, -1} {
		w, err := NewWallet("W001", "Alice", 500)
		if err != nil {
			t.Fatalf("NewWallet() error = %v", err)
		}

		before := w.balance
		err = w.Withdraw(amount)

		if err == nil {
			t.Errorf("Withdraw(%d) expected error, got nil", amount)
		}
		if !errors.Is(err, ErrNonPositiveWithdraw) {
			t.Errorf("Withdraw(%d) error = %v, want ErrNonPositiveWithdraw", amount, err)
		}
		if w.balance != before {
			t.Errorf("Withdraw(%d) changed balance: got %d, want %d", amount, w.balance, before)
		}
	}
}

func TestIsZero(t *testing.T) {
	zeroWallet, err := NewWallet("W001", "Alice", 0)
	if err != nil {
		t.Fatalf("NewWallet() error = %v", err)
	}
	if !zeroWallet.IsZero() {
		t.Error("IsZero() = false, want true for zero balance")
	}

	nonZeroWallet, err := NewWallet("W002", "Bob", 1)
	if err != nil {
		t.Fatalf("NewWallet() error = %v", err)
	}
	if nonZeroWallet.IsZero() {
		t.Error("IsZero() = true, want false for non-zero balance")
	}
}

func TestInfoContainsFieldsAndDoesNotModifyWallet(t *testing.T) {
	w, err := NewWallet("W001", "Alice", 1050)
	if err != nil {
		t.Fatalf("NewWallet() error = %v", err)
	}

	beforeID, beforeOwner, beforeBalance := w.id, w.owner, w.balance

	info := w.Info()
	if !strings.Contains(info, "W001") ||
		!strings.Contains(info, "Alice") ||
		!strings.Contains(info, "1050") {
		t.Errorf("Info() = %q, want it to contain id, owner and balance", info)
	}

	if w.id != beforeID || w.owner != beforeOwner || w.balance != beforeBalance {
		t.Error("Info() modified wallet state")
	}
}
