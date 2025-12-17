package account

import (
	"testing"
)

// TestGetPlatformInfo tests platform info retrieval
func TestGetPlatformInfo(t *testing.T) {
	tests := []struct {
		platform     string
		expectedIcon string
		expectedName string
	}{
		{PlatformGitHub, IconGitHub, "GitHub"},
		{PlatformGitLab, IconGitLab, "GitLab"},
		{PlatformBitbucket, IconBitbucket, "Bitbucket"},
		{PlatformGitea, IconGitea, "Gitea"},
		{PlatformOther, IconOther, "Other"},
		{"unknown", IconOther, "Other"}, // Unknown defaults to Other
		{"GITHUB", IconGitHub, "GitHub"}, // Case insensitive
	}

	for _, tt := range tests {
		t.Run(tt.platform, func(t *testing.T) {
			info := GetPlatformInfo(tt.platform)
			if info.Icon != tt.expectedIcon {
				t.Errorf("GetPlatformInfo(%s).Icon = %s, expected %s", tt.platform, info.Icon, tt.expectedIcon)
			}
			if info.Name != tt.expectedName {
				t.Errorf("GetPlatformInfo(%s).Name = %s, expected %s", tt.platform, info.Name, tt.expectedName)
			}
		})
	}
}

// TestGetPlatformIcon tests icon retrieval
func TestGetPlatformIcon(t *testing.T) {
	if GetPlatformIcon("github") != IconGitHub {
		t.Error("Expected GitHub icon")
	}
	if GetPlatformIcon("gitlab") != IconGitLab {
		t.Error("Expected GitLab icon")
	}
	if GetPlatformIcon("unknown") != IconOther {
		t.Error("Expected Other icon for unknown platform")
	}
}

// TestGetPlatformName tests name retrieval
func TestGetPlatformName(t *testing.T) {
	if GetPlatformName("github") != "GitHub" {
		t.Error("Expected 'GitHub'")
	}
	if GetPlatformName("gitlab") != "GitLab" {
		t.Error("Expected 'GitLab'")
	}
}

// TestGetPlatformDisplay tests display string generation
func TestGetPlatformDisplay(t *testing.T) {
	// Without custom domain
	display := GetPlatformDisplay("github", "")
	if display != IconGitHub+" GitHub" {
		t.Errorf("Expected '%s GitHub', got '%s'", IconGitHub, display)
	}

	// With custom domain
	display = GetPlatformDisplay("gitea", "git.company.com")
	if display != IconGitea+" git.company.com" {
		t.Errorf("Expected '%s git.company.com', got '%s'", IconGitea, display)
	}
}

// TestDetectPlatformFromURL tests platform detection from URLs
func TestDetectPlatformFromURL(t *testing.T) {
	tests := []struct {
		url      string
		expected string
	}{
		// GitHub
		{"https://github.com/user/repo.git", PlatformGitHub},
		{"git@github.com:user/repo.git", PlatformGitHub},
		{"ssh://git@github.com/user/repo.git", PlatformGitHub},

		// GitLab
		{"https://gitlab.com/user/repo.git", PlatformGitLab},
		{"git@gitlab.com:user/repo.git", PlatformGitLab},

		// Bitbucket
		{"https://bitbucket.org/user/repo.git", PlatformBitbucket},
		{"git@bitbucket.org:user/repo.git", PlatformBitbucket},

		// Gitea / Codeberg
		{"https://codeberg.org/user/repo.git", PlatformGitea},
		{"https://gitea.example.com/user/repo.git", PlatformGitea},

		// Other
		{"https://custom.git.server/user/repo.git", PlatformOther},
		{"git@custom.server:user/repo.git", PlatformOther},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			result := DetectPlatformFromURL(tt.url)
			if result != tt.expected {
				t.Errorf("DetectPlatformFromURL(%s) = %s, expected %s", tt.url, result, tt.expected)
			}
		})
	}
}

// TestGetSupportedPlatforms tests supported platforms list
func TestGetSupportedPlatforms(t *testing.T) {
	platforms := GetSupportedPlatforms()

	if len(platforms) != 5 {
		t.Errorf("Expected 5 supported platforms, got %d", len(platforms))
	}

	// Check all expected platforms are present
	expected := map[string]bool{
		PlatformGitHub:    false,
		PlatformGitLab:    false,
		PlatformBitbucket: false,
		PlatformGitea:     false,
		PlatformOther:     false,
	}

	for _, p := range platforms {
		if _, ok := expected[p]; ok {
			expected[p] = true
		}
	}

	for p, found := range expected {
		if !found {
			t.Errorf("Expected platform %s not found in supported list", p)
		}
	}
}

// TestIsValidPlatform tests platform validation
func TestIsValidPlatform(t *testing.T) {
	tests := []struct {
		platform string
		expected bool
	}{
		{"github", true},
		{"gitlab", true},
		{"bitbucket", true},
		{"gitea", true},
		{"other", true},
		{"GITHUB", true}, // Case insensitive
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.platform, func(t *testing.T) {
			result := IsValidPlatform(tt.platform)
			if result != tt.expected {
				t.Errorf("IsValidPlatform(%s) = %v, expected %v", tt.platform, result, tt.expected)
			}
		})
	}
}

// TestPlatformConstants tests platform constants
func TestPlatformConstants(t *testing.T) {
	if PlatformGitHub != "github" {
		t.Error("PlatformGitHub should be 'github'")
	}
	if PlatformGitLab != "gitlab" {
		t.Error("PlatformGitLab should be 'gitlab'")
	}
	if PlatformBitbucket != "bitbucket" {
		t.Error("PlatformBitbucket should be 'bitbucket'")
	}
	if PlatformGitea != "gitea" {
		t.Error("PlatformGitea should be 'gitea'")
	}
	if PlatformOther != "other" {
		t.Error("PlatformOther should be 'other'")
	}
}

// TestIconConstants tests icon constants
func TestIconConstants(t *testing.T) {
	if IconGitHub != "🐙" {
		t.Error("IconGitHub should be 🐙")
	}
	if IconGitLab != "🦊" {
		t.Error("IconGitLab should be 🦊")
	}
	if IconBitbucket != "🪣" {
		t.Error("IconBitbucket should be 🪣")
	}
	if IconGitea != "🍵" {
		t.Error("IconGitea should be 🍵")
	}
	if IconOther != "🔗" {
		t.Error("IconOther should be 🔗")
	}
}
