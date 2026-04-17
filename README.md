# go-git-fixtures [![GoDoc](https://godoc.org/gopkg.in/go-git/go-git-fixtures.v6?status.svg)](https://pkg.go.dev/github.com/go-git/go-git-fixtures/v6) [![Test](https://github.com/go-git/go-git-fixtures/workflows/Test/badge.svg)](https://github.com/go-git/go-git-fixtures/actions?query=workflow%3ATest) [![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/go-git/go-git-fixtures/badge)](https://scorecard.dev/viewer/?uri=github.com/go-git/go-git-fixtures)

git repository fixtures used by [go-git](https://github.com/go-git/go-git)

## Adding new Fixtures

### Adding new pack fixtures

1. Get the `.idx`, `.rev` and `.pack` files from the repository:

```sh
git clone https://<repository>
cd <repository_name>
git checkout <REF>
git gc

ls .git/objects/pack/
```

2. Copy them into `/data`.
3. Add a new entry in `fixtures.go`:

```
{
	Tags:         []string{"packfile", "<TAG_TO_REFER_TO>"},
	PackfileHash: "<PACK_HASH>",
}
```

### Adding new bitmap fixtures

Bitmap fixtures are `.bitmap` files that accompany existing pack files.
Git does not support generating pack bitmaps without also repacking. Instead,
generate a multi-pack-index (MIDX) bitmap and convert it using the
`doctor-bitmap` tool:

1. In a repository that contains the target pack, generate a MIDX bitmap:

```sh
git multi-pack-index write --bitmap

ls .git/objects/pack/multi-pack-index-*.bitmap
```

2. Convert the MIDX bitmap to a pack bitmap:

```sh
go run ./internal/cmd/doctor-bitmap \
  .git/objects/pack/multi-pack-index-<HASH>.bitmap \
  .git/objects/pack/pack-<PACK_HASH>.pack \
  data/pack-<PACK_HASH>.bitmap
```

3. Add the `"bitmap"` tag to the fixture's `Tags` in `fixtures.go`.

### Adding new dot fixtures

1. Tarball the contents of .git from a git repository:

```sh
git clone https://<repository>
cd <repository_name>
git checkout <REF>
git gc

tar -czf git.tgz -C .git .
```

2. Get the sha1/sha256 of the file: `sha1sum < git.tgz`.
3. Move the file using the checksum to `data/git-<checksum>.tgz`
4. Add a new entry in `fixtures.go`:

```
{
	Tags:         []string{".git", "<TAG_TO_REFER_TO>"},
	DotGitHash: "<GIT_TAR_HASH>",
}
```



### Adding new worktree fixtures

1. Tarball the contents of the cloned repository:

```sh
git clone https://<repository> <repository_name>

tar -czf worktree.tgz -C <repository_name> .
```

2. Get the sha1/sha256 of the file: `sha256sum < worktree.tgz`.
3. Move the file using the checksum to `data/worktree-<checksum>.tgz`
4. Add a new entry in `fixtures.go`:

```
{
	Tags:         []string{"worktree", "<TAG_TO_REFER_TO>"},
	WorktreeHash: "<WORKTREE_TAR_HASH>",
}
```
