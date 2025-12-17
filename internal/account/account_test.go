package account

import (
	"testing"

	"github.com/dwirx/ghex/internal/config"
)

// TestNewManager tests manager creation
func TestNewManager(t *testing.T) {
	cfg := config.NewAppConfig()
	manager := NewManager(cfg)

	if manager == nil {
		t.Fatal("Expected manager to be created")
	}

	if manager.cfg != cfg {
		t.Error("Expected manager to hold the config reference")
	}
}

// TestAddAccount tests adding accounts
func TestAddAccount(t *testing.T) {
	cfg := config.NewAppConfig()
	manager := NewManager(cfg)

	acc := config.Account{
		Name:        "test-account",
		GitUserName: "Test User",
		GitEmail:    "test@example.com",
	}

	err := manager.Add(acc)
	if err != nil {
		t.Fatalf("Failed to add account: %v", err)
	}

	if len(cfg.Accounts) != 1 {
		t.Errorf("Expected 1 account, got %d", len(cfg.Accounts))
	}

	if cfg.Accounts[0].Name != "test-account" {
		t.Errorf("Expected account name 'test-account', got '%s'", cfg.Accounts[0].Name)
	}
}

// TestAddDuplicateAccount tests that duplicate names are rejected
func TestAddDuplicateAccount(t *testing.T) {
	cfg := config.NewAppConfig()
	manager := NewManager(cfg)

	acc1 := config.Account{Name: "test-account"}
	acc2 := config.Account{Name: "Test-Account"} // Case-insensitive duplicate

	err := manager.Add(acc1)
	if err != nil {
		t.Fatalf("Failed to add first account: %v", err)
	}

	err = manager.Add(acc2)
	if err == nil {
		t.Error("Expected error when adding duplicate account name")
	}
}

// TestFindAccount tests finding accounts by name
func TestFindAccount(t *testing.T) {
	cfg := config.NewAppConfig()
	manager := NewManager(cfg)

	acc := config.Account{
		Name:        "find-test",
		GitUserName: "Find User",
	}
	manager.Add(acc)

	// Test exact match
	found := manager.Find("find-test")
	if found == nil {
		t.Fatal("Expected to find account")
	}
	if found.GitUserName != "Find User" {
		t.Errorf("Expected GitUserName 'Find User', got '%s'", found.GitUserName)
	}

	// Test case-insensitive match
	found = manager.Find("FIND-TEST")
	if found == nil {
		t.Fatal("Expected to find account with case-insensitive search")
	}

	// Test not found
	found = manager.Find("nonexistent")
	if found != nil {
		t.Error("Expected nil for nonexistent account")
	}
}

// TestRemoveAccount tests removing accounts
func TestRemoveAccount(t *testing.T) {
	cfg := config.NewAppConfig()
	manager := NewManager(cfg)

	manager.Add(config.Account{Name: "account1"})
	manager.Add(config.Account{Name: "account2"})
	manager.Add(config.Account{Name: "account3"})

	if len(cfg.Accounts) != 3 {
		t.Fatalf("Expected 3 accounts, got %d", len(cfg.Accounts))
	}

	// Remove middle account
	err := manager.Remove("account2")
	if err != nil {
		t.Fatalf("Failed to remove account: %v", err)
	}

	if len(cfg.Accounts) != 2 {
		t.Errorf("Expected 2 accounts after removal, got %d", len(cfg.Accounts))
	}

	// Verify account2 is gone
	if manager.Find("account2") != nil {
		t.Error("Expected account2 to be removed")
	}

	// Verify others remain
	if manager.Find("account1") == nil {
		t.Error("Expected account1 to remain")
	}
	if manager.Find("account3") == nil {
		t.Error("Expected account3 to remain")
	}

	// Test removing nonexistent
	err = manager.Remove("nonexistent")
	if err == nil {
		t.Error("Expected error when removing nonexistent account")
	}
}

// TestUpdateAccount tests updating accounts
func TestUpdateAccount(t *testing.T) {
	cfg := config.NewAppConfig()
	manager := NewManager(cfg)

	manager.Add(config.Account{
		Name:        "update-test",
		GitUserName: "Original Name",
		GitEmail:    "original@example.com",
	})

	updated := config.Account{
		Name:        "update-test",
		GitUserName: "Updated Name",
		GitEmail:    "updated@example.com",
	}

	err := manager.Update("update-test", updated)
	if err != nil {
		t.Fatalf("Failed to update account: %v", err)
	}

	found := manager.Find("update-test")
	if found.GitUserName != "Updated Name" {
		t.Errorf("Expected GitUserName 'Updated Name', got '%s'", found.GitUserName)
	}
	if found.GitEmail != "updated@example.com" {
		t.Errorf("Expected GitEmail 'updated@example.com', got '%s'", found.GitEmail)
	}

	// Test updating nonexistent
	err = manager.Update("nonexistent", updated)
	if err == nil {
		t.Error("Expected error when updating nonexistent account")
	}
}

// TestListAccounts tests listing all accounts
func TestListAccounts(t *testing.T) {
	cfg := config.NewAppConfig()
	manager := NewManager(cfg)

	// Empty list
	list := manager.List()
	if len(list) != 0 {
		t.Errorf("Expected empty list, got %d accounts", len(list))
	}

	// Add accounts
	manager.Add(config.Account{Name: "acc1"})
	manager.Add(config.Account{Name: "acc2"})

	list = manager.List()
	if len(list) != 2 {
		t.Errorf("Expected 2 accounts, got %d", len(list))
	}
}

// TestLogActivity tests activity logging
func TestLogActivity(t *testing.T) {
	cfg := config.NewAppConfig()
	manager := NewManager(cfg)

	entry := config.ActivityLogEntry{
		Action:      "test",
		AccountName: "test-account",
		Success:     true,
	}

	manager.LogActivity(entry)

	if len(cfg.ActivityLog) != 1 {
		t.Fatalf("Expected 1 activity log entry, got %d", len(cfg.ActivityLog))
	}

	if cfg.ActivityLog[0].Timestamp == "" {
		t.Error("Expected timestamp to be set automatically")
	}
}

// TestGetRecentActivity tests getting recent activity
func TestGetRecentActivity(t *testing.T) {
	cfg := config.NewAppConfig()
	manager := NewManager(cfg)

	// Add 5 activities
	for i := 1; i <= 5; i++ {
		manager.LogActivity(config.ActivityLogEntry{
			Action:      "test",
			AccountName: "account",
		})
	}

	// Get last 3
	recent := manager.GetRecentActivity(3)
	if len(recent) != 3 {
		t.Errorf("Expected 3 recent activities, got %d", len(recent))
	}

	// Get more than available
	recent = manager.GetRecentActivity(10)
	if len(recent) != 5 {
		t.Errorf("Expected 5 activities (all available), got %d", len(recent))
	}

	// Get with 0 limit
	recent = manager.GetRecentActivity(0)
	if len(recent) != 5 {
		t.Errorf("Expected all activities with 0 limit, got %d", len(recent))
	}
}

// TestSwitchMethodConstants tests switch method constants
func TestSwitchMethodConstants(t *testing.T) {
	if MethodSSH != "ssh" {
		t.Errorf("Expected MethodSSH to be 'ssh', got '%s'", MethodSSH)
	}
	if MethodToken != "token" {
		t.Errorf("Expected MethodToken to be 'token', got '%s'", MethodToken)
	}
}
