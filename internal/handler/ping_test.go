package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingHandler(t *testing.T) {
	tests := []struct {
		name string
		want struct {
			response   string
			statusCode int
		}
	}{
		{
			name: "it response with pong and success status code",
			want: struct {
				response   string
				statusCode int
			}{
				response:   `{"message": "pong"}`,
				statusCode: http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/ping", nil)
			rr := httptest.NewRecorder()

			handler := PingHandler()
			handler(rr, req)
			assert.Equal(t, tt.want.statusCode, rr.Code)
			assert.JSONEq(t, tt.want.response, rr.Body.String())
		})
	}
}
