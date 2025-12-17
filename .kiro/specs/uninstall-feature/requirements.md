# Requirements Document

## Introduction

Fitur uninstall untuk GHEX CLI tool yang memungkinkan pengguna menghapus aplikasi GHEX beserta file konfigurasi terkait dari sistem mereka. Fitur ini mendukung platform Linux, macOS, dan Windows dengan opsi untuk menghapus konfigurasi atau mempertahankannya.

## Glossary

- **GHEX**: GitHub Account Switcher & Universal Downloader CLI tool
- **Config Directory**: Direktori yang menyimpan konfigurasi GHEX (biasanya `~/.config/ghex` atau `%APPDATA%\ghex`)
- **Binary**: File executable GHEX yang terinstall di sistem
- **Install Directory**: Lokasi binary GHEX (`/usr/local/bin` untuk Linux/macOS, `%LOCALAPPDATA%\ghex` untuk Windows)

## Requirements

### Requirement 1

**User Story:** As a user, I want to uninstall GHEX from my system, so that I can cleanly remove the application when I no longer need it.

#### Acceptance Criteria

1. WHEN a user runs `ghex uninstall` THEN the System SHALL display a confirmation prompt before proceeding with uninstallation
2. WHEN a user confirms uninstallation THEN the System SHALL remove the GHEX binary from the install directory
3. WHEN a user runs `ghex uninstall --force` THEN the System SHALL skip the confirmation prompt and proceed with uninstallation
4. WHEN uninstallation completes successfully THEN the System SHALL display a success message with details of removed files

### Requirement 2

**User Story:** As a user, I want to optionally remove my configuration files during uninstall, so that I can choose whether to keep my settings for potential reinstallation.

#### Acceptance Criteria

1. WHEN a user runs `ghex uninstall` THEN the System SHALL prompt whether to remove configuration files
2. WHEN a user runs `ghex uninstall --purge` THEN the System SHALL remove both binary and configuration directory
3. WHEN a user runs `ghex uninstall --keep-config` THEN the System SHALL remove only the binary and preserve configuration files
4. WHEN configuration files are removed THEN the System SHALL display the path of removed configuration directory

### Requirement 3

**User Story:** As a user, I want the uninstall process to work on my operating system, so that I can remove GHEX regardless of my platform.

#### Acceptance Criteria

1. WHEN running on Linux or macOS THEN the System SHALL remove binary from `/usr/local/bin/ghex`
2. WHEN running on Windows THEN the System SHALL remove binary from `%LOCALAPPDATA%\ghex\ghex.exe`
3. WHEN running on Windows THEN the System SHALL remove the install directory from PATH environment variable
4. WHEN elevated permissions are required THEN the System SHALL inform the user and provide instructions for manual removal

### Requirement 4

**User Story:** As a user, I want to have an uninstall script available, so that I can uninstall GHEX even if the binary is corrupted or missing.

#### Acceptance Criteria

1. WHEN a user runs the uninstall shell script on Linux/macOS THEN the System SHALL remove GHEX binary and optionally config files
2. WHEN a user runs the uninstall PowerShell script on Windows THEN the System SHALL remove GHEX binary, directory, and PATH entry
3. WHEN the uninstall script runs THEN the System SHALL detect the current platform and use appropriate removal methods
4. WHEN files cannot be removed due to permissions THEN the System SHALL display clear error messages with remediation steps

### Requirement 5

**User Story:** As a user, I want to see what will be removed before uninstalling, so that I can verify the uninstall will not affect other files.

#### Acceptance Criteria

1. WHEN a user runs `ghex uninstall --dry-run` THEN the System SHALL display all files and directories that would be removed without actually removing them
2. WHEN displaying removal preview THEN the System SHALL show binary path, config directory path, and any PATH modifications
3. WHEN dry-run completes THEN the System SHALL not modify any files or system settings
