package routemanagement

import "errors"

var (
	// ErrRouteNotFound indicates that the requested route was not found
	ErrRouteNotFound = errors.New("route not found")

	// ErrRouteExists indicates that a route with the same ID already exists
	ErrRouteExists = errors.New("route already exists")

	// ErrInvalidRoute indicates that the route configuration is invalid
	ErrInvalidRoute = errors.New("invalid route configuration")

	// ErrRouteDisabled indicates that the route is currently disabled
	ErrRouteDisabled = errors.New("route is disabled")

	// ErrRouteUnhealthy indicates that the route is currently unhealthy
	ErrRouteUnhealthy = errors.New("route is unhealthy")

	// ErrNoActiveRoutes indicates that there are no active routes available
	ErrNoActiveRoutes = errors.New("no active routes available")

	// ErrNoEligibleRoutes indicates that no routes meet the specified criteria
	ErrNoEligibleRoutes = errors.New("no routes meet the specified criteria")
)
