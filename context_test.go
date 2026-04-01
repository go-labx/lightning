package lightning

import (
	"bytes"
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
	"text/template"

	"github.com/valyala/fasthttp"
)

func createTestContext(method, path string, body []byte) (*Context, *fasthttp.RequestCtx) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(method)
	ctx.Request.Header.SetRequestURI(path)
	if body != nil {
		ctx.Request.SetBody(body)
	}

	c := &Context{
		ctx:   ctx,
		index: -1,
		data:  contextData{},
	}
	c.req = newRequest(ctx)
	c.res = newResponse(ctx)
	c.Method = c.req.method()
	c.Path = c.req.path()

	return c, ctx
}

func newTestCtxForApp(method, path string) *fasthttp.RequestCtx {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(method)
	ctx.Request.Header.SetRequestURI(path)
	return ctx
}

func TestNewContext(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.Header.SetRequestURI("/test")

	c := NewContext(ctx)
	if c == nil {
		t.Fatal("NewContext returned nil")
	}
	if c.ctx != ctx {
		t.Error("RequestCtx not set correctly")
	}
}

func TestContextReset(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.Header.SetRequestURI("/test")

	c := NewContext(ctx)
	c.reset()

	if c.ctx != nil {
		t.Errorf("Expected ctx to be nil after reset")
	}
	if c.req != nil {
		t.Errorf("Expected req to be nil after reset")
	}
	if c.res != nil {
		t.Errorf("Expected res to be nil after reset")
	}
	if c.handlers != nil {
		t.Errorf("Expected handlers to be nil after reset")
	}
	if c.index != -1 {
		t.Errorf("Expected index to be -1 after reset, got %d", c.index)
	}
}

func TestContext_Next(t *testing.T) {
	called := false
	ctx := &Context{
		handlers: []HandlerFunc{
			func(c *Context) {
				called = true
			},
		},
		index: -1,
	}

	ctx.Next()

	if !called {
		t.Error("Handler was not called")
	}
	if ctx.index != 0 {
		t.Errorf("Expected index to be 0, but got %d", ctx.index)
	}
}

func TestContext_NextMultiple(t *testing.T) {
	order := []int{}
	ctx := &Context{
		handlers: []HandlerFunc{
			func(c *Context) {
				order = append(order, 1)
				c.Next()
				order = append(order, 4)
			},
			func(c *Context) {
				order = append(order, 2)
				c.Next()
				order = append(order, 3)
			},
		},
		index: -1,
	}

	ctx.Next()

	expected := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(order, expected) {
		t.Errorf("Expected order %v, got %v", expected, order)
	}
}

func TestRawBody(t *testing.T) {
	reqBody := []byte("test request body")
	c, _ := createTestContext("POST", "/test", reqBody)

	rawBody := c.RawBody()

	if !bytes.Equal(rawBody, reqBody) {
		t.Errorf("RawBody() = %v, want %v", rawBody, reqBody)
	}
}

func TestStringBody(t *testing.T) {
	reqBody := []byte("test body")
	c, _ := createTestContext("POST", "/test", reqBody)

	body := c.StringBody()
	if body != "test body" {
		t.Errorf("expected body to be 'test body', but got '%s'", body)
	}
}

func TestJSONBody(t *testing.T) {
	c, _ := createTestContext("POST", "/test", []byte(`{"name": "John", "age": 30}`))

	type Person struct {
		Name string `json:"name" validate:"required"`
		Age  int    `json:"age" validate:"gte=0"`
	}
	var p Person

	err := c.JSONBody(&p, true)
	if err != nil {
		t.Fatalf("Error parsing JSON body: %v", err)
	}

	if p.Name != "John" {
		t.Errorf("Expected name to be 'John', got '%s'", p.Name)
	}
	if p.Age != 30 {
		t.Errorf("Expected age to be 30, got %d", p.Age)
	}
}

func TestJSONBodyValidation(t *testing.T) {
	c, _ := createTestContext("POST", "/test", []byte(`{"name": "John", "age": "thirty"}`))

	type Person struct {
		Name string `json:"name" validate:"required"`
		Age  int    `json:"age" validate:"gte=0"`
	}
	var p Person

	err := c.JSONBody(&p, true)
	if err == nil {
		t.Error("Expected error when parsing invalid JSON")
	}
}

func TestJSONBodyInvalidJSON(t *testing.T) {
	c, _ := createTestContext("POST", "/test", []byte(`invalid json`))

	type Person struct {
		Name string `json:"name"`
	}
	var p Person

	err := c.JSONBody(&p)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestSetHandlers(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)

	handlers := []HandlerFunc{
		func(c *Context) {},
		func(c *Context) {},
	}

	c.setHandlers(handlers)

	if len(c.handlers) != len(handlers) {
		t.Errorf("expected %d handlers, got %d", len(handlers), len(c.handlers))
	}
}

func TestContext_Param(t *testing.T) {
	c, _ := createTestContext("GET", "/users/123", nil)
	params := map[string]string{"id": "123"}
	c.setParams(params)

	got := c.Param("id")
	want := "123"
	if got != want {
		t.Errorf("ctx.Param(\"id\") = %q, want %q", got, want)
	}

	gotParams := c.Params()
	if !reflect.DeepEqual(gotParams, params) {
		t.Errorf("ctx.Params() = %v, want %v", gotParams, params)
	}
}

