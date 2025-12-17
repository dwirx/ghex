package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/dwirx/ghex/internal/account"
	"github.com/dwirx/ghex/internal/config"
	"github.com/dwirx/ghex/internal/git"
	"github.com/dwirx/ghex/internal/ssh"
	"github.com/dwirx/ghex/internal/ui"
	"github.com/spf13/cobra"
)

// NewListCmd creates the list command
func NewListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all configured accounts",
		Run: func(cmd *cobra.Command, args []string) {
			runList()
		},
	}
}

// NewStatusCmd creates the status command
func NewStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show current repository status",
		Run: func(cmd *cobra.Command, args []string) {
			runStatus()
		},
	}
}

// NewSwitchCmd creates the switch command
func NewSwitchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "switch [account]",
		Short: "Switch to a specific account",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				runSwitchTo(args[0])
			} else {
				runSwitch()
			}
		},
	}
}

// NewAddCmd creates the add command
func NewAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add",
		Short: "Add a new account",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, _ := config.Load()
			runAddAccount(cfg)
		},
	}
}

// NewRemoveCmd creates the remove command
func NewRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove [account]",
		Short: "Remove an account",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cfg, _ := config.Load()
			runRemoveAccount(cfg)
		},
	}
}

// NewEditCmd creates the edit command
func NewEditCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "edit [account]",
		Short: "Edit an account",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cfg, _ := config.Load()
			runEditAccount(cfg)
		},
	}
}

func runStatus() {
	cfg, _ := config.Load()
	cwd, _ := os.Getwd()

	if !git.IsGitRepo(cwd) {
		ui.ShowError("Not in a git repository")
		return
	}

	fmt.Println()
	fmt.Println(ui.Primary("📊 Repository Status"))
	ui.ShowSeparator()

	manager := account.NewManager(cfg)
	
	// Use enhanced detection with scoring
	matchScore, _ := manager.DetectActiveWithScore(cwd)
	remoteInfo, _ := account.GetRemoteInfo(cwd)
	userName, userEmail, _ := git.GetCurrentUser(cwd)
	branch, _ := git.GetCurrentBranch(cwd)

	if remoteInfo != nil {
		// Show platform with icon
		platformDisplay := account.GetPlatformDisplay(remoteInfo.Platform, "")
		ui.ShowKeyValue("Repository", remoteInfo.RepoPath)
		ui.ShowKeyValue("Remote URL", remoteInfo.RemoteURL)
		ui.ShowKeyValue("Auth Type", strings.ToUpper(remoteInfo.AuthType))
		ui.ShowKeyValue("Platform", platformDisplay)
	}

	fmt.Println()
	fmt.Println(ui.Primary("👤 Git Identity"))
	ui.ShowSeparator()
	ui.ShowKeyValue("Name", userName)
	ui.ShowKeyValue("Email", userEmail)

	fmt.Println()
	fmt.Println(ui.Primary("🔐 Active Account"))
	ui.ShowSeparator()
	if matchScore != nil && matchScore.IsActive {
		ui.ShowKeyValue("Account", ui.Success(matchScore.AccountName))
		ui.ShowKeyValue("Confidence", fmt.Sprintf("%d%% (%s)", matchScore.Score, strings.Join(matchScore.MatchedFields, ", ")))
	} else {
		ui.ShowWarning("No matching account detected")
		if userName != "" || userEmail != "" {
			ui.ShowInfo(fmt.Sprintf("Current identity: %s <%s>", userName, userEmail))
		}
	}

	if branch != "" {
		fmt.Println()
		fmt.Println(ui.Primary("🌿 Current Branch"))
		ui.ShowSeparator()
		ui.ShowKeyValue("Branch", branch)
	}
}

func runList() {
	cfg, err := config.Load()
	if err != nil {
		ui.ShowError(fmt.Sprintf("Failed to load config: %v", err))
		return
	}

	if len(cfg.Accounts) == 0 {
		fmt.Println()
		fmt.Println(ui.RenderEmptyAccountList())
		return
	}

	fmt.Println()
	fmt.Println(ui.Primary("📋 Configured Accounts"))
	ui.ShowSeparator()

	manager := account.NewManager(cfg)
	cwd, _ := os.Getwd()
	activeAccount, _ := manager.DetectActive(cwd)

	// Build health status map
	healthStatuses := make(map[string]*config.HealthStatus)
	for i := range cfg.HealthChecks {
		healthStatuses[cfg.HealthChecks[i].AccountName] = &cfg.HealthChecks[i]
	}

	// Render enhanced table
	fmt.Println()
	fmt.Print(ui.RenderAccountTable(cfg.Accounts, activeAccount, healthStatuses))

	fmt.Println()
	fmt.Println(ui.RenderAccountSummary(len(cfg.Accounts), activeAccount))
}

