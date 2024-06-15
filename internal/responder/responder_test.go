package responder

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/shahbaz275817/prismo/pkg/errors"
)

func TestNewSuccessResponse(t *testing.T) {
	type args struct {
		res interface{}
	}
	tests := []struct {
		name string
		args args
		want *Response
	}{
		{
			name: "Return new response object",
			args: args{res: "test"},
			want: &Response{
				Success: true,
				Data:    "test",
				Errors:  []ErrorItem{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSuccessResponse(tt.args.res); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSuccessResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriteError(t *testing.T) {
	type args struct {
		w   *httptest.ResponseRecorder
		r   *http.Request
		err error
	}
	type wantRes struct {
		statusCode int
		res        *Response
	}

	req := httptest.NewRequest("Get", "/test", nil)

	tests := []struct {
		name string
		args args
		want wantRes
	}{
		{
			name: "Return bad request if the error is of type bad request",
			args: args{
				w:   httptest.NewRecorder(),
				r:   req,
				err: errors.NewBadRequestError("some bad req message", nil),
			},
			want: wantRes{
				statusCode: http.StatusBadRequest,
				res: &Response{
					Success: false,
					Data:    nil,
					Errors: []ErrorItem{
						{
							Message:      "some bad req message",
							Code:         "some bad req message",
							MessageTitle: "Bad Request",
						},
					},
				},
			},
		},
		{
			name: "Return not found if the error is of type not found",
			args: args{
				w:   httptest.NewRecorder(),
				r:   req,
				err: errors.NewNotFoundError("some not found message", nil),
			},
			want: wantRes{
				statusCode: http.StatusNotFound,
				res: &Response{
					Success: false,
					Data:    nil,
					Errors: []ErrorItem{
						{
							Message:      "some not found message",
							Code:         "some not found message",
							MessageTitle: "Not Found",
						},
					},
				},
			},
		},
		{
			name: "Return unprocessable entity if the error is of type unprocessable",
			args: args{
				w:   httptest.NewRecorder(),
				r:   req,
				err: errors.NewUnprocessableEntityError("some unprocessable message"),
			},
			want: wantRes{
				statusCode: http.StatusUnprocessableEntity,
				res: &Response{
					Success: false,
					Data:    nil,
					Errors: []ErrorItem{
						{
							Message:      "some unprocessable message",
							Code:         "some unprocessable message",
							MessageTitle: "Unprocessable Entity",
						},
					},
				},
			},
		},
		{
			name: "Return unauthorized if the error is of type unauthorized",
			args: args{
				w:   httptest.NewRecorder(),
				r:   req,
				err: errors.NewUnauthorizedError("some unauthorized message", nil),
			},
			want: wantRes{
				statusCode: http.StatusUnauthorized,
				res: &Response{
					Success: false,
					Data:    nil,
					Errors: []ErrorItem{
						{
							Message:      "some unauthorized message",
							Code:         "some unauthorized message",
							MessageTitle: "Unauthorized",
						},
					},
				},
			},
		},
		{
			name: "Return forbidden if the error is of type forbidden",
			args: args{
				w:   httptest.NewRecorder(),
				r:   req,
				err: errors.NewForbiddenError("some forbidden message", nil),
			},
			want: wantRes{
				statusCode: http.StatusForbidden,
				res: &Response{
					Success: false,
					Data:    nil,
					Errors: []ErrorItem{
						{
							Message:      "some forbidden message",
							Code:         "some forbidden message",
							MessageTitle: "Forbidden",
						},
					},
				},
			},
		},
		{
			name: "Return not found if get external API call error with status not found",
			args: args{
				w:   httptest.NewRecorder(),
				r:   req,
				err: errors.NewExternalAPICallError(http.StatusNotFound, "not found", nil),
			},
			want: wantRes{
				statusCode: http.StatusNotFound,
				res: &Response{
					Success: false,
					Data:    nil,
					Errors: []ErrorItem{
						{
							Message:      "not found",
							Code:         "not found",
							MessageTitle: "Not Found",
						},
					},
				},
			},
		},
		{
			name: "Return unprocessable entity if get external API call error with status unprocessable entity",
			args: args{
				w:   httptest.NewRecorder(),
				r:   req,
				err: errors.NewExternalAPICallError(http.StatusUnprocessableEntity, "Unprocessable entity", nil),
			},
			want: wantRes{
				statusCode: http.StatusUnprocessableEntity,
				res: &Response{
					Success: false,
					Data:    nil,
					Errors: []ErrorItem{
						{
							Message:      "Unprocessable entity",
							Code:         "Unprocessable entity",
							MessageTitle: "Unprocessable Entity",
						},
					},
				},
			},
		},
		{
			name: "Return internal error if get external API call error with status bad request",
			args: args{
				w:   httptest.NewRecorder(),
				r:   req,
				err: errors.NewExternalAPICallError(http.StatusBadRequest, "bad request", nil),
			},
			want: wantRes{
				statusCode: http.StatusInternalServerError,
				res: &Response{
					Success: false,
					Data:    nil,
					Errors: []ErrorItem{
						{
							Message:      "bad request",
							Code:         "amphibian:service:internal_error",
							MessageTitle: "Internal Server Error",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteError(tt.args.w, tt.args.r, tt.args.err)
			resp := tt.args.w.Result()
			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			expectedBody, err := json.Marshal(tt.want.res)
			assert.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, string(expectedBody), strings.Trim(string(body), "\n"))
		})
	}
}

func TestWriteResponse(t *testing.T) {
	type args struct {
		w            *httptest.ResponseRecorder
		res          *Response
		responseCode []int
	}
	type wantRes struct {
		statusCode int
		res        *Response
	}

	res := &Response{
		Success: true,
		Data:    "test",
		Errors:  nil,
	}

	tests := []struct {
		name string
		args args
		want wantRes
	}{
		{
			name: "Return response with given status code and data",
			args: args{
				w:            httptest.NewRecorder(),
				res:          res,
				responseCode: []int{http.StatusCreated},
			},
			want: wantRes{
				statusCode: http.StatusCreated,
				res:        res,
			},
		},
		{
			name: "Return response with status ok if not status code given",
			args: args{
				w:            httptest.NewRecorder(),
				res:          res,
				responseCode: []int{},
			},
			want: wantRes{
				statusCode: http.StatusOK,
				res:        res,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteResponse(context.Background(), tt.args.w, tt.args.res, tt.args.responseCode...)
			resp := tt.args.w.Result()
			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			expectedBody, err := json.Marshal(tt.want.res)
			assert.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, string(expectedBody), strings.Trim(string(body), "\n"))
		})
	}
}
