package lightning

import (
	"net/http"
	"os"
	"strings"
)

// parsePattern splits a route pattern string into its individual parts.
func parsePattern(pattern string) []string {
	parts := strings.Split(pattern, "/")
	result := make([]string, 0)
	for _, part := range parts {
		if part != "" {
			result = append(result, part)
		}
	}

	return result
}

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

// defaultNotFound is the default handler function for 404 Not Found error
func defaultNotFound(ctx *Context) {
	ctx.Text(http.StatusNotFound, http.StatusText(http.StatusNotFound))
}

// defaultInternalServerError is the default handler function for 500 Internal Server Error
func defaultInternalServerError(ctx *Context) {
	ctx.Text(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}
