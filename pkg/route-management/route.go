package routemanagement

import (
	"sort"
	"sync"
	"time"
)

// Route represents an SMPP route configuration
type Route struct {
	ID              string
	Name            string
	Priority        int
	Weight          int
	Host            string
	Port            int
	SystemID        string
	Password        string
	BindType        string
	Status          string
	MaxConnections  int
	RateLimit       int
	CostPerMessage  float64
	Enabled         bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
	LastHealthCheck time.Time
	HealthStatus    string
	Metrics         *RouteMetrics
}

// RouteMetrics tracks performance metrics for a route
type RouteMetrics struct {
	TotalMessages     int64
	SuccessMessages   int64
	FailedMessages    int64
	AverageLatency    time.Duration
	LastLatency       time.Duration
	LastError         error
	LastErrorTime     time.Time
	ConsecutiveErrors int
}

// RouteManager handles SMPP route management
type RouteManager struct {
	mu           sync.RWMutex
	routes       map[string]*Route
	activeRoutes map[string]bool
	routeOrder   []string
}

// NewRouteManager creates a new route manager
func NewRouteManager() *RouteManager {
	return &RouteManager{
		routes:       make(map[string]*Route),
		activeRoutes: make(map[string]bool),
		routeOrder:   make([]string, 0),
	}
}

// AddRoute adds a new route
func (rm *RouteManager) AddRoute(route *Route) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if _, exists := rm.routes[route.ID]; exists {
		return ErrRouteExists
	}

	route.CreatedAt = time.Now()
	route.UpdatedAt = time.Now()
	route.Metrics = &RouteMetrics{}

	rm.routes[route.ID] = route
	rm.updateRouteOrder()
	return nil
}

// UpdateRoute updates an existing route
func (rm *RouteManager) UpdateRoute(route *Route) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	existing, exists := rm.routes[route.ID]
	if !exists {
		return ErrRouteNotFound
	}

	route.CreatedAt = existing.CreatedAt
	route.UpdatedAt = time.Now()
	route.Metrics = existing.Metrics

	rm.routes[route.ID] = route
	rm.updateRouteOrder()
	return nil
}

// GetRoute returns a route by ID
func (rm *RouteManager) GetRoute(routeID string) (*Route, error) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	route, exists := rm.routes[routeID]
	if !exists {
		return nil, ErrRouteNotFound
	}
	return route, nil
}

// GetActiveRoutes returns all active routes in priority order
func (rm *RouteManager) GetActiveRoutes() []*Route {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	var activeRoutes []*Route
	for _, routeID := range rm.routeOrder {
		if route, exists := rm.routes[routeID]; exists && route.Enabled && rm.activeRoutes[routeID] {
			activeRoutes = append(activeRoutes, route)
		}
	}
	return activeRoutes
}

// UpdateRouteStatus updates a route's health status
func (rm *RouteManager) UpdateRouteStatus(routeID string, status string, err error) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	route, exists := rm.routes[routeID]
	if !exists {
		return ErrRouteNotFound
	}

	route.HealthStatus = status
	route.LastHealthCheck = time.Now()

	if err != nil {
		route.Metrics.LastError = err
		route.Metrics.LastErrorTime = time.Now()
		route.Metrics.ConsecutiveErrors++
		rm.activeRoutes[routeID] = false
	} else {
		route.Metrics.ConsecutiveErrors = 0
		rm.activeRoutes[routeID] = true
	}

	return nil
}

// RecordMetrics records performance metrics for a route
func (rm *RouteManager) RecordMetrics(routeID string, success bool, latency time.Duration) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	route, exists := rm.routes[routeID]
	if !exists {
		return ErrRouteNotFound
	}

	route.Metrics.TotalMessages++
	route.Metrics.LastLatency = latency

	// Update average latency
	if route.Metrics.AverageLatency == 0 {
		route.Metrics.AverageLatency = latency
	} else {
		route.Metrics.AverageLatency = (route.Metrics.AverageLatency + latency) / 2
	}

	if success {
		route.Metrics.SuccessMessages++
	} else {
		route.Metrics.FailedMessages++
	}

	return nil
}

// updateRouteOrder updates the ordered list of routes based on priority and weight
func (rm *RouteManager) updateRouteOrder() {
	routes := make([]*Route, 0, len(rm.routes))
	for _, route := range rm.routes {
		routes = append(routes, route)
	}

	// Sort routes by priority (higher first) and weight (higher first)
	sort.Slice(routes, func(i, j int) bool {
		if routes[i].Priority == routes[j].Priority {
			return routes[i].Weight > routes[j].Weight
		}
		return routes[i].Priority > routes[j].Priority
	})

	rm.routeOrder = make([]string, len(routes))
	for i, route := range routes {
		rm.routeOrder[i] = route.ID
	}
}
