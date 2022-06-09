package slicepackage

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Package struct {
	Name         string   `json:"name"`
	Version      string   `json:"version"`
	Architecture string   `json:"architecture"`
	Maintainer   string   `json:"maintainer"`
	Dependencies []string `json:"dependencies"`
	Description  string   `json:"description"`
}

func CompressData(data []byte) []byte {
	buffer := new(bytes.Buffer)

	writer, _ := gzip.NewWriterLevel(buffer, gzip.BestCompression)

	writer.Write(data)

	writer.Close()

	return buffer.Bytes()
}

func DecompressData(data []byte) []byte {
	reader, _ := gzip.NewReader(bytes.NewReader(data))

	decompressedData, _ := ioutil.ReadAll(reader)

	return decompressedData
}

func CreatePackageTarball(source string) ([]byte, error) {
	buffer := new(bytes.Buffer)

	tarball := tar.NewWriter(buffer)

	defer tarball.Close()

	_, err := os.Stat(source)

	if err != nil {
		return nil, nil
	}

	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
	
		var link string
		
		if info.Mode() & os.ModeSymlink == os.ModeSymlink {
			if link, err = os.Readlink(path); err != nil {
				return err
			}
		}
	
		header, err := tar.FileInfoHeader(info, link)

		if err != nil {
			return err
		}
	
		header.Name = strings.TrimPrefix(path, source)

		if len(header.Name) > 0 {
			header.Name = header.Name[1:]
		}

		if err = tarball.WriteHeader(header); err != nil {
			return err
		}
	
		if !info.Mode().IsRegular() {
			return nil
		}
	
		file, err := os.Open(path)

		if err != nil {
			return err
		}

		defer file.Close()

		_, err = io.Copy(tarball, file)

		return err
	})

	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
