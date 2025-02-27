package fixtures

import (
	"github.com/go-git/go-billy/v6"
	"github.com/go-git/go-billy/v6/osfs"
	"github.com/go-git/go-git-fixtures/v5/internal/tgz"
)

type Option func(*options)

type options struct {
	fsFactory func() (billy.Filesystem, error)
}

func newOptions() *options {
	return &options{
		fsFactory: tgz.MemFactory,
	}
}

// WithMemFS returns the option of using memfs for the fs created for Fixtures.
func WithMemFS() Option {
	return func(o *options) {
		o.fsFactory = tgz.MemFactory
	}
}

// WithTargetDir returns the option of using an OS-based filesystem based on a target dir.
// The target dir will be based on the name returned from dirName, which aligns with tempdir
// functions in different testing frameworks (e.g. t.TempDir, c.MkDir).
//
// The caller is responsible for removing the dir from disk. Therefore, it is recommended
// to delegate that to the testing framework:
//
// Go:
//
//	WithTargetDir(t.TempDir)
//
// Check Framework:
//
//	WithTargetDir(c.Mkdir)
func WithTargetDir(dirName func() string) Option {
	return func(o *options) {
		o.fsFactory = func() (billy.Filesystem, error) {
			return osfs.New(dirName(), osfs.WithChrootOS()), nil
		}
	}
}
