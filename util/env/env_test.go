package genv

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestAccessor(t *testing.T) {
	// Set up test cases
	accessor := NewAccessor("PREFIX_")

	// Test Get method
	key := "KEY"
	expectedValue := "value"
	os.Setenv("PREFIX_KEY", expectedValue)
	result := accessor.Get(key)
	if result.String() != expectedValue {
		t.Errorf("Get(%s) = %s, want %s", key, result, expectedValue)
	}

	// Test Lookup method
	lookupKey := "LOOKUP_KEY"
	lookupExpectedValue := "lookup_value"
	os.Setenv("PREFIX_LOOKUP_KEY", lookupExpectedValue)
	lookupResult, found := accessor.Lookup(lookupKey)
	if !found || lookupResult.String() != lookupExpectedValue {
		t.Errorf("Lookup(%s) = %s, found = %t, want %s, true", lookupKey, lookupResult, found, lookupExpectedValue)
	}
	lookupResult, found = accessor.Lookup(lookupKey + "none")
	assert.False(t, found)
	assert.Nil(t, lookupResult)

	// Test Set method
	setKey := "SET_KEY"
	setValue := "set_value"
	accessor.Set(setKey, setValue)
	setResult := os.Getenv("PREFIX_SET_KEY")
	if setResult != setValue {
		t.Errorf("Set(%s, %s) did not set the environment variable correctly", setKey, setValue)
	}

	// Test Unset method
	unsetKey := "UNSET_KEY"
	unsetValue := "unset_value"
	os.Setenv("PREFIX_UNSET_KEY", unsetValue)
	accessor.Unset(unsetKey)
	unsetResult := os.Getenv("PREFIX_UNSET_KEY")
	if unsetResult != "" {
		t.Errorf("Unset(%s) did not unset the environment variable", unsetKey)
	}
}
