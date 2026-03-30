package payment

import (
	"errors"
	"fmt"
	"testing"
)

func TestAuthorise(t *testing.T) {
	result, _ := NewAuthorisationService(100).Authorise(10)
	expected := Authorisation{true, "Payment authorised"}
	if result != expected {
		t.Errorf("Authorise returned unexpected result: got %v want %v",
			result, expected)
	}
}

func TestFailOverCertainAmount(t *testing.T) {
	declineAmount := float32(10)
	result, _ := NewAuthorisationService(declineAmount).Authorise(100)
	expected := Authorisation{false, fmt.Sprintf("Payment declined: amount exceeds %.2f", declineAmount)}
	if result != expected {
		t.Errorf("Authorise returned unexpected result: got %v want %v",
			result, expected)
	}
}

func TestFailIfAmountIsZero(t *testing.T) {
	_, err := NewAuthorisationService(10).Authorise(0)
	if err == nil {
		t.Fatal("Authorise expected error for zero amount, got nil")
	}
	if !errors.Is(err, ErrInvalidPaymentAmount) {
		t.Errorf("Authorise expected ErrInvalidPaymentAmount for zero amount, got %v", err)
	}
}

func TestFailIfAmountNegative(t *testing.T) {
	_, err := NewAuthorisationService(10).Authorise(-1)
	if err == nil {
		t.Fatal("Authorise expected error for negative amount, got nil")
	}
	if !errors.Is(err, ErrInvalidPaymentAmount) {
		t.Errorf("Authorise expected ErrInvalidPaymentAmount for negative amount, got %v", err)
	}
}

// TestFailIfAmountNegativeLarge verifies that large negative amounts (e.g. from a
// catalogue item with a corrupted/negative price) are also rejected.
func TestFailIfAmountNegativeLarge(t *testing.T) {
	_, err := NewAuthorisationService(1000).Authorise(-500)
	if err == nil {
		t.Fatal("Authorise expected error for large negative amount, got nil")
	}
	if !errors.Is(err, ErrInvalidPaymentAmount) {
		t.Errorf("Authorise expected ErrInvalidPaymentAmount for large negative amount, got %v", err)
	}
}
