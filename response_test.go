package lightning

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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

	if resp.body != nil {
		t.Errorf("Expected data to be nil, but got %v", resp.body)
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

func TestSetBody(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()

	r := newResponse(req, res)
	body := []byte("test data")
	r.setBody(body)

	if !bytes.Equal(r.body, body) {
		t.Errorf("expected body to be %v, but got %v", body, r.body)
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

func TestResponse_File(t *testing.T) {
	// create a temporary file
	file, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	// set the file path using the file method
	resp := newResponse(nil, nil)
	err = resp.file(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	// check if the fileUrl field is set to the correct absolute path
	absPath, err := filepath.Abs(file.Name())
	if err != nil {
		t.Fatal(err)
	}
	if resp.fileUrl != absPath {
		t.Errorf("fileUrl field is %s, expected %s", resp.fileUrl, absPath)
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

func TestSendFile(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/file.txt", nil)
	res := httptest.NewRecorder()

	// Create a temporary file to serve
	file, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	_, err = file.WriteString("test content")
	if err != nil {
		t.Fatal(err)
	}
	err = file.Close()
	if err != nil {
		t.Fatal(err)
	}

	resp := newResponse(req, res)
	err = resp.file(file.Name())
	if err != nil {
		t.Fatal(err)
	}
	resp.sendFile()

	// Check that the Content-Disposition header was set correctly
	expectedHeader := "attachment; filename=" + filepath.Base(file.Name())
	if res.Header().Get(HeaderContentDisposition) != expectedHeader {
		t.Errorf("Expected Content-Disposition header %q, got %q", expectedHeader, res.Header().Get(HeaderContentDisposition))
	}

	// Check that the file was served
	expectedBody := "test content"
	if res.Body.String() != expectedBody {
		t.Errorf("Expected response body %q, got %q", expectedBody, res.Body.String())
	}
}
