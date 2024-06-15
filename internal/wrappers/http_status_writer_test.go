package wrappers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpStatusWriter(t *testing.T) {
	tests := []struct {
		name       string
		h          func(w http.ResponseWriter, r *http.Request)
		statusCode int
		content    string
	}{
		{
			name: "200 Header with only WriteHeader Called",
			h: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			},
			statusCode: 200,
			content:    "",
		},
		{
			name: "200 Header with only Write Called",
			h: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("hi"))
			},
			statusCode: 200,
			content:    "hi",
		},
		{
			name: "200 Header with WriteHeader Called multiple times",
			h: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				w.WriteHeader(400)
				w.WriteHeader(500)
			},
			statusCode: 200,
			content:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := http.NewRequest("GET", "http://gosendtest.com", nil)
			h := http.HandlerFunc(tt.h)
			rr := httptest.NewRecorder()
			rec := NewHTTPStatusWriter(rr)
			h.ServeHTTP(rec, r)
			assert.Equal(t, tt.statusCode, rr.Code)
			assert.Equal(t, tt.statusCode, rec.statusCode)
			assert.Equal(t, tt.content, rr.Body.String())
		})
	}

}
