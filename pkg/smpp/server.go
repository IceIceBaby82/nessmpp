package smpp

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/tarik/nessmpp/pkg/pdu"
)

// Server represents an SMPP server
type Server struct {
	addr      string
	listener  net.Listener
	tlsConfig *tls.Config
	sessions  map[string]*Session
	mu        sync.RWMutex
	wg        sync.WaitGroup
	done      chan struct{}
	handlers  map[uint32]PDUHandler
	pool      *ConnectionPool
	security  *SecurityManager
	config    *ServerConfig
}

// ServerConfig holds server configuration parameters
type ServerConfig struct {
	Address        string
	TLSConfig      *TLSConfig
	MaxConnections int
	IdleTimeout    time.Duration
	MaxLifetime    time.Duration
	RateLimit      float64
	RateBurst      int
	DefaultAllow   bool
	SystemID       string
	Password       string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	BindTimeout    time.Duration
}

// SessionState represents the state of an SMPP session
type SessionState int

const (
	// OPEN indicates the session is open but not bound
	OPEN SessionState = iota
	// BOUND_TX indicates the session is bound as transmitter
	BOUND_TX
	// BOUND_RX indicates the session is bound as receiver
	BOUND_RX
	// BOUND_TRX indicates the session is bound as transceiver
	BOUND_TRX
	// CLOSED indicates the session is closed
	CLOSED
)

// Session represents an SMPP connection session
type Session struct {
	mu        sync.RWMutex
	conn      net.Conn
	server    *Server
	state     SessionState
	systemID  string
	createdAt time.Time
	lastUsed  time.Time
	busy      bool
}

// PDUHandler is a function type that handles specific PDU types
type PDUHandler func(*Session, pdu.Header) error

// NewServer creates a new SMPP server
func NewServer(config *ServerConfig) (*Server, error) {
	if config == nil {
		return nil, fmt.Errorf("server config is nil")
	}

	var tlsConfig *tls.Config
	if config.TLSConfig != nil {
		if err := ValidateTLSConfig(config.TLSConfig); err != nil {
			return nil, fmt.Errorf("invalid TLS config: %v", err)
		}
		var err error
		tlsConfig, err = NewTLSConfig(config.TLSConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create TLS config: %v", err)
		}
	}

	server := &Server{
		addr:      config.Address,
		tlsConfig: tlsConfig,
		sessions:  make(map[string]*Session),
		done:      make(chan struct{}),
		handlers:  make(map[uint32]PDUHandler),
		pool:      NewConnectionPool(config.MaxConnections, config.IdleTimeout, config.MaxLifetime),
		security:  NewSecurityManager(config.DefaultAllow, config.RateLimit, config.RateBurst),
		config:    config,
	}

	// Register default handlers
	server.registerDefaultHandlers()

	return server, nil
}

// Start starts the SMPP server
func (s *Server) Start() error {
	var err error
	if s.tlsConfig != nil {
		s.listener, err = tls.Listen("tcp", s.addr, s.tlsConfig)
	} else {
		s.listener, err = net.Listen("tcp", s.addr)
	}

	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}

	go s.acceptConnections()
	return nil
}

// Stop stops the SMPP server
func (s *Server) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

// acceptConnections accepts incoming connections
func (s *Server) acceptConnections() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			// Check if server is shutting down
			select {
			case <-s.done:
				return
			default:
				log.Printf("Failed to accept connection: %v", err)
				continue
			}
		}

		// Get client IP
		ip, _, err := net.SplitHostPort(conn.RemoteAddr().String())
		if err != nil {
			log.Printf("Failed to get client IP: %v", err)
			conn.Close()
			continue
		}

		// Check security
		if err := s.security.CheckIP(ip); err != nil {
			log.Printf("Connection rejected: %v", err)
			conn.Close()
			continue
		}

		s.wg.Add(1)
		go func(c net.Conn) {
			defer s.wg.Done()
			s.handleConnection(c)
		}(conn)
	}
}

