package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"
)

// RateLimitExceededResponse is the JSON response returned when rate limit is exceeded
type RateLimitExceededResponse struct {
	Error          string `json:"error"`
	RequestsPerSec int    `json:"requests_per_sec"`
	ResetAfter     int    `json:"reset_after_sec"`
	RetryAfter     int    `json:"retry_after_sec"`
}

// RateLimiter implements a token bucket rate limiting algorithm
type RateLimiter struct {
	tokens          map[string]int     // Number of tokens available for each user/IP
	lastRefill      map[string]time.Time // Time of last token refill for each user/IP
	tokensPerSecond int                // Rate at which tokens are replenished
	resetAfter      time.Duration      // Duration after which token count is reset
	mu              sync.Mutex         // Mutex to protect shared state
	logger          *zap.Logger        // Logger for rate limiting events
}

// newRateLimiter creates a new rate limiter with the specified settings
func newRateLimiter(tokensPerSecond int, resetAfter time.Duration, logger *zap.Logger) *RateLimiter {
	return &RateLimiter{
		tokens:          make(map[string]int),
		lastRefill:      make(map[string]time.Time),
		tokensPerSecond: tokensPerSecond,
		resetAfter:      resetAfter,
		logger:          logger,
	}
}

// getClientID returns a unique identifier for the client (currently IP address)
func getClientID(r *http.Request) string {
	// Try to get the real IP from X-Forwarded-For or similar headers first
	clientIP := r.Header.Get("X-Real-IP")
	if clientIP == "" {
		clientIP = r.Header.Get("X-Forwarded-For")
	}
	if clientIP == "" {
		clientIP = r.RemoteAddr
	}

	// Could be extended to use user ID from context if authenticated
	userID, err := ExtractUserID(r)
	if err == nil && userID != "" {
		return userID
	}

	return clientIP
}

// allow checks if a request should be allowed based on rate limits
func (rl *RateLimiter) allow(clientID string) (bool, int, time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Check if we should reset tokens
	lastRefill, exists := rl.lastRefill[clientID]
	if !exists || now.Sub(lastRefill) > rl.resetAfter {
		// Reset tokens for this client
		rl.tokens[clientID] = rl.tokensPerSecond
		rl.lastRefill[clientID] = now
	} else {
		// Calculate tokens to replenish based on time since last refill
		elapsed := now.Sub(lastRefill).Seconds()
		tokensToAdd := int(elapsed * float64(rl.tokensPerSecond))
		
		if tokensToAdd > 0 {
			// Add tokens up to the maximum
			rl.tokens[clientID] = min(rl.tokens[clientID]+tokensToAdd, rl.tokensPerSecond)
			rl.lastRefill[clientID] = now
		}
	}

	// Check if the client has available tokens
	if rl.tokens[clientID] > 0 {
		rl.tokens[clientID]--
		return true, rl.tokens[clientID], 0
	}

	// Calculate time until next token is available
	nextTokenTime := time.Second / time.Duration(rl.tokensPerSecond)
	return false, 0, nextTokenTime
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// RateLimitMiddleware creates a middleware that enforces rate limits
func RateLimitMiddleware(tokensPerSecond int, resetAfter time.Duration, logger *zap.Logger) func(http.Handler) http.Handler {
	rl := newRateLimiter(tokensPerSecond, resetAfter, logger)

	logger.Info("Rate limiting middleware initialized",
		zap.Int("tokens_per_second", tokensPerSecond),
		zap.Duration("reset_after", resetAfter),
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientID := getClientID(r)
			
			allowed, remainingTokens, retryAfter := rl.allow(clientID)
			
			// Set rate limiting headers
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(tokensPerSecond))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remainingTokens))
			
			if !allowed {
				w.Header().Set("Retry-After", strconv.Itoa(int(retryAfter.Seconds())))
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte("Rate limit exceeded. Please try again later."))
				
				logger.Info("Rate limit exceeded",
					zap.String("client_id", clientID),
					zap.Duration("retry_after", retryAfter),
				)
				return
			}
			
			next.ServeHTTP(w, r)
		})
	}
}

// getClientIP extracts the client IP address from a request
func getClientIP(r *http.Request) string {
	// Check for X-Forwarded-For header
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return forwarded
	}

	// Check for X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
} 