func runSwitch() {
	cfg, err := config.Load()
	if err != nil {
		ui.ShowError(fmt.Sprintf("Failed to load config: %v", err))
		return
	}

	cwd, _ := os.Getwd()
	if !git.IsGitRepo(cwd) {
		ui.ShowError("Not in a git repository")
		return
	}

	if len(cfg.Accounts) == 0 {
		ui.ShowWarning("No accounts configured")
		return
	}

	// Build account items for selector with platform icons
	manager := account.NewManager(cfg)
	activeAccount, _ := manager.DetectActive(cwd)

	items := make([]ui.SelectorItem, len(cfg.Accounts))
	for i, acc := range cfg.Accounts {
		methods := []string{}
		if acc.SSH != nil {
			methods = append(methods, "🔑SSH")
		}
		if acc.Token != nil {
			methods = append(methods, "🔐Token")
		}

		// Get platform icon
		platformType := account.PlatformGitHub
		if acc.Platform != nil && acc.Platform.Type != "" {
			platformType = acc.Platform.Type
		}
		platformIcon := account.GetPlatformIcon(platformType)

		desc := platformIcon + " " + strings.Join(methods, " ")
		if acc.Name == activeAccount {
			desc = "✓ ACTIVE • " + desc
		}

		items[i] = ui.SelectorItem{
			Title:       acc.Name,
			Description: desc,
			Value:       acc.Name,
		}
	}

	// Run interactive selector
	idx, err := ui.RunSelector("Select Account (↑/k ↓/j to navigate, enter/l to select)", items)
	if err != nil {
		ui.ShowError(fmt.Sprintf("Selection error: %v", err))
		return
	}

	if idx < 0 {
		ui.ShowInfo("Cancelled")
		return
	}

	acc := cfg.Accounts[idx]

	// Select method if both available
	method := account.MethodSSH
	if acc.SSH != nil && acc.Token != nil {
		methodStr, err := ui.SelectMethodInteractive(acc.SSH != nil, acc.Token != nil)
		if err != nil {
			ui.ShowError(fmt.Sprintf("Selection error: %v", err))
			return
		}
		if methodStr == "" {
			ui.ShowInfo("Cancelled")
			return
		}
		if methodStr == "token" {
			method = account.MethodToken
		}
	} else if acc.Token != nil {
		method = account.MethodToken
	}

	if err := manager.Switch(acc.Name, method, cwd); err != nil {
		ui.ShowError(fmt.Sprintf("Failed to switch account: %v", err))
		return
	}

	if err := config.Save(cfg); err != nil {
		ui.ShowWarning(fmt.Sprintf("Failed to save config: %v", err))
	}

	ui.ShowSuccess(fmt.Sprintf("Switched to account: %s (%s)", acc.Name, method))
}

func runSwitchTo(accountName string) {
	cfg, err := config.Load()
	if err != nil {
		ui.ShowError(fmt.Sprintf("Failed to load config: %v", err))
		return
	}

	cwd, _ := os.Getwd()
	if !git.IsGitRepo(cwd) {
		ui.ShowError("Not in a git repository")
		return
	}

	manager := account.NewManager(cfg)
	acc := manager.Find(accountName)
	if acc == nil {
		ui.ShowError(fmt.Sprintf("Account '%s' not found", accountName))
		return
	}

	method := account.MethodSSH
	if acc.SSH == nil && acc.Token != nil {
		method = account.MethodToken
	}

	if err := manager.Switch(acc.Name, method, cwd); err != nil {
		ui.ShowError(fmt.Sprintf("Failed to switch account: %v", err))
		return
	}

	if err := config.Save(cfg); err != nil {
		ui.ShowWarning(fmt.Sprintf("Failed to save config: %v", err))
	}

	ui.ShowSuccess(fmt.Sprintf("Switched to account: %s", acc.Name))
}

