package fixtures_test

import (
	"testing"

	fixtures "github.com/go-git/go-git-fixtures/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntries(t *testing.T) {
	t.Parallel()

	for _, f := range fixtures.ByTag("packfile-entries") {
		t.Run(f.PackfileHash, func(t *testing.T) {
			t.Parallel()

			entries := f.Entries()
			require.NotNil(t, entries)
			assert.NotEmpty(t, entries)

			for hash, offset := range entries {
				assert.NotEmpty(t, hash)
				assert.GreaterOrEqual(t, offset, int64(0))
			}
		})
	}
}

func TestEntriesReturnsNilForUnregisteredPackfile(t *testing.T) {
	t.Parallel()

	f := fixtures.ByTag("notes").One()
	require.NotNil(t, f)
	assert.Nil(t, f.Entries())
}

func TestEntriesReturnsFreshMap(t *testing.T) {
	t.Parallel()

	f := fixtures.ByTag("packfile-entries").One()
	require.NotNil(t, f)

	a := f.Entries()
	b := f.Entries()

	require.NotNil(t, a)
	require.NotNil(t, b)

	assert.Equal(t, a, b)

	// Mutating one must not affect the other.
	for k := range a {
		a[k] = -1

		break
	}

	assert.NotEqual(t, a, b)
}

func TestScannerEntries(t *testing.T) {
	t.Parallel()

	for _, f := range fixtures.ByTag("scanner-entries") {
		t.Run(f.PackfileHash, func(t *testing.T) {
			t.Parallel()

			entries := f.ScannerEntries()
			require.NotNil(t, entries)
			assert.NotEmpty(t, entries)

			var prevOffset int64
			for i, e := range entries {
				assert.Greater(t, e.Offset, prevOffset,
					"entry %d: offset %d must be greater than previous %d", i, e.Offset, prevOffset)
				prevOffset = e.Offset

				assert.GreaterOrEqual(t, e.Size, int64(0))
				assert.NotZero(t, e.Type)
				assert.NotZero(t, e.CRC32)
			}
		})
	}
}

func TestScannerEntriesReturnsNilForUnregisteredPackfile(t *testing.T) {
	t.Parallel()

	f := fixtures.ByTag("notes").One()
	require.NotNil(t, f)
	assert.Nil(t, f.ScannerEntries())
}

func TestScannerEntriesReturnsFreshSlice(t *testing.T) {
	t.Parallel()

	f := fixtures.ByTag("scanner-entries").One()
	require.NotNil(t, f)

	a := f.ScannerEntries()
	b := f.ScannerEntries()

	require.NotNil(t, a)
	require.NotNil(t, b)

	assert.Equal(t, a, b)

	a[0].Offset = -1
	assert.NotEqual(t, a, b)
}

func TestEntriesConsistentWithScannerEntries(t *testing.T) {
	t.Parallel()

	// Fixtures tagged with both should have consistent hash/offset pairs.
	ff := fixtures.ByTag("packfile-entries").ByTag("scanner-entries")
	require.NotEmpty(t, ff)

	for _, f := range ff {
		t.Run(f.PackfileHash, func(t *testing.T) {
			t.Parallel()

			entries := f.Entries()
			scannerEntries := f.ScannerEntries()

			require.NotNil(t, entries)
			require.NotNil(t, scannerEntries)

			for _, se := range scannerEntries {
				if se.Hash == "" {
					continue // delta objects have no hash
				}

				offset, ok := entries[se.Hash]
				require.True(t, ok, "scanner entry hash %s not found in packfile entries", se.Hash)
				assert.Equal(t, se.Offset, offset,
					"offset mismatch for hash %s", se.Hash)
			}
		})
	}
}
