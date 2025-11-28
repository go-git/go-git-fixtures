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

func (f *OSFixture) Packfile() billy.File {
	return embedToOsfs(f.dir, f.Fixture.Packfile())
}

func (f *OSFixture) Idx() billy.File {
	return embedToOsfs(f.dir, f.Fixture.Idx())
}

func (f *OSFixture) Rev() billy.File {
	return embedToOsfs(f.dir, f.Fixture.Rev())
}

func (f *OSFixture) DotGit(opts ...Option) billy.Filesystem {
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
		},
		dir: f.dir,
	}

	return nf
}

func (f *OSFixture) Worktree(opts ...Option) billy.Filesystem {
	return f.Fixture.Worktree(opts...)
}

func embedToOsfs(dir string, f billy.File) billy.File {
	defer f.Close()

	fs := osfs.New(dir)

	out, err := fs.TempFile("", "embed")
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(out, f)
	if err != nil {
		panic(err)
	}

	_, err = out.Seek(0, io.SeekStart)
	if err != nil {
		panic(err)
	}

	return out
}
