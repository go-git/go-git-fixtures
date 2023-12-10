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

func TestEmbeddedFiles(t *testing.T) {
	for i, f := range fixtures {
		if f.PackfileHash != "" {
			if f.Packfile() == nil {
				assert.Fail(t, "failed to get pack file", i)
			}
			// skip pack file ee4fef0 as it does not have an idx file.
			if f.PackfileHash != "ee4fef0ef8be5053ebae4ce75acf062ddf3031fb" && f.Idx() == nil {
				assert.Fail(t, "failed to get idx file", i)
			}
		}

		if f.WorktreeHash != "" {
			if f.Worktree() == nil {
				assert.Fail(t, "failed to get worktree", i)
			}
		}

		if f.DotGitHash != "" {
			if f.DotGit() == nil {
				assert.Fail(t, "failed to get dotgit", i)
			}
		}
	}
}

func TestRevFiles(t *testing.T) {
	f := ByTag("packfile-sha256").One()

	if f.Rev() == nil {
		assert.Fail(t, "failed to get rev file")
	}
}
