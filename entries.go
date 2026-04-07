package fixtures

// PackfileEntry maps an object hash (hex-encoded) to its byte offset in the packfile.
type PackfileEntry = map[string]int64

type packfileData struct {
	hashes  []string
	offsets []int64
}

// Entries returns the expected object entries for this fixture's packfile.
// Each entry maps an object hash (hex-encoded) to its byte offset in the packfile.
// Returns nil if no entries are registered for this fixture's packfile.
func (f *Fixture) Entries() PackfileEntry {
	d, ok := packfileEntries[f.PackfileHash]
	if !ok {
		return nil
	}

	m := make(PackfileEntry, len(d.hashes))
	for i, h := range d.hashes {
		m[h] = d.offsets[i]
	}

	return m
}

// ScannerEntry represents a single object as read from the packfile stream.
// Fields use plain types so go-git-fixtures does not depend on go-git.
//
// Object type constants match plumbing.ObjectType:
// 1=commit, 2=tree, 3=blob, 4=tag, 6=ofs-delta, 7=ref-delta.
type ScannerEntry struct {
	Type            int
	Offset          int64
	Size            int64
	Hash            string // hex-encoded; empty for delta objects
	Reference       string // hex-encoded; base object hash for ref-delta
	OffsetReference int64  // base object offset for ofs-delta
	CRC32           uint32
}

// ScannerEntries returns the expected scanner output for this fixture's packfile,
// ordered by pack stream offset. Returns nil if not registered.
func (f *Fixture) ScannerEntries() []ScannerEntry {
	entries, ok := scannerEntries[f.PackfileHash]
	if !ok {
		return nil
	}

	out := make([]ScannerEntry, len(entries))
	copy(out, entries)

	return out
}

// basicSHA1Hashes are the 31 object hashes in basic.git (SHA1 format), sorted lexicographically.
// Shared across ofs-delta and ref-delta packfiles which contain the same objects at different offsets.
//
//nolint:gochecknoglobals
var basicSHA1Hashes = [...]string{
	"1669dce138d9b841a518c64b10914d88f5e488ea",
	"32858aad3c383ed1ff0a0f9bdf231d54a00c9e88",
	"35e85108805c84807bc66a02d91535e1e24b38b9",
	"49c6bb89b17060d7b4deacb7b338fcc6ea2352a9",
	"4d081c50e250fa32ea8b1313cf8bb7c2ad7627fd",
	"586af567d0bb5e771e49bdd9434f5e0fb76d25fa",
	"5a877e6a906a2743ad6e45d99c1793642aaf8eda",
	"6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
	"7e59600739c96546163833214c36459e324bad0a",
	"880cd14280f4b9b6ed3986d6671f907d7cc2a198",
	"8dcef98b1d52143e1e2dbc458ffe38f925786bf2",
	"918c48b83bd081e863dbe1b80f8998f058cd8294",
	"9a48f23120e880dfbe41f7c9b7b708e9ee62a492",
	"9dea2395f5403188298c1dabe8bdafe562c491e3",
	"a39771a7651f97faf5c72e08224d857fc35133db",
	"a5b8b09e2f8fcb0bb99d3ccb0958157b40890d69",
	"a8d315b2b1c615d43042c3a62402b8a54288cf5c",
	"aa9b383c260e1d05fbbf6b30a02914555e20c725",
	"af2d6a6954d532f8ffb47615169c8fdf9d383a1a",
	"b029517f6300c2da0f4b651b8642506cd6aaf45d",
	"b8e471f58bcbca63b07bda20e428190409c2db47",
	"c192bd6a24ea1ab01d78686e417c8bdc7c3d197f",
	"c2d30fa8ef288618f65f6eed6e168e0d514886f4",
	"c8f1d8c61f9da76f4cb49fd86322b6e685dba956",
	"cf4aa3b38974fb7d81f367c0830f7d78d65ab86b",
	"d3ff53e0564a9f87d8e84b6e28e5060e517008aa",
	"d5c0f4ab811897cadf03aec358ae60d21f91c50d",
	"dbd3641b371024f44d0e469a9c8f5457b0660de1",
	"e8d3ffab552895c19b9fcf7aa264d277cde33881",
	"eba74343e2f15d62adedfd8c883ee0262b5c8021",
	"fb72698cab7617ac416264415f13224dfd7a165e",
}

