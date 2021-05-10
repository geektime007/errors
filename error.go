package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

// ErrorCode is the code identifier of the error.
type ErrorCode uint32

// 特殊错误码列表
const (
	// 成功
	ErrorCodeNil ErrorCode = 0
	// 格式错误
	ErrorCodeBadFormat ErrorCode = 1
	// 执行命令失败
	ErrorRunCommandError ErrorCode = 2
	// 超时
	ErrorCodeTimeout ErrorCode = 3
	// 未实现
	ErrorCodeUnimplemented ErrorCode = 4
	// 未知错误
	ErrorCodeUnknown ErrorCode = 111111
)

// 错误的抽象

// Error is the abstraction of marsaxlokk's error.
// Which implements the standard error interface.
type Error struct {
	// Status 表示 http status or grpc code
	Status int `json:"-"`
	// The code of this error.
	Code ErrorCode `json:"code"`
	// The message of this error.
	Message string `json:"message"`
	// The original error for unknown error.
	Origin error `json:"origin"`
	// The hint message
	Hint string `json:"hint"`
}

func (e *Error) String() string { return e.Error() }

// GetCode returns the error code.
func (e *Error) GetCode() ErrorCode {
	if e == nil {
		return ErrorCodeNil
	}
	return e.Code
}

// GetMessage returns the error message.
func (e *Error) GetMessage() string {
	switch e.GetCode() {
	case ErrorCodeNil:
		return "ok"
	case ErrorCodeUnknown:
		return "未知错误"
	default:
		return e.Message
	}
}

// GetStatus return the http status error.
func (e *Error) GetStatus() int {
	if e == nil {
		return 0
	}
	return e.Status
}

// GetOrigin returns the original error.
func (e *Error) GetOrigin() error {
	if e == nil {
		return nil
	}
	return e.Origin
}

// GetHint returns the hint.
func (e *Error) GetHint() string {
	if e == nil {
		// 默认走 message
		return e.GetMessage()
	}
	return e.Hint
}

// Error implements the standard error interface.
// https://golang.org/pkg/builtin/#error
func (e *Error) Error() string {
	status := e.GetStatus()
	code := e.GetCode()
	origin := e.GetOrigin()
	message := e.GetMessage()

	if code == ErrorCodeNil {
		// 成功
		return "ok"
	}

	if origin == nil {
		// 不带额外的原始错误信息
		return fmt.Sprintf("%s[%v][%v]", message, status, code)
	}

	// 带额外的原始错误信息
	return fmt.Sprintf("%s[%v][%v] (%v)", message, status, code, origin)
}

// IsNil returns true if this error is Nil.
func (e *Error) IsNil() bool { return e.GetCode() == ErrorCodeNil }

// IsNotNil return true if this error is not Nil.
func (e *Error) IsNotNil() bool { return e.GetCode() != ErrorCodeNil }

// Ok returns true if this error is Nil.
func (e *Error) Ok() bool { return e.IsNil() }

// Clone this error to get a new Error instance with attaching an origin error.
//
// Example:
//
//	import "github.com/geektime007/errors"
//	db, err := gorm.Open(...)
//	if err != nil {
//		return errors.DatabaseOpenError.CloneWithOriginError(err)
//	}
//
func (e *Error) CloneWithOriginError(origin error) *Error {
	return &Error{
		Status:  e.GetStatus(),
		Code:    e.GetCode(),
		Message: e.GetMessage(),
		Origin:  errors.WithStack(origin),
		Hint:    e.GetHint(),
	}
}

// Clone this error and format it with fmt.Sprintf.
func (e *Error) Format(a ...interface{}) *Error {
	return &Error{
		Status:  e.GetStatus(),
		Code:    e.GetCode(),
		Message: fmt.Sprintf(e.GetMessage(), a...),
		Origin:  errors.WithStack(e.GetOrigin()),
		Hint:    e.Hint,
	}
}

// CloneWithHint
// Clone this error to get a new Error instance with attaching a hint message.
func (e *Error) CloneWithHint(hint string) *Error {
	return &Error{
		Status:  e.GetStatus(),
		Code:    e.GetCode(),
		Message: e.GetMessage(),
		Origin:  errors.WithMessage(e.GetOrigin(), hint),
		Hint:    hint,
	}
}

// Internal errors registry.
var registry = make(map[ErrorCode]*Error)

// NewError returns a new Error.
// Internal: will register an error onto registry.
func NewError(code ErrorCode, status int, message string) *Error {
	// Check code conflicts.
	_, ok := registry[code]
	if ok {
		panic(fmt.Sprintf("conflicts error code: %v", code))
	}
	// Create.
	e := &Error{status, code, message, nil, ""}
	// Register.
	registry[code] = e
	return e
}
