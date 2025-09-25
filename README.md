# go-git-fixtures [![GoDoc](https://godoc.org/gopkg.in/go-git/go-git-fixtures.v6?status.svg)](https://pkg.go.dev/github.com/go-git/go-git-fixtures/v6) [![Test](https://github.com/go-git/go-git-fixtures/workflows/Test/badge.svg)](https://github.com/go-git/go-git-fixtures/actions?query=workflow%3ATest) [![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/go-git/go-git-fixtures/badge)](https://scorecard.dev/viewer/?uri=github.com/go-git/go-git-fixtures)

git repository fixtures used by [go-git](https://github.com/go-git/go-git)

## Adding new Fixtures

### Adding new pack fixtures

1. Get the `.idx` and `.pack` files from the repository:

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

### Adding new dot fixtures

1. Tarball the contents of .git from a git repository:

```sh
git clone https://<repository>
cd <repository_name>
git checkout <REF>
git gc

tar -czf git.tgz -C .git .
```

2. Get the sha1 of the file: `shasum < git-.tgz`.
3. Move the file using the checksum to `data/git-<checksum>.tgz`
4. Add a new entry in `fixtures.go`:

```
{
	Tags:         []string{"packfile", "<TAG_TO_REFER_TO>"},
	PackfileHash: "<PACK_HASH>",
}
```
