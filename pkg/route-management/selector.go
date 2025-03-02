package routemanagement

import (
	"math/rand"
	"time"
)

// RouteSelector handles the logic for selecting the best route for message delivery
type RouteSelector struct {
	manager *RouteManager
}

// NewRouteSelector creates a new route selector
func NewRouteSelector(manager *RouteManager) *RouteSelector {
	return &RouteSelector{
		manager: manager,
	}
}

// SelectRoute chooses the best route based on various criteria
func (rs *RouteSelector) SelectRoute(criteria *RouteCriteria) (*Route, error) {
	activeRoutes := rs.manager.GetActiveRoutes()
	if len(activeRoutes) == 0 {
		return nil, ErrNoActiveRoutes
	}

	// Filter routes based on criteria
	eligibleRoutes := rs.filterRoutes(activeRoutes, criteria)
	if len(eligibleRoutes) == 0 {
		return nil, ErrNoEligibleRoutes
	}

	// Group routes by priority
	priorityGroups := make(map[int][]*Route)
	highestPriority := -1
	for _, route := range eligibleRoutes {
		if route.Priority > highestPriority {
			highestPriority = route.Priority
		}
		priorityGroups[route.Priority] = append(priorityGroups[route.Priority], route)
	}

	// Get routes with highest priority
	highPriorityRoutes := priorityGroups[highestPriority]

	// Calculate total weight of high priority routes
	totalWeight := 0
	for _, route := range highPriorityRoutes {
		totalWeight += route.Weight
	}

	// Select route based on weighted random selection
	if totalWeight > 0 {
		randomWeight := rand.Intn(totalWeight)
		currentWeight := 0
		for _, route := range highPriorityRoutes {
			currentWeight += route.Weight
			if randomWeight < currentWeight {
				return route, nil
			}
		}
	}

	// If all weights are 0, select randomly
	return highPriorityRoutes[rand.Intn(len(highPriorityRoutes))], nil
}

// RouteCriteria defines the criteria for route selection
type RouteCriteria struct {
	MaxCost           float64
	RequiredBindType  string
	MinSuccessRate    float64
	MaxLatency        time.Duration
	DestinationRegion string
}

// filterRoutes filters routes based on the given criteria
func (rs *RouteSelector) filterRoutes(routes []*Route, criteria *RouteCriteria) []*Route {
	if criteria == nil {
		return routes
	}

	var filtered []*Route
	for _, route := range routes {
		if !rs.meetsCriteria(route, criteria) {
			continue
		}
		filtered = append(filtered, route)
	}
	return filtered
}

// meetsCriteria checks if a route meets the specified criteria
func (rs *RouteSelector) meetsCriteria(route *Route, criteria *RouteCriteria) bool {
	// Check cost constraint
	if criteria.MaxCost > 0 && route.CostPerMessage > criteria.MaxCost {
		return false
	}

	// Check bind type
	if criteria.RequiredBindType != "" && route.BindType != criteria.RequiredBindType {
		return false
	}

	// Check success rate
	if criteria.MinSuccessRate > 0 {
		totalMessages := float64(route.Metrics.TotalMessages)
		if totalMessages > 0 {
			successRate := float64(route.Metrics.SuccessMessages) / totalMessages
			if successRate < criteria.MinSuccessRate {
				return false
			}
		}
	}

	// Check latency
	if criteria.MaxLatency > 0 && route.Metrics.AverageLatency > criteria.MaxLatency {
		return false
	}

	return true
}

// init initializes the random number generator
func init() {
	rand.Seed(time.Now().UnixNano())
}
