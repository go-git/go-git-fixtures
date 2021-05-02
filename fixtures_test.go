package fixtures

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDotGit(t *testing.T) {
	fs := Basic().One().DotGit()
	files, err := fs.ReadDir("/")
	assert.NoError(t, err)
	assert.True(t, len(files) > 1)
}
