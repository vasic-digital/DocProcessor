# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [1.0.0] - 2026-03-19

### Added
- `pkg/loader` -- Document loading for Markdown, YAML, HTML, AsciiDoc, RST
- `pkg/feature` -- Feature extraction, FeatureMap building, category/platform matrix
- `pkg/coverage` -- Thread-safe coverage tracking with RWMutex
- `pkg/docgraph` -- Inter-document link graph with JSON/Mermaid export
- `pkg/llm` -- LLMAgent interface and versioned prompt templates
- `pkg/config` -- Configuration from .env files
- `cmd/docprocessor` -- CLI entry point
- 190+ tests: unit, integration, stress, security, E2E, automation
- Makefile with build, test, coverage, vet, fmt targets
- 4-remote Upstreams scripts
- Full documentation: README, ARCHITECTURE, API_REFERENCE, USER_GUIDE
- CONTRIBUTING, CHANGELOG, CLAUDE, AGENTS
- VIDEO_COURSE with 5 episode scripts
- Apache 2.0 license
