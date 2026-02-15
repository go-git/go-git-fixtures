package fixtures

import (
	"io"
	"slices"

	"github.com/go-git/go-billy/v6"
	"github.com/go-git/go-billy/v6/osfs"
)

// OSFixture wraps a Fixture and provides OS filesystem-based access to fixture
// files instead of using the embedded filesystem. This is useful when tests
// require real filesystem operations to exercise a specific execution path.
type OSFixture struct {
	*Fixture

	// dir is the base directory where temporary files will be created.
	dir string
}

// NewOSFixture converts a Fixture which is based on embedfs, into
// an OS-based fixture. The dir parameter specifies the base directory
// where temporary files will be created when accessing packfiles, indexes,
// or rev files.
func NewOSFixture(f *Fixture, dir string) *OSFixture {
	return &OSFixture{
		Fixture: f,
		dir:     dir,
	}
}

// Is reports whether the fixture has the specified tag.
func (f *OSFixture) Is(tag string) bool {
	return f.Fixture.Is(tag)
}

// Packfile returns the packfile as an OS-based file.
func (f *OSFixture) Packfile() (billy.File, error) {
	file, err := f.Fixture.Packfile()
	if err != nil {
		return nil, err
	}

	return embedToOsfs(f.dir, file)
}

// Idx returns the packfile index as an OS-based file.
func (f *OSFixture) Idx() (billy.File, error) {
	file, err := f.Fixture.Idx()
	if err != nil {
		return nil, err
	}

	return embedToOsfs(f.dir, file)
}

// Rev returns the reverse index file as an OS-based file.
func (f *OSFixture) Rev() (billy.File, error) {
	file, err := f.Fixture.Rev()
	if err != nil {
		return nil, err
	}

	return embedToOsfs(f.dir, file)
}

// DotGit returns the .git directory filesystem. This delegates to the
// underlying Fixture's DotGit method.
func (f *OSFixture) DotGit(opts ...Option) (billy.Filesystem, error) {
	return f.Fixture.DotGit(opts...)
}

// Clone creates a deep copy of the OSFixture.
func (f *OSFixture) Clone() *OSFixture {
	nf := &OSFixture{
		Fixture: &Fixture{
			URL:          f.URL,
			DotGitHash:   f.DotGitHash,
			Head:         f.Head,
			PackfileHash: f.PackfileHash,
			WorktreeHash: f.WorktreeHash,
			ObjectsCount: f.ObjectsCount,
			Tags:         slices.Clone(f.Tags),
			ObjectFormat: f.ObjectFormat,
		},
		dir: f.dir,
	}

	return nf
}

// Worktree returns the worktree filesystem. This delegates to the
// underlying Fixture's Worktree method.
func (f *OSFixture) Worktree(opts ...Option) (billy.Filesystem, error) {
	return f.Fixture.Worktree(opts...)
}

func embedToOsfs(dir string, f billy.File) (billy.File, error) {
	defer f.Close()

	fs := osfs.New(dir)

	out, err := fs.TempFile("", "embed")
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(out, f)
	if err != nil {
		return nil, err
	}

	_, err = out.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	return out, nil
}
