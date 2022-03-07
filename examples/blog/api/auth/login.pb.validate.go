// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: auth/login.proto

package auth

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on GetTokenRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *GetTokenRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetTokenRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetTokenRequestMultiError, or nil if none found.
func (m *GetTokenRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetTokenRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if l := len(m.GetUsername()); l < 2 || l > 128 {
		err := GetTokenRequestValidationError{
			field:  "Username",
			reason: "value length must be between 2 and 128 bytes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for Password

	if len(errors) > 0 {
		return GetTokenRequestMultiError(errors)
	}

	return nil
}

// GetTokenRequestMultiError is an error wrapping multiple validation errors
// returned by GetTokenRequest.ValidateAll() if the designated constraints
// aren't met.
type GetTokenRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetTokenRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetTokenRequestMultiError) AllErrors() []error { return m }

// GetTokenRequestValidationError is the validation error returned by
// GetTokenRequest.Validate if the designated constraints aren't met.
type GetTokenRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetTokenRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetTokenRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetTokenRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetTokenRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetTokenRequestValidationError) ErrorName() string { return "GetTokenRequestValidationError" }

// Error satisfies the builtin error interface
func (e GetTokenRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetTokenRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetTokenRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetTokenRequestValidationError{}

// Validate checks the field values on GetTokenReply with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *GetTokenReply) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetTokenReply with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in GetTokenReplyMultiError, or
// nil if none found.
func (m *GetTokenReply) ValidateAll() error {
	return m.validate(true)
}

func (m *GetTokenReply) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return GetTokenReplyMultiError(errors)
	}

	return nil
}

// GetTokenReplyMultiError is an error wrapping multiple validation errors
// returned by GetTokenReply.ValidateAll() if the designated constraints
// aren't met.
type GetTokenReplyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetTokenReplyMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetTokenReplyMultiError) AllErrors() []error { return m }

// GetTokenReplyValidationError is the validation error returned by
// GetTokenReply.Validate if the designated constraints aren't met.
type GetTokenReplyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetTokenReplyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetTokenReplyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetTokenReplyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetTokenReplyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetTokenReplyValidationError) ErrorName() string { return "GetTokenReplyValidationError" }

// Error satisfies the builtin error interface
func (e GetTokenReplyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetTokenReply.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetTokenReplyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetTokenReplyValidationError{}