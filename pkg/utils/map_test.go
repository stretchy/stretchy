package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchy/stretchy/pkg/utils"
)

func TestMapHasKey_NoSliceProvided(t *testing.T) {
	_, err := utils.MapHasKey(1, 1)
	assert.Error(t, err)
}

func TestMapHasKey(t *testing.T) {
	testCases := []struct {
		name     string
		key      interface{}
		haystack interface{}
		expected bool
	}{
		{
			name: "string fail",
			key:  "key-1",
			haystack: map[interface{}]interface{}{
				"a":       1,
				"key-123": 2,
			},
			expected: false,
		},
		{
			name: "string success",
			key:  "key-123",
			haystack: map[interface{}]interface{}{
				"a":       1,
				"key-123": 2,
			},
			expected: true,
		},
		{
			name: "int fail",
			key:  2,
			haystack: map[interface{}]interface{}{
				1: 1,
				4: 2,
			},
			expected: false,
		},
		{
			name: "int success",
			key:  4,
			haystack: map[interface{}]interface{}{
				1: 1,
				4: 2,
			},
			expected: true,
		},
		{
			name: "float fail",
			key:  2.1,
			haystack: map[interface{}]interface{}{
				1.1: 1,
				2.4: 2,
			},
			expected: false,
		},
		{
			name: "float success",
			key:  1.1,
			haystack: map[interface{}]interface{}{
				1.1: 1,
				2.4: 2,
			},
			expected: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result, err := utils.MapHasKey(tc.key, tc.haystack)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}