func (s *Server) registerDefaultHandlers() {
	// Bind operations
	s.handlers[pdu.BIND_TRANSMITTER] = s.handleBindTransmitter
	s.handlers[pdu.BIND_RECEIVER] = s.handleBindReceiver
	s.handlers[pdu.BIND_TRANSCEIVER] = s.handleBindTransceiver

	// Messaging operations
	s.handlers[pdu.SUBMIT_SM] = handleSubmitSM
	s.handlers[pdu.DELIVER_SM] = handleDeliverSM
	s.handlers[pdu.DATA_SM] = handleDataSM

	// Query operations
	s.handlers[pdu.QUERY_SM] = handleQuerySM
	s.handlers[pdu.QUERY_BROADCAST_SM] = handleQueryBroadcastSM

	// Cancel operations
	s.handlers[pdu.CANCEL_SM] = handleCancelSM
	s.handlers[pdu.CANCEL_BROADCAST_SM] = handleCancelBroadcastSM

	// Replace operation
	s.handlers[pdu.REPLACE_SM] = handleReplaceSM

	// Broadcast operations
	s.handlers[pdu.BROADCAST_SM] = handleBroadcastSM

	// Connection management
	s.handlers[pdu.UNBIND] = handleUnbind
	s.handlers[pdu.ENQUIRE_LINK] = handleEnquireLink
	s.handlers[pdu.GENERIC_NACK] = handleGenericNack
}

// handleConnection handles a new connection
func (s *Server) handleConnection(conn net.Conn) {
	session := NewSession(conn, s)

	// Add session to pool
	if err := s.pool.Add(session); err != nil {
		log.Printf("Failed to add session to pool: %v", err)
		conn.Close()
		return
	}

	defer func() {
		s.pool.Remove(session)
		conn.Close()
	}()

	// Start reading PDUs
	headerBuf := make([]byte, 16)
	for {
		select {
		case <-s.done:
			return
		default:
		}

		// Set read deadline
		if err := conn.SetReadDeadline(time.Now().Add(time.Second * 30)); err != nil {
			log.Printf("Failed to set read deadline: %v", err)
			return
		}

		// Read PDU header
		if _, err := io.ReadFull(conn, headerBuf); err != nil {
			if err != io.EOF {
				log.Printf("Failed to read PDU header: %v", err)
			}
			return
		}

		header := &pdu.Header{}
		if err := header.Unmarshal(headerBuf); err != nil {
			log.Printf("Failed to unmarshal PDU header: %v", err)
			if err := session.sendGenericNack(0, pdu.ESME_RINVMSGLEN); err != nil {
				log.Printf("Failed to send generic_nack: %v", err)
			}
			continue
		}

		// Get handler for command
		handler, exists := s.handlers[header.CommandID]
		if !exists {
			log.Printf("No handler for command: %d", header.CommandID)
			if err := session.sendGenericNack(header.SequenceNumber, pdu.ESME_RINVCMDID); err != nil {
				log.Printf("Failed to send generic_nack: %v", err)
			}
			continue
		}

		// Handle PDU
		if err := handler(session, *header); err != nil {
			log.Printf("Failed to handle PDU: %v", err)
			if err := session.sendGenericNack(header.SequenceNumber, pdu.ESME_RSYSERR); err != nil {
				log.Printf("Failed to send generic_nack: %v", err)
			}
			continue
		}

		// Update last used time
		session.lastUsed = time.Now()
	}
}

// Handler implementations

// handleBindTransmitter handles a bind_transmitter PDU
func (s *Server) handleBindTransmitter(sess *Session, header pdu.Header) error {
	// Read body data
	bodyLen := int(header.CommandLength) - 16
	bodyData := make([]byte, bodyLen)
	if _, err := io.ReadFull(sess.conn, bodyData); err != nil {
		resp := &pdu.Header{
			CommandLength:  16,
			CommandID:      pdu.BIND_TRANSMITTER_RESP,
			CommandStatus:  pdu.ESME_RSYSERR,
			SequenceNumber: header.SequenceNumber,
		}
		return sess.sendPDU(resp, nil)
	}

	// Unmarshal bind request
	bindReq := &pdu.BindTransmitter{}
	if err := bindReq.Unmarshal(bodyData); err != nil {
		resp := &pdu.Header{
			CommandLength:  16,
			CommandID:      pdu.BIND_TRANSMITTER_RESP,
			CommandStatus:  pdu.ESME_RSYSERR,
			SequenceNumber: header.SequenceNumber,
		}
		return sess.sendPDU(resp, nil)
	}

	// Validate credentials
	if !s.validateCredentials(bindReq.SystemID, bindReq.Password) {
		resp := &pdu.Header{
			CommandLength:  16,
			CommandID:      pdu.BIND_TRANSMITTER_RESP,
			CommandStatus:  pdu.ESME_RINVPASWD,
			SequenceNumber: header.SequenceNumber,
		}
		return sess.sendPDU(resp, nil)
	}

	// Update session state
	sess.mu.Lock()
	sess.systemID = bindReq.SystemID
	sess.state = BOUND_TX
	sess.mu.Unlock()

	// Send success response
	resp := &pdu.Header{
		CommandLength:  16,
		CommandID:      pdu.BIND_TRANSMITTER_RESP,
		CommandStatus:  pdu.ESME_ROK,
		SequenceNumber: header.SequenceNumber,
	}

	return sess.sendPDU(resp, nil)
}

