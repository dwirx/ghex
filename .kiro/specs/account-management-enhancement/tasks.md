# Implementation Plan

## Account Management Enhancement

- [x] 1. Create DuplicateValidator Component

  - [x] 1.1 Create `internal/account/validator.go` with DuplicateValidator struct and ValidationResult type


    - Implement CheckNameDuplicate with case-insensitive comparison
    - Implement CheckEmailDuplicate for same platform check
    - Implement CheckSSHKeyDuplicate for key path check
    - Implement CheckTokenDuplicate for same platform check
    - Implement ValidateNew that aggregates all checks
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5_
  - [ ]* 1.2 Write property test for name uniqueness enforcement
    - **Property 6: Name Uniqueness Enforcement**
    - **Validates: Requirements 3.1**
  - [ ]* 1.3 Write property test for case-insensitive comparison
    - **Property 10: Case-Insensitive Comparison**
    - **Validates: Requirements 3.5**
  - [ ]* 1.4 Write property test for duplicate warnings (email, SSH, token)
    - **Property 7: Email Duplicate Warning**
    - **Property 8: SSH Key Duplicate Warning**
    - **Property 9: Token Username Duplicate Warning**
    - **Validates: Requirements 3.2, 3.3, 3.4**



- [x] 2. Enhance PlatformFormatter Component

  - [x] 2.1 Create `internal/account/platform.go` with PlatformInfo struct and GetPlatformInfo function

    - Define platform icons mapping (🐙 GitHub, 🦊 GitLab, 🪣 Bitbucket, 🍵 Gitea, 🔗 Other)
    - Implement GetPlatformInfo to return icon, name, and domain for each platform type
    - Implement DetectPlatformFromURL to identify platform from remote URL
    - _Requirements: 1.2, 4.4_
  - [ ]* 2.2 Write property test for platform icon mapping
    - **Property 2: Platform Icon Mapping Consistency**
    - **Validates: Requirements 1.2**
  - [ ]* 2.3 Write property test for platform detection from URL
    - **Property 12: Platform Detection from URL**


    - **Validates: Requirements 4.4**

- [x] 3. Enhance ActiveDetector with Scoring System

  - [x] 3.1 Add MatchScore struct and DetectActiveWithScore method to `internal/account/detect.go`
    - Implement scoring: user.name match (30pts), user.email match (30pts), SSH key (20pts), platform match (20pts)
    - Return best matching account with confidence score and matched fields
    - Handle edge cases: not in git repo, no remote, no matching account
    - _Requirements: 2.1, 2.2, 2.3, 2.4_
  - [ ]* 3.2 Write property test for active account detection accuracy
    - **Property 5: Active Account Detection Accuracy**
    - **Validates: Requirements 2.2**




- [x] 4. Checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.


- [x] 5. Implement Health Indicators
  - [x] 5.1 Create `internal/account/health.go` with HealthIndicators struct
    - Implement CheckSSHKeyHealth to verify key file exists and is valid
    - Implement CheckTokenHealth placeholder (token validation requires API call)
    - Implement IsStale to check if health data is older than 24 hours
    - Implement GetHealthIndicator to return correct symbol (✓, ✗, ?)
    - _Requirements: 5.1, 5.2, 5.3, 5.4_
  - [ ]* 5.2 Write property test for health indicator correctness
    - **Property 13: Health Indicator Correctness**

    - **Validates: Requirements 5.1, 5.2**
  - [ ]* 5.3 Write property test for stale health data detection
    - **Property 14: Stale Health Data Detection**
    - **Validates: Requirements 5.4**


- [x] 6. Implement Enhanced Table Renderer
  - [x] 6.1 Create `internal/ui/table.go` with AccountTableRow struct and RenderAccountTable function

    - Implement table rendering with columns: status, name, platform (with icon), git identity, auth methods, health
    - Implement active account highlighting with ✓ ACTIVE indicator
    - Handle empty account list with helpful message
    - _Requirements: 1.1, 1.3, 1.4, 2.1_
  - [ ]* 6.2 Write property test for table rendering completeness
    - **Property 1: Account Table Rendering Completeness**
    - **Validates: Requirements 1.1**
  - [ ]* 6.3 Write property test for dual authentication display
    - **Property 3: Dual Authentication Display**
    - **Validates: Requirements 1.3**
  - [ ]* 6.4 Write property test for active account indicator
    - **Property 4: Active Account Indicator Presence**
    - **Validates: Requirements 2.1**


- [x] 7. Checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.


- [x] 8. Implement Platform URL Builder Enhancement
  - [x] 8.1 Enhance `internal/git/remote.go` with platform-specific URL building


    - Implement BuildRemoteURL for each platform (GitHub, GitLab, Bitbucket, Gitea)
    - Support both SSH and HTTPS URL formats
    - Handle custom domains for self-hosted instances
    - _Requirements: 4.3, 4.5_
  - [ ]* 8.2 Write property test for platform URL format correctness
    - **Property 11: Platform URL Format Correctness**
    - **Validates: Requirements 4.3, 4.5**


- [x] 9. Implement Account Serialization Round-Trip
  - [x] 9.1 Ensure `internal/config/config.go` properly handles JSON serialization


    - Verify all fields are preserved during serialize/deserialize
    - Handle missing optional fields gracefully with defaults
    - Implement debug JSON output for troubleshooting
    - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5_
  - [ ]* 9.2 Write property test for account configuration round-trip
    - **Property 15: Account Configuration Round-Trip**
    - **Validates: Requirements 6.1, 6.2, 6.3, 6.5**
  - [ ]* 9.3 Write property test for graceful optional field handling
    - **Property 16: Graceful Optional Field Handling**

    - **Validates: Requirements 6.4**

- [x] 10. Integrate Components into CLI Commands

  - [x] 10.1 Update `cmd/ghex/commands/account.go` to use new components

    - Integrate DuplicateValidator into runAddAccount
    - Update runList to use RenderAccountTable with health indicators
    - Update runStatus to show enhanced active detection with score
    - Add duplicate warnings during account addition
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 2.1, 2.2, 2.3, 2.4, 3.1, 3.2, 3.3, 3.4_


- [x] 11. Final Checkpoint - Ensure all tests pass

  - Ensure all tests pass, ask the user if questions arise.
