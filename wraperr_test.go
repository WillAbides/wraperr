package wraperr

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"
)

func TestWrapNil(t *testing.T) {
	assert.Nil(t, Wrap(nil, "no error"))
}

func TestWrap(t *testing.T) {
	tests := []struct {
		err     error
		message string
		want    string
	}{
		{io.EOF, "read error", "read error: EOF"},
		{Wrap(io.EOF, "read error"), "client error", "client error: read error: EOF"},
	}
	for _, tt := range tests {
		got := Wrap(tt.err, tt.message).Error()
		assert.Equal(t, tt.want, got)
	}
}

type nilError struct{}

func (nilError) Error() string { return "nil error" }

func TestCause(t *testing.T) {
	x := xerrors.New("error")
	tests := []struct {
		err  error
		want error
	}{{
		// nil error is nil
		err:  nil,
		want: nil,
	}, {
		// explicit nil error is nil
		err:  (error)(nil),
		want: nil,
	}, {
		// typed nil is nil
		err:  (*nilError)(nil),
		want: (*nilError)(nil),
	}, {
		// uncaused error is unaffected
		err:  io.EOF,
		want: io.EOF,
	}, {
		// caused error returns cause
		err:  Wrap(io.EOF, "ignored"),
		want: io.EOF,
	}, {
		err:  x, // return from errors.New
		want: x,
	}, {
		err:  WithMessage(nil, "whoops"),
		want: nil,
	}, {
		err:  WithMessage(io.EOF, "whoops"),
		want: io.EOF,
	}, {
		err:  WithStack(nil),
		want: nil,
	}, {
		err:  WithStack(io.EOF),
		want: io.EOF,
	}}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, tt.want, Cause(tt.err))
		})
	}
}

func TestWrapfNil(t *testing.T) {
	got := Wrapf(nil, "no error")
	assert.Nil(t, got)
}

func TestWrapf(t *testing.T) {
	tests := []struct {
		err     error
		message string
		want    string
	}{
		{io.EOF, "read error", "read error: EOF"},
		{Wrapf(io.EOF, "read error without format specifiers"), "client error", "client error: read error without format specifiers: EOF"},
		{Wrapf(io.EOF, "read error with %d format specifier", 1), "client error", "client error: read error with 1 format specifier: EOF"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := Wrapf(tt.err, tt.message).Error()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWithStackNil(t *testing.T) {
	got := WithStack(nil)
	assert.Nil(t, got)
}

func TestWithStack(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{io.EOF, "EOF"},
		{WithStack(io.EOF), "EOF"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := WithStack(tt.err).Error()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWithMessageNil(t *testing.T) {
	got := WithMessage(nil, "no error")
	assert.Nil(t, got)
}

func TestWithMessage(t *testing.T) {
	tests := []struct {
		err     error
		message string
		want    string
	}{
		{io.EOF, "read error", "read error: EOF"},
		{WithMessage(io.EOF, "read error"), "client error", "client error: read error: EOF"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := WithMessage(tt.err, tt.message).Error()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWithMessagefNil(t *testing.T) {
	got := WithMessagef(nil, "no error")
	assert.Nil(t, got)
}

func TestWithMessagef(t *testing.T) {
	tests := []struct {
		err     error
		message string
		want    string
	}{
		{io.EOF, "read error", "read error: EOF"},
		{WithMessagef(io.EOF, "read error without format specifier"), "client error", "client error: read error without format specifier: EOF"},
		{WithMessagef(io.EOF, "read error with %d format specifier", 1), "client error", "client error: read error with 1 format specifier: EOF"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := WithMessagef(tt.err, tt.message).Error()
			assert.Equal(t, tt.want, got)
		})
	}
}

// errors.New, etc values are not expected to be compared by value
// but the change in errors#27 made them incomparable. Assert that
// various kinds of errors have a functional equality operator, even
// if the result of that equality is always false.
func TestErrorEquality(t *testing.T) {
	vals := []error{
		nil,
		io.EOF,
		errors.New("EOF"),
		xerrors.New("EOF"),
		xerrors.Errorf("EOF"),
		Wrap(io.EOF, "EOF"),
		Wrapf(io.EOF, "EOF%d", 2),
		WithMessage(nil, "whoops"),
		WithMessage(io.EOF, "whoops"),
		WithStack(io.EOF),
		WithStack(nil),
	}

	for i := range vals {
		for j := range vals {
			_ = vals[i] == vals[j] // mustn't panic
		}
	}
}
