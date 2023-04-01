package lightning

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewResponse(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	resp := newResponse(req, res)

	if resp.originReq != req {
		t.Errorf("Expected originReq to be %v, but got %v", req, resp.originReq)
	}

	if resp.originRes != res {
		t.Errorf("Expected originRes to be %v, but got %v", res, resp.originRes)
	}

	if resp.statusCode != http.StatusNotFound {
		t.Errorf("Expected statusCode to be %v, but got %v", http.StatusNotFound, resp.statusCode)
	}

	if len(resp.cookies) != 0 {
		t.Errorf("Expected cookies to be empty, but got %v", resp.cookies)
	}

	if resp.data != nil {
		t.Errorf("Expected data to be nil, but got %v", resp.data)
	}

	if resp.redirectUrl != "" {
		t.Errorf("Expected redirectUrl to be empty, but got %v", resp.redirectUrl)
	}

	if resp.fileUrl != "" {
		t.Errorf("Expected fileUrl to be empty, but got %v", resp.fileUrl)
	}
}

func TestSetStatus(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()
	r := newResponse(req, res)

	r.setStatus(http.StatusOK)

	if r.statusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, r.statusCode)
	}
}

func TestResponseJson(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()

	r := newResponse(req, res)
	obj := map[string]string{"foo": "bar"}

	r.setStatus(http.StatusOK)
	if err := r.json(obj); err != nil {
		t.Fatal(err)
	}
	r.flush()

	if res.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, res.Code)
	}

	expectedContentType := "application/json"
	if res.Header().Get("Content-Type") != expectedContentType {
		t.Errorf("expected content type %q, got %q", expectedContentType, res.Header().Get("Content-Type"))
	}

	expectedBody := `{"foo":"bar"}`
	if res.Body.String() != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, res.Body.String())
	}
}

func TestResponse_Text(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()

	r := newResponse(req, res)
	text := "hello world"

	r.setStatus(http.StatusOK)
	r.text(text)
	r.flush()

	if res.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, res.Code)
	}

	expectedContentType := "text/plain"
	fmt.Println(res.Header().Get("Content-Type"))
	if res.Header().Get("Content-Type") != expectedContentType {
		t.Errorf("expected content type %q, got %q", expectedContentType, res.Header().Get("Content-Type"))
	}

	expectedBody := text
	if res.Body.String() != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, res.Body.String())
	}
}

func TestResponseRaw(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()

	resp := newResponse(req, res)
	data := []byte("test data")
	resp.raw(data)
	resp.flush()

	if !bytes.Equal(res.Body.Bytes(), data) {
		t.Errorf("expected response body to be %v, but got %v", data, res.Body.Bytes())
	}
}

func TestResponse_Redirect(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	resp := newResponse(req, w)

	resp.redirect(http.StatusFound, "http://example.com/bar")

	resp.flush()

	if w.Code != http.StatusFound {
		t.Errorf("expected status code %d, got %d", http.StatusFound, w.Code)
	}

	if w.Header().Get("Location") != "http://example.com/bar" {
		t.Errorf("expected Location header %q, got %q", "http://example.com/bar", w.Header().Get("Location"))
	}
}

func TestResponse_AddHeader(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	resp := newResponse(req, res)

	resp.addHeader("X-Test-Header", "test-value")
	resp.flush()

	if res.Header().Get("X-Test-Header") != "test-value" {
		t.Errorf("Expected header X-Test-Header to be set to test-value, but got %s", res.Header().Get("X-Test-Header"))
	}
}

func TestResponse_SetHeader(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	resp := newResponse(req, res)

	resp.setHeader("X-Test-Header", "test-value")
	resp.flush()

	if res.Header().Get("X-Test-Header") != "test-value" {
		t.Errorf("Expected header X-Test-Header to be set to test-value, but got %s", res.Header().Get("X-Test-Header"))
	}
}

func TestResponse_DelHeader(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	resp := newResponse(req, res)

	resp.addHeader("X-Test-Header", "test-value")
	resp.delHeader("X-Test-Header")
	resp.flush()

	if res.Header().Get("X-Test-Header") != "" {
		t.Errorf("Expected header X-Test-Header to be set to test-value, but got %s", res.Header().Get("X-Test-Header"))
	}
}
