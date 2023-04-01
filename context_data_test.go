package lightning

import (
	"reflect"
	"testing"
)

func TestContextData_Del(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		c    ContextData
		args args
		want interface{}
	}{
		{
			name: "delete existing key",
			c:    ContextData{"foo": "bar"},
			args: args{"foo"},
			want: ContextData{},
		},
		{
			name: "delete non-existing key",
			c:    ContextData{"foo": "bar"},
			args: args{"baz"},
			want: ContextData{"foo": "bar"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.Del(tt.args.key)
			if !reflect.DeepEqual(tt.c, tt.want) {
				t.Errorf("Del() = %v, want %v", tt.c, tt.want)
			}
		})
	}
}

func TestContextData_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		c    ContextData
		args args
		want interface{}
	}{
		{
			name: "get existing key",
			c:    ContextData{"foo": "bar"},
			args: args{"foo"},
			want: "bar",
		},
		{
			name: "get non-existing key",
			c:    ContextData{"foo": "bar"},
			args: args{"baz"},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Get(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContextData_Set(t *testing.T) {
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name string
		c    ContextData
		args args
		want interface{}
	}{
		{
			name: "set new key-value pair in empty ContextData",
			c:    ContextData{},
			args: args{"foo", "bar"},
			want: ContextData{"foo": "bar"},
		},
		{
			name: "set new key-value pair in non-empty ContextData",
			c:    ContextData{"foo": "bar"},
			args: args{"baz", 123},
			want: ContextData{"foo": "bar", "baz": 123},
		},
		{
			name: "set existing key to new value",
			c:    ContextData{"foo": "bar"},
			args: args{"foo", 123},
			want: ContextData{"foo": 123},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.Set(tt.args.key, tt.args.value)
			if !reflect.DeepEqual(tt.c, tt.want) {
				t.Errorf("Set() = %v, want %v", tt.c, tt.want)
			}
		})
	}
}
