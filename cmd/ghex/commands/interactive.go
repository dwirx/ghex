package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/dwirx/ghex/internal/account"
	"github.com/dwirx/ghex/internal/config"
	"github.com/dwirx/ghex/internal/git"
	"github.com/dwirx/ghex/internal/ui"
)

func runInteractive() {
	ui.ShowTitle()

	cfg, err := config.Load()
	if err != nil {
		ui.ShowError(fmt.Sprintf("Failed to load config: %v", err))
		return
	}

	for {
		showRepositoryContext(cfg)

		items := []ui.SelectorItem{
			{Title: "🔄 Switch account", Description: "Switch account for current repository", Value: "switch"},
			{Title: "📋 List accounts", Description: "Show all configured accounts", Value: "list"},
			{Title: "➕ Add account", Description: "Add a new GitHub account", Value: "add"},
			{Title: "✏️  Edit account", Description: "Modify an existing account", Value: "edit"},
			{Title: "🗑️  Remove account", Description: "Delete an account", Value: "remove"},
			{Title: "🔑 SSH Management", Description: "Generate, import, or manage SSH keys", Value: "ssh"},
			{Title: "🌐 Switch SSH globally", Description: "Change global SSH configuration", Value: "globalssh"},
			{Title: "📥 Download (dlx)", Description: "Download files from URLs or Git repos", Value: "dlx"},
			{Title: "🧪 Test connection", Description: "Test SSH/Token authentication", Value: "test"},
			{Title: "🏥 Health check", Description: "Check all account connections", Value: "health"},
			{Title: "📜 Activity log", Description: "View recent activity", Value: "log"},
			{Title: "🚪 Exit", Description: "Quit GHEX", Value: "exit"},
		}

		idx, err := ui.RunSelector("Main Menu (↑/k ↓/j navigate, enter/l select, q quit)", items)
		if err != nil || idx < 0 {
			ui.ShowSeparator()
			ui.ShowSuccess("Thank you for using GHEX! 👋")
			return
		}

		fmt.Println()

		switch items[idx].Value {
		case "switch":
			runSwitch()
		case "list":
			runList()
		case "add":
			runAddAccount(cfg)
		case "edit":
			runEditAccount(cfg)
		case "remove":
			runRemoveAccount(cfg)
		case "ssh":
			runSSHMenu(cfg)
		case "globalssh":
			runSwitchGlobalSSH(cfg)
		case "dlx":
			runDlxMenu()
		case "test":
			runTestConnection(cfg)
		case "health":
			runHealthCheck()
		case "log":
			runActivityLog()
		case "exit":
			ui.ShowSeparator()
			ui.ShowSuccess("Thank you for using GHEX! 👋")
			return
		}

		fmt.Println()
		ui.Prompt("Press Enter to continue...")

		// Reload config in case it changed
		cfg, _ = config.Load()
	}
}

func showRepositoryContext(cfg *config.AppConfig) {
	cwd, _ := os.Getwd()

	if !git.IsGitRepo(cwd) {
		ui.ShowBox(ui.Muted("Run ghex inside a Git repository to see active account details."), ui.BoxOptions{
			Title: "Repository Context",
			Type:  "info",
		})
		return
	}

	manager := account.NewManager(cfg)
	activeAccount, _ := manager.DetectActive(cwd)
	remoteInfo, _ := account.GetRemoteInfo(cwd)
	userName, userEmail, _ := git.GetCurrentUser(cwd)

	var lines []string
	if remoteInfo != nil {
		lines = append(lines, fmt.Sprintf("Repository: %s", ui.Accent(remoteInfo.RepoPath)))
		lines = append(lines, fmt.Sprintf("Auth Type: %s", ui.Secondary(strings.ToUpper(remoteInfo.AuthType))))
	}
	if userName != "" {
		lines = append(lines, fmt.Sprintf("Git User: %s", userName))
	}
	if userEmail != "" {
		lines = append(lines, fmt.Sprintf("Git Email: %s", userEmail))
	}
	if activeAccount != "" {
		lines = append(lines, fmt.Sprintf("Active Account: %s", ui.Success(activeAccount)))
	} else {
		lines = append(lines, ui.Warning("Active account could not be detected"))
	}

	boxType := "success"
	if activeAccount == "" {
		boxType = "warning"
	}

	ui.ShowBox(strings.Join(lines, "\n"), ui.BoxOptions{
		Title: "Repository Context",
		Type:  boxType,
	})
}

func isGitURL(s string) bool {
	return strings.HasPrefix(s, "http://") ||
		strings.HasPrefix(s, "https://") ||
		strings.HasPrefix(s, "git@") ||
		strings.HasPrefix(s, "ssh://")
}