func runAddAccount(cfg *config.AppConfig) {
	ui.ShowSection("Add Account")

	name := ui.Prompt("Account label (e.g., work, personal)")
	if name == "" {
		ui.ShowError("Account name is required")
		return
	}

	// Validate for duplicate name early
	validator := account.NewDuplicateValidator(cfg.Accounts)
	if validator.CheckNameDuplicate(name) {
		ui.ShowError(fmt.Sprintf("Account with name '%s' already exists", name))
		return
	}

	gitUserName := ui.Prompt("Git user.name (optional)")
	gitEmail := ui.Prompt("Git user.email (optional)")

	// Interactive platform selection with icons
	platformItems := []ui.SelectorItem{
		{Title: account.IconGitHub + " GitHub", Description: "github.com", Value: account.PlatformGitHub},
		{Title: account.IconGitLab + " GitLab", Description: "gitlab.com", Value: account.PlatformGitLab},
		{Title: account.IconBitbucket + " Bitbucket", Description: "bitbucket.org", Value: account.PlatformBitbucket},
		{Title: account.IconGitea + " Gitea", Description: "Self-hosted Gitea", Value: account.PlatformGitea},
		{Title: account.IconCodeberg + " Codeberg", Description: "codeberg.org", Value: account.PlatformCodeberg},
		{Title: account.IconOther + " Other", Description: "Other Git platform", Value: account.PlatformOther},
	}

	platformIdx, err := ui.RunSelector("Select Platform", platformItems)
	if err != nil {
		ui.ShowError(fmt.Sprintf("Selection error: %v", err))
		return
	}
	if platformIdx < 0 {
		ui.ShowInfo("Cancelled")
		return
	}

	platformType := platformItems[platformIdx].Value

	// Prompt for custom domain if needed
	customDomain := ""
	if platformType == account.PlatformGitea || platformType == account.PlatformOther {
		customDomain = ui.Prompt("Custom domain (e.g., git.company.com)")
	}

	// Interactive method selection
	methodItems := []ui.SelectorItem{
		{Title: "🔑 SSH only", Description: "Use SSH key authentication", Value: "1"},
		{Title: "🔐 Token only", Description: "Use Personal Access Token", Value: "2"},
		{Title: "🔑🔐 Both", Description: "Configure both SSH and Token", Value: "3"},
	}

	methodIdx, err := ui.RunSelector("Select Authentication Method", methodItems)
	if err != nil {
		ui.ShowError(fmt.Sprintf("Selection error: %v", err))
		return
	}
	if methodIdx < 0 {
		ui.ShowInfo("Cancelled")
		return
	}

	methodChoice := methodItems[methodIdx].Value

	acc := config.Account{
		Name:        name,
		GitUserName: gitUserName,
		GitEmail:    gitEmail,
		Platform:    &config.PlatformConfig{Type: platformType, Domain: customDomain},
	}

	if methodChoice == "1" || methodChoice == "3" {
		// Show existing SSH keys for selection
		keys, _ := ssh.ListPrivateKeys()
		if len(keys) > 0 {
			fmt.Println()
			ui.ShowInfo("Existing SSH keys found. Select one or enter a new path:")

			keyItems := make([]ui.SelectorItem, len(keys)+1)
			for i, key := range keys {
				keyItems[i] = ui.SelectorItem{
					Title: key,
					Value: key,
				}
			}
			keyItems[len(keys)] = ui.SelectorItem{
				Title:       "📝 Enter custom path",
				Description: "Type a new SSH key path",
				Value:       "__custom__",
			}

			keyIdx, err := ui.RunSelector("Select SSH Key", keyItems)
			if err == nil && keyIdx >= 0 {
				selectedKey := keyItems[keyIdx].Value
				
				// Check for SSH key duplicate
				if selectedKey != "__custom__" {
					if conflictAcc := validator.CheckSSHKeyDuplicate(selectedKey); conflictAcc != nil {
						ui.ShowWarning(fmt.Sprintf("SSH key is already used by account '%s'", conflictAcc.Name))
						if !ui.Confirm("Continue anyway?") {
							ui.ShowInfo("Cancelled")
							return
						}
					}
				}
				
				if selectedKey == "__custom__" {
					acc.SSH = &config.SshConfig{
						KeyPath:   ui.PromptWithDefault("SSH key path", fmt.Sprintf("~/.ssh/id_ed25519_%s", name)),
						HostAlias: ui.PromptWithDefault("SSH host alias", fmt.Sprintf("%s-%s", platformType, name)),
					}
				} else {
					acc.SSH = &config.SshConfig{
						KeyPath:   selectedKey,
						HostAlias: ui.PromptWithDefault("SSH host alias", fmt.Sprintf("%s-%s", platformType, name)),
					}
				}
			}
		} else {
			acc.SSH = &config.SshConfig{
				KeyPath:   ui.PromptWithDefault("SSH key path", fmt.Sprintf("~/.ssh/id_ed25519_%s", name)),
				HostAlias: ui.PromptWithDefault("SSH host alias", fmt.Sprintf("%s-%s", platformType, name)),
			}
		}
	}

	if methodChoice == "2" || methodChoice == "3" {
		username := ui.Prompt(fmt.Sprintf("%s username", account.GetPlatformName(platformType)))
		
		// Check for token username duplicate
		if conflictAcc := validator.CheckTokenDuplicate(username, platformType); conflictAcc != nil {
			ui.ShowWarning(fmt.Sprintf("Token username '%s' is already used by account '%s' on %s", 
				username, conflictAcc.Name, platformType))
			if !ui.Confirm("Continue anyway?") {
				ui.ShowInfo("Cancelled")
				return
			}
		}
		
		token := ui.PromptPassword("Personal Access Token")
		acc.Token = &config.TokenConfig{
			Username: username,
			Token:    token,
		}
	}

	// Check for email duplicate
	if gitEmail != "" {
		if conflictAcc := validator.CheckEmailDuplicate(gitEmail, platformType); conflictAcc != nil {
			ui.ShowWarning(fmt.Sprintf("Email '%s' is already used by account '%s' on %s", 
				gitEmail, conflictAcc.Name, platformType))
		}
	}

	manager := account.NewManager(cfg)
	if err := manager.Add(acc); err != nil {
		ui.ShowError(fmt.Sprintf("Failed to add account: %v", err))
		return
	}

	if err := config.Save(cfg); err != nil {
		ui.ShowError(fmt.Sprintf("Failed to save config: %v", err))
		return
	}

	ui.ShowSuccess(fmt.Sprintf("Account '%s' added successfully", name))
}

