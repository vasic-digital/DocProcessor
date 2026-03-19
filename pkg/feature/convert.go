// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Milos Vasic

package feature

import (
	"digital.vasic.docprocessor/pkg/llm"
)

// toLLMFeature converts a Feature to an llm.Feature.
func toLLMFeature(f Feature) llm.Feature {
	return llm.Feature{
		ID:          f.ID,
		Name:        f.Name,
		Description: f.Description,
		Category:    string(f.Category),
		Platforms:   f.Platforms,
	}
}

// toLLMFeatures converts a slice of Feature to llm.Feature.
func toLLMFeatures(features []Feature) []llm.Feature {
	result := make([]llm.Feature, len(features))
	for i, f := range features {
		result[i] = toLLMFeature(f)
	}
	return result
}

// fromLLMExpectedScreen converts an llm.ExpectedScreen to ExpectedScreen.
func fromLLMExpectedScreen(s llm.ExpectedScreen) ExpectedScreen {
	return ExpectedScreen{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		Features:    s.Features,
		Platforms:   s.Platforms,
	}
}

// fromLLMExpectedScreens converts a slice of llm.ExpectedScreen to ExpectedScreen.
func fromLLMExpectedScreens(screens []llm.ExpectedScreen) []ExpectedScreen {
	result := make([]ExpectedScreen, len(screens))
	for i, s := range screens {
		result[i] = fromLLMExpectedScreen(s)
	}
	return result
}

// fromLLMTestStep converts an llm.TestStep to TestStep.
func fromLLMTestStep(s llm.TestStep) TestStep {
	return TestStep{
		Order:       s.Order,
		Action:      s.Action,
		Expected:    s.Expected,
		ScreenID:    s.ScreenID,
		Platform:    s.Platform,
		Description: s.Description,
	}
}

// fromLLMTestSteps converts a slice of llm.TestStep to TestStep.
func fromLLMTestSteps(steps []llm.TestStep) []TestStep {
	result := make([]TestStep, len(steps))
	for i, s := range steps {
		result[i] = fromLLMTestStep(s)
	}
	return result
}
