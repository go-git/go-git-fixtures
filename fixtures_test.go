package fixtures_test

import (
	"io"
	"strconv"
	"testing"

	"github.com/go-git/go-billy/v6/osfs"
	fixtures "github.com/go-git/go-git-fixtures/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDotGit(t *testing.T) {
	t.Parallel()

	fs, err := fixtures.Basic().One().DotGit(fixtures.WithTargetDir(t.TempDir))
	require.NoError(t, err)

	files, err := fs.ReadDir("/")
	require.NoError(t, err)
	assert.Greater(t, len(files), 1)

	fs, err = fixtures.Basic().One().DotGit(fixtures.WithMemFS())
	require.NoError(t, err)

	files, err = fs.ReadDir("/")
	require.NoError(t, err)
	assert.Greater(t, len(files), 1)
}

func TestEmbeddedFiles(t *testing.T) {
	t.Parallel()

	for i, f := range fixtures.All() {
		if f.PackfileHash != "" {
			file, err := f.Packfile()
			require.NoError(t, err)
			assert.NotNil(t, file, "failed to get pack file", i)
		}

		if f.WorktreeHash != "" {
			wt, err := f.Worktree(fixtures.WithMemFS())
			require.NoError(t, err)
			assert.NotNil(t, wt, "[mem] failed to get worktree", i)

			wt, err = f.Worktree(fixtures.WithTargetDir(t.TempDir))
			require.NoError(t, err)
			assert.NotNil(t, wt, "[tempdir] failed to get worktree", i)
		}

		if f.DotGitHash != "" {
			dot, err := f.DotGit(fixtures.WithMemFS())
			require.NoError(t, err)
			assert.NotNil(t, dot, "[mem] failed to get dotgit", i)

			dot, err = f.DotGit(fixtures.WithTargetDir(t.TempDir))
			require.NoError(t, err)
			assert.NotNil(t, dot, "[tempdir] failed to get dotgit", i)
		}
	}
}

func TestRevFiles(t *testing.T) {
	t.Parallel()

	f := fixtures.ByTag("rev").One()
	require.NotNil(t, f)

	file, err := f.Rev()
	require.NoError(t, err)
	assert.NotNil(t, file, "failed to get rev file")
}

func TestAll(t *testing.T) {
	t.Parallel()

	fs := fixtures.All()

	assert.Len(t, fs, 39)
}

func TestByTag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		tag string
		len int
	}{
		{tag: "packfile", len: 21},
		{tag: "ofs-delta", len: 3},
		{tag: ".git", len: 13},
		{tag: "merge-conflict", len: 1},
		{tag: "worktree", len: 6},
		{tag: "submodule", len: 1},
		{tag: "tags", len: 1},
		{tag: "notes", len: 1},
		{tag: "multi-packfile", len: 1},
		{tag: "diff-tree", len: 7},
	}

	for _, tc := range tests {
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
		t.Run("#"+strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()

			index, err := f.Idx()
			require.NoError(t, err)
			require.NotNil(t, index)

			err = index.Close()
			assert.NoError(t, err)
		})
	}
}

func TestWithMemFS(t *testing.T) {
	t.Parallel()

	f := fixtures.Basic().One()
	require.NotNil(t, f)

	fs, err := f.DotGit(fixtures.WithMemFS())
	require.NoError(t, err)
	require.NotNil(t, fs)

	files, err := fs.ReadDir("/")
	require.NoError(t, err)
	assert.NotEmpty(t, files)

	stat, err := fs.Stat("config")
	require.NoError(t, err)
	assert.NotNil(t, stat)
}

func TestWithTargetDir(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		options []osfs.Option
	}{
		{
			name:    "no options",
			options: nil,
		},
		{
			name:    "with chroot",
			options: []osfs.Option{osfs.WithChrootOS()},
		},
		{
			name:    "with bound",
			options: []osfs.Option{osfs.WithBoundOS()},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := fixtures.Basic().One()
			require.NotNil(t, f)

			fs, err := f.DotGit(fixtures.WithTargetDir(t.TempDir, tc.options...))
			require.NoError(t, err)
			require.NotNil(t, fs)

			files, err := fs.ReadDir("/")
			require.NoError(t, err)
			assert.NotEmpty(t, files)

			stat, err := fs.Stat("config")
			require.NoError(t, err)
			assert.NotNil(t, stat)
		})
	}
}

func TestByObjectFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		objectFormat string
		tag          string
		expectedLen  int
	}{
		{
			name:         "sha1",
			objectFormat: "sha1",
			expectedLen:  37,
		},
		{
			name:         "sha256",
			objectFormat: "sha256",
			expectedLen:  2,
		},
		{
			name:         "sha1 with .git tag",
			objectFormat: "sha1",
			tag:          ".git",
			expectedLen:  12,
		},
		{
			name:         "sha256 with .git tag",
			objectFormat: "sha256",
			tag:          ".git",
			expectedLen:  1,
		},
		{
			name:         "sha1 with packfile tag",
			objectFormat: "sha1",
			tag:          "packfile",
			expectedLen:  20,
		},
		{
			name:         "sha256 with packfile tag",
			objectFormat: "sha256",
			tag:          "packfile",
			expectedLen:  1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := fixtures.ByObjectFormat(tc.objectFormat)

			if tc.tag != "" {
				f = f.ByTag(tc.tag)
			}

			assert.Len(t, f, tc.expectedLen)
		})
	}
}

func TestEnsureIsBare(t *testing.T) {
	t.Parallel()

	f := fixtures.Basic().One()
	require.NotNil(t, f)

	fs, err := f.DotGit(fixtures.WithMemFS())
	require.NoError(t, err)

	err = fixtures.EnsureIsBare(fs)
	require.NoError(t, err)

	cfg, err := fs.Open("config")
	require.NoError(t, err)

	t.Cleanup(func() {
		cfg.Close()
	})

	content, err := io.ReadAll(cfg)
	require.NoError(t, err)

	assert.Contains(t, string(content), "bare = true")
}