func runEditAccount(cfg *config.AppConfig) {
	if len(cfg.Accounts) == 0 {
		ui.ShowWarning("No accounts to edit")
		return
	}

	// Build items for selector
	items := make([]ui.SelectorItem, len(cfg.Accounts))
	for i, acc := range cfg.Accounts {
		desc := ""
		if acc.GitEmail != "" {
			desc = acc.GitEmail
		}
		items[i] = ui.SelectorItem{
			Title:       acc.Name,
			Description: desc,
			Value:       acc.Name,
		}
	}

	idx, err := ui.RunSelector("Select Account to Edit", items)
	if err != nil {
		ui.ShowError(fmt.Sprintf("Selection error: %v", err))
		return
	}
	if idx < 0 {
		ui.ShowInfo("Cancelled")
		return
	}

	acc := &cfg.Accounts[idx]

	fmt.Println()
	acc.Name = ui.PromptWithDefault("Account label", acc.Name)
	acc.GitUserName = ui.PromptWithDefault("Git user.name", acc.GitUserName)
	acc.GitEmail = ui.PromptWithDefault("Git user.email", acc.GitEmail)

	if err := config.Save(cfg); err != nil {
		ui.ShowError(fmt.Sprintf("Failed to save config: %v", err))
		return
	}

	ui.ShowSuccess(fmt.Sprintf("Account '%s' updated", acc.Name))
}

func runRemoveAccount(cfg *config.AppConfig) {
	if len(cfg.Accounts) == 0 {
		ui.ShowWarning("No accounts to remove")
		return
	}

	// Build items for selector
	items := make([]ui.SelectorItem, len(cfg.Accounts))
	for i, acc := range cfg.Accounts {
		desc := ""
		if acc.GitEmail != "" {
			desc = acc.GitEmail
		}
		items[i] = ui.SelectorItem{
			Title:       acc.Name,
			Description: desc,
			Value:       acc.Name,
		}
	}

	idx, err := ui.RunSelector("Select Account to Remove", items)
	if err != nil {
		ui.ShowError(fmt.Sprintf("Selection error: %v", err))
		return
	}
	if idx < 0 {
		ui.ShowInfo("Cancelled")
		return
	}

	acc := cfg.Accounts[idx]
	fmt.Println()
	if !ui.Confirm(fmt.Sprintf("Remove account '%s'?", acc.Name)) {
		ui.ShowInfo("Cancelled")
		return
	}

	manager := account.NewManager(cfg)
	if err := manager.Remove(acc.Name); err != nil {
		ui.ShowError(fmt.Sprintf("Failed to remove account: %v", err))
		return
	}

	if err := config.Save(cfg); err != nil {
		ui.ShowError(fmt.Sprintf("Failed to save config: %v", err))
		return
	}

	ui.ShowSuccess(fmt.Sprintf("Account '%s' removed", acc.Name))
}
