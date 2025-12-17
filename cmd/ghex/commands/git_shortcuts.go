package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/dwirx/ghex/internal/git"
	"github.com/dwirx/ghex/internal/shell"
	"github.com/dwirx/ghex/internal/ui"
	"github.com/spf13/cobra"
)

// AddGitShortcuts adds all git shortcut commands
func AddGitShortcuts(rootCmd *cobra.Command) {
	// Git status
	rootCmd.AddCommand(&cobra.Command{
		Use:   "gs",
		Short: "git status",
		Run: func(cmd *cobra.Command, args []string) {
			shell.RunInteractive("git", "status")
		},
	})

	// Git branch
	rootCmd.AddCommand(&cobra.Command{
		Use:   "gb",
		Short: "git branch",
		Run: func(cmd *cobra.Command, args []string) {
			shell.RunInteractive("git", "branch")
		},
	})

	// Git branch -a
	rootCmd.AddCommand(&cobra.Command{
		Use:   "gba",
		Short: "git branch -a",
		Run: func(cmd *cobra.Command, args []string) {
			shell.RunInteractive("git", "branch", "-a")
		},
	})

	// Git branch -r
	rootCmd.AddCommand(&cobra.Command{
		Use:   "gbr",
		Short: "git branch -r",
		Run: func(cmd *cobra.Command, args []string) {
			shell.RunInteractive("git", "branch", "-r")
		},
	})

	// Git fetch
	rootCmd.AddCommand(&cobra.Command{
		Use:   "gf",
		Short: "git fetch origin",
		Run: func(cmd *cobra.Command, args []string) {
			ui.ShowInfo("Fetching from origin...")
			shell.RunInteractive("git", "fetch", "origin")
			ui.ShowSuccess("Fetch completed")
		},
	})

	// Git pull
	rootCmd.AddCommand(&cobra.Command{
		Use:   "gp",
		Short: "git pull",
		Run: func(cmd *cobra.Command, args []string) {
			ui.ShowInfo("Pulling from remote...")
			shell.RunInteractive("git", "pull")
			ui.ShowSuccess("Pull completed")
		},
	})

	// Git pull --rebase
	rootCmd.AddCommand(&cobra.Command{
		Use:   "gpr",
		Short: "git pull --rebase",
		Run: func(cmd *cobra.Command, args []string) {
			ui.ShowInfo("Pulling with rebase...")
			shell.RunInteractive("git", "pull", "--rebase")
			ui.ShowSuccess("Pull completed")
		},
	})

	// Git checkout
	rootCmd.AddCommand(&cobra.Command{
		Use:   "gco [branch]",
		Short: "git checkout",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			shell.RunInteractive("git", "checkout", args[0])
		},
	})

	// Git checkout -b
	rootCmd.AddCommand(&cobra.Command{
		Use:   "gcb [branch]",
		Short: "git checkout -b (create new branch)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			shell.RunInteractive("git", "checkout", "-b", args[0])
		},
	})

	// Git log
	rootCmd.AddCommand(&cobra.Command{
		Use:   "gl",
		Short: "git log --oneline",
		Run: func(cmd *cobra.Command, args []string) {
			shell.RunInteractive("git", "log", "--oneline", "-20")
		},
	})

	// Git diff
	rootCmd.AddCommand(&cobra.Command{
		Use:   "gd",
		Short: "git diff",
		Run: func(cmd *cobra.Command, args []string) {
			shell.RunInteractive("git", "diff")
		},
	})

	// Git diff staged
	rootCmd.AddCommand(&cobra.Command{
		Use:   "gds",
		Short: "git diff --staged",
		Run: func(cmd *cobra.Command, args []string) {
			shell.RunInteractive("git", "diff", "--staged")
		},
	})

	// Shove command (add, commit, push)
	rootCmd.AddCommand(&cobra.Command{
		Use:   "shove [message]",
		Short: "git add, commit, and push",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			message := strings.Join(args, " ")
			runShove(message, false)
		},
	})

	// Shove no-confirm
	rootCmd.AddCommand(&cobra.Command{
		Use:   "shovenc [message]",
		Short: "git add, commit, and push (no confirm)",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			message := strings.Join(args, " ")
			runShove(message, true)
		},
	})

	// Set name
	rootCmd.AddCommand(&cobra.Command{
		Use:   "setname [name]",
		Short: "Set global git user.name",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := strings.Join(args, " ")
			if err := git.SetGlobalIdentity(name, ""); err != nil {
				ui.ShowError(fmt.Sprintf("Failed: %v", err))
				return
			}
			ui.ShowSuccess(fmt.Sprintf("Git user.name set to: %s", name))
		},
	})

	// Set email
	rootCmd.AddCommand(&cobra.Command{
		Use:   "setmail [email]",
		Short: "Set global git user.email",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := git.SetGlobalIdentity("", args[0]); err != nil {
				ui.ShowError(fmt.Sprintf("Failed: %v", err))
				return
			}
			ui.ShowSuccess(fmt.Sprintf("Git user.email set to: %s", args[0]))
		},
	})

	// Show config
	rootCmd.AddCommand(&cobra.Command{
		Use:   "showconfig",
		Short: "Show git configuration",
		Run: func(cmd *cobra.Command, args []string) {
			output, err := git.GetConfigList()
			if err != nil {
				ui.ShowError(fmt.Sprintf("Failed: %v", err))
				return
			}
			fmt.Println(output)
		},
	})

	// Git stash
	rootCmd.AddCommand(&cobra.Command{
		Use:   "gst",
		Short: "git stash",
		Run: func(cmd *cobra.Command, args []string) {
			shell.RunInteractive("git", "stash")
		},
	})

	// Git stash pop
	rootCmd.AddCommand(&cobra.Command{
		Use:   "gstp",
		Short: "git stash pop",
		Run: func(cmd *cobra.Command, args []string) {
			shell.RunInteractive("git", "stash", "pop")
		},
	})

	// Git reset
	rootCmd.AddCommand(&cobra.Command{
		Use:   "greset",
		Short: "git reset HEAD",
		Run: func(cmd *cobra.Command, args []string) {
			shell.RunInteractive("git", "reset", "HEAD")
		},
	})
}

func runShove(message string, noConfirm bool) {
	cwd, _ := os.Getwd()

	if !git.IsGitRepo(cwd) {
		ui.ShowError("Not in a git repository")
		return
	}

	if message == "" {
		ui.ShowError("Commit message is required")
		return
	}

	// Git add
	ui.ShowInfo("Adding files...")
	if _, err := shell.Run("git", "add", "."); err != nil {
		ui.ShowError(fmt.Sprintf("Failed to add files: %v", err))
		return
	}
	ui.ShowSuccess("Files added")

	// Git commit
	ui.ShowInfo(fmt.Sprintf("Committing: %s", message))
	if _, err := shell.Run("git", "commit", "-m", message); err != nil {
		ui.ShowError(fmt.Sprintf("Failed to commit: %v", err))
		return
	}
	ui.ShowSuccess("Committed successfully")

	// Push
	if noConfirm || ui.Confirm("Push to origin?") {
		ui.ShowInfo("Pushing to origin...")
		if err := shell.RunInteractive("git", "push", "origin"); err != nil {
			ui.ShowError(fmt.Sprintf("Failed to push: %v", err))
			return
		}
		ui.ShowSuccess("Pushed successfully")
	} else {
		ui.ShowWarning("Push cancelled")
	}
}
