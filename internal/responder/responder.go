package responder

import (
	"context"
	"encoding/json"
	"net/http"

	contextWrapper "github.com/shahbaz275817/prismo/pkg/context"
	"github.com/shahbaz275817/prismo/pkg/errors"
	"github.com/shahbaz275817/prismo/pkg/logger"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Errors  []ErrorItem `json:"errors"`
}

type ErrorItem struct {
	Message      string      `json:"message"`
	MessageTitle string      `json:"title"`
	Code         string      `json:"code"`
	Data         interface{} `json:"data,omitempty"`
}

func NewSuccessResponse(res interface{}) *Response {
	return &Response{
		Success: true,
		Data:    res,
		Errors:  []ErrorItem{},
	}
}

func WriteSuccessResponse(ctx context.Context, w http.ResponseWriter, res interface{}) {
	WriteAnyResponse(ctx, w, NewSuccessResponse(res), http.StatusOK)
}

func WriteResponse(ctx context.Context, w http.ResponseWriter, res *Response, responseCode ...int) {
	WriteAnyResponse(ctx, w, res, responseCode...)
}

func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Add("Content-Type", "application/json; utf8")

	language := contextWrapper.Language(r.Context())
	statusCode, res := handleError(r, err, nil, language)

	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.WithContext(r.Context()).Errorf("unable to write response to client %v", err)
	}
}

func WriteErrorWithData(w http.ResponseWriter, r *http.Request, err error, data interface{}) {
	w.Header().Add("Content-Type", "application/json; utf8")

	language := contextWrapper.Language(r.Context())
	statusCode, res := handleError(r, err, data, language)

	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.WithContext(r.Context()).Errorf("unable to write response to client %v", err)
	}
}

func newErrorResponse(err errors.GenericError, title string, data interface{}, language string) *Response {
	return &Response{
		Success: false,
		Errors: []ErrorItem{
			{
				Code:         err.ErrorID(),
				Message:      errMessage(err, language),
				MessageTitle: errTitle(err, title, language),
				Data:         data,
			},
		},
	}
}

func errMessage(err errors.GenericError, language string) string {
	msg := err.Error()
	if msg == "" {
		msg = err.ErrorID()
	}
	return msg
}

func errTitle(err errors.GenericError, title string, language string) string {
	t := err.ErrTitle()
	if t == "" {
		t = title
	}
	return t
}

func WriteAnyResponse(ctx context.Context, w http.ResponseWriter, res interface{}, responseCode ...int) {
	w.Header().Add("Content-Type", "application/json; utf8")

	if len(responseCode) == 0 {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(responseCode[0])
	}

	if res != nil {
		if err := json.NewEncoder(w).Encode(res); err != nil {
			logger.WithContext(ctx).Errorf("unable to write response to client %v", err)
		}
	}
}

func handleError(r *http.Request, err error, data interface{}, language string) (int, *Response) {
	l := logger.WithRequest(r)
	switch errorType := err.(type) {
	case errors.BadRequestError:
		l.Warnf("%v", err)
		return http.StatusBadRequest, newErrorResponse(errorType, "Bad Request", data, language)
	case errors.NotFoundError:
		l.Warnf("%v", err)
		return http.StatusNotFound, newErrorResponse(errorType, "Not Found", data, language)
	case errors.UnprocessableEntityError:
		l.Warnf("%v", err)
		return http.StatusUnprocessableEntity, newErrorResponse(errorType, "Unprocessable Entity", data, language)
	case errors.UnauthorizedError:
		l.Warnf("%v", err)
		return http.StatusUnauthorized, newErrorResponse(errorType, "Unauthorized", data, language)
	case errors.ForbiddenError:
		l.Warnf("%v", err)
		return http.StatusForbidden, newErrorResponse(errorType, "Forbidden", data, language)
	case *errors.ValidationError:
		l.Warnf("%v", err)
		return http.StatusUnprocessableEntity, newErrorResponse(errorType, "Unprocessable Entity", data, language)
	case errors.TooManyRequestsError:
		l.Warnf("%v", err)
		return http.StatusTooManyRequests, newErrorResponse(errorType, "Too Many Requests", data, language)
	case errors.NotAcceptableError:
		l.Warnf("%v", err)
		return http.StatusNotAcceptable, newErrorResponse(errorType, "Not Acceptable", data, language)
	case errors.DuplicatePackageError:
		return http.StatusUnprocessableEntity, newErrorResponse(errorType, "Duplicate package error", data, language)
	case errors.EntityLockedError:
		l.Warnf("%v", err)
		return http.StatusLocked, newErrorResponse(errorType, "Entity Locked Error", data, language)
	case errors.InternalServerError:
		l.Warnf("%v", err)
		return http.StatusInternalServerError, newErrorResponse(errorType, "Internal Server Error", data, language)
	case errors.ExternalAPICallError:
		l.Warnf("%v", err)
		return newErrorResponseByHTTPStatus(errorType, data, language)
	default:
		l.Warnf("Something Wrong: %v", err)
		return internalErr(errorType.Error(), "Internal Server Error", data)
	}
}

func internalErr(msg string, title string, data interface{}) (int, *Response) {
	return http.StatusInternalServerError, &Response{
		Success: false,
		Errors: []ErrorItem{
			{
				Code:         "amphibian:service:internal_error",
				MessageTitle: "Internal Server Error",
				Message:      msg,
				Data:         data,
			},
		},
	}
}

func newErrorResponseByHTTPStatus(err errors.ExternalAPICallError, data interface{}, language string) (int, *Response) {
	switch err.GetStatusCode() {
	// 404
	case http.StatusNotFound:
		return http.StatusNotFound, newErrorResponse(err, "Not Found", data, language)
	// 406, 422
	case http.StatusNotAcceptable, http.StatusUnprocessableEntity:
		return http.StatusUnprocessableEntity, newErrorResponse(err, "Unprocessable Entity", data, language)
	// 500
	default:
		return internalErr(err.Error(), "Internal Server Error", data)
	}
}