func TestContext_ParamInt(t *testing.T) {
	c, _ := createTestContext("GET", "/users/123", nil)
	params := map[string]string{"id": "123"}
	c.setParams(params)

	got, err := c.ParamInt("id")
	if err != nil {
		t.Errorf("ctx.ParamInt(\"id\") returned an error: %v", err)
	}
	want := 123
	if got != want {
		t.Errorf("ctx.ParamInt(\"id\") = %d, want %d", got, want)
	}
}

func TestContext_ParamIntWithException(t *testing.T) {
	c, _ := createTestContext("GET", "/users/abc", nil)
	params := map[string]string{"id": "abc"}
	c.setParams(params)

	_, err := c.ParamInt("id")
	if err == nil {
		t.Error("ctx.ParamInt(\"id\") did not return an error")
	}
}

func TestContext_ParamInt64(t *testing.T) {
	c, _ := createTestContext("GET", "/users/123", nil)
	params := map[string]string{"id": "123"}
	c.setParams(params)

	got, err := c.ParamInt64("id")
	if err != nil {
		t.Errorf("ctx.ParamInt64(\"id\") returned an error: %v", err)
	}
	want := int64(123)
	if got != want {
		t.Errorf("ctx.ParamInt64(\"id\") = %d, want %d", got, want)
	}
}

func TestContext_ParamUInt(t *testing.T) {
	c, _ := createTestContext("GET", "/users/123", nil)
	params := map[string]string{"id": "123"}
	c.setParams(params)

	got, err := c.ParamUInt("id")
	if err != nil {
		t.Errorf("ctx.ParamUInt(\"id\") returned an error: %v", err)
	}
	want := uint(123)
	if got != want {
		t.Errorf("ctx.ParamUInt(\"id\") = %d, want %d", got, want)
	}
}

func TestContext_ParamUInt64(t *testing.T) {
	c, _ := createTestContext("GET", "/users/123", nil)
	params := map[string]string{"id": "123"}
	c.setParams(params)

	got, err := c.ParamUInt64("id")
	if err != nil {
		t.Errorf("ctx.ParamUInt64(\"id\") returned an error: %v", err)
	}
	want := uint64(123)
	if got != want {
		t.Errorf("ctx.ParamUInt64(\"id\") = %d, want %d", got, want)
	}
}

func TestContext_ParamFloat32(t *testing.T) {
	c, _ := createTestContext("GET", "/users/123.456", nil)
	params := map[string]string{"id": "123.456"}
	c.setParams(params)

	got, err := c.ParamFloat32("id")
	if err != nil {
		t.Errorf("ctx.ParamFloat32(\"id\") returned an error: %v", err)
	}
	want := float32(123.456)
	if got != want {
		t.Errorf("ctx.ParamFloat32(\"id\") = %f, want %f", got, want)
	}
}

func TestContext_ParamFloat64(t *testing.T) {
	c, _ := createTestContext("GET", "/users/123.456", nil)
	params := map[string]string{"id": "123.456"}
	c.setParams(params)

	got, err := c.ParamFloat64("id")
	if err != nil {
		t.Errorf("ctx.ParamFloat64(\"id\") returned an error: %v", err)
	}
	want := float64(123.456)
	if got != want {
		t.Errorf("ctx.ParamFloat64(\"id\") = %f, want %f", got, want)
	}
}

func TestContext_Query(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=value", nil)

	got := c.Query("key")
	want := "value"
	if got != want {
		t.Errorf("Query() = %q, want %q", got, want)
	}
}

func TestContext_QueryBool(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=true", nil)

	got, err := c.QueryBool("key")
	if err != nil {
		t.Errorf("ctx.QueryBool(\"key\") returned an error: %v", err)
	}
	want := true
	if got != want {
		t.Errorf("QueryBool() = %v, want %v", got, want)
	}
}

func TestContext_QueryBoolWithException(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=notabool", nil)

	_, err := c.QueryBool("key")
	if err == nil {
		t.Error("ctx.QueryBool(\"key\") did not return an error")
	}
}

func TestContext_QueryBoolWithEmptyKey(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=", nil)

	got, err := c.QueryBool("key")
	if err != nil {
		t.Errorf("ctx.QueryBool(\"key\") returned an error: %v", err)
	}
	want := false
	if got != want {
		t.Errorf("QueryBool() = %v, want %v", got, want)
	}
}

func TestContext_QueryInt(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=123", nil)

	got, err := c.QueryInt("key")
	if err != nil {
		t.Errorf("ctx.QueryInt(\"key\") returned an error: %v", err)
	}
	want := 123
	if got != want {
		t.Errorf("QueryInt() = %d, want %d", got, want)
	}
}

func TestContext_QueryIntWithException(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=abc", nil)

	_, err := c.QueryInt("key")
	if err == nil {
		t.Error("ctx.QueryInt(\"key\") did not return an error")
	}
}

func TestContext_QueryIntWithEmptyKey(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=", nil)

	got, err := c.QueryInt("key")
	if err != nil {
		t.Errorf("ctx.QueryInt(\"key\") returned an error: %v", err)
	}
	want := 0
	if got != want {
		t.Errorf("QueryInt() = %d, want %d", got, want)
	}
}

