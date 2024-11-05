package fixtures_test

import (
	"strconv"
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

func TestAll(t *testing.T) {
	fs := fixtures.All()

	assert.Len(t, fs, 38)
}

func TestByTag(t *testing.T) {
	t.Parallel()
	tests := []struct {
		tag string
		len int
	}{
		{tag: "packfile", len: 20},
		{tag: "ofs-delta", len: 3},
		{tag: ".git", len: 12},
		{tag: "merge-conflict", len: 1},
		{tag: "worktree", len: 6},
		{tag: "submodule", len: 1},
		{tag: "tags", len: 1},
		{tag: "notes", len: 1},
		{tag: "multi-packfile", len: 1},
		{tag: "diff-tree", len: 7},
		{tag: "packfile-sha256", len: 1},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.tag, func(t *testing.T) {
			t.Parallel()

			f := fixtures.ByTag(tc.tag)
			assert.Len(t, f, tc.len)
		})
	}
}

func TestByURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		URL string
		len int
	}{
		{URL: "https://github.com/git-fixtures/root-references.git", len: 1},
		{URL: "https://github.com/git-fixtures/basic.git", len: 9},
		{URL: "https://github.com/git-fixtures/submodule.git", len: 1},
		{URL: "https://github.com/src-d/go-git.git", len: 1},
		{URL: "https://github.com/git-fixtures/tags.git", len: 1},
		{URL: "https://github.com/spinnaker/spinnaker.git", len: 1},
		{URL: "https://github.com/jamesob/desk.git", len: 1},
		{URL: "https://github.com/cpcs499/Final_Pres_P.git", len: 1},
		{URL: "https://github.com/github/gem-builder.git", len: 1},
		{URL: "https://github.com/githubtraining/example-branches.git", len: 1},
		{URL: "https://github.com/rumpkernel/rumprun-xen.git", len: 1},
		{URL: "https://github.com/mcuadros/skeetr.git", len: 1},
		{URL: "https://github.com/dezfowler/LiteMock.git", len: 1},
		{URL: "https://github.com/tyba/storable.git", len: 1},
		{URL: "https://github.com/toqueteos/ts3.git", len: 1},
		{URL: "https://github.com/git-fixtures/empty.git", len: 1},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.URL, func(t *testing.T) {
			t.Parallel()

			f := fixtures.ByURL(tc.URL)
			assert.Len(t, f, tc.len)
		})
	}
}

func TestIdx(t *testing.T) {
	t.Parallel()

	for i, f := range fixtures.ByTag("packfile") {
		f := f
		t.Run("#"+strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()

			index := f.Idx()
			assert.NotNil(t, index)

			err := index.Close()
			assert.NoError(t, err)
		})
	}
}
