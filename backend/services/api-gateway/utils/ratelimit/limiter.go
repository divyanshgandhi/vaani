package ratelimit

import (
	"sync"
	"time"
)

// TokenBucket represents a token bucket rate limiter for a single user
type TokenBucket struct {
	capacity       int           // Maximum number of tokens
	tokens         int           // Current token count
	refillRate     int           // Tokens per second to refill
	lastRefillTime time.Time     // Last time tokens were refilled
	resetAfter     time.Duration // Duration after which to fully reset the bucket
	lastResetTime  time.Time     // Last time a full reset was performed
}

// UserRateLimiter maintains rate limiters for multiple users
type UserRateLimiter struct {
	buckets         map[string]*TokenBucket // Map of user IDs to token buckets
	capacity        int                     // Maximum tokens per bucket
	refillRate      int                     // Refill rate in tokens per second
	resetAfter      time.Duration           // Duration after which to fully reset buckets
	mutex           sync.RWMutex            // Mutex for thread-safe access
	cleanupInterval time.Duration           // How often to clean up inactive buckets
	lastCleanup     time.Time               // Last time cleanup was performed
}

// NewUserRateLimiter creates a new user rate limiter
func NewUserRateLimiter(capacity, refillRate int, resetAfter time.Duration) *UserRateLimiter {
	limiter := &UserRateLimiter{
		buckets:         make(map[string]*TokenBucket),
		capacity:        capacity,
		refillRate:      refillRate,
		resetAfter:      resetAfter,
		cleanupInterval: 10 * time.Minute, // Clean up inactive buckets every 10 minutes
		lastCleanup:     time.Now(),
	}
	return limiter
}

// Allow checks if a request from a user should be allowed
// Returns true if the request is allowed, false otherwise
func (l *UserRateLimiter) Allow(userID string) bool {
	now := time.Now()

	// Periodically clean up inactive buckets
	l.cleanupIfNeeded(now)

	l.mutex.Lock()
	defer l.mutex.Unlock()

	// Get or create bucket for this user
	bucket, exists := l.buckets[userID]
	if !exists {
		bucket = &TokenBucket{
			capacity:       l.capacity,
			tokens:         l.capacity, // Start with full capacity
			refillRate:     l.refillRate,
			lastRefillTime: now,
			resetAfter:     l.resetAfter,
			lastResetTime:  now,
		}
		l.buckets[userID] = bucket

		// Consume a token for this request
		bucket.tokens--
		return true
	}

	// Check if we should perform a full reset
	if now.Sub(bucket.lastResetTime) >= l.resetAfter {
		bucket.tokens = l.capacity
		bucket.lastRefillTime = now
		bucket.lastResetTime = now

		// Consume a token for this request
		bucket.tokens--
		return true
	}

	// Refill tokens based on time elapsed
	elapsedSeconds := int(now.Sub(bucket.lastRefillTime).Seconds())
	if elapsedSeconds > 0 {
		newTokens := elapsedSeconds * bucket.refillRate
		bucket.tokens = min(bucket.capacity, bucket.tokens+newTokens)
		bucket.lastRefillTime = now
	}

	// Check if request can be allowed
	if bucket.tokens > 0 {
		bucket.tokens--
		return true
	}

	return false
}

// TokensRemaining returns the number of tokens remaining for a user
func (l *UserRateLimiter) TokensRemaining(userID string) int {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	bucket, exists := l.buckets[userID]
	if !exists {
		return l.capacity
	}

	now := time.Now()

	// Check if we should perform a full reset
	if now.Sub(bucket.lastResetTime) >= l.resetAfter {
		return l.capacity
	}

	// Calculate tokens with refill
	elapsedSeconds := int(now.Sub(bucket.lastRefillTime).Seconds())
	if elapsedSeconds > 0 {
		newTokens := elapsedSeconds * bucket.refillRate
		return min(bucket.capacity, bucket.tokens+newTokens)
	}

	return bucket.tokens
}

// ResetUser resets the rate limit for a specific user
func (l *UserRateLimiter) ResetUser(userID string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	delete(l.buckets, userID)
}

// ResetAll resets all rate limits
func (l *UserRateLimiter) ResetAll() {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.buckets = make(map[string]*TokenBucket)
}

// cleanupIfNeeded removes inactive buckets to prevent memory leaks
func (l *UserRateLimiter) cleanupIfNeeded(now time.Time) {
	if now.Sub(l.lastCleanup) < l.cleanupInterval {
		return
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	// Reset lastCleanup time
	l.lastCleanup = now

	// Remove buckets that haven't been used recently
	for userID, bucket := range l.buckets {
		if now.Sub(bucket.lastRefillTime) > l.resetAfter*2 {
			delete(l.buckets, userID)
		}
	}
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
