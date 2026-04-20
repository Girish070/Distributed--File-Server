package store

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"
)

type PathKey struct {
	PathName string
	FileName string
}

// FristPathName returns the root folder of this path (Useful for cleanUo later)
func (p *PathKey) FirstPathName() string {
	paths := strings.Split(p.PathName, "/")
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

// FullPath returns the complete path incuding the fileName
func (p *PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.PathName, p.FileName)
}

// PathTransFromFunc is blueprint for functions that decide where files gp
type PathTransformFunc func(string) PathKey

// DeFaultPathTranformFunc places everything in the root directory  (Not recommended for prod)
var DefaultPathTransfromFunc = func(key string) PathKey {
	return PathKey{
		PathName: key,
		FileName: key,
	}
}

// CASPathtransformFunc creates a deeply nested structure using SHA-1 hashing
func CASPathTransfromFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hashStr) / blockSize

	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashStr[from:to]
	}
	return PathKey{
		PathName: strings.Join(paths, "/"),
		FileName: hashStr,
	}
}

// StoreOpts allows to configure where and how files are saved
type StoreOpts struct {
	Root              string //The main folder where everything lives
	PathTransformFunc PathTransformFunc
}

// DefaultStoreOpts provide sensible defaults if we don't specify them
var DefaultStoreOpts = StoreOpts{
	Root:              "network_data",
	PathTransformFunc: CASPathTransfromFunc,
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	//If no transform function provided, use the default (CAS)
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultStoreOpts.PathTransformFunc
	}

	//If no root is provided use the default
	if len(opts.Root) == 0 {
		opts.Root = DefaultStoreOpts.Root
	}

	return &Store{
		StoreOpts: opts,
	}
}

// WriteStream reads from the io.Reader and writes the bytes to disk
func (s *Store) WriteStream(key string, r io.Reader) error {
	pathkey := s.PathTransformFunc(key)

	//Combine the Root folder with the generated hash path
	pathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathkey.PathName)

	//1. Create the nested folders! os.MakedirAll is like "mkdir -p" in linux
	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
		return err
	}

	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathkey.FullPath())

	//2. Create the actual file
	f, err := os.Create(fullPathWithRoot)
	if err != nil {
		return err
	}
	defer f.Close()

	//3. Stream the bytes from the reader directly to the hard drive
	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("Written (%d) bytes to the disk: %s", n, fullPathWithRoot)
	return nil
}

// ReadStream locates the file by its key and returns a stream to read it
// Note: it returns an io.ReadCloser so the caller is responsible for calling close() when finished!
func (s *Store) ReadStream(key string) (io.ReadCloser, error) {
	//1 Calculate the exact path using CAS function
	pathKey := s.PathTransformFunc(key)

	//2 Combine it with root directory
	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())

	// Open the file and return Both the stream and the error
	return os.Open(fullPathWithRoot)
}

// Delete completely removes the file and its nested CAS directories from the disk
func (s *Store) Delete(key string) error {
	pathKey := s.PathTransformFunc(key)

	//Delete from the root of the hash to clean up all empty folders
	firstPathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FirstPathName())

	err := os.RemoveAll(firstPathWithRoot)
	if err != nil {
		return err
	}
	log.Printf("Deleted file and cleaned up directories: %s", firstPathWithRoot)
	return nil
}

// Has checks if a file already exists in the local CAS storage 
func (s *Store) Has(key string) bool {
	pathKey := s.PathTransformFunc(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())

	_, err := os.Stat(fullPathWithRoot)
	return !errors.Is(err, fs.ErrNotExist)
}
