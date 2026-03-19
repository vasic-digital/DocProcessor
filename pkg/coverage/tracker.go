// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Milos Vasic

// Package coverage provides thread-safe feature verification tracking
// with coverage reporting by platform and category.
package coverage

import (
	"sync"
	"time"
)

// Evidence represents proof that a feature was verified.
type Evidence struct {
	ScreenshotPath string
	VideoPath      string
	VideoOffset    time.Duration
	LogPath        string
	Timestamp      time.Time
	Description    string
}

// Issue represents a problem found during verification.
type Issue struct {
	Type        string // "visual", "ux", "accessibility", "functional", "performance", "crash"
	Severity    string // "critical", "high", "medium", "low"
	Title       string
	Description string
	ScreenID    string
	Evidence    []string // screenshot paths
}

// VerificationState represents the state of a feature verification.
type VerificationState string

const (
	StateUnverified VerificationState = "unverified"
	StateVerified   VerificationState = "verified"
	StateFailed     VerificationState = "failed"
	StateSkipped    VerificationState = "skipped"
)

// VerificationStatus holds the verification state for a feature on a specific platform.
type VerificationStatus struct {
	State    VerificationState
	Evidence Evidence
	Issue    Issue
	Reason   string // for skipped
}

// Feature is a minimal feature representation used by the tracker.
type Feature struct {
	ID       string
	Name     string
	Category string
}

// CoverageReport contains aggregate coverage statistics.
type CoverageReport struct {
	Total      int
	Verified   int
	Failed     int
	Skipped    int
	Unverified int
	OverallPct float64
	ByPlatform map[string]float64
	ByCategory map[string]float64
	Issues     []Issue
}

// CoverageSnapshot is a serializable snapshot of all coverage data.
type CoverageSnapshot struct {
	Features map[string]map[string]VerificationStatus `json:"features"` // featureID -> platform -> status
	Report   CoverageReport                           `json:"report"`
}

// CoverageTracker defines the interface for thread-safe coverage tracking.
type CoverageTracker interface {
	// RegisterFeature registers a feature for tracking.
	RegisterFeature(f Feature, platforms []string)

	// MarkVerified marks a feature as verified on a platform with evidence.
	MarkVerified(featureID string, platform string, evidence Evidence)

	// MarkFailed marks a feature as failed on a platform with an issue.
	MarkFailed(featureID string, platform string, issue Issue)

	// MarkSkipped marks a feature as skipped on a platform with a reason.
	MarkSkipped(featureID string, platform string, reason string)

	// Coverage returns the aggregate coverage report.
	Coverage() CoverageReport

	// CoverageByPlatform returns coverage percentage for a specific platform.
	CoverageByPlatform(platform string) float64

	// CoverageByCategory returns coverage percentage for a specific category.
	CoverageByCategory(category string) float64

	// Unverified returns all features that have not been verified on any platform.
	Unverified() []Feature

	// Export returns a serializable snapshot of all coverage data.
	Export() CoverageSnapshot
}

// featureStatus tracks the verification status of a feature across platforms.
type featureStatus struct {
	Feature   Feature
	Platforms map[string]VerificationStatus
}

// tracker is the thread-safe implementation of CoverageTracker.
type tracker struct {
	features map[string]*featureStatus
	mu       sync.RWMutex
}

// NewTracker creates a new CoverageTracker.
func NewTracker() CoverageTracker {
	return &tracker{
		features: make(map[string]*featureStatus),
	}
}

// RegisterFeature registers a feature for tracking on the given platforms.
func (t *tracker) RegisterFeature(f Feature, platforms []string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	status := &featureStatus{
		Feature:   f,
		Platforms: make(map[string]VerificationStatus),
	}
	for _, p := range platforms {
		status.Platforms[p] = VerificationStatus{State: StateUnverified}
	}
	t.features[f.ID] = status
}

// MarkVerified marks a feature as verified on a platform.
func (t *tracker) MarkVerified(featureID string, platform string, evidence Evidence) {
	t.mu.Lock()
	defer t.mu.Unlock()

	fs, ok := t.features[featureID]
	if !ok {
		return
	}
	fs.Platforms[platform] = VerificationStatus{
		State:    StateVerified,
		Evidence: evidence,
	}
}

