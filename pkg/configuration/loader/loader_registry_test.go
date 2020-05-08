package loader_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchy/stretchy/pkg/configuration/loader"
)

func TestLoaderRegistry_GetByFormat_Unknown(t *testing.T) {
	loaderRegistry := loader.NewRegistry("")
	l, err := loaderRegistry.GetByFormat("unknown-format")
	assert.Nil(t, l)
	assert.Error(t, err)
}

func TestLoaderRegistry_GetByFormat(t *testing.T) {
	loaderRegistry := loader.NewRegistry("")

	testCases := []struct {
		format         string
		expectedLoader loader.Loader
	}{
		{
			format:         "yml",
			expectedLoader: &loader.YAMLLoader{},
		},
		{
			format:         "yaml",
			expectedLoader: &loader.YAMLLoader{},
		},
		{
			format:         "json",
			expectedLoader: &loader.JSONLoader{},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.format, func(t *testing.T) {
			l, err := loaderRegistry.GetByFormat(testCase.format)
			assert.NoError(t, err)
			assert.IsType(t, testCase.expectedLoader, l)
		})
	}
}
