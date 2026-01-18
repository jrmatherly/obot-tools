# Obot-Tools Patterns Changelog

All notable changes to implementation patterns will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned

- Credential store implementation pattern
- Tool testing and validation patterns
- Multi-language tool patterns (Go, Python, TypeScript comparison)

## [1.0.0] - 2026-01-16

### Added

- **model-provider-creation.md** - OpenAI-compatible model provider implementation pattern
- **auth-provider-creation.md** - OAuth2 authentication provider creation pattern
- **gptscript-tool-development.md** - Complete GPTScript .gpt file format and tool development guide
- Pattern README with version tracking and selection guide
- Pattern template for consistent structure
- Canonical examples reference section

### Documentation

- Created patterns/ directory structure
- Established semantic versioning for patterns
- Added CHANGELOG.md for tracking pattern evolution
- Integrated with existing Serena memory (gptscript_tool_format.md)

---

## Version Guidelines

### Major Version (X.0.0)

- Breaking changes to .gpt file format interpretation
- Incompatible API or server implementation changes
- GPTScript version updates requiring pattern rewrites

### Minor Version (X.Y.0)

- New pattern features or implementation options
- Additional troubleshooting guidance or examples
- Enhanced code snippets or canonical references
- Backwards-compatible improvements to procedures

### Patch Version (X.Y.Z)

- Bug fixes in code examples or procedures
- Typo corrections and documentation clarifications
- Updated component versions or dependency info
- Minor improvements to testing or validation steps

---

## Deprecation Policy

Patterns marked as deprecated will:

1. Remain available for 2 minor versions
2. Include migration guide to replacement pattern
3. Display deprecation notice at top of document
4. Be moved to `patterns/archive/` after removal

---

## Pattern Dependencies

**GPTScript Version:** Patterns compatible with GPTScript as of January 2026
**Go Version:** Go 1.23+ for model/auth providers
**Python Version:** Python 3.13+ for Python-based tools
**Node.js Version:** Node.js 18+ for TypeScript tools

---

**Maintainer:** Obot-Tools Development Team
**Last Updated:** 2026-01-16