func TestContext_QueryInt8(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=123", nil)

	got, err := c.QueryInt8("key")
	if err != nil {
		t.Errorf("ctx.QueryInt8(\"key\") returned an error: %v", err)
	}
	want := int8(123)
	if got != want {
		t.Errorf("QueryInt8() = %d, want %d", got, want)
	}
}

func TestContext_QueryInt32(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=123", nil)

	got, err := c.QueryInt32("key")
	if err != nil {
		t.Errorf("ctx.QueryInt32(\"key\") returned an error: %v", err)
	}
	want := int32(123)
	if got != want {
		t.Errorf("QueryInt32() = %d, want %d", got, want)
	}
}

func TestContext_QueryInt64(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=123", nil)

	got, err := c.QueryInt64("key")
	if err != nil {
		t.Errorf("ctx.QueryInt64(\"key\") returned an error: %v", err)
	}
	want := int64(123)
	if got != want {
		t.Errorf("QueryInt64() = %d, want %d", got, want)
	}
}

func TestContext_QueryUInt(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=123", nil)

	got, err := c.QueryUInt("key")
	if err != nil {
		t.Errorf("ctx.QueryUInt(\"key\") returned an error: %v", err)
	}
	want := uint(123)
	if got != want {
		t.Errorf("QueryUInt() = %d, want %d", got, want)
	}
}

func TestContext_QueryUInt8(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=123", nil)

	got, err := c.QueryUInt8("key")
	if err != nil {
		t.Errorf("ctx.QueryUInt8(\"key\") returned an error: %v", err)
	}
	want := uint8(123)
	if got != want {
		t.Errorf("QueryUInt8() = %d, want %d", got, want)
	}
}

func TestContext_QueryUInt32(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=123", nil)

	got, err := c.QueryUInt32("key")
	if err != nil {
		t.Errorf("ctx.QueryUInt32(\"key\") returned an error: %v", err)
	}
	want := uint32(123)
	if got != want {
		t.Errorf("QueryUInt32() = %d, want %d", got, want)
	}
}

func TestContext_QueryUInt64(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=123", nil)

	got, err := c.QueryUInt64("key")
	if err != nil {
		t.Errorf("ctx.QueryUInt64(\"key\") returned an error: %v", err)
	}
	want := uint64(123)
	if got != want {
		t.Errorf("QueryUInt64() = %d, want %d", got, want)
	}
}

func TestContext_QueryFloat32(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=3.1415", nil)

	got, err := c.QueryFloat32("key")
	if err != nil {
		t.Errorf("ctx.QueryFloat32(\"key\") returned an error: %v", err)
	}
	want := float32(3.1415)
	if got != want {
		t.Errorf("QueryFloat32() = %f, want %f", got, want)
	}
}

func TestContext_QueryFloat64(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=3.1415", nil)

	got, err := c.QueryFloat64("key")
	if err != nil {
		t.Errorf("ctx.QueryFloat64(\"key\") returned an error: %v", err)
	}
	want := float64(3.1415)
	if got != want {
		t.Errorf("QueryFloat64() = %f, want %f", got, want)
	}
}

func TestContextQueries(t *testing.T) {
	c, _ := createTestContext("GET", "/path?foo=bar&baz=qux", nil)

	queries := c.Queries()
	if len(queries["foo"]) != 1 || queries["foo"][0] != "bar" {
		t.Errorf("got %v, want foo=bar", queries)
	}
	if len(queries["baz"]) != 1 || queries["baz"][0] != "qux" {
		t.Errorf("got %v, want baz=qux", queries)
	}
}

func TestContext_Status(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)

	status := c.Status()
	if status != StatusNotFound {
		t.Errorf("Expected status code %d, but got %d", StatusNotFound, status)
	}

	c.SetStatus(StatusOK)
	status = c.Status()
	if status != StatusOK {
		t.Errorf("Expected status code %d, but got %d", StatusOK, status)
	}
}

func TestHeader(t *testing.T) {
	c, _ := createTestContext("GET", "/path", nil)
	c.ctx.Request.Header.Set("Content-Type", "application/json")

	if got := c.Header("Content-Type"); got != "application/json" {
		t.Errorf("Header() = %q, want %q", got, "application/json")
	}
}

func TestHeaders(t *testing.T) {
	c, _ := createTestContext("GET", "/path", nil)
	c.ctx.Request.Header.Set("Content-Type", "application/json")
	c.ctx.Request.Header.Set("X-Request-ID", "12345")

	headers := c.Headers()
	if len(headers) == 0 {
		t.Error("Expected headers to be non-empty")
	}
}

func TestAddHeader(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)

	c.AddHeader("X-Custom-Header", "value1")
	c.AddHeader("X-Custom-Header", "value2")

	hdr := string(c.ctx.Response.Header.Peek("X-Custom-Header"))
	if hdr == "" {
		t.Errorf("Expected X-Custom-Header to be set, got empty")
	}
}

func TestSetHeader(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)

	c.SetHeader("Content-Type", "application/json")

	if got := string(c.ctx.Response.Header.Peek("Content-Type")); got != "application/json" {
		t.Errorf("Expected Content-Type to be 'application/json', got %s", got)
	}
}

