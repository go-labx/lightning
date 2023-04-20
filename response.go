package lightning

import (
	"net/http"
	"os"
	"path/filepath"
)

// response Declaring the response structure that will be used to hold HTTP response body.
type response struct {
	originReq   *http.Request       // A pointer to an HTTP request.
	originRes   http.ResponseWriter // An HTTP response writer.
	statusCode  int                 // The status code of the HTTP response (e.g. 200, 404, 500, etc.).
	cookies     cookiesMap          // An array of cookies to be sent with the HTTP response.
	body        []byte              // The response body to be sent.
	redirectUrl string              // The URL to redirect to.
	fileUrl     string              // The file to send.
}

// newResponse A constructor function for the response structure.
func newResponse(req *http.Request, res http.ResponseWriter) *response {
	return &response{
		originReq:   req,
		originRes:   res,
		statusCode:  http.StatusNotFound,
		cookies:     cookiesMap{},
		body:        nil,
		redirectUrl: "",
	}
}

// setStatus sets the status code of the HTTP response.
func (r *response) setStatus(code int) {
	r.statusCode = code
}

// setBody sets the response body to be sent.
func (r *response) setBody(body []byte) {
	r.body = body
}

// redirect sets a redirect URL.
func (r *response) redirect(code int, url string) {
	r.statusCode = code
	r.redirectUrl = url
}

// file serves a file.
func (r *response) file(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	if _, err = os.Stat(absPath); err != nil {
		return err
	}

	r.fileUrl = absPath
	return nil
}

// addHeader adds a new header key-value pair to the response.
func (r *response) addHeader(key, value string) {
	r.originRes.Header().Add(key, value)
}

// setHeader sets the value of a given header in the response.
func (r *response) setHeader(key string, value string) {
	r.originRes.Header().Set(key, value)
}

// delHeader deletes a given header from the response.
func (r *response) delHeader(key string) {
	r.originRes.Header().Del(key)
}

// sendFile sends the file as an attachment.
func (r *response) sendFile() {
	base := filepath.Base(r.fileUrl)
	r.originRes.Header().Set(HeaderContentDisposition, "attachment; filename="+base)
	http.ServeFile(r.originRes, r.originReq, r.fileUrl)
}

// flush sends the HTTP response.
func (r *response) flush() {
	for _, v := range r.cookies {
		http.SetCookie(r.originRes, v)
	}

	if len(r.fileUrl) > 0 {
		r.sendFile()
	} else if len(r.redirectUrl) > 0 {
		http.Redirect(r.originRes, r.originReq, r.redirectUrl, r.statusCode)
	} else {
		r.originRes.WriteHeader(r.statusCode)
		_, err := r.originRes.Write(r.body)
		if err != nil {
			return
		}
	}
}
