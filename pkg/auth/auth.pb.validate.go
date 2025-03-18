// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: auth.proto

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

// define the regex for a UUID once up-front
var _auth_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on RegisterRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *RegisterRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on RegisterRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// RegisterRequestMultiError, or nil if none found.
func (m *RegisterRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *RegisterRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if l := utf8.RuneCountInString(m.GetUsername()); l < 3 || l > 20 {
		err := RegisterRequestValidationError{
			field:  "Username",
			reason: "value length must be between 3 and 20 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if !_RegisterRequest_Username_Pattern.MatchString(m.GetUsername()) {
		err := RegisterRequestValidationError{
			field:  "Username",
			reason: "value does not match regex pattern \"^[a-zA-Z0-9_]+$\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if l := utf8.RuneCountInString(m.GetPassword()); l < 8 || l > 50 {
		err := RegisterRequestValidationError{
			field:  "Password",
			reason: "value length must be between 8 and 50 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if err := m._validateEmail(m.GetEmail()); err != nil {
		err = RegisterRequestValidationError{
			field:  "Email",
			reason: "value must be a valid email address",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return RegisterRequestMultiError(errors)
	}

	return nil
}

func (m *RegisterRequest) _validateHostname(host string) error {
	s := strings.ToLower(strings.TrimSuffix(host, "."))

	if len(host) > 253 {
		return errors.New("hostname cannot exceed 253 characters")
	}

	for _, part := range strings.Split(s, ".") {
		if l := len(part); l == 0 || l > 63 {
			return errors.New("hostname part must be non-empty and cannot exceed 63 characters")
		}

		if part[0] == '-' {
			return errors.New("hostname parts cannot begin with hyphens")
		}

		if part[len(part)-1] == '-' {
			return errors.New("hostname parts cannot end with hyphens")
		}

		for _, r := range part {
			if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' {
				return fmt.Errorf("hostname parts can only contain alphanumeric characters or hyphens, got %q", string(r))
			}
		}
	}

	return nil
}

func (m *RegisterRequest) _validateEmail(addr string) error {
	a, err := mail.ParseAddress(addr)
	if err != nil {
		return err
	}
	addr = a.Address

	if len(addr) > 254 {
		return errors.New("email addresses cannot exceed 254 characters")
	}

	parts := strings.SplitN(addr, "@", 2)

	if len(parts[0]) > 64 {
		return errors.New("email address local phrase cannot exceed 64 characters")
	}

	return m._validateHostname(parts[1])
}

// RegisterRequestMultiError is an error wrapping multiple validation errors
// returned by RegisterRequest.ValidateAll() if the designated constraints
// aren't met.
type RegisterRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m RegisterRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m RegisterRequestMultiError) AllErrors() []error { return m }

// RegisterRequestValidationError is the validation error returned by
// RegisterRequest.Validate if the designated constraints aren't met.
type RegisterRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RegisterRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RegisterRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RegisterRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RegisterRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RegisterRequestValidationError) ErrorName() string { return "RegisterRequestValidationError" }

// Error satisfies the builtin error interface
func (e RegisterRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRegisterRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RegisterRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RegisterRequestValidationError{}

var _RegisterRequest_Username_Pattern = regexp.MustCompile("^[a-zA-Z0-9_]+$")

// Validate checks the field values on RegisterResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *RegisterResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on RegisterResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// RegisterResponseMultiError, or nil if none found.
func (m *RegisterResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *RegisterResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for UserId

	// no validation rules for Message

	if len(errors) > 0 {
		return RegisterResponseMultiError(errors)
	}

	return nil
}

// RegisterResponseMultiError is an error wrapping multiple validation errors
// returned by RegisterResponse.ValidateAll() if the designated constraints
// aren't met.
type RegisterResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m RegisterResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m RegisterResponseMultiError) AllErrors() []error { return m }

// RegisterResponseValidationError is the validation error returned by
// RegisterResponse.Validate if the designated constraints aren't met.
type RegisterResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RegisterResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RegisterResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RegisterResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RegisterResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RegisterResponseValidationError) ErrorName() string { return "RegisterResponseValidationError" }

