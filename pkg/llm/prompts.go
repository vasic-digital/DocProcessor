// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Milos Vasic

package llm

import "fmt"

// PromptVersion tracks the version of prompt templates for reproducibility.
const PromptVersion = "1.0.0"

// SummarizePrompt generates a prompt for document summarization.
func SummarizePrompt(text string) string {
	return fmt.Sprintf(`Summarize the following documentation in 2-3 concise paragraphs.
Focus on the key features, capabilities, and technical details.

---
%s
---

Provide a clear, factual summary.`, text)
}

// ExtractFeaturesPrompt generates a prompt for feature extraction.
func ExtractFeaturesPrompt(text string) string {
	return fmt.Sprintf(`Extract all application features from the following documentation.
Return a JSON array of objects with these fields:
- name: feature name (string)
- description: brief description (string)
- category: one of "format", "ui", "network", "settings", "storage", "auth", "editor", "other" (string)
- platforms: array of platforms, any of "android", "desktop", "web", "ios" (string array)
- priority: one of "critical", "high", "medium", "low" (string)

---
%s
---

Return ONLY the JSON array, no other text.`, text)
}

// ClassifyFeaturePrompt generates a prompt for feature classification.
func ClassifyFeaturePrompt(feature RawFeature) string {
	return fmt.Sprintf(`Classify the following feature into exactly one category.
Categories: format, ui, network, settings, storage, auth, editor, other

Feature name: %s
Feature description: %s

Return ONLY the category name, nothing else.`, feature.Name, feature.Description)
}

// InferScreensPrompt generates a prompt for screen inference.
func InferScreensPrompt(features []Feature) string {
	var featureList string
	for _, f := range features {
		featureList += fmt.Sprintf("- %s: %s\n", f.Name, f.Description)
	}
	return fmt.Sprintf(`Given these application features, infer the expected screens/views.
Return a JSON array of objects with these fields:
- id: screen identifier (string, e.g., "screen-settings")
- name: human-readable name (string)
- description: what the screen shows (string)
- features: array of feature IDs visible on this screen (string array)
- platforms: array of platforms where this screen exists (string array)

Features:
%s

Return ONLY the JSON array, no other text.`, featureList)
}

// GenerateTestStepsPrompt generates a prompt for test step generation.
func GenerateTestStepsPrompt(feature Feature) string {
	return fmt.Sprintf(`Generate verification test steps for this feature.
Return a JSON array of objects with these fields:
- order: step number starting at 1 (int)
- action: what to do (string)
- expected: expected result (string)
- screen_id: screen where this happens (string)
- platform: which platform, or "all" (string)
- description: brief description of the step (string)

Feature: %s
Description: %s
Platforms: %v

Return ONLY the JSON array, no other text.`, feature.Name, feature.Description, feature.Platforms)
}
