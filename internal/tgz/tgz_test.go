package tgz

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"sort"
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-billy/v5/util"
)

func TestExtractError(t *testing.T) {
	for i, test := range [...]struct {
		tgz    string
		errRgx *regexp.Regexp
	}{
		{
			tgz:    "not-found",
			errRgx: regexp.MustCompile("open not-found: no such file .*"),
		}, {
			tgz:    "fixtures/invalid-gzip.tgz",
			errRgx: regexp.MustCompile("gzip: invalid header"),
		}, {
			tgz:    "fixtures/not-a-tar.tgz",
			errRgx: regexp.MustCompile("unexpected EOF"),
		},
	} {
		com := fmt.Sprintf("%d) tgz path = %s", i, test.tgz)
		tmpFs, err := Extract(osfs.New(""), test.tgz)
		if err == nil {
			t.Errorf("%s: expect an error, but none was returned", com)
		} else if errorNotMatch(err, test.errRgx) {
			t.Errorf("%s:\n\treceived error: %s\n\texpected regexp: %s\n",
				com, err, test.errRgx)
		}

		if tmpFs == nil {
			continue
		}

		err = util.RemoveAll(tmpFs, "")
		if err != nil {
			t.Fatalf("%s: unexpected error removing temporal path: %s", com, err)
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
			tgz: "fixtures/test-01.tgz",
			tree: []string{
				"foo.txt",
			},
		}, {
			tgz: "fixtures/test-02.tgz",
			tree: []string{
				"baz.txt",
				"bla.txt",
				"foo.txt",
			},
		}, {
			tgz: "fixtures/test-03.tgz",
			tree: []string{
				"bar",
				"bar/baz.txt",
				"bar/foo.txt",
				"baz",
				"baz/bar",
				"baz/bar/foo.txt",
				"baz/baz",
				"baz/baz/baz",
				"baz/baz/baz/foo.txt",
				"foo.txt",
			},
		},
	} {
		com := fmt.Sprintf("%d) tgz path = %s", i, test.tgz)

		tmpFs, err := Extract(osfs.New(""), test.tgz)
		if err != nil {
			t.Fatalf("%s: unexpected error extracting: %s", test.tgz, err)
		}
		obt, err := relativeTree(tmpFs)
		if err != nil {
			t.Errorf("%s: unexpected error calculating relative path: %s", com, err)
		}

		sort.Strings(test.tree)
		if !reflect.DeepEqual(obt, test.tree) {
			t.Fatalf("%s:\n\tobtained: %v\n\t expected: %v", com, obt, test.tree)
		}

		err = util.RemoveAll(tmpFs, "")
		if err != nil {
			t.Fatalf("%s: unexpected error removing temporal path: %s", com, err)
		}
	}
}

// relativeTree returns the list of relative paths to the files and
// directories inside a given directory, recursively.
func relativeTree(fs billy.Filesystem) ([]string, error) {
	relPaths := []string{}
	walkFn := func(path string, _ os.FileInfo, _ error) error {
		if path != "" {
			relPaths = append(relPaths, path)
		}
		return nil
	}

	err := util.Walk(fs, "", walkFn)
	return relPaths, err
}
