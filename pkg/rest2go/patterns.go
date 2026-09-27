package rest2go

import "strings"

func antToRegex(pattern string) string {
	var builder strings.Builder

	for index := 0; index < len(pattern); index++ {
		char := pattern[index]

		if char == '*' {
			if index+1 < len(pattern) && pattern[index+1] == '*' {
				builder.WriteString(".*")
				index++
				continue
			}

			builder.WriteString(`[^/]*`)
			continue
		}

		if char == '?' {
			builder.WriteString(`[^/]`)
			continue
		}

		switch char {
		case '.', '+', '(', ')', '[', ']', '{', '}', '^', '$', '|', '\\':
			builder.WriteByte('\\')
			builder.WriteByte(char)

		default:
			builder.WriteByte(char)
		}
	}

	regex := builder.String()
	regex = strings.ReplaceAll(regex, "/.*/", "(?:/.*)?/")

	if strings.HasSuffix(pattern, "/**") {
		regex = strings.TrimSuffix(regex, "/.*") + "(?:/.*)?"
	}

	if !strings.HasSuffix(regex, "/") {
		regex += "/?"
	}

	return regex
}
