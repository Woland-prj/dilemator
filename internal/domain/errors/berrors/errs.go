package berrors

import (
	"fmt"
)

type Error struct {
	Op      string
	Message string
	Err     error
}

// Error реализует интерфейс error.
func (e *Error) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("%s: %s", e.Op, e.Message)
	}

	return fmt.Sprintf("%s: %s: %v", e.Op, e.Message, e.Err)
}

// Unwrap позволяет использовать errors.Is / errors.As.
func (e *Error) Unwrap() error {
	return e.Err
}

// New создает базовую ошибку.
func New(op, message string) *Error {
	return &Error{Op: op, Message: message}
}

// Wrap оборачивает другую ошибку с добавлением контекста.
func Wrap(op, message string, err error) *Error {
	if err == nil {
		return New(op, message)
	}

	return &Error{
		Op:      op,
		Message: message,
		Err:     err,
	}
}

// FromErr создает Error из уже существующей ошибки, добавляя операцию.
func FromErr(op string, err error) *Error {
	if err == nil {
		return nil
	}

	return &Error{Op: op, Message: err.Error(), Err: err}
}
