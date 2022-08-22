// Package errors 错误处理,pkg/errors 直接转调
// @author: xs
// @date: 2022/8/8
// @Description: errors
package errors

import (
	"fmt"
	pkgErrs "github.com/pkg/errors"
)

const (
	// UnknownCode is unknown code for error info.
	UnknownCode = 500
	// UnknownReason is unknown reason for error info.
	UnknownReason = ""
)

type Status struct {
	Code    int32  `json:"code"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

// Error is a status error.
type Error struct {
	Status
	cause error
}

func (e *Error) Error() string {
	return fmt.Sprintf("error: code = %d reason = %s message = %s cause = %v", e.Code, e.Reason, e.Message, e.cause)
}

// New returns an error object for the code, message.
func New(code int, reason, message string) *Error {
	return &Error{
		Status: Status{
			Code:    int32(code),
			Message: message,
			Reason:  reason,
		},
	}
}

// Newf New(code fmt.Sprintf(format, a...))
func Newf(code int, reason, format string, a ...interface{}) *Error {
	return New(code, reason, fmt.Sprintf(format, a...))
}

// Errorf returns an error object for the code, message and error info.
func Errorf(code int, reason, format string, a ...interface{}) error {
	return New(code, reason, fmt.Sprintf(format, a...))
}

func FromError(err error) *Error {
	if err == nil {
		return nil
	}
	if se := new(Error); pkgErrs.As(err, &se) {
		return se
	}
	return New(UnknownCode, UnknownReason, err.Error())
}

// pkg errors

func As(err error, target interface{}) bool {
	return pkgErrs.As(err, target)
}

func Cause(err error) error {
	return pkgErrs.Cause(err)
}

func Is(err, target error) bool {
	return pkgErrs.Is(err, target)
}

func Unwrap(err error) error {
	return pkgErrs.Unwrap(err)
}

func WithMessage(err error, message string) error {
	return pkgErrs.WithMessage(err, message)
}
func WithMessagef(err error, format string, args ...interface{}) error {
	return pkgErrs.WithMessagef(err, format, args...)
	return nil
}
func WithStack(err error) error {
	return pkgErrs.WithStack(err)
}
func Wrap(err error, message string) error {
	return pkgErrs.Wrap(err, message)
}
func Wrapf(err error, format string, args ...interface{}) error {
	return pkgErrs.Wrapf(err, format, args...)
}
