package clients_hh

import (
	"html"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

// cleanHTML removes HTML tags and unescapes HTML entities from the input string.
func cleanHTML(input string) string {
	policy := bluemonday.StripTagsPolicy()
	cleaned := policy.Sanitize(input)
	return html.UnescapeString(strings.TrimSpace(cleaned))
}