func TestDelHeader(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)

	c.SetHeader("X-Custom", "value")
	c.DelHeader("X-Custom")

	if got := string(c.ctx.Response.Header.Peek("X-Custom")); got != "" {
		t.Errorf("Expected X-Custom to be deleted, got %s", got)
	}
}

func TestCookie(t *testing.T) {
	c, _ := createTestContext("GET", "/path", nil)
	c.ctx.Request.Header.SetCookie("test", "value")

	cookie := c.Cookie("test")
	if cookie != nil {
		if string(cookie.Key()) != "test" {
			t.Errorf("Expected cookie key 'test', got '%s'", string(cookie.Key()))
		}
	}
}

func TestCookieNotFound(t *testing.T) {
	c, _ := createTestContext("GET", "/path", nil)

	cookie := c.Cookie("nonexistent")
	if cookie != nil {
		t.Errorf("Expected nil for nonexistent cookie, got %v", cookie)
	}
}

func TestSetCookie(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)

	c.SetCookie("test", "value")
	c.flush()

	cookie := string(c.ctx.Response.Header.Peek("Set-Cookie"))
	if !strings.Contains(cookie, "test=value") {
		t.Errorf("Expected Set-Cookie to contain 'test=value', got %s", cookie)
	}
}

func TestContextBody(t *testing.T) {
	body := []byte("test body")
	c := &Context{
		res: &response{
			body: body,
		},
	}

	result := c.Body()
	if !bytes.Equal(result, body) {
		t.Errorf("expected body %v, but got %v", body, result)
	}
}

func TestContextSetBody(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)

	body := []byte("test body")
	c.SetBody(body)

	if !bytes.Equal(c.res.body, body) {
		t.Errorf("expected body %v, got %v", body, c.res.body)
	}
}

func TestJSON(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)

	c.JSON(200, map[string]string{"message": "hello world"})
	c.flush()

	if c.Status() != 200 {
		t.Errorf("Expected status code 200 but got %d", c.Status())
	}

	expectedBody := `{"message":"hello world"}`
	if string(c.res.body) != expectedBody {
		t.Errorf("Expected body %s but got %s", expectedBody, string(c.res.body))
	}
}

func TestText(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)

	c.Text(200, "hello world")
	c.flush()

	if c.Status() != 200 {
		t.Errorf("Expected status code 200 but got %d", c.Status())
	}

	expectedBody := "hello world"
	if string(c.res.body) != expectedBody {
		t.Errorf("Expected body %s but got %s", expectedBody, string(c.res.body))
	}
}

func TestContext_XML(t *testing.T) {
	c, _ := createTestContext("GET", "/xml", nil)

	type person struct {
		Name string `xml:"name"`
		Age  int    `xml:"age"`
	}
	obj := &person{
		Name: "John",
		Age:  30,
	}

	c.XML(StatusOK, obj)
	c.flush()

	contentType := string(c.ctx.Response.Header.Peek("Content-Type"))
	if !strings.Contains(contentType, MIMEApplicationXML) {
		t.Errorf("Expected Content-Type header to contain %s, but got %s", MIMEApplicationXML, contentType)
	}

	expectedBody := `<person><name>John</name><age>30</age></person>`
	if string(c.res.body) != expectedBody {
		t.Errorf("Expected response body to be %s, but got %s", expectedBody, string(c.res.body))
	}
}

func TestContext_GetData(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)

	if c.GetData("nonexistent") != nil {
		t.Errorf("expected nil value for nonexistent key, got %v", c.GetData("nonexistent"))
	}

	c.SetData("key", "value")
	if c.GetData("key") != "value" {
		t.Errorf("expected value 'value' for key 'key', got %v", c.GetData("key"))
	}

	c.DelData("key")
	if c.GetData("key") != nil {
		t.Errorf("expected nil value for deleted key 'key', got %v", c.GetData("key"))
	}
}

func TestContext_Redirect(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)

	redirectUrl := "/new"
	c.Redirect(StatusMovedPermanently, redirectUrl)
	c.flush()

	if c.ctx.Response.StatusCode() != StatusMovedPermanently {
		t.Errorf("expected status code %d, got %d", StatusMovedPermanently, c.ctx.Response.StatusCode())
	}

	location := string(c.ctx.Response.Header.Peek("Location"))
	if location == "" {
		t.Error("expected Location header to be set")
	}
}

func TestUserAgent(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)
	ua := "my-user-agent"
	c.ctx.Request.Header.SetUserAgent(ua)

	if userAgent := c.UserAgent(); userAgent != ua {
		t.Errorf("expected user agent %q, got %q", ua, userAgent)
	}
}

func TestReferer(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)
	ref := "https://example.com"
	c.ctx.Request.Header.SetReferer(ref)

	if referer := c.Referer(); referer != ref {
		t.Errorf("expected referer %q, got %q", ref, referer)
	}
}

func TestContext_Success(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)

	testData := map[string]string{"foo": "bar"}
	c.Success(testData)
	c.flush()

	if c.ctx.Response.StatusCode() != StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", c.ctx.Response.StatusCode(), StatusOK)
	}

	expected := `{"code":0,"data":{"foo":"bar"},"message":"ok"}`
	if string(c.res.body) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", string(c.res.body), expected)
	}
}

