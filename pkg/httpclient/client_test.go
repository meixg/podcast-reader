package httpclient

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewRetryableClient(t *testing.T) {
	client := NewRetryableClient(10*time.Second, 3, 1*time.Second)

	if client == nil {
		t.Fatal("NewRetryableClient returned nil")
	}

	if client.client == nil {
		t.Error("client.client is nil")
	}

	if client.maxRetries != 3 {
		t.Errorf("maxRetries = %d, want 3", client.maxRetries)
	}

	if client.retryDelay != 1*time.Second {
		t.Errorf("retryDelay = %v, want 1s", client.retryDelay)
	}

	if client.client.Timeout != 10*time.Second {
		t.Errorf("timeout = %v, want 10s", client.client.Timeout)
	}
}

func TestRetryableClient_Do_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))
	defer server.Close()

	client := NewRetryableClient(5*time.Second, 3, 10*time.Millisecond)
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	ctx := context.Background()
	resp, err := client.Do(ctx, req)

	if err != nil {
		t.Errorf("Do() error = %v", err)
	}

	if resp == nil {
		t.Fatal("resp is nil")
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	resp.Body.Close()

	if string(body) != "success" {
		t.Errorf("body = %q, want %q", string(body), "success")
	}
}

func TestRetryableClient_Do_ClientErrorNoRetry(t *testing.T) {
	attemptCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	}))
	defer server.Close()

	client := NewRetryableClient(5*time.Second, 3, 10*time.Millisecond)
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	ctx := context.Background()
	resp, err := client.Do(ctx, req)

	if err != nil {
		t.Errorf("Do() error = %v", err)
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusNotFound)
	}

	if attemptCount != 1 {
		t.Errorf("attemptCount = %d, want 1 (should not retry 4xx errors)", attemptCount)
	}
}

func TestRetryableClient_Do_ServerErrorWithRetry(t *testing.T) {
	attemptCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		if attemptCount < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success after retry"))
		}
	}))
	defer server.Close()

	client := NewRetryableClient(5*time.Second, 3, 10*time.Millisecond)
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	ctx := context.Background()
	resp, err := client.Do(ctx, req)

	if err != nil {
		t.Errorf("Do() error = %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	if attemptCount != 3 {
		t.Errorf("attemptCount = %d, want 3", attemptCount)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	resp.Body.Close()

	if string(body) != "success after retry" {
		t.Errorf("body = %q, want %q", string(body), "success after retry")
	}
}

func TestRetryableClient_Do_AllRetriesExhausted(t *testing.T) {
	attemptCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("service unavailable"))
	}))
	defer server.Close()

	client := NewRetryableClient(5*time.Second, 3, 10*time.Millisecond)
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	ctx := context.Background()
	resp, err := client.Do(ctx, req)

	if err == nil {
		t.Error("Do() should return error after all retries exhausted")
	}

	if resp != nil {
		t.Error("resp should be nil when all retries fail")
	}

	if !strings.Contains(err.Error(), "all 3 retry attempts failed") {
		t.Errorf("error message should mention retry exhaustion, got: %v", err)
	}

	// Should attempt 1 initial + 3 retries = 4 total
	if attemptCount != 4 {
		t.Errorf("attemptCount = %d, want 4", attemptCount)
	}
}

func TestRetryableClient_Do_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewRetryableClient(5*time.Second, 3, 10*time.Millisecond)
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	// Cancel immediately
	cancel()

	resp, err := client.Do(ctx, req)

	if err == nil {
		t.Error("Do() should return error when context is cancelled")
	}

	if resp != nil {
		t.Error("resp should be nil when context is cancelled")
	}

	if !errors.Is(err, context.Canceled) {
		t.Errorf("error should be context.Canceled, got: %v", err)
	}
}

func TestRetryableClient_Do_NetworkError(t *testing.T) {
	// Use a closed server to simulate network error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	serverURL := server.URL
	server.Close() // Close immediately to cause connection error

	client := NewRetryableClient(1*time.Second, 2, 10*time.Millisecond)
	req, err := http.NewRequest("GET", serverURL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	ctx := context.Background()
	resp, err := client.Do(ctx, req)

	if err == nil {
		t.Error("Do() should return error on network failure")
	}

	if resp != nil {
		t.Error("resp should be nil on network failure")
	}

	if !strings.Contains(err.Error(), "all 2 retry attempts failed") {
		t.Errorf("error message should mention retry exhaustion, got: %v", err)
	}
}

func TestRetryableClient_Do_PostRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Method = %s, want POST", r.Method)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Failed to read request body: %v", err)
		}
		r.Body.Close()

		if string(body) != "test body" {
			t.Errorf("request body = %q, want %q", string(body), "test body")
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("created"))
	}))
	defer server.Close()

	client := NewRetryableClient(5*time.Second, 3, 10*time.Millisecond)
	req, err := http.NewRequest("POST", server.URL, strings.NewReader("test body"))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	ctx := context.Background()
	resp, err := client.Do(ctx, req)

	if err != nil {
		t.Errorf("Do() error = %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusCreated)
	}
}
