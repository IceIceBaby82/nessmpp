package smpp

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
)

// TLSConfig holds TLS configuration parameters
type TLSConfig struct {
	CertFile   string
	KeyFile    string
	CAFile     string
	ClientAuth tls.ClientAuthType
	MinVersion uint16
}

// NewTLSConfig creates a new TLS configuration
func NewTLSConfig(config *TLSConfig) (*tls.Config, error) {
	if config == nil {
		return nil, fmt.Errorf("TLS config is nil")
	}

	cert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load certificate: %v", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   config.ClientAuth,
		MinVersion:   config.MinVersion,
	}

	// Load CA if specified
	if config.CAFile != "" {
		caCert, err := ioutil.ReadFile(config.CAFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA file: %v", err)
		}

		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to parse CA certificate")
		}

		tlsConfig.ClientCAs = caCertPool
	}

	return tlsConfig, nil
}

// DefaultTLSConfig returns a default TLS configuration
func DefaultTLSConfig() *TLSConfig {
	return &TLSConfig{
		ClientAuth: tls.RequireAndVerifyClientCert,
		MinVersion: tls.VersionTLS12,
	}
}

// ValidateTLSConfig validates TLS configuration parameters
func ValidateTLSConfig(config *TLSConfig) error {
	if config == nil {
		return fmt.Errorf("TLS config is nil")
	}

	if config.CertFile == "" {
		return fmt.Errorf("certificate file is required")
	}

	if config.KeyFile == "" {
		return fmt.Errorf("private key file is required")
	}

	if config.ClientAuth == tls.RequireAndVerifyClientCert && config.CAFile == "" {
		return fmt.Errorf("CA file is required when client authentication is enabled")
	}

	return nil
}
