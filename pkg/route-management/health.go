package routemanagement

import (
	"context"
	"sync"
	"time"
)

// HealthStatus represents the health status of a route
type HealthStatus string

const (
	// HealthStatusUp indicates the route is healthy and operational
	HealthStatusUp HealthStatus = "UP"

	// HealthStatusDown indicates the route is not operational
	HealthStatusDown HealthStatus = "DOWN"

	// HealthStatusDegraded indicates the route is operational but experiencing issues
	HealthStatusDegraded HealthStatus = "DEGRADED"

	// HealthStatusMaintenance indicates the route is under maintenance
	HealthStatusMaintenance HealthStatus = "MAINTENANCE"
)

// HealthChecker monitors the health of routes
type HealthChecker struct {
	manager       *RouteManager
	checkInterval time.Duration
	timeout       time.Duration
	checkers      map[string]RouteHealthCheck
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
}

// RouteHealthCheck is a function type that performs health check for a route
type RouteHealthCheck func(context.Context, *Route) (HealthStatus, error)

// NewHealthChecker creates a new health checker
func NewHealthChecker(manager *RouteManager, checkInterval, timeout time.Duration) *HealthChecker {
	ctx, cancel := context.WithCancel(context.Background())
	return &HealthChecker{
		manager:       manager,
		checkInterval: checkInterval,
		timeout:       timeout,
		checkers:      make(map[string]RouteHealthCheck),
		ctx:           ctx,
		cancel:        cancel,
	}
}

// RegisterChecker registers a health check function for a route type
func (hc *HealthChecker) RegisterChecker(routeType string, checker RouteHealthCheck) {
	hc.checkers[routeType] = checker
}

// Start begins health checking for all routes
func (hc *HealthChecker) Start() {
	hc.wg.Add(1)
	go hc.runChecks()
}

// Stop stops all health checks
func (hc *HealthChecker) Stop() {
	hc.cancel()
	hc.wg.Wait()
}

// runChecks periodically runs health checks for all routes
func (hc *HealthChecker) runChecks() {
	defer hc.wg.Done()

	ticker := time.NewTicker(hc.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-hc.ctx.Done():
			return
		case <-ticker.C:
			hc.checkAllRoutes()
		}
	}
}

// checkAllRoutes performs health checks on all routes
func (hc *HealthChecker) checkAllRoutes() {
	routes := hc.manager.GetActiveRoutes()

	for _, route := range routes {
		hc.wg.Add(1)
		go func(r *Route) {
			defer hc.wg.Done()
			hc.checkRoute(r)
		}(route)
	}
}

// checkRoute performs a health check on a single route
func (hc *HealthChecker) checkRoute(route *Route) {
	checker, exists := hc.checkers[route.BindType]
	if !exists {
		return
	}

	ctx, cancel := context.WithTimeout(hc.ctx, hc.timeout)
	defer cancel()

	status, err := checker(ctx, route)
	if err != nil {
		hc.manager.UpdateRouteStatus(route.ID, string(HealthStatusDown), err)
		return
	}

	hc.manager.UpdateRouteStatus(route.ID, string(status), nil)
}

// DefaultSMPPHealthCheck provides a default health check for SMPP routes
func DefaultSMPPHealthCheck(ctx context.Context, route *Route) (HealthStatus, error) {
	// Check consecutive errors
	if route.Metrics.ConsecutiveErrors >= 3 {
		return HealthStatusDown, ErrRouteUnhealthy
	}

	// Check success rate
	if route.Metrics.TotalMessages > 100 {
		successRate := float64(route.Metrics.SuccessMessages) / float64(route.Metrics.TotalMessages)
		if successRate < 0.9 {
			return HealthStatusDegraded, nil
		}
	}

	// Check average latency
	if route.Metrics.AverageLatency > time.Second {
		return HealthStatusDegraded, nil
	}

	return HealthStatusUp, nil
}

// DefaultHTTPHealthCheck provides a default health check for HTTP routes
func DefaultHTTPHealthCheck(ctx context.Context, route *Route) (HealthStatus, error) {
	// Similar to SMPP health check but with HTTP-specific metrics
	if route.Metrics.ConsecutiveErrors >= 3 {
		return HealthStatusDown, ErrRouteUnhealthy
	}

	if route.Metrics.TotalMessages > 100 {
		successRate := float64(route.Metrics.SuccessMessages) / float64(route.Metrics.TotalMessages)
		if successRate < 0.95 {
			return HealthStatusDegraded, nil
		}
	}

	if route.Metrics.AverageLatency > 500*time.Millisecond {
		return HealthStatusDegraded, nil
	}

	return HealthStatusUp, nil
}
