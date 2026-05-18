package config

import (
	"crypto/tls"
	"encoding/json"
	"fmt"

	"github.com/google/go-jsonnet"
	"github.com/os-webui/os-webui/internal/utils"
)

// Config represents the core configuration architecture
type Config struct {
	Dev     bool         `json:"dev"`
	Web     WebConfig    `json:"web"`
	Plugins PluginConfig `json:"plugins"`
}

// WebConfig defines the underlying web service routing
type WebConfig struct {
	Network string    `json:"network"` // tcp, tcp4, unix
	Addr    string    `json:"addr"`    // listen address
	TLS     TLSConfig `json:"tls"`
}

// TLSConfig handles certificates and multiplexing protocols
type TLSConfig struct {
	CertFile string   `json:"certFile"`
	KeyFile  string   `json:"keyFile"`
	CertText string   `json:"certText"`
	CertKey  string   `json:"certKey"`
	ALPN     []string `json:"alpn"`
}

// PluginConfig handles directory bindings for isolation
type PluginConfig struct {
	Install string `json:"install"` // read-only plugin code dir
	Data    string `json:"data"`    // stateful plugin data dir
	Config  string `json:"config"`  // stateful plugin config dir
}

// LoadConfig parses the JavaScript configuration file from the given path into the provided cfg pointer
func LoadConfig(path string, cfg *Config) error {
	// Guard against nil pointer assignment traps
	if cfg == nil {
		return fmt.Errorf("configuration destination pointer cannot be nil")
	}
	vm := jsonnet.MakeVM()
	jsonStr, err := vm.EvaluateFile(path)
	if err != nil {
		return fmt.Errorf("jsonnet evaluation failed: %w", err)
	}
	if err = json.Unmarshal(utils.StringToBytes(jsonStr), cfg); err != nil {
		return fmt.Errorf("failed to unmarshal JSON payload into config struct: %w", err)
	}
	return nil
}

// MakeTLSConfig evaluates dual-track certificates and compiles native *tls.Config
func (t *TLSConfig) MakeTLSConfig() (*tls.Config, error) {
	// Return nil if no TLS properties are specified (fallback to insecure HTTP)
	if t.CertFile == "" && t.KeyFile == "" && t.CertText == "" && t.CertKey == "" {
		return nil, nil
	}

	var cert tls.Certificate
	var err error

	// Dual-track evaluation: prioritize in-memory strings over local file paths
	if t.CertText != "" && t.CertKey != "" {
		// Defensive Design: Standard Go slice memory allocation []byte(s) is mandated here.
		// tls.X509KeyPair natively mutates the underlying byte array during PEM block parsing.
		// Bypassing StringToBytes here prevents kernel-level read-only memory SIGSEGV segmentation faults.
		cert, err = tls.X509KeyPair([]byte(t.CertText), []byte(t.CertKey))
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS certificate from memory (certText): %w", err)
		}
	} else if t.CertFile != "" && t.KeyFile != "" {
		// Standard localized file parsing
		cert, err = tls.LoadX509KeyPair(t.CertFile, t.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS certificate from file (certFile): %w", err)
		}
	} else {
		return nil, fmt.Errorf("incomplete TLS payload: missing corresponding file/text pairs")
	}

	// Compile structural native tls.Config
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   t.ALPN,           // Bind ALPN (e.g., ['h2', 'http/1.1']) for multiplexing streams
		MinVersion:   tls.VersionTLS12, // Strict baseline: reject insecure legacy connections
	}

	return tlsConfig, nil
}
