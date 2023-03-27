package lightning

import (
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
			if got := ParsePattern(tt.pattern); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParsePattern() = %v, want %v", got, tt.want)
			}
		})
	}
}
