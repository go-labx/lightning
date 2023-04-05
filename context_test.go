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

func TestFlushResponse(t *testing.T) {
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

	// Call the flushResponse function
	ctx.flushResponse()

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
	// Create a new request with a JSON body
	reqBody := []byte(`{"name": "John", "age": 30}`)
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a new context with the request and response writer
	w := httptest.NewRecorder()
	ctx, err := NewContext(w, req)
	if err != nil {
		t.Fatal(err)
	}

	// Parse the JSON body into a struct
	var user struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	err = ctx.JSONBody(&user)
	if err != nil {
		t.Fatal(err)
	}

	// Assert that the parsed JSON object matches the expected output
	expectedUser := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "John",
		Age:  30,
	}
	if !reflect.DeepEqual(user, expectedUser) {
		t.Errorf("got %v, want %v", user, expectedUser)
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
	ctx.flushResponse()

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
	ctx.flushResponse()

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
	ctx.flushResponse()

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
	ctx.flushResponse()

	if ctx.Status() != 200 {
		t.Errorf("Expected status code %d but got %d", 200, ctx.Status())
	}

	expectedBody := "hello world"
	if string(ctx.res.body) != expectedBody {
		t.Errorf("Expected body %s but got %s", expectedBody, string(ctx.res.body))
	}
}

//func TestFile(t *testing.T) {
//	req, err := http.NewRequest("GET", "/test", nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//	rr := httptest.NewRecorder()
//	ctx, err := NewContext(rr, req)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	ctx.File("test.txt")
//	if status := rr.Code; status != http.StatusOK {
//		t.Errorf("handler returned wrong status code: got %v want %v",
//			status, http.StatusOK)
//	}
//}

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
	ctx.flushResponse()

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
	ctx.flushResponse()

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	expected := `{"code":0,"data":{"foo":"bar"},"msg":"ok"}`
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
	ctx.flushResponse()

	// Check that the response status code and body are correct
	if w.Code != 200 {
		t.Errorf("expected status code 200, got %d", w.Code)
	}
	expectedBody := `{"code":500,"msg":"Internal Server Error"}`
	if w.Body.String() != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, w.Body.String())
	}
}
