package tgz_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git-fixtures/v5/internal/tgz"
	"github.com/stretchr/testify/require"
)

func TestExtractError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		tgz      string
		notFound bool
		wantErr  string
	}{
		{
			tgz:      "not-found",
			notFound: true,
		},
		{
			tgz:     "invalid-gzip.tgz",
			wantErr: "gzip: invalid header",
		},
		{
			tgz:     "not-a-tar.tgz",
			wantErr: "unexpected EOF",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run("tgz path = "+tc.tgz, func(t *testing.T) {
			t.Parallel()

			d, err := os.Getwd()
			require.NoError(t, err)

			source := osfs.New(d + "/fixtures")
			f, err := source.Open(tc.tgz)
			if tc.notFound {
				require.ErrorIs(t, err, os.ErrNotExist)
			} else {
				fs, err := tgz.MemFactory()
				if err != nil {
					panic(err)
				}

				err = tgz.Extract(f, fs)
				require.ErrorContains(t, err, tc.wantErr)
			}
		})
	}
}

func TestExtract(t *testing.T) {
	t.Parallel()

	tests := []struct {
		tgz  string
		tree []string
	}{
		{
			tgz: "test-01.tgz",
			tree: []string{
				"foo.txt",
			},
		}, {
			tgz: "test-02.tgz",
			tree: []string{
				"baz.txt",
				"bla.txt",
				"foo.txt",
			},
		}, {
			tgz: "test-03.tgz",
			tree: []string{
				"bar",
				filepath.Join("bar", "baz.txt"),
				filepath.Join("bar", "foo.txt"),
				"baz",
				filepath.Join("baz", "bar"),
				filepath.Join("baz", "bar", "foo.txt"),
				filepath.Join("baz", "baz"),
				filepath.Join("baz", "baz", "baz"),
				filepath.Join("baz", "baz", "baz", "foo.txt"),
				"foo.txt",
			},
		},
	}

	factories := []struct {
		name    string
		factory func() (billy.Filesystem, error)
	}{
		{name: "mem", factory: tgz.MemFactory},
		{name: "osfs-temp", factory: func() (billy.Filesystem, error) {
			return osfs.New(t.TempDir(), osfs.WithChrootOS()), nil
		}},
	}

	for _, ff := range factories {
		ff := ff
		for _, tc := range tests {
			t.Run(fmt.Sprintf("[%s] tgz path = %s", ff.name, tc.tgz), func(t *testing.T) {
				t.Parallel()

				source := osfs.New("fixtures", osfs.WithChrootOS())
				f, err := source.Open(tc.tgz)
				require.NoError(t, err)

				fs, err := ff.factory()
				if err != nil {
					panic(err)
				}

				err = tgz.Extract(f, fs)
				require.NoError(t, err, "%s: unexpected error extracting: %s", tc.tgz, err)

				for _, path := range tc.tree {
					_, err = fs.Stat(path)
					require.NoError(t, err)
				}
			})
		}
	}
}
