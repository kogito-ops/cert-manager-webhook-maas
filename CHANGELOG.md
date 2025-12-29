# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.15] - 2025-12-29

### Fixed

- Updated golang.org/x/crypto to v0.45.0 (CVE-2025-47914, CVE-2025-58181)
- Updated golang.org/x/oauth2 to v0.27.0 (CVE-2025-22868)

### Added

- CI workflow for pull request builds and tests
- Integration test skip when TEST_ZONE_NAME not set

### Changed

- Go version updated to 1.24.0
- Docker build workflow uses go-version-file instead of hardcoded version
- Added OCI image description label to Dockerfile for GitHub Packages

## [1.0.14] - 2025-12-18

### Fixed
- Fixed DNS record name extraction for nested subdomains
- Record name is now computed relative to the configured zone instead of splitting on first dot
- Example: `_acme-challenge.lb-test.k8s.example.com` with zone `example.com` now yields name `_acme-challenge.lb-test.k8s`

### Changed
- Zone configuration is now required (no longer optional)
- Updated markdownlint configuration formatting

## [1.0.13] - 2025-08-01

### Changed
- Updated Helm chart version to 1.0.13
- Updated container image tag to v1.0.13

## [1.0.12] - 2025-07-31

### Fixed
- Fixed DNS record creation by using name and domain fields separately instead of FQDN
- Resolved MAAS/BIND "bad name (check-names)" error when creating TXT records
- Properly split FQDN into name and domain components for MAAS API compatibility

## [1.0.11] - 2025-07-31

### Fixed
- Fixed MAAS API DNS record creation by using FQDN field directly instead of splitting into name and domain
- Removed trailing dots from FQDNs as MAAS API does not expect them
- Removed obsolete recordName function that was causing empty name parameters

## [1.0.10] - 2025-07-31

### Fixed
- Fixed MAAS API parameter conflict when creating DNS records
- Changed DNS record creation to use name+domain parameters instead of both fqdn and name
- Added default GROUP_NAME environment variable to prevent panic when not set

## [1.0.9] - 2025-07-31

### Changed
- Updated Helm chart version to 1.0.9
- Updated container image tag to v1.0.9

## [1.0.0] - 2025-07-29

### Added
- Initial implementation of cert-manager webhook for Canonical MAAS DNS
- Support for DNS01 ACME challenges using MAAS DNS API
- Integration with gomaasclient v0.15.0 for MAAS API communication
- Automatic zone detection and TXT record management
- Helm chart for easy deployment
- GitHub Actions workflows for CI/CD:
  - Automated testing and linting
  - Multi-architecture Docker image builds
  - Automatic Helm chart publishing to GitHub Pages
- Container image publishing to GitHub Container Registry
- Comprehensive documentation and installation guides
- Support for configurable MAAS API versions
- Kubernetes RBAC and security configurations

### Features
- **DNS Challenge Solver**: Implements cert-manager DNS01 webhook interface
- **MAAS Integration**: Uses official Canonical gomaasclient library
- **Multi-Architecture**: Supports linux/amd64 and linux/arm64 platforms
- **Helm Chart**: Production-ready chart with configurable values
- **CI/CD Pipeline**: Automated build, test, and release workflows
- **Documentation**: Complete installation and configuration guides

### Configuration Options
- `secretName`: Kubernetes secret containing MAAS API credentials
- `apiUrl`: MAAS API endpoint URL
- `zoneName`: DNS zone name (optional, auto-detected if not specified)
- `apiVersion`: MAAS API version (default: "2.0")

### Installation Methods
- Helm chart via GitHub Pages repository
- Direct kubectl deployment
- Source code installation

[1.0.15]: https://github.com/kogito-ops/cert-manager-webhook-maas/releases/tag/v1.0.15
[1.0.14]: https://github.com/kogito-ops/cert-manager-webhook-maas/releases/tag/v1.0.14
[1.0.13]: https://github.com/kogito-ops/cert-manager-webhook-maas/releases/tag/v1.0.13
[1.0.12]: https://github.com/kogito-ops/cert-manager-webhook-maas/releases/tag/v1.0.12
[1.0.11]: https://github.com/kogito-ops/cert-manager-webhook-maas/releases/tag/v1.0.11
[1.0.10]: https://github.com/kogito-ops/cert-manager-webhook-maas/releases/tag/v1.0.10
[1.0.9]: https://github.com/kogito-ops/cert-manager-webhook-maas/releases/tag/v1.0.9
[1.0.0]: https://github.com/kogito-ops/cert-manager-webhook-maas/releases/tag/v1.0.0