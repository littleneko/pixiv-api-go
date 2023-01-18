package pixiv_api_go

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
)

func CheckAndMkdir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

type HashType int

const (
	NoHash HashType = -1
	Sha1   HashType = iota
	Sha256
)

func WriteFileAndCalHash(reader io.Reader, filename string, hashType HashType) (int64, string, error) {
	dirName := filepath.Dir(filename)
	err := CheckAndMkdir(dirName)
	if err != nil {
		return 0, "", err
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return 0, "", err
	}
	defer func() {
		_ = file.Close()
	}()

	if hashType == NoHash {
		size, err := io.Copy(file, reader)
		return size, "", err
	}

	var h hash.Hash
	switch hashType {
	case Sha256:
		h = sha256.New()
		break
	case Sha1:
	default:
		h = sha1.New()
		break
	}
	r := io.TeeReader(reader, h)

	size, err := io.Copy(file, r)
	sum := fmt.Sprintf("%x", h.Sum(nil))
	return size, sum, nil
}

func WriteFile(reader io.Reader, filename string) (int64, error) {
	size, _, err := WriteFileAndCalHash(reader, filename, NoHash)
	return size, err
}

func WriteFIleCalSha1(reader io.Reader, filename string) (int64, string, error) {
	return WriteFileAndCalHash(reader, filename, Sha1)
}
