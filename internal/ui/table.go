package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dwirx/ghex/internal/account"
	"github.com/dwirx/ghex/internal/config"
)

// Table column widths
const (
	colWidthStatus   = 10
	colWidthName     = 15
	colWidthPlatform = 12
	colWidthIdentity = 30
	colWidthAuth     = 12
	colWidthHealth   = 10
)

// AccountTableRow represents a single row in the account table
type AccountTableRow struct {
	Status      string // Active indicator
	Name        string
	Platform    string // With icon
	GitIdentity string // name <email>
	AuthMethods string // SSH/Token indicators
	Health      string // Health status
}

// Table styles
var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(AccentColor).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(MutedColor)

	activeRowStyle = lipgloss.NewStyle().
			Foreground(SuccessColor).
			Bold(true)

	normalRowStyle = lipgloss.NewStyle().
			Foreground(TextColor)

	cellStyle = lipgloss.NewStyle().
			PaddingRight(2)
)

// RenderAccountTable renders accounts in a formatted table
func RenderAccountTable(accounts []config.Account, activeAccount string, healthStatuses map[string]*config.HealthStatus) string {
	if len(accounts) == 0 {
		return RenderEmptyAccountList()
	}

	var sb strings.Builder

	// Render header
	sb.WriteString(renderTableHeader())
	sb.WriteString("\n")

	// Render rows
	for _, acc := range accounts {
		isActive := strings.EqualFold(acc.Name, activeAccount)
		
		// Get health status
		var healthStatus *config.HealthStatus
		if healthStatuses != nil {
			healthStatus = healthStatuses[acc.Name]
		}
		
		row := buildAccountRow(acc, isActive, healthStatus)
		sb.WriteString(renderTableRow(row, isActive))
		sb.WriteString("\n")
	}

	return sb.String()
}

// RenderEmptyAccountList returns message for empty account list
func RenderEmptyAccountList() string {
	return WarningStyle.Render("No accounts configured") + "\n" +
		AccentStyle.Render("ℹ ") + TextStyle.Render("Run 'ghex add' to add your first account")
}

// buildAccountRow creates a table row from account data
func buildAccountRow(acc config.Account, isActive bool, healthStatus *config.HealthStatus) AccountTableRow {
	row := AccountTableRow{}

	// Status
	if isActive {
		row.Status = SuccessStyle.Render("✓ ACTIVE")
	} else {
		row.Status = MutedStyle.Render("○")
	}

	// Name
	row.Name = acc.Name

	// Platform with icon
	platformType := "github"
	customDomain := ""
	if acc.Platform != nil {
		platformType = acc.Platform.Type
		customDomain = acc.Platform.Domain
	}
	row.Platform = account.GetPlatformDisplay(platformType, customDomain)

	// Git Identity
	if acc.GitUserName != "" || acc.GitEmail != "" {
		if acc.GitUserName != "" && acc.GitEmail != "" {
			row.GitIdentity = fmt.Sprintf("%s <%s>", acc.GitUserName, acc.GitEmail)
		} else if acc.GitUserName != "" {
			row.GitIdentity = acc.GitUserName
		} else {
			row.GitIdentity = acc.GitEmail
		}
	} else {
		row.GitIdentity = MutedStyle.Render("-")
	}

	// Auth Methods
	authMethods := []string{}
	if acc.SSH != nil {
		authMethods = append(authMethods, "🔑SSH")
	}
	if acc.Token != nil {
		authMethods = append(authMethods, "🔐Token")
	}
	if len(authMethods) > 0 {
		row.AuthMethods = strings.Join(authMethods, " ")
	} else {
		row.AuthMethods = MutedStyle.Render("-")
	}

	// Health
	health := account.GetAccountHealth(acc, healthStatus)
	row.Health = account.FormatHealthDisplay(health)

	return row
}

// renderTableHeader renders the table header
func renderTableHeader() string {
	headers := []string{
		padRight("STATUS", colWidthStatus),
		padRight("NAME", colWidthName),
		padRight("PLATFORM", colWidthPlatform),
		padRight("GIT IDENTITY", colWidthIdentity),
		padRight("AUTH", colWidthAuth),
		padRight("HEALTH", colWidthHealth),
	}

	headerLine := strings.Join(headers, "")
	return headerStyle.Render(headerLine)
}

// renderTableRow renders a single table row
func renderTableRow(row AccountTableRow, isActive bool) string {
	cells := []string{
		padRight(row.Status, colWidthStatus),
		padRight(row.Name, colWidthName),
		padRight(row.Platform, colWidthPlatform),
		padRight(truncate(row.GitIdentity, colWidthIdentity-2), colWidthIdentity),
		padRight(row.AuthMethods, colWidthAuth),
		padRight(row.Health, colWidthHealth),
	}

	rowLine := strings.Join(cells, "")

	if isActive {
		return activeRowStyle.Render(rowLine)
	}
	return normalRowStyle.Render(rowLine)
}

// padRight pads a string to the right with spaces
func padRight(s string, width int) string {
	// Calculate visible length (excluding ANSI codes)
	visibleLen := visibleLength(s)
	if visibleLen >= width {
		return s
	}
	return s + strings.Repeat(" ", width-visibleLen)
}

// truncate truncates a string to max length with ellipsis
func truncate(s string, maxLen int) string {
	if visibleLength(s) <= maxLen {
		return s
	}
	// Simple truncation - doesn't handle ANSI codes perfectly
	runes := []rune(s)
	if len(runes) > maxLen-3 {
		return string(runes[:maxLen-3]) + "..."
	}
	return s
}

// visibleLength returns the visible length of a string (excluding ANSI codes)
func visibleLength(s string) int {
	// Simple implementation - count runes excluding ANSI escape sequences
	inEscape := false
	length := 0
	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		length++
	}
	return length
}

// RenderAccountSummary renders a summary line for accounts
func RenderAccountSummary(total int, activeAccount string) string {
	summary := fmt.Sprintf("Total accounts: %d", total)
	if activeAccount != "" {
		summary += fmt.Sprintf(" | Active: %s", SuccessStyle.Render(activeAccount))
	}
	return AccentStyle.Render("ℹ ") + TextStyle.Render(summary)
}