// basicSHA256Hashes are the 36 object hashes in basic.git (SHA256 format), sorted lexicographically.
//
//nolint:gochecknoglobals
var basicSHA256Hashes = [...]string{
	"011218223f6e9e4a7f7ed704999158d6a3d080bedff536983c0d0e03d262c664",
	"030d8320428f364839a75c1fe8d4cc2cdada2b683dcaffcc94d9770640302dd1",
	"176d63c1aa704b4021d82cc75c6a8a7bbd96c7b30d774d7e629773581cfd4501",
	"1e7242fb7dfbf84896c05ee1f2fde2d591103cc5f6e5b9c7f8562b51e9e1732b",
	"1f307724f91af43be1570b77aeef69c5010e8136e50bef83c28de2918a08f494",
	"23e6fd7f90c8b0d3dddd48d8a87a00662972b8a917d4d7fc2d0967921b65f3ed",
	"2849f40d9cd298ce2a85d6dc603e84c99e6c6bcbf798740b57bc7deaaa913360",
	"2a246d3eaea67b7c4ac36d96d1dc9dad2a4dc24486c4d67eb7cb73963f522481",
	"2a7543a59f760f7ca41784bc898057799ae960323733cab1175c21960a750f72",
	"2ad4c66a3680b32a04547b749c105edb5421ddffd4fac791fd24368b23c7ffc6",
	"33a5013ed4af64b6e54076c986a4733c2c11ce8ab27ede79f21366e8722ac5ed",
	"38ad2967b54c80797487d45a5db951406d72927580faeb224a678576f962bcef",
	"40b7c05726c9da78c3d5a705c2a48a120261b36f521302ce06bad41916d000f7",
	"4c61794e77ff8c7ab7f07404cdb1bc0e989b27530e37a6be6d2ef73639aaff6d",
	"4c75fc7e8111ca9dc7f087c85f1a6f969f6e394034acd203d4f07562e668bb95",
	"4fef4adac3be863b9b94613016bdd8e53f67f6d7577234e028bc9d24c5a6a27c",
	"53a77fb45346d765b3d7054ab6ed3e7a227cb3b81ac2c60c83d0a8308a15264c",
	"5dd3e66d32270068b4ed56cedc1b82b9b39e2dde6df9aa724092879a4cddad6b",
	"65bb8b5ad068a89499ce27b1e0397fb4c027c013d7c407671bb8c70777f78e13",
	"665e33431d9b88280d7c1837680fdb66664c4cb4b394c9057cdbd07f3b4acff8",
	"6e8d71fbfd367c34968d31ef8886929a9862b02de4616bfc569583b3f5a76808",
	"73660d98a4c6c8951f86bb8c4744a0b4837a6dd5f796c314064c1615781c400c",
	"789c9f4220d167b66020b46bacddcad0ab5bb12f0f469576aa60bb59d98293dc",
	"80d53c7b7196c44b0abd4d102772dedeb33069b617e5df2f0becc2563a37e1b0",
	"8cc70e96f2ee81cdad77361933640703a42ee3a04fade68578e836714f535d76",
	"9768a9bcb42f35dc598a517bd98a5cbba79052b980a8a015f3be5577ebd9f201",
	"abb8b2cb27cba10e236e97e06e4a5a0acd6b50e89e8a6e974439a39cb4c3de00",
	"ac16b517cae0a031a218f0edb988ae0df4ee267a531a8a35f17dba1787ea0422",
	"b8bdc620cb4859cf6e48768fd67f526229f3a57aa417740024bf7e6af5fdb04c",
	"c74a1ff56ec2c88a7e214436a560e30c0c4b699e92cb62449df487d4707bea3d",
	"cbaa8eafbf007764f1ef3681261384359976a9edff18e906cbd6802fcdec6f75",
	"e6ee53c7eb0e33417ee04110b84b304ff2da5c1b856f320b61ad9f2ef56c6e4e",
	"e725c2efbb1bb3e5ff39d5b1cb6c38e33c7f294259974b367c157d811425776a",
	"ee4e96e4a1684b5ad691c752be98c517bb4f71fbbef6c35e743c4accdbc1f231",
	"ef36d9a576158df19554d50c9180d503ded2b86a85956d3da9bf1369449f34ab",
	"fa60c322a88283ab1e9d872f4782eb4f4da7f98179e574ba85f58b992d918d6a",
}

// packfileEntries maps packfile hashes to their object data.
// Hashes are shared across packfiles of the same repository and format;
// only the offsets differ between delta encodings.
//
//nolint:gochecknoglobals
var packfileEntries = map[string]packfileData{
	// basic.git ofs-delta (sha1)
	"a3fed42da1e8189a077c0e6846c040dcf73fc9dd": {
		hashes: basicSHA1Hashes[:],
		offsets: []int64{
			615, 1524, 1063, 78882, 84688, 84559, 84479, 186,
			84653, 78050, 84741, 286, 80998, 84032, 84430, 838,
			84375, 84760, 449, 1392, 1230, 1713, 84725, 80725,
			84608, 1685, 2351, 84115, 12, 84708, 84671,
		},
	},
	// basic.git ref-delta (sha1)
	"c544593473465e6315ad4182d04d366c4592b829": {
		hashes: basicSHA1Hashes[:],
		offsets: []int64{
			633, 1542, 1243, 79129, 85262, 81265, 79049, 186,
			85244, 78117, 85448, 304, 81314, 84797, 78068, 856,
			84880, 85485, 467, 1410, 1081, 1731, 85335, 80972,
			84752, 1703, 2369, 85176, 12, 85300, 85141,
		},
	},
	// basic.git (sha256)
	"c88dfe1663bd216e278d5bb3c8decd0a4bb174a6204585dc44b7c7a05fceed55": {
		hashes: basicSHA256Hashes[:],
		offsets: []int64{
			299, 1927, 85787, 85748, 85711, 85473, 1199, 82044,
			85540, 85483, 79369, 1464, 2674, 80201, 85417, 411,
			2267, 85805, 85826, 3501, 12, 82383, 2863, 85729,
			810, 2120, 79304, 79200, 608, 1732, 82317, 2835,
			1003, 85623, 85646, 85769,
		},
	},
}
