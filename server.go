package main

import (
	"fmt"

	"github.com/girish/storage/p2p"
)

type FileServerOpts struct {
	StorageRoot string //StorageRoot is the folder on hard drive where this node will save files
	Transport   p2p.Transport
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
	fmt.Printf("FileServer received %d bytes from %s\n", len(rpc.Payload), rpc.From)

	// TODO: Next, we will parse rpc.Payload to figure out if the client
	// wants to STORE a file or GET a file!
	fmt.Printf("Payload: %s\n", string(rpc.Payload))
}