func TestContextFail(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)

	c.Fail(500, "Internal Server Error")
	c.flush()

	if c.ctx.Response.StatusCode() != StatusOK {
		t.Errorf("expected status code 200, got %d", c.ctx.Response.StatusCode())
	}
	expectedBody := `{"code":500,"message":"Internal Server Error"}`
	if string(c.res.body) != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, string(c.res.body))
	}
}

func TestContext_JSONError(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)

	c.JSONError(StatusBadRequest, "invalid request")
	c.flush()

	if c.ctx.Response.StatusCode() != StatusBadRequest {
		t.Errorf("expected status code %d, got %d", StatusBadRequest, c.ctx.Response.StatusCode())
	}
	expected := `{"code":400,"message":"invalid request"}`
	if string(c.res.body) != expected {
		t.Errorf("expected body %q, got %q", expected, string(c.res.body))
	}
}

func TestContext_IsAjax(t *testing.T) {
	tests := []struct {
		name     string
		header   string
		expected bool
	}{
		{"XMLHttpRequest", "XMLHttpRequest", true},
		{"empty", "", false},
		{"other", "fetch", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := createTestContext("GET", "/test", nil)
			if tt.header != "" {
				c.ctx.Request.Header.Set("X-Requested-With", tt.header)
			}

			if got := c.IsAjax(); got != tt.expected {
				t.Errorf("IsAjax() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestContext_IsWebSocket(t *testing.T) {
	tests := []struct {
		name     string
		header   string
		expected bool
	}{
		{"websocket", "websocket", true},
		{"empty", "", false},
		{"http", "http/1.1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := createTestContext("GET", "/test", nil)
			if tt.header != "" {
				c.ctx.Request.Header.Set("Upgrade", tt.header)
			}

			if got := c.IsWebSocket(); got != tt.expected {
				t.Errorf("IsWebSocket() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestContext_ContentType(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)
	c.ctx.Request.Header.SetContentType("application/json")

	if got := c.ContentType(); got != "application/json" {
		t.Errorf("ContentType() = %v, want %v", got, "application/json")
	}
}

func TestContext_AcceptedLanguages(t *testing.T) {
	tests := []struct {
		name     string
		header   string
		expected []string
	}{
		{"single", "en-US", []string{"en-US"}},
		{"multiple", "en-US, zh-CN, fr", []string{"en-US", "zh-CN", "fr"}},
		{"with quality", "en-US;q=0.9, zh-CN;q=0.8", []string{"en-US", "zh-CN"}},
		{"empty", "", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := createTestContext("GET", "/test", nil)
			if tt.header != "" {
				c.ctx.Request.Header.Set("Accept-Language", tt.header)
			}

			got := c.AcceptedLanguages()
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("AcceptedLanguages() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestContext_File(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	content := []byte("test content")
	if _, err := tmpFile.Write(content); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	c, _ := createTestContext("GET", "/test", nil)

	err = c.File(tmpFile.Name())
	if err != nil {
		t.Errorf("File() returned error: %v", err)
	}
}

func TestContext_FileNotFound(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)

	err := c.File("/nonexistent/path/file.txt")
	if err == nil {
		t.Error("File() expected error for nonexistent file")
	}
}

func TestContext_HTML(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "templates")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tmplPath := tmpDir + "/test.html"
	if err := os.WriteFile(tmplPath, []byte("<html>{{.Name}}</html>"), 0644); err != nil {
		t.Fatal(err)
	}

	app := NewApp()
	app.SetFuncMap(template.FuncMap{})
	app.LoadHTMLGlob(tmpDir + "/*.html")

	app.Get("/test", func(ctx *Context) {
		ctx.HTML(StatusOK, "test.html", map[string]string{"Name": "World"})
	})

	c := newTestCtxForApp(MethodGet, "/test")
	app.serveRequest(c)

	if c.Response.StatusCode() != StatusOK {
		t.Errorf("expected status %d, got %d", StatusOK, c.Response.StatusCode())
	}
	if !strings.Contains(string(c.Response.Body()), "World") {
		t.Errorf("expected body to contain 'World', got %q", string(c.Response.Body()))
	}
}

func TestSkipFlush(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)

	c.SkipFlush()

	c.flush()
}

func TestContext_Cookies(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.Header.SetRequestURI("/test")
	ctx.Request.Header.SetCookie("session", "abc")
	ctx.Request.Header.SetCookie("theme", "dark")

	c := &Context{
		ctx:   ctx,
		index: -1,
		data:  contextData{},
	}
	c.req = newRequest(ctx)

	cookies := c.Cookies()
	if len(cookies) == 0 {
		t.Error("Expected cookies to be non-empty")
	}
}

func TestContext_RemoteAddr(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.Header.SetRequestURI("/test")

	c := &Context{
		ctx:   ctx,
		index: -1,
		data:  contextData{},
	}
	c.req = newRequest(ctx)

	addr := c.RemoteAddr()
	if addr == "" {
		t.Error("Expected non-empty remote address")
	}
}

func TestContext_RemoteAddrWithXRealIP(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.Header.SetRequestURI("/test")
	ctx.Request.Header.Set("X-Real-IP", "1.2.3.4")

	c := &Context{
		ctx:   ctx,
		index: -1,
		data:  contextData{},
	}
	c.req = newRequest(ctx)

	addr := c.RemoteAddr()
	if addr != "1.2.3.4" {
		t.Errorf("Expected X-Real-IP, got %s", addr)
	}
}

func TestContext_RemoteAddrWithXForwardedFor(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.Header.SetRequestURI("/test")
	ctx.Request.Header.Set("X-Forwarded-For", "1.2.3.4, 5.6.7.8")

	c := &Context{
		ctx:   ctx,
		index: -1,
		data:  contextData{},
	}
	c.req = newRequest(ctx)

	addr := c.RemoteAddr()
	if addr != "1.2.3.4" {
		t.Errorf("Expected first IP from X-Forwarded-For, got %s", addr)
	}
}

func TestContext_JSONBodyWithCustomDecoder(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("POST")
	ctx.Request.Header.SetRequestURI("/test")
	ctx.Request.SetBody([]byte(`{"name":"custom"}`))

	app := NewApp()
	called := false
	app.Config.JSONDecoder = func(data []byte, v any) error {
		called = true
		return json.Unmarshal(data, v)
	}

	c := &Context{
		ctx:   ctx,
		index: -1,
		data:  contextData{},
		App:   app,
	}
	c.req = newRequest(ctx)
	c.res = newResponse(ctx)
	c.Method = c.req.method()
	c.Path = c.req.path()

	var result map[string]string
	err := c.JSONBody(&result)
	if err != nil {
		t.Fatalf("JSONBody returned error: %v", err)
	}
	if !called {
		t.Error("Expected custom decoder to be called")
	}
	if result["name"] != "custom" {
		t.Errorf("Expected name 'custom', got '%s'", result["name"])
	}
}

func TestContext_JSONWithCustomEncoder(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.Header.SetRequestURI("/test")

	app := NewApp()
	app.Config.JSONEncoder = func(v any) ([]byte, error) {
		return []byte(`{"custom":true}`), nil
	}

	c := &Context{
		ctx:   ctx,
		index: -1,
		data:  contextData{},
		App:   app,
	}
	c.req = newRequest(ctx)
	c.res = newResponse(ctx)
	c.Method = c.req.method()
	c.Path = c.req.path()

	c.JSON(StatusOK, map[string]bool{"test": true})
	c.flush()

	if string(c.res.body) != `{"custom":true}` {
		t.Errorf("Expected custom encoded body, got %s", string(c.res.body))
	}
}

func TestContext_HTMLWithTemplateError(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "templates_error")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	if err := os.WriteFile(tmpDir+"/bad.html", []byte("<html>{{.Name}}</html>"), 0644); err != nil {
		t.Fatal(err)
	}

	app := NewApp()
	app.SetFuncMap(template.FuncMap{})
	app.LoadHTMLGlob(tmpDir + "/*.html")

	app.Get("/bad", func(ctx *Context) {
		ctx.HTML(StatusOK, "bad.html", nil)
	})

	c := newTestCtxForApp(MethodGet, "/bad")
	app.serveRequest(c)

	if c.Response.StatusCode() != StatusOK {
		t.Errorf("Expected status %d, got %d", StatusOK, c.Response.StatusCode())
	}
}

func TestContext_XMLWithMarshalError(t *testing.T) {
	c, _ := createTestContext("GET", "/xml", nil)

	type BadXML struct {
		Ch chan int `xml:"ch"`
	}
	c.XML(StatusOK, &BadXML{Ch: make(chan int)})
	c.flush()

	if len(c.res.body) != 0 {
		t.Errorf("Expected empty body for XML marshal error, got %s", string(c.res.body))
	}
}

func TestContext_QueryInt8Error(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=invalid", nil)

	_, err := c.QueryInt8("key")
	if err == nil {
		t.Error("Expected error for invalid int8")
	}
}

func TestContext_QueryInt32Error(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=invalid", nil)

	_, err := c.QueryInt32("key")
	if err == nil {
		t.Error("Expected error for invalid int32")
	}
}

func TestContext_QueryInt64Error(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=invalid", nil)

	_, err := c.QueryInt64("key")
	if err == nil {
		t.Error("Expected error for invalid int64")
	}
}

func TestContext_QueryUIntError(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=invalid", nil)

	_, err := c.QueryUInt("key")
	if err == nil {
		t.Error("Expected error for invalid uint")
	}
}

func TestContext_QueryUInt8Error(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=invalid", nil)

	_, err := c.QueryUInt8("key")
	if err == nil {
		t.Error("Expected error for invalid uint8")
	}
}

func TestContext_QueryUInt32Error(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=invalid", nil)

	_, err := c.QueryUInt32("key")
	if err == nil {
		t.Error("Expected error for invalid uint32")
	}
}

func TestContext_QueryUInt64Error(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=invalid", nil)

	_, err := c.QueryUInt64("key")
	if err == nil {
		t.Error("Expected error for invalid uint64")
	}
}

func TestContext_QueryFloat32Error(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=invalid", nil)

	_, err := c.QueryFloat32("key")
	if err == nil {
		t.Error("Expected error for invalid float32")
	}
}

func TestContext_QueryFloat64Error(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=invalid", nil)

	_, err := c.QueryFloat64("key")
	if err == nil {
		t.Error("Expected error for invalid float64")
	}
}

func TestContext_QueryBoolEmpty(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=", nil)

	got, err := c.QueryBool("key")
	if err != nil {
		t.Errorf("QueryBool returned error: %v", err)
	}
	if got != false {
		t.Errorf("Expected false for empty key, got %v", got)
	}
}

func TestContext_QueryIntEmpty(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=", nil)

	got, err := c.QueryInt("key")
	if err != nil {
		t.Errorf("QueryInt returned error: %v", err)
	}
	if got != 0 {
		t.Errorf("Expected 0 for empty key, got %d", got)
	}
}

func TestContext_QueryInt8Empty(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=", nil)

	got, err := c.QueryInt8("key")
	if err != nil {
		t.Errorf("QueryInt8 returned error: %v", err)
	}
	if got != 0 {
		t.Errorf("Expected 0 for empty key, got %d", got)
	}
}

func TestContext_QueryInt32Empty(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=", nil)

	got, err := c.QueryInt32("key")
	if err != nil {
		t.Errorf("QueryInt32 returned error: %v", err)
	}
	if got != 0 {
		t.Errorf("Expected 0 for empty key, got %d", got)
	}
}

func TestContext_QueryInt64Empty(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=", nil)

	got, err := c.QueryInt64("key")
	if err != nil {
		t.Errorf("QueryInt64 returned error: %v", err)
	}
	if got != 0 {
		t.Errorf("Expected 0 for empty key, got %d", got)
	}
}

func TestContext_QueryUIntEmpty(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=", nil)

	got, err := c.QueryUInt("key")
	if err != nil {
		t.Errorf("QueryUInt returned error: %v", err)
	}
	if got != 0 {
		t.Errorf("Expected 0 for empty key, got %d", got)
	}
}

func TestContext_QueryUInt8Empty(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=", nil)

	got, err := c.QueryUInt8("key")
	if err != nil {
		t.Errorf("QueryUInt8 returned error: %v", err)
	}
	if got != 0 {
		t.Errorf("Expected 0 for empty key, got %d", got)
	}
}

func TestContext_QueryUInt32Empty(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=", nil)

	got, err := c.QueryUInt32("key")
	if err != nil {
		t.Errorf("QueryUInt32 returned error: %v", err)
	}
	if got != 0 {
		t.Errorf("Expected 0 for empty key, got %d", got)
	}
}

func TestContext_QueryUInt64Empty(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=", nil)

	got, err := c.QueryUInt64("key")
	if err != nil {
		t.Errorf("QueryUInt64 returned error: %v", err)
	}
	if got != 0 {
		t.Errorf("Expected 0 for empty key, got %d", got)
	}
}

func TestContext_QueryFloat32Empty(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=", nil)

	got, err := c.QueryFloat32("key")
	if err != nil {
		t.Errorf("QueryFloat32 returned error: %v", err)
	}
	if got != 0 {
		t.Errorf("Expected 0 for empty key, got %f", got)
	}
}

func TestContext_QueryFloat64Empty(t *testing.T) {
	c, _ := createTestContext("GET", "/path?key=", nil)

	got, err := c.QueryFloat64("key")
	if err != nil {
		t.Errorf("QueryFloat64 returned error: %v", err)
	}
	if got != 0 {
		t.Errorf("Expected 0 for empty key, got %f", got)
	}
}

func TestContext_setParams(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)

	params := map[string]string{"id": "123", "name": "test"}
	c.setParams(params)

	if c.Param("id") != "123" {
		t.Errorf("Expected param id=123, got %s", c.Param("id"))
	}
	if c.Param("name") != "test" {
		t.Errorf("Expected param name=test, got %s", c.Param("name"))
	}
}

func TestContext_setHandlers(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)

	handlers := []HandlerFunc{
		func(c *Context) {},
		func(c *Context) {},
		func(c *Context) {},
	}
	c.setHandlers(handlers)

	if len(c.handlers) != 3 {
		t.Errorf("Expected 3 handlers, got %d", len(c.handlers))
	}
}

func TestContext_setApp(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)

	app := NewApp()
	c.setApp(app)

	if c.App != app {
		t.Error("App not set correctly")
	}
}

func TestContext_NextNoMoreHandlers(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)
	c.index = 0

	c.Next()

	if c.index != 1 {
		t.Errorf("Expected index to be 1, got %d", c.index)
	}
}

func TestContext_StringBodyEmpty(t *testing.T) {
	c, _ := createTestContext("POST", "/test", nil)

	body := c.StringBody()
	if body != "" {
		t.Errorf("Expected empty body, got %s", body)
	}
}

func TestContext_JSONBodyEmpty(t *testing.T) {
	c, _ := createTestContext("POST", "/test", nil)

	var result map[string]string
	err := c.JSONBody(&result)
	if err == nil {
		t.Error("Expected error for empty body")
	}
}

func TestContext_Body(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)

	body := []byte("response body")
	c.SetBody(body)

	if !bytes.Equal(c.Body(), body) {
		t.Errorf("Expected body %v, got %v", body, c.Body())
	}
}

func TestContext_DelHeader(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)

	c.SetHeader("X-Custom", "value")
	c.DelHeader("X-Custom")

	if got := string(c.ctx.Response.Header.Peek("X-Custom")); got != "" {
		t.Errorf("Expected header to be deleted, got %s", got)
	}
}

func TestContext_SetCookie(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)

	c.SetCookie("session", "abc123")
	c.flush()

	cookie := string(c.ctx.Response.Header.Peek("Set-Cookie"))
	if !strings.Contains(cookie, "session=abc123") {
		t.Errorf("Expected cookie to contain 'session=abc123', got %s", cookie)
	}
}

