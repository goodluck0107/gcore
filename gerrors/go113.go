package gerrors

import (
	"errors"
)

// New Wrapping for gerrors.New standard library
func New(text string) error { return errors.New(text) }

// Is Wrapping for gerrors.Is standard library
func Is(err, target error) bool { return errors.Is(err, target) }

// As Wrapping for gerrors.As standard library
func As(err error, target any) bool { return errors.As(err, target) }

// Unwrap Wrapping for gerrors.Unwrap standard library
func Unwrap(err error) error { return errors.Unwrap(err) }
