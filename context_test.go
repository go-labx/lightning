package lightning

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestNewContext(t *testing.T) {
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	ctx, err := NewContext(rr, req)
	if err != nil {
		t.Fatal(err)
	}
	if ctx.Method != "GET" {
		t.Errorf("Expected method to be GET, but got %s", ctx.Method)
	}
	if ctx.Path != "/test" {
		t.Errorf("Expected path to be /test, but got %s", ctx.Path)
	}
}

func TestNewContextWithError(t *testing.T) {
	req := httptest.NewRequest("GET", "/path", &errorReader{})
	rr := httptest.NewRecorder()
	_, err := NewContext(rr, req)
	if err == nil {
		t.Error("Expected error, but got nil")
	}
}

func TestSkipFlush(t *testing.T) {
	// Create a new request and response recorder
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Create a new context with the request and response recorder
	ctx, err := NewContext(rr, req)
	if err != nil {
		t.Fatal(err)
	}

	// Call the SkipFlush function
	ctx.SkipFlush()

	// Check if the skipFlush flag is set to true
	if !ctx.skipFlush {
		t.Errorf("SkipFlush did not set skipFlush flag to true")
	}
}

func TestContext_Next(t *testing.T) {
	// Create a new context with a mock handler function
	ctx := &Context{
		handlers: []HandlerFunc{
			func(c *Context) {
				// Do nothing
			},
		},
		index: -1,
	}

	// Call the Next method
	ctx.Next()

	// Check that the index has been incremented
	if ctx.index != 0 {
		t.Errorf("Expected index to be 0, but got %d", ctx.index)
	}
}

func TestContext_Flush(t *testing.T) {
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new mock response writer
	w := httptest.NewRecorder()

	// Create a new context object with the mock response writer
	ctx, err := NewContext(w, req)
	if err != nil {
		t.Fatal(err)
	}

	// Call the flush function
	ctx.flush()

	// Check that the response writer was flushed correctly
	if w.Code != http.StatusNotFound {
		t.Errorf("expected status code %d but got %d", http.StatusOK, w.Code)
	}
	if w.Body.String() != "" {
		t.Errorf("expected empty response body but got %s", w.Body.String())
	}
}

func TestRawBody(t *testing.T) {
	reqBody := "test request body"
	req, err := http.NewRequest("POST", "/test", strings.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()

	ctx, err := NewContext(res, req)
	if err != nil {
		t.Fatal(err)
	}

	rawBody := ctx.RawBody()
	expectedRawBody := []byte(reqBody)

	if !bytes.Equal(rawBody, expectedRawBody) {
		t.Errorf("RawBody() = %v, want %v", rawBody, expectedRawBody)
	}
}

func TestStringBody(t *testing.T) {
	req, err := http.NewRequest("GET", "/path", strings.NewReader("test body"))
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	ctx, err := NewContext(res, req)
	if err != nil {
		t.Fatal(err)
	}
	body := ctx.StringBody()
	if body != "test body" {
		t.Errorf("expected body to be 'test body', but got '%s'", body)
	}
}

func TestJSONBody(t *testing.T) {
	// Create a new context with a mock request and response
	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(`{"name": "John", "age": 30}`))
	res := httptest.NewRecorder()
	ctx, err := NewContext(res, req)
	if err != nil {
		t.Fatalf("Error creating context: %v", err)
	}

	// Define a struct to unmarshal the JSON into
	type Person struct {
		Name string `json:"name" validate:"required"`
		Age  int    `json:"age" validate:"gte=0"`
	}
	var p Person

	// Call the JSONBody function with the struct and validation flag
	err = ctx.JSONBody(&p, true)
	if err != nil {
		t.Fatalf("Error parsing JSON body: %v", err)
	}

	// Check that the struct was populated correctly
	if p.Name != "John" {
		t.Errorf("Expected name to be 'John', got '%s'", p.Name)
	}
	if p.Age != 30 {
		t.Errorf("Expected age to be 30, got %d", p.Age)
	}

	// Check that the function returns an error when given invalid JSON
	req = httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(`{"name": "John", "age": "thirty"}`))
	res = httptest.NewRecorder()
	ctx, err = NewContext(res, req)
	if err != nil {
		t.Fatalf("Error creating context: %v", err)
	}

	err = ctx.JSONBody(&p, true)
	if err == nil {
		t.Error("Expected error when parsing invalid JSON")
	}
}

