package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/dwirx/ghex/internal/platform"
)

// Manager handles configuration loading and saving
type Manager struct {
	primaryPath string
	legacyPath  string
}

// NewManager creates a new configuration manager
func NewManager() *Manager {
	configDir := platform.GetConfigDir("ghe")
	legacyDir := platform.GetConfigDir("github-switch")

	return &Manager{
		primaryPath: filepath.Join(configDir, "config.json"),
		legacyPath:  filepath.Join(legacyDir, "config.json"),
	}
}

// GetConfigPath returns the primary configuration file path
func (m *Manager) GetConfigPath() string {
	return m.primaryPath
}

// Load reads the configuration from disk
// It tries the primary path first, then falls back to legacy path
func (m *Manager) Load() (*AppConfig, error) {
	paths := []string{m.primaryPath, m.legacyPath}

	for _, path := range paths {
		cfg, err := m.loadFromPath(path)
		if err == nil {
			// If loaded from legacy path, migrate to new location
			if path == m.legacyPath {
				_ = m.Save(cfg) // Ignore migration errors
			}
			return cfg, nil
		}

		// If file doesn't exist, try next path
		if os.IsNotExist(err) {
			continue
		}

		// For other errors, return them
		return nil, err
	}

	// No config file found, return empty config
	return NewAppConfig(), nil
}

// loadFromPath loads configuration from a specific path
func (m *Manager) loadFromPath(path string) (*AppConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Ensure accounts is not nil
	if cfg.Accounts == nil {
		cfg.Accounts = []Account{}
	}

	return &cfg, nil
}

// Save writes the configuration to disk
func (m *Manager) Save(cfg *AppConfig) error {
	// Ensure directory exists
	dir := filepath.Dir(m.primaryPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Marshal with indentation for readability
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	// Add trailing newline
	data = append(data, '\n')

	return os.WriteFile(m.primaryPath, data, 0644)
}

// Global manager instance
var defaultManager *Manager

// GetManager returns the default configuration manager
func GetManager() *Manager {
	if defaultManager == nil {
		defaultManager = NewManager()
	}
	return defaultManager
}

// Load is a convenience function to load configuration
func Load() (*AppConfig, error) {
	return GetManager().Load()
}

// Save is a convenience function to save configuration
func Save(cfg *AppConfig) error {
	return GetManager().Save(cfg)
}

// ToJSON serializes an Account to JSON string for debugging
func (a *Account) ToJSON() (string, error) {
	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FromJSON deserializes an Account from JSON string
func AccountFromJSON(jsonStr string) (*Account, error) {
	var acc Account
	if err := json.Unmarshal([]byte(jsonStr), &acc); err != nil {
		return nil, err
	}
	return &acc, nil
}

// ToJSON serializes AppConfig to JSON string for debugging
func (c *AppConfig) ToJSON() (string, error) {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// AppConfigFromJSON deserializes AppConfig from JSON string
func AppConfigFromJSON(jsonStr string) (*AppConfig, error) {
	var cfg AppConfig
	if err := json.Unmarshal([]byte(jsonStr), &cfg); err != nil {
		return nil, err
	}
	
	// Ensure slices are not nil
	if cfg.Accounts == nil {
		cfg.Accounts = []Account{}
	}
	if cfg.ActivityLog == nil {
		cfg.ActivityLog = []ActivityLogEntry{}
	}
	if cfg.HealthChecks == nil {
		cfg.HealthChecks = []HealthStatus{}
	}
	
	return &cfg, nil
}

// Clone creates a deep copy of an Account
func (a *Account) Clone() Account {
	clone := Account{
		Name:        a.Name,
		GitUserName: a.GitUserName,
		GitEmail:    a.GitEmail,
	}
	
	if a.SSH != nil {
		clone.SSH = &SshConfig{
			KeyPath:   a.SSH.KeyPath,
			HostAlias: a.SSH.HostAlias,
		}
	}
	
	if a.Token != nil {
		clone.Token = &TokenConfig{
			Username: a.Token.Username,
			Token:    a.Token.Token,
		}
	}
	
	if a.Platform != nil {
		clone.Platform = &PlatformConfig{
			Type:   a.Platform.Type,
			Domain: a.Platform.Domain,
			ApiUrl: a.Platform.ApiUrl,
		}
	}
	
	return clone
}

// Equals checks if two accounts are equal
func (a *Account) Equals(other *Account) bool {
	if a == nil || other == nil {
		return a == other
	}
	
	if a.Name != other.Name || a.GitUserName != other.GitUserName || a.GitEmail != other.GitEmail {
		return false
	}
	
	// Compare SSH
	if (a.SSH == nil) != (other.SSH == nil) {
		return false
	}
	if a.SSH != nil {
		if a.SSH.KeyPath != other.SSH.KeyPath || a.SSH.HostAlias != other.SSH.HostAlias {
			return false
		}
	}
	
	// Compare Token
	if (a.Token == nil) != (other.Token == nil) {
		return false
	}
	if a.Token != nil {
		if a.Token.Username != other.Token.Username || a.Token.Token != other.Token.Token {
			return false
		}
	}
	
	// Compare Platform
	if (a.Platform == nil) != (other.Platform == nil) {
		return false
	}
	if a.Platform != nil {
		if a.Platform.Type != other.Platform.Type || a.Platform.Domain != other.Platform.Domain || a.Platform.ApiUrl != other.Platform.ApiUrl {
			return false
		}
	}
	
	return true
}
