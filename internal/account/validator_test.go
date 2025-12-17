package account

import (
	"testing"

	"github.com/dwirx/ghex/internal/config"
)

// TestNewDuplicateValidator tests validator creation
func TestNewDuplicateValidator(t *testing.T) {
	accounts := []config.Account{
		{Name: "test1"},
		{Name: "test2"},
	}

	validator := NewDuplicateValidator(accounts)
	if validator == nil {
		t.Fatal("Expected validator to be created")
	}
}

// TestCheckNameDuplicate tests name duplicate checking
func TestCheckNameDuplicate(t *testing.T) {
	accounts := []config.Account{
		{Name: "existing-account"},
		{Name: "Another-Account"},
	}

	validator := NewDuplicateValidator(accounts)

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"exact match", "existing-account", true},
		{"case insensitive", "EXISTING-ACCOUNT", true},
		{"mixed case", "Existing-Account", true},
		{"another account", "another-account", true},
		{"not found", "new-account", false},
		{"partial match", "existing", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.CheckNameDuplicate(tt.input)
			if result != tt.expected {
				t.Errorf("CheckNameDuplicate(%s) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestCheckEmailDuplicate tests email duplicate checking on same platform
func TestCheckEmailDuplicate(t *testing.T) {
	accounts := []config.Account{
		{
			Name:     "github-account",
			GitEmail: "user@example.com",
			Platform: &config.PlatformConfig{Type: "github"},
		},
		{
			Name:     "gitlab-account",
			GitEmail: "user@example.com",
			Platform: &config.PlatformConfig{Type: "gitlab"},
		},
	}

	validator := NewDuplicateValidator(accounts)

	// Same email, same platform - should find duplicate
	result := validator.CheckEmailDuplicate("user@example.com", "github")
	if result == nil {
		t.Error("Expected to find duplicate email on github")
	}
	if result != nil && result.Name != "github-account" {
		t.Errorf("Expected conflict with 'github-account', got '%s'", result.Name)
	}

	// Same email, different platform - should find duplicate on that platform
	result = validator.CheckEmailDuplicate("user@example.com", "gitlab")
	if result == nil {
		t.Error("Expected to find duplicate email on gitlab")
	}

	// Same email, platform not in list - should not find duplicate
	result = validator.CheckEmailDuplicate("user@example.com", "bitbucket")
	if result != nil {
		t.Error("Expected no duplicate on bitbucket")
	}

	// Different email - should not find duplicate
	result = validator.CheckEmailDuplicate("other@example.com", "github")
	if result != nil {
		t.Error("Expected no duplicate for different email")
	}

	// Case insensitive email
	result = validator.CheckEmailDuplicate("USER@EXAMPLE.COM", "github")
	if result == nil {
		t.Error("Expected case-insensitive email match")
	}
}

// TestCheckSSHKeyDuplicate tests SSH key duplicate checking
func TestCheckSSHKeyDuplicate(t *testing.T) {
	accounts := []config.Account{
		{
			Name: "account1",
			SSH:  &config.SshConfig{KeyPath: "~/.ssh/id_ed25519_work"},
		},
		{
			Name: "account2",
			SSH:  &config.SshConfig{KeyPath: "/home/user/.ssh/id_rsa"},
		},
		{
			Name: "account3",
			// No SSH config
		},
	}

	validator := NewDuplicateValidator(accounts)

	// Exact match
	result := validator.CheckSSHKeyDuplicate("~/.ssh/id_ed25519_work")
	if result == nil {
		t.Error("Expected to find duplicate SSH key")
	}

	// Different path
	result = validator.CheckSSHKeyDuplicate("~/.ssh/id_ed25519_personal")
	if result != nil {
		t.Error("Expected no duplicate for different path")
	}

	// Case insensitive path (normalized)
	result = validator.CheckSSHKeyDuplicate("/HOME/USER/.SSH/ID_RSA")
	if result == nil {
		t.Error("Expected case-insensitive path match")
	}
}

// TestCheckTokenDuplicate tests token username duplicate checking
func TestCheckTokenDuplicate(t *testing.T) {
	accounts := []config.Account{
		{
			Name:     "github-work",
			Token:    &config.TokenConfig{Username: "workuser"},
			Platform: &config.PlatformConfig{Type: "github"},
		},
		{
			Name:     "gitlab-work",
			Token:    &config.TokenConfig{Username: "workuser"},
			Platform: &config.PlatformConfig{Type: "gitlab"},
		},
	}

	validator := NewDuplicateValidator(accounts)

	// Same username, same platform
	result := validator.CheckTokenDuplicate("workuser", "github")
	if result == nil {
		t.Error("Expected to find duplicate token username on github")
	}

	// Same username, different platform
	result = validator.CheckTokenDuplicate("workuser", "bitbucket")
	if result != nil {
		t.Error("Expected no duplicate on bitbucket")
	}

	// Different username
	result = validator.CheckTokenDuplicate("otheruser", "github")
	if result != nil {
		t.Error("Expected no duplicate for different username")
	}

	// Case insensitive username
	result = validator.CheckTokenDuplicate("WORKUSER", "github")
	if result == nil {
		t.Error("Expected case-insensitive username match")
	}
}

// TestValidateNew tests full validation of new account
func TestValidateNew(t *testing.T) {
	accounts := []config.Account{
		{
			Name:     "existing",
			GitEmail: "existing@example.com",
			SSH:      &config.SshConfig{KeyPath: "~/.ssh/existing"},
			Token:    &config.TokenConfig{Username: "existinguser"},
			Platform: &config.PlatformConfig{Type: "github"},
		},
	}

	validator := NewDuplicateValidator(accounts)

	// Test duplicate name (error)
	newAcc := config.Account{
		Name:     "existing",
		Platform: &config.PlatformConfig{Type: "github"},
	}
	result := validator.ValidateNew(newAcc)
	if result.IsValid {
		t.Error("Expected validation to fail for duplicate name")
	}
	if len(result.Errors) == 0 {
		t.Error("Expected error message for duplicate name")
	}

	// Test duplicate email (warning)
	newAcc = config.Account{
		Name:     "new-account",
		GitEmail: "existing@example.com",
		Platform: &config.PlatformConfig{Type: "github"},
	}
	result = validator.ValidateNew(newAcc)
	if !result.IsValid {
		t.Error("Expected validation to pass with warning for duplicate email")
	}
	if len(result.Warnings) == 0 {
		t.Error("Expected warning for duplicate email")
	}

	// Test duplicate SSH key (warning)
	newAcc = config.Account{
		Name:     "new-account",
		SSH:      &config.SshConfig{KeyPath: "~/.ssh/existing"},
		Platform: &config.PlatformConfig{Type: "github"},
	}
	result = validator.ValidateNew(newAcc)
	if !result.IsValid {
		t.Error("Expected validation to pass with warning for duplicate SSH key")
	}
	if len(result.Warnings) == 0 {
		t.Error("Expected warning for duplicate SSH key")
	}

	// Test completely new account (no errors or warnings)
	newAcc = config.Account{
		Name:     "brand-new",
		GitEmail: "new@example.com",
		SSH:      &config.SshConfig{KeyPath: "~/.ssh/new"},
		Token:    &config.TokenConfig{Username: "newuser"},
		Platform: &config.PlatformConfig{Type: "gitlab"},
	}
	result = validator.ValidateNew(newAcc)
	if !result.IsValid {
		t.Error("Expected validation to pass for new account")
	}
	if len(result.Errors) > 0 {
		t.Errorf("Expected no errors, got: %v", result.Errors)
	}
	if len(result.Warnings) > 0 {
		t.Errorf("Expected no warnings, got: %v", result.Warnings)
	}
}

// TestEqualFoldStrings tests case-insensitive string comparison
func TestEqualFoldStrings(t *testing.T) {
	tests := []struct {
		a, b     string
		expected bool
	}{
		{"hello", "hello", true},
		{"Hello", "hello", true},
		{"HELLO", "hello", true},
		{"hello", "world", false},
		{"", "", true},
		{"test", "", false},
	}

	for _, tt := range tests {
		result := EqualFoldStrings(tt.a, tt.b)
		if result != tt.expected {
			t.Errorf("EqualFoldStrings(%s, %s) = %v, expected %v", tt.a, tt.b, result, tt.expected)
		}
	}
}
