// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: sync.proto

package sync

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

// define the regex for a UUID once up-front
var _sync_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on GetChangesRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *GetChangesRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetChangesRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetChangesRequestMultiError, or nil if none found.
func (m *GetChangesRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetChangesRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if err := m._validateUuid(m.GetUserId()); err != nil {
		err = GetChangesRequestValidationError{
			field:  "UserId",
			reason: "value must be a valid UUID",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for LastSyncTime

	if len(errors) > 0 {
		return GetChangesRequestMultiError(errors)
	}

	return nil
}

func (m *GetChangesRequest) _validateUuid(uuid string) error {
	if matched := _sync_uuidPattern.MatchString(uuid); !matched {
		return errors.New("invalid uuid format")
	}

	return nil
}

// GetChangesRequestMultiError is an error wrapping multiple validation errors
// returned by GetChangesRequest.ValidateAll() if the designated constraints
// aren't met.
type GetChangesRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetChangesRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetChangesRequestMultiError) AllErrors() []error { return m }

// GetChangesRequestValidationError is the validation error returned by
// GetChangesRequest.Validate if the designated constraints aren't met.
type GetChangesRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetChangesRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetChangesRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetChangesRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetChangesRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetChangesRequestValidationError) ErrorName() string {
	return "GetChangesRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetChangesRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetChangesRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetChangesRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetChangesRequestValidationError{}

// Validate checks the field values on GetChangesResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetChangesResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetChangesResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetChangesResponseMultiError, or nil if none found.
func (m *GetChangesResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GetChangesResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetChanges() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, GetChangesResponseValidationError{
						field:  fmt.Sprintf("Changes[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, GetChangesResponseValidationError{
						field:  fmt.Sprintf("Changes[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return GetChangesResponseValidationError{
					field:  fmt.Sprintf("Changes[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return GetChangesResponseMultiError(errors)
	}

	return nil
}

// GetChangesResponseMultiError is an error wrapping multiple validation errors
// returned by GetChangesResponse.ValidateAll() if the designated constraints
// aren't met.
type GetChangesResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetChangesResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetChangesResponseMultiError) AllErrors() []error { return m }

// GetChangesResponseValidationError is the validation error returned by
// GetChangesResponse.Validate if the designated constraints aren't met.
type GetChangesResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetChangesResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetChangesResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetChangesResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetChangesResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetChangesResponseValidationError) ErrorName() string {
	return "GetChangesResponseValidationError"
}

// Error satisfies the builtin error interface
func (e GetChangesResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetChangesResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetChangesResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetChangesResponseValidationError{}

// Validate checks the field values on PushChangesRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *PushChangesRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PushChangesRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// PushChangesRequestMultiError, or nil if none found.
func (m *PushChangesRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *PushChangesRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if err := m._validateUuid(m.GetUserId()); err != nil {
		err = PushChangesRequestValidationError{
			field:  "UserId",
			reason: "value must be a valid UUID",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	for idx, item := range m.GetChanges() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, PushChangesRequestValidationError{
						field:  fmt.Sprintf("Changes[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, PushChangesRequestValidationError{
						field:  fmt.Sprintf("Changes[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return PushChangesRequestValidationError{
					field:  fmt.Sprintf("Changes[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return PushChangesRequestMultiError(errors)
	}

	return nil
}

func (m *PushChangesRequest) _validateUuid(uuid string) error {
	if matched := _sync_uuidPattern.MatchString(uuid); !matched {
		return errors.New("invalid uuid format")
	}

	return nil
}

// PushChangesRequestMultiError is an error wrapping multiple validation errors
// returned by PushChangesRequest.ValidateAll() if the designated constraints
// aren't met.
type PushChangesRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PushChangesRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PushChangesRequestMultiError) AllErrors() []error { return m }

// PushChangesRequestValidationError is the validation error returned by
// PushChangesRequest.Validate if the designated constraints aren't met.
type PushChangesRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PushChangesRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PushChangesRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PushChangesRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PushChangesRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PushChangesRequestValidationError) ErrorName() string {
	return "PushChangesRequestValidationError"
}

// Error satisfies the builtin error interface
func (e PushChangesRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPushChangesRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PushChangesRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PushChangesRequestValidationError{}

// Validate checks the field values on PushChangesResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *PushChangesResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PushChangesResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// PushChangesResponseMultiError, or nil if none found.
func (m *PushChangesResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *PushChangesResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Success

	// no validation rules for Message

	if len(errors) > 0 {
		return PushChangesResponseMultiError(errors)
	}

	return nil
}

// PushChangesResponseMultiError is an error wrapping multiple validation
// errors returned by PushChangesResponse.ValidateAll() if the designated
// constraints aren't met.
type PushChangesResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PushChangesResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PushChangesResponseMultiError) AllErrors() []error { return m }

// PushChangesResponseValidationError is the validation error returned by
// PushChangesResponse.Validate if the designated constraints aren't met.
type PushChangesResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PushChangesResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PushChangesResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PushChangesResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PushChangesResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PushChangesResponseValidationError) ErrorName() string {
	return "PushChangesResponseValidationError"
}

// Error satisfies the builtin error interface
func (e PushChangesResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPushChangesResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PushChangesResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PushChangesResponseValidationError{}

// Validate checks the field values on DataChange with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *DataChange) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DataChange with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in DataChangeMultiError, or
// nil if none found.
func (m *DataChange) ValidateAll() error {
	return m.validate(true)
}

func (m *DataChange) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for DataId

	// no validation rules for Type

	// no validation rules for Data

	// no validation rules for Metadata

	// no validation rules for Timestamp

	if len(errors) > 0 {
		return DataChangeMultiError(errors)
	}

	return nil
}

// DataChangeMultiError is an error wrapping multiple validation errors
// returned by DataChange.ValidateAll() if the designated constraints aren't met.
type DataChangeMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DataChangeMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DataChangeMultiError) AllErrors() []error { return m }

// DataChangeValidationError is the validation error returned by
// DataChange.Validate if the designated constraints aren't met.
type DataChangeValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DataChangeValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DataChangeValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DataChangeValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DataChangeValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DataChangeValidationError) ErrorName() string { return "DataChangeValidationError" }

// Error satisfies the builtin error interface
func (e DataChangeValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDataChange.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DataChangeValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DataChangeValidationError{}
