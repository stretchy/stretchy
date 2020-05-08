package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchy/stretchy/pkg/utils"
)

func TestIsAScalar(t *testing.T) {
	testCases := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{
			name:     "int",
			value:    int(123),
			expected: true,
		},
		{
			name:     "int64",
			value:    int64(123),
			expected: true,
		},
		{
			name:     "float32",
			value:    float32(123.1),
			expected: true,
		},
		{
			name:     "float64",
			value:    float64(123.1),
			expected: true,
		},
		{
			name:     "string",
			value:    "af168987-4fc0-4b0e-8377-354e3346599b",
			expected: true,
		},
		{
			name:     "map",
			value:    map[string]interface{}{"a": 1, "b": "662abf66-9ea7-4b82-bde6-a7421bb6c3cb"},
			expected: false,
		},
		{
			name:     "int slice",
			value:    []int{1, 2, 3},
			expected: false,
		},
		{
			name:     "int64 slice",
			value:    []int64{1, 2, 3},
			expected: false,
		},
		{
			name:     "float32 slice",
			value:    []float32{123.1, 1.2, 19.4},
			expected: false,
		},
		{
			name:     "float64 slice",
			value:    []float64{123.1, 1.2, 19.4},
			expected: false,
		},
		{
			name:     "string slice",
			value:    []string{"af168987-4fc0-4b0e-8377-354e3346599b", "a52f5af3-5a4f-4bcc-a27a-63a9e09ef8e2"},
			expected: false,
		},

		{
			name:     "bool true",
			value:    true,
			expected: true,
		},
		{
			name:     "bool false",
			value:    true,
			expected: true,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expected, utils.IsAScalar(testCase.value))
		})
	}
}

func TestPrintScalar(t *testing.T) {
	testCases := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{
			name:     "int",
			value:    int(123),
			expected: "123",
		},
		{
			name:     "int64",
			value:    int64(123),
			expected: "123",
		},
		{
			name:     "float32",
			value:    float32(123.1),
			expected: "123.1",
		},
		{
			name:     "float64",
			value:    float64(123.1),
			expected: "123.1",
		},
		{
			name:     "string",
			value:    "af168987-4fc0-4b0e-8377-354e3346599b",
			expected: "af168987-4fc0-4b0e-8377-354e3346599b",
		},
		{
			name:     "map",
			value:    map[string]interface{}{"a": 1, "b": "662abf66-9ea7-4b82-bde6-a7421bb6c3cb"},
			expected: "",
		},
		{
			name:     "int slice",
			value:    []int{1, 2, 3},
			expected: "",
		},
		{
			name:     "int64 slice",
			value:    []int64{1, 2, 3},
			expected: "",
		},
		{
			name:     "float32 slice",
			value:    []float32{123.1, 1.2, 19.4},
			expected: "",
		},
		{
			name:     "float64 slice",
			value:    []float64{123.1, 1.2, 19.4},
			expected: "",
		},
		{
			name:     "string slice",
			value:    []string{"af168987-4fc0-4b0e-8377-354e3346599b", "a52f5af3-5a4f-4bcc-a27a-63a9e09ef8e2"},
			expected: "",
		},
		{
			name:     "bool true",
			value:    true,
			expected: "true",
		},
		{
			name:     "bool false",
			value:    false,
			expected: "false",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expected, utils.PrintScalar(testCase.value))
		})
	}
}