func TestContext_AddHeader(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)

	c.AddHeader("X-Custom", "value1")

	hdr := string(c.ctx.Response.Header.Peek("X-Custom"))
	if hdr == "" {
		t.Error("Expected header to be set")
	}
}

func TestContext_SetHeader(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)

	c.SetHeader("Content-Type", "text/plain")

	if got := string(c.ctx.Response.Header.Peek("Content-Type")); got != "text/plain" {
		t.Errorf("Expected 'text/plain', got %s", got)
	}
}

func TestContext_Header(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)
	c.ctx.Request.Header.Set("X-Custom", "header-value")

	if got := c.Header("X-Custom"); got != "header-value" {
		t.Errorf("Expected 'header-value', got %s", got)
	}
}

func TestContext_Headers(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)
	c.ctx.Request.Header.Set("Content-Type", "application/json")

	headers := c.Headers()
	if len(headers) == 0 {
		t.Error("Expected headers to be non-empty")
	}
}

func TestContext_Text(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)

	c.Text(StatusOK, "hello")
	c.flush()

	if c.Status() != StatusOK {
		t.Errorf("Expected status %d, got %d", StatusOK, c.Status())
	}
	if string(c.res.body) != "hello" {
		t.Errorf("Expected body 'hello', got %s", string(c.res.body))
	}
}

