package lightning

import (
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestParsePattern(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		want    []string
	}{
		{
			name:    "empty pattern",
			pattern: "",
			want:    []string{},
		},
		{
			name:    "root pattern",
			pattern: "/",
			want:    []string{},
		},
		{
			name:    "double root pattern",
			pattern: "//",
			want:    []string{},
		},
		{
			name:    "single level pattern",
			pattern: "/foo",
			want:    []string{"foo"},
		},
		{
			name:    "single level pattern with double root",
			pattern: "//foo",
			want:    []string{"foo"},
		},
		{
			name:    "single level pattern with trailing slash",
			pattern: "/foo/",
			want:    []string{"foo"},
		},
		{
			name:    "multi-level pattern",
			pattern: "/foo/bar/baz",
			want:    []string{"foo", "bar", "baz"},
		},
		{
			name:    "pattern with named parameter",
			pattern: "/api/user/:userId",
			want:    []string{"api", "user", ":userId"},
		},
		{
			name:    "pattern with wildcard",
			pattern: "/api/*",
			want:    []string{"api", "*"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parsePattern(tt.pattern); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parsePattern() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResolveAddress(t *testing.T) {
	// Test case for when PORT environment variable is set
	os.Setenv("PORT", "1234")
	defer os.Unsetenv("PORT")
	expected := ":1234"
	result := resolveAddress([]string{})
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
	os.Unsetenv("PORT")

	// Test case for when no parameters are passed in
	expected = ":6789"
	result = resolveAddress([]string{})
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}

	// Test case for when one parameter is passed in
	expected = "localhost:8080"
	result = resolveAddress([]string{"localhost:8080"})
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}

	// Test case for when more than one parameter is passed in
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but did not get one")
		}
	}()
	resolveAddress([]string{"localhost:8080", "localhost:8081"})
}

func TestDefaultNotFound(t *testing.T) {
	req, _ := http.NewRequest("GET", "/foo", nil)

	// Create a new context with a mock response writer
	w := httptest.NewRecorder()
	ctx, _ := NewContext(w, req)

	// Call the defaultNotFound function
	defaultNotFound(ctx)
	ctx.Flush()

	// Verify that the response status code is 404
	if w.Code != http.StatusNotFound {
		t.Errorf("expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	// Verify that the response body is "Not Found"
	if w.Body.String() != http.StatusText(http.StatusNotFound) {
		t.Errorf("expected body %q, got %q", http.StatusText(http.StatusNotFound), w.Body.String())
	}
}

func TestDefaultInternalServerError(t *testing.T) {
	req, _ := http.NewRequest("GET", "/foo", nil)

	// Create a new context with a mock response writer
	w := httptest.NewRecorder()
	ctx, _ := NewContext(w, req)

	// Call the defaultInternalServerError function
	defaultInternalServerError(ctx)
	ctx.Flush()

	// Verify that the response status code is 500
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	// Verify that the response body is "Internal Server Error"
	if w.Body.String() != http.StatusText(http.StatusInternalServerError) {
		t.Errorf("expected body %q, got %q", http.StatusText(http.StatusInternalServerError), w.Body.String())
	}
}

func TestIsValidHTTPMethod(t *testing.T) {
	validMethods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	invalidMethods := []string{"", "get", "post", "put", "patch", "delete", "head", "options", "TRACE", "CONNECT"}

	for _, method := range validMethods {
		if !isValidHTTPMethod(method) {
			t.Errorf("isValidHTTPMethod(%q) = false, want true", method)
		}
	}

	for _, method := range invalidMethods {
		if isValidHTTPMethod(method) {
			t.Errorf("isValidHTTPMethod(%q) = true, want false", method)
		}
	}
}
