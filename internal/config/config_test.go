package config

import (
	"encoding/json"
	"testing"
)

// TestNewAppConfig tests creating new app config
func TestNewAppConfig(t *testing.T) {
	cfg := NewAppConfig()

	if cfg == nil {
		t.Fatal("Expected config to be created")
	}

	if cfg.Accounts == nil {
		t.Error("Expected Accounts to be initialized")
	}

	if cfg.ActivityLog == nil {
		t.Error("Expected ActivityLog to be initialized")
	}

	if cfg.HealthChecks == nil {
		t.Error("Expected HealthChecks to be initialized")
	}
}

// TestDefaultPlatform tests default platform creation
func TestDefaultPlatform(t *testing.T) {
	platform := DefaultPlatform()

	if platform == nil {
		t.Fatal("Expected platform to be created")
	}

	if platform.Type != "github" {
		t.Errorf("Expected default platform type 'github', got '%s'", platform.Type)
	}
}

// TestAccountToJSON tests account serialization
func TestAccountToJSON(t *testing.T) {
	acc := Account{
		Name:        "test-account",
		GitUserName: "Test User",
		GitEmail:    "test@example.com",
		SSH: &SshConfig{
			KeyPath:   "~/.ssh/id_ed25519",
			HostAlias: "github-test",
		},
		Token: &TokenConfig{
			Username: "testuser",
			Token:    "ghp_xxxx",
		},
		Platform: &PlatformConfig{
			Type:   "github",
			Domain: "",
		},
	}

	jsonStr, err := acc.ToJSON()
	if err != nil {
		t.Fatalf("Failed to serialize account: %v", err)
	}

	if jsonStr == "" {
		t.Error("Expected non-empty JSON string")
	}

	// Verify it's valid JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		t.Errorf("Produced invalid JSON: %v", err)
	}
}

// TestAccountFromJSON tests account deserialization
func TestAccountFromJSON(t *testing.T) {
	jsonStr := `{
		"name": "test-account",
		"gitUserName": "Test User",
		"gitEmail": "test@example.com",
		"ssh": {
			"keyPath": "~/.ssh/id_ed25519",
			"hostAlias": "github-test"
		},
		"token": {
			"username": "testuser",
			"token": "ghp_xxxx"
		},
		"platform": {
			"type": "github"
		}
	}`

	acc, err := AccountFromJSON(jsonStr)
	if err != nil {
		t.Fatalf("Failed to deserialize account: %v", err)
	}

	if acc.Name != "test-account" {
		t.Errorf("Expected name 'test-account', got '%s'", acc.Name)
	}

	if acc.SSH == nil {
		t.Error("Expected SSH config to be deserialized")
	}

	if acc.Token == nil {
		t.Error("Expected Token config to be deserialized")
	}

	if acc.Platform == nil {
		t.Error("Expected Platform config to be deserialized")
	}
}

// TestAccountRoundTrip tests serialization round-trip
func TestAccountRoundTrip(t *testing.T) {
	original := Account{
		Name:        "round-trip-test",
		GitUserName: "Round Trip User",
		GitEmail:    "roundtrip@example.com",
		SSH: &SshConfig{
			KeyPath:   "~/.ssh/id_ed25519_rt",
			HostAlias: "github-rt",
		},
		Token: &TokenConfig{
			Username: "rtuser",
			Token:    "ghp_roundtrip",
		},
		Platform: &PlatformConfig{
			Type:   "gitlab",
			Domain: "gitlab.company.com",
			ApiUrl: "https://gitlab.company.com/api/v4",
		},
	}

	// Serialize
	jsonStr, err := original.ToJSON()
	if err != nil {
		t.Fatalf("Failed to serialize: %v", err)
	}

	// Deserialize
	restored, err := AccountFromJSON(jsonStr)
	if err != nil {
		t.Fatalf("Failed to deserialize: %v", err)
	}

	// Compare
	if !original.Equals(restored) {
		t.Error("Round-trip failed: restored account doesn't match original")
	}
}

// TestAccountClone tests account cloning
func TestAccountClone(t *testing.T) {
	original := Account{
		Name:        "clone-test",
		GitUserName: "Clone User",
		GitEmail:    "clone@example.com",
		SSH: &SshConfig{
			KeyPath:   "~/.ssh/id_ed25519_clone",
			HostAlias: "github-clone",
		},
		Token: &TokenConfig{
			Username: "cloneuser",
			Token:    "ghp_clone",
		},
		Platform: &PlatformConfig{
			Type:   "github",
			Domain: "",
		},
	}

	clone := original.Clone()

	// Verify clone equals original
	if !original.Equals(&clone) {
		t.Error("Clone doesn't match original")
	}

	// Verify it's a deep copy (modifying clone doesn't affect original)
	clone.Name = "modified"
	if original.Name == "modified" {
		t.Error("Clone is not a deep copy - modifying clone affected original")
	}

	clone.SSH.KeyPath = "modified-path"
	if original.SSH.KeyPath == "modified-path" {
		t.Error("Clone SSH is not a deep copy")
	}
}

