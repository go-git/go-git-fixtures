package fixtures_test

import (
	"io"
	"testing"

	fixtures "github.com/go-git/go-git-fixtures/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOSFixturePackfile(t *testing.T) {
	t.Parallel()

	f := fixtures.ByTag("packfile").One()
	require.NotNil(t, f)

	osf := fixtures.NewOSFixture(f, t.TempDir())
	file, err := osf.Packfile()
	require.NoError(t, err)
	require.NotNil(t, file)

	t.Cleanup(func() {
		file.Close()
	})

	content, err := io.ReadAll(file)
	require.NoError(t, err)
	assert.NotEmpty(t, content)
}

func TestOSFixtureIdx(t *testing.T) {
	t.Parallel()

	f := fixtures.ByTag("packfile").One()
	require.NotNil(t, f)

	osf := fixtures.NewOSFixture(f, t.TempDir())
	file, err := osf.Idx()
	require.NoError(t, err)
	require.NotNil(t, file)

	t.Cleanup(func() {
		file.Close()
	})

	content, err := io.ReadAll(file)
	require.NoError(t, err)
	assert.NotEmpty(t, content)
}

func TestOSFixtureRev(t *testing.T) {
	t.Parallel()

	f := fixtures.ByTag("rev").One()
	require.NotNil(t, f)

	osf := fixtures.NewOSFixture(f, t.TempDir())
	file, err := osf.Rev()
	require.NoError(t, err)
	require.NotNil(t, file)

	t.Cleanup(func() {
		file.Close()
	})

	content, err := io.ReadAll(file)
	require.NoError(t, err)
	assert.NotEmpty(t, content)
}

func TestOSFixtureIs(t *testing.T) {
	t.Parallel()

	f := fixtures.Basic().One()
	require.NotNil(t, f)

	osf := fixtures.NewOSFixture(f, t.TempDir())

	assert.True(t, osf.Is("packfile"))
	assert.True(t, osf.Is(".git"))
	assert.False(t, osf.Is("nonexistent-tag"))
}

func TestOSFixtureClone(t *testing.T) {
	t.Parallel()

	f := fixtures.Basic().One()
	require.NotNil(t, f)

	osf := fixtures.NewOSFixture(f, t.TempDir())

	clone := osf.Clone()
	require.NotNil(t, clone)

	assert.Equal(t, osf.URL, clone.URL)
	assert.Equal(t, osf.Head, clone.Head)
	assert.Equal(t, osf.PackfileHash, clone.PackfileHash)
	assert.Equal(t, osf.DotGitHash, clone.DotGitHash)
	assert.Equal(t, osf.ObjectsCount, clone.ObjectsCount)
	assert.Equal(t, osf.ObjectFormat, clone.ObjectFormat)
}

func TestOSFixtureDotGit(t *testing.T) {
	t.Parallel()

	f := fixtures.Basic().One()
	require.NotNil(t, f)

	osf := fixtures.NewOSFixture(f, t.TempDir())

	fs, err := osf.DotGit(fixtures.WithMemFS())
	require.NoError(t, err)
	require.NotNil(t, fs)

	files, err := fs.ReadDir("/")
	require.NoError(t, err)
	assert.NotEmpty(t, files)
}
