package lightning

import (
	"net/http"
	"reflect"
	"testing"
)

func TestCookie_Del(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		cookies cookiesMap
		args    args
	}{
		{
			name:    "TestCookie_Del",
			cookies: cookiesMap{"test": &http.Cookie{Name: "test", Value: "value"}},
			args:    args{key: "test"},
		},
		{
			name:    "TestCookie_Del_NotExist",
			cookies: cookiesMap{},
			args:    args{key: "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cookies.del(tt.args.key)
			if _, ok := tt.cookies[tt.args.key]; ok {
				t.Errorf("del() failed to delete key %s", tt.args.key)
			}
		})
	}
}

func TestCookie_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		cookies cookiesMap
		args    args
		want    *http.Cookie
	}{
		{
			name:    "TestCookie_Get",
			cookies: cookiesMap{"test": &http.Cookie{Name: "test", Value: "test"}},
			args:    args{key: "test"},
			want:    &http.Cookie{Name: "test", Value: "test"},
		},
		{
			name:    "TestCookie_Get_NotExist",
			cookies: cookiesMap{},
			args:    args{key: "test"},
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cookies.get(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCookie_Set(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name    string
		cookies cookiesMap
		args    args
	}{
		{
			name:    "TestCookie_Set",
			cookies: cookiesMap{},
			args:    args{key: "test", value: "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cookies.set(tt.args.key, tt.args.value)
			if _, ok := tt.cookies[tt.args.key]; !ok {
				t.Errorf("set() failed to set key %s", tt.args.key)
			}
			if tt.cookies[tt.args.key].Value != tt.args.value {
				t.Errorf("set() failed to set value %s for key %s", tt.args.value, tt.args.key)
			}
		})
	}
}

func TestCookie_SetCustom(t *testing.T) {
	type args struct {
		cookie *http.Cookie
	}
	tests := []struct {
		name    string
		cookies cookiesMap
		args    args
	}{
		{
			name:    "TestCookie_SetCustom",
			cookies: cookiesMap{},
			args:    args{cookie: &http.Cookie{Name: "test", Value: "test"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cookies.setCustom(tt.args.cookie)
			if _, ok := tt.cookies[tt.args.cookie.Name]; !ok {
				t.Errorf("setCustom() failed to set key %s", tt.args.cookie.Name)
			}
			if tt.cookies[tt.args.cookie.Name] != tt.args.cookie {
				t.Errorf("setCustom() failed to set value %s for key %s", tt.args.cookie.Value, tt.args.cookie.Name)
			}
		})
	}
}