// Error satisfies the builtin error interface
func (e RegisterResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRegisterResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RegisterResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RegisterResponseValidationError{}

// Validate checks the field values on LoginRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *LoginRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on LoginRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in LoginRequestMultiError, or
// nil if none found.
func (m *LoginRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *LoginRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if l := utf8.RuneCountInString(m.GetUsername()); l < 3 || l > 20 {
		err := LoginRequestValidationError{
			field:  "Username",
			reason: "value length must be between 3 and 20 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if l := utf8.RuneCountInString(m.GetPassword()); l < 8 || l > 50 {
		err := LoginRequestValidationError{
			field:  "Password",
			reason: "value length must be between 8 and 50 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return LoginRequestMultiError(errors)
	}

	return nil
}

// LoginRequestMultiError is an error wrapping multiple validation errors
// returned by LoginRequest.ValidateAll() if the designated constraints aren't met.
type LoginRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m LoginRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m LoginRequestMultiError) AllErrors() []error { return m }

// LoginRequestValidationError is the validation error returned by
// LoginRequest.Validate if the designated constraints aren't met.
type LoginRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LoginRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LoginRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LoginRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LoginRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LoginRequestValidationError) ErrorName() string { return "LoginRequestValidationError" }

// Error satisfies the builtin error interface
func (e LoginRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLoginRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LoginRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LoginRequestValidationError{}

// Validate checks the field values on LoginResponse with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *LoginResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on LoginResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in LoginResponseMultiError, or
// nil if none found.
func (m *LoginResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *LoginResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Token

	// no validation rules for ExpiresIn

	if len(errors) > 0 {
		return LoginResponseMultiError(errors)
	}

	return nil
}

// LoginResponseMultiError is an error wrapping multiple validation errors
// returned by LoginResponse.ValidateAll() if the designated constraints
// aren't met.
type LoginResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m LoginResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m LoginResponseMultiError) AllErrors() []error { return m }

// LoginResponseValidationError is the validation error returned by
// LoginResponse.Validate if the designated constraints aren't met.
type LoginResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LoginResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LoginResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LoginResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LoginResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LoginResponseValidationError) ErrorName() string { return "LoginResponseValidationError" }

// Error satisfies the builtin error interface
func (e LoginResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLoginResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LoginResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LoginResponseValidationError{}

// Validate checks the field values on ValidateTokenRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ValidateTokenRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ValidateTokenRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ValidateTokenRequestMultiError, or nil if none found.
func (m *ValidateTokenRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *ValidateTokenRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetToken()) < 10 {
		err := ValidateTokenRequestValidationError{
			field:  "Token",
			reason: "value length must be at least 10 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return ValidateTokenRequestMultiError(errors)
	}

	return nil
}

// ValidateTokenRequestMultiError is an error wrapping multiple validation
// errors returned by ValidateTokenRequest.ValidateAll() if the designated
// constraints aren't met.
type ValidateTokenRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ValidateTokenRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ValidateTokenRequestMultiError) AllErrors() []error { return m }

// ValidateTokenRequestValidationError is the validation error returned by
// ValidateTokenRequest.Validate if the designated constraints aren't met.
type ValidateTokenRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ValidateTokenRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ValidateTokenRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ValidateTokenRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ValidateTokenRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ValidateTokenRequestValidationError) ErrorName() string {
	return "ValidateTokenRequestValidationError"
}

// Error satisfies the builtin error interface
func (e ValidateTokenRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sValidateTokenRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ValidateTokenRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ValidateTokenRequestValidationError{}

// Validate checks the field values on ValidateTokenResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ValidateTokenResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ValidateTokenResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ValidateTokenResponseMultiError, or nil if none found.
func (m *ValidateTokenResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *ValidateTokenResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Valid

	// no validation rules for UserId

	if len(errors) > 0 {
		return ValidateTokenResponseMultiError(errors)
	}

	return nil
}

// ValidateTokenResponseMultiError is an error wrapping multiple validation
// errors returned by ValidateTokenResponse.ValidateAll() if the designated
// constraints aren't met.
type ValidateTokenResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ValidateTokenResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ValidateTokenResponseMultiError) AllErrors() []error { return m }

// ValidateTokenResponseValidationError is the validation error returned by
// ValidateTokenResponse.Validate if the designated constraints aren't met.
type ValidateTokenResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ValidateTokenResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ValidateTokenResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ValidateTokenResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ValidateTokenResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ValidateTokenResponseValidationError) ErrorName() string {
	return "ValidateTokenResponseValidationError"
}

