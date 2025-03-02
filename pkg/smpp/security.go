package smpp

import (
	"fmt"
	"net"
	"sync"

	"golang.org/x/time/rate"
)

// IPFilter manages IP whitelisting and blacklisting
type IPFilter struct {
	mu           sync.RWMutex
	whitelist    map[string]struct{}
	blacklist    map[string]struct{}
	cidrs        []*net.IPNet
	defaultAllow bool
}

// NewIPFilter creates a new IP filter
func NewIPFilter(defaultAllow bool) *IPFilter {
	return &IPFilter{
		whitelist:    make(map[string]struct{}),
		blacklist:    make(map[string]struct{}),
		cidrs:        make([]*net.IPNet, 0),
		defaultAllow: defaultAllow,
	}
}

// AddToWhitelist adds an IP or CIDR to the whitelist
func (f *IPFilter) AddToWhitelist(ipOrCIDR string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, cidr, err := net.ParseCIDR(ipOrCIDR); err == nil {
		f.cidrs = append(f.cidrs, cidr)
		return nil
	}

	if ip := net.ParseIP(ipOrCIDR); ip != nil {
		f.whitelist[ip.String()] = struct{}{}
		return nil
	}

	return fmt.Errorf("invalid IP or CIDR: %s", ipOrCIDR)
}

// AddToBlacklist adds an IP to the blacklist
func (f *IPFilter) AddToBlacklist(ip string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if parsedIP := net.ParseIP(ip); parsedIP != nil {
		f.blacklist[parsedIP.String()] = struct{}{}
		return nil
	}

	return fmt.Errorf("invalid IP: %s", ip)
}

// IsAllowed checks if an IP is allowed
func (f *IPFilter) IsAllowed(ip string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	// Check blacklist first
	if _, blacklisted := f.blacklist[parsedIP.String()]; blacklisted {
		return false
	}

	// Check whitelist
	if _, whitelisted := f.whitelist[parsedIP.String()]; whitelisted {
		return true
	}

	// Check CIDR ranges
	for _, cidr := range f.cidrs {
		if cidr.Contains(parsedIP) {
			return true
		}
	}

	return f.defaultAllow
}

// RateLimiter manages rate limiting per IP
type RateLimiter struct {
	mu       sync.RWMutex
	limiters map[string]*rate.Limiter
	rate     rate.Limit
	burst    int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(r float64, burst int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Limit(r),
		burst:    burst,
	}
}

// GetLimiter gets or creates a limiter for an IP
func (rl *RateLimiter) GetLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[ip] = limiter
	}

	return limiter
}

// Allow checks if a request from an IP should be allowed
func (rl *RateLimiter) Allow(ip string) bool {
	return rl.GetLimiter(ip).Allow()
}

// SecurityManager manages IP filtering and rate limiting
type SecurityManager struct {
	ipFilter    *IPFilter
	rateLimiter *RateLimiter
}

// NewSecurityManager creates a new security manager
func NewSecurityManager(defaultAllow bool, ratePerSecond float64, burst int) *SecurityManager {
	return &SecurityManager{
		ipFilter:    NewIPFilter(defaultAllow),
		rateLimiter: NewRateLimiter(ratePerSecond, burst),
	}
}

// CheckIP checks if an IP is allowed and not rate limited
func (sm *SecurityManager) CheckIP(ip string) error {
	if !sm.ipFilter.IsAllowed(ip) {
		return fmt.Errorf("IP %s is not allowed", ip)
	}

	if !sm.rateLimiter.Allow(ip) {
		return fmt.Errorf("rate limit exceeded for IP %s", ip)
	}

	return nil
}
