package git

import (
	"fmt"
	"regexp"
	"strings"
)

// URLInfo contains parsed information from a git URL
type URLInfo struct {
	URL      string
	IsSSH    bool
	Host     string
	Owner    string
	Repo     string
	Platform string // github, gitlab, bitbucket, gitea, other
}

// ParseRepoFromURL extracts owner/repo from a git URL
func ParseRepoFromURL(rawURL string) (owner, repo string, err error) {
	if rawURL == "" {
		return "", "", fmt.Errorf("empty URL")
	}

	rawURL = strings.TrimSpace(rawURL)

	// SSH format: git@host:owner/repo.git
	sshPattern := regexp.MustCompile(`^git@([^:]+):(.+?)(?:\.git)?$`)
	if matches := sshPattern.FindStringSubmatch(rawURL); len(matches) == 3 {
		parts := strings.Split(matches[2], "/")
		if len(parts) >= 2 {
			return parts[0], strings.TrimSuffix(parts[len(parts)-1], ".git"), nil
		}
	}

	// SSH format: ssh://git@host/owner/repo.git
	sshURLPattern := regexp.MustCompile(`^ssh://git@([^/]+)/(.+?)(?:\.git)?$`)
	if matches := sshURLPattern.FindStringSubmatch(rawURL); len(matches) == 3 {
		parts := strings.Split(matches[2], "/")
		if len(parts) >= 2 {
			return parts[0], strings.TrimSuffix(parts[len(parts)-1], ".git"), nil
		}
	}

	// HTTPS format: https://host/owner/repo.git
	httpsPattern := regexp.MustCompile(`^https?://([^/]+)/(.+?)(?:\.git)?$`)
	if matches := httpsPattern.FindStringSubmatch(rawURL); len(matches) == 3 {
		parts := strings.Split(matches[2], "/")
		if len(parts) >= 2 {
			return parts[0], strings.TrimSuffix(parts[len(parts)-1], ".git"), nil
		}
	}

	return "", "", fmt.Errorf("unable to parse URL: %s", rawURL)
}

// NormalizeURL normalizes a git URL and adds .git suffix if missing
func NormalizeURL(rawURL string) (normalized string, isSSH bool, err error) {
	if rawURL == "" {
		return "", false, fmt.Errorf("empty URL")
	}

	rawURL = strings.TrimSpace(rawURL)
	rawURL = strings.TrimSuffix(rawURL, "#") // Remove trailing #

	// SSH format: git@host:path or ssh://git@host/path
	if strings.HasPrefix(rawURL, "git@") || strings.HasPrefix(rawURL, "ssh://") {
		if !strings.HasSuffix(rawURL, ".git") {
			rawURL += ".git"
		}
		return rawURL, true, nil
	}

	// HTTPS format
	if strings.HasPrefix(rawURL, "http://") || strings.HasPrefix(rawURL, "https://") {
		if !strings.HasSuffix(rawURL, ".git") {
			rawURL += ".git"
		}
		return rawURL, false, nil
	}

	return "", false, fmt.Errorf("invalid git URL format: %s", rawURL)
}

// ParseURL parses a git URL and returns detailed information
func ParseURL(rawURL string) (*URLInfo, error) {
	normalized, isSSH, err := NormalizeURL(rawURL)
	if err != nil {
		return nil, err
	}

	owner, repo, err := ParseRepoFromURL(normalized)
	if err != nil {
		return nil, err
	}

	host := detectHost(normalized)
	platform := detectPlatform(host)

	return &URLInfo{
		URL:      normalized,
		IsSSH:    isSSH,
		Host:     host,
		Owner:    owner,
		Repo:     repo,
		Platform: platform,
	}, nil
}

// detectHost extracts the host from a URL
func detectHost(rawURL string) string {
	// SSH format: git@host:path
	if strings.HasPrefix(rawURL, "git@") {
		parts := strings.SplitN(rawURL[4:], ":", 2)
		if len(parts) > 0 {
			return parts[0]
		}
	}

	// SSH URL format: ssh://git@host/path
	if strings.HasPrefix(rawURL, "ssh://git@") {
		parts := strings.SplitN(rawURL[10:], "/", 2)
		if len(parts) > 0 {
			return parts[0]
		}
	}

	// HTTPS format
	httpsPattern := regexp.MustCompile(`^https?://([^/]+)`)
	if matches := httpsPattern.FindStringSubmatch(rawURL); len(matches) == 2 {
		return matches[1]
	}

	return ""
}