// handleBindReceiver handles bind_receiver PDU
func (s *Server) handleBindReceiver(sess *Session, header pdu.Header) error {
	// Read body data
	bodyLen := int(header.CommandLength) - 16
	bodyData := make([]byte, bodyLen)
	if _, err := io.ReadFull(sess.conn, bodyData); err != nil {
		resp := &pdu.Header{
			CommandLength:  16,
			CommandID:      pdu.BIND_RECEIVER_RESP,
			CommandStatus:  pdu.ESME_RSYSERR,
			SequenceNumber: header.SequenceNumber,
		}
		return sess.sendPDU(resp, nil)
	}

	// Unmarshal bind request
	bindReq := &pdu.BindReceiver{}
	if err := bindReq.Unmarshal(bodyData); err != nil {
		resp := &pdu.Header{
			CommandLength:  16,
			CommandID:      pdu.BIND_RECEIVER_RESP,
			CommandStatus:  pdu.ESME_RSYSERR,
			SequenceNumber: header.SequenceNumber,
		}
		return sess.sendPDU(resp, nil)
	}

	// Validate credentials
	if !s.validateCredentials(bindReq.SystemID, bindReq.Password) {
		resp := &pdu.Header{
			CommandLength:  16,
			CommandID:      pdu.BIND_RECEIVER_RESP,
			CommandStatus:  pdu.ESME_RINVPASWD,
			SequenceNumber: header.SequenceNumber,
		}
		return sess.sendPDU(resp, nil)
	}

	// Update session state
	sess.mu.Lock()
	sess.systemID = bindReq.SystemID
	sess.state = BOUND_RX
	sess.mu.Unlock()

	// Send success response
	resp := &pdu.Header{
		CommandLength:  16,
		CommandID:      pdu.BIND_RECEIVER_RESP,
		CommandStatus:  pdu.ESME_ROK,
		SequenceNumber: header.SequenceNumber,
	}

	return sess.sendPDU(resp, nil)
}

// handleBindTransceiver handles bind_transceiver PDU
func (s *Server) handleBindTransceiver(sess *Session, header pdu.Header) error {
	// Read body data
	bodyLen := int(header.CommandLength) - 16
	bodyData := make([]byte, bodyLen)
	if _, err := io.ReadFull(sess.conn, bodyData); err != nil {
		resp := &pdu.Header{
			CommandLength:  16,
			CommandID:      pdu.BIND_TRANSCEIVER_RESP,
			CommandStatus:  pdu.ESME_RSYSERR,
			SequenceNumber: header.SequenceNumber,
		}
		return sess.sendPDU(resp, nil)
	}

	// Unmarshal bind request
	bindReq := &pdu.BindTransceiver{}
	if err := bindReq.Unmarshal(bodyData); err != nil {
		resp := &pdu.Header{
			CommandLength:  16,
			CommandID:      pdu.BIND_TRANSCEIVER_RESP,
			CommandStatus:  pdu.ESME_RSYSERR,
			SequenceNumber: header.SequenceNumber,
		}
		return sess.sendPDU(resp, nil)
	}

	// Validate credentials
	if !s.validateCredentials(bindReq.SystemID, bindReq.Password) {
		resp := &pdu.Header{
			CommandLength:  16,
			CommandID:      pdu.BIND_TRANSCEIVER_RESP,
			CommandStatus:  pdu.ESME_RINVPASWD,
			SequenceNumber: header.SequenceNumber,
		}
		return sess.sendPDU(resp, nil)
	}

	// Update session state
	sess.mu.Lock()
	sess.systemID = bindReq.SystemID
	sess.state = BOUND_TRX
	sess.mu.Unlock()

	// Send success response
	resp := &pdu.Header{
		CommandLength:  16,
		CommandID:      pdu.BIND_TRANSCEIVER_RESP,
		CommandStatus:  pdu.ESME_ROK,
		SequenceNumber: header.SequenceNumber,
	}

	return sess.sendPDU(resp, nil)
}

