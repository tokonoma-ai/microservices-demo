package payment

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/opentracing/opentracing-go"
)

func TestComponent(t *testing.T) {
	handler, logger := WireUp(float32(99.99), opentracing.GlobalTracer(), "test")

	ts := httptest.NewServer(handler)
	defer ts.Close()

	var request AuthoriseRequest
	request.Amount = 9.99
	requestBytes, err := json.Marshal(request)
	if err != nil {
		t.Fatal("ERROR", err)
	}

	res, err := http.Post(ts.URL+"/paymentAuth", "application/json", bytes.NewReader(requestBytes))
	if err != nil {
		t.Fatal("ERROR", err)
	}
	greeting, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal("ERROR", err)
	}
	var response Authorisation
	json.Unmarshal(greeting, &response)

	logger.Log("Authorised", response.Authorised)

	expected := true
	if response.Authorised != expected {
		t.Errorf("Authorise returned unexpected result: got %v expected %v",
			response.Authorised, expected)
	}
}

func TestNegativeAmountReturnsBadRequest(t *testing.T) {
	handler, _ := WireUp(float32(99.99), opentracing.GlobalTracer(), "test")

	ts := httptest.NewServer(handler)
	defer ts.Close()

	request := AuthoriseRequest{Amount: -5.00}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		t.Fatal("ERROR", err)
	}

	res, err := http.Post(ts.URL+"/paymentAuth", "application/json", bytes.NewReader(requestBytes))
	if err != nil {
		t.Fatal("ERROR", err)
	}
	res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request for negative amount, got %d", res.StatusCode)
	}
}

func TestZeroAmountReturnsBadRequest(t *testing.T) {
	handler, _ := WireUp(float32(99.99), opentracing.GlobalTracer(), "test")

	ts := httptest.NewServer(handler)
	defer ts.Close()

	request := AuthoriseRequest{Amount: 0}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		t.Fatal("ERROR", err)
	}

	res, err := http.Post(ts.URL+"/paymentAuth", "application/json", bytes.NewReader(requestBytes))
	if err != nil {
		t.Fatal("ERROR", err)
	}
	res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request for zero amount, got %d", res.StatusCode)
	}
}
