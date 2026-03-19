// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Milos Vasic

package coverage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTracker(t *testing.T) {
	tracker := NewTracker()
	assert.NotNil(t, tracker)

	report := tracker.Coverage()
	assert.Equal(t, 0, report.Total)
}

func TestRegisterFeature(t *testing.T) {
	tracker := NewTracker()
	tracker.RegisterFeature(
		Feature{ID: "f1", Name: "Feature 1", Category: "ui"},
		[]string{"android", "desktop"},
	)

	report := tracker.Coverage()
	assert.Equal(t, 2, report.Total) // 1 feature x 2 platforms
	assert.Equal(t, 2, report.Unverified)
	assert.Equal(t, 0.0, report.OverallPct)
}

func TestMarkVerified(t *testing.T) {
	tracker := NewTracker()
	tracker.RegisterFeature(
		Feature{ID: "f1", Name: "Feature 1", Category: "ui"},
		[]string{"android", "desktop"},
	)

	tracker.MarkVerified("f1", "android", Evidence{
		ScreenshotPath: "/screenshots/f1-android.png",
		Timestamp:      time.Now(),
	})

	report := tracker.Coverage()
	assert.Equal(t, 1, report.Verified)
	assert.Equal(t, 1, report.Unverified)
	assert.Equal(t, 0.5, report.OverallPct)
}

func TestMarkFailed(t *testing.T) {
	tracker := NewTracker()
	tracker.RegisterFeature(
		Feature{ID: "f1", Name: "Feature 1", Category: "ui"},
		[]string{"android"},
	)

	tracker.MarkFailed("f1", "android", Issue{
		Type:     "functional",
		Severity: "high",
		Title:    "Button doesn't work",
	})

	report := tracker.Coverage()
	assert.Equal(t, 0, report.Verified)
	assert.Equal(t, 1, report.Failed)
	assert.Len(t, report.Issues, 1)
	assert.Equal(t, "Button doesn't work", report.Issues[0].Title)
}

func TestMarkSkipped(t *testing.T) {
	tracker := NewTracker()
	tracker.RegisterFeature(
		Feature{ID: "f1", Name: "Feature 1", Category: "ui"},
		[]string{"android"},
	)

	tracker.MarkSkipped("f1", "android", "Not applicable on this device")

	report := tracker.Coverage()
	assert.Equal(t, 1, report.Skipped)
	assert.Equal(t, 0, report.Unverified)
}

func TestMarkVerified_NonexistentFeature(t *testing.T) {
	tracker := NewTracker()
	// Should not panic
	tracker.MarkVerified("nonexistent", "android", Evidence{})
}

func TestMarkFailed_NonexistentFeature(t *testing.T) {
	tracker := NewTracker()
	// Should not panic
	tracker.MarkFailed("nonexistent", "android", Issue{})
}

func TestMarkSkipped_NonexistentFeature(t *testing.T) {
	tracker := NewTracker()
	// Should not panic
	tracker.MarkSkipped("nonexistent", "android", "reason")
}

func TestCoverageByPlatform(t *testing.T) {
	tracker := NewTracker()
	tracker.RegisterFeature(
		Feature{ID: "f1", Category: "ui"},
		[]string{"android", "desktop"},
	)
	tracker.RegisterFeature(
		Feature{ID: "f2", Category: "format"},
		[]string{"android", "desktop"},
	)

	tracker.MarkVerified("f1", "android", Evidence{})
	tracker.MarkVerified("f2", "android", Evidence{})
	tracker.MarkVerified("f1", "desktop", Evidence{})

	assert.Equal(t, 1.0, tracker.CoverageByPlatform("android"))
	assert.Equal(t, 0.5, tracker.CoverageByPlatform("desktop"))
	assert.Equal(t, 0.0, tracker.CoverageByPlatform("web"))
}

func TestCoverageByCategory(t *testing.T) {
	tracker := NewTracker()
	tracker.RegisterFeature(
		Feature{ID: "f1", Category: "ui"},
		[]string{"android"},
	)
	tracker.RegisterFeature(
		Feature{ID: "f2", Category: "ui"},
		[]string{"android"},
	)
	tracker.RegisterFeature(
		Feature{ID: "f3", Category: "format"},
		[]string{"android"},
	)

	tracker.MarkVerified("f1", "android", Evidence{})

	assert.Equal(t, 0.5, tracker.CoverageByCategory("ui"))
	assert.Equal(t, 0.0, tracker.CoverageByCategory("format"))
	assert.Equal(t, 0.0, tracker.CoverageByCategory("nonexistent"))
}

