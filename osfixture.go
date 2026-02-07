package fixtures

import (
	"io"
	"slices"

	"github.com/go-git/go-billy/v6"
	"github.com/go-git/go-billy/v6/osfs"
)

type OSFixture struct {
	*Fixture

	dir string
}

// NewOSFixture converts a Fixture which is based on embedfs, into
// an OS based fixture.
func NewOSFixture(f *Fixture, dir string) *OSFixture {
	return &OSFixture{
		Fixture: f,
		dir:     dir,
	}
}

func (f *OSFixture) Is(tag string) bool {
	return f.Fixture.Is(tag)
}

func (f *OSFixture) Packfile() (billy.File, error) {
	file, err := f.Fixture.Packfile()
	if err != nil {
		return nil, err
	}

	return embedToOsfs(f.dir, file)
}

func (f *OSFixture) Idx() (billy.File, error) {
	file, err := f.Fixture.Idx()
	if err != nil {
		return nil, err
	}

	return embedToOsfs(f.dir, file)
}

func (f *OSFixture) Rev() (billy.File, error) {
	file, err := f.Fixture.Rev()
	if err != nil {
		return nil, err
	}

	return embedToOsfs(f.dir, file)
}

func (f *OSFixture) DotGit(opts ...Option) (billy.Filesystem, error) {
	return f.Fixture.DotGit(opts...)
}

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