// TestAccountEquals tests account equality
func TestAccountEquals(t *testing.T) {
	acc1 := &Account{
		Name:        "test",
		GitUserName: "User",
		GitEmail:    "user@example.com",
	}

	acc2 := &Account{
		Name:        "test",
		GitUserName: "User",
		GitEmail:    "user@example.com",
	}

	if !acc1.Equals(acc2) {
		t.Error("Expected equal accounts to be equal")
	}

	// Different name
	acc2.Name = "different"
	if acc1.Equals(acc2) {
		t.Error("Expected accounts with different names to not be equal")
	}

	// Nil comparison
	if acc1.Equals(nil) {
		t.Error("Expected non-nil account to not equal nil")
	}

	var nilAcc *Account
	if nilAcc.Equals(acc1) {
		t.Error("Expected nil account to not equal non-nil")
	}

	if !nilAcc.Equals(nil) {
		t.Error("Expected nil to equal nil")
	}
}

// TestAccountEqualsWithOptionalFields tests equality with optional fields
func TestAccountEqualsWithOptionalFields(t *testing.T) {
	// One with SSH, one without
	acc1 := &Account{
		Name: "test",
		SSH:  &SshConfig{KeyPath: "~/.ssh/key"},
	}
	acc2 := &Account{
		Name: "test",
	}

	if acc1.Equals(acc2) {
		t.Error("Expected accounts with different SSH configs to not be equal")
	}

	// Both with SSH but different values
	acc2.SSH = &SshConfig{KeyPath: "~/.ssh/different"}
	if acc1.Equals(acc2) {
		t.Error("Expected accounts with different SSH key paths to not be equal")
	}

	// Same SSH
	acc2.SSH.KeyPath = "~/.ssh/key"
	if !acc1.Equals(acc2) {
		t.Error("Expected accounts with same SSH to be equal")
	}
}

// TestAppConfigToJSON tests app config serialization
func TestAppConfigToJSON(t *testing.T) {
	cfg := NewAppConfig()
	cfg.Accounts = append(cfg.Accounts, Account{Name: "test"})

	jsonStr, err := cfg.ToJSON()
	if err != nil {
		t.Fatalf("Failed to serialize config: %v", err)
	}

	if jsonStr == "" {
		t.Error("Expected non-empty JSON string")
	}
}

// TestAppConfigFromJSON tests app config deserialization
func TestAppConfigFromJSON(t *testing.T) {
	jsonStr := `{
		"accounts": [
			{"name": "account1"},
			{"name": "account2"}
		],
		"activityLog": [],
		"healthChecks": []
	}`

	cfg, err := AppConfigFromJSON(jsonStr)
	if err != nil {
		t.Fatalf("Failed to deserialize config: %v", err)
	}

	if len(cfg.Accounts) != 2 {
		t.Errorf("Expected 2 accounts, got %d", len(cfg.Accounts))
	}
}

// TestAppConfigFromJSONWithMissingFields tests graceful handling of missing fields
func TestAppConfigFromJSONWithMissingFields(t *testing.T) {
	// Minimal JSON with only accounts
	jsonStr := `{"accounts": [{"name": "test"}]}`

	cfg, err := AppConfigFromJSON(jsonStr)
	if err != nil {
		t.Fatalf("Failed to deserialize minimal config: %v", err)
	}

	// Should have initialized slices
	if cfg.ActivityLog == nil {
		t.Error("Expected ActivityLog to be initialized")
	}

	if cfg.HealthChecks == nil {
		t.Error("Expected HealthChecks to be initialized")
	}
}

// TestAccountFromJSONWithMissingOptionalFields tests graceful handling
func TestAccountFromJSONWithMissingOptionalFields(t *testing.T) {
	// Minimal account JSON
	jsonStr := `{"name": "minimal"}`

	acc, err := AccountFromJSON(jsonStr)
	if err != nil {
		t.Fatalf("Failed to deserialize minimal account: %v", err)
	}

	if acc.Name != "minimal" {
		t.Errorf("Expected name 'minimal', got '%s'", acc.Name)
	}

	// Optional fields should be nil/empty
	if acc.SSH != nil {
		t.Error("Expected SSH to be nil for minimal account")
	}

	if acc.Token != nil {
		t.Error("Expected Token to be nil for minimal account")
	}

	if acc.Platform != nil {
		t.Error("Expected Platform to be nil for minimal account")
	}
}
