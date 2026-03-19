// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Milos Vasic

package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadFromEnv_PathTraversalInValues(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")

	// Values containing path traversal are stored as-is
	// (validation is the caller's responsibility)
	content := `HELIX_DOCS_ROOT=../../etc/passwd
`
	err := os.WriteFile(envFile, []byte(content), 0644)
	require.NoError(t, err)

	cfg, err := LoadFromEnv(envFile)
	require.NoError(t, err)

	// The config stores the raw value; callers should validate
	assert.Equal(t, "../../etc/passwd", cfg.DocsRoot)
}

func TestLoadFromEnv_LargeFile(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")

	// Create a large env file (100K)
	var sb strings.Builder
	for i := 0; i < 10000; i++ {
		sb.WriteString("# comment line\n")
	}
	sb.WriteString("HELIX_DOCS_ROOT=/valid\n")

	err := os.WriteFile(envFile, []byte(sb.String()), 0644)
	require.NoError(t, err)

	cfg, err := LoadFromEnv(envFile)
	require.NoError(t, err)
	assert.Equal(t, "/valid", cfg.DocsRoot)
}

func TestLoadFromEnv_SpecialCharactersInValues(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")

	content := `HELIX_DOCS_ROOT=/path/with spaces/and-special_chars.123
`
	err := os.WriteFile(envFile, []byte(content), 0644)
	require.NoError(t, err)

	cfg, err := LoadFromEnv(envFile)
	require.NoError(t, err)
	assert.Equal(t, "/path/with spaces/and-special_chars.123", cfg.DocsRoot)
}

func TestLoadFromEnv_EqualsInValue(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")

	content := `HELIX_DOCS_ROOT=/path=with=equals
`
	err := os.WriteFile(envFile, []byte(content), 0644)
	require.NoError(t, err)

	cfg, err := LoadFromEnv(envFile)
	require.NoError(t, err)
	assert.Equal(t, "/path=with=equals", cfg.DocsRoot)
}

func TestLoadFromEnv_PermissionDenied(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")

	err := os.WriteFile(envFile, []byte("HELIX_DOCS_ROOT=/test"), 0000)
	require.NoError(t, err)

	_, err = LoadFromEnv(envFile)
	assert.Error(t, err)
}