func TestContext_SetStatus(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)

	c.SetStatus(StatusCreated)

	if c.Status() != StatusCreated {
		t.Errorf("Expected status %d, got %d", StatusCreated, c.Status())
	}
}

func TestContext_DefaultStatus(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)

	if c.Status() != StatusNotFound {
		t.Errorf("Expected default status %d, got %d", StatusNotFound, c.Status())
	}
}

func TestContext_CookieNotFound(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)

	cookie := c.Cookie("nonexistent")
	if cookie != nil {
		t.Errorf("Expected nil for nonexistent cookie, got %v", cookie)
	}
}

func TestContext_DelData(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)

	c.SetData("key", "value")
	c.DelData("key")

	if c.GetData("key") != nil {
		t.Error("Expected nil after DelData")
	}
}

func TestContext_Fail(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)

	c.Fail(500, "internal error")
	c.flush()

	if c.Status() != StatusOK {
		t.Errorf("Expected status %d, got %d", StatusOK, c.Status())
	}
	if !strings.Contains(string(c.res.body), "internal error") {
		t.Errorf("Expected body to contain 'internal error', got %s", string(c.res.body))
	}
}

func TestContext_JSONErrorStatus(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)

	c.JSONError(StatusBadRequest, "bad request")
	c.flush()

	if c.Status() != StatusBadRequest {
		t.Errorf("Expected status %d, got %d", StatusBadRequest, c.Status())
	}
}

