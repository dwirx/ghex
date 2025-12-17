package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewVersionCmd creates the version command
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			showVersion()
		},
	}
}

func showVersion() {
	fmt.Printf("ghex v%s\n", Version)
	fmt.Println("Beautiful GitHub Account Switcher & Universal Downloader")
	fmt.Println("Interactive CLI tool for managing multiple GitHub accounts per repository")
	fmt.Println()
	fmt.Println("GitHub: https://github.com/dwirx/ghex")
}
