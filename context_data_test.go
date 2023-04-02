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
		c    contextData
		args args
		want interface{}
	}{
		{
			name: "delete existing key",
			c:    contextData{"foo": "bar"},
			args: args{"foo"},
			want: contextData{},
		},
		{
			name: "delete non-existing key",
			c:    contextData{"foo": "bar"},
			args: args{"baz"},
			want: contextData{"foo": "bar"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.del(tt.args.key)
			if !reflect.DeepEqual(tt.c, tt.want) {
				t.Errorf("del() = %v, want %v", tt.c, tt.want)
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
		c    contextData
		args args
		want interface{}
	}{
		{
			name: "get existing key",
			c:    contextData{"foo": "bar"},
			args: args{"foo"},
			want: "bar",
		},
		{
			name: "get non-existing key",
			c:    contextData{"foo": "bar"},
			args: args{"baz"},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.get(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("get() = %v, want %v", got, tt.want)
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
		c    contextData
		args args
		want interface{}
	}{
		{
			name: "set new key-value pair in empty contextData",
			c:    contextData{},
			args: args{"foo", "bar"},
			want: contextData{"foo": "bar"},
		},
		{
			name: "set new key-value pair in non-empty contextData",
			c:    contextData{"foo": "bar"},
			args: args{"baz", 123},
			want: contextData{"foo": "bar", "baz": 123},
		},
		{
			name: "set existing key to new value",
			c:    contextData{"foo": "bar"},
			args: args{"foo", 123},
			want: contextData{"foo": 123},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.set(tt.args.key, tt.args.value)
			if !reflect.DeepEqual(tt.c, tt.want) {
				t.Errorf("set() = %v, want %v", tt.c, tt.want)
			}
		})
	}
}
