package commands

import (
	"fmt"

	"github.com/dwirx/ghex/internal/account"
	"github.com/dwirx/ghex/internal/config"
	"github.com/dwirx/ghex/internal/git"
	"github.com/dwirx/ghex/internal/ui"
)

func runClone(repoURL, targetDir string) {
	cfg, _ := config.Load()

	ui.ShowTitle()
	ui.ShowInfo(fmt.Sprintf("Cloning: %s", repoURL))

	urlInfo, err := git.ParseURL(repoURL)
	if err != nil {
		ui.ShowError(fmt.Sprintf("Invalid URL: %v", err))
		return
	}

	if len(cfg.Accounts) > 0 {
		fmt.Println(ui.Primary("Select account (or press Enter to skip):"))
		for i, acc := range cfg.Accounts {
			fmt.Printf("  %s %s\n", ui.Dim(fmt.Sprintf("[%d]", i+1)), acc.Name)
		}
		fmt.Printf("  %s Skip account setup\n", ui.Dim("[0]"))

		choice := ui.Prompt("Enter choice")
		var idx int
		_, _ = fmt.Sscanf(choice, "%d", &idx)

		if idx > 0 && idx <= len(cfg.Accounts) {
			acc := cfg.Accounts[idx-1]

			spinner := ui.NewSpinner("Cloning repository...")
			spinner.Start()

			clonedDir, err := git.CloneWithIdentity(repoURL, targetDir, acc.GitUserName, acc.GitEmail)
			if err != nil {
				spinner.StopWithError(fmt.Sprintf("Clone failed: %v", err))
				return
			}

			spinner.StopWithSuccess(fmt.Sprintf("Cloned to: %s", clonedDir))

			manager := account.NewManager(cfg)
			method := account.MethodSSH
			if acc.SSH == nil && acc.Token != nil {
				method = account.MethodToken
			}

			if err := manager.Switch(acc.Name, method, clonedDir); err != nil {
				ui.ShowWarning(fmt.Sprintf("Failed to set up account: %v", err))
			} else {
				ui.ShowSuccess(fmt.Sprintf("Account '%s' configured", acc.Name))
			}

			_ = config.Save(cfg)
			return
		}
	}

	spinner := ui.NewSpinner("Cloning repository...")
	spinner.Start()

	clonedDir, err := git.Clone(repoURL, targetDir)
	if err != nil {
		spinner.StopWithError(fmt.Sprintf("Clone failed: %v", err))
		return
	}

	spinner.StopWithSuccess(fmt.Sprintf("Cloned to: %s", clonedDir))
	ui.ShowInfo(fmt.Sprintf("Repository: %s/%s", urlInfo.Owner, urlInfo.Repo))
}
