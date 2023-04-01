package lightning

import (
	"strings"
)

// ParsePattern splits a route pattern string into its individual parts.
func ParsePattern(pattern string) []string {
	parts := strings.Split(pattern, "/")
	result := make([]string, 0)
	for _, part := range parts {
		if part != "" {
			result = append(result, part)
		}
	}

	return result
}

// Map is a shortcut for map[string]interface{}
type Map map[string]any
