package user

import "errors"

const (
	DefaultMinLen = 8
)

var (
	ErrPasswordTooSmall = errors.New("password length is too small")
)

type SimplePasswordValidator struct {
	minLen int
}

func NewSimplePasswordValidator(minLen int) *SimplePasswordValidator {
	return &SimplePasswordValidator{minLen: minLen}
}

func (v *SimplePasswordValidator) Validate(password string) error {
	if len(password) < v.minLen {
		return ErrPasswordTooSmall
	}

	return nil
}
