package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/dwirx/ghex/internal/ui"
	"github.com/dwirx/ghex/internal/uninstall"
	"github.com/spf13/cobra"
)

// NewUninstallCmd creates the uninstall command
func NewUninstallCmd() *cobra.Command {
	var force bool
	var purge bool
	var keepConfig bool
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall GHEX from your system",
		Long:  "Remove GHEX binary and optionally configuration files from your system",
		Run: func(cmd *cobra.Command, args []string) {
			runUninstall(force, purge, keepConfig, dryRun)
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation prompts")
	cmd.Flags().BoolVarP(&purge, "purge", "p", false, "Remove configuration files as well")
	cmd.Flags().BoolVar(&keepConfig, "keep-config", false, "Keep configuration files (default when not using --purge)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be removed without actually removing")

	return cmd
}

func runUninstall(force, purge, keepConfig, dryRun bool) {
	svc := uninstall.NewService()

	// Show banner
	ui.ShowSection("GHEX Uninstaller")
	fmt.Println()

	// Get preview
	preview := svc.GetPreview()

	// Show what will be removed
	fmt.Println("The following will be removed:")
	fmt.Println()

	if svc.BinaryExists() {
		fmt.Printf("  Binary: %s\n", preview.BinaryPath)
	} else {
		fmt.Printf("  Binary: (not found)\n")
	}

	if svc.ConfigExists() {
		fmt.Printf("  Config: %s\n", preview.ConfigPath)
		if preview.LegacyConfig != "" {
			fmt.Printf("  Legacy Config: %s\n", preview.LegacyConfig)
		}
	}

	if preview.PathEntry != "" {
		fmt.Printf("  PATH entry: %s\n", preview.PathEntry)
	}

	fmt.Println()

	// Dry run - just show preview and exit
	if dryRun {
		ui.ShowInfo("Dry run mode - no files will be removed")
		return
	}

	// Confirm uninstallation
	if !force {
		if !confirm("Do you want to uninstall GHEX?") {
			ui.ShowInfo("Uninstallation cancelled")
			return
		}
	}

	// Ask about config removal if not specified
	removeConfig := purge
	if !purge && !keepConfig && !force {
		fmt.Println()
		removeConfig = confirm("Do you want to remove configuration files as well?")
	}

	// Execute uninstallation
	opts := uninstall.Options{
		Force:      force,
		Purge:      removeConfig,
		KeepConfig: keepConfig && !removeConfig,
		DryRun:     false,
	}

	result := svc.Execute(opts)

	// Show results
	fmt.Println()

	if result.BinaryRemoved {
		ui.ShowSuccess("Binary removed")
	} else if !svc.BinaryExists() {
		ui.ShowInfo("Binary was not installed")
	} else {
		ui.ShowError("Failed to remove binary")
		fmt.Println()
		fmt.Println(svc.GetManualRemovalInstructions())
	}

	if result.ConfigRemoved {
		ui.ShowSuccess("Configuration files removed")
	} else if opts.Purge && svc.ConfigExists() {
		ui.ShowWarning("Some configuration files could not be removed")
	} else if !opts.Purge {
		ui.ShowInfo("Configuration files preserved")
	}

	if result.PathUpdated {
		ui.ShowSuccess("PATH updated")
	}

	// Show errors if any
	if len(result.Errors) > 0 {
		fmt.Println()
		ui.ShowWarning("Some operations failed:")
		for _, err := range result.Errors {
			fmt.Printf("  - %s\n", err)
		}
	}

	// Final message
	fmt.Println()
	if result.Success || result.BinaryRemoved {
		ui.ShowSuccess("GHEX has been uninstalled!")
		fmt.Println()
		fmt.Println("Thank you for using GHEX! 👋")
	} else {
		ui.ShowError("Uninstallation incomplete")
		fmt.Println()
		fmt.Println(svc.GetManualRemovalInstructions())
	}
}

func confirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [y/N]: ", prompt)
	
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}
