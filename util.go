package lightning

import (
	"os"
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

func resolveAddress(addr []string) string {
	if port := os.Getenv("PORT"); port != "" {
		logger.Debug("Environment variable PORT=\"%s\"", port)
		return ":" + port
	}

	switch len(addr) {
	case 0:
		logger.Debug("Using port :6789 by default")
		return ":6789"
	case 1:
		return addr[0]
	default:
		panic("too many parameters")
	}
}
