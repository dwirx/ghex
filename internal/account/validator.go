package account

import (
	"strings"

	"github.com/dwirx/ghex/internal/config"
)

// ValidationResult contains duplicate check results
type ValidationResult struct {
	IsValid  bool
	Errors   []string // Hard errors (must fix)
	Warnings []string // Soft warnings (can proceed)
}

// DuplicateValidator validates account uniqueness
type DuplicateValidator struct {
	accounts []config.Account
}

// NewDuplicateValidator creates a new validator with existing accounts
func NewDuplicateValidator(accounts []config.Account) *DuplicateValidator {
	return &DuplicateValidator{accounts: accounts}
}

// ValidateNew checks if a new account would create duplicates
func (v *DuplicateValidator) ValidateNew(account config.Account) ValidationResult {
	result := ValidationResult{
		IsValid:  true,
		Errors:   []string{},
		Warnings: []string{},
	}

	// Check name duplicate (hard error)
	if v.CheckNameDuplicate(account.Name) {
		result.IsValid = false
		result.Errors = append(result.Errors, "Account name '"+account.Name+"' already exists")
	}

	// Get platform type
	platformType := "github"
	if account.Platform != nil && account.Platform.Type != "" {
		platformType = account.Platform.Type
	}

	// Check email duplicate on same platform (warning)
	if account.GitEmail != "" {
		if conflictAcc := v.CheckEmailDuplicate(account.GitEmail, platformType); conflictAcc != nil {
			result.Warnings = append(result.Warnings,
				"Email '"+account.GitEmail+"' is already used by account '"+conflictAcc.Name+"' on "+platformType)
		}
	}

	// Check SSH key duplicate (warning)
	if account.SSH != nil && account.SSH.KeyPath != "" {
		if conflictAcc := v.CheckSSHKeyDuplicate(account.SSH.KeyPath); conflictAcc != nil {
			result.Warnings = append(result.Warnings,
				"SSH key '"+account.SSH.KeyPath+"' is already used by account '"+conflictAcc.Name+"'")
		}
	}

	// Check token username duplicate on same platform (warning)
	if account.Token != nil && account.Token.Username != "" {
		if conflictAcc := v.CheckTokenDuplicate(account.Token.Username, platformType); conflictAcc != nil {
			result.Warnings = append(result.Warnings,
				"Token username '"+account.Token.Username+"' is already used by account '"+conflictAcc.Name+"' on "+platformType)
		}
	}

	return result
}

// CheckNameDuplicate checks for duplicate account names (case-insensitive)
func (v *DuplicateValidator) CheckNameDuplicate(name string) bool {
	for _, acc := range v.accounts {
		if strings.EqualFold(acc.Name, name) {
			return true
		}
	}
	return false
}

// CheckEmailDuplicate checks for duplicate email on same platform
func (v *DuplicateValidator) CheckEmailDuplicate(email, platform string) *config.Account {
	for i, acc := range v.accounts {
		if strings.EqualFold(acc.GitEmail, email) {
			accPlatform := "github"
			if acc.Platform != nil && acc.Platform.Type != "" {
				accPlatform = acc.Platform.Type
			}
			if strings.EqualFold(accPlatform, platform) {
				return &v.accounts[i]
			}
		}
	}
	return nil
}

// CheckSSHKeyDuplicate checks if SSH key is already used
func (v *DuplicateValidator) CheckSSHKeyDuplicate(keyPath string) *config.Account {
	normalizedPath := normalizePath(keyPath)
	for i, acc := range v.accounts {
		if acc.SSH != nil && acc.SSH.KeyPath != "" {
			if normalizePath(acc.SSH.KeyPath) == normalizedPath {
				return &v.accounts[i]
			}
		}
	}
	return nil
}

// CheckTokenDuplicate checks for duplicate token username on same platform
func (v *DuplicateValidator) CheckTokenDuplicate(username, platform string) *config.Account {
	for i, acc := range v.accounts {
		if acc.Token != nil && strings.EqualFold(acc.Token.Username, username) {
			accPlatform := "github"
			if acc.Platform != nil && acc.Platform.Type != "" {
				accPlatform = acc.Platform.Type
			}
			if strings.EqualFold(accPlatform, platform) {
				return &v.accounts[i]
			}
		}
	}
	return nil
}

// normalizePath normalizes file paths for comparison
func normalizePath(path string) string {
	// Convert to lowercase and normalize separators
	path = strings.ToLower(path)
	path = strings.ReplaceAll(path, "\\", "/")
	return path
}

// EqualFoldStrings checks if two strings are equal (case-insensitive)
func EqualFoldStrings(a, b string) bool {
	return strings.EqualFold(a, b)
}