func TestContext_IsAjaxTrue(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)
	c.ctx.Request.Header.Set("X-Requested-With", "XMLHttpRequest")

	if !c.IsAjax() {
		t.Error("Expected IsAjax to return true")
	}
}

func TestContext_IsNotAjax(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)

	if c.IsAjax() {
		t.Error("Expected IsAjax to return false")
	}
}

func TestContext_IsWebSocketTrue(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)
	c.ctx.Request.Header.Set("Upgrade", "websocket")

	if !c.IsWebSocket() {
		t.Error("Expected IsWebSocket to return true")
	}
}

func TestContext_IsNotWebSocket(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)

	if c.IsWebSocket() {
		t.Error("Expected IsWebSocket to return false")
	}
}

func TestContext_ContentTypeSet(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)
	c.ctx.Request.Header.SetContentType("application/json")

	if c.ContentType() != "application/json" {
		t.Errorf("Expected 'application/json', got %s", c.ContentType())
	}
}

func TestContext_AcceptedLanguagesEmpty(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)

	langs := c.AcceptedLanguages()
	if langs != nil {
		t.Errorf("Expected nil for empty Accept-Language, got %v", langs)
	}
}

func TestContext_AcceptedLanguagesMultiple(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)
	c.ctx.Request.Header.Set("Accept-Language", "en-US, zh-CN, fr")

	langs := c.AcceptedLanguages()
	if len(langs) != 3 {
		t.Errorf("Expected 3 languages, got %d", len(langs))
	}
}

func TestContext_AcceptedLanguagesWithQuality(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)
	c.ctx.Request.Header.Set("Accept-Language", "en-US;q=0.9, zh-CN;q=0.8")

	langs := c.AcceptedLanguages()
	if len(langs) != 2 {
		t.Errorf("Expected 2 languages, got %d", len(langs))
	}
}

func TestContext_FileExists(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	content := []byte("test content")
	if _, err := tmpFile.Write(content); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	c, _ := createTestContext("GET", "/test", nil)

	err = c.File(tmpFile.Name())
	if err != nil {
		t.Errorf("File() returned error: %v", err)
	}
}