// detectPlatform detects the git platform from the host
func detectPlatform(host string) string {
	host = strings.ToLower(host)

	if strings.Contains(host, "github") {
		return "github"
	}
	if strings.Contains(host, "gitlab") {
		return "gitlab"
	}
	if strings.Contains(host, "bitbucket") {
		return "bitbucket"
	}
	if strings.Contains(host, "codeberg") {
		return "codeberg"
	}
	if strings.Contains(host, "gitea") {
		return "gitea"
	}

	return "other"
}

// PlatformURLConfig holds URL configuration for a platform
type PlatformURLConfig struct {
	SSHFormat   string // Format string for SSH URL (e.g., "git@%s:%s")
	HTTPSFormat string // Format string for HTTPS URL (e.g., "https://%s/%s")
	DefaultHost string // Default host for the platform
}

// platformURLConfigs holds URL configurations for each platform
var platformURLConfigs = map[string]PlatformURLConfig{
	"github": {
		SSHFormat:   "git@%s:%s",
		HTTPSFormat: "https://%s/%s",
		DefaultHost: "github.com",
	},
	"gitlab": {
		SSHFormat:   "git@%s:%s",
		HTTPSFormat: "https://%s/%s",
		DefaultHost: "gitlab.com",
	},
	"bitbucket": {
		SSHFormat:   "git@%s:%s",
		HTTPSFormat: "https://%s/%s",
		DefaultHost: "bitbucket.org",
	},
	"gitea": {
		SSHFormat:   "git@%s:%s",
		HTTPSFormat: "https://%s/%s",
		DefaultHost: "", // Gitea requires custom domain
	},
	"codeberg": {
		SSHFormat:   "git@%s:%s",
		HTTPSFormat: "https://%s/%s",
		DefaultHost: "codeberg.org",
	},
	"other": {
		SSHFormat:   "git@%s:%s",
		HTTPSFormat: "https://%s/%s",
		DefaultHost: "",
	},
}

// GetPlatformURLConfig returns URL configuration for a platform
func GetPlatformURLConfig(platform string) PlatformURLConfig {
	platform = strings.ToLower(platform)
	if config, ok := platformURLConfigs[platform]; ok {
		return config
	}
	return platformURLConfigs["other"]
}

// BuildRemoteURL builds a remote URL for a given platform
func BuildRemoteURL(platform, domain, repoPath string, useSSH bool) string {
	config := GetPlatformURLConfig(platform)

	if domain == "" {
		domain = config.DefaultHost
		if domain == "" {
			domain = "github.com" // Fallback
		}
	}

	// Ensure repo path has .git suffix
	if !strings.HasSuffix(repoPath, ".git") {
		repoPath += ".git"
	}

	if useSSH {
		return fmt.Sprintf(config.SSHFormat, domain, repoPath)
	}

	return fmt.Sprintf(config.HTTPSFormat, domain, repoPath)
}

// BuildSSHRemoteURL builds an SSH remote URL for a platform
func BuildSSHRemoteURL(platform, domain, repoPath string) string {
	return BuildRemoteURL(platform, domain, repoPath, true)
}

// BuildHTTPSRemoteURL builds an HTTPS remote URL for a platform
func BuildHTTPSRemoteURL(platform, domain, repoPath string) string {
	return BuildRemoteURL(platform, domain, repoPath, false)
}

// GetDefaultDomain returns the default domain for a platform
func GetDefaultDomain(platform string) string {
	config := GetPlatformURLConfig(platform)
	return config.DefaultHost
}

// WithGitSuffix ensures a repo path has .git suffix
func WithGitSuffix(repoPath string) string {
	if strings.HasSuffix(repoPath, ".git") {
		return repoPath
	}
	return repoPath + ".git"
}

// GetPlatformSSHHost returns the SSH host for a platform
func GetPlatformSSHHost(platform, domain string) string {
	if domain != "" {
		return domain
	}

	switch platform {
	case "github":
		return "github.com"
	case "gitlab":
		return "gitlab.com"
	case "bitbucket":
		return "bitbucket.org"
	default:
		return "github.com"
	}
}