func TestSetHandlers(t *testing.T) {
	req, err := http.NewRequest("POST", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new context with the request and response writer
	w := httptest.NewRecorder()

	ctx, err := NewContext(w, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	handlers := []HandlerFunc{
		func(c *Context) {},
		func(c *Context) {},
	}

	ctx.setHandlers(handlers)

	if len(ctx.handlers) != len(handlers) {
		t.Errorf("expected %d handlers, got %d", len(handlers), len(ctx.handlers))
	}
}

func TestContext_Param(t *testing.T) {
	req, err := http.NewRequest("GET", "/users/123", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	params := map[string]string{"id": "123"}
	ctx.setParams(params)

	got := ctx.Param("id")
	want := "123"
	if got != want {
		t.Errorf("ctx.Param(\"id\") = %q, want %q", got, want)
	}

	gotParams := ctx.Params()
	if !reflect.DeepEqual(gotParams, params) {
		t.Errorf("ctx.Params() = %q, want %q", gotParams, want)
	}
}

func TestContext_ParamInt(t *testing.T) {
	req, err := http.NewRequest("GET", "/users/123", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	params := map[string]string{"id": "123"}
	ctx.setParams(params)

	got, err := ctx.ParamInt("id")
	if err != nil {
		t.Errorf("ctx.ParamInt(\"id\") returned an error: %v", err)
	}
	want := 123
	if got != want {
		t.Errorf("ctx.ParamInt(\"id\") = %d, want %d", got, want)
	}
}

func TestContext_ParamIntWithException(t *testing.T) {
	req, err := http.NewRequest("GET", "/users/abc", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	params := map[string]string{"id": "abc"}
	ctx.setParams(params)

	_, err = ctx.ParamInt("id")
	if err == nil {
		t.Error("ctx.ParamInt(\"id\") did not return an error")
	}
}

func TestContext_ParamUInt(t *testing.T) {
	req, err := http.NewRequest("GET", "/users/123", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	params := map[string]string{"id": "123"}
	ctx.setParams(params)

	got, err := ctx.ParamUInt("id")
	if err != nil {
		t.Errorf("ctx.ParamUInt(\"id\") returned an error: %v", err)
	}
	want := uint(123)
	if got != want {
		t.Errorf("ctx.ParamUInt(\"id\") = %d, want %d", got, want)
	}
}

func TestContext_ParamUIntWithException(t *testing.T) {
	req, err := http.NewRequest("GET", "/users/abc", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	params := map[string]string{"id": "abc"}
	ctx.setParams(params)

	_, err = ctx.ParamUInt("id")
	if err == nil {
		t.Error("ctx.ParamUInt(\"id\") did not return an error")
	}
}

func TestContext_ParamInt64(t *testing.T) {
	req, err := http.NewRequest("GET", "/users/123", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	params := map[string]string{"id": "123"}
	ctx.setParams(params)

	got, err := ctx.ParamInt64("id")
	if err != nil {
		t.Errorf("ctx.ParamInt64(\"id\") returned an error: %v", err)
	}
	want := int64(123)
	if got != want {
		t.Errorf("ctx.ParamInt64(\"id\") = %d, want %d", got, want)
	}
}

func TestContext_ParamInt64WithException(t *testing.T) {
	req, err := http.NewRequest("GET", "/users/abc", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	params := map[string]string{"id": "abc"}
	ctx.setParams(params)

	_, err = ctx.ParamInt64("id")
	if err == nil {
		t.Error("ctx.ParamInt64(\"id\") did not return an error")
	}
}

func TestContext_ParamUInt64(t *testing.T) {
	req, err := http.NewRequest("GET", "/users/123", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	params := map[string]string{"id": "123"}
	ctx.setParams(params)

	got, err := ctx.ParamUInt64("id")
	if err != nil {
		t.Errorf("ctx.ParamUInt64(\"id\") returned an error: %v", err)
	}
	want := uint64(123)
	if got != want {
		t.Errorf("ctx.ParamUInt64(\"id\") = %d, want %d", got, want)
	}
}

func TestContext_ParamUInt64WithException(t *testing.T) {
	req, err := http.NewRequest("GET", "/users/abc", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	params := map[string]string{"id": "abc"}
	ctx.setParams(params)

	_, err = ctx.ParamUInt64("id")
	if err == nil {
		t.Error("ctx.ParamUInt64(\"id\") did not return an error")
	}
}

func TestContext_ParamFloat32(t *testing.T) {
	req, err := http.NewRequest("GET", "/users/123.456", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	params := map[string]string{"id": "123.456"}
	ctx.setParams(params)

	got, err := ctx.ParamFloat32("id")
	if err != nil {
		t.Errorf("ctx.ParamFloat32(\"id\") returned an error: %v", err)
	}
	want := float32(123.456)
	if got != want {
		t.Errorf("ctx.ParamFloat32(\"id\") = %f, want %f", got, want)
	}
}

func TestContext_ParamFloat32WithException(t *testing.T) {
	req, err := http.NewRequest("GET", "/users/abc", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	params := map[string]string{"id": "abc"}
	ctx.setParams(params)

	_, err = ctx.ParamFloat32("id")
	if err == nil {
		t.Error("ctx.ParamFloat32(\"id\") did not return an error")
	}
}

func TestContext_ParamFloat64(t *testing.T) {
	req, err := http.NewRequest("GET", "/users/123.456", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	params := map[string]string{"id": "123.456"}
	ctx.setParams(params)

	got, err := ctx.ParamFloat64("id")
	if err != nil {
		t.Errorf("ctx.ParamFloat64(\"id\") returned an error: %v", err)
	}
	want := float64(123.456)
	if got != want {
		t.Errorf("ctx.ParamFloat64(\"id\") = %f, want %f", got, want)
	}
}

func TestContext_ParamFloat64WithException(t *testing.T) {
	req, err := http.NewRequest("GET", "/users/abc", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	params := map[string]string{"id": "abc"}
	ctx.setParams(params)

	_, err = ctx.ParamFloat64("id")
	if err == nil {
		t.Error("ctx.ParamFloat64(\"id\") did not return an error")
	}
}

func TestContext_ParamString(t *testing.T) {
	req, err := http.NewRequest("GET", "/users/123", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	params := map[string]string{"id": "123"}
	ctx.setParams(params)

	got := ctx.ParamString("id")
	want := "123"
	if got != want {
		t.Errorf("ctx.ParamString(\"id\") = %s, want %s", got, want)
	}
}

func TestContext_Query(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?key=value", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	got := ctx.Query("key")
	want := "value"
	if got != want {
		t.Errorf("Query() = %q, want %q", got, want)
	}
}

func TestContext_QueryString(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?key=value", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	got := ctx.QueryString("key")
	want := "value"
	if got != want {
		t.Errorf("QueryString() = %q, want %q", got, want)
	}
}

func TestContext_QueryBool(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?key=true", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ctx.QueryBool("key")
	if err != nil {
		t.Errorf("ctx.QueryBool(\"key\") returned an error: %v", err)
	}
	want := true
	if got != want {
		t.Errorf("QueryBool() = %v, want %v", got, want)
	}
}

func TestContext_QueryBoolWithException(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?key=notabool", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ctx.QueryBool("key")
	if err == nil {
		t.Error("ctx.QueryBool(\"key\") did not return an error")
	}
}

func TestContext_QueryBoolWithEmptyKey(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?key=", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ctx.QueryBool("key")
	if err != nil {
		t.Errorf("ctx.QueryBool(\"key\") returned an error: %v", err)
	}
	want := false
	if got != want {
		t.Errorf("QueryBool() = %v, want %v", got, want)
	}
}

func TestContext_QueryInt(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?key=123", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ctx.QueryInt("key")
	if err != nil {
		t.Errorf("ctx.QueryInt(\"key\") returned an error: %v", err)
	}
	want := 123
	if got != want {
		t.Errorf("QueryInt() = %q, want %q", got, want)
	}
}

func TestContext_QueryIntWithException(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?key=abc", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ctx.QueryInt("key")
	if err == nil {
		t.Error("ctx.QueryInt(\"id\") did not return an error")
	}
}

func TestContext_QueryIntWithEmptyKey(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?key=", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ctx.QueryInt("key")
	if err != nil {
		t.Errorf("ctx.QueryInt(\"key\") returned an error: %v", err)
	}
	want := 0
	if got != want {
		t.Errorf("QueryInt() = %d, want %d", got, want)
	}
}

func TestContext_QueryUInt(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?key=123", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ctx.QueryUInt("key")
	if err != nil {
		t.Errorf("ctx.QueryUInt(\"key\") returned an error: %v", err)
	}
	want := uint(123)
	if got != want {
		t.Errorf("QueryUInt() = %q, want %q", got, want)
	}
}

func TestContext_QueryUIntWithException(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?key=abc", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ctx.QueryUInt("key")
	if err == nil {
		t.Error("ctx.QueryUInt(\"id\") did not return an error")
	}
}

func TestContext_QueryUIntWithEmptyKey(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?key=", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ctx.QueryUInt("key")
	if err != nil {
		t.Errorf("ctx.QueryUInt(\"key\") returned an error: %v", err)
	}
	want := uint(0)
	if got != want {
		t.Errorf("QueryUInt() = %d, want %d", got, want)
	}
}

func TestContext_QueryInt8(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?key=123", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ctx.QueryInt8("key")
	if err != nil {
		t.Errorf("ctx.QueryInt8(\"key\") returned an error: %v", err)
	}
	want := int8(123)
	if got != want {
		t.Errorf("QueryInt8() = %q, want %q", got, want)
	}
}

func TestContext_QueryInt8WithException(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?key=abc", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ctx.QueryInt8("key")
	if err == nil {
		t.Error("ctx.QueryInt8(\"id\") did not return an error")
	}
}

func TestContext_QueryInt8WithEmptyKey(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?key=", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ctx.QueryInt8("key")
	if err != nil {
		t.Errorf("ctx.QueryInt8(\"key\") returned an error: %v", err)
	}
	want := int8(0)
	if got != want {
		t.Errorf("QueryInt8() = %d, want %d", got, want)
	}
}

func TestContext_QueryUInt8(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?key=123", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ctx.QueryUInt8("key")
	if err != nil {
		t.Errorf("ctx.QueryUInt8(\"key\") returned an error: %v", err)
	}
	want := uint8(123)
	if got != want {
		t.Errorf("QueryUInt8() = %q, want %q", got, want)
	}
}

func TestContext_QueryUInt8WithException(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?key=abc", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ctx.QueryUInt8("key")
	if err == nil {
		t.Error("ctx.QueryUInt8(\"id\") did not return an error")
	}
}

func TestContext_QueryUInt8WithEmptyKey(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?key=", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ctx.QueryUInt8("key")
	if err != nil {
		t.Errorf("ctx.QueryUInt8(\"key\") returned an error: %v", err)
	}
	want := uint8(0)
	if got != want {
		t.Errorf("QueryUInt8() = %d, want %d", got, want)
	}
}

func TestContext_QueryInt32(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?key=123", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ctx.QueryInt32("key")
	if err != nil {
		t.Errorf("ctx.QueryInt32(\"key\") returned an error: %v", err)
	}
	want := int32(123)
	if got != want {
		t.Errorf("QueryInt32() = %q, want %q", got, want)
	}
}

func TestContext_QueryInt32WithException(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?key=abc", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ctx.QueryInt32("key")
	if err == nil {
		t.Error("ctx.QueryInt32(\"id\") did not return an error")
	}
}

func TestContext_QueryInt32WithEmptyKey(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?key=", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ctx.QueryInt32("key")
	if err != nil {
		t.Errorf("ctx.QueryInt32(\"key\") returned an error: %v", err)
	}
	want := int32(0)
	if got != want {
		t.Errorf("QueryInt32() = %d, want %d", got, want)
	}
}

func TestContext_QueryUInt32(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?key=123", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ctx.QueryUInt32("key")
	if err != nil {
		t.Errorf("ctx.QueryUInt32(\"key\") returned an error: %v", err)
	}
	want := uint32(123)
	if got != want {
		t.Errorf("QueryUInt32() = %q, want %q", got, want)
	}
}

func TestContext_QueryUInt32WithException(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?key=abc", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ctx.QueryUInt32("key")
	if err == nil {
		t.Error("ctx.QueryUInt32(\"id\") did not return an error")
	}
}

func TestContext_QueryUInt32WithEmptyKey(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?key=", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ctx.QueryUInt32("key")
	if err != nil {
		t.Errorf("ctx.QueryUInt32(\"key\") returned an error: %v", err)
	}
	want := uint32(0)
	if got != want {
		t.Errorf("QueryUInt32() = %d, want %d", got, want)
	}
}

func TestContext_QueryInt64(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?key=123", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ctx.QueryInt64("key")
	if err != nil {
		t.Errorf("ctx.QueryInt64(\"key\") returned an error: %v", err)
	}
	want := int64(123)
	if got != want {
		t.Errorf("QueryInt64() = %q, want %q", got, want)
	}
}

func TestContext_QueryInt64WithException(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?key=abc", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ctx.QueryInt64("key")
	if err == nil {
		t.Error("ctx.QueryInt64(\"id\") did not return an error")
	}
}

func TestContext_QueryInt64WithEmptyKey(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?key=", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ctx.QueryInt64("key")
	if err != nil {
		t.Errorf("ctx.QueryInt64(\"key\") returned an error: %v", err)
	}
	want := int64(0)
	if got != want {
		t.Errorf("QueryInt64() = %d, want %d", got, want)
	}
}

func TestContext_QueryUInt64(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?key=123", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ctx.QueryUInt64("key")
	if err != nil {
		t.Errorf("ctx.QueryUInt64(\"key\") returned an error: %v", err)
	}
	want := uint64(123)
	if got != want {
		t.Errorf("QueryUInt64() = %q, want %q", got, want)
	}
}

func TestContext_QueryUInt64WithException(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?key=abc", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ctx.QueryUInt64("key")
	if err == nil {
		t.Error("ctx.QueryUInt64(\"id\") did not return an error")
	}
}

func TestContext_QueryUInt64WithEmptyKey(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?key=", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ctx.QueryUInt64("key")
	if err != nil {
		t.Errorf("ctx.QueryUInt64(\"key\") returned an error: %v", err)
	}
	want := uint64(0)
	if got != want {
		t.Errorf("QueryUInt64() = %d, want %d", got, want)
	}
}

func TestContext_QueryFloat32(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?key=3.1415", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ctx.QueryFloat32("key")
	if err != nil {
		t.Errorf("ctx.QueryFloat32(\"key\") returned an error: %v", err)
	}
	want := float32(3.1415)
	if got != want {
		t.Errorf("QueryFloat32() = %f, want %f", got, want)
	}
}

func TestContext_QueryFloat32WithException(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?key=abc", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ctx.QueryFloat32("key")
	if err == nil {
		t.Error("ctx.QueryFloat32(\"id\") did not return an error")
	}
}

func TestContext_QueryFloat32WithEmptyKey(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?key=", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ctx.QueryFloat32("key")
	if err != nil {
		t.Errorf("ctx.QueryFloat32(\"key\") returned an error: %v", err)
	}
	want := float32(0)
	if got != want {
		t.Errorf("QueryFloat32() = %f, want %f", got, want)
	}
}

func TestContext_QueryFloat64(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?key=3.1415", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ctx.QueryFloat64("key")
	if err != nil {
		t.Errorf("ctx.QueryFloat64(\"key\") returned an error: %v", err)
	}
	want := float64(3.1415)
	if got != want {
		t.Errorf("QueryFloat64() = %f, want %f", got, want)
	}
}

func TestContext_QueryFloat64WithException(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?key=abc", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ctx.QueryFloat64("key")
	if err == nil {
		t.Error("ctx.QueryFloat64(\"id\") did not return an error")
	}
}

func TestContext_QueryFloat64WithEmptyKey(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?key=", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ctx.QueryFloat64("key")
	if err != nil {
		t.Errorf("ctx.QueryFloat64(\"key\") returned an error: %v", err)
	}
	want := float64(0)
	if got != want {
		t.Errorf("QueryFloat64() = %f, want %f", got, want)
	}
}

func TestContextQueries(t *testing.T) {
	req, err := http.NewRequest("GET", "/path?foo=bar&baz=qux", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := NewContext(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	queries := ctx.Queries()
	expected := map[string][]string{
		"foo": {"bar"},
		"baz": {"qux"},
	}
	if !reflect.DeepEqual(queries, expected) {
		t.Errorf("got %v, want %v", queries, expected)
	}
}

func TestContext_Status(t *testing.T) {
	// Create a new HTTP request and response
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Create a new context object
	ctx, err := NewContext(rr, req)
	if err != nil {
		t.Fatal(err)
	}

	// Test the Status method
	status := ctx.Status()
	if status != http.StatusNotFound {
		t.Errorf("Expected status code %d, but got %d", http.StatusNotFound, status)
	}

	// Test the SetStatus method
	ctx.SetStatus(http.StatusOK)
	status = ctx.Status()
	if status != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, status)
	}
}

func TestContext_SetStatus(t *testing.T) {
	// Create a new context object
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	ctx, err := NewContext(rr, req)
	if err != nil {
		t.Fatal(err)
	}

	// Call the SetStatus method with a status code of 200
	ctx.SetStatus(http.StatusOK)

	// Check that the response status code is 200
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestHeader(t *testing.T) {
	req, err := http.NewRequest("GET", "/path", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	ctx, err := NewContext(w, req)
	if err != nil {
		t.Fatal(err)
	}
	key := "Content-Type"
	value := "application/json"
	req.Header.Set(key, value)
	if got := ctx.Header(key); got != value {
		t.Errorf("Header(%q) = %q, want %q", key, got, value)
	}
}

func TestHeaders(t *testing.T) {
	req, err := http.NewRequest("GET", "/path", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	ctx, err := NewContext(rr, req)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Request-ID", "12345")
	headers := ctx.Headers()
	if len(headers) != 2 {
		t.Errorf("Expected 2 headers, but got %d", len(headers))
	}
	if headers.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type header to be 'application/json', but got '%s'", headers.Get("Content-Type"))
	}
	if headers.Get("X-Request-ID") != "12345" {
		t.Errorf("Expected X-Request-ID header to be '12345', but got '%s'", headers.Get("X-Request-ID"))
	}
}

func TestAddHeader(t *testing.T) {
	// Create a new context object
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	ctx, err := NewContext(res, req)
	if err != nil {
		t.Fatal(err)
	}

	// Call the AddHeader method
	ctx.AddHeader("Content-Type", "application/json")

	// Check the response headers
	headers := res.Header()
	contentType := headers.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("unexpected content type: got %v want %v", contentType, "application/json")
	}
}

func TestSetHeader(t *testing.T) {
	// Create a new context object
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	ctx, err := NewContext(res, req)
	if err != nil {
		t.Fatal(err)
	}

	// Call the SetHeader method
	ctx.SetHeader("Content-Type", "application/json")

	// Check the response headers
	headers := res.Header()
	contentType := headers.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("unexpected content type: got %v want %v", contentType, "application/json")
	}
}

func TestDelHeader(t *testing.T) {
	// Create a new context object
	req, err := http.NewRequest("GET", "/path", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	ctx, err := NewContext(res, req)
	if err != nil {
		t.Fatal(err)
	}

	// Set a header in the response
	ctx.SetHeader("key", "value")

	// Call the DelHeader method
	ctx.DelHeader("key")

	// Check that the header was deleted
	if res.Header().Get("key") != "" {
		t.Errorf("Header was not deleted")
	}
}

func TestContextCookie(t *testing.T) {
	req, err := http.NewRequest("GET", "/path", nil)
	if err != nil {
		t.Fatal(err)
	}
	cookie := &http.Cookie{Name: "test", Value: "value"}
	req.AddCookie(cookie)

	res := httptest.NewRecorder()
	ctx, err := NewContext(res, req)
	if err != nil {
		t.Fatal(err)
	}

	if c := ctx.Cookie("test"); c == nil || c.Value != "value" {
		t.Errorf("Cookie() = %v, want %v", c, cookie)
	}
}

func TestCookies(t *testing.T) {
	req, err := http.NewRequest("GET", "/path", nil)
	if err != nil {
		t.Fatal(err)
	}
	cookie := &http.Cookie{Name: "test", Value: "value"}
	req.AddCookie(cookie)
	w := httptest.NewRecorder()
	ctx, err := NewContext(w, req)
	if err != nil {
		t.Fatal(err)
	}

	cookies := ctx.Cookies()
	if len(cookies) != 1 {
		t.Errorf("Expected 1 cookie, got %d", len(cookies))
	}
	if cookies[0].Name != "test" {
		t.Errorf("Expected cookie name 'test', got '%s'", cookies[0].Name)
	}
	if cookies[0].Value != "value" {
		t.Errorf("Expected cookie value 'value', got '%s'", cookies[0].Value)
	}
}

func TestSetCookie(t *testing.T) {
	// Create a new context object
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	ctx, err := NewContext(w, req)
	if err != nil {
		t.Fatal(err)
	}

	// Call the SetCookie function
	ctx.SetCookie("test", "value")
	ctx.flush()

	// Check that the cookie was set correctly
	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Errorf("expected 1 cookie, got %d", len(cookies))
	}
	if cookies[0].Name != "test" {
		t.Errorf("expected cookie name 'test', got '%s'", cookies[0].Name)
	}
	if cookies[0].Value != "value" {
		t.Errorf("expected cookie value 'value', got '%s'", cookies[0].Value)
	}
}

func TestSetCustomCookie(t *testing.T) {
	// Create a new context object
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	ctx, err := NewContext(w, req)
	if err != nil {
		t.Fatal(err)
	}

	// Call the SetCustomCookie function
	cookie := &http.Cookie{Name: "test", Value: "value"}
	ctx.SetCustomCookie(cookie)
	ctx.flush()

	// Check that the cookie was set correctly
	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Errorf("expected 1 cookie, got %d", len(cookies))
	}
	if cookies[0].Name != "test" {
		t.Errorf("expected cookie name 'test', got '%s'", cookies[0].Name)
	}
	if cookies[0].Value != "value" {
		t.Errorf("expected cookie value 'value', got '%s'", cookies[0].Value)
	}
}

func TestContextBody(t *testing.T) {
	// create a new context with a response body
	body := []byte("test body")
	ctx := &Context{
		res: &response{
			body: body,
		},
	}

	// call the Body() function and check the result
	result := ctx.Body()
	if !bytes.Equal(result, body) {
		t.Errorf("expected body %v, but got %v", body, result)
	}
}

func TestContextSetBody(t *testing.T) {
	// create a new context
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	ctx, _ := NewContext(res, req)

	// set the body using SetBody
	body := []byte("test body")
	ctx.SetBody(body)

	// check that the body was set correctly
	if !bytes.Equal(ctx.res.body, body) {
		t.Errorf("expected body %v, got %v", body, ctx.res.body)
	}
}

func TestJSON(t *testing.T) {
	// Create a new context object
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	ctx, err := NewContext(w, req)
	if err != nil {
		t.Fatal(err)
	}

	ctx.JSON(200, map[string]string{"message": "hello world"})
	ctx.flush()

	if ctx.Status() != 200 {
		t.Errorf("Expected status code %d but got %d", 200, ctx.Status())
	}

	expectedBody := `{"message":"hello world"}`
	if string(ctx.res.body) != expectedBody {
		t.Errorf("Expected body %s but got %s", expectedBody, string(ctx.res.body))
	}
}

func TestText(t *testing.T) {
	// Create a new context object
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	ctx, err := NewContext(w, req)
	if err != nil {
		t.Fatal(err)
	}

	ctx.Text(200, "hello world")
	ctx.flush()

	if ctx.Status() != 200 {
		t.Errorf("Expected status code %d but got %d", 200, ctx.Status())
	}

	expectedBody := "hello world"
	if string(ctx.res.body) != expectedBody {
		t.Errorf("Expected body %s but got %s", expectedBody, string(ctx.res.body))
	}
}

func TestContext_XML(t *testing.T) {
	// Create a new context
	req := httptest.NewRequest("GET", "/xml", nil)
	res := httptest.NewRecorder()
	ctx, err := NewContext(res, req)
	if err != nil {
		t.Fatal(err)
	}

	// Create a test object
	type person struct {
		Name string `xml:"name"`
		Age  int    `xml:"age"`
	}
	obj := &person{
		Name: "John",
		Age:  30,
	}

	// Call the XML function with the test object
	ctx.XML(http.StatusOK, obj)
	ctx.flush()

	// Check the response headers
	if res.Header().Get(HeaderContentType) != MIMEApplicationXML {
		t.Errorf("Expected Content-Type header to be %s, but got %s", MIMEApplicationXML, res.Header().Get(HeaderContentType))
	}

	// Check the response status code
	if res.Result().StatusCode != http.StatusOK {
		t.Errorf("Expected status code to be %d, but got %d", http.StatusOK, res.Result().StatusCode)
	}

	// Check the response body
	expectedBody := `<person><name>John</name><age>30</age></person>`
	if res.Body.String() != expectedBody {
		t.Errorf("Expected response body to be %s, but got %s", expectedBody, res.Body.String())
	}
}

func TestContext_GetData(t *testing.T) {
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	ctx, err := NewContext(rr, req)
	if err != nil {
		t.Fatal(err)
	}

	// Test getting a non-existent key
	if ctx.GetData("nonexistent") != nil {
		t.Errorf("expected nil value for nonexistent key, got %v", ctx.GetData("nonexistent"))
	}

	// Test setting and getting a key
	ctx.SetData("key", "value")
	if ctx.GetData("key") != "value" {
		t.Errorf("expected value 'value' for key 'key', got %v", ctx.GetData("key"))
	}

	// Test deleting a key
	ctx.DelData("key")
	if ctx.GetData("key") != nil {
		t.Errorf("expected nil value for deleted key 'key', got %v", ctx.GetData("key"))
	}
}

func TestContext_Redirect(t *testing.T) {
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	ctx, err := NewContext(rr, req)
	if err != nil {
		t.Fatal(err)
	}

	// Call the Redirect method with a test URL and status code
	redirectUrl := "https://example.com"
	ctx.Redirect(http.StatusMovedPermanently, redirectUrl)
	ctx.flush()

	// Verify that the response status code and location header are set correctly
	if rr.Result().StatusCode != http.StatusMovedPermanently {
		t.Errorf("expected status code %d, got %d", http.StatusMovedPermanently, rr.Result().StatusCode)
	}
	url := rr.Header().Get("Location")
	if url != redirectUrl {
		t.Errorf("expected Location header %q, got %q", redirectUrl, url)
	}
}

func TestUserAgent(t *testing.T) {
	ctx := createMockContext(t)
	ua := "my-user-agent"
	ctx.req.originReq.Header.Add("user-agent", ua)

	if userAgent := ctx.UserAgent(); userAgent != "my-user-agent" {
		t.Errorf("expected user agent %q, got %q", ua, userAgent)
	}
}

func TestReferer(t *testing.T) {
	ctx := createMockContext(t)
	ref := "https://example.com"
	ctx.req.originReq.Header.Add("referer", ref)

	if referer := ctx.Referer(); referer != ref {
		t.Errorf("expected referer %q, got %q", ref, referer)
	}
}

func TestContext_RemoteAddr(t *testing.T) {
	ctx := createMockContext(t)
	expected := "1.2.3.4:5678"
	ctx.req.originReq.RemoteAddr = expected

	if addr := ctx.RemoteAddr(); addr != expected {
		t.Errorf("Expected RemoteAddr to return %q, but got %q", expected, addr)
	}
}

func createMockContext(t *testing.T) *Context {
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	ctx, err := NewContext(rr, req)
	if err != nil {
		t.Fatal(err)
	}

	return ctx
}

func TestContext_Success(t *testing.T) {
	// Create a new request with an empty body
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Create a new context with the request and response recorder
	ctx, err := NewContext(rr, req)
	if err != nil {
		t.Fatal(err)
	}

	// Call the Success method with some test data
	testData := map[string]string{"foo": "bar"}
	ctx.Success(testData)
	ctx.flush()

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	expected := `{"code":0,"data":{"foo":"bar"},"message":"ok"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestContextFail(t *testing.T) {
	// Create a new context with a mock response writer and request
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	ctx, err := NewContext(w, req)
	if err != nil {
		t.Fatal(err)
	}

	// Call the Fail method with a custom code and message
	ctx.Fail(500, "Internal Server Error")
	ctx.flush()

	// Check that the response status code and body are correct
	if w.Code != 200 {
		t.Errorf("expected status code 200, got %d", w.Code)
	}
	expectedBody := `{"code":500,"message":"Internal Server Error"}`
	if w.Body.String() != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, w.Body.String())
	}
}
