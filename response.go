package lightning

import (
	"os"
	"path/filepath"

	"github.com/valyala/fasthttp"
)

type response struct {
	ctx        *fasthttp.RequestCtx
	statusCode int
	body       []byte
	redirectTo string
	filePath   string
	cookies    cookiesMap
}

func newResponse(ctx *fasthttp.RequestCtx) *response {
	return &response{
		ctx:        ctx,
		statusCode: StatusNotFound,
		cookies:    make(cookiesMap),
	}
}

func (r *response) setStatus(code int) {
	r.statusCode = code
}

func (r *response) setBody(body []byte) {
	r.body = body
}

func (r *response) redirect(code int, url string) {
	r.statusCode = code
	r.redirectTo = url
}

func (r *response) file(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	if _, err = os.Stat(absPath); err != nil {
		return err
	}

	r.filePath = absPath
	return nil
}

func (r *response) addHeader(key, value string) {
	r.ctx.Response.Header.Add(key, value)
}

func (r *response) setHeader(key string, value string) {
	r.ctx.Response.Header.Set(key, value)
}

func (r *response) delHeader(key string) {
	r.ctx.Response.Header.Del(key)
}

func (r *response) sendFile() {
	base := filepath.Base(r.filePath)
	r.ctx.Response.Header.Set(HeaderContentDisposition, "attachment; filename="+base)
	r.ctx.SendFile(r.filePath)
}

func (r *response) flush() {
	for name, value := range r.cookies {
		var c fasthttp.Cookie
		c.SetKey(name)
		c.SetValue(value)
		r.ctx.Response.Header.SetCookie(&c)
	}

	if len(r.filePath) > 0 {
		r.sendFile()
	} else if len(r.redirectTo) > 0 {
		r.ctx.Redirect(r.redirectTo, r.statusCode)
	} else {
		r.ctx.Response.SetStatusCode(r.statusCode)
		r.ctx.Response.SetBody(r.body)
	}
}
