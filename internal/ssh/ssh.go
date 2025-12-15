package ssh

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dwirx/ghex/internal/platform"
	"github.com/dwirx/ghex/internal/shell"
)

// GenerateKey generates a new Ed25519 SSH key pair
func GenerateKey(keyPath, comment string) error {
	// Expand path
	keyPath = platform.ExpandPath(keyPath)

	// Ensure directory exists
	dir := filepath.Dir(keyPath)
	if err := platform.EnsureDir(dir, 0700); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Generate key using ssh-keygen
	args := []string{
		"-t", "ed25519",
		"-f", keyPath,
		"-C", comment,
		"-N", "", // Empty passphrase
	}

	_, err := shell.Run("ssh-keygen", args...)
	if err != nil {
		return fmt.Errorf("failed to generate SSH key: %w", err)
	}

	// Set permissions
	if err := SetKeyPermissions(keyPath); err != nil {
		return fmt.Errorf("failed to set key permissions: %w", err)
	}

	return nil
}

// ImportKey copies an SSH private key to a new location
func ImportKey(srcPath, destPath string) error {
	srcPath = platform.ExpandPath(srcPath)
	destPath = platform.ExpandPath(destPath)

	// Check source exists
	if !platform.FileExists(srcPath) {
		return fmt.Errorf("source key not found: %s", srcPath)
	}

	// Ensure destination directory exists
	dir := filepath.Dir(destPath)
	if err := platform.EnsureDir(dir, 0700); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Copy file
	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source key: %w", err)
	}
	defer src.Close()

	dst, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to create destination key: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to copy key: %w", err)
	}

	// Set permissions
	return SetKeyPermissions(destPath)
}

// EnsurePublicKey generates a public key from a private key if it doesn't exist
func EnsurePublicKey(privateKeyPath string) (string, error) {
	privateKeyPath = platform.ExpandPath(privateKeyPath)
	pubPath := privateKeyPath + ".pub"

	// Check if public key already exists
	if platform.FileExists(pubPath) {
		return pubPath, nil
	}

	// Generate public key from private key
	output, err := shell.Run("ssh-keygen", "-y", "-f", privateKeyPath)
	if err != nil {
		return "", fmt.Errorf("failed to generate public key: %w", err)
	}

	// Write public key
	pubKey := strings.TrimSpace(output) + "\n"
	if err := os.WriteFile(pubPath, []byte(pubKey), 0644); err != nil {
		return "", fmt.Errorf("failed to write public key: %w", err)
	}

	return pubPath, nil
}

// SetKeyPermissions sets proper permissions on SSH key files
func SetKeyPermissions(keyPath string) error {
	keyPath = platform.ExpandPath(keyPath)

	// Set private key permissions (600)
	if err := os.Chmod(keyPath, 0600); err != nil {
		// On Windows, chmod might not work as expected
		if !platform.IsWindows() {
			return fmt.Errorf("failed to set private key permissions: %w", err)
		}
	}

	// Set public key permissions (644) if exists
	pubPath := keyPath + ".pub"
	if platform.FileExists(pubPath) {
		if err := os.Chmod(pubPath, 0644); err != nil {
			if !platform.IsWindows() {
				return fmt.Errorf("failed to set public key permissions: %w", err)
			}
		}
	}

	return nil
}

// TestConnection tests SSH connection to a host
func TestConnection(host string) (bool, string, error) {
	if host == "" {
		host = "github.com"
	}

	args := []string{
		"-T",
		"-o", "StrictHostKeyChecking=no",
		"-o", "ConnectTimeout=10",
		"-o", "BatchMode=yes",
		fmt.Sprintf("git@%s", host),
	}

	output, err := shell.Exec("ssh", args...)

	// SSH -T returns exit code 1 for successful auth on GitHub
	// Check output for success patterns
	successPatterns := []string{
		"successfully authenticated",
		"Hi .+! You've successfully authenticated",
		"Welcome to GitLab",
		"logged in as",
		"authenticated via",
		"You can use git",
		"Hi there,",
	}

	for _, pattern := range successPatterns {
		matched, _ := regexp.MatchString("(?i)"+pattern, output)
		if matched {
			// Extract username if possible
			userRe := regexp.MustCompile(`Hi\s+([^!,]+)[!,]|logged in as\s+(\S+)|@(\S+)|Hi there,?\s+([^!]+)!`)
			if matches := userRe.FindStringSubmatch(output); len(matches) > 1 {
				for _, m := range matches[1:] {
					if m != "" {
						return true, fmt.Sprintf("Successfully authenticated as %s", m), nil
					}
				}
			}
			return true, "Successfully authenticated", nil
		}
	}

	if err != nil {
		return false, output, err
	}

	return false, output, nil
}

// ListPrivateKeys returns a list of SSH private keys in the SSH directory
func ListPrivateKeys() ([]string, error) {
	sshDir := platform.GetSSHDir()

	entries, err := os.ReadDir(sshDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	excludeFiles := map[string]bool{
		"known_hosts":       true,
		"known_hosts.old":   true,
		"config":            true,
		"authorized_keys":   true,
		"authorized_keys2":  true,
	}

	var keys []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		// Skip excluded files
		if excludeFiles[name] {
			continue
		}

		// Skip public keys
		if strings.HasSuffix(name, ".pub") {
			continue
		}

		keys = append(keys, filepath.Join(sshDir, name))
	}

	return keys, nil
}

// SuggestKeyFilenames suggests destination filenames for SSH keys
func SuggestKeyFilenames(username, label string) []string {
	base := username
	if base == "" {
		base = label
	}
	if base == "" {
		base = "github"
	}

	// Clean the base name
	base = strings.ToLower(base)
	base = regexp.MustCompile(`[^a-zA-Z0-9_-]+`).ReplaceAllString(base, "")

	candidates := []string{
		fmt.Sprintf("id_ed25519_%s", base),
		fmt.Sprintf("id_ecdsa_%s", base),
		fmt.Sprintf("id_rsa_%s", base),
		"id_ed25519_github",
	}

	// Remove duplicates
	seen := make(map[string]bool)
	var unique []string
	for _, c := range candidates {
		if !seen[c] {
			seen[c] = true
			unique = append(unique, c)
		}
	}

	return unique
}
