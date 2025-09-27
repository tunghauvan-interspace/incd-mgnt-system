package idempotency

import (
	"crypto/md5"
	"fmt"
	"sync"
	"time"
)

// IdempotencyKey represents a unique key for request deduplication
type IdempotencyKey struct {
	Key       string
	ProcessedAt time.Time
	ExpiresAt   time.Time
}

// IdempotencyStore handles idempotency tracking
type IdempotencyStore interface {
	IsProcessed(key string) bool
	MarkProcessed(key string, ttl time.Duration) error
}

// MemoryIdempotencyStore is an in-memory implementation of IdempotencyStore
type MemoryIdempotencyStore struct {
	mu   sync.RWMutex
	keys map[string]*IdempotencyKey
}

// NewMemoryIdempotencyStore creates a new memory-based idempotency store
func NewMemoryIdempotencyStore() *MemoryIdempotencyStore {
	store := &MemoryIdempotencyStore{
		keys: make(map[string]*IdempotencyKey),
	}
	
	// Start cleanup goroutine
	go store.cleanupExpiredKeys()
	
	return store
}

// IsProcessed checks if a key has already been processed
func (s *MemoryIdempotencyStore) IsProcessed(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	idempotencyKey, exists := s.keys[key]
	if !exists {
		return false
	}
	
	// Check if key has expired
	if time.Now().After(idempotencyKey.ExpiresAt) {
		// Key has expired, remove it (will be cleaned up later)
		return false
	}
	
	return true
}

// MarkProcessed marks a key as processed with a TTL
func (s *MemoryIdempotencyStore) MarkProcessed(key string, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	now := time.Now()
	s.keys[key] = &IdempotencyKey{
		Key:       key,
		ProcessedAt: now,
		ExpiresAt:   now.Add(ttl),
	}
	
	return nil
}

// cleanupExpiredKeys periodically removes expired keys
func (s *MemoryIdempotencyStore) cleanupExpiredKeys() {
	ticker := time.NewTicker(5 * time.Minute) // Cleanup every 5 minutes
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			s.mu.Lock()
			now := time.Now()
			for key, idempotencyKey := range s.keys {
				if now.After(idempotencyKey.ExpiresAt) {
					delete(s.keys, key)
				}
			}
			s.mu.Unlock()
		}
	}
}

// GenerateKeyFromPayload generates an idempotency key from webhook payload
func GenerateKeyFromPayload(payload []byte) string {
	hash := md5.Sum(payload)
	return fmt.Sprintf("%x", hash)
}

// WebhookIdempotencyManager handles webhook idempotency
type WebhookIdempotencyManager struct {
	store IdempotencyStore
	ttl   time.Duration
}

// NewWebhookIdempotencyManager creates a new webhook idempotency manager
func NewWebhookIdempotencyManager(store IdempotencyStore, ttl time.Duration) *WebhookIdempotencyManager {
	return &WebhookIdempotencyManager{
		store: store,
		ttl:   ttl,
	}
}

// IsAlreadyProcessed checks if a webhook payload has already been processed
func (m *WebhookIdempotencyManager) IsAlreadyProcessed(payload []byte) bool {
	key := GenerateKeyFromPayload(payload)
	return m.store.IsProcessed(key)
}

// MarkAsProcessed marks a webhook payload as processed
func (m *WebhookIdempotencyManager) MarkAsProcessed(payload []byte) error {
	key := GenerateKeyFromPayload(payload)
	return m.store.MarkProcessed(key, m.ttl)
}