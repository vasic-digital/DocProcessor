// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Milos Vasic

// Package config provides configuration loading for DocProcessor.
package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Config holds DocProcessor configuration values.
type Config struct {
	DocsRoot     string   // Root directory for documentation files
	AutoDiscover bool     // Whether to auto-discover docs by well-known patterns
	Formats      []string // Supported file formats (e.g., md, yaml, html, adoc, rst)
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		DocsRoot:     "./docs",
		AutoDiscover: true,
		Formats:      []string{"md", "yaml", "html", "adoc", "rst"},
	}
}

// LoadFromEnv loads configuration from a .env file at the given path.
// Lines starting with # are treated as comments. Blank lines are skipped.
// Format: KEY=VALUE (no quoting required).
func LoadFromEnv(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %s: %w", path, err)
	}
	defer f.Close()

	env := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		env[key] = value
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("config: read %s: %w", path, err)
	}

	cfg := DefaultConfig()

	if v, ok := env["HELIX_DOCS_ROOT"]; ok {
		cfg.DocsRoot = v
	}
	if v, ok := env["HELIX_DOCS_AUTO_DISCOVER"]; ok {
		cfg.AutoDiscover = v == "true"
	}
	if v, ok := env["HELIX_DOCS_FORMATS"]; ok {
		formats := strings.Split(v, ",")
		cleaned := make([]string, 0, len(formats))
		for _, f := range formats {
			f = strings.TrimSpace(f)
			if f != "" {
				cleaned = append(cleaned, f)
			}
		}
		if len(cleaned) > 0 {
			cfg.Formats = cleaned
		}
	}

	return cfg, nil
}

// LoadFromMap creates a Config from a key-value map (useful for testing).
func LoadFromMap(env map[string]string) *Config {
	cfg := DefaultConfig()

	if v, ok := env["HELIX_DOCS_ROOT"]; ok {
		cfg.DocsRoot = v
	}
	if v, ok := env["HELIX_DOCS_AUTO_DISCOVER"]; ok {
		cfg.AutoDiscover = v == "true"
	}
	if v, ok := env["HELIX_DOCS_FORMATS"]; ok {
		formats := strings.Split(v, ",")
		cleaned := make([]string, 0, len(formats))
		for _, f := range formats {
			f = strings.TrimSpace(f)
			if f != "" {
				cleaned = append(cleaned, f)
			}
		}
		if len(cleaned) > 0 {
			cfg.Formats = cleaned
		}
	}

	return cfg
}
