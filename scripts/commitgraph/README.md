# Generating the commit-graph fixtures

Run from the fixture repository:

```sh
sh scripts/commitgraph/generate.sh
```

The script can also be invoked by absolute path from another directory.

## Requirements

Generation requires Git 2.54 or later, Go (`gofmt`), GNU tar, gzip, sha256sum and
standard POSIX shell utilities. Git 2.54.0 was used for the checked-in fixtures.
The Git version requirement covers SHA-256 commit graphs and acceptance of a
split graph with 256 bases whose one-byte header count wraps to zero.

Consumers load the embedded archives with `DotGit(fixtures.WithMemFS())` and
need no Git executable. Expected commit, tree and parent IDs are available
through `CommitGraphEntries()` in oldest-first order.

## Fixture contents

| Tag | Object format | Commits | Graph layout |
| --- | --- | --- | --- |
| `commit-graph-sha1` | SHA-1 | 3 | Standalone |
| `commit-graph-chain-sha1` | SHA-1 | 3 | Three layers |
| `commit-graph-sha256` | SHA-256 | 3 | Standalone |
| `commit-graph-chain-sha256` | SHA-256 | 3 | Three layers |
| `commit-graph-chain-sha1-257` | SHA-1 | 257 | 257 layers |

Each archive contains a complete bare repository with config, symbolic HEAD,
`refs/heads/main`, loose commit and tree objects, and graphs under `objects/info`.
The three-commit standalone and split fixtures have identical objects and refs.
The 257-layer SHA-1 chain extends the same history with one commit per layer.
Its newest graph contains 256 SHA-1 graph IDs in its BASE chunk and zero in
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

Archives use the existing `data/git-<sha256-of-archive>.tgz` layout. GNU tar
normalizes entry ordering, timestamps and ownership; `gzip -n` removes gzip
timestamps, and a fixed umask stabilizes permissions. Repeating generation with
the same tool versions produces the same archive hashes.

`commit_graph_entries.go` is regenerated with literal commit, tree and parent IDs
in oldest-first order, keyed by archive hash. The script prints each fixture's
unique tag, `DotGitHash`, `Head`, `ObjectsCount` and `ObjectFormat`. If a newer Git
changes graph bytes, update the corresponding registrations in `fixtures.go`
using this output and remove only the superseded commit-graph archives. Do not
replace unrelated fixtures.

After generation, run `make test` and `make validate`. Validation requires a
clean tracked worktree, so run its dirty/pack checks after committing changes.
