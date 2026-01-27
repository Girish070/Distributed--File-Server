package p2p

import (
	"encoding/gob"
	"io"
)

type Decoder interface {
	Decode(io.Reader, any) error
}

type GOBDecoder struct{}

func (dec GOBDecoder) Decode(r io.Reader, v any) error {
	return gob.NewDecoder(r).Decode(v)
}

type NOPDecoder struct{}

func (doc NOPDecoder) Decoder(r io.Reader, v any) error {
	buf := make([]byte, 1028)
	m, err := r.Read(buf)
	if err != nil {
		return err
	}
	return nil
}
