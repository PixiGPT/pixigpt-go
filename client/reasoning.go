package client

import "regexp"

var thinkPattern = regexp.MustCompile(`(?i)<think(?:ing)?>(.*?)</think(?:ing)?>`)

// extractReasoning extracts chain of thought reasoning from content.
// Returns (mainContent, reasoningContent).
func extractReasoning(content string) (string, string) {
	matches := thinkPattern.FindStringSubmatch(content)
	if len(matches) < 2 {
		return content, ""
	}

	reasoning := matches[1]
	// Remove thinking tags from main content
	mainContent := thinkPattern.ReplaceAllString(content, "")

	return mainContent, reasoning
}
