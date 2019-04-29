// Package wraperr provides some of the useful helpers from github.com/pkg/errors implemented for xerrors

package wraperr

import (
	"fmt"

	"golang.org/x/xerrors"
)

// WithStack annotates err with a stack trace at the point WithStack was called.
// If err is nil, WithStack returns nil.
func WithStack(err error) error {
	if err != nil {
		err = &withStack{
			err:   err,
			frame: xerrors.Caller(1),
		}
	}
	return err
}

type withStack struct {
	err   error
	frame xerrors.Frame
}

func (w *withStack) Error() string {
	return w.err.Error()
}

func (e *withStack) FormatError(p xerrors.Printer) (next error) {
	p.Print(e.err)
	e.frame.Format(p)
	return e.err
}

func (w *withStack) Format(s fmt.State, v rune) { xerrors.FormatError(w, s, v) }

func (w *withStack) Unwrap() error {
	return w.err
}

// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	err = &withMessage{
		err: err,
		msg: msg,
	}
	return &withStack{
		err:   err,
		frame: xerrors.Caller(1),
	}
}

// Wrapf returns an error annotating err with a stack trace
// at the point Wrapf is called, and the format specifier.
// If err is nil, Wrapf returns nil.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	err = &withMessage{
		err: err,
		msg: fmt.Sprintf(format, args...),
	}
	return &withStack{
		err:   err,
		frame: xerrors.Caller(1),
	}
}

// WithMessage annotates err with a new message.
// If err is nil, WithMessage returns nil.
func WithMessage(err error, message string) error {
	if err == nil {
		return nil
	}
	return &withMessage{
		err: err,
		msg: message,
	}
}

// WithMessagef annotates err with the format specifier.
// If err is nil, WithMessagef returns nil.
func WithMessagef(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &withMessage{
		err: err,
		msg: fmt.Sprintf(format, args...),
	}
}

type withMessage struct {
	err error
	msg string
}

func (w *withMessage) Error() string { return w.msg + ": " + w.err.Error() }

func (w *withMessage) Format(s fmt.State, v rune) { xerrors.FormatError(w, s, v) }

func (w *withMessage) FormatError(p xerrors.Printer) (next error) {
	p.Print(w.msg)
	return w.err
}

func (w *withMessage) Unwrap() error {
	return w.err
}

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the xerrors.Wrapper
// interface.
//
// If the error does not implement xerrors.Wrapper, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Cause(err error) error {
	if err == nil {
		return nil
	}
	for err != nil {
		cause, ok := err.(xerrors.Wrapper)
		if !ok {
			break
		}
		err = cause.Unwrap()
	}
	return err
}
