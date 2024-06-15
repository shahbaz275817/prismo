package wrappers

import "net/http"

type HTTPStatusWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewHTTPStatusWriter(w http.ResponseWriter) *HTTPStatusWriter {
	return &HTTPStatusWriter{
		ResponseWriter: w,
		statusCode:     0,
	}
}

func (h *HTTPStatusWriter) WriteHeader(status int) {
	//As per the http implementation, only write header if it is not called before
	if h.writeHeaderNotCalled() {
		h.ResponseWriter.WriteHeader(status)
		h.statusCode = status
	}
}

// If WriteHeader is not called explicitly, the first call to Write
// will trigger an implicit WriteHeader(http.StatusOK).
// Thus explicit calls to WriteHeader are mainly used to
// send error codes.
func (h *HTTPStatusWriter) Write(b []byte) (int, error) {
	if h.writeHeaderNotCalled() {
		h.WriteHeader(http.StatusOK)
	}
	return h.ResponseWriter.Write(b)
}

func (h *HTTPStatusWriter) Header() http.Header {
	return h.ResponseWriter.Header()
}

func (h *HTTPStatusWriter) writeHeaderNotCalled() bool {
	return h.statusCode == 0
}
