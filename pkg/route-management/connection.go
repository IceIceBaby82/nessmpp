package routemanagement

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

// ConnectionState represents the state of a route connection
type ConnectionState string

const (
	// ConnectionStateDisconnected indicates the connection is not established
	ConnectionStateDisconnected ConnectionState = "DISCONNECTED"

	// ConnectionStateConnecting indicates a connection attempt is in progress
	ConnectionStateConnecting ConnectionState = "CONNECTING"

	// ConnectionStateConnected indicates the connection is established
	ConnectionStateConnected ConnectionState = "CONNECTED"

	// ConnectionStateBinding indicates SMPP bind is in progress
	ConnectionStateBinding ConnectionState = "BINDING"

	// ConnectionStateBound indicates SMPP bind is complete
	ConnectionStateBound ConnectionState = "BOUND"

	// ConnectionStateError indicates the connection is in error state
	ConnectionStateError ConnectionState = "ERROR"
)

// ConnectionPool represents a pool of connections for a route
type ConnectionPool struct {
	route       *Route
	connections []*Connection
	mu          sync.RWMutex
}

// Connection represents a single connection to an SMPP route
type Connection struct {
	ID            string
	Route         *Route
	State         ConnectionState
	Conn          net.Conn
	LastActivity  time.Time
	LastError     error
	LastErrorTime time.Time
	RetryCount    int
	mu            sync.RWMutex
}

// ConnectionManager manages connections to SMPP routes
type ConnectionManager struct {
	manager      *RouteManager
	pools        map[string]*ConnectionPool
	maxRetries   int
	retryBackoff time.Duration
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	mu           sync.RWMutex
}

// NewConnectionManager creates a new connection manager
func NewConnectionManager(manager *RouteManager, maxRetries int, retryBackoff time.Duration) *ConnectionManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &ConnectionManager{
		manager:      manager,
		pools:        make(map[string]*ConnectionPool),
		maxRetries:   maxRetries,
		retryBackoff: retryBackoff,
		ctx:          ctx,
		cancel:       cancel,
	}
}

// Start begins connection management for all routes
func (cm *ConnectionManager) Start() {
	cm.wg.Add(1)
	go cm.manageConnections()
}

// Stop stops all connection management
func (cm *ConnectionManager) Stop() {
	cm.cancel()
	cm.wg.Wait()

	// Close all connections
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for _, pool := range cm.pools {
		pool.mu.Lock()
		for _, conn := range pool.connections {
			if conn.Conn != nil {
				conn.Conn.Close()
			}
		}
		pool.mu.Unlock()
	}
}

// manageConnections periodically checks and manages connections
func (cm *ConnectionManager) manageConnections() {
	defer cm.wg.Done()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-cm.ctx.Done():
			return
		case <-ticker.C:
			cm.checkConnections()
		}
	}
}

// checkConnections verifies and manages all connections
func (cm *ConnectionManager) checkConnections() {
	routes := cm.manager.GetActiveRoutes()

	for _, route := range routes {
		cm.ensureConnectionPool(route)
	}

	// Check each pool
	cm.mu.RLock()
	for routeID, pool := range cm.pools {
		pool.mu.RLock()
		activeConns := 0
		for _, conn := range pool.connections {
			if conn.State == ConnectionStateBound {
				activeConns++
			}
		}
		pool.mu.RUnlock()

		// If we need more connections, create them
		if activeConns < pool.route.MaxConnections {
			cm.wg.Add(1)
			go func(p *ConnectionPool) {
				defer cm.wg.Done()
				cm.createConnection(p)
			}(pool)
		}

		// Remove pools for inactive routes
		if !cm.manager.activeRoutes[routeID] {
			cm.removeConnectionPool(routeID)
		}
	}
	cm.mu.RUnlock()
}

// ensureConnectionPool creates a connection pool if it doesn't exist
func (cm *ConnectionManager) ensureConnectionPool(route *Route) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if _, exists := cm.pools[route.ID]; !exists {
		cm.pools[route.ID] = &ConnectionPool{
			route:       route,
			connections: make([]*Connection, 0),
		}
	}
}

