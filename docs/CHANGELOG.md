# Changelog

All notable changes to this project will be documented in this file.

This project adheres to [Semantic Versioning](https://semver.org/).

[中文版](./CHANGELOG_CN.md)

## [v1.0.0] - 2024-03-26

### Added
- Initial release
- Basic AlloyDB connection functionality
- GORM integration
- Environment variable configuration support
- Connection pool configuration support (using Go standard library defaults)
- Version management functionality

### Features
- Using Go standard library default connection pool settings
- Support for customizing connection pool parameters via environment variables
- GetVersion() function for version checking
- Comprehensive documentation in both English and Chinese

### Technical Details
- Maximum open connections default: 0 (unlimited)
- Maximum idle connections default: 2
- Maximum connection lifetime default: 0 (unlimited)
- Maximum idle connection lifetime default: 0 (unlimited)
