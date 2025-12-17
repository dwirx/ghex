# Changelog

All notable changes to GHEX will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Enhanced account management with duplicate validation
- Platform icons for GitHub (🐙), GitLab (🦊), Bitbucket (🪣), Gitea (🍵)
- Active account detection with confidence scoring
- Health indicators for SSH keys and tokens
- Enhanced table display for account listing
- Case-insensitive duplicate checking
- Support for custom domains (self-hosted GitLab, Gitea, etc.)
- Comprehensive test suite

### Changed
- Improved account switching with platform-specific URL handling
- Better error messages and warnings for duplicate accounts
- Enhanced status display with match confidence percentage

### Fixed
- Case-sensitive account name comparison
- SSH key path normalization for duplicate detection

## [1.0.0] - 2024-XX-XX

### Added
- Initial release
- Multi-account management for Git platforms
- SSH and Token authentication support
- Interactive account switching
- Repository status display
- Activity logging
- Support for GitHub, GitLab, Bitbucket, Gitea
- Cross-platform support (Windows, Linux, macOS)

---

## Version History

| Version | Date | Description |
|---------|------|-------------|
| 1.0.0 | TBD | Initial stable release |

## Upgrade Guide

### From 0.x to 1.0

No breaking changes. Simply replace the binary with the new version.

```bash
# Linux/macOS
curl -sSL https://raw.githubusercontent.com/dwirx/ghex/main/scripts/install.sh | bash

# Windows
iwr -useb https://raw.githubusercontent.com/dwirx/ghex/main/scripts/install.ps1 | iex
```
