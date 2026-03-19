// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Milos Vasic

// Package feature provides types and logic for building structured feature maps
// from project documentation.
package feature

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"digital.vasic.docprocessor/pkg/docgraph"
)

// FeatureCategory classifies features by their area.
type FeatureCategory string

const (
	CategoryFormat   FeatureCategory = "format"
	CategoryUI       FeatureCategory = "ui"
	CategoryNetwork  FeatureCategory = "network"
	CategorySettings FeatureCategory = "settings"
	CategoryStorage  FeatureCategory = "storage"
	CategoryAuth     FeatureCategory = "auth"
	CategoryEditor   FeatureCategory = "editor"
	CategoryOther    FeatureCategory = "other"
)

// AllCategories returns all valid feature categories.
func AllCategories() []FeatureCategory {
	return []FeatureCategory{
		CategoryFormat, CategoryUI, CategoryNetwork, CategorySettings,
		CategoryStorage, CategoryAuth, CategoryEditor, CategoryOther,
	}
}

// ValidCategory checks if a string is a valid FeatureCategory.
func ValidCategory(s string) bool {
	for _, c := range AllCategories() {
		if string(c) == s {
			return true
		}
	}
	return false
}

// TestStep describes a single test step for verifying a feature.
type TestStep struct {
	Order       int
	Action      string
	Expected    string
	ScreenID    string
	Platform    string
	Description string
}

// RawFeature is an unprocessed feature extracted from documentation.
type RawFeature struct {
	Name        string
	Description string
	Category    string
	Platforms   []string
	Priority    string
	SourceDoc   string
	Section     string
}

// Feature represents a fully processed, verified feature entry.
type Feature struct {
	ID            string
	Name          string
	Description   string
	Category      FeatureCategory
	Platforms     []string
	Priority      string
	Screens       []string
	TestSteps     []TestStep
	SourceDoc     string
	SourceSection string
}

// ExpectedScreen represents an expected application screen.
type ExpectedScreen struct {
	ID          string
	Name        string
	Description string
	Features    []string // Feature IDs visible on this screen
	Platforms   []string
}

// Workflow represents a multi-step user workflow.
type Workflow struct {
	ID          string
	Name        string
	Description string
	Steps       []WorkflowStep
	Features    []string // Feature IDs involved
	Platforms   []string
}

// WorkflowStep is a single step in a workflow.
type WorkflowStep struct {
	Order       int
	Description string
	ScreenID    string
	FeatureID   string
}

// FeatureMap is the central structured output of the documentation processor.
type FeatureMap struct {
	Features       []Feature
	Screens        []ExpectedScreen
	Workflows      []Workflow
	Categories     map[FeatureCategory][]Feature
	PlatformMatrix map[string][]Feature
	DocGraph       *docgraph.DocGraph
	GeneratedAt    time.Time
	ProjectRoot    string
}

// GenerateID creates a deterministic feature ID from a name.
func GenerateID(name string) string {
	normalized := strings.ToLower(strings.TrimSpace(name))
	normalized = strings.ReplaceAll(normalized, " ", "-")
	// Remove non-alphanumeric characters except hyphens
	var result strings.Builder
	for _, r := range normalized {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	slug := result.String()
	// Truncate long slugs but append hash suffix for uniqueness
	if len(slug) > 40 {
		h := sha256.Sum256([]byte(normalized))
		slug = slug[:32] + fmt.Sprintf("-%x", h[:4])
	}
	return "feat-" + slug
}

// NewFeatureMap creates an initialized FeatureMap.
func NewFeatureMap(projectRoot string) *FeatureMap {
	return &FeatureMap{
		Categories:     make(map[FeatureCategory][]Feature),
		PlatformMatrix: make(map[string][]Feature),
		DocGraph:       docgraph.New(),
		GeneratedAt:    time.Now(),
		ProjectRoot:    projectRoot,
	}
}

// AddFeature adds a feature to the map and updates categories and platform matrix.
func (fm *FeatureMap) AddFeature(f Feature) {
	fm.Features = append(fm.Features, f)
	fm.Categories[f.Category] = append(fm.Categories[f.Category], f)
	for _, p := range f.Platforms {
		fm.PlatformMatrix[p] = append(fm.PlatformMatrix[p], f)
	}
}

// FeatureByID returns a feature by its ID, or nil if not found.
func (fm *FeatureMap) FeatureByID(id string) *Feature {
	for i := range fm.Features {
		if fm.Features[i].ID == id {
			return &fm.Features[i]
		}
	}
	return nil
}

// FeaturesForPlatform returns all features available on a given platform.
func (fm *FeatureMap) FeaturesForPlatform(platform string) []Feature {
	return fm.PlatformMatrix[platform]
}

// FeaturesForCategory returns all features in a given category.
func (fm *FeatureMap) FeaturesForCategory(category FeatureCategory) []Feature {
	return fm.Categories[category]
}
