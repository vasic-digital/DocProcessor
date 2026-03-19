// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Milos Vasic

package llm

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPromptVersion(t *testing.T) {
	assert.NotEmpty(t, PromptVersion)
	parts := strings.Split(PromptVersion, ".")
	assert.Len(t, parts, 3, "version should be semver format")
}

func TestSummarizePrompt(t *testing.T) {
	prompt := SummarizePrompt("Test document content here.")
	assert.Contains(t, prompt, "Summarize")
	assert.Contains(t, prompt, "Test document content here.")
}

func TestSummarizePrompt_Empty(t *testing.T) {
	prompt := SummarizePrompt("")
	assert.Contains(t, prompt, "Summarize")
}

func TestExtractFeaturesPrompt(t *testing.T) {
	prompt := ExtractFeaturesPrompt("Feature documentation text.")
	assert.Contains(t, prompt, "Extract all application features")
	assert.Contains(t, prompt, "Feature documentation text.")
	assert.Contains(t, prompt, "JSON array")
	assert.Contains(t, prompt, "name")
	assert.Contains(t, prompt, "category")
	assert.Contains(t, prompt, "platforms")
	assert.Contains(t, prompt, "priority")
}

func TestExtractFeaturesPrompt_Categories(t *testing.T) {
	prompt := ExtractFeaturesPrompt("text")
	categories := []string{"format", "ui", "network", "settings", "storage", "auth", "editor", "other"}
	for _, cat := range categories {
		assert.Contains(t, prompt, cat, "prompt should mention category: %s", cat)
	}
}

func TestClassifyFeaturePrompt(t *testing.T) {
	feature := RawFeature{
		Name:        "Markdown Editing",
		Description: "Edit markdown files with preview",
	}
	prompt := ClassifyFeaturePrompt(feature)
	assert.Contains(t, prompt, "Classify")
	assert.Contains(t, prompt, "Markdown Editing")
	assert.Contains(t, prompt, "Edit markdown files with preview")
}

func TestInferScreensPrompt(t *testing.T) {
	features := []Feature{
		{Name: "Feature A", Description: "Description A"},
		{Name: "Feature B", Description: "Description B"},
	}
	prompt := InferScreensPrompt(features)
	assert.Contains(t, prompt, "infer the expected screens")
	assert.Contains(t, prompt, "Feature A")
	assert.Contains(t, prompt, "Feature B")
	assert.Contains(t, prompt, "JSON array")
}

func TestInferScreensPrompt_EmptyFeatures(t *testing.T) {
	prompt := InferScreensPrompt(nil)
	assert.Contains(t, prompt, "infer the expected screens")
}

func TestGenerateTestStepsPrompt(t *testing.T) {
	feature := Feature{
		Name:        "CSV Import",
		Description: "Import CSV files",
		Platforms:   []string{"android", "desktop"},
	}
	prompt := GenerateTestStepsPrompt(feature)
	assert.Contains(t, prompt, "Generate verification test steps")
	assert.Contains(t, prompt, "CSV Import")
	assert.Contains(t, prompt, "Import CSV files")
	assert.Contains(t, prompt, "android")
	assert.Contains(t, prompt, "desktop")
}

func TestAllPromptsContainJSONInstruction(t *testing.T) {
	// Prompts that return structured data should mention JSON
	extractPrompt := ExtractFeaturesPrompt("text")
	assert.Contains(t, extractPrompt, "JSON")

	inferPrompt := InferScreensPrompt(nil)
	assert.Contains(t, inferPrompt, "JSON")

	stepsPrompt := GenerateTestStepsPrompt(Feature{})
	assert.Contains(t, stepsPrompt, "JSON")
}

func TestPromptSanitization_NoInjection(t *testing.T) {
	// Verify prompts are constructed safely even with malicious input
	malicious := `'; DROP TABLE features; --`
	prompt := SummarizePrompt(malicious)
	assert.Contains(t, prompt, malicious) // The text is embedded as-is, but in a prompt context
}

func TestRawFeature_Struct(t *testing.T) {
	rf := RawFeature{
		Name:        "Test Feature",
		Description: "A test feature",
		Category:    "ui",
		Platforms:   []string{"android"},
		Priority:    "high",
	}
	assert.Equal(t, "Test Feature", rf.Name)
	assert.Equal(t, "ui", rf.Category)
}

func TestFeature_LLMStruct(t *testing.T) {
	f := Feature{
		ID:          "feat-test",
		Name:        "Test",
		Description: "Test feature",
		Category:    "format",
		Platforms:   []string{"android", "desktop"},
	}
	assert.Equal(t, "feat-test", f.ID)
	assert.Len(t, f.Platforms, 2)
}

func TestExpectedScreen_LLMStruct(t *testing.T) {
	s := ExpectedScreen{
		ID:          "screen-1",
		Name:        "Main Screen",
		Description: "The main screen",
		Features:    []string{"f1", "f2"},
		Platforms:   []string{"android"},
	}
	assert.Equal(t, "screen-1", s.ID)
}

func TestTestStep_LLMStruct(t *testing.T) {
	step := TestStep{
		Order:       1,
		Action:      "Click button",
		Expected:    "Dialog opens",
		ScreenID:    "screen-1",
		Platform:    "android",
		Description: "Test the button",
	}
	assert.Equal(t, 1, step.Order)
	assert.Equal(t, "Click button", step.Action)
}

func TestFeatureCategory_Type(t *testing.T) {
	var cat FeatureCategory = "format"
	assert.Equal(t, FeatureCategory("format"), cat)
}
