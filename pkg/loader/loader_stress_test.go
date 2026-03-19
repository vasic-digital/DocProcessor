// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Milos Vasic

package loader

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadDir_ManyFiles(t *testing.T) {
	dir := t.TempDir()

	// Create 100 markdown files
	for i := 0; i < 100; i++ {
		name := filepath.Join(dir, fmt.Sprintf("doc-%03d.md", i))
		content := fmt.Sprintf("# Document %d\n\nContent for document %d with enough text to extract.", i, i)
		err := os.WriteFile(name, []byte(content), 0644)
		require.NoError(t, err)
	}

	l := NewDefaultLoader([]string{"md"})
	docs, err := l.LoadDir(context.Background(), dir)
	require.NoError(t, err)

	assert.Len(t, docs, 100)
}

func TestLoadDir_ConcurrentReads(t *testing.T) {
	dir := t.TempDir()

	// Create some test files
	for i := 0; i < 20; i++ {
		name := filepath.Join(dir, fmt.Sprintf("doc-%02d.md", i))
		content := fmt.Sprintf("# Document %d\n\nContent for document number %d here.", i, i)
		err := os.WriteFile(name, []byte(content), 0644)
		require.NoError(t, err)
	}

	l := NewDefaultLoader([]string{"md"})

	// Concurrent LoadDir calls should be safe
	var wg sync.WaitGroup
	errors := make([]error, 10)
	results := make([][]Document, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			docs, err := l.LoadDir(context.Background(), dir)
			errors[idx] = err
			results[idx] = docs
		}(i)
	}

	wg.Wait()

	for i := 0; i < 10; i++ {
		assert.NoError(t, errors[i])
		assert.Len(t, results[i], 20)
	}
}

func TestLoadFile_ConcurrentLoadSameFile(t *testing.T) {
	dir := t.TempDir()
	mdFile := filepath.Join(dir, "shared.md")
	err := os.WriteFile(mdFile, []byte("# Shared\n\nShared content is here now."), 0644)
	require.NoError(t, err)

	l := NewDefaultLoader([]string{"md"})

	var wg sync.WaitGroup
	errors := make([]error, 50)
	docs := make([]Document, 50)

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			doc, err := l.LoadFile(context.Background(), mdFile)
			errors[idx] = err
			docs[idx] = doc
		}(i)
	}

	wg.Wait()

	for i := 0; i < 50; i++ {
		assert.NoError(t, errors[i])
		assert.Equal(t, "Shared", docs[i].Title)
	}
}

func TestLoadDir_DeepNesting(t *testing.T) {
	dir := t.TempDir()

	// Create 10 levels of nesting with one file each
	current := dir
	for i := 0; i < 10; i++ {
		current = filepath.Join(current, fmt.Sprintf("level%d", i))
		err := os.MkdirAll(current, 0755)
		require.NoError(t, err)

		name := filepath.Join(current, fmt.Sprintf("doc-level-%d.md", i))
		content := fmt.Sprintf("# Level %d\n\nContent at level %d with some extra text.", i, i)
		err = os.WriteFile(name, []byte(content), 0644)
		require.NoError(t, err)
	}

	l := NewDefaultLoader([]string{"md"})
	docs, err := l.LoadDir(context.Background(), dir)
	require.NoError(t, err)

	assert.Len(t, docs, 10)
}

func TestLoadFile_LargeMarkdownContent(t *testing.T) {
	dir := t.TempDir()
	mdFile := filepath.Join(dir, "large.md")

	// Create a markdown file with many sections
	var content string
	content += "# Large Document\n\n"
	for i := 0; i < 100; i++ {
		content += fmt.Sprintf("## Section %d\n\nParagraph for section %d with enough content to parse.\n\n", i, i)
	}

	err := os.WriteFile(mdFile, []byte(content), 0644)
	require.NoError(t, err)

	l := NewDefaultLoader([]string{"md"})
	doc, err := l.LoadFile(context.Background(), mdFile)
	require.NoError(t, err)

	assert.Equal(t, "Large Document", doc.Title)
	assert.Len(t, doc.Sections, 101) // 1 h1 + 100 h2
}
