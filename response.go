package lightning

import (
	"net/http"
	"os"
	"path/filepath"
)

type response struct {
	req        *http.Request
	writer     http.ResponseWriter
	statusCode int
	cookies    cookiesMap
	body       []byte
	redirectTo string
	filePath   string
}

func newResponse(req *http.Request, writer http.ResponseWriter) *response {
	return &response{
		req:        req,
		writer:     writer,
		statusCode: http.StatusNotFound,
		cookies:    cookiesMap{},
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
	r.writer.Header().Add(key, value)
}

func (r *response) setHeader(key string, value string) {
	r.writer.Header().Set(key, value)
}

func (r *response) delHeader(key string) {
	r.writer.Header().Del(key)
}

func (r *response) sendFile() {
	base := filepath.Base(r.filePath)
	r.writer.Header().Set(HeaderContentDisposition, "attachment; filename="+base)
	http.ServeFile(r.writer, r.req, r.filePath)
}

func (r *response) flush() {
	for _, v := range r.cookies {
		http.SetCookie(r.writer, v)
	}

	if len(r.filePath) > 0 {
		r.sendFile()
	} else if len(r.redirectTo) > 0 {
		http.Redirect(r.writer, r.req, r.redirectTo, r.statusCode)
	} else {
		r.writer.WriteHeader(r.statusCode)
		r.writer.Write(r.body)
	}
}
