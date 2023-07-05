//go:build linux || darwin

package uds

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsAvailable(t *testing.T) {
	// Should always be available on Linux and Darwin.
	assert.NoError(t, IsAvailable())
}
