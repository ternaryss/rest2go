package rest2go

import "strings"

func antToRegex(pattern string) string {
	marker := "\x00DOUBLESTAR\x00"
	pattern = strings.ReplaceAll(pattern, "**", marker)
	var builder strings.Builder

	for _, char := range pattern {
		switch char {
		case '.', '+', '(', ')', '[', ']', '{', '}', '^', '$', '|', '\\':
			builder.WriteRune('\\')
			builder.WriteRune(char)

		default:
			builder.WriteRune(char)
		}
	}

	regex := builder.String()
	regex = strings.ReplaceAll(regex, marker, ".*")
	regex = strings.ReplaceAll(regex, "*", `[^/]*`)
	regex = strings.ReplaceAll(regex, "?", `[^/]`)

	if !strings.HasSuffix(regex, "/") {
		regex += "/?"
	}

	return regex
}
