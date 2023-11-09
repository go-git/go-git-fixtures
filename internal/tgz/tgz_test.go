package tgz

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"testing"

	"github.com/go-git/go-billy/v5/osfs"
)

func TestExtractError(t *testing.T) {
	for i, test := range [...]struct {
		tgz    string
		errRgx *regexp.Regexp
	}{
		{
			tgz:    "not-found",
			errRgx: regexp.MustCompile("open not-found: (The system cannot find the file specified|no such file).*"),
		}, {
			tgz:    filepath.Join("fixtures", "invalid-gzip.tgz"),
			errRgx: regexp.MustCompile("gzip: invalid header"),
		}, {
			tgz:    filepath.Join("fixtures", "not-a-tar.tgz"),
			errRgx: regexp.MustCompile("unexpected EOF"),
		},
	} {
		com := fmt.Sprintf("%d) tgz path = %s", i, test.tgz)
		_, err, cleanup := Extract(osfs.New(""), test.tgz)
		defer cleanup() // Clean up temporary dirs.

		if err == nil {
			t.Errorf("%s: expect an error, but none was returned", com)
		} else if errorNotMatch(err, test.errRgx) {
			t.Errorf("%s:\n\treceived error: %s\n\texpected regexp: %s\n",
				com, err, test.errRgx)
		}
	}
}

func errorNotMatch(err error, regexp *regexp.Regexp) bool {
	return !regexp.MatchString(err.Error())
}

func TestExtract(t *testing.T) {
	for i, test := range [...]struct {
		tgz  string
		tree []string
	}{
		{
			tgz: filepath.Join("fixtures", "test-01.tgz"),
			tree: []string{
				"foo.txt",
			},
		}, {
			tgz: filepath.Join("fixtures", "test-02.tgz"),
			tree: []string{
				"baz.txt",
				"bla.txt",
				"foo.txt",
			},
		}, {
			tgz: filepath.Join("fixtures", "test-03.tgz"),
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
	} {
		com := fmt.Sprintf("%d) tgz path = %s", i, test.tgz)

		path, err, cleanup := Extract(osfs.New(""), test.tgz)
		defer cleanup() // Clean up temporary dirs.

		if err != nil {
			t.Fatalf("%s: unexpected error extracting: %s", test.tgz, err)
		}

		obt, err := relativeTree(path.Root())
		if err != nil {
			t.Errorf("%s: unexpected error calculating relative path: %s", com, err)
		}

		sort.Strings(test.tree)
		if !reflect.DeepEqual(obt, test.tree) {
			t.Fatalf("%s:\n\tobtained: %v\n\t expected: %v", com, obt, test.tree)
		}
	}
}

// relativeTree returns the list of relative paths to the files and
// directories inside a given directory, recursively.
func relativeTree(dir string) ([]string, error) {
	dir = filepath.Clean(dir)

	absPaths := []string{}
	walkFn := func(path string, _ os.FileInfo, _ error) error {
		absPaths = append(absPaths, path)
		return nil
	}

	_ = filepath.Walk(dir, walkFn)

	return toRelative(absPaths[1:], dir)
}

// toRelative returns the relative paths (form b) of the list of paths in l.
func toRelative(l []string, b string) ([]string, error) {
	r := []string{}
	for _, p := range l {
		rel, err := filepath.Rel(b, p)
		if err != nil {
			return nil, err
		}
		r = append(r, rel)
	}

	return r, nil
}
