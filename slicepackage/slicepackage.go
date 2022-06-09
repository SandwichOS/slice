package slicepackage

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
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

		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
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

func GetPackageMetadata(source []byte) (Package, error) {
	buffer := new(bytes.Buffer)
	buffer.Write(source)

	tarball := tar.NewReader(buffer)

	metadataBuffer := new(bytes.Buffer)

	for {
		header, err := tarball.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return Package{"", "", "", "", []string{}, ""}, err
		}

		if header.Name != "metadata.json" {
			continue
		}

		if _, err := io.Copy(metadataBuffer, tarball); err != nil {
			return Package{"", "", "", "", []string{}, ""}, err
		}

		break
	}

	var packageMetadata Package

	err := json.Unmarshal(metadataBuffer.Bytes(), &packageMetadata)

	if err != nil {
		return Package{"", "", "", "", []string{}, ""}, err
	}

	return packageMetadata, nil
}

func ExtractPackageTarball(source []byte, target string) error {
	buffer := new(bytes.Buffer)
	buffer.Write(source)

	tarReader := tar.NewReader(buffer)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path := filepath.Join(target, header.Name)
		info := header.FileInfo()

		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}

			continue
		}

		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			os.Symlink(filepath.Join(filepath.Dir(path), header.Linkname), path)
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())

		if err != nil {
			return err
		}

		defer file.Close()
		_, err = io.Copy(file, tarReader)

		if err != nil {
			return err
		}
	}

	return nil
}
