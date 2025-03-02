package clientgroups

import (
	"net"
	"time"
)

// ClientGroup represents a group of SMPP clients with shared configuration
type ClientGroup struct {
	ID               string
	Name             string
	SMPPVersion      string
	MaxBinds         int
	TLSEnabled       bool
	AllowedIPs       []net.IPNet
	GeoRestrictions  []string
	CharacterSets    []string
	CustomErrorCodes map[uint32]string
	RateLimit        int
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// Client represents an SMPP client within a group
type Client struct {
	ID           string
	GroupID      string
	SystemID     string
	Password     string
	IPAddress    net.IP
	BindCount    int
	LastBindTime time.Time
	LastPingTime time.Time
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// ClientManager handles client group operations
type ClientManager struct {
	groups  map[string]*ClientGroup
	clients map[string]*Client
}

// NewClientManager creates a new client manager
func NewClientManager() *ClientManager {
	return &ClientManager{
		groups:  make(map[string]*ClientGroup),
		clients: make(map[string]*Client),
	}
}

// AddGroup adds a new client group
func (cm *ClientManager) AddGroup(group *ClientGroup) error {
	if _, exists := cm.groups[group.ID]; exists {
		return ErrGroupExists
	}
	cm.groups[group.ID] = group
	return nil
}

// AddClient adds a new client to a group
func (cm *ClientManager) AddClient(client *Client) error {
	if _, exists := cm.groups[client.GroupID]; !exists {
		return ErrGroupNotFound
	}
	if _, exists := cm.clients[client.ID]; exists {
		return ErrClientExists
	}
	cm.clients[client.ID] = client
	return nil
}

// ValidateClient validates client credentials and restrictions
func (cm *ClientManager) ValidateClient(systemID, password string, ip net.IP) error {
	client, exists := cm.clients[systemID]
	if !exists {
		return ErrClientNotFound
	}

	group, exists := cm.groups[client.GroupID]
	if !exists {
		return ErrGroupNotFound
	}

	// Validate password
	if client.Password != password {
		return ErrInvalidCredentials
	}

	// Validate IP restrictions
	if !cm.isIPAllowed(ip, group.AllowedIPs) {
		return ErrIPNotAllowed
	}

	// Validate bind count
	if client.BindCount >= group.MaxBinds {
		return ErrMaxBindsExceeded
	}

	return nil
}

// isIPAllowed checks if the IP is in the allowed IP ranges
func (cm *ClientManager) isIPAllowed(ip net.IP, allowedNets []net.IPNet) bool {
	if len(allowedNets) == 0 {
		return true
	}
	for _, net := range allowedNets {
		if net.Contains(ip) {
			return true
		}
	}
	return false
}

// UpdateClientStatus updates a client's status
func (cm *ClientManager) UpdateClientStatus(clientID string, status string) error {
	client, exists := cm.clients[clientID]
	if !exists {
		return ErrClientNotFound
	}
	client.Status = status
	client.UpdatedAt = time.Now()
	return nil
}

// GetClientGroup returns the group for a given client
func (cm *ClientManager) GetClientGroup(clientID string) (*ClientGroup, error) {
	client, exists := cm.clients[clientID]
	if !exists {
		return nil, ErrClientNotFound
	}
	group, exists := cm.groups[client.GroupID]
	if !exists {
		return nil, ErrGroupNotFound
	}
	return group, nil
}
