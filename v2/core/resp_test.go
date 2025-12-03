package core

import (
	"reflect"
	"testing"
)

func TestDecode(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected any
	}{
		{
			name:     "Simple String",
			input:    "+OK\r\n",
			expected: "OK",
		},
		{
			name:     "Error",
			input:    "-Error message\r\n",
			expected: "Error message",
		},
		{
			name:     "Integer",
			input:    ":1000\r\n",
			expected: int64(1000),
		},
		{
			name:     "Bulk String",
			input:    "$6\r\nfoobar\r\n",
			expected: "foobar",
		},
		{
			name:     "Array",
			input:    "*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n",
			expected: []any{"foo", "bar"},
		},
		{
			name:     "Nested Array",
			input:    "*2\r\n*2\r\n:12\r\n$3\r\nbar\r\n$3\r\nbaz\r\n",
			expected: []any{[]any{int64(12), "bar"}, "baz"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := Decode([]byte(c.input))
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if !reflect.DeepEqual(got, c.expected) {
				t.Errorf("Expected %v, got %v", c.expected, got)
			}
		})
	}
}
