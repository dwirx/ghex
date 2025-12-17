package account

import (
	"strings"
)

// Platform type constants
const (
	PlatformGitHub    = "github"
	PlatformGitLab    = "gitlab"
	PlatformBitbucket = "bitbucket"
	PlatformGitea     = "gitea"
	PlatformCodeberg  = "codeberg"
	PlatformOther     = "other"
)

// Platform icons
const (
	IconGitHub    = "🐙"
	IconGitLab    = "🦊"
	IconBitbucket = "🪣"
	IconGitea     = "🍵"
	IconCodeberg  = "🏔️"
	IconOther     = "🔗"
)

// PlatformInfo contains display information for a platform
type PlatformInfo struct {
	Type   string
	Icon   string
	Name   string
	Domain string
}

// platformRegistry holds all supported platforms
var platformRegistry = map[string]PlatformInfo{
	PlatformGitHub: {
		Type:   PlatformGitHub,
		Icon:   IconGitHub,
		Name:   "GitHub",
		Domain: "github.com",
	},
	PlatformGitLab: {
		Type:   PlatformGitLab,
		Icon:   IconGitLab,
		Name:   "GitLab",
		Domain: "gitlab.com",
	},
	PlatformBitbucket: {
		Type:   PlatformBitbucket,
		Icon:   IconBitbucket,
		Name:   "Bitbucket",
		Domain: "bitbucket.org",
	},
	PlatformGitea: {
		Type:   PlatformGitea,
		Icon:   IconGitea,
		Name:   "Gitea",
		Domain: "",
	},
	PlatformCodeberg: {
		Type:   PlatformCodeberg,
		Icon:   IconCodeberg,
		Name:   "Codeberg",
		Domain: "codeberg.org",
	},
	PlatformOther: {
		Type:   PlatformOther,
		Icon:   IconOther,
		Name:   "Other",
		Domain: "",
	},
}

// GetPlatformInfo returns display info for a platform type
func GetPlatformInfo(platformType string) PlatformInfo {
	platformType = strings.ToLower(platformType)
	if info, ok := platformRegistry[platformType]; ok {
		return info
	}
	return platformRegistry[PlatformOther]
}

// GetPlatformIcon returns the icon for a platform type
func GetPlatformIcon(platformType string) string {
	return GetPlatformInfo(platformType).Icon
}

// GetPlatformName returns the display name for a platform type
func GetPlatformName(platformType string) string {
	return GetPlatformInfo(platformType).Name
}

// GetPlatformDisplay returns formatted platform display string with icon
func GetPlatformDisplay(platformType string, customDomain string) string {
	info := GetPlatformInfo(platformType)
	if customDomain != "" {
		return info.Icon + " " + customDomain
	}
	if info.Domain != "" {
		return info.Icon + " " + info.Name
	}
	return info.Icon + " " + info.Name
}

// DetectPlatformFromURL identifies platform type from remote URL
func DetectPlatformFromURL(url string) string {
	url = strings.ToLower(url)

	// Check for known domains
	if strings.Contains(url, "github.com") || strings.Contains(url, "github:") {
		return PlatformGitHub
	}
	if strings.Contains(url, "gitlab.com") || strings.Contains(url, "gitlab:") {
		return PlatformGitLab
	}
	if strings.Contains(url, "bitbucket.org") || strings.Contains(url, "bitbucket:") {
		return PlatformBitbucket
	}
	if strings.Contains(url, "codeberg.org") {
		return PlatformCodeberg
	}
	if strings.Contains(url, "gitea") {
		return PlatformGitea
	}

	return PlatformOther
}

// GetSupportedPlatforms returns list of supported platform types
func GetSupportedPlatforms() []string {
	return []string{
		PlatformGitHub,
		PlatformGitLab,
		PlatformBitbucket,
		PlatformGitea,
		PlatformCodeberg,
		PlatformOther,
	}
}

// IsValidPlatform checks if a platform type is valid
func IsValidPlatform(platformType string) bool {
	_, ok := platformRegistry[strings.ToLower(platformType)]
	return ok
}
