package fixtures_test

import (
	"testing"

	fixtures "github.com/go-git/go-git-fixtures/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDotGit(t *testing.T) {
	t.Parallel()

	fs := fixtures.Basic().One().DotGit(fixtures.WithTargetDir(t.TempDir))
	files, err := fs.ReadDir("/")
	require.NoError(t, err)
	assert.Greater(t, len(files), 1)

	fs = fixtures.Basic().One().DotGit(fixtures.WithMemFS())
	files, err = fs.ReadDir("/")
	require.NoError(t, err)
	assert.Greater(t, len(files), 1)
}

//nolint:cyclop
func TestEmbeddedFiles(t *testing.T) {
	t.Parallel()

	for i, f := range fixtures.All() {
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
			if f.Worktree(fixtures.WithMemFS()) == nil {
				assert.Fail(t, "[mem] failed to get worktree", i)
			}

			if f.Worktree(fixtures.WithTargetDir(t.TempDir)) == nil {
				assert.Fail(t, "[tempdir] failed to get worktree", i)
			}
		}

		if f.DotGitHash != "" {
			if f.DotGit(fixtures.WithMemFS()) == nil {
				assert.Fail(t, "[mem] failed to get dotgit", i)
			}

			if f.DotGit(fixtures.WithTargetDir(t.TempDir)) == nil {
				assert.Fail(t, "[tempdir] failed to get dotgit", i)
			}
		}
	}
}

func TestRevFiles(t *testing.T) {
	t.Parallel()

	f := fixtures.ByTag("packfile-sha256").One()

	assert.NotNil(t, f)
	assert.NotNil(t, f.Rev(), "failed to get rev file")
}
