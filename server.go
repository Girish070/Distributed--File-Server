package main

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/girish/storage/p2p"
	"github.com/girish/storage/store"
)

type FileServerOpts struct {
	StorageRoot string //StorageRoot is the folder on hard drive where this node will save files
	Transport   p2p.Transport
	Store       *store.Store
}

// DataMessage is wire protocol payload!
type DataMessage struct {
	Key  string
	Data []byte
}

type FileServer struct {
	FileServerOpts
	quitCh chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	return &FileServer{
		FileServerOpts: opts,
		quitCh:         make(chan struct{}),
	}
}

// Start boots up the transport layer and begins processing messages
func (s *FileServer) Start() error {
	fmt.Printf("Starting File Server... (Storage Root: %s)\n", s.StorageRoot)

	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	s.loop()
	return nil
}

func (s *FileServer) loop() {
	for {
		select {
		case rpc := <-s.Transport.Consume():
			s.handleMessage(rpc)
		case <-s.quitCh:
			fmt.Println("File Server shutting down")
			return
		}
	}
}

func (s *FileServer) handleMessage(rpc p2p.RPC) {
	//1 Decode the binary network payload back into Structured DataMessage
	var msg DataMessage
	if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg); err != nil {
		fmt.Printf("Failed to decode network payload: %s\n", err)
		return
	}
	fmt.Printf("FileServer received command to store file: '%s' (%d bytes)\n", msg.Key, len(msg.Data))

	//2 Write the file to disk using CAS store engine
	err := s.Store.WriteStream(msg.Key, bytes.NewReader(msg.Data))
	if err != nil {
		fmt.Printf("Error Storing file to disk: %s\n", err)
		return
	}
	fmt.Println("File successfully saved from the network!")
}
