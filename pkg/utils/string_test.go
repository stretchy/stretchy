package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchy/stretchy/pkg/utils"
)

func Test_RandomString(t *testing.T) {
	testCases := []struct {
		name   string
		length int
	}{
		{
			name:   "empty",
			length: 0,
		},
		{
			name:   "with 5 chars",
			length: 5,
		},
		{
			name:   "with 50 chars",
			length: 50,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			result := utils.RandomString(testCase.length)
			assert.Len(t, result, testCase.length)
		})
	}
}
