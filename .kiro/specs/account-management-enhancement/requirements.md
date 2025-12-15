# Requirements Document

## Introduction

Dokumen ini mendefinisikan requirements untuk penyempurnaan fitur manajemen akun pada aplikasi GHEX (GitHub Account Switcher & Universal Downloader). Fitur ini bertujuan untuk meningkatkan tampilan daftar akun, menandai akun aktif dengan lebih jelas, mencegah duplikasi akun, dan memastikan kompatibilitas dengan berbagai platform Git (GitHub, GitLab, Bitbucket, Gitea, dll).

## Glossary

- **GHEX**: Aplikasi CLI untuk mengelola multiple akun Git dan melakukan operasi download/clone
- **Account**: Konfigurasi akun Git yang berisi informasi identitas, SSH key, dan/atau token
- **Active Account**: Akun yang sedang digunakan pada repository Git saat ini berdasarkan git identity dan remote URL
- **Platform**: Layanan hosting Git seperti GitHub, GitLab, Bitbucket, atau Gitea
- **SSH Config**: Konfigurasi autentikasi menggunakan SSH key
- **Token Config**: Konfigurasi autentikasi menggunakan Personal Access Token (PAT)
- **Duplicate Account**: Akun yang memiliki kombinasi identik dari nama, email, SSH key path, atau token username pada platform yang sama

## Requirements

### Requirement 1

**User Story:** As a user, I want to see a clear and organized list of my configured accounts, so that I can easily identify and manage them.

#### Acceptance Criteria

1. WHEN a user runs the list command THEN the GHEX System SHALL display all accounts in a formatted table with columns for status, name, platform, git identity, and authentication methods
2. WHEN displaying account information THEN the GHEX System SHALL show the platform type with an appropriate icon (🐙 GitHub, 🦊 GitLab, 🪣 Bitbucket, 🍵 Gitea)
3. WHEN an account has both SSH and Token configured THEN the GHEX System SHALL display both authentication methods with their respective indicators
4. WHEN the account list is empty THEN the GHEX System SHALL display a helpful message guiding the user to add an account

### Requirement 2

**User Story:** As a user, I want to clearly see which account is currently active, so that I know which identity will be used for Git operations.

#### Acceptance Criteria

1. WHEN displaying the account list THEN the GHEX System SHALL mark the active account with a prominent visual indicator (✓ ACTIVE) and highlight styling
2. WHEN detecting the active account THEN the GHEX System SHALL match based on git user.name, git user.email, and remote URL authentication type
3. WHEN no account matches the current repository configuration THEN the GHEX System SHALL display "No matching account detected" with the current git identity
4. WHEN the user is not in a git repository THEN the GHEX System SHALL display the account list without active status indicators

### Requirement 3

**User Story:** As a user, I want the system to prevent duplicate accounts, so that I maintain a clean and organized configuration.

#### Acceptance Criteria

1. WHEN a user attempts to add an account with an existing name THEN the GHEX System SHALL reject the addition and display an error message
2. WHEN a user attempts to add an account with identical git email on the same platform THEN the GHEX System SHALL warn the user about potential duplicate
3. WHEN a user attempts to add an account with the same SSH key path THEN the GHEX System SHALL warn the user that the key is already associated with another account
4. WHEN a user attempts to add an account with the same token username on the same platform THEN the GHEX System SHALL warn the user about potential duplicate
5. WHEN validating for duplicates THEN the GHEX System SHALL perform case-insensitive comparison for names and emails

### Requirement 4

**User Story:** As a user, I want the system to work seamlessly with different Git platforms, so that I can manage accounts across GitHub, GitLab, Bitbucket, and other services.

#### Acceptance Criteria

1. WHEN adding a new account THEN the GHEX System SHALL allow selection from supported platforms (GitHub, GitLab, Bitbucket, Gitea, Other)
2. WHEN a platform requires a custom domain THEN the GHEX System SHALL prompt for the domain URL
3. WHEN switching accounts THEN the GHEX System SHALL configure the correct SSH host and remote URL format for the target platform
4. WHEN detecting the active account THEN the GHEX System SHALL consider the platform type from the remote URL
5. WHEN building remote URLs THEN the GHEX System SHALL use the correct format for each platform (SSH and HTTPS variants)

### Requirement 5

**User Story:** As a user, I want to see account health status at a glance, so that I can quickly identify any authentication issues.

#### Acceptance Criteria

1. WHEN displaying the account list THEN the GHEX System SHALL show health indicators for SSH key validity (✓ valid, ✗ invalid, ? unknown)
2. WHEN displaying the account list THEN the GHEX System SHALL show health indicators for token validity (✓ valid, ✗ invalid/expired, ? unknown)
3. WHEN an SSH key file does not exist THEN the GHEX System SHALL display a warning indicator next to the account
4. WHEN health check data is stale (older than 24 hours) THEN the GHEX System SHALL indicate that a refresh is recommended

### Requirement 6

**User Story:** As a user, I want to serialize and deserialize account configuration reliably, so that my settings persist correctly across sessions.

#### Acceptance Criteria

1. WHEN saving account configuration THEN the GHEX System SHALL serialize the data to JSON format
2. WHEN loading account configuration THEN the GHEX System SHALL deserialize the JSON data and reconstruct account objects
3. WHEN serializing account data THEN the GHEX System SHALL preserve all fields including optional SSH, Token, and Platform configurations
4. WHEN deserializing account data THEN the GHEX System SHALL validate the structure and handle missing optional fields gracefully
5. WHEN printing account configuration for debugging THEN the GHEX System SHALL format the output as valid JSON that can be parsed back

