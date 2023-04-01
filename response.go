package lightning

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"os"
	"path/filepath"
)

// response Declaring the response structure that will be used to hold HTTP response data.
type response struct {
	req         *http.Request       // A pointer to an HTTP request.
	res         http.ResponseWriter // An HTTP response writer.
	status      int                 // The status code of the HTTP response (e.g. 200, 404, 500, etc.).
	cookies     Cookies             // An array of cookies to be sent with the HTTP response.
	data        []byte              // The response data to be sent.
	redirectUrl string              // The URL to redirect to.
	fileUrl     string              // The file to send.
}

// newResponse A constructor function for the response structure.
func newResponse(req *http.Request, res http.ResponseWriter) *response {
	return &response{
		req:         req,
		res:         res,
		status:      http.StatusNotFound,
		cookies:     Cookies{},
		data:        nil,
		redirectUrl: "",
	}
}

// setStatus sets the status code of the HTTP response.
func (r *response) setStatus(code int) {
	r.status = code
}

// json marshals a json object and sets the appropriate headers.
func (r *response) json(obj interface{}) error {
	encodeData, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	r.res.Header().Set("Content-Type", "application/json")

	r.raw(encodeData)
	return nil
}

// xml marshals a xml object and sets the appropriate headers.
func (r *response) xml(obj interface{}) error {
	encodeData, err := xml.Marshal(obj)
	if err != nil {
		return err
	}
	r.res.Header().Set("Content-Type", "application/xml")
	r.raw(encodeData)
	return nil
}

// text sets plain text as the response data.
func (r *response) text(text string) {
	r.raw([]byte(text))
}

// raw sets the response data directly.
func (r *response) raw(data []byte) {
	r.data = data
}

// redirect sets a redirect URL.
func (r *response) redirect(code int, url string) {
	r.status = code
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
	r.res.Header().Add(key, value)
}

// setHeader sets the value of a given header in the response.
func (r *response) setHeader(key string, value string) {
	r.res.Header().Set(key, value)
}

// delHeader deletes a given header from the response.
func (r *response) delHeader(key string) {
	r.res.Header().Del(key)
}

// sendFile sends the file as an attachment.
func (r *response) sendFile() {
	base := filepath.Base(r.fileUrl)
	r.res.Header().Set("Content-Disposition", "attachment; filename="+base)
	http.ServeFile(r.res, r.req, r.fileUrl)
}

// flush sends the HTTP response.
func (r *response) flush() {
	for _, v := range r.cookies {
		http.SetCookie(r.res, v)
	}

	if len(r.fileUrl) > 0 {
		r.sendFile()
	} else if len(r.redirectUrl) > 0 {
		http.Redirect(r.res, r.req, r.redirectUrl, r.status)
	} else {
		r.res.WriteHeader(r.status)
		_, err := r.res.Write(r.data)
		if err != nil {
			return
		}
	}
}
