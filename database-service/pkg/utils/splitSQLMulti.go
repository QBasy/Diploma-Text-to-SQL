package utils

import "strings"

func removeLineComments(query string) string {
	lines := strings.Split(query, "\n")
	var cleaned []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "--") {
			continue
		}
		if strings.Contains(line, "--") {
			parts := strings.Split(line, "--")
			line = parts[0]
		}
		cleaned = append(cleaned, line)
	}
	return strings.Join(cleaned, "\n")
}

func SplitSQLStatements(query string) []string {
	query = removeLineComments(query)
	statements := strings.Split(query, ";")
	var result []string
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt != "" {
			result = append(result, stmt)
		}
	}
	return result
}
