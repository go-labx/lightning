package lightning

import (
	"testing"

	"github.com/valyala/fasthttp"
)

func TestRecovery(t *testing.T) {
	app := NewApp()

	app.Use(Recovery(func(ctx *Context) {
		ctx.Text(StatusInternalServerError, "Internal Server Error")
	}))

	app.Get("/", func(ctx *Context) {
		panic("test panic")
	})

	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(MethodGet)
	ctx.Request.Header.SetRequestURI("/")

	app.serveRequest(ctx)

	if ctx.Response.StatusCode() != StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			ctx.Response.StatusCode(), StatusInternalServerError)
	}

	expected := "Internal Server Error"
	if body := string(ctx.Response.Body()); body != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			body, expected)
	}
}
