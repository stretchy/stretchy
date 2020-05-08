package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchy/stretchy/pkg/utils"
)

func TestInSlice_NoSliceProvided(t *testing.T) {
	_, err := utils.InSlice(1, 1)
	assert.Error(t, err)
}

func TestInSlice(t *testing.T) {
	testCases := []struct {
		name     string
		element  interface{}
		slice    interface{}
		expected bool
	}{
		{
			name:     "int success",
			element:  1,
			slice:    []int{1, 2, 3, 4, 5},
			expected: true,
		},
		{
			name:     "int fail",
			element:  18,
			slice:    []int{1, 2, 3, 4, 5},
			expected: false,
		},
		{
			name:     "float32 success",
			element:  float32(5.7),
			slice:    []float32{7.9, 2.4, 5.7},
			expected: true,
		},
		{
			name:     "float32 fail",
			element:  float32(1.8),
			slice:    []float32{7.9, 2.4, 5.7},
			expected: false,
		},
		{
			name:     "float64 success",
			element:  float64(5.7),
			slice:    []float64{7.9, 2.4, 5.7},
			expected: true,
		},
		{
			name:     "float64 fail",
			element:  float64(1.8),
			slice:    []float64{7.9, 2.4, 5.7},
			expected: false,
		},
		{
			name:     "string success",
			element:  "dolor",
			slice:    []string{"Lorem", "ipsum", "dolor", "sit", "amet"},
			expected: true,
		},
		{
			name:     "string fail",
			element:  "dolors",
			slice:    []string{"Lorem", "ipsum", "dolor", "sit", "amet"},
			expected: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result, err := utils.InSlice(tc.element, tc.slice)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}
