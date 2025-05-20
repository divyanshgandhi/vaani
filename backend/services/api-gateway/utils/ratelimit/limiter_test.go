package ratelimit

import (
	"testing"
	"time"
)

func TestUserRateLimiter_Allow(t *testing.T) {
	// Create a rate limiter with 5 requests per second and 60 second reset
	limiter := NewUserRateLimiter(5, 1, 60*time.Second)
	userID := "test-user"

	// First 5 requests should be allowed
	for i := 0; i < 5; i++ {
		if !limiter.Allow(userID) {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// 6th request should be denied
	if limiter.Allow(userID) {
		t.Errorf("Request 6 should be denied")
	}

	// Wait for 1 second to refill 1 token
	time.Sleep(1 * time.Second)

	// Next request should be allowed after refill
	if !limiter.Allow(userID) {
		t.Errorf("Request after 1 second should be allowed")
	}

	// But immediately after, another should be denied
	if limiter.Allow(userID) {
		t.Errorf("Request after token used should be denied")
	}
}

func TestUserRateLimiter_TokensRemaining(t *testing.T) {
	limiter := NewUserRateLimiter(5, 1, 60*time.Second)
	userID := "test-user"

	// A new user should have full capacity
	if tokens := limiter.TokensRemaining(userID); tokens != 5 {
		t.Errorf("New user should have 5 tokens, got %d", tokens)
	}

	// After one request, should have 4 tokens
	limiter.Allow(userID)
	if tokens := limiter.TokensRemaining(userID); tokens != 4 {
		t.Errorf("User should have 4 tokens after one request, got %d", tokens)
	}

	// Use all remaining tokens
	for i := 0; i < 4; i++ {
		limiter.Allow(userID)
	}

	// Should have 0 tokens
	if tokens := limiter.TokensRemaining(userID); tokens != 0 {
		t.Errorf("User should have 0 tokens after using all, got %d", tokens)
	}
}

func TestUserRateLimiter_ResetUser(t *testing.T) {
	limiter := NewUserRateLimiter(5, 1, 60*time.Second)
	userID := "test-user"

	// Use some tokens
	for i := 0; i < 3; i++ {
		limiter.Allow(userID)
	}

	// Should have 2 tokens remaining
	if tokens := limiter.TokensRemaining(userID); tokens != 2 {
		t.Errorf("User should have 2 tokens, got %d", tokens)
	}

	// Reset the user
	limiter.ResetUser(userID)

	// Should have full capacity again
	if tokens := limiter.TokensRemaining(userID); tokens != 5 {
		t.Errorf("User should have 5 tokens after reset, got %d", tokens)
	}
}

func TestUserRateLimiter_ResetAll(t *testing.T) {
	limiter := NewUserRateLimiter(5, 1, 60*time.Second)
	
	// Use some tokens for multiple users
	limiter.Allow("user1")
	limiter.Allow("user1")
	limiter.Allow("user2")
	
	// Check tokens remaining
	if tokens := limiter.TokensRemaining("user1"); tokens != 3 {
		t.Errorf("User1 should have 3 tokens, got %d", tokens)
	}
	if tokens := limiter.TokensRemaining("user2"); tokens != 4 {
		t.Errorf("User2 should have 4 tokens, got %d", tokens)
	}
	
	// Reset all users
	limiter.ResetAll()
	
	// All users should have full capacity
	if tokens := limiter.TokensRemaining("user1"); tokens != 5 {
		t.Errorf("User1 should have 5 tokens after reset, got %d", tokens)
	}
	if tokens := limiter.TokensRemaining("user2"); tokens != 5 {
		t.Errorf("User2 should have 5 tokens after reset, got %d", tokens)
	}
}

func TestUserRateLimiter_FullReset(t *testing.T) {
	// Create a rate limiter with 5 capacity, 1 token per second refill,
	// but with a much shorter reset time for testing
	resetDuration := 2 * time.Second
	limiter := NewUserRateLimiter(5, 1, resetDuration)
	userID := "test-user"
	
	// Use all tokens
	for i := 0; i < 5; i++ {
		limiter.Allow(userID)
	}
	
	// Verify no tokens left
	if tokens := limiter.TokensRemaining(userID); tokens != 0 {
		t.Errorf("User should have 0 tokens, got %d", tokens)
	}
	
	// Wait for the full reset duration
	time.Sleep(resetDuration + 100*time.Millisecond) // Add a little buffer
	
	// Bucket should be reset to full capacity
	if tokens := limiter.TokensRemaining(userID); tokens != 5 {
		t.Errorf("User should have 5 tokens after full reset, got %d", tokens)
	}
}

func TestUserRateLimiter_Cleanup(t *testing.T) {
	// Create limiter with short cleanup interval
	limiter := NewUserRateLimiter(5, 1, 60*time.Second)
	limiter.cleanupInterval = 1 * time.Second // Override for testing
	
	// Use tokens for two users
	limiter.Allow("user1")
	limiter.Allow("user2")
	
	// Hack: modify last refill time for user1 to be in the past
	limiter.mutex.Lock()
	bucket := limiter.buckets["user1"]
	bucket.lastRefillTime = time.Now().Add(-121 * time.Second) // 2x resetAfter + a bit more
	limiter.mutex.Unlock()
	
	// Wait for cleanup to run
	time.Sleep(1200 * time.Millisecond)
	
	// Force a check that will trigger cleanup
	limiter.Allow("user3")
	
	// Check that user1 bucket was cleaned up, but user2 remains
	limiter.mutex.RLock()
	_, user1Exists := limiter.buckets["user1"]
	_, user2Exists := limiter.buckets["user2"]
	limiter.mutex.RUnlock()
	
	if user1Exists {
		t.Errorf("User1 bucket should have been cleaned up")
	}
	
	if !user2Exists {
		t.Errorf("User2 bucket should still exist")
	}
}

func BenchmarkUserRateLimiter_Allow(b *testing.B) {
	limiter := NewUserRateLimiter(1000, 100, 60*time.Second)
	userID := "benchmark-user"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limiter.Allow(userID)
	}
} 