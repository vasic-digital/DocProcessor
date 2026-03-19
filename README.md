# DocProcessor

Documentation processing and feature map extraction for QA automation.

## Overview

DocProcessor is a standalone Go module that loads project documentation, builds structured feature maps, and tracks verification coverage. It is designed to work with LLM agents for intelligent feature extraction, but also includes heuristic-based extraction for offline use.

## Quick Start

```bash
# Clone
git clone git@github.com:vasic-digital/DocProcessor.git
cd DocProcessor

# Build
go build ./...

# Test
go test ./... -race -count=1

# Run
go run ./cmd/docprocessor /path/to/docs
```

## Architecture

DocProcessor is organized into 6 packages:

| Package | Purpose |
|---------|---------|
| `pkg/loader` | Document loading and parsing (Markdown, YAML, HTML, AsciiDoc, RST) |
| `pkg/feature` | Feature extraction, FeatureMap building, FeatureMapBuilder |
| `pkg/coverage` | Thread-safe coverage tracking with RWMutex |
| `pkg/docgraph` | Inter-document link graph with JSON/Mermaid export |
| `pkg/llm` | LLMAgent interface and prompt templates |
| `pkg/config` | Configuration from .env files |

### Processing Pipeline

```
Load Docs -> Parse Sections -> Extract Features -> Build FeatureMap -> Enrich (LLM) -> Track Coverage
```

1. **Load & Parse** -- Scan project tree for documentation files
2. **Extract Features** -- Heuristic extraction or LLM-powered extraction
3. **Build Feature Map** -- Structured, queryable map with categories and platform matrix
4. **Enrich** -- LLM infers screens and generates test steps
5. **Track Coverage** -- Thread-safe per-platform verification tracking

## Key Interfaces

- `loader.Loader` -- Load documents from filesystem
- `feature.FeatureMapBuilder` -- Build feature maps from documents
- `coverage.CoverageTracker` -- Track feature verification status
- `llm.LLMAgent` -- Injected LLM for intelligent extraction (no hard dependency)

## Configuration

Copy `.env.example` to `.env` and customize:

```bash
HELIX_DOCS_ROOT=./docs
HELIX_DOCS_AUTO_DISCOVER=true
HELIX_DOCS_FORMATS=md,yaml,html,adoc,rst
```

## Testing

```bash
make test          # Run tests
make test-race     # Run with race detection
make test-cover    # Run with coverage report
```

190+ tests across 6 test types: unit, integration, stress, security, E2E, automation.

## License

Apache License 2.0. See [LICENSE](LICENSE).
