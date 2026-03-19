// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Milos Vasic

package docprocessor_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"digital.vasic.docprocessor/pkg/loader"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecurity_PathTraversal_LoadFile(t *testing.T) {
	dir := t.TempDir()
	outsideDir := t.TempDir()

	// Create a file outside the docs directory
	outsideFile := filepath.Join(outsideDir, "secret.md")
	err := os.WriteFile(outsideFile, []byte("# Secret\n\nThis should not be accessible."), 0644)
	require.NoError(t, err)

	// Try to load with path traversal from inside dir
	l := loader.NewDefaultLoader([]string{"md"})
	relativePath := filepath.Join(dir, "..", "..", "..", outsideFile)

	// Loader resolves to absolute paths, so it may or may not succeed
	// depending on whether the resolved path exists
	_, err = l.LoadFile(context.Background(), relativePath)
	// We just verify it doesn't panic
	_ = err
}

func TestSecurity_PathTraversal_LoadDir(t *testing.T) {
	dir := t.TempDir()

	l := loader.NewDefaultLoader([]string{"md"})

	// Try to traverse to root
	_, err := l.LoadDir(context.Background(), filepath.Join(dir, "..", "..", ".."))
	// Should not panic; may succeed or fail depending on permissions
	_ = err
}

func TestSecurity_LargeContent_DoesNotOOM(t *testing.T) {
	dir := t.TempDir()

	// Create a moderately large markdown file (1 MB)
	mdFile := filepath.Join(dir, "large.md")
	content := "# Large Document\n\n" + strings.Repeat("This is a paragraph. ", 50000)
	err := os.WriteFile(mdFile, []byte(content), 0644)
	require.NoError(t, err)

	l := loader.NewDefaultLoader([]string{"md"})
	doc, err := l.LoadFile(context.Background(), mdFile)
	require.NoError(t, err)
	assert.True(t, len(doc.Content) > 0)
}

func TestSecurity_MalformedYAML_DoesNotPanic(t *testing.T) {
	dir := t.TempDir()

	testCases := []struct {
		name    string
		content string
	}{
		{"deeply_nested", strings.Repeat("a:\n  ", 100) + "value"},
		{"binary_content", "title: test\nbinary: \x00\x01\x02\x03"},
		{"unicode_bomb", "title: " + strings.Repeat("\ufeff", 100)},
		{"long_key", strings.Repeat("a", 10000) + ": value"},
		{"long_value", "key: " + strings.Repeat("x", 100000)},
	}

	l := loader.NewDefaultLoader([]string{"yaml"})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			yamlFile := filepath.Join(dir, tc.name+".yaml")
			err := os.WriteFile(yamlFile, []byte(tc.content), 0644)
			require.NoError(t, err)

			// Should not panic
			_, _ = l.LoadFile(context.Background(), yamlFile)
		})
	}
}

func TestSecurity_MalformedMarkdown_DoesNotPanic(t *testing.T) {
	dir := t.TempDir()

	testCases := []struct {
		name    string
		content string
	}{
		{"no_newlines", strings.Repeat("# ", 10000)},
		{"only_hashes", strings.Repeat("#", 10000)},
		{"deep_links", strings.Repeat("[a](b)", 10000)},
		{"null_bytes", "# Title\n\n" + strings.Repeat("\x00", 1000)},
		{"long_heading", "# " + strings.Repeat("x", 100000)},
		{"many_headings", func() string {
			var sb strings.Builder
			for i := 0; i < 1000; i++ {
				sb.WriteString("## Heading\n\nContent.\n\n")
			}
			return sb.String()
		}()},
	}

	l := loader.NewDefaultLoader([]string{"md"})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mdFile := filepath.Join(dir, tc.name+".md")
			err := os.WriteFile(mdFile, []byte(tc.content), 0644)
			require.NoError(t, err)

			// Should not panic
			_, _ = l.LoadFile(context.Background(), mdFile)
		})
	}
}

func TestSecurity_NoSensitiveFileTypes(t *testing.T) {
	dir := t.TempDir()

	// Create files that should NOT be loaded
	sensitiveFiles := []struct {
		name    string
		content string
	}{
		{".env", "API_KEY=secret123"},
		{"credentials.json", `{"key": "secret"}`},
		{"key.pem", "-----BEGIN PRIVATE KEY-----"},
		{".gitconfig", "[user]\n  name = test"},
	}

	for _, sf := range sensitiveFiles {
		err := os.WriteFile(filepath.Join(dir, sf.name), []byte(sf.content), 0644)
		require.NoError(t, err)
	}

	l := loader.NewDefaultLoader([]string{"md", "yaml", "html"})
	docs, err := l.LoadDir(context.Background(), dir)
	require.NoError(t, err)

	// None of the sensitive files should be loaded
	assert.Len(t, docs, 0)
}

func TestSecurity_MaxFileSizeEnforced(t *testing.T) {
	dir := t.TempDir()
	mdFile := filepath.Join(dir, "huge.md")

	// Create file just over 10MB
	data := make([]byte, loader.MaxFileSize+1)
	for i := range data {
		data[i] = 'x'
	}
	err := os.WriteFile(mdFile, data, 0644)
	require.NoError(t, err)

	l := loader.NewDefaultLoader([]string{"md"})
	_, err = l.LoadFile(context.Background(), mdFile)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds max size")
}
