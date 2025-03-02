package smpp

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	// ErrPoolFull indicates that the connection pool is at maximum capacity
	ErrPoolFull = errors.New("connection pool is full")

	// ErrConnectionBusy indicates that the connection is currently in use
	ErrConnectionBusy = errors.New("connection is busy")

	// ErrInvalidConn indicates that the connection is not in the pool
	ErrInvalidConn = errors.New("invalid connection")
)

// ConnectionPool manages a pool of SMPP connections
type ConnectionPool struct {
	mu          sync.RWMutex
	sessions    map[*Session]struct{}
	maxSize     int
	idleTimeout time.Duration
	maxLifetime time.Duration
}

// NewConnectionPool creates a new connection pool
func NewConnectionPool(maxSize int, idleTimeout, maxLifetime time.Duration) *ConnectionPool {
	return &ConnectionPool{
		sessions:    make(map[*Session]struct{}),
		maxSize:     maxSize,
		idleTimeout: idleTimeout,
		maxLifetime: maxLifetime,
	}
}

// Add adds a session to the pool
func (p *ConnectionPool) Add(session *Session) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.sessions) >= p.maxSize {
		return ErrPoolFull
	}

	p.sessions[session] = struct{}{}
	return nil
}

// Remove removes a session from the pool
func (p *ConnectionPool) Remove(session *Session) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.sessions, session)
}

// Size returns the current number of sessions in the pool
func (p *ConnectionPool) Size() int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return len(p.sessions)
}

// CleanupIdleSessions removes idle sessions from the pool
func (p *ConnectionPool) CleanupIdleSessions() {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	for session := range p.sessions {
		if now.Sub(session.lastUsed) > p.idleTimeout || now.Sub(session.createdAt) > p.maxLifetime {
			session.Close()
			delete(p.sessions, session)
		}
	}
}

// Close closes all sessions in the pool
func (p *ConnectionPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for session := range p.sessions {
		session.Close()
		delete(p.sessions, session)
	}
}

// Acquire gets an available session from the pool
func (p *ConnectionPool) Acquire(ctx context.Context) (*Session, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Check if pool is closed
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Find an available connection
	for session := range p.sessions {
		if !session.IsBusy() {
			session.SetBusy(true)
			session.lastUsed = time.Now()
			return session, nil
		}
	}

	// Create new connection if pool not full
	if len(p.sessions) < p.maxSize {
		// Note: This is a placeholder. The actual session creation should be handled by the caller
		// since we need both the connection and server instance.
		return nil, errors.New("session creation must be handled by caller")
	}

	return nil, ErrPoolFull
}

// Release returns a session to the pool
func (p *ConnectionPool) Release(sess *Session) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Check if connection exists in pool
	_, exists := p.sessions[sess]
	if !exists {
		return ErrInvalidConn
	}

	if !sess.IsBusy() {
		return ErrConnectionBusy
	}

	sess.SetBusy(false)
	sess.lastUsed = time.Now()
	return nil
}

// Available returns the number of available sessions in the pool
func (p *ConnectionPool) Available() int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	available := 0
	for session := range p.sessions {
		if !session.IsBusy() {
			available++
		}
	}
	return available
}
