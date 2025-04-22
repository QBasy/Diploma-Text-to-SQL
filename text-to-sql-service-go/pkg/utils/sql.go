package utils

import (
	"regexp"
	"strings"
)

func StripSQLMarkdown(input string) string {
	re := regexp.MustCompile("(?s)```sql\\s*(.*?)\\s*```")
	matches := re.FindStringSubmatch(input)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return strings.TrimSpace(input)
}
