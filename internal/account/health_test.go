package account

import (
	"testing"
	"time"

	"github.com/dwirx/ghex/internal/config"
)

// TestHealthIndicatorConstants tests health indicator constants
func TestHealthIndicatorConstants(t *testing.T) {
	if HealthValid != "✓" {
		t.Errorf("HealthValid should be '✓', got '%s'", HealthValid)
	}
	if HealthInvalid != "✗" {
		t.Errorf("HealthInvalid should be '✗', got '%s'", HealthInvalid)
	}
	if HealthUnknown != "?" {
		t.Errorf("HealthUnknown should be '?', got '%s'", HealthUnknown)
	}
}

// TestGetHealthIndicator tests health indicator symbol retrieval
func TestGetHealthIndicator(t *testing.T) {
	tests := []struct {
		state    HealthState
		expected string
	}{
		{HealthStateValid, HealthValid},
		{HealthStateInvalid, HealthInvalid},
		{HealthStateUnknown, HealthUnknown},
	}

	for _, tt := range tests {
		result := GetHealthIndicator(tt.state)
		if result != tt.expected {
			t.Errorf("GetHealthIndicator(%v) = %s, expected %s", tt.state, result, tt.expected)
		}
	}
}

// TestGetHealthIndicatorFromBool tests health indicator from bool pointer
func TestGetHealthIndicatorFromBool(t *testing.T) {
	trueVal := true
	falseVal := false

	// Valid
	result := GetHealthIndicatorFromBool(&trueVal)
	if result != HealthValid {
		t.Errorf("Expected '%s' for true, got '%s'", HealthValid, result)
	}

	// Invalid
	result = GetHealthIndicatorFromBool(&falseVal)
	if result != HealthInvalid {
		t.Errorf("Expected '%s' for false, got '%s'", HealthInvalid, result)
	}

	// Unknown (nil)
	result = GetHealthIndicatorFromBool(nil)
	if result != HealthUnknown {
		t.Errorf("Expected '%s' for nil, got '%s'", HealthUnknown, result)
	}
}

// TestIsStaleCheck tests stale check detection
func TestIsStaleCheck(t *testing.T) {
	// Zero time is stale
	if !IsStaleCheck(time.Time{}) {
		t.Error("Zero time should be stale")
	}

	// Recent time is not stale
	recent := time.Now().Add(-1 * time.Hour)
	if IsStaleCheck(recent) {
		t.Error("1 hour ago should not be stale")
	}

	// Old time is stale
	old := time.Now().Add(-25 * time.Hour)
	if !IsStaleCheck(old) {
		t.Error("25 hours ago should be stale")
	}

	// Exactly 24 hours is not stale (boundary)
	boundary := time.Now().Add(-24 * time.Hour)
	// This might be flaky due to timing, so we just check it doesn't panic
	_ = IsStaleCheck(boundary)
}

// TestCheckSSHKeyHealth tests SSH key health checking
func TestCheckSSHKeyHealth(t *testing.T) {
	// Test with non-existent path
	indicators := CheckSSHKeyHealth("/nonexistent/path/to/key")
	if indicators.SSHKeyExists {
		t.Error("Expected SSHKeyExists to be false for non-existent path")
	}
	if indicators.SSHKeyValid == nil || *indicators.SSHKeyValid {
		t.Error("Expected SSHKeyValid to be false for non-existent path")
	}

	// Test with empty path
	indicators = CheckSSHKeyHealth("")
	if indicators.SSHKeyExists {
		t.Error("Expected SSHKeyExists to be false for empty path")
	}
}

// TestCheckTokenHealth tests token health checking
func TestCheckTokenHealth(t *testing.T) {
	// Test with nil token
	indicators := CheckTokenHealth(nil, "github")
	if indicators.TokenValid == nil || *indicators.TokenValid {
		t.Error("Expected TokenValid to be false for nil token")
	}

	// Test with empty token
	emptyToken := &config.TokenConfig{Username: "user", Token: ""}
	indicators = CheckTokenHealth(emptyToken, "github")
	if indicators.TokenValid == nil || *indicators.TokenValid {
		t.Error("Expected TokenValid to be false for empty token")
	}

	// Test with valid token (should be unknown since we can't verify without API)
	validToken := &config.TokenConfig{Username: "user", Token: "ghp_xxxx"}
	indicators = CheckTokenHealth(validToken, "github")
	if indicators.TokenValid != nil {
		t.Error("Expected TokenValid to be nil (unknown) for token that can't be verified")
	}
}

// TestGetAccountHealth tests account health retrieval
func TestGetAccountHealth(t *testing.T) {
	// Account without SSH or Token
	acc := config.Account{Name: "test"}
	indicators := GetAccountHealth(acc, nil)
	if indicators.SSHKeyExists {
		t.Error("Expected SSHKeyExists to be false for account without SSH")
	}

	// Account with SSH
	acc = config.Account{
		Name: "test",
		SSH:  &config.SshConfig{KeyPath: "/nonexistent/key"},
	}
	indicators = GetAccountHealth(acc, nil)
	if indicators.SSHKeyExists {
		t.Error("Expected SSHKeyExists to be false for non-existent key")
	}

	// With cached health status
	lastChecked := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)
	trueVal := true
	healthStatus := &config.HealthStatus{
		AccountName: "test",
		TokenValid:  &trueVal,
		LastChecked: lastChecked,
	}
	indicators = GetAccountHealth(acc, healthStatus)
	if indicators.TokenValid == nil || !*indicators.TokenValid {
		t.Error("Expected TokenValid from cached status")
	}
	if indicators.IsStale {
		t.Error("Expected not stale for recent health check")
	}
}

// TestFormatHealthDisplay tests health display formatting
func TestFormatHealthDisplay(t *testing.T) {
	trueVal := true
	falseVal := false

	// SSH valid
	indicators := HealthIndicators{
		SSHKeyValid: &trueVal,
	}
	display := FormatHealthDisplay(indicators)
	if display != "SSH:✓" {
		t.Errorf("Expected 'SSH:✓', got '%s'", display)
	}

	// SSH invalid
	indicators = HealthIndicators{
		SSHKeyValid: &falseVal,
	}
	display = FormatHealthDisplay(indicators)
	if display != "SSH:✗" {
		t.Errorf("Expected 'SSH:✗', got '%s'", display)
	}

	// Both SSH and Token
	indicators = HealthIndicators{
		SSHKeyValid: &trueVal,
		TokenValid:  &trueVal,
	}
	display = FormatHealthDisplay(indicators)
	if display != "SSH:✓ Token:✓" {
		t.Errorf("Expected 'SSH:✓ Token:✓', got '%s'", display)
	}

	// Stale indicator
	indicators = HealthIndicators{
		SSHKeyValid: &trueVal,
		IsStale:     true,
		LastChecked: time.Now().Add(-25 * time.Hour),
	}
	display = FormatHealthDisplay(indicators)
	if display != "SSH:✓ (stale)" {
		t.Errorf("Expected 'SSH:✓ (stale)', got '%s'", display)
	}

	// Unknown (no indicators)
	indicators = HealthIndicators{}
	display = FormatHealthDisplay(indicators)
	if display != HealthUnknown {
		t.Errorf("Expected '%s', got '%s'", HealthUnknown, display)
	}
}

// TestStaleThreshold tests stale threshold constant
func TestStaleThreshold(t *testing.T) {
	expected := 24 * time.Hour
	if StaleThreshold != expected {
		t.Errorf("StaleThreshold should be %v, got %v", expected, StaleThreshold)
	}
}
