// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Milos Vasic

package loader

import (
	"regexp"
	"strings"
)

var (
	markdownHeadingRe = regexp.MustCompile(`(?m)^(#{1,6})\s+(.+)$`)
	markdownLinkRe    = regexp.MustCompile(`\[([^\]]*)\]\(([^)]+)\)`)
)

// parseMarkdown extracts sections and links from Markdown content.
func parseMarkdown(content string) (title string, sections []Section, links []string) {
	// Extract headings as sections
	matches := markdownHeadingRe.FindAllStringSubmatchIndex(content, -1)
	for i, loc := range matches {
		level := loc[3] - loc[2] // length of # prefix
		heading := content[loc[4]:loc[5]]

		// Determine section content: from end of heading line to next heading or EOF
		sectionStart := loc[1]
		var sectionEnd int
		if i+1 < len(matches) {
			sectionEnd = matches[i+1][0]
		} else {
			sectionEnd = len(content)
		}
		sectionContent := strings.TrimSpace(content[sectionStart:sectionEnd])

		// Line number (1-based)
		line := strings.Count(content[:loc[0]], "\n") + 1

		sections = append(sections, Section{
			Title:   heading,
			Level:   level,
			Content: sectionContent,
			Line:    line,
		})

		// First h1 is the document title
		if title == "" && level == 1 {
			title = heading
		}
	}

	// Extract links
	linkMatches := markdownLinkRe.FindAllStringSubmatch(content, -1)
	seen := make(map[string]bool)
	for _, m := range linkMatches {
		link := m[2]
		if !seen[link] {
			links = append(links, link)
			seen[link] = true
		}
	}

	return title, sections, links
}
