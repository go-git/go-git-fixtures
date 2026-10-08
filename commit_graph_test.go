package fixtures_test

import (
	"compress/zlib"
	"crypto/sha1" //nolint:gosec // Git's SHA-1 object format is part of the fixture contract.
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/go-git/go-billy/v6"
	fixtures "github.com/go-git/go-git-fixtures/v6"
	"github.com/stretchr/testify/require"
)

func TestCommitGraphFixtures(t *testing.T) {
	t.Parallel()

	for _, tag := range []string{
		"commit-graph-sha1",
		"commit-graph-chain-sha1",
		"commit-graph-sha256",
		"commit-graph-chain-sha256",
		"commit-graph-chain-sha1-257",
	} {
		t.Run(tag, func(t *testing.T) {
			t.Parallel()

			matches := fixtures.ByTag(tag)
			require.Len(t, matches, 1)
			fixture := matches.One()
			require.True(t, fixture.Is(".git"))
			fs, err := fixture.DotGit(fixtures.WithMemFS())
			require.NoError(t, err)
			_, err = fs.Stat("refs/heads/main")
			require.NoError(t, err)
			verifyCommitGraphObjects(t, fixture, fs)
			verifyCommitGraphFiles(t, fixture, fs)
		})
	}
}

func verifyCommitGraphObjects(t *testing.T, fixture *fixtures.Fixture, fs billy.Filesystem) {
	t.Helper()

	count := 3
	if fixture.Is("commit-graph-chain-sha1-257") {
		count = 257
	}

	entries := fixture.CommitGraphEntries()
	require.Len(t, entries, count)
	require.EqualValues(t, count+1, fixture.ObjectsCount)
	require.Equal(t, fixture.Head+"\n", string(readGraphFixtureFile(t, fs, "refs/heads/main")))
	require.Equal(t, "ref: refs/heads/main\n", string(readGraphFixtureFile(t, fs, "HEAD")))
	require.Equal(t, fixture.Head, entries[count-1].Hash)

	for i, entry := range entries {
		parents := ""
		if i == 0 {
			require.Empty(t, entry.Parents)
		} else {
			require.Equal(t, []string{entries[i-1].Hash}, entry.Parents)
			parents = "parent " + entries[i-1].Hash + "\n"
		}

		commit := readGraphFixtureObject(t, fs, entry.Hash, fixture.ObjectFormat)
		body := fmt.Sprintf("tree %s\n%sauthor Test <test@example.com> %d +0000\n"+
			"committer Test <test@example.com> %d +0000\n\ncommit %d\n",
			entry.Tree, parents, 1700000000+i, 1700000000+i, i)
		require.Equal(t, fmt.Sprintf("commit %d\x00%s", len(body), body), string(commit))
		tree := readGraphFixtureObject(t, fs, entry.Tree, fixture.ObjectFormat)
		require.Equal(t, "tree 0\x00", string(tree))
	}
}

func readGraphFixtureFile(t *testing.T, fs billy.Filesystem, name string) []byte {
	t.Helper()

	file, err := fs.Open(name)
	require.NoError(t, err)
	content, err := io.ReadAll(file)
	require.NoError(t, err)
	require.NoError(t, file.Close())

	return content
}

func readGraphFixtureObject(t *testing.T, fs billy.Filesystem, id, format string) []byte {
	t.Helper()

	file, err := fs.Open("objects/" + id[:2] + "/" + id[2:])
	require.NoError(t, err)
	reader, err := zlib.NewReader(file)
	require.NoError(t, err)
	content, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.NoError(t, reader.Close())
	require.NoError(t, file.Close())
	require.Equal(t, id, graphFixtureHash(content, format))

	return content
}

func graphFixtureHash(content []byte, format string) string {
	if format == "sha256" {
		sum := sha256.Sum256(content)

		return hex.EncodeToString(sum[:])
	}

	sum := sha1.Sum(content) //nolint:gosec // Verify Git SHA-1 fixture objects and graphs.

	return hex.EncodeToString(sum[:])
}

func verifyCommitGraphFiles(t *testing.T, fixture *fixtures.Fixture, fs billy.Filesystem) {
	t.Helper()

	paths := []string{"objects/info/commit-graph"}
	if fixture.Is("commit-graph-split") {
		chain := readGraphFixtureFile(t, fs, "objects/info/commit-graphs/commit-graph-chain")
		ids := strings.Fields(string(chain))
		require.Len(t, ids, len(fixture.CommitGraphEntries()))
		paths = make([]string, 0, len(ids))
		for _, id := range ids {
			paths = append(paths, "objects/info/commit-graphs/graph-"+id+".graph")
		}
		_, err := fs.Stat("objects/info/commit-graph")
		require.Error(t, err)
	} else {
		_, err := fs.Stat("objects/info/commit-graphs/commit-graph-chain")
		require.Error(t, err)
	}

	hashSize, hashVersion := 20, byte(1)
	if fixture.ObjectFormat == "sha256" {
		hashSize, hashVersion = 32, 2
	}

	var commits []string
	for layer, path := range paths {
		graph := readGraphFixtureFile(t, fs, path)
		require.Equal(t, "CGPH", string(graph[:4]))
		require.Equal(t, byte(1), graph[4])
		require.Equal(t, hashVersion, graph[5])
		require.EqualValues(t, layer%256, graph[7])
		require.Equal(t, hex.EncodeToString(graph[len(graph)-hashSize:]),
			graphFixtureHash(graph[:len(graph)-hashSize], fixture.ObjectFormat))

		var baseSize int
		for chunk := range int(graph[6]) {
			position := 8 + chunk*12
			start := binary.BigEndian.Uint64(graph[position+4 : position+12])
			end := binary.BigEndian.Uint64(graph[position+16 : position+24])
			require.LessOrEqual(t, start, end)
			require.LessOrEqual(t, end, uint64(len(graph)-hashSize)) //nolint:gosec // Graph length exceeds its trailer.
			switch string(graph[position : position+4]) {
			case "BASE":
				baseSize = len(graph[start:end])
			case "OIDL":
				for offset := start; offset < end; offset += uint64(hashSize) {
					commits = append(commits, hex.EncodeToString(graph[offset:offset+uint64(hashSize)]))
				}
			}
		}
		require.Equal(t, layer*hashSize, baseSize)
	}

	expected := make([]string, 0, len(fixture.CommitGraphEntries()))
	for _, entry := range fixture.CommitGraphEntries() {
		expected = append(expected, entry.Hash)
	}
	require.ElementsMatch(t, expected, commits)
}
