package config

import (
	"crypto/tls"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the core configuration architecture
type Config struct {
	Web     WebConfig    `yaml:"web"`
	Plugins PluginConfig `yaml:"plugins"`
}

// WebConfig defines the underlying web service routing
type WebConfig struct {
	Network string    `yaml:"network"` // tcp, tcp4, unix
	Addr    string    `yaml:"addr"`    // listen address
	TLS     TLSConfig `yaml:"tls"`
}

// TLSConfig handles certificates and multiplexing protocols
type TLSConfig struct {
	CertFile string   `yaml:"certFile"`
	KeyFile  string   `yaml:"keyFile"`
	CertText string   `yaml:"certText"`
	CertKey  string   `yaml:"certKey"`
	ALPN     []string `yaml:"alpn"`
}

// PluginConfig handles directory bindings for isolation
type PluginConfig struct {
	Install string `yaml:"install"` // read-only plugin code dir
	Data    string `yaml:"data"`    // stateful plugin data dir
}

// LoadConfig parses the YAML configuration file from the given path
func LoadConfig(path string, cfg *Config) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	cfg.Web.Network = "tcp"
	cfg.Web.Addr = ":9026"
	cfg.Plugins.Install = "plugins"
	cfg.Plugins.Data = "plugins-data"

	// Parse YAML payload
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return fmt.Errorf("invalid configuration syntax: %w", err)
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
		// Zero-IO: direct memory initialization to protect disk longevity
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
