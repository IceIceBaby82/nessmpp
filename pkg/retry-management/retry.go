package retrymanagement

import (
	"context"
	"math/rand"
	"sync"
	"time"
)

// RetryProfile represents a retry configuration profile
type RetryProfile struct {
	ID                string
	Name              string
	MaxAttempts       int
	InitialDelay      time.Duration
	MaxDelay          time.Duration
	BackoffMultiplier float64
	Jitter            float64
	ExpiryDuration    time.Duration
	ErrorCategories   []string
}

// RetryManager handles retry operations for failed messages
type RetryManager struct {
	mu       sync.RWMutex
	profiles map[string]*RetryProfile
	attempts map[string]*RetryAttempt
}

// RetryAttempt tracks retry attempts for a message
type RetryAttempt struct {
	MessageID     string
	ProfileID     string
	AttemptCount  int
	FirstAttempt  time.Time
	LastAttempt   time.Time
	NextAttempt   time.Time
	LastError     error
	ErrorCategory string
}

// NewRetryManager creates a new retry manager
func NewRetryManager() *RetryManager {
	return &RetryManager{
		profiles: make(map[string]*RetryProfile),
		attempts: make(map[string]*RetryAttempt),
	}
}

// AddProfile adds a new retry profile
func (rm *RetryManager) AddProfile(profile *RetryProfile) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.profiles[profile.ID] = profile
}

// RecordFailure records a message failure and schedules retry
func (rm *RetryManager) RecordFailure(ctx context.Context, messageID, profileID string, err error, errorCategory string) (*RetryAttempt, error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	profile, exists := rm.profiles[profileID]
	if !exists {
		return nil, ErrProfileNotFound
	}

	attempt, exists := rm.attempts[messageID]
	if !exists {
		attempt = &RetryAttempt{
			MessageID:     messageID,
			ProfileID:     profileID,
			FirstAttempt:  time.Now(),
			AttemptCount:  1,
			ErrorCategory: errorCategory,
		}
		rm.attempts[messageID] = attempt
	} else {
		attempt.AttemptCount++
		attempt.ErrorCategory = errorCategory
	}

	attempt.LastAttempt = time.Now()
	attempt.LastError = err

	// Check if max attempts reached
	if attempt.AttemptCount >= profile.MaxAttempts {
		return attempt, ErrMaxAttemptsReached
	}

	// Calculate next attempt time
	delay := profile.InitialDelay * time.Duration(attempt.AttemptCount)
	if profile.BackoffMultiplier > 1 {
		delay = time.Duration(float64(delay) * profile.BackoffMultiplier)
	}
	if delay > profile.MaxDelay {
		delay = profile.MaxDelay
	}

	// Add jitter
	if profile.Jitter > 0 {
		jitterRange := float64(delay) * profile.Jitter
		delay = time.Duration(float64(delay) + (jitterRange * (0.5 - rand.Float64())))
	}

	attempt.NextAttempt = time.Now().Add(delay)
	return attempt, nil
}

// GetPendingRetries returns all messages pending retry
func (rm *RetryManager) GetPendingRetries() []*RetryAttempt {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	now := time.Now()
	var pending []*RetryAttempt

	for _, attempt := range rm.attempts {
		if attempt.NextAttempt.Before(now) {
			pending = append(pending, attempt)
		}
	}

	return pending
}

// CleanupExpired removes expired retry attempts
func (rm *RetryManager) CleanupExpired() {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	now := time.Now()
	for messageID, attempt := range rm.attempts {
		profile := rm.profiles[attempt.ProfileID]
		if profile == nil {
			continue
		}

		if now.Sub(attempt.FirstAttempt) > profile.ExpiryDuration {
			delete(rm.attempts, messageID)
		}
	}
}

// RemoveAttempt removes a retry attempt
func (rm *RetryManager) RemoveAttempt(messageID string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	delete(rm.attempts, messageID)
}

// GetAttempt returns retry attempt information for a message
func (rm *RetryManager) GetAttempt(messageID string) (*RetryAttempt, bool) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	attempt, exists := rm.attempts[messageID]
	return attempt, exists
}
