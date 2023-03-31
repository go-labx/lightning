package lightning

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"os"
	"path/filepath"
)

// Response Declaring the Response structure that will be used to hold HTTP response data.
type Response struct {
	req      *http.Request       // A pointer to an HTTP request.
	res      http.ResponseWriter // An HTTP response writer.
	status   int                 // The status code of the HTTP response (e.g. 200, 404, 500, etc.).
	Cookies  Cookies             // An array of cookies to be sent with the HTTP response.
	data     []byte              // The response data to be sent.
	redirect string              // The URL to redirect to.
	file     string              // The file to send.
}

// NewResponse A constructor function for the Response structure.
func NewResponse(req *http.Request, res http.ResponseWriter) *Response {
	response := &Response{
		req:      req,
		res:      res,
		status:   http.StatusNotFound,
		Cookies:  Cookies{},
		data:     nil,
		redirect: "",
	}
	return response
}

// SetStatus sets the status code of the HTTP response.
func (r *Response) SetStatus(code int) {
	r.status = code
}

// JSON marshals a JSON object and sets the appropriate headers.
func (r *Response) JSON(obj interface{}) error {
	encodeData, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	r.res.Header().Set("Content-Type", "application/json")

	r.Raw(encodeData)
	return nil
}

// XML marshals an XML object and sets the appropriate headers.
func (r *Response) XML(obj interface{}) error {
	encodeData, err := xml.Marshal(obj)
	if err != nil {
		return err
	}
	r.res.Header().Set("Content-Type", "application/xml")
	r.Raw(encodeData)
	return nil
}

// Text sets plain text as the response data.
func (r *Response) Text(text string) {
	r.Raw([]byte(text))
}

// Raw sets the response data directly.
func (r *Response) Raw(data []byte) {
	r.data = data
}

// Redirect sets a redirect URL.
func (r *Response) Redirect(code int, url string) {
	r.status = code
	r.redirect = url
}

// File serves a file.
func (r *Response) File(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	if _, err = os.Stat(absPath); err != nil {
		return err
	}

	r.file = absPath
	return nil
}

// AddHeader adds a new header key-value pair to the response.
func (r *Response) AddHeader(key, value string) {
	r.res.Header().Add(key, value)
}

// SetHeader sets the value of a given header in the response.
func (r *Response) SetHeader(key string, value string) {
	r.res.Header().Set(key, value)
}

// DelHeader deletes a given header from the response.
func (r *Response) DelHeader(key string) {
	r.res.Header().Del(key)
}

// sendFile sends the file as an attachment.
func (r *Response) sendFile() {
	base := filepath.Base(r.file)
	r.res.Header().Set("Content-Disposition", "attachment; filename="+base)
	http.ServeFile(r.res, r.req, r.file)
}

// flush sends the HTTP response.
func (r *Response) flush() {
	for _, v := range r.Cookies {
		http.SetCookie(r.res, v)
	}

	if len(r.file) > 0 {
		r.sendFile()
	} else if len(r.redirect) > 0 {
		http.Redirect(r.res, r.req, r.redirect, r.status)
	} else {
		r.res.WriteHeader(r.status)
		_, err := r.res.Write(r.data)
		if err != nil {
			return
		}
	}
}
}
