package fixtures

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDotGit(t *testing.T) {
	fs := Basic().One().DotGit(WithTargetDir(t.TempDir))
	files, err := fs.ReadDir("/")
	assert.NoError(t, err)
	assert.True(t, len(files) > 1)

	fs = Basic().One().DotGit(WithMemFS())
	files, err = fs.ReadDir("/")
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
			if f.Worktree(WithMemFS()) == nil {
				assert.Fail(t, "[mem] failed to get worktree", i)
			}

			if f.Worktree(WithTargetDir(t.TempDir)) == nil {
				assert.Fail(t, "[tempdir] failed to get worktree", i)
			}
		}

		if f.DotGitHash != "" {
			if f.DotGit(WithMemFS()) == nil {
				assert.Fail(t, "[mem] failed to get dotgit", i)
			}

			if f.DotGit(WithTargetDir(t.TempDir)) == nil {
				assert.Fail(t, "[tempdir] failed to get dotgit", i)
			}
		}
	}
}

func TestRevFiles(t *testing.T) {
	f := ByTag("packfile-sha256").One()

	assert.NotNil(t, f)
	assert.NotNil(t, f.Rev(), "failed to get rev file")
}
