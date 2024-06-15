package errors

import (
	"errors"

	pkgerrors "github.com/pkg/errors"
)

const (
	Title    string = "ErrTitle"
	fieldKey        = "Field"
)

var (
	ErrNotFound = New("not found")
)

func New(msg string) error {
	return errors.New(msg)
}

func Wrap(err error, message string) error {
	return pkgerrors.Wrap(err, message)
}

func WithStack(err error) error {
	return pkgerrors.WithStack(err)
}

func Errorf(format string, args ...interface{}) error {
	return pkgerrors.Errorf(format, args...)
}

type GenericError interface {
	Error() string
	ErrorID() string
	GetParams() map[string]interface{}
	ErrTitle() string
	SetErrTitle(titleCode string)
}

type genericErrorImpl struct {
	id      string
	message string
	params  map[string]interface{}
}

func (ed genericErrorImpl) Error() string {
	return ed.message
}

func (ed genericErrorImpl) ErrorID() string {
	return ed.id
}

func (ed genericErrorImpl) GetParams() map[string]interface{} {
	return ed.params
}

func (ed genericErrorImpl) ErrTitle() string {
	t, ok := ed.params[Title]
	if !ok {
		return ""
	}
	st, ok := t.(string)
	if !ok {
		return ""
	}
	return st
}

func (ed genericErrorImpl) SetErrTitle(titleCode string) {
	ed.params[Title] = titleCode
}

func (ed genericErrorImpl) SetErrDetails(details *ErrDetails) {
	if details != nil {
		ed.message = details.Message
		ed.params = details.Params
	}
}

func newGenericError(code string, details *ErrDetails) genericErrorImpl {
	msg := ""
	var params map[string]interface{}
	if details != nil {
		msg = details.Message
		params = details.Params
	}
	if msg == "" {
		msg = code
	}
	return genericErrorImpl{
		id:      code,
		message: msg,
		params:  params,
	}
}

type ValidationError struct {
	genericErrorImpl
}

type ErrDetails struct {
	Message string
	Params  map[string]interface{}
}

func NewValidationError(code string, details *ErrDetails) ValidationError {
	return ValidationError{newGenericError(code, details)}
}

type LegacyAuthError struct {
	genericErrorImpl
}

func NewLegacyAuthError(code string, details *ErrDetails) LegacyAuthError {
	return LegacyAuthError{newGenericError(code, details)}
}

type BadRequestError struct {
	genericErrorImpl
}

func NewBadRequestError(code string, details *ErrDetails) BadRequestError {
	return BadRequestError{newGenericError(code, details)}
}

type ForbiddenError struct {
	genericErrorImpl
}

func NewForbiddenError(code string, details *ErrDetails) ForbiddenError {
	return ForbiddenError{newGenericError(code, details)}
}

type MissingFieldError struct {
	genericErrorImpl
}

func NewMissingFieldError(code string, field string) MissingFieldError {
	details := ErrDetails{
		Params: map[string]interface{}{
			"Field": field,
		},
	}
	return MissingFieldError{newGenericError(code, &details)}
}

type InvalidFieldError struct {
	genericErrorImpl
}

func NewInvalidFieldError(code string, field string) InvalidFieldError {
	details := ErrDetails{
		Params: map[string]interface{}{
			fieldKey: field,
		},
	}
	return InvalidFieldError{newGenericError(code, &details)}
}

func (i InvalidFieldError) GetField() string {
	f, ok := i.params[fieldKey]
	if !ok {
		return ""
	}
	if fStr, ok := f.(string); ok {
		return fStr
	}
	return ""

}

type UnknownError struct { // Should  map to 500 status code
	genericErrorImpl
}

type DuplicatePackageError struct {
	genericErrorImpl
}

func NewDuplicatePackageError(code string) DuplicatePackageError {
	return DuplicatePackageError{newGenericError(code, nil)}
}

func NewUnknownError(code string) UnknownError {
	return UnknownError{newGenericError(code, nil)}
}

func NewUnknownErrorWithDetails(code string, details *ErrDetails) UnknownError {
	return UnknownError{newGenericError(code, details)}
}

type UnprocessableEntityError struct { // Should  map to 422 status code
	genericErrorImpl
}

func NewUnprocessableEntityError(code string) UnprocessableEntityError {
	return UnprocessableEntityError{newGenericError(code, nil)}
}

type EntityLockedError struct { // Should  map to 423 status code
	genericErrorImpl
}

func NewEntityLockedError(code string, details *ErrDetails) EntityLockedError {
	return EntityLockedError{newGenericError(code, details)}
}

func NewStatusUnprocessableEntity(code string, details *ErrDetails) UnprocessableEntityError {
	return UnprocessableEntityError{newGenericError(code, details)}
}

type NotFoundError struct {
	genericErrorImpl
}

func NewNotFoundError(code string, details *ErrDetails) NotFoundError {
	return NotFoundError{newGenericError(code, details)}
}

type NotImplementedError struct {
	genericErrorImpl
}

func NewNotImplementedError(code string) NotImplementedError {
	return NotImplementedError{newGenericError(code, nil)}
}

type LegacyNotFoundError struct {
	Err string
}

func (e LegacyNotFoundError) Error() string {
	return e.Err
}

type TooManyRequestsError struct {
	genericErrorImpl
}

func NewTooManyRequestsError(code string, details *ErrDetails) TooManyRequestsError {
	return TooManyRequestsError{newGenericError(code, details)}
}

type InternalServerError struct {
	genericErrorImpl
}

func NewInternalServerError(code string, details *ErrDetails) InternalServerError {
	return InternalServerError{newGenericError(code, details)}
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

type UnauthorizedError struct {
	genericErrorImpl
}

func NewUnauthorizedError(msg string, details *ErrDetails) UnauthorizedError {
	return UnauthorizedError{newGenericError(msg, details)}
}

type NotAcceptableError struct {
	genericErrorImpl
}

func NewNotAcceptableError(code string, details *ErrDetails) NotAcceptableError {
	return NotAcceptableError{newGenericError(code, details)}
}

type ExternalAPICallError struct {
	statusCode int
	genericErrorImpl
}

func (e ExternalAPICallError) Error() string {
	return e.message
}

func (e ExternalAPICallError) GetStatusCode() int {
	return e.statusCode
}

func NewExternalAPICallError(statusCode int, code string, details *ErrDetails) ExternalAPICallError {
	return ExternalAPICallError{statusCode, newGenericError(code, details)}
}