// Error satisfies the builtin error interface
func (e ValidateTokenResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sValidateTokenResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ValidateTokenResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ValidateTokenResponseValidationError{}

// Validate checks the field values on GenerateOTPRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GenerateOTPRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GenerateOTPRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GenerateOTPRequestMultiError, or nil if none found.
func (m *GenerateOTPRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GenerateOTPRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if err := m._validateUuid(m.GetUserId()); err != nil {
		err = GenerateOTPRequestValidationError{
			field:  "UserId",
			reason: "value must be a valid UUID",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return GenerateOTPRequestMultiError(errors)
	}

	return nil
}

func (m *GenerateOTPRequest) _validateUuid(uuid string) error {
	if matched := _auth_uuidPattern.MatchString(uuid); !matched {
		return errors.New("invalid uuid format")
	}

	return nil
}

// GenerateOTPRequestMultiError is an error wrapping multiple validation errors
// returned by GenerateOTPRequest.ValidateAll() if the designated constraints
// aren't met.
type GenerateOTPRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GenerateOTPRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GenerateOTPRequestMultiError) AllErrors() []error { return m }

// GenerateOTPRequestValidationError is the validation error returned by
// GenerateOTPRequest.Validate if the designated constraints aren't met.
type GenerateOTPRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GenerateOTPRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GenerateOTPRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GenerateOTPRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GenerateOTPRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GenerateOTPRequestValidationError) ErrorName() string {
	return "GenerateOTPRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GenerateOTPRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGenerateOTPRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GenerateOTPRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GenerateOTPRequestValidationError{}

// Validate checks the field values on GenerateOTPResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GenerateOTPResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GenerateOTPResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GenerateOTPResponseMultiError, or nil if none found.
func (m *GenerateOTPResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GenerateOTPResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for OtpCode

	if len(errors) > 0 {
		return GenerateOTPResponseMultiError(errors)
	}

	return nil
}

// GenerateOTPResponseMultiError is an error wrapping multiple validation
// errors returned by GenerateOTPResponse.ValidateAll() if the designated
// constraints aren't met.
type GenerateOTPResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GenerateOTPResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GenerateOTPResponseMultiError) AllErrors() []error { return m }

// GenerateOTPResponseValidationError is the validation error returned by
// GenerateOTPResponse.Validate if the designated constraints aren't met.
type GenerateOTPResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GenerateOTPResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GenerateOTPResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GenerateOTPResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GenerateOTPResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GenerateOTPResponseValidationError) ErrorName() string {
	return "GenerateOTPResponseValidationError"
}

// Error satisfies the builtin error interface
func (e GenerateOTPResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGenerateOTPResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GenerateOTPResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GenerateOTPResponseValidationError{}

// Validate checks the field values on ValidateOTPRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ValidateOTPRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ValidateOTPRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ValidateOTPRequestMultiError, or nil if none found.
func (m *ValidateOTPRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *ValidateOTPRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if err := m._validateUuid(m.GetUserId()); err != nil {
		err = ValidateOTPRequestValidationError{
			field:  "UserId",
			reason: "value must be a valid UUID",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetOtpCode()) != 6 {
		err := ValidateOTPRequestValidationError{
			field:  "OtpCode",
			reason: "value length must be 6 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)

	}

	if len(errors) > 0 {
		return ValidateOTPRequestMultiError(errors)
	}

	return nil
}

func (m *ValidateOTPRequest) _validateUuid(uuid string) error {
	if matched := _auth_uuidPattern.MatchString(uuid); !matched {
		return errors.New("invalid uuid format")
	}

	return nil
}

// ValidateOTPRequestMultiError is an error wrapping multiple validation errors
// returned by ValidateOTPRequest.ValidateAll() if the designated constraints
// aren't met.
type ValidateOTPRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ValidateOTPRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ValidateOTPRequestMultiError) AllErrors() []error { return m }

// ValidateOTPRequestValidationError is the validation error returned by
// ValidateOTPRequest.Validate if the designated constraints aren't met.
type ValidateOTPRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ValidateOTPRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ValidateOTPRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ValidateOTPRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ValidateOTPRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ValidateOTPRequestValidationError) ErrorName() string {
	return "ValidateOTPRequestValidationError"
}

// Error satisfies the builtin error interface
func (e ValidateOTPRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sValidateOTPRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ValidateOTPRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ValidateOTPRequestValidationError{}

// Validate checks the field values on ValidateOTPResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ValidateOTPResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ValidateOTPResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ValidateOTPResponseMultiError, or nil if none found.
func (m *ValidateOTPResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *ValidateOTPResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Valid

	if len(errors) > 0 {
		return ValidateOTPResponseMultiError(errors)
	}

	return nil
}

// ValidateOTPResponseMultiError is an error wrapping multiple validation
// errors returned by ValidateOTPResponse.ValidateAll() if the designated
// constraints aren't met.
type ValidateOTPResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ValidateOTPResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ValidateOTPResponseMultiError) AllErrors() []error { return m }

// ValidateOTPResponseValidationError is the validation error returned by
// ValidateOTPResponse.Validate if the designated constraints aren't met.
type ValidateOTPResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ValidateOTPResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ValidateOTPResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ValidateOTPResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ValidateOTPResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ValidateOTPResponseValidationError) ErrorName() string {
	return "ValidateOTPResponseValidationError"
}

// Error satisfies the builtin error interface
func (e ValidateOTPResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sValidateOTPResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ValidateOTPResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ValidateOTPResponseValidationError{}
