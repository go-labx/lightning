package lightning

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecovery(t *testing.T) {
	app := NewApp()

	app.Use(Recovery())

	app.Get("/", func(ctx *Context) {
		panic("test panic")
	})

	req := httptest.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	app.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	expected := "Internal Server Error"
	if body := res.Body.String(); body != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			body, expected)
	}
}
