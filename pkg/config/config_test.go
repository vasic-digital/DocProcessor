// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Milos Vasic

package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	assert.Equal(t, "./docs", cfg.DocsRoot)
	assert.True(t, cfg.AutoDiscover)
	assert.Equal(t, []string{"md", "yaml", "html", "adoc", "rst"}, cfg.Formats)
}

func TestLoadFromEnv(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")

	content := `# DocProcessor config
HELIX_DOCS_ROOT=./documentation
HELIX_DOCS_AUTO_DISCOVER=true
HELIX_DOCS_FORMATS=md,yaml,html
`
	err := os.WriteFile(envFile, []byte(content), 0644)
	require.NoError(t, err)

	cfg, err := LoadFromEnv(envFile)
	require.NoError(t, err)

	assert.Equal(t, "./documentation", cfg.DocsRoot)
	assert.True(t, cfg.AutoDiscover)
	assert.Equal(t, []string{"md", "yaml", "html"}, cfg.Formats)
}

func TestLoadFromEnv_AutoDiscoverFalse(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")

	content := `HELIX_DOCS_AUTO_DISCOVER=false
`
	err := os.WriteFile(envFile, []byte(content), 0644)
	require.NoError(t, err)

	cfg, err := LoadFromEnv(envFile)
	require.NoError(t, err)

	assert.False(t, cfg.AutoDiscover)
}

func TestLoadFromEnv_Defaults(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")

	// Empty env file: all defaults
	err := os.WriteFile(envFile, []byte(""), 0644)
	require.NoError(t, err)

	cfg, err := LoadFromEnv(envFile)
	require.NoError(t, err)

	assert.Equal(t, "./docs", cfg.DocsRoot)
	assert.True(t, cfg.AutoDiscover)
	assert.Equal(t, []string{"md", "yaml", "html", "adoc", "rst"}, cfg.Formats)
}

func TestLoadFromEnv_Comments(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")

	content := `# This is a comment
HELIX_DOCS_ROOT=/opt/docs
# Another comment
`
	err := os.WriteFile(envFile, []byte(content), 0644)
	require.NoError(t, err)

	cfg, err := LoadFromEnv(envFile)
	require.NoError(t, err)

	assert.Equal(t, "/opt/docs", cfg.DocsRoot)
}

func TestLoadFromEnv_FileNotFound(t *testing.T) {
	_, err := LoadFromEnv("/nonexistent/path/.env")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "config: open")
}

func TestLoadFromEnv_MalformedLines(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")

	content := `HELIX_DOCS_ROOT=/opt/docs
NOEQUALSSIGN
=NOKEY
HELIX_DOCS_AUTO_DISCOVER=false
`
	err := os.WriteFile(envFile, []byte(content), 0644)
	require.NoError(t, err)

	cfg, err := LoadFromEnv(envFile)
	require.NoError(t, err)

	assert.Equal(t, "/opt/docs", cfg.DocsRoot)
	assert.False(t, cfg.AutoDiscover)
}

func TestLoadFromMap(t *testing.T) {
	env := map[string]string{
		"HELIX_DOCS_ROOT":          "/custom/path",
		"HELIX_DOCS_AUTO_DISCOVER": "true",
		"HELIX_DOCS_FORMATS":       "md,rst",
	}

	cfg := LoadFromMap(env)
	assert.Equal(t, "/custom/path", cfg.DocsRoot)
	assert.True(t, cfg.AutoDiscover)
	assert.Equal(t, []string{"md", "rst"}, cfg.Formats)
}

func TestLoadFromMap_Empty(t *testing.T) {
	cfg := LoadFromMap(map[string]string{})
	assert.Equal(t, DefaultConfig().DocsRoot, cfg.DocsRoot)
	assert.Equal(t, DefaultConfig().AutoDiscover, cfg.AutoDiscover)
	assert.Equal(t, DefaultConfig().Formats, cfg.Formats)
}

func TestLoadFromMap_FormatsWithSpaces(t *testing.T) {
	env := map[string]string{
		"HELIX_DOCS_FORMATS": "md, yaml, html",
	}

	cfg := LoadFromMap(env)
	assert.Equal(t, []string{"md", "yaml", "html"}, cfg.Formats)
}

func TestLoadFromMap_EmptyFormatsKeepsDefault(t *testing.T) {
	env := map[string]string{
		"HELIX_DOCS_FORMATS": "",
	}

	cfg := LoadFromMap(env)
	assert.Equal(t, DefaultConfig().Formats, cfg.Formats)
}

func TestLoadFromEnv_BlankLines(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")

	content := `

HELIX_DOCS_ROOT=/opt/docs

HELIX_DOCS_AUTO_DISCOVER=true

`
	err := os.WriteFile(envFile, []byte(content), 0644)
	require.NoError(t, err)

	cfg, err := LoadFromEnv(envFile)
	require.NoError(t, err)

	assert.Equal(t, "/opt/docs", cfg.DocsRoot)
	assert.True(t, cfg.AutoDiscover)
}