// Helper function to send bind response
func sendBindResp(sess *Session, cmdID uint32, status uint32, systemID string, tlvs map[uint16]*pdu.TLVParam) error {
	resp := &pdu.Header{
		CommandLength:  16, // Base header size
		CommandID:      cmdID,
		CommandStatus:  status,
		SequenceNumber: sess.nextSequenceNumber(),
	}

	if len(systemID) > 0 {
		resp.CommandLength += uint32(len(systemID) + 1) // +1 for null terminator
	}

	if tlvs != nil {
		for _, tlv := range tlvs {
			resp.CommandLength += 4 + uint32(len(tlv.Value)) // Tag(2) + Length(2) + Value(n)
		}
	}

	return sess.sendPDU(resp, nil)
}

// validateCredentials validates the system_id and password
func validateCredentials(systemID, password string) error {
	// TODO: Implement actual credential validation
	// This should check against configured valid credentials
	return nil
}

func handleSubmitSM(sess *Session, header pdu.Header) error {
	// TODO: Implement submit_sm handling
	return nil
}

func handleDeliverSM(sess *Session, header pdu.Header) error {
	// TODO: Implement deliver_sm handling
	return nil
}

func handleDataSM(sess *Session, header pdu.Header) error {
	// TODO: Implement data_sm handling
	return nil
}

func handleQuerySM(sess *Session, header pdu.Header) error {
	// TODO: Implement query_sm handling
	return nil
}

func handleQueryBroadcastSM(sess *Session, header pdu.Header) error {
	// TODO: Implement query_broadcast_sm handling
	return nil
}

func handleCancelSM(sess *Session, header pdu.Header) error {
	// TODO: Implement cancel_sm handling
	return nil
}

func handleCancelBroadcastSM(sess *Session, header pdu.Header) error {
	// TODO: Implement cancel_broadcast_sm handling
	return nil
}

func handleReplaceSM(sess *Session, header pdu.Header) error {
	// TODO: Implement replace_sm handling
	return nil
}

func handleBroadcastSM(sess *Session, header pdu.Header) error {
	// TODO: Implement broadcast_sm handling
	return nil
}

func handleUnbind(sess *Session, header pdu.Header) error {
	// TODO: Implement unbind handling
	return nil
}

func handleEnquireLink(sess *Session, header pdu.Header) error {
	// TODO: Implement enquire_link handling
	return nil
}

func handleGenericNack(sess *Session, header pdu.Header) error {
	// TODO: Implement generic_nack handling
	return nil
}

// Helper methods for Session
func (sess *Session) sendPDU(header *pdu.Header, body []byte) error {
	sess.mu.Lock()
	defer sess.mu.Unlock()

	// Write header
	headerBytes, err := header.Marshal()
	if err != nil {
		return err
	}
	if _, err := sess.conn.Write(headerBytes); err != nil {
		return err
	}

	// Write body if present
	if len(body) > 0 {
		if _, err := sess.conn.Write(body); err != nil {
			return err
		}
	}

	return nil
}

func (sess *Session) nextSequenceNumber() uint32 {
	sess.mu.Lock()
	defer sess.mu.Unlock()
	sess.seqNumber++
	return sess.seqNumber
}

// Helper methods for Server
func (s *Server) addSession(systemID string, sess *Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[systemID] = sess
}

func (s *Server) removeSession(systemID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, systemID)
}

// GetState returns the current session state
func (s *Session) GetState() SessionState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state
}

// SetState sets the session state
func (s *Session) SetState(state SessionState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state = state
}

// GetSystemID returns the system ID of the bound session
func (s *Session) GetSystemID() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.systemID
}

// SetSystemID sets the system ID for the session
func (s *Session) SetSystemID(systemID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.systemID = systemID
}

// Close closes the session and its connection
func (s *Session) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.state = CLOSED
	return s.conn.Close()
}

// sendGenericNack sends a generic NACK PDU
func (s *Session) sendGenericNack(sequence uint32, status uint32) error {
	header := &pdu.Header{
		CommandLength:  16, // Only header, no body
		CommandID:      pdu.GENERIC_NACK,
		CommandStatus:  status,
		SequenceNumber: sequence,
	}
	return s.sendPDU(header, nil)
}

// validateCredentials validates the system ID and password
func (s *Server) validateCredentials(systemID, password string) bool {
	return systemID == s.config.SystemID && password == s.config.Password
}

// NewSession creates a new Session instance
func NewSession(conn net.Conn, server *Server) *Session {
	return &Session{
		conn:      conn,
		server:    server,
		state:     OPEN,
		createdAt: time.Now(),
		lastUsed:  time.Now(),
		busy:      false,
	}
}

// IsBusy returns whether the session is currently in use
func (s *Session) IsBusy() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.busy
}

// SetBusy sets the busy status of the session
func (s *Session) SetBusy(busy bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.busy = busy
}