// removeConnectionPool removes a connection pool and closes all connections
func (cm *ConnectionManager) removeConnectionPool(routeID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if pool, exists := cm.pools[routeID]; exists {
		pool.mu.Lock()
		for _, conn := range pool.connections {
			if conn.Conn != nil {
				conn.Conn.Close()
			}
		}
		pool.mu.Unlock()
		delete(cm.pools, routeID)
	}
}

// createConnection establishes a new connection for a route
func (cm *ConnectionManager) createConnection(pool *ConnectionPool) {
	conn := &Connection{
		ID:    fmt.Sprintf("%s-%d", pool.route.ID, time.Now().UnixNano()),
		Route: pool.route,
		State: ConnectionStateDisconnected,
	}

	pool.mu.Lock()
	pool.connections = append(pool.connections, conn)
	pool.mu.Unlock()

	for conn.RetryCount < cm.maxRetries {
		if err := cm.connect(conn); err != nil {
			conn.mu.Lock()
			conn.LastError = err
			conn.LastErrorTime = time.Now()
			conn.State = ConnectionStateError
			conn.RetryCount++
			conn.mu.Unlock()

			select {
			case <-cm.ctx.Done():
				return
			case <-time.After(cm.retryBackoff * time.Duration(conn.RetryCount)):
				continue
			}
		}
		break
	}
}

// connect establishes a TCP connection and performs SMPP bind
func (cm *ConnectionManager) connect(conn *Connection) error {
	conn.mu.Lock()
	conn.State = ConnectionStateConnecting
	conn.mu.Unlock()

	// Establish TCP connection
	address := fmt.Sprintf("%s:%d", conn.Route.Host, conn.Route.Port)
	tcpConn, err := net.DialTimeout("tcp", address, 10*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %v", address, err)
	}

	conn.mu.Lock()
	conn.Conn = tcpConn
	conn.State = ConnectionStateConnected
	conn.LastActivity = time.Now()
	conn.mu.Unlock()

	// TODO: Implement SMPP bind process
	// This will be implemented when we integrate with the SMPP client
	conn.mu.Lock()
	conn.State = ConnectionStateBound
	conn.mu.Unlock()

	return nil
}

// GetConnection returns an available connection from the pool
func (cm *ConnectionManager) GetConnection(routeID string) (*Connection, error) {
	cm.mu.RLock()
	pool, exists := cm.pools[routeID]
	cm.mu.RUnlock()

	if !exists {
		return nil, ErrRouteNotFound
	}

	pool.mu.RLock()
	defer pool.mu.RUnlock()

	// Find a bound connection with the least recent activity
	var selectedConn *Connection
	var oldestActivity time.Time

	for _, conn := range pool.connections {
		if conn.State == ConnectionStateBound {
			if selectedConn == nil || conn.LastActivity.Before(oldestActivity) {
				selectedConn = conn
				oldestActivity = conn.LastActivity
			}
		}
	}

	if selectedConn == nil {
		return nil, fmt.Errorf("no available connections for route %s", routeID)
	}

	selectedConn.mu.Lock()
	selectedConn.LastActivity = time.Now()
	selectedConn.mu.Unlock()

	return selectedConn, nil
}

// UpdateConnectionState updates the state of a connection
func (cm *ConnectionManager) UpdateConnectionState(conn *Connection, state ConnectionState) {
	conn.mu.Lock()
	defer conn.mu.Unlock()

	conn.State = state
	conn.LastActivity = time.Now()
}

// RecordConnectionError records an error for a connection
func (cm *ConnectionManager) RecordConnectionError(conn *Connection, err error) {
	conn.mu.Lock()
	defer conn.mu.Unlock()

	conn.LastError = err
	conn.LastErrorTime = time.Now()
	conn.State = ConnectionStateError
}
