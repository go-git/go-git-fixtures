# Generating the commit-graph fixtures

Run from the fixture repository:

```sh
sh scripts/commitgraph/generate.sh
```

The script can also be invoked by absolute path from another directory.

## Requirements

Generation requires Git 2.54 or later, Go (`gofmt`), GNU tar, gzip, sha256sum and
standard POSIX shell utilities. Git 2.54.0 was used for the checked-in fixtures.
The version check is not about graph features: sha256 and split graphs predate
it, and Git's loader never reads the base count in header byte 7, so the wrapped
count needs no particular release. Graph bytes, and with them archive hashes,
can differ between Git versions, so the check keeps regeneration on the release
the registered archives were produced with, or a newer one.

Consumers load the embedded archives with `DotGit(fixtures.WithMemFS())` and
need no Git executable. Expected commit, tree and parent IDs are available
through `CommitGraphEntries()` in oldest-first order.

## Fixture contents

Tags name the graph layout and never the object format, so the same test code
covers both formats by combining `ByTag` with `ByObjectFormat`:

```go
fixture := fixtures.ByTag("commit-graph-split").ByObjectFormat("sha256").One()
```

| Tag | Object format | Commits | Graph layout |
| --- | --- | --- | --- |
| `commit-graph-standalone` | sha1 | 3 | Standalone |
| `commit-graph-standalone` | sha256 | 3 | Standalone |
| `commit-graph-split` | sha1 | 3 | Three layers |
| `commit-graph-split` | sha256 | 3 | Three layers |
| `commit-graph-split-257` | sha1 | 257 | 257 layers |

Each archive contains a complete bare repository with config, symbolic HEAD,
`refs/heads/main`, loose commit and tree objects, and graphs under `objects/info`.
The three-commit standalone and split fixtures have identical objects and refs.
The 257-layer sha1 chain extends the same history with one commit per layer.
Its newest graph contains 256 sha1 graph IDs in its BASE chunk and zero in
header byte 7; its chain file lists all 257 graphs, oldest first.

Both author and committer are `Test <test@example.com>`. Commit `i` has message
`commit i\n`, timestamp `1700000000 + i` and timezone `+0000`. Every commit uses
the empty tree and, except the root, the preceding commit as its only parent.
The script uses `commit-tree`, disables global/system Git configuration and
templates, and writes split graphs with `--split=no-merge`.

## Regeneration and validation

The script verifies every generated repository and extracted archive with:

```sh
git --git-dir=<repository> commit-graph verify
git --git-dir=<repository> fsck --strict
```

Generation aborts if verification fails. It also checks chain lengths and the
wrapped header byte. The in-memory tests check deterministic loose objects,
graph checksums, commit IDs and the 256-base boundary.

Archives use the existing `data/git-<sha256-of-archive>.tgz` layout. The script
builds each one with `tar --format=gnu --sort=name --mtime=@1700000000
--owner=0 --group=0 --numeric-owner`, then `gzip -n`. The tar format is pinned
because a tar build default can select another one, which changes the archive
bytes. A fixed umask stabilizes permissions, and `gzip -n` removes gzip
timestamps.

The script unsets every inherited `GIT_*` variable, `TAR_OPTIONS` and `GZIP`,
and creates the repositories with `--ref-format=files`. Left in place, they can
change commit bytes (`i18n.commitEncoding` through `GIT_CONFIG_*`), add graph
chunks (`GIT_TEST_COMMIT_GRAPH_CHANGED_PATHS`), switch the ref storage
(`GIT_DEFAULT_REF_FORMAT`), drop objects from the archives
(`GIT_ALTERNATE_OBJECT_DIRECTORIES`), or change the tar and gzip output.

Repeating generation with the same Git version produces the same object IDs and
graph bytes. The archive hashes can still differ: Git stores loose objects
zlib-compressed, so the hashes also depend on the zlib build Git uses, and on
the gzip that compresses the archives.

`commit_graph_entries.go` is regenerated with literal commit, tree and parent IDs
in oldest-first order, keyed by archive hash. The script prints each fixture's
tag, `ObjectFormat`, `DotGitHash`, `Head` and `ObjectsCount`. If a newer Git
changes graph bytes, update the corresponding registrations in `fixtures.go`
using this output and remove only the superseded commit-graph archives. Do not
replace unrelated fixtures.

After generation, run `make test` and `make validate`. Validation requires a
clean tracked worktree, so run its dirty/pack checks after committing changes.
