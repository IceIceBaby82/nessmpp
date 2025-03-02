package retrymanagement

import (
	"context"
	"sync"
	"time"
)

// RetryScheduler handles scheduling and execution of retries
type RetryScheduler struct {
	manager       *RetryManager
	interval      time.Duration
	workerCount   int
	stopChan      chan struct{}
	workerWg      sync.WaitGroup
	retryCallback func(context.Context, *RetryAttempt) error
}

// NewRetryScheduler creates a new retry scheduler
func NewRetryScheduler(manager *RetryManager, interval time.Duration, workerCount int) *RetryScheduler {
	return &RetryScheduler{
		manager:     manager,
		interval:    interval,
		workerCount: workerCount,
		stopChan:    make(chan struct{}),
	}
}

// SetRetryCallback sets the callback function to be called for each retry attempt
func (rs *RetryScheduler) SetRetryCallback(callback func(context.Context, *RetryAttempt) error) {
	rs.retryCallback = callback
}

// Start starts the retry scheduler
func (rs *RetryScheduler) Start(ctx context.Context) {
	// Start worker pool
	for i := 0; i < rs.workerCount; i++ {
		rs.workerWg.Add(1)
		go rs.worker(ctx)
	}

	// Start scheduler
	ticker := time.NewTicker(rs.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-rs.stopChan:
			return
		case <-ticker.C:
			rs.scheduleRetries(ctx)
		}
	}
}

// Stop stops the retry scheduler
func (rs *RetryScheduler) Stop() {
	close(rs.stopChan)
	rs.workerWg.Wait()
}

// worker processes retry attempts
func (rs *RetryScheduler) worker(ctx context.Context) {
	defer rs.workerWg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case <-rs.stopChan:
			return
		default:
			attempts := rs.manager.GetPendingRetries()
			for _, attempt := range attempts {
				if rs.retryCallback != nil {
					if err := rs.retryCallback(ctx, attempt); err != nil {
						// Record failure and schedule next retry
						rs.manager.RecordFailure(ctx, attempt.MessageID, attempt.ProfileID, err, attempt.ErrorCategory)
					} else {
						// Remove successful attempt
						rs.manager.RemoveAttempt(attempt.MessageID)
					}
				}
			}
			time.Sleep(rs.interval)
		}
	}
}

// scheduleRetries schedules pending retries
func (rs *RetryScheduler) scheduleRetries(ctx context.Context) {
	// Cleanup expired attempts
	rs.manager.CleanupExpired()

	// Get and process pending retries
	attempts := rs.manager.GetPendingRetries()
	for _, attempt := range attempts {
		if rs.retryCallback != nil {
			if err := rs.retryCallback(ctx, attempt); err != nil {
				// Record failure and schedule next retry
				rs.manager.RecordFailure(ctx, attempt.MessageID, attempt.ProfileID, err, attempt.ErrorCategory)
			} else {
				// Remove successful attempt
				rs.manager.RemoveAttempt(attempt.MessageID)
			}
		}
	}
}
