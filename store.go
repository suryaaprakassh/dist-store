package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

// any function that gives the path of the file
// according to the key
type PathTransformFunc func(string) PathKey

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])
	blockSize := 5
	sliceLen := len(hashStr) / blockSize
	paths := make([]string, sliceLen)
	for i := 0; i < sliceLen; i++ {
		from := i * blockSize
		//out of bound so min of len and calculated end
		to := min(from+blockSize, len(hashStr))
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		Path:     strings.Join(paths, "/"),
		FileName: hashStr,
	}
}

func DefaultPathTransform(key string) PathKey {
	return PathKey{
		Path:     "",
		FileName: key,
	}
}

type PathKey struct {
	Path     string
	FileName string
}

func (p *PathKey) fileNameWithPath() string {
	return fmt.Sprintf("%s/%s", p.Path, p.FileName)
}

func (p *PathKey) firstParentDir() string {
	paths := strings.Split(p.Path, "/")

	//TODO: handle the thing gracefully!
	if len(paths) < 1 {
		panic("Path string is not valid")
	}
	return paths[0]
}

type StorageOpts struct {
	PathTransformFunc
}

// used to store the stream data to the disk
// can be a db not so sure
type Store struct {
	StorageOpts
}

func NewStore(opts StorageOpts) *Store {
	return &Store{
		StorageOpts: opts,
	}
}

func (s *Store) Has(key string) bool {
	pathKey := s.PathTransformFunc(key)
	_, err := os.Stat(pathKey.fileNameWithPath())
	if err != nil {
		return false
	}
	return true
}

func (s *Store) Delete(key string) error {
	pathKey := s.PathTransformFunc(key)

	//NOTE: directories will also get removed
	//check for any possible side effects
	return os.RemoveAll(pathKey.firstParentDir())
}

func (s *Store) Write(key string, r io.Reader) error {
	pathName := s.PathTransformFunc(key)

	//TODO: check os permissions
	if err := os.MkdirAll(pathName.Path, os.ModePerm); err != nil {
		return err
	}

	//file creation
	fileNameWithPath := pathName.fileNameWithPath()
	f, err := os.Create(fileNameWithPath)

	if err != nil {
		return err
	}

	//copying the contents to the disk
	_, err = io.Copy(f, r)
	if err != nil {
		return err
	}

	slog.Info("Created New File", "key", key)
	return nil
}

// function returns a reader with the contents of the file
// If and only if file exists
func (s *Store) Read(key string) (io.Reader, error) {
	pathKey := s.PathTransformFunc(key)
	filePath := pathKey.fileNameWithPath()
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf := new(bytes.Buffer)

	_, err = io.Copy(buf, file)

	if err != nil {
		return nil, err
	}

	return buf, nil
}
