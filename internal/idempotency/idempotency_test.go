package idempotency

import (
	"testing"
	"time"
)

func TestMemoryIdempotencyStore(t *testing.T) {
	store := NewMemoryIdempotencyStore()

	// Test initial state
	if store.IsProcessed("test-key") {
		t.Error("Expected key to not be processed initially")
	}

	// Test marking as processed
	err := store.MarkProcessed("test-key", 1*time.Minute)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Test key is now processed
	if !store.IsProcessed("test-key") {
		t.Error("Expected key to be processed after marking")
	}

	// Test different key is not processed
	if store.IsProcessed("different-key") {
		t.Error("Expected different key to not be processed")
	}
}

func TestWebhookIdempotencyManager(t *testing.T) {
	store := NewMemoryIdempotencyStore()
	manager := NewWebhookIdempotencyManager(store, 5*time.Minute)

	payload1 := []byte(`{"test": "payload1"}`)
	payload2 := []byte(`{"test": "payload2"}`)

	// Test initial state
	if manager.IsAlreadyProcessed(payload1) {
		t.Error("Expected payload1 to not be processed initially")
	}

	// Test marking as processed
	err := manager.MarkAsProcessed(payload1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Test payload1 is now processed
	if !manager.IsAlreadyProcessed(payload1) {
		t.Error("Expected payload1 to be processed after marking")
	}

	// Test payload2 is not processed
	if manager.IsAlreadyProcessed(payload2) {
		t.Error("Expected payload2 to not be processed")
	}

	// Test same payload content generates same key
	payload1Copy := []byte(`{"test": "payload1"}`)
	if !manager.IsAlreadyProcessed(payload1Copy) {
		t.Error("Expected identical payload content to be considered processed")
	}
}

func TestGenerateKeyFromPayload(t *testing.T) {
	payload1 := []byte(`{"test": "payload"}`)
	payload2 := []byte(`{"test": "payload"}`)
	payload3 := []byte(`{"test": "different"}`)

	key1 := GenerateKeyFromPayload(payload1)
	key2 := GenerateKeyFromPayload(payload2)
	key3 := GenerateKeyFromPayload(payload3)

	// Same payload should generate same key
	if key1 != key2 {
		t.Error("Expected same payload to generate same key")
	}

	// Different payload should generate different key
	if key1 == key3 {
		t.Error("Expected different payload to generate different key")
	}

	// Keys should be non-empty
	if key1 == "" || key3 == "" {
		t.Error("Expected keys to be non-empty")
	}
}