func TestUnverified(t *testing.T) {
	tracker := NewTracker()
	tracker.RegisterFeature(
		Feature{ID: "f1", Name: "Verified Feature", Category: "ui"},
		[]string{"android"},
	)
	tracker.RegisterFeature(
		Feature{ID: "f2", Name: "Unverified Feature", Category: "format"},
		[]string{"android"},
	)
	tracker.RegisterFeature(
		Feature{ID: "f3", Name: "Partially Verified", Category: "ui"},
		[]string{"android", "desktop"},
	)

	tracker.MarkVerified("f1", "android", Evidence{})
	tracker.MarkVerified("f3", "android", Evidence{})

	unverified := tracker.Unverified()
	assert.Len(t, unverified, 1)
	assert.Equal(t, "f2", unverified[0].ID)
}

func TestUnverified_AllVerified(t *testing.T) {
	tracker := NewTracker()
	tracker.RegisterFeature(
		Feature{ID: "f1", Category: "ui"},
		[]string{"android"},
	)

	tracker.MarkVerified("f1", "android", Evidence{})

	unverified := tracker.Unverified()
	assert.Len(t, unverified, 0)
}

func TestExport(t *testing.T) {
	tracker := NewTracker()
	tracker.RegisterFeature(
		Feature{ID: "f1", Name: "Feature 1", Category: "ui"},
		[]string{"android", "desktop"},
	)

	tracker.MarkVerified("f1", "android", Evidence{ScreenshotPath: "/test.png"})

	snapshot := tracker.Export()
	require.Contains(t, snapshot.Features, "f1")
	require.Contains(t, snapshot.Features["f1"], "android")
	assert.Equal(t, StateVerified, snapshot.Features["f1"]["android"].State)
	assert.Equal(t, StateUnverified, snapshot.Features["f1"]["desktop"].State)
	assert.Equal(t, 1, snapshot.Report.Verified)
}

func TestCoverage_MultipleFeatures(t *testing.T) {
	tracker := NewTracker()

	for i := 0; i < 10; i++ {
		tracker.RegisterFeature(
			Feature{ID: "f" + string(rune('0'+i)), Category: "ui"},
			[]string{"android", "desktop", "web"},
		)
	}

	report := tracker.Coverage()
	assert.Equal(t, 30, report.Total) // 10 features x 3 platforms
	assert.Equal(t, 30, report.Unverified)
	assert.Equal(t, 0.0, report.OverallPct)
}

func TestCoverage_AllStates(t *testing.T) {
	tracker := NewTracker()
	tracker.RegisterFeature(Feature{ID: "f1", Category: "ui"}, []string{"android"})
	tracker.RegisterFeature(Feature{ID: "f2", Category: "ui"}, []string{"android"})
	tracker.RegisterFeature(Feature{ID: "f3", Category: "ui"}, []string{"android"})
	tracker.RegisterFeature(Feature{ID: "f4", Category: "ui"}, []string{"android"})

	tracker.MarkVerified("f1", "android", Evidence{})
	tracker.MarkFailed("f2", "android", Issue{Title: "Bug"})
	tracker.MarkSkipped("f3", "android", "N/A")
	// f4 stays unverified

	report := tracker.Coverage()
	assert.Equal(t, 4, report.Total)
	assert.Equal(t, 1, report.Verified)
	assert.Equal(t, 1, report.Failed)
	assert.Equal(t, 1, report.Skipped)
	assert.Equal(t, 1, report.Unverified)
	assert.Equal(t, 0.25, report.OverallPct)
}

func TestMarkVerified_OverwritesPreviousState(t *testing.T) {
	tracker := NewTracker()
	tracker.RegisterFeature(Feature{ID: "f1", Category: "ui"}, []string{"android"})

	tracker.MarkFailed("f1", "android", Issue{Title: "Bug"})
	assert.Equal(t, 1, tracker.Coverage().Failed)

	tracker.MarkVerified("f1", "android", Evidence{})
	assert.Equal(t, 1, tracker.Coverage().Verified)
	assert.Equal(t, 0, tracker.Coverage().Failed)
}

func TestCoverageReport_ByPlatformAndCategory(t *testing.T) {
	tracker := NewTracker()
	tracker.RegisterFeature(Feature{ID: "f1", Category: "ui"}, []string{"android", "desktop"})
	tracker.RegisterFeature(Feature{ID: "f2", Category: "format"}, []string{"android"})

	tracker.MarkVerified("f1", "android", Evidence{})
	tracker.MarkVerified("f2", "android", Evidence{})

	report := tracker.Coverage()
	assert.Equal(t, 1.0, report.ByPlatform["android"])
	assert.Equal(t, 0.0, report.ByPlatform["desktop"])
	assert.Equal(t, 0.5, report.ByCategory["ui"])
	assert.Equal(t, 1.0, report.ByCategory["format"])
}
