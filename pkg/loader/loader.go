// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Milos Vasic

// Package loader provides document loading and parsing for project documentation.
package loader

import (
	"context"
	"time"
)

// Section represents a titled section within a document.
type Section struct {
	Title   string
	Level   int
	Content string
	Line    int
}

// Document represents a loaded and parsed documentation file.
type Document struct {
	Path       string
	Format     string
	Title      string
	Content    string
	Sections   []Section
	Links      []string
	ModifiedAt time.Time
}

// Loader defines the interface for loading documentation files.
type Loader interface {
	// LoadDir recursively loads all supported documents from a directory.
	LoadDir(ctx context.Context, path string) ([]Document, error)

	// LoadFile loads a single documentation file.
	LoadFile(ctx context.Context, path string) (Document, error)

	// SupportedFormats returns the list of file extensions this loader supports.
	SupportedFormats() []string
}
