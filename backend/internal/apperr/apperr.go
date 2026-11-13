package apperr

import (
	"errors"
	"fmt"
)

// Kind 领域错误类别，决定 HTTP 状态码
type Kind int

const (
	KindInvalid Kind = iota
	KindUnauthorized
	KindForbidden
	KindNotFound
	KindConflict
	KindInternal
)

// Error 带错误码的领域错误
type Error struct {
	Kind    Kind
	Code    string
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *Error) Unwrap() error { return e.Err }

func New(kind Kind, code, message string) *Error {
	return &Error{Kind: kind, Code: code, Message: message}
}

func Invalid(code, message string) *Error {
	return New(KindInvalid, code, message)
}

func Unauthorized(code, message string) *Error {
	return New(KindUnauthorized, code, message)
}

func Forbidden(code, message string) *Error {
	return New(KindForbidden, code, message)
}

func NotFound(code, message string) *Error {
	return New(KindNotFound, code, message)
}

func Conflict(code, message string) *Error {
	return New(KindConflict, code, message)
}

func Internal(err error) *Error {
	return &Error{Kind: KindInternal, Code: "internal_error", Message: "服务器内部错误", Err: err}
}

// From 从任意错误链中提取领域错误
func From(err error) (*Error, bool) {
	var e *Error
	if errors.As(err, &e) {
		return e, true
	}
	return nil, false
}

func CodeOf(err error) string {
	if e, ok := From(err); ok {
		return e.Code
	}
	return ""
}

func KindOf(err error) (Kind, bool) {
	if e, ok := From(err); ok {
		return e.Kind, true
	}
	return KindInternal, false
}
