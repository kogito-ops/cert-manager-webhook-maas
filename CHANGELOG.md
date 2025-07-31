# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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

[1.0.10]: https://github.com/kogito-ops/cert-manager-webhook-maas/releases/tag/v1.0.10
[1.0.9]: https://github.com/kogito-ops/cert-manager-webhook-maas/releases/tag/v1.0.9
[1.0.0]: https://github.com/kogito-ops/cert-manager-webhook-maas/releases/tag/v1.0.0