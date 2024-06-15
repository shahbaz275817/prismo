package middleware

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/shahbaz275817/prismo/internal/config"
)

func TestWithHTTPAuth(t *testing.T) {
	type args struct {
		authorization string
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Returns ok when Authorization is valid",
			args: args{
				authorization: "Basic bG1zX3BvcnRhbDptdFVOVzAzdkdjT0lUVGN6bHlmcld3PT0=",
			},
			want: http.StatusOK,
		},
		{
			name: "Returns status unauthorized when Authorization is not valid",
			args: args{
				authorization: "Basic bG1zX3BvcnRhbDptdFVOVzAzdkdj",
			},
			want: http.StatusUnauthorized,
		},
		{
			name: "Returns status unauthorized when Authorization is null",
			args: args{
				authorization: "",
			},
			want: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			config.Load()

			if got := makeAPICallWithHTTPAuth(tt.args.authorization); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithHTTPAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func makeAPICallWithHTTPAuth(authorization string) (responseCode int) {
	req := httptest.NewRequest("GET", "/end_point", nil)
	w := httptest.NewRecorder()
	req.Header.Set("Authorization", authorization)
	dummyHandler := func(w http.ResponseWriter, r *http.Request) {}
	WithHTTPAuth(http.HandlerFunc(dummyHandler)).ServeHTTP(w, req)

	return w.Code
}
