package store

import (
	"bytes"
	"testing"
)

// ... (keep your existing TestPathTransformFunc) ...

func TestStore(t *testing.T) {
	// 1. Initialize our Store
	opts := StoreOpts{
		Root:              "my_test_network",
		PathTransformFunc: CASPathTransfromFunc,
	}
	s := NewStore(opts)

	// 2. Pretend this is a file a user is uploading
	key := "my_special_picture.png"
	data := []byte("These are some fake image bytes. Pretend this is a cool PNG!")
	
	// Create an io.Reader out of our bytes
	sourceStream := bytes.NewReader(data)

	// 3. Write it to disk!
	err := s.WriteStream(key, sourceStream)
	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}
}