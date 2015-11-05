// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package file

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/nanobox-io/nanobox/config"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Copy
func Copy(dst, src string) error {

	// ensure src exists
	sfi, err := os.Stat(src)
	if err != nil {
		return err
	}

	// create dest dir
	if err = os.MkdirAll(dst, sfi.Mode()); err != nil {
		return err
	}

	//
	return copyDir(src, dst)
}

// Tar
func Tar(path string, writers ...io.Writer) error {

	//
	mw := io.MultiWriter(writers...)

	//
	gzw := gzip.NewWriter(mw)
	defer gzw.Close()

	//
	tw := tar.NewWriter(gzw)
	defer tw.Close()

	//
	return filepath.Walk(path, func(file string, fi os.FileInfo, err error) error {

		//
		if err != nil {
			return err
		}

		// only tar files (not dirs)
		if fi.Mode().IsRegular() {

			// create header for this file
			header := &tar.Header{
				Name: file,
				Mode: int64(fi.Mode()),
				Size: fi.Size(),
				// ModTime:  fi.ModTime(),
				Typeflag: tar.TypeReg,
			}

			// write the header to the tarball archive
			if err := tw.WriteHeader(header); err != nil {
				return err
			}

			// open the file for taring...
			f, err := os.Open(file)
			defer f.Close()
			if err != nil {
				return err
			}

			// copy from file data into tar writer
			if _, err := io.Copy(tw, f); err != nil {
				return err
			}
		}

		return nil
	})
}

// Untar
func Untar(dst string, r io.Reader) error {

	//
	gzr, err := gzip.NewReader(r)
	defer gzr.Close()
	if err != nil {
		return err
	}

	//
	tr := tar.NewReader(gzr)

	//
	for {
		header, err := tr.Next()

		//
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		}

		//
		switch header.Typeflag {

		// if its a dir, make it
		case tar.TypeDir:
			if err := os.MkdirAll(header.Name, os.FileMode(header.Mode)); err != nil {
				return err
			}

		// if its a file, add it to the dir
		case tar.TypeReg:
			f, err := os.Create(header.Name)
			defer f.Close()
			if err != nil {
				return err
			}

			// copy from tar reader into file
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

		//
		default:
			return fmt.Errorf("Unhandled type (%v) for %v", header.Typeflag, header.Name)
		}
	}

}

// Download
func Download(path string, w io.Writer) error {
	res, err := http.Get(path)
	defer res.Body.Close()
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		config.Fatal("[util/file/file] ioutil.ReadAll() failed - ", err.Error())
	}

	w.Write(b)

	return nil
}

// Progress
func Progress(path string, w io.Writer) error {

	//
	download, err := http.Get(path)
	defer download.Body.Close()
	if err != nil {
		return err
	}

	var percent float64
	var down int

	// format the response content length to be more 'friendly'
	total := float64(download.ContentLength) / math.Pow(1024, 2)

	// create a 'buffer' to read into
	p := make([]byte, 2048)

	//
	for {

		// read the response body (streaming)
		n, err := download.Body.Read(p)

		// write to our buffer
		w.Write(p[:n])

		// update the total bytes read
		down += n

		// update the percent downloaded
		percent = (float64(down) / float64(download.ContentLength)) * 100

		// show download progress: down/totalMB [*** progress *** %]
		fmt.Printf("\r   %.2f/%.2fMB [%-41s %.2f%%]", float64(down)/math.Pow(1024, 2), total, strings.Repeat("*", int(percent/2.5)), percent)

		// detect EOF and break the 'stream'
		if err != nil {
			if err == io.EOF {
				fmt.Println("")
				break
			} else {
				return err
			}
		}
	}

	return nil
}

// copyFile
func copyFile(src, dst string) error {

	sf, err := os.Open(src)
	defer sf.Close()
	if err != nil {
		return err
	}

	df, err := os.Create(dst)
	defer df.Close()
	if err != nil {
		return err
	}

	if _, err := io.Copy(df, sf); err != nil {
		return err
	}

	fi, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, fi.Mode())
}

// copyDir
func copyDir(src, dst string) error {

	// get properties of source dir
	fi, err := os.Stat(src)
	if err != nil {
		return err
	}

	// create dest dir
	err = os.MkdirAll(dst, fi.Mode())
	if err != nil {
		return err
	}

	//
	dir, err := os.Open(src)
	if err != nil {
		return err
	}

	//
	fis, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	for _, fi := range fis {

		srcPath := filepath.Join(src, fi.Name())
		dstPath := filepath.Join(dst, fi.Name())

		switch {

		// create sub-directories - recursively
		case fi.Mode().IsDir():
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}

		// perform copy
		case fi.Mode().IsRegular():
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}
