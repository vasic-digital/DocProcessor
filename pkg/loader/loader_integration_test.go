// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Milos Vasic

package loader

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_LoadDir_MixedFormats(t *testing.T) {
	dir := t.TempDir()

	// Create a realistic project docs directory
	docsDir := filepath.Join(dir, "docs")
	err := os.MkdirAll(docsDir, 0755)
	require.NoError(t, err)

	// README.md
	err = os.WriteFile(filepath.Join(dir, "README.md"), []byte(`# Project Name

## Overview

This is a cross-platform text editor supporting multiple formats.

## Quick Start

Install and run the application with the instructions below.
`), 0644)
	require.NoError(t, err)

	// docs/architecture.md
	err = os.WriteFile(filepath.Join(docsDir, "architecture.md"), []byte(`# Architecture

## Module Structure

The project uses a shared module pattern.

## Build System

Gradle with Kotlin DSL.

See [README](../README.md) for details.
`), 0644)
	require.NoError(t, err)

	// docs/config.yaml
	err = os.WriteFile(filepath.Join(docsDir, "config.yaml"), []byte(`title: Project Config
version: 1.0.0
platforms:
  - android
  - desktop
  - web
`), 0644)
	require.NoError(t, err)

	// docs/page.html
	err = os.WriteFile(filepath.Join(docsDir, "page.html"), []byte(`<html>
<head><title>Documentation</title></head>
<body><h1>Docs</h1></body>
</html>`), 0644)
	require.NoError(t, err)

	l := NewDefaultLoader([]string{"md", "yaml", "html"})
	docs, err := l.LoadDir(context.Background(), dir)
	require.NoError(t, err)

	assert.Len(t, docs, 4)

	// Verify each format was parsed correctly
	var mdCount, yamlCount, htmlCount int
	for _, doc := range docs {
		switch doc.Format {
		case "md":
			mdCount++
			assert.NotEmpty(t, doc.Title)
			assert.NotEmpty(t, doc.Sections)
		case "yaml":
			yamlCount++
			assert.Equal(t, "Project Config", doc.Title)
		case "html":
			htmlCount++
		}
	}

	assert.Equal(t, 2, mdCount)
	assert.Equal(t, 1, yamlCount)
	assert.Equal(t, 1, htmlCount)
}

func TestIntegration_LoadFile_MarkdownWithLinks(t *testing.T) {
	dir := t.TempDir()

	// Create two linked markdown files
	file1 := filepath.Join(dir, "doc1.md")
	err := os.WriteFile(file1, []byte(`# Document 1

See [Document 2](./doc2.md) for more info.
Also check [external](https://example.com).
`), 0644)
	require.NoError(t, err)

	file2 := filepath.Join(dir, "doc2.md")
	err = os.WriteFile(file2, []byte(`# Document 2

Referenced from [Document 1](./doc1.md).
`), 0644)
	require.NoError(t, err)

	l := NewDefaultLoader([]string{"md"})

	doc1, err := l.LoadFile(context.Background(), file1)
	require.NoError(t, err)
	assert.Contains(t, doc1.Links, "./doc2.md")
	assert.Contains(t, doc1.Links, "https://example.com")

	doc2, err := l.LoadFile(context.Background(), file2)
	require.NoError(t, err)
	assert.Contains(t, doc2.Links, "./doc1.md")
}

func TestIntegration_LoadDir_WellKnownPatterns(t *testing.T) {
	dir := t.TempDir()

	// Create well-known doc files
	wellKnown := []string{
		"README.md", "CONTRIBUTING.md", "ARCHITECTURE.md",
		"USER_GUIDE.md", "API_REFERENCE.md", "CHANGELOG.md",
	}

	for _, name := range wellKnown {
		content := "# " + name + "\n\nContent for " + name + " with sufficient text to parse."
		err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644)
		require.NoError(t, err)
	}

	l := NewDefaultLoader([]string{"md"})
	docs, err := l.LoadDir(context.Background(), dir)
	require.NoError(t, err)

	assert.Len(t, docs, len(wellKnown))
}

func TestIntegration_LoadDir_ContextCancelMidway(t *testing.T) {
	dir := t.TempDir()

	// Create enough files to have a chance of mid-cancel
	for i := 0; i < 50; i++ {
		name := filepath.Join(dir, filepath.Join("sub", "doc"+string(rune('0'+i%10))+".md"))
		os.MkdirAll(filepath.Dir(name), 0755)
		content := "# Doc\n\nSome content here for the doc."
		os.WriteFile(name, []byte(content), 0644)
	}

	ctx, cancel := context.WithCancel(context.Background())
	// Cancel immediately
	cancel()

	l := NewDefaultLoader([]string{"md"})
	_, err := l.LoadDir(ctx, dir)
	assert.ErrorIs(t, err, context.Canceled)
}
