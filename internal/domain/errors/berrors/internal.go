package berrors

var ErrInternal = New("internal", "internal service error")

// InternalFromErr создает ErrInternal из уже существующей ошибки, добавляя операцию.
func InternalFromErr(op string, err error) *Error {
	if err == nil {
		return nil
	}

	return &Error{Op: op, Message: ErrInternal.Error(), Err: err}
}
