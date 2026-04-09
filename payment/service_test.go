package payment

import (
	"context"
	"fmt"
	"testing"
)

func TestAuthorise(t *testing.T) {
	result, _ := NewAuthorisationService(100).Authorise(context.Background(), 10)
	expected := Authorisation{true, "Payment authorised"}
	if result != expected {
		t.Errorf("Authorise returned unexpected result: got %v want %v",
			result, expected)
	}
}

func TestFailOverCertainAmount(t *testing.T) {
	declineAmount := float32(10)
	result, _ := NewAuthorisationService(declineAmount).Authorise(context.Background(), 100)
	expected := Authorisation{false, fmt.Sprintf("Payment declined: amount exceeds %.2f", declineAmount)}
	if result != expected {
		t.Errorf("Authorise returned unexpected result: got %v want %v",
			result, expected)
	}
}

func TestFailIfAmountIsZero(t *testing.T) {
	_, err := NewAuthorisationService(10).Authorise(context.Background(), 0)
	_, ok := err.(error)
	if !ok {
		t.Errorf("Authorise returned unexpected result: got %v want %v",
			err, "Zero payment")
	}
}

func TestFailIfAmountNegative(t *testing.T) {
	_, err := NewAuthorisationService(10).Authorise(context.Background(), -1)
	_, ok := err.(error)
	if !ok {
		t.Errorf("Authorise returned unexpected result: got %v want %v",
			err, "Negative payment")
	}
}
