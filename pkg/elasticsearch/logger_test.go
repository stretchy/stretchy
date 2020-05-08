package elasticsearch

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func Test_newLogger(t *testing.T) {
	logger := newLogger()

	assert.IsType(t, &log.Logger{}, logger)
}