// MarkFailed marks a feature as failed on a platform.
func (t *tracker) MarkFailed(featureID string, platform string, issue Issue) {
	t.mu.Lock()
	defer t.mu.Unlock()

	fs, ok := t.features[featureID]
	if !ok {
		return
	}
	fs.Platforms[platform] = VerificationStatus{
		State: StateFailed,
		Issue: issue,
	}
}

// MarkSkipped marks a feature as skipped on a platform.
func (t *tracker) MarkSkipped(featureID string, platform string, reason string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	fs, ok := t.features[featureID]
	if !ok {
		return
	}
	fs.Platforms[platform] = VerificationStatus{
		State:  StateSkipped,
		Reason: reason,
	}
}

// Coverage returns the aggregate coverage report.
func (t *tracker) Coverage() CoverageReport {
	t.mu.RLock()
	defer t.mu.RUnlock()

	report := CoverageReport{
		ByPlatform: make(map[string]float64),
		ByCategory: make(map[string]float64),
	}

	platformCounts := make(map[string][2]int) // [verified, total]
	categoryCounts := make(map[string][2]int) // [verified, total]

	for _, fs := range t.features {
		for platform, vs := range fs.Platforms {
			report.Total++

			pc := platformCounts[platform]
			pc[1]++

			cc := categoryCounts[fs.Feature.Category]
			cc[1]++

			switch vs.State {
			case StateVerified:
				report.Verified++
				pc[0]++
				cc[0]++
			case StateFailed:
				report.Failed++
				report.Issues = append(report.Issues, vs.Issue)
			case StateSkipped:
				report.Skipped++
			case StateUnverified:
				report.Unverified++
			}

			platformCounts[platform] = pc
			categoryCounts[fs.Feature.Category] = cc
		}
	}

	if report.Total > 0 {
		report.OverallPct = float64(report.Verified) / float64(report.Total)
	}

	for platform, counts := range platformCounts {
		if counts[1] > 0 {
			report.ByPlatform[platform] = float64(counts[0]) / float64(counts[1])
		}
	}

	for category, counts := range categoryCounts {
		if counts[1] > 0 {
			report.ByCategory[category] = float64(counts[0]) / float64(counts[1])
		}
	}

	return report
}

// CoverageByPlatform returns coverage percentage for a specific platform.
func (t *tracker) CoverageByPlatform(platform string) float64 {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var verified, total int
	for _, fs := range t.features {
		if vs, ok := fs.Platforms[platform]; ok {
			total++
			if vs.State == StateVerified {
				verified++
			}
		}
	}

	if total == 0 {
		return 0
	}
	return float64(verified) / float64(total)
}

// CoverageByCategory returns coverage percentage for a specific category.
func (t *tracker) CoverageByCategory(category string) float64 {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var verified, total int
	for _, fs := range t.features {
		if fs.Feature.Category != category {
			continue
		}
		for _, vs := range fs.Platforms {
			total++
			if vs.State == StateVerified {
				verified++
			}
		}
	}

	if total == 0 {
		return 0
	}
	return float64(verified) / float64(total)
}

// Unverified returns features that are not verified on any platform.
func (t *tracker) Unverified() []Feature {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var result []Feature
	for _, fs := range t.features {
		allUnverified := true
		for _, vs := range fs.Platforms {
			if vs.State == StateVerified {
				allUnverified = false
				break
			}
		}
		if allUnverified {
			result = append(result, fs.Feature)
		}
	}
	return result
}

// Export returns a serializable snapshot.
func (t *tracker) Export() CoverageSnapshot {
	// Compute report first (acquires its own RLock)
	report := t.Coverage()

	t.mu.RLock()
	defer t.mu.RUnlock()

	snapshot := CoverageSnapshot{
		Features: make(map[string]map[string]VerificationStatus),
		Report:   report,
	}

	for featureID, fs := range t.features {
		platforms := make(map[string]VerificationStatus)
		for p, vs := range fs.Platforms {
			platforms[p] = vs
		}
		snapshot.Features[featureID] = platforms
	}

	return snapshot
}
