package store

import (
	"bytes"
	"io"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "my_best_picture.png"
	pathKey := CASPathTransfromFunc(key)

	expectedFilename := "e09bba016f88f2b92e824d9a4142a9c81c6440ae"
	expectedPathName := "e09bb/a016f/88f2b/92e82/4d9a4/142a9/c81c6/440ae"

	if pathKey.PathName != expectedPathName {
		t.Errorf("have %s want %s", pathKey.PathName, expectedPathName)
	}

	if pathKey.FileName != expectedFilename {
		t.Errorf("have %s want %s", pathKey.FileName, expectedFilename)
	}
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		Root:              "my_test_network",
		PathTransformFunc: CASPathTransfromFunc,
	}
	s := NewStore(opts)

	key := "my_special_picture.png"
	data := []byte("These are some fake image bytes. Pretend this is a cool PNG!")

	sourceStream := bytes.NewReader(data)

	err := s.WriteStream(key, sourceStream)
	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}
}

func TestStoreRead(t *testing.T) {
	opts := StoreOpts{
		Root:              "my_test_network",
		PathTransformFunc: CASPathTransfromFunc,
	}
	s := NewStore(opts)

	key := "test_read_file.txt"
	data := []byte("We need to make sure we can read this back from the disk!")

	// 1. Write the file
	err := s.WriteStream(key, bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to write file: %s", err)
	}

	// 2. Read the file back
	r, err := s.ReadStream(key)
	if err != nil {
		t.Fatalf("failed to read file: %s", err)
	}
	defer r.Close()

	// 3. Read the bytes from the stream
	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read bytes from stream: %s", err)
	}

	// 4. Verify it matches what we wrote!
	if string(b) != string(data) {
		t.Errorf("want %s, got %s", string(data), string(b))
	}
}
