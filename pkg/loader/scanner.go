// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Milos Vasic

package loader

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// MaxFileSize is the maximum file size the loader will process (10 MB).
const MaxFileSize = 10 * 1024 * 1024

// DefaultLoader implements the Loader interface for local filesystem documents.
type DefaultLoader struct {
	formats []string
}

// NewDefaultLoader creates a Loader that supports the given file extensions.
// Extensions should be provided without the leading dot (e.g., "md", "yaml").
func NewDefaultLoader(formats []string) *DefaultLoader {
	return &DefaultLoader{formats: formats}
}

// SupportedFormats returns the file extensions this loader handles.
func (l *DefaultLoader) SupportedFormats() []string {
	return l.formats
}

// LoadFile loads and parses a single documentation file.
func (l *DefaultLoader) LoadFile(ctx context.Context, path string) (Document, error) {
	select {
	case <-ctx.Done():
		return Document{}, ctx.Err()
	default:
	}

	// Security: resolve and validate path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return Document{}, fmt.Errorf("loader: resolve path %s: %w", path, err)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return Document{}, fmt.Errorf("loader: stat %s: %w", absPath, err)
	}
	if info.IsDir() {
		return Document{}, fmt.Errorf("loader: %s is a directory, not a file", absPath)
	}
	if info.Size() > MaxFileSize {
		return Document{}, fmt.Errorf("loader: file %s exceeds max size (%d > %d bytes)", absPath, info.Size(), MaxFileSize)
	}

	ext := strings.TrimPrefix(filepath.Ext(absPath), ".")
	if ext == "yml" {
		ext = "yaml"
	}
	if !l.isSupported(ext) {
		return Document{}, fmt.Errorf("loader: unsupported format %q for %s", ext, absPath)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return Document{}, fmt.Errorf("loader: read %s: %w", absPath, err)
	}

	content := string(data)

	doc := Document{
		Path:       absPath,
		Format:     ext,
		Content:    content,
		ModifiedAt: info.ModTime(),
	}

	switch ext {
	case "md":
		doc.Title, doc.Sections, doc.Links = parseMarkdown(content)
	case "yaml":
		title, sections, links, parseErr := parseYAML(content)
		if parseErr != nil {
			return Document{}, fmt.Errorf("loader: parse %s: %w", absPath, parseErr)
		}
		doc.Title = title
		doc.Sections = sections
		doc.Links = links
	default:
		// For html, adoc, rst: store raw content, extract title from first line
		lines := strings.SplitN(content, "\n", 2)
		if len(lines) > 0 {
			doc.Title = strings.TrimSpace(lines[0])
		}
		doc.Sections = []Section{
			{Title: doc.Title, Level: 1, Content: content, Line: 1},
		}
	}

	if doc.Title == "" {
		// Fallback: use filename without extension
		doc.Title = strings.TrimSuffix(filepath.Base(absPath), filepath.Ext(absPath))
	}

	return doc, nil
}

// LoadDir recursively scans a directory for supported documentation files.
func (l *DefaultLoader) LoadDir(ctx context.Context, path string) ([]Document, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("loader: resolve path %s: %w", path, err)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("loader: stat %s: %w", absPath, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("loader: %s is not a directory", absPath)
	}

	var docs []Document
	err = filepath.Walk(absPath, func(p string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if info.IsDir() {
			// Skip hidden directories
			if strings.HasPrefix(filepath.Base(p), ".") && p != absPath {
				return filepath.SkipDir
			}
			return nil
		}

		ext := strings.TrimPrefix(filepath.Ext(p), ".")
		if ext == "yml" {
			ext = "yaml"
		}
		if !l.isSupported(ext) {
			return nil
		}

		doc, err := l.LoadFile(ctx, p)
		if err != nil {
			// Skip files that fail to parse but log/continue
			return nil
		}
		docs = append(docs, doc)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("loader: walk %s: %w", absPath, err)
	}

	return docs, nil
}

func (l *DefaultLoader) isSupported(ext string) bool {
	for _, f := range l.formats {
		if strings.EqualFold(f, ext) {
			return true
		}
	}
	return false
}
