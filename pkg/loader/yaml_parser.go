// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Milos Vasic

package loader

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// parseYAML extracts a title and sections from YAML content.
// It parses the YAML into a generic structure and creates sections
// from top-level keys.
func parseYAML(content string) (title string, sections []Section, links []string, err error) {
	var data map[string]interface{}
	if err := yaml.Unmarshal([]byte(content), &data); err != nil {
		return "", nil, nil, fmt.Errorf("yaml parse: %w", err)
	}

	// Title from "title", "name", or first key
	if t, ok := data["title"]; ok {
		title = fmt.Sprintf("%v", t)
	} else if t, ok := data["name"]; ok {
		title = fmt.Sprintf("%v", t)
	}

	// Each top-level key becomes a section
	line := 1
	for key, val := range data {
		valStr := formatYAMLValue(val)
		sections = append(sections, Section{
			Title:   key,
			Level:   1,
			Content: valStr,
			Line:    line,
		})
		line += strings.Count(valStr, "\n") + 2
	}

	return title, sections, nil, nil
}

// formatYAMLValue converts a YAML value to a readable string representation.
func formatYAMLValue(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case []interface{}:
		var parts []string
		for _, item := range val {
			parts = append(parts, fmt.Sprintf("- %v", item))
		}
		return strings.Join(parts, "\n")
	case map[string]interface{}:
		var parts []string
		for k, v := range val {
			parts = append(parts, fmt.Sprintf("%s: %v", k, v))
		}
		return strings.Join(parts, "\n")
	default:
		return fmt.Sprintf("%v", v)
	}
}
