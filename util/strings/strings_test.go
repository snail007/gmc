package gstrings

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBytesRef(t *testing.T) {
	a := "abc"
	b := StringToBytes(a)
	assert.Equal(t, []byte("abc"), b)
}

func TestStringRef(t *testing.T) {
	a := []byte("abc")
	b := BytesToString(a)
	assert.Equal(t, "abc", b)
}

func TestHasPrefixAny(t *testing.T) {
	assert.True(t, HasPrefixAny("http://", "https://", "http://"))
	assert.False(t, HasPrefixAny("ftp://", "https://", "http://"))
}

func TestHasHTTPPrefix(t *testing.T) {
	assert.True(t, HasHTTPPrefix("http://"))
	assert.False(t, HasHTTPPrefix("ftp://"))
}

func TestHasSuffix(t *testing.T) {
	assert.True(t, HasSuffixAny("a.txt", ".log", ".txt"))
	assert.True(t, HasSuffixAny("a.txt", "", ".txt"))
	assert.False(t, HasSuffixAny("a.log", ".log1", ".txt1"))
}

func TestReplace(t *testing.T) {
	testCases := []struct {
		input    string
		oldNew   []string
		expected string
	}{
		{"Hello, world!", []string{"world", "Go"}, "Hello, Go!"},
		{"foo bar foo", []string{"foo", "baz", "bar", "qux"}, "baz qux baz"},
		{"abc", []string{"x", "y"}, "abc"}, // No replacement
		{"", []string{}, ""},               // Empty string and no replacements
	}

	for _, testCase := range testCases {
		result := Replace(testCase.input, testCase.oldNew...)
		if result != testCase.expected {
			t.Errorf("For input '%s' with replacements %v, expected '%s', but got '%s'", testCase.input, testCase.oldNew, testCase.expected, result)
		}
	}
}

func TestContainsAny(t *testing.T) {
	assert.False(t, ContainsAny("abcd", "aa", "aaa", "aaa"))
	assert.True(t, ContainsAny("abcd", "ab", "bc", "cde"))
	assert.True(t, ContainsAny("abcd"))
}

func TestContainsAll(t *testing.T) {
	assert.False(t, ContainsAll("abcd", "a", "aaa", "aaa"))
	assert.True(t, ContainsAll("abcd", "ab", "bc", "cd"))
	assert.True(t, ContainsAll("abcd"))
}

func TestRender(t *testing.T) {
	// Define test cases
	tests := []struct {
		input    string
		data     map[string]string
		expected string
	}{
		{
			input:    "Hello, ${NAME}!",
			data:     map[string]string{"NAME": "John"},
			expected: "Hello, John!",
		},
		{
			input:    "This is $CITY, ${COUNTRY}.",
			data:     map[string]string{"CITY": "New York", "COUNTRY": "USA"},
			expected: "This is New York, USA.",
		},
		{
			input:    "No placeholders here!",
			data:     map[string]string{"KEY": "VALUE"},
			expected: "No placeholders here!",
		},
	}

	// Run test cases
	for _, test := range tests {
		result := Render(test.input, test.data)
		if result != test.expected {
			t.Errorf("Render(%s, %v) = %s, want %s", test.input, test.data, result, test.expected)
		}
	}
}
