package p2p

import (
	"encoding/binary"
	"encoding/gob"
	"io"
)

type Decoder interface {
	Decode(io.Reader, any) error
}

type GOBDecoder struct{}

type LengthPrefixDecoder struct{}

func (doc LengthPrefixDecoder) Decode(r io.Reader, v any) error {
	msg, ok := v.(*Message)
	if !ok {
		return io.EOF
	}

	//1. Read the length Prefix (4 bytes)
	lengthBuf := make([]byte, 4)
	if _, err := io.ReadFull(r, lengthBuf); err != nil {
		return err
	}

	//Convert the 4 bytes into an interger
	msgLength := binary.LittleEndian.Uint32(lengthBuf)

	//2. Read the exact number of bytes specified by the msgLength
	payloadBuf := make([]byte, msgLength)
	if _, err := io.ReadFull(r, payloadBuf); err != nil {
		return err
	}

	msg.Payload = payloadBuf
	return nil
}

func (dec GOBDecoder) Decode(r io.Reader, v any) error {
	return gob.NewDecoder(r).Decode(v)
}

type NOPDecoder struct{}

func (doc NOPDecoder) Decoder(r io.Reader, v any) error {
	buf := make([]byte, 1028)
	_, err := r.Read(buf)
	if err != nil {
		return err
	}
	return nil
}
