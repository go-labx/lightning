package lightning

import (
	"os"
	"reflect"
	"testing"

	"github.com/valyala/fasthttp"
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
	os.Setenv("PORT", "1234")
	defer os.Unsetenv("PORT")
	expected := ":1234"
	result := resolveAddress([]string{})
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
	os.Unsetenv("PORT")

	expected = ":6789"
	result = resolveAddress([]string{})
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}

	expected = "localhost:8080"
	result = resolveAddress([]string{"localhost:8080"})
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but did not get one")
		}
	}()
	resolveAddress([]string{"localhost:8080", "localhost:8081"})
}

func TestDefaultNotFound(t *testing.T) {
	fctx := &fasthttp.RequestCtx{}
	fctx.Request.Header.SetMethod("GET")
	fctx.Request.Header.SetRequestURI("/foo")

	c := &Context{
		ctx:   fctx,
		index: -1,
		res:   newResponse(fctx),
	}
	defaultNotFound(c)
	c.flush()

	if fctx.Response.StatusCode() != StatusNotFound {
		t.Errorf("expected status code %d, got %d", StatusNotFound, fctx.Response.StatusCode())
	}

	if string(fctx.Response.Body()) != "Not Found" {
		t.Errorf("expected body %q, got %q", "Not Found", string(fctx.Response.Body()))
	}
}

func TestDefaultInternalServerError(t *testing.T) {
	fctx := &fasthttp.RequestCtx{}
	fctx.Request.Header.SetMethod("GET")
	fctx.Request.Header.SetRequestURI("/foo")

	c := &Context{
		ctx:   fctx,
		index: -1,
		res:   newResponse(fctx),
	}
	defaultInternalServerError(c)
	c.flush()

	if fctx.Response.StatusCode() != StatusInternalServerError {
		t.Errorf("expected status code %d, got %d", StatusInternalServerError, fctx.Response.StatusCode())
	}

	if string(fctx.Response.Body()) != "Internal Server Error" {
		t.Errorf("expected body %q, got %q", "Internal Server Error", string(fctx.Response.Body()))
	}
}
