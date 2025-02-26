package tgz

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"

	"github.com/go-git/go-billy/v6"
	"github.com/go-git/go-billy/v6/memfs"
)

//nolint:gochecknoglobals
var MemFactory = func() (billy.Filesystem, error) {
	return memfs.New(), nil
}

// Extract decompress a gziped tarball into the fs billy.Filesystem.
//
// A non-nil error is returned if the method fails to complete.
func Extract(tgz billy.File, fs billy.Filesystem) (err error) {
	defer func() {
		errClose := tgz.Close()
		if err == nil {
			err = errClose
		}
	}()

	tar, err := zipTarReader(tgz)
	if err != nil {
		return
	}

	err = unTar(fs, tar)
	if err != nil {
		return
	}

	return
}

func zipTarReader(r io.Reader) (*tar.Reader, error) {
	zip, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return tar.NewReader(zip), nil
}

func filemode(mode int64) (fs.FileMode, error) {
	if mode < 0 {
		return 0, errors.New("mode cannot be negative")
	}
	if mode > math.MaxUint32 {
		return 0, errors.New("mode cannot be greater than max uint32")
	}

	return os.FileMode(mode), nil
}

func unTar(fs billy.Filesystem, src *tar.Reader) error {
	for {
		header, err := src.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return err
		}

		dst := header.Name
		mode, err := filemode(header.Mode)
		if err != nil {
			return err
		}
		switch header.Typeflag {
		case tar.TypeDir:
			err := fs.MkdirAll(dst, mode)
			if err != nil {
				return err
			}
		case tar.TypeReg:
			err := makeFile(fs, dst, mode, src)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unable to untar type: %c in file %s",
				header.Typeflag, header.Name)
		}
	}

	return nil
}

func makeFile(fs billy.Filesystem, path string, mode os.FileMode, contents io.Reader) (err error) {
	w, err := fs.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		errClose := w.Close()
		if err == nil {
			err = errClose
		}
	}()

	_, err = io.Copy(w, contents)
	if err != nil {
		return err
	}

	if fs, ok := fs.(billy.Change); ok {
		if err = fs.Chmod(path, mode); err != nil {
			return err
		}
	}

	return nil
}
