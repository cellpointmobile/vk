// Copyright © 2019 Anders Bruun Olsen <anders@bruun-olsen.net>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package file

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// ExtractFromTar extracts a file from a tarball
func ExtractFromTar(source string, target string, destination string) error {
	resp, err := http.Get(source)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var in io.Reader

	if strings.HasSuffix(source, "gz") {
		in, err = gzip.NewReader(resp.Body)
	} else if strings.HasSuffix(source, "bz2") {
		in = bzip2.NewReader(resp.Body)
	}
	if err != nil {
		return err
	}

	tr := tar.NewReader(in)
	for {
		header, err := tr.Next()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}
		if header.Name == target {
			out, err := os.Create(destination)
			if err != nil {
				return err
			}
			defer out.Close()
			_, err = io.Copy(out, tr)
			if err != nil {
				return err
			}
		}
	}
}

// ExtractFromZip extracts a file from a zip-file
func ExtractFromZip(source string, target string, destination string) error {
	resp, err := http.Get(source)
	if err != nil {
		return err
	}
	contentZipped, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	readerAt := bytes.NewReader(contentZipped)
	zr, err := zip.NewReader(readerAt, int64(len(contentZipped)))
	if err != nil {
		return err
	}
	for _, f := range zr.File {
		if f.Name == target {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()
			outFile, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}

			_, err = io.Copy(outFile, rc)
			// Close the file without defer to close before next iteration of loop
			outFile.Close()
			if err != nil {
				return err
			}
		}

	}
	return nil
}
