package lightning

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type errorReader struct{}

func (r *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func TestNewRequestWithBodyReadError(t *testing.T) {
	req := httptest.NewRequest("GET", "/path", &errorReader{})
	_, err := newRequest(req)
	if err == nil {
		t.Error("Expected error, but got nil")
	}
}

func TestNewRequest(t *testing.T) {
	type args struct {
		req    *http.Request
		params map[string]string
	}
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	params := map[string]string{"key1": "value1", "key2": "value2"}
	tests := []struct {
		name string
		args args
		want *request
	}{
		{
			name: "Test_NewRequest",
			args: args{
				req:    req,
				params: params,
			},
			want: &request{
				originReq: req,
				paramsMap: params,
				method:    req.Method,
				path:      req.URL.Path,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := newRequest(tt.args.req)
			got.setParams(tt.args.params)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_Cookie(t *testing.T) {
	type fields struct {
		req    *http.Request
		params map[string]string
		method string
		path   string
	}
	type args struct {
		name string
	}
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.AddCookie(&http.Cookie{Name: "cookie1", Value: "value1"})
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *http.Cookie
	}{
		{
			name: "Test_Request_Cookie",
			fields: fields{
				req: req,
			},
			args: args{
				name: "cookie1",
			},
			want: &http.Cookie{Name: "cookie1", Value: "value1"},
		},
		{
			name: "Test_Request_Cookie_Invalid",
			fields: fields{
				req: req,
			},
			args: args{
				name: "cookie_invalid",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &request{
				originReq: tt.fields.req,
				paramsMap: tt.fields.params,
				method:    tt.fields.method,
				path:      tt.fields.path,
			}
			if got := r.cookie(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cookie() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_Cookies(t *testing.T) {
	type fields struct {
		req    *http.Request
		params map[string]string
		method string
		path   string
	}
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	cookie1 := &http.Cookie{Name: "cookie1", Value: "value1"}
	cookie2 := &http.Cookie{Name: "cookie2", Value: "value2"}
	req.AddCookie(cookie1)
	req.AddCookie(cookie2)
	tests := []struct {
		name   string
		fields fields
		want   []*http.Cookie
	}{
		{
			name: "Test_Request_Cookies",
			fields: fields{
				req: req,
			},
			want: []*http.Cookie{cookie1, cookie2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &request{
				originReq: tt.fields.req,
				paramsMap: tt.fields.params,
				method:    tt.fields.method,
				path:      tt.fields.path,
			}
			if got := r.cookies(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cookies() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_Header(t *testing.T) {
	type fields struct {
		req    *http.Request
		params map[string]string
		method string
		path   string
	}
	type args struct {
		key string
	}
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("header1", "value1")
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Test_Request_Header",
			fields: fields{
				req: req,
			},
			args: args{
				key: "header1",
			},
			want: "value1",
		},
		{
			name: "Test_Request_Header_Invalid",
			fields: fields{
				req: req,
			},
			args: args{
				key: "header_invalid",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &request{
				originReq: tt.fields.req,
				paramsMap: tt.fields.params,
				method:    tt.fields.method,
				path:      tt.fields.path,
			}
			if got := r.header(tt.args.key); got != tt.want {
				t.Errorf("header() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_Headers(t *testing.T) {
	type fields struct {
		req    *http.Request
		params map[string]string
		method string
		path   string
	}
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("header1", "value1")
	req.Header.Set("header2", "value2")
	want := http.Header{"Header1": []string{"value1"}, "Header2": []string{"value2"}}
	tests := []struct {
		name   string
		fields fields
		want   http.Header
	}{
		{
			name: "Test_Request_Headers",
			fields: fields{
				req: req,
			},
			want: want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &request{
				originReq: tt.fields.req,
				paramsMap: tt.fields.params,
				method:    tt.fields.method,
				path:      tt.fields.path,
			}
			if got := r.headers(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("headers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_Param(t *testing.T) {
	type fields struct {
		req    *http.Request
		params map[string]string
		method string
		path   string
	}
	type args struct {
		key string
	}
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	params := map[string]string{"param1": "value1"}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Test_Request_Param",
			fields: fields{
				req:    req,
				params: params,
			},
			args: args{
				key: "param1",
			},
			want: "value1",
		},
		{
			name: "Test_Request_Param_Invalid",
			fields: fields{
				req:    req,
				params: params,
			},
			args: args{
				key: "param_invalid",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &request{
				originReq: tt.fields.req,
				paramsMap: tt.fields.params,
				method:    tt.fields.method,
				path:      tt.fields.path,
			}
			if got := r.param(tt.args.key); got != tt.want {
				t.Errorf("param() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_Params(t *testing.T) {
	type fields struct {
		req    *http.Request
		params map[string]string
		method string
		path   string
	}
	req, _ := http.NewRequest("GET", "<http://example.com?param1=value1&param2=value2>", nil)
	params := map[string]string{"param1": "value1", "param2": "value2"}
	tests := []struct {
		name   string
		fields fields
		want   map[string]string
	}{
		{
			name: "Test_Request_Params",
			fields: fields{
				req:    req,
				params: params,
			},
			want: params,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &request{
				originReq: tt.fields.req,
				paramsMap: tt.fields.params,
				method:    tt.fields.method,
				path:      tt.fields.path,
			}
			if got := r.params(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("params() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_Queries(t *testing.T) {
	type fields struct {
		req    *http.Request
		params map[string]string
		method string
		path   string
	}
	req, _ := http.NewRequest("GET", "http://example.com?param1=value1&param2=value2", nil)
	want := map[string][]string{"param1": []string{"value1"}, "param2": []string{"value2"}}
	tests := []struct {
		name   string
		fields fields
		want   map[string][]string
	}{
		{
			name: "Test_Request_Queries",
			fields: fields{
				req: req,
			},
			want: want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &request{
				originReq: tt.fields.req,
				paramsMap: tt.fields.params,
				method:    tt.fields.method,
				path:      tt.fields.path,
			}
			if got := r.queries(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("queries() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_Query(t *testing.T) {
	type fields struct {
		req    *http.Request
		params map[string]string
		method string
		path   string
	}
	type args struct {
		key string
	}
	req, _ := http.NewRequest("GET", "http://example.com?param1=value1", nil)
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Test_Request_Query",
			fields: fields{
				req: req,
			},
			args: args{
				key: "param1",
			},
			want: "value1",
		},
		{
			name: "Test_Request_Query_Invalid",
			fields: fields{
				req: req,
			},
			args: args{
				key: "param_invalid",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &request{
				originReq: tt.fields.req,
				paramsMap: tt.fields.params,
				method:    tt.fields.method,
				path:      tt.fields.path,
			}
			if got := r.query(tt.args.key); got != tt.want {
				t.Errorf("query() = %v, want %v", got, tt.want)
			}
		})
	}
}
