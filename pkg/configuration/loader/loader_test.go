package loader_test

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func getScenarioPath(t *testing.T, scenarioName string) string {
	currentDir, err := os.Getwd()
	assert.NoError(t, err)

	return filepath.Join(currentDir, "test_scenarios", scenarioName)
}
