package gold_test

import (
	"testing"

	. "github.com/mfcochauxlaberge/gold"

	"github.com/stretchr/testify/assert"
)

func TestNewRunner(t *testing.T) {
	assert := assert.New(t)

	runner := NewRunner("some/path")
	assert.Equal("some/path", runner.Directory)
}
