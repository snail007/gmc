package gerror

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewError(t *testing.T) {
	err := New("test error")
	if err == nil {
		t.Error("New error should not be nil")
	}

	if err.Error() != "[test error]" {
		t.Errorf("Unexpected error message. Expected: 'test error'. Got: '%s'", err.Error())
	}

	stackFrames := err.StackFrames()
	if len(stackFrames) == 0 {
		t.Error("StackFrames should not be empty")
	}
}

func TestWrapError(t *testing.T) {
	innerErr := errors.New("inner error")
	err := New(innerErr).Wrap("outer error")

	if err.Error() != "outer error" {
		t.Errorf("Unexpected error message. Expected: 'outer error: inner error'. Got: '%s'", err.Error())
	}

	stackFrames := err.StackFrames()
	if len(stackFrames) == 0 {
		t.Error("StackFrames should not be empty")
	}
}

func TestWrapPrefixError(t *testing.T) {
	innerErr := errors.New("inner error")
	err := New(innerErr).WrapPrefix("outer error", "prefix", 0)

	if err.Error() != "prefix: outer error" {
		t.Errorf("Unexpected error message. Expected: 'prefix: outer error: inner error'. Got: '%s'", err.Error())
	}

	stackFrames := err.StackFrames()
	if len(stackFrames) == 0 {
		t.Error("StackFrames should not be empty")
	}
}

func TestErrorf(t *testing.T) {
	err := New().Errorf("error: %s", "test")

	if err.Error() != "error: test" {
		t.Errorf("Unexpected error message. Expected: 'error: test'. Got: '%s'", err.Error())
	}

	stackFrames := err.StackFrames()
	if len(stackFrames) == 0 {
		t.Error("StackFrames should not be empty")
	}
}

func TestErrorStack(t *testing.T) {
	err := New("test error").WrapPrefix("inner error", "prefix", 0)
	expectedErrorStack := fmt.Sprintf("%s test error\n%s", err.TypeName(), string(err.Stack()))
	assert.Contains(t, expectedErrorStack, "*errors.errorString test error")
}
