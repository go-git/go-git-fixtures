package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// Doctor a MIDX bitmap into a pack bitmap.
// Usage: doctor_bitmap <midx-bitmap-in> <pack-file> <pack-bitmap-out>
func main() {
	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "usage: doctor_bitmap <midx-bitmap> <pack> <out-bitmap>")
		os.Exit(2)
	}

	midxBm, err := os.ReadFile(os.Args[1])
	check(err)

	pf, err := os.Open(os.Args[2])
	check(err)
	defer pf.Close()
	info, err := pf.Stat()
	check(err)

	packTrailer := make([]byte, 32)
	_, err = pf.ReadAt(packTrailer, info.Size()-32)
	check(err)

	if len(midxBm) < 12+32+32 {
		check(fmt.Errorf("midx bitmap too small: %d", len(midxBm)))
	}

	body := make([]byte, len(midxBm))
	copy(body, midxBm)
	copy(body[12:44], packTrailer)

	h := sha256.New()
	_, err = h.Write(body[:len(body)-32])
	check(err)
	copy(body[len(body)-32:], h.Sum(nil))

	out, err := os.Create(os.Args[3])
	check(err)
	defer out.Close()

	_, err = io.Copy(out, readerOf(body))
	check(err)
	fmt.Printf("wrote %s (%d bytes)\n", os.Args[3], len(body))
}

func readerOf(b []byte) io.Reader { return &sliceReader{b: b} }

type sliceReader struct {
	b   []byte
	off int
}

func (r *sliceReader) Read(p []byte) (int, error) {
	if r.off >= len(r.b) {
		return 0, io.EOF
	}
	n := copy(p, r.b[r.off:])
	r.off += n
	return n, nil
